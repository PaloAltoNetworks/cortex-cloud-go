package types

// CreateOutpostTemplateRequest is the request for the CreateOutpostTemplate endpoint.
type CreateOutpostTemplateRequest struct {
	CloudProvider      string `json:"cloud_provider"`
	CustomResourceTags []Tag  `json:"custom_resources_tags,omitempty"`
}

// UpdateOutpostRequest is the request for the UpdateOutpost endpoint.
// The OpenAPI spec for this endpoint is faulty, so this struct is a best guess.
type UpdateOutpostRequest struct {
	OutpostID          string `json:"outpost_id"`
	CloudProvider      string `json:"cloud_provider"`
	CustomResourceTags []Tag  `json:"custom_resources_tags,omitempty"`
}

// ListOutpostsRequest is the request for the ListOutposts endpoint.
type ListOutpostsRequest = ListIntegrationInstancesRequest

// Outpost represents an outpost object.
type Outpost struct {
	CloudProvider string `json:"cloud_provider"`
	OutpostID     string `json:"outpost_id"`
	CreatedAt     int    `json:"created_at"`
	Type          string `json:"type"`
}

// ListOutpostsResponse is the response for the ListOutposts endpoint.
type ListOutpostsResponse struct {
	Data        []Outpost `json:"DATA"`
	FilterCount int       `json:"FILTER_COUNT"`
	TotalCount  int       `json:"TOTAL_COUNT"`
}
