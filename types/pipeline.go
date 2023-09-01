package types

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

	fluentbitconfig "github.com/calyptia/go-fluentbit-config/v2"
)

const DefaultHealthCheckPipelinePort = 2020

// Pipeline model.
type Pipeline struct {
	ID                           string           `json:"id" yaml:"id"`
	Name                         string           `json:"name" yaml:"name"`
	Kind                         PipelineKind     `json:"kind" yaml:"kind"`
	Config                       PipelineConfig   `json:"config" yaml:"config"`
	ConfigSections               []ConfigSection  `json:"configSections" yaml:"configSections"`
	Image                        *string          `json:"image" yaml:"image"`
	Status                       PipelineStatus   `json:"status" yaml:"status"`
	ResourceProfile              ResourceProfile  `json:"resourceProfile" yaml:"resourceProfile"`
	TracingEnabled               bool             `json:"tracingEnabled" yaml:"tracingEnabled"`
	WaitForChecksBeforeDeploying bool             `json:"waitForChecksBeforeDeploying" yaml:"waitForChecksBeforeDeploying"`
	ReplicasCount                uint             `json:"replicasCount" yaml:"replicasCount"`
	Tags                         []string         `json:"tags" yaml:"tags"`
	Metadata                     *json.RawMessage `json:"metadata" yaml:"metadata"`
	ChecksTotal                  uint             `json:"checksTotal" yaml:"checksTotal"`
	ChecksOK                     uint             `json:"checksOK" yaml:"checksOK"`
	ChecksRunning                uint             `json:"checksRunning" yaml:"checksRunning"`
	CreatedAt                    time.Time        `json:"createdAt" yaml:"createdAt"`
	UpdatedAt                    time.Time        `json:"updatedAt" yaml:"updatedAt"`
	Ports                        []PipelinePort   `json:"ports,omitempty" yaml:"ports,omitempty"`
	Files                        []PipelineFile   `json:"files,omitempty" yaml:"files,omitempty"`
	Secrets                      []PipelineSecret `json:"secrets,omitempty" yaml:"secrets,omitempty"`
}

func (p *Pipeline) ApplyConfigSections() error {
	if len(p.ConfigSections) == 0 {
		return nil
	}

	format := fluentbitconfig.Format(p.Config.ConfigFormat)
	c, err := fluentbitconfig.ParseAs(p.Config.RawConfig, format)
	if err != nil {
		return err
	}

	for _, section := range p.ConfigSections {
		c.AddSection(fluentbitconfig.SectionKind(section.Kind), section.Properties.AsProperties())
	}

	raw, err := c.DumpAs(format)
	if err != nil {
		return err
	}

	p.Config.RawConfig = raw
	p.Status.Config.RawConfig = raw

	return nil
}

type PipelineKind string

const (
	PipelineKindDaemonSet  PipelineKind = "daemonSet"
	PipelineKindDeployment PipelineKind = "deployment"
)

// AllPipelineKindTypes all valid pipeline kinds.
var AllPipelineKindTypes = [...]PipelineKind{
	PipelineKindDaemonSet,
	PipelineKindDeployment,
}

