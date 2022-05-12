package types

import "time"

const (
	DefaultAggregatorVersion = "v0.1.12"
)

// Aggregator model.
type Aggregator struct {
	ID              string    `json:"id" yaml:"id"`
	Token           string    `json:"token" yaml:"token"`
	Name            string    `json:"name" yaml:"name"`
	EnvironmentName string    `json:"environmentName" yaml:"environmentName"`
	Version         string    `json:"version" yaml:"version"`
	PipelinesCount  uint64    `json:"pipelinesCount" yaml:"pipelinesCount"`
	Tags            []string  `json:"tags" yaml:"tags"`
	CreatedAt       time.Time `json:"createdAt" yaml:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt" yaml:"updatedAt"`
}

// Aggregators paginated list.
type Aggregators struct {
	Items     []Aggregator
	EndCursor *string
}

// CreateAggregator request payload for creating a new aggregator.
type CreateAggregator struct {
	Name                    string   `json:"name"`
	Version                 string   `json:"version" yaml:"version"`
	AddHealthCheckPipeline  bool     `json:"addHealthCheckPipeline"`
	HealthCheckPipelinePort uint     `json:"healthCheckPipelinePort"`
	Tags                    []string `json:"tags" yaml:"tags"`
}

// CreatedAggregator response payload after creating an aggregator successfully.
type CreatedAggregator struct {
	ID            string    `json:"id"`
	Token         string    `json:"token"`
	PrivateRSAKey []byte    `json:"privateRSAKey"`
	PublicRSAKey  []byte    `json:"publicRSAKey"`
	Name          string    `json:"name"`
	Version       string    `json:"version"`
	CreatedAt     time.Time `json:"createdAt"`
	Tags          []string  `json:"tags"`

	HealthCheckPipeline *Pipeline         `json:"healthCheckPipeline"`
	ResourceProfiles    []ResourceProfile `json:"resourceProfiles"`
}

// AggregatorsParams request payload for querying aggregators.
type AggregatorsParams struct {
	Last   *uint64
	Before *string
	Name   *string
	Tags   *string
}

// UpdateAggregator request payload for updating an aggregator.
type UpdateAggregator struct {
	Name    *string `json:"name"`
	Version *string `json:"version"`
}
