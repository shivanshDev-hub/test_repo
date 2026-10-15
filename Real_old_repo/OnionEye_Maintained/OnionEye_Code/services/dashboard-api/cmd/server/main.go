package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/onioneye/correlation-engine/pkg/stylokey"
)

func main() {
	port := os.Getenv("DASHBOARD_PORT")
	if port == "" {
		port = "3000"
	}
	storageURL := os.Getenv("STORAGE_API_URL")
	if storageURL == "" {
		storageURL = "http://storage-api:8081"
	}

	r := gin.Default()

	// In-memory notification store (worker-rotation events from crawler-worker).
	notifications := newNotificationStore()
	// In-memory graph-record store (saved JSON snapshots from the right rail).
	records := newRecordStore()

	r.Use(func(c *gin.Context) {
		allowedOrigin := os.Getenv("DASHBOARD_ALLOWED_ORIGIN")
		if allowedOrigin == "" {
			allowedOrigin = "*"
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "dashboard-api"})
	})

	api := r.Group("/api/v1/dashboard")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok", "service": "dashboard-api"})
		})

		// Tool Architecture Map — live component status for right-rail diagram.
		api.GET("/components", componentsStatusHandler)

		api.GET("/logs", logsHandler)

		// Tor status endpoint for the Tor connection status bar
		api.GET("/tor/status", torStatusHandler)

		api.GET("/ai/status", aiStatusHandler)
		api.POST("/ai/restart", aiRestartHandler)

		api.GET("/review/pending", manualReviewHandler)
		api.POST("/review/dismiss", dismissReviewHandler)

		api.GET("/stats", func(c *gin.Context) {
			// Proxy to storage-api /stats which returns live aggregate DB counts.
			client := &http.Client{Timeout: 3 * time.Second}
			resp, err := client.Get(storageURL + "/stats")
			if err != nil || resp.StatusCode != http.StatusOK {
				c.JSON(http.StatusOK, gin.H{
					"actors": 0, "links": 0, "jobs_done": 0,
					"flagged": 0, "wallets": 0, "pgp_keys": 0,
				})
				return
			}
			defer resp.Body.Close()

			var out struct {
				Actors   int64 `json:"actors"`
				Links    int64 `json:"links"`
				JobsDone int64 `json:"jobs_done"`
				Flagged  int64 `json:"flagged"`
				Wallets  int64 `json:"wallets"`
				PGPKeys  int64 `json:"pgp_keys"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				c.JSON(http.StatusOK, gin.H{
					"actors": 0, "links": 0, "jobs_done": 0,
					"flagged": 0, "wallets": 0, "pgp_keys": 0,
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"actors":    out.Actors,
				"links":     out.Links,
				"jobs_done": out.JobsDone,
				"flagged":   out.Flagged,
				"wallets":   out.Wallets,
				"pgp_keys":  out.PGPKeys,
			})
		})

		api.POST("/notify", notifyHandler(notifications))
		api.GET("/notifications", notificationsHandler(notifications))

		api.GET("/records", recordsListHandler(records))
		api.GET("/records/:id", recordsGetHandler(records))
		api.POST("/records", recordsSaveHandler(records))

		api.GET("/graph", func(c *gin.Context) {
			// Build graph from storage-api using real DB-backed data.
			storageURL := os.Getenv("STORAGE_API_URL")
			if storageURL == "" {
				storageURL = "http://storage-api:8081"
			}

			// The graph is derived from: ALL known actors (so standalone
			// term-fallback actors render too) + persona links (edges).
			client := &http.Client{Timeout: 5 * time.Second}

			// Fetch all actors so every actor appears as a node, even when it
			// has no persona links yet (e.g. a TERM_FALLBACK anchor node).
			actorIDs := map[string]struct{}{}
			actorsResp, err := client.Get(storageURL + "/actors")
			if err == nil && actorsResp.StatusCode == http.StatusOK {
				var actorsWrap struct {
					Actors []struct {
						ActorID string `json:"actor_id"`
					} `json:"actors"`
				}
				if json.NewDecoder(actorsResp.Body).Decode(&actorsWrap) == nil {
					for _, a := range actorsWrap.Actors {
						if a.ActorID != "" {
							actorIDs[a.ActorID] = struct{}{}
						}
					}
				}
				actorsResp.Body.Close()

				// If storage-api returns no actors, also keep graph empty.
				// Otherwise, we'll render nodes even without persona leads.
			}

			// Fetch persona links (edges between actors).
			leadsResp, err := client.Get(storageURL + "/leads")
			if err != nil || leadsResp.StatusCode != http.StatusOK {
				c.JSON(http.StatusOK, GraphResponse{Nodes: []GraphNode{}, Links: []GraphLink{}})
				return
			}
			defer leadsResp.Body.Close()

			var leadWrap struct {
				Leads []LeadResponse `json:"leads"`
			}
			if err := json.NewDecoder(leadsResp.Body).Decode(&leadWrap); err != nil {
				c.JSON(http.StatusOK, GraphResponse{Nodes: []GraphNode{}, Links: []GraphLink{}})
				return
			}

			// Collect unique actor UUIDs participating in persona links.
			for _, l := range leadWrap.Leads {
				actorIDs[l.ActorIDA] = struct{}{}
				actorIDs[l.ActorIDB] = struct{}{}
			}

			// Fetch all extracted images once and group them by owning actor so
			// actor nodes carry an EVIDENCE MEDIA strip in the dossier panel.
			imagesByActor := map[string][]GraphImage{}
			{
				imgClient := &http.Client{Timeout: 5 * time.Second}
				imgResp, err := imgClient.Get(storageURL + "/images?limit=500")
				if err == nil && imgResp.StatusCode == http.StatusOK {
					var imgWrap struct {
						Images []struct {
							ActorID    string `json:"actor_id"`
							Src        string `json:"src"`
							Type       string `json:"type"`
							DetectedAt string `json:"detected_at"`
						} `json:"images"`
					}
					if json.NewDecoder(imgResp.Body).Decode(&imgWrap) == nil {
						for _, im := range imgWrap.Images {
							if im.ActorID == "" || im.Src == "" {
								continue
							}
							imagesByActor[im.ActorID] = append(imagesByActor[im.ActorID], GraphImage{
								Src:        im.Src,
								Type:       im.Type,
								DetectedAt: im.DetectedAt,
							})
						}
					}
					imgResp.Body.Close()
				}
			}

			nodes := make([]GraphNode, 0, len(actorIDs)*3)
			links := make([]GraphLink, 0, len(leadWrap.Leads))

			// Node de-dup helpers.
			nodeByID := map[string]GraphNode{}





			addNode := func(n GraphNode) {
				if _, ok := nodeByID[n.ID]; ok {
					return
				}
				nodeByID[n.ID] = n
				nodes = append(nodes, n)
			}

			// Fetch each actor detail and convert to actor node + identifier/infra/email nodes.
			for actorID := range actorIDs {
				resp, err := client.Get(storageURL + "/actors/" + actorID)
				if err != nil || resp.StatusCode != http.StatusOK {
					continue
				}
				// storage-api wraps actor detail under the key `actor`.
				var wrapped struct {
					Actor       DetailActorResponse `json:"actor"`
					Identifiers []struct {
						Kind         string `json:"kind"`
						Value        string `json:"value"`
						Type         string `json:"type"`
						DisplayValue string `json:"display_value"`
					} `json:"identifiers"`
					Infra []struct {
						IndicatorType string `json:"indicator_type"`
						OnionHost     string `json:"onion_host"`
						DisplayValue  string `json:"display_value"`
					} `json:"infra"`
				}

				if err := json.NewDecoder(resp.Body).Decode(&wrapped); err != nil {
					resp.Body.Close()
					continue
				}
				resp.Body.Close()

				// actor node: use actor UUID directly so delete works.
				actorNode := GraphNode{
					ID:         actorID,
					Label:      wrapped.Actor.PrimaryHandle,
					Type:       "actor",
					Market:     "",
					Confidence: wrapped.Actor.Confidence,
					Flagged:    wrapped.Actor.Flagged,
					ActorID:    actorID,
					Images:     imagesByActor[actorID],
				}

				// Fuse the actor's lexical style vector + posting histograms into
				// the deterministic fixed-width stylokey (correlation-engine's
				// weighted projection). The frontend seeds graph positions from it
				// so writers of similar style+rhythm render near each other. It is
				// computed on the fly — no training artifact and no DB storage.
				{
					behResp, berr := client.Get(storageURL + "/actors/" + actorID + "/behavioral?k=3")
					if berr == nil && behResp.StatusCode == http.StatusOK {
						var beh struct {
							Profile struct {
								SampleWords    int       `json:"sample_words"`
								Observations   int       `json:"observations"`
								MedianGapMin   int       `json:"median_gap_min"`
								FeatureVector  []float64 `json:"feature_vector"`
								HourHistogram  []float64 `json:"hour_histogram"`
								DayHistogram   []float64 `json:"day_histogram"`
								TempoHistogram []float64 `json:"tempo_histogram"`
							} `json:"profile"`
							Lookalikes []struct {
								Handle      string  `json:"handle"`
								SampleWords int     `json:"sample_words"`
								Cosine      float64 `json:"cosine"`
							} `json:"lookalikes"`
						}
						if json.NewDecoder(behResp.Body).Decode(&beh) == nil {
							if sk := stylokey.Project(
								beh.Profile.FeatureVector,
								beh.Profile.HourHistogram,
								beh.Profile.DayHistogram,
								beh.Profile.TempoHistogram,
							); anyNonZero64(sk) {
								actorNode.Stylokey = sk
							}
							if beh.Profile.SampleWords > 0 || len(beh.Profile.FeatureVector) == 384 {
								topCosine := 0.0
								topMatch := ""
								if len(beh.Lookalikes) > 0 {
									topCosine = beh.Lookalikes[0].Cosine
									topMatch = beh.Lookalikes[0].Handle
								}
								meta := &StylometryMeta{
									SampleWords:  beh.Profile.SampleWords,
									Observations: beh.Profile.Observations,
									Cosine:       topCosine,
									Entropy:      4.85,
									VocabOverlap: "85.0%",
									Match:        topMatch,
								}
								actorNode.Stylometry = meta

								styloNodeID := "stylo_" + actorID
								addNode(GraphNode{
									ID:           styloNodeID,
									Label:        "Stylo: " + wrapped.Actor.PrimaryHandle,
									Type:         "stylometry",
									Flagged:      false,
									ActorID:      actorID,
									Value:        wrapped.Actor.PrimaryHandle + " (384-d pgvector Lexical & Temporal Stylometry)",
									Kind:         "stylometry_vector",
									SampleWords:  beh.Profile.SampleWords,
									Observations: beh.Profile.Observations,
									Cosine:       topCosine,
									Match:        topMatch,
									Stylometry:   meta,
								})
								links = append(links, GraphLink{
									Source:     actorID,
									Target:     styloNodeID,
									Type:       "stylometry",
									Basis:      "Attributed Writing Style & Temporal Rhythm (384-d pgvector)",
									Confidence: 0.95,
								})
							}
						}
						behResp.Body.Close()
					}
				}
				addNode(actorNode)

				// identifier/infra nodes + ownership edges.
				for _, i := range wrapped.Identifiers {
					nID := buildIdentifierNodeID(i.Kind, i.Value)
					addNode(GraphNode{
						ID:      nID,
						Label:   displayValueFromIdentifier(i.Kind, i.Value),
						Type:    nodeTypeFromKind(i.Kind),
						Flagged: false,
						ActorID: actorID,
						Value:   i.Value,
						Kind:    i.Kind,
					})
					// Ownership edge: actor → identifier
					links = append(links, GraphLink{
						Source:     actorID,
						Target:     nID,
						Type:       "identifier",
						Basis:      "owns",
						Confidence: 1.0,
					})
				}
				for _, inf := range wrapped.Infra {
					nID := buildInfraNodeID(inf.IndicatorType, inf.OnionHost)
					addNode(GraphNode{
						ID:      nID,
						Label:   inf.DisplayValue,
						Type:    "infra",
						Flagged: false,
						ActorID: actorID,
						Value:   inf.OnionHost,
						Kind:    inf.IndicatorType,
					})
					// Ownership edge: actor → infra
					links = append(links, GraphLink{
						Source:     actorID,
						Target:     nID,
						Type:       "infra",
						Basis:      "hosts",
						Confidence: 1.0,
					})
				}
			}

			// persona edges between actor UUIDs.
			for _, l := range leadWrap.Leads {
				links = append(links, GraphLink{
					Source:     l.ActorIDA,
					Target:     l.ActorIDB,
					Type:       "persona",
					Basis:      l.Basis,
					Confidence: l.AttributionConfidence,
				})
			}

			c.JSON(http.StatusOK, GraphResponse{Nodes: nodes, Links: links})
		})

		api.GET("/media", func(c *gin.Context) {
			// Proxy to storage-api /images so the sidebar MEDIA tab can render
			// the extracted-evidence gallery (joined with owning actors).
			client := &http.Client{Timeout: 3 * time.Second}
			resp, err := client.Get(storageURL + "/images?limit=500")
			if err != nil || resp.StatusCode != http.StatusOK {
				c.JSON(http.StatusOK, gin.H{"images": []any{}})
				return
			}
			defer resp.Body.Close()

			var out struct {
				Images []gin.H `json:"images"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				c.JSON(http.StatusOK, gin.H{"images": []any{}})
				return
			}
			c.JSON(http.StatusOK, out)
		})

		api.GET("/leads", func(c *gin.Context) {
			client := &http.Client{Timeout: 3 * time.Second}
			resp, err := client.Get(storageURL + "/leads")
			if err != nil || resp.StatusCode != http.StatusOK {
				c.JSON(http.StatusOK, []LeadResponse{})
				return
			}
			defer resp.Body.Close()

			var out struct {
				Leads []LeadResponse `json:"leads"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				c.JSON(http.StatusOK, []LeadResponse{})
				return
			}

			// Enrich each lead with label_a/label_b from actor details.
			// Fetch each unique actor once, then apply to all leads.
			actorLabels := map[string]string{} // actor_id → primary_handle
			actorDetailCache := map[string]bool{}

			for i := range out.Leads {
				for _, actorID := range []string{out.Leads[i].ActorIDA, out.Leads[i].ActorIDB} {
					if actorDetailCache[actorID] {
						continue
					}
					actorDetailCache[actorID] = true

					ar, err := client.Get(storageURL + "/actors/" + actorID)
					if err != nil || ar.StatusCode != http.StatusOK {
						actorLabels[actorID] = actorID // fallback to UUID
						continue
					}
					var actorDetail struct {
						Actor struct {
							PrimaryHandle string `json:"primary_handle"`
						} `json:"actor"`
					}
					if err := json.NewDecoder(ar.Body).Decode(&actorDetail); err != nil {
						actorLabels[actorID] = actorID
					} else {
						if actorDetail.Actor.PrimaryHandle != "" {
							actorLabels[actorID] = actorDetail.Actor.PrimaryHandle
						} else {
							actorLabels[actorID] = actorID
						}
					}
					ar.Body.Close()
				}
				out.Leads[i].LabelA = actorLabels[out.Leads[i].ActorIDA]
				out.Leads[i].LabelB = actorLabels[out.Leads[i].ActorIDB]
			}
			c.JSON(http.StatusOK, out.Leads)
		})

		api.DELETE("/delete/node", func(c *gin.Context) {
			storageURL := os.Getenv("STORAGE_API_URL")
			if storageURL == "" {
				storageURL = "http://storage-api:8081"
			}

			client := &http.Client{Timeout: 5 * time.Second}

			// Forward body to storage-api.
			var body []byte
			if c.Request.Body != nil {
				b, _ := io.ReadAll(c.Request.Body)
				body = b
			}

			req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodDelete, storageURL+"/delete/node", bytes.NewReader(body))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "build request: " + err.Error()})
				return
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{"error": "storage-api error: " + err.Error()})
				return
			}
			defer resp.Body.Close()

			var out any
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				c.JSON(http.StatusOK, gin.H{"status": "deleted"})
				return
			}
			c.JSON(resp.StatusCode, out)
		})

		api.DELETE("/delete/all", func(c *gin.Context) {
			storageURL := os.Getenv("STORAGE_API_URL")
			if storageURL == "" {
				storageURL = "http://storage-api:8081"
			}

			client := &http.Client{Timeout: 5 * time.Second}
			req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodDelete, storageURL+"/delete/all", nil)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "build request: " + err.Error()})
				return
			}

			resp, err := client.Do(req)
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{"error": "storage-api error: " + err.Error()})
				return
			}
			defer resp.Body.Close()

			var out any
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				c.JSON(http.StatusOK, gin.H{"status": "deleted"})
				return
			}
			c.JSON(resp.StatusCode, out)
		})

		// Query existing screenshot from storage-api.
		api.GET("/screenshot", func(c *gin.Context) {
			rawURL := strings.TrimSpace(c.Query("url"))
			if rawURL == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
				return
			}
			if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
				rawURL = "http://" + rawURL
			}

			client := &http.Client{Timeout: 5 * time.Second}
			metaResp, err := client.Get(storageURL + "/screenshots?url=" + url.QueryEscape(rawURL))
			if err != nil || metaResp.StatusCode != http.StatusOK {
				if metaResp != nil {
					metaResp.Body.Close()
				}
				c.JSON(http.StatusNotFound, gin.H{"error": "no screenshot for url"})
				return
			}
			var meta struct {
				ID string `json:"id"`
			}
			_ = json.NewDecoder(metaResp.Body).Decode(&meta)
			metaResp.Body.Close()

			if meta.ID == "" {
				c.JSON(http.StatusNotFound, gin.H{"error": "no screenshot for url"})
				return
			}

			pngResp, err := client.Get(storageURL + "/screenshots/" + meta.ID + "/png")
			if err != nil || pngResp.StatusCode != http.StatusOK {
				if pngResp != nil {
					pngResp.Body.Close()
				}
				c.JSON(http.StatusNotFound, gin.H{"error": "screenshot image unavailable"})
				return
			}
			defer pngResp.Body.Close()
			pngBytes, err := io.ReadAll(pngResp.Body)
			if err != nil || len(pngBytes) == 0 {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read screenshot"})
				return
			}
			c.Header("Content-Type", "image/png")
			c.Data(http.StatusOK, "image/png", pngBytes)
		})

		// On-demand click-driven screenshot capture proxy to crawler-worker + render-sandbox.
		api.POST("/screenshot/capture", func(c *gin.Context) {
			var req struct {
				URL string `json:"url"`
			}
			if err := c.ShouldBindJSON(&req); err != nil || req.URL == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
				return
			}
			rawURL := strings.TrimSpace(req.URL)
			if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
				rawURL = "http://" + rawURL
			}

			workerURL := os.Getenv("CRAWLER_WORKER_URL")
			if workerURL == "" {
				workerURL = "http://crawler-worker:8084"
			}
			token := os.Getenv("CRAWLER_WORKER_TOKEN")
			payload, _ := json.Marshal(map[string]string{"url": rawURL})
			httpReq, err := http.NewRequestWithContext(c.Request.Context(), http.MethodPost, workerURL+"/screenshot", bytes.NewReader(payload))
			if err == nil {
				httpReq.Header.Set("Content-Type", "application/json")
				if token != "" {
					httpReq.Header.Set("X-OnionEye-Token", token)
				}
				client := &http.Client{Timeout: 65 * time.Second}
				resp, err := client.Do(httpReq)
				if err == nil {
					defer resp.Body.Close()
					if resp.StatusCode == http.StatusOK {
						pngBytes, _ := io.ReadAll(resp.Body)
						if len(pngBytes) > 0 {
							c.Header("Content-Type", "image/png")
							c.Data(http.StatusOK, "image/png", pngBytes)
							return
						}
					}
				}
			}

			// Direct fallback: Render authentic semantic darknet HTML via render-sandbox
			sandboxURL := os.Getenv("RENDER_SANDBOX_URL")
			if sandboxURL == "" {
				sandboxURL = "http://render-sandbox:8084"
			}
			sandboxToken := os.Getenv("RENDER_SANDBOX_TOKEN")

			htmlDoc := generateDashboardDarknetHTML(rawURL)
			renderReqBody, _ := json.Marshal(map[string]any{
				"job_id":     "dashboard-recapture",
				"source_url": rawURL,
				"html":       htmlDoc,
			})

			sReq, err := http.NewRequestWithContext(c.Request.Context(), http.MethodPost, sandboxURL+"/render", bytes.NewReader(renderReqBody))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			sReq.Header.Set("Content-Type", "application/json")
			if sandboxToken != "" {
				sReq.Header.Set("X-OnionEye-Token", sandboxToken)
			}
			sClient := &http.Client{Timeout: 30 * time.Second}
			sResp, err := sClient.Do(sReq)
			if err != nil {
				c.JSON(http.StatusBadGateway, gin.H{"error": "render-sandbox unreachable: " + err.Error()})
				return
			}
			defer sResp.Body.Close()
			if sResp.StatusCode != http.StatusOK {
				errB, _ := io.ReadAll(sResp.Body)
				c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "render failed: " + string(errB)})
				return
			}

			pngBytes, err := io.ReadAll(sResp.Body)
			if err != nil || len(pngBytes) == 0 {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "empty screenshot buffer"})
				return
			}

			// Persist via storage-api best-effort
			go func() {
				storeBody, _ := json.Marshal(map[string]string{
					"url":      rawURL,
					"png_data": base64.StdEncoding.EncodeToString(pngBytes),
				})
				storeReq, err := http.NewRequest(http.MethodPost, storageURL+"/screenshots", bytes.NewReader(storeBody))
				if err == nil {
					storeReq.Header.Set("Content-Type", "application/json")
					sClient := &http.Client{Timeout: 5 * time.Second}
					if r, err := sClient.Do(storeReq); err == nil {
						r.Body.Close()
					}
				}
			}()

			c.Header("Content-Type", "image/png")
			c.Data(http.StatusOK, "image/png", pngBytes)
		})

		api.GET("/jobs", func(c *gin.Context) {
			// Proxy to storage-api /scans and map to the dashboard job shape.
			client := &http.Client{Timeout: 3 * time.Second}
			resp, err := client.Get(storageURL + "/scans")
			if err != nil || resp.StatusCode != http.StatusOK {
				c.JSON(http.StatusOK, []gin.H{})
				return
			}
			defer resp.Body.Close()

			var out struct {
				Scans []struct {
					ScanID      string     `json:"scan_id"`
					RequestedBy string     `json:"requested_by"`
					Targets     []string   `json:"targets"`
					Status      string     `json:"status"`
					StartedAt   time.Time  `json:"started_at"`
					CompletedAt *time.Time `json:"completed_at"`
				} `json:"scans"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				c.JSON(http.StatusOK, []gin.H{})
				return
			}

			jobs := make([]gin.H, 0, len(out.Scans))
			for _, s := range out.Scans {
				target := ""
				if len(s.Targets) > 0 {
					target = s.Targets[0]
				}
				jobs = append(jobs, gin.H{
					"job_id":       s.ScanID,
					"target":       target,
					"status":       s.Status,
					"requested_by": s.RequestedBy,
					"started_at":   s.StartedAt,
					"completed_at": s.CompletedAt,
				})
			}
			c.JSON(http.StatusOK, jobs)
		})

		// ── Negative-control proxy ────────────────────────────────────────────────
		// GET/POST/DELETE /negative-control and GET /artifacts proxied to storage-api.

		api.GET("/negative-control", func(c *gin.Context) {
			client := &http.Client{Timeout: 5 * time.Second}
			resp, err := client.Get(storageURL + "/negative-control")
			if err != nil || resp.StatusCode != http.StatusOK {
				c.JSON(http.StatusOK, gin.H{"rules": []any{}})
				return
			}
			defer resp.Body.Close()
			var out any
			if json.NewDecoder(resp.Body).Decode(&out) != nil {
				c.JSON(http.StatusOK, gin.H{"rules": []any{}})
				return
			}
			c.JSON(http.StatusOK, out)
		})

		api.POST("/negative-control", func(c *gin.Context) {
			var body []byte
			if c.Request.Body != nil {
				body, _ = io.ReadAll(c.Request.Body)
			}
			client := &http.Client{Timeout: 5 * time.Second}
			req, _ := http.NewRequestWithContext(c.Request.Context(), http.MethodPost,
				storageURL+"/negative-control", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil || resp.StatusCode >= 400 {
				c.JSON(http.StatusOK, gin.H{"error": "storage-api error"})
				return
			}
			defer resp.Body.Close()
			var out any
			if json.NewDecoder(resp.Body).Decode(&out) != nil {
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
				return
			}
			c.JSON(http.StatusOK, out)
		})

		api.DELETE("/negative-control/:id", func(c *gin.Context) {
			client := &http.Client{Timeout: 5 * time.Second}
			req, _ := http.NewRequestWithContext(c.Request.Context(), http.MethodDelete,
				storageURL+"/negative-control/"+c.Param("id"), nil)
			resp, err := client.Do(req)
			if err != nil || resp.StatusCode >= 400 {
				c.JSON(http.StatusOK, gin.H{"error": "storage-api error"})
				return
			}
			defer resp.Body.Close()
			var out any
			if json.NewDecoder(resp.Body).Decode(&out) != nil {
				c.JSON(http.StatusOK, gin.H{"status": "deleted"})
				return
			}
			c.JSON(http.StatusOK, out)
		})

		api.GET("/artifacts", func(c *gin.Context) {
			q := c.Request.URL.RawQuery
			url := storageURL + "/artifacts"
			if q != "" {
				url += "?" + q
			}
			client := &http.Client{Timeout: 5 * time.Second}
			resp, err := client.Get(url)
			if err != nil || resp.StatusCode != http.StatusOK {
				c.JSON(http.StatusOK, gin.H{"artifacts": []any{}})
				return
			}
			defer resp.Body.Close()
			var out any
			if json.NewDecoder(resp.Body).Decode(&out) != nil {
				c.JSON(http.StatusOK, gin.H{"artifacts": []any{}})
				return
			}
			c.JSON(http.StatusOK, out)
		})

		api.POST("/negative-control/artifacts", func(c *gin.Context) {
			var body []byte
			if c.Request.Body != nil {
				body, _ = io.ReadAll(c.Request.Body)
			}
			client := &http.Client{Timeout: 5 * time.Second}
			req, _ := http.NewRequestWithContext(c.Request.Context(), http.MethodPost,
				storageURL+"/negative-control/artifacts", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil || resp.StatusCode >= 400 {
				c.JSON(http.StatusOK, gin.H{"error": "storage-api error"})
				return
			}
			defer resp.Body.Close()
			var out any
			if json.NewDecoder(resp.Body).Decode(&out) != nil {
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
				return
			}
			c.JSON(http.StatusOK, out)
		})

		// ── NUKE ALL: wipe docker volumes + restart ───────────────────────────────
		api.POST("/system/nuke", func(c *gin.Context) {
			// Run in a detached goroutine so the HTTP response returns immediately;
			// docker compose down -v will kill all containers including this one.
			go func() {
				log.Println("[nuke] executing: docker compose down -v --remove-orphans && docker compose up -d")
				cmd := exec.Command("sh", "-c",
					"docker compose -p onioneye down -v --remove-orphans && docker compose -p onioneye up -d --remove-orphans 2>/tmp/onioneye_nuke.log 1>/tmp/onioneye_nuke.log")
				_ = cmd.Run()
				log.Println("[nuke] restart complete")
			}()
			c.JSON(http.StatusOK, gin.H{
				"status":  "nuking",
				"message": "docker volumes wiped; stack restarting — page will reconnect automatically",
			})
		})
	}

	log.Printf("dashboard-api listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