// CreatePipeline request payload for creating a new pipeline.
type CreatePipeline struct {
	Name                 string                 `json:"name"`
	Kind                 PipelineKind           `json:"kind"`
	ReplicasCount        uint                   `json:"replicasCount"`
	RawConfig            string                 `json:"rawConfig"`
	ConfigFormat         ConfigFormat           `json:"configFormat"`
	Secrets              []CreatePipelineSecret `json:"secrets"`
	Files                []CreatePipelineFile   `json:"files"`
	ResourceProfileName  string                 `json:"resourceProfile"`
	Image                *string                `json:"image"`
	SkipConfigValidation bool                   `json:"skipConfigValidation"`
	Metadata             *json.RawMessage       `json:"metadata"`
	Tags                 []string               `json:"tags"`

	// Deprecated: in favor of NoAutoCreateEndpointsFromConfig
	AutoCreatePortsFromConfig bool `json:"autoCreatePortsFromConfig"`
	// Deprecated: in favor of NoAutoCreateChecksFromConfig
	AutoCreateChecksFromConfig bool `json:"autoCreateChecksFromConfig"`

	// no automatically create endpoints from config
	NoAutoCreateEndpointsFromConfig bool `json:"noAutoCreateEndpointsFromConfig"`
	// no automatically create checks based on the output configuration.
	NoAutoCreateChecksFromConfig bool `json:"noAutoCreateChecksFromConfig"`

	// WaitForChecksBeforeDeploying is a conditional variable that defines behavior on the
	// pipeline deployment
	//
	// If set to true:
	//
	// If all checks associated with the pipeline run successfully, the status of
	// the pipeline will be switched to CHECKS_OK and the deployment will be executed.
	//
	// If any of the checks associated with the pipeline fails, the status of
	// the pipeline will be switched to CHECKS_FAILED and the deployment of the pipeline
	// will be blocked.
	//
	// If set to false (default):
	//
	// If all checks associated with the pipeline run successfully, the status of
	// the pipeline will be switched to CHECKS_OK and the deployment will be executed.
	//
	// If any of the checks associated with the pipeline fails, the status of
	// the pipeline will be switched to CHECKS_FAILED and the deployment of the pipeline
	// will be executed.
	WaitForChecksBeforeDeploying bool `json:"waitForChecksBeforeDeploying"`

	status PipelineStatusKind
	// Internal denotes that this pipeline was created by the system.
	// That's the case for the "health-check-*" automated pipeline
	// with each new core instance.
	// We use it to not take these into account for project quotas.
	internal bool
	// ClusterLogging denotes that this pipeline is internal.
	// There should be only one cluster-logging pipeline in the system.
	clusterLogging bool
}

func (in *CreatePipeline) SetStatus(status PipelineStatusKind) {
	in.status = status
}

func (in *CreatePipeline) SetInternal(internal bool) {
	in.internal = internal
}

func (in *CreatePipeline) SetClusterLogging(clusterLogging bool) {
	in.clusterLogging = clusterLogging
}

func (in CreatePipeline) Status() PipelineStatusKind {
	return in.status
}

func (in CreatePipeline) Internal() bool {
	return in.internal
}

func (in CreatePipeline) ClusterLogging() bool {
	return in.clusterLogging
}

// CreatedPipeline response payload after creating a pipeline successfully.
type CreatedPipeline struct {
	ID                           string           `json:"id"`
	Name                         string           `json:"name"`
	Kind                         PipelineKind     `json:"kind"`
	Config                       PipelineConfig   `json:"config"`
	Secrets                      []PipelineSecret `json:"secrets"`
	Files                        []PipelineFile   `json:"files"`
	Status                       PipelineStatus   `json:"status"`
	ResourceProfile              ResourceProfile  `json:"resourceProfile"`
	Checks                       []PipelineCheck  `json:"checks"`
	ReplicasCount                uint             `json:"replicasCount"`
	WaitForChecksBeforeDeploying bool             `json:"waitForChecksBeforeDeploying"`
	CreatedAt                    time.Time        `json:"createdAt"`
}

// PipelinesParams represents the request payload for querying pipelines.
type PipelinesParams struct {
	ProjectID      *string
	CoreInstanceID *string

	Last                     *uint
	Before                   *string
	Kind                     *PipelineKind
	Name                     *string
	TagsQuery                *string
	ConfigFormat             *ConfigFormat
	IncludeObjects           *PipelineObjectsParams
	RenderWithConfigSections bool
}

func (p PipelinesParams) Tags() *[]string {
	if p.TagsQuery == nil {
		return nil
	}

	tags := strings.Split(*p.TagsQuery, " AND ")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}
	return &tags
}

// PipelineObjectsParams represents the options for including different types of pipeline objects in the response.
type PipelineObjectsParams struct {
	Files   bool // include files in the response.
	Secrets bool // include secrets in the response.
	Ports   bool // include ports in the response.
}

