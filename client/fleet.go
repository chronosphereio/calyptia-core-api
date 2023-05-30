package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
	fluentbitconfig "github.com/calyptia/go-fluentbit-config/v2"
)

func (c *Client) CreateFleet(ctx context.Context, in types.CreateFleet) (types.CreatedFleet, error) {
	var out types.CreatedFleet
	return out, c.do(ctx, http.MethodPost, "/v1/projects/"+url.PathEscape(in.ProjectID)+"/fleets", in, &out)
}

func (c *Client) Fleets(ctx context.Context, params types.FleetsParams) (types.Fleets, error) {
	q := url.Values{}
	if params.Name != nil {
		q.Set("name", *params.Name)
	}
	if params.TagsQuery != nil {
		q.Set("tags_query", *params.TagsQuery)
	}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(uint64(*params.Last), uintBase))
	}
	if params.Before != nil {
		q.Set("before", *params.Before)
	}

	var out types.Fleets
	path := "/v1/projects/" + url.PathEscape(params.ProjectID) + "/fleets?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out.Items, withCursor(&out.EndCursor))
}

func (c *Client) Fleet(ctx context.Context, fleetID string) (types.Fleet, error) {
	var out types.Fleet
	return out, c.do(ctx, http.MethodGet, "/v1/fleets/"+url.PathEscape(fleetID), nil, &out)
}

func (c *Client) FleetConfig(ctx context.Context, fleetID string, params types.FleetConfigParams) (*fluentbitconfig.Config, error) {
	q := url.Values{}
	if params.ConfigFormat != nil {
		q.Set("format", string(*params.ConfigFormat))
	}
	path := "/v1/fleets/" + url.PathEscape(fleetID) + "/config?" + q.Encode()

	switch *params.ConfigFormat {
	case types.ConfigFormatINI:
		var out string
		err := c.do(ctx, http.MethodGet, path, nil, &out)
		if err != nil {
			return nil, err
		}
		cfg, err := fluentbitconfig.ParseAs(out, fluentbitconfig.FormatClassic)
		if err != nil {
			return nil, err
		}
		return &cfg, nil
	case types.ConfigFormatYAML:
		fallthrough
	case types.ConfigFormatJSON:
		var out fluentbitconfig.Config
		return &out, c.do(ctx, http.MethodGet, path, nil, &out)
	}
	return nil, nil
}

func (c *Client) UpdateFleet(ctx context.Context, in types.UpdateFleet) (types.UpdatedFleet, error) {
	var out types.UpdatedFleet
	return out, c.do(ctx, http.MethodPatch, "/v1/fleets/"+url.PathEscape(in.ID), in, &out)
}

func (c *Client) DeleteFleet(ctx context.Context, fleetID string) (types.DeletedFleet, error) {
	var out types.DeletedFleet
	return out, c.do(ctx, http.MethodDelete, "/v1/fleets/"+url.PathEscape(fleetID), nil, &out)
}
