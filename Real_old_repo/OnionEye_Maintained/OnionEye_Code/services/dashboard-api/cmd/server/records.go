package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const maxStoredRecords = 50

// graphRecord is one saved OnionEye graph snapshot (JSON Record format).
type graphRecord struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	CreatedAt string         `json:"created_at"`
	NodeCount int            `json:"node_count"`
	LinkCount int            `json:"link_count"`
	Graph     map[string]any `json:"graph,omitempty"`
}

type recordStore struct {
	mu    sync.Mutex
	items []graphRecord
}

func newRecordStore() *recordStore {
	return &recordStore{items: make([]graphRecord, 0, maxStoredRecords)}
}

func (s *recordStore) list() []graphRecord {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]graphRecord, len(s.items))
	for i, r := range s.items {
		// Strip the payload from the list view — only metadata.
		out[i] = graphRecord{
			ID: r.ID, Name: r.Name, CreatedAt: r.CreatedAt,
			NodeCount: r.NodeCount, LinkCount: r.LinkCount,
		}
	}
	return out
}

func (s *recordStore) get(id string) (graphRecord, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, r := range s.items {
		if r.ID == id {
			return r, true
		}
	}
	return graphRecord{}, false
}

func (s *recordStore) push(r graphRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append([]graphRecord{r}, s.items...)
	if len(s.items) > maxStoredRecords {
		s.items = s.items[:maxStoredRecords]
	}
}

// recordsListHandler GET /api/v1/dashboard/records
func recordsListHandler(store *recordStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"records": store.list()})
	}
}

// recordsGetHandler GET /api/v1/dashboard/records/:id
func recordsGetHandler(store *recordStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		r, ok := store.get(c.Param("id"))
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
			return
		}
		c.JSON(http.StatusOK, r)
	}
}

type saveRecordRequest struct {
	Name  string         `json:"name"`
	Graph map[string]any `json:"graph"`
}

// recordsSaveHandler POST /api/v1/dashboard/records
func recordsSaveHandler(store *recordStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req saveRecordRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.Graph == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload — need {name, graph}"})
			return
		}
		if req.Name == "" {
			req.Name = "untitled"
		}
		nodeCount, linkCount := 0, 0
		if nodes, ok := req.Graph["nodes"].([]any); ok {
			nodeCount = len(nodes)
		}
		if links, ok := req.Graph["links"].([]any); ok {
			linkCount = len(links)
		}
		rec := graphRecord{
			ID:        time.Now().UTC().Format("20060102T150405.000000000"),
			Name:      req.Name,
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
			NodeCount: nodeCount,
			LinkCount: linkCount,
			Graph:     req.Graph,
		}
		store.push(rec)
		c.JSON(http.StatusOK, gin.H{"status": "ok", "id": rec.ID, "name": rec.Name})
	}
}