// NewPipelineObjectsParams creates and returns a new PipelineObjectsParams object based on the "include" parameter in the given request.
func NewPipelineObjectsParams(r *http.Request) *PipelineObjectsParams {
	include := strings.Split(r.URL.Query().Get("include"), ",")
	return &PipelineObjectsParams{
		Files:   slices.Contains(include, "files"),
		Secrets: slices.Contains(include, "secrets"),
		Ports:   slices.Contains(include, "ports"),
	}
}

// Pipelines paginated list.
type Pipelines struct {
	Items     []Pipeline `json:"items"`
	EndCursor *string    `json:"endCursor"`
	Count     int        `json:"count"`
}

// UpdatePipeline request payload for updating a pipeline.
type UpdatePipeline struct {
	Name            *string             `json:"name"`
	Kind            *PipelineKind       `json:"kind"`
	Status          *PipelineStatusKind `json:"status"`
	ConfigFormat    *ConfigFormat       `json:"configFormat"`
	ReplicasCount   *uint               `json:"replicasCount"`
	RawConfig       *string             `json:"rawConfig"`
	ResourceProfile *string             `json:"resourceProfile"`
	Image           *string             `json:"image"`
	// Deprecated: in favor of NoAutoCreateEndpointsFromConfig
	AutoCreatePortsFromConfig *bool `json:"autoCreatePortsFromConfig"`
	// Deprecated: in favor of NoAutoCreateChecksFromConfig
	AutoCreateChecksFromConfig *bool `json:"autoCreateChecksFromConfig"`

	// no automatically create endpoints from config
	NoAutoCreateEndpointsFromConfig bool `json:"noAutoCreateEndpointsFromConfig"`
	// no automatically create checks based on the output configuration.
	NoAutoCreateChecksFromConfig bool `json:"noAutoCreateChecksFromConfig"`

	// this defines behavior; await for checks to complete before reporting the status back.
	WaitForChecksBeforeDeploying *bool `json:"waitForChecksBeforeDeploying"`

	SkipConfigValidation bool                   `json:"skipConfigValidation"`
	Metadata             *json.RawMessage       `json:"metadata"`
	Secrets              []UpdatePipelineSecret `json:"secrets"`
	Files                []UpdatePipelineFile   `json:"files"`
	Events               []PipelineEvent        `json:"events"`

	clusterLogging    *bool
	resourceProfileID *string
}

func (in *UpdatePipeline) SetClusterLogging(clusterLogging bool) {
	in.clusterLogging = &clusterLogging
}

func (in *UpdatePipeline) SetResourceProfileID(resourceProfileID string) {
	in.resourceProfileID = &resourceProfileID
}

func (in UpdatePipeline) ClusterLogging() *bool {
	return in.clusterLogging
}

func (in UpdatePipeline) ResourceProfileID() *string {
	return in.resourceProfileID
}

// PipelineParams request payload for querying a single pipeline.
type PipelineParams struct {
	ConfigFormat             *ConfigFormat
	RenderWithConfigSections bool
}

// UpdatedPipeline response payload after updating a pipeline successfully.
type UpdatedPipeline struct {
	AddedPorts   []PipelinePort `json:"addedPorts"`
	RemovedPorts []PipelinePort `json:"removedPorts"`

	// Pipeline checks that have been created/updated based on AutoCreatePreChecksFromConfig changes.
	AddedChecks   []PipelineCheck `json:"addedChecks"`
	RemovedChecks []PipelineCheck `json:"removedChecks"`
}

// UpdatePipelineClusterObjects update cluster objects associated to a pipeline.
type UpdatePipelineClusterObjects struct {
	ClusterObjectsIDs []string `json:"clusterObjectsIDs"`
}

// PipelineClusterObjectsParams request payload to filter cluster objects belonging to a pipeline.
type PipelineClusterObjectsParams struct {
	Last   *uint
	Before *string
}

// PipelineMetadata is the default metadata format for a pipeline.
type PipelineMetadata map[string]any

// UpdatePipelineMetadata request payload to store a key on the metadata field with the given value (json serializable).
type UpdatePipelineMetadata struct {
	Key   *string          `json:"key"`
	Value *json.RawMessage `json:"value"`
}

// PipelineMetadataParams request payload for listing metadata from keys.
type PipelineMetadataParams struct {
	Keys *[]string
}
