package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

// CreateResourceProfile within an aggregator.
// A resource profile is a specification of a resource used during the deployment of a pipeline.
// By default, when you setup an aggregator, Calyptia Cloud will generate 3 resource profiles for you:
// - high-performance-guaranteed-delivery.
// - high-performance-optimal-throughput.
// - best-effort-low-resource.
func (c *Client) CreateResourceProfile(ctx context.Context, aggregatorID string, payload types.CreateResourceProfile) (types.CreatedResourceProfile, error) {
	var out types.CreatedResourceProfile
	return out, c.do(ctx, http.MethodPost, "/v1/aggregators/"+url.PathEscape(aggregatorID)+"/resource_profiles", payload, &out)
}

// ResourceProfiles from an aggregator in descending order.
func (c *Client) ResourceProfiles(ctx context.Context, aggregatorID string, params types.ResourceProfilesParams) ([]types.ResourceProfile, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(*params.Last, defaultUintFormatBase))
	}

	var out []types.ResourceProfile
	path := "/v1/aggregators/" + url.PathEscape(aggregatorID) + "/resource_profiles?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// ResourceProfile by ID.
func (c *Client) ResourceProfile(ctx context.Context, resourceProfileID string) (types.ResourceProfile, error) {
	var out types.ResourceProfile
	return out, c.do(ctx, http.MethodGet, "/v1/resource_profiles/"+url.PathEscape(resourceProfileID), nil, &out)
}

// UpdateResourceProfile by its ID.
func (c *Client) UpdateResourceProfile(ctx context.Context, resourceProfileID string, opts types.UpdateResourceProfile) error {
	return c.do(ctx, http.MethodPatch, "/v1/resource_profiles/"+url.PathEscape(resourceProfileID), opts, nil)
}

// DeleteResourceProfile by its ID.
// The profile cannot be deleted if some pipeline is still referencing it;
// you must delete the pipeline first if you want to delete the profile.
func (c *Client) DeleteResourceProfile(ctx context.Context, resourceProfileID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/resource_profiles/"+url.PathEscape(resourceProfileID), nil, nil)
}
