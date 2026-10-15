package handler

import (
"fmt"
"net/http"
"time"

"github.com/gin-gonic/gin"
"github.com/onioneye/correlation-engine/internal/engine"
"github.com/onioneye/correlation-engine/internal/negcontrol"
)

// EvaluateActorRequest is the payload to evaluate an actor's composite score across multi-modal vectors.
type EvaluateActorRequest struct {
ActorID   string                     `json:"actor_id"`
NCPassed  int                        `json:"nc_passed"`
NCFlagged int                        `json:"nc_flagged"`
Evidence  []EvidenceInput            `json:"evidence,omitempty"`
Artifacts []engine.ExtractedArtifact `json:"artifacts,omitempty"`
}

type EvidenceInput struct {
ImageSrc        string    `json:"image_src"`
ClassifierScore float64   `json:"classifier_score"`
IsHashMatch     bool      `json:"is_hash_match"`
FlaggedAt       time.Time `json:"flagged_at"`
}

// EvaluateJobRequest evaluates a single job for Negative Control state.
type EvaluateJobRequest struct {
JobID        string             `json:"job_id"`
Timeouts     int                `json:"timeouts"`
CircuitBroke bool               `json:"circuit_broke"`
Images       []ImageResultInput `json:"images"`
}

type ImageResultInput struct {
ImageSrc       string `json:"image_src"`
HashMatch      bool   `json:"hash_match"`
ClassifierHit  bool   `json:"classifier_hit"`
ShortCircuited bool   `json:"short_circuited"`
}

// EvaluateStylometryRequest takes two vectors and their word counts.
type EvaluateStylometryRequest struct {
VectorA      []float64 `json:"vector_a"`
VectorB      []float64 `json:"vector_b"`
SampleWordsA int       `json:"sample_words_a"`
SampleWordsB int       `json:"sample_words_b"`
}

func RegisterRoutes(r gin.IRouter) {
r.POST("/evaluate/actor", handleEvaluateActor)
r.POST("/evaluate/job", handleEvaluateJob)
r.POST("/evaluate/stylometry", handleEvaluateStylometry)
r.POST("/evaluate/vector", handleEvaluateVector)
}

func handleEvaluateActor(c *gin.Context) {
var req EvaluateActorRequest
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid JSON: %v", err)})
return
}

profile := engine.ActorProfile{
ActorID:   req.ActorID,
NCPassed:  req.NCPassed,
NCFlagged: req.NCFlagged,
Evidence:  make([]engine.Evidence, len(req.Evidence)),
Artifacts: req.Artifacts,
}

for i, ev := range req.Evidence {
profile.Evidence[i] = engine.Evidence{
ImageSrc:        ev.ImageSrc,
ClassifierScore: ev.ClassifierScore,
IsHashMatch:     ev.IsHashMatch,
FlaggedAt:       ev.FlaggedAt,
}
}

result := engine.EvaluateProfile(profile, engine.DefaultConfig, time.Now())

c.JSON(http.StatusOK, gin.H{
"actor_id":         result.ActorID,
"score":            result.CompositeScore,
"alert":            result.IsAlert,
"threshold":        result.Threshold,
"nc_dampener":      result.NCDampener,
"scored_artifacts": result.ScoredArtifacts,
})
}

func handleEvaluateJob(c *gin.Context) {
var req EvaluateJobRequest
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid JSON: %v", err)})
return
}

outcome := negcontrol.JobOutcome{
JobID:        req.JobID,
Timeouts:     req.Timeouts,
CircuitBroke: req.CircuitBroke,
Images:       make([]negcontrol.ImageResult, len(req.Images)),
}

for i, img := range req.Images {
outcome.Images[i] = negcontrol.ImageResult{
ImageSrc:       img.ImageSrc,
HashMatch:      img.HashMatch,
ClassifierHit:  img.ClassifierHit,
ShortCircuited: img.ShortCircuited,
}
}

state := negcontrol.EvaluateJob(outcome)

c.JSON(http.StatusOK, gin.H{
"job_id": outcome.JobID,
"state":  state,
})
}

func handleEvaluateStylometry(c *gin.Context) {
var req EvaluateStylometryRequest
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid JSON: %v", err)})
return
}

result := engine.EvaluateStylometry(req.VectorA, req.VectorB, req.SampleWordsA, req.SampleWordsB)
c.JSON(http.StatusOK, result)
}

func handleEvaluateVector(c *gin.Context) {
var art engine.ExtractedArtifact
if err := c.ShouldBindJSON(&art); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid JSON: %v", err)})
return
}

scored := engine.ScoreArtifact(art, time.Now(), engine.DefaultConfig.HalfLifeDays)
c.JSON(http.StatusOK, scored)
}
