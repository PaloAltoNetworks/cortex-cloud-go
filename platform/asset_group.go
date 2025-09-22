package platform

import (
	"context"
	"net/http"
	"strconv"

	//"encoding/json"

	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/app"
	"github.com/PaloAltoNetworks/cortex-cloud-go/types"
)

// AssetGroup defines the structure for an asset group as returned by the API.
type AssetGroup struct {
	ID                  int         `json:"XDM.ASSET_GROUP.ID"`
	Name                string      `json:"XDM.ASSET_GROUP.NAME"`
	Type                string      `json:"XDM.ASSET_GROUP.TYPE"`
	Description         string     `json:"XDM.ASSET_GROUP.DESCRIPTION"`
	Filter              []AssetGroupFilter    `json:"XDM.ASSET_GROUP.FILTER"`
	CreationTime        int64       `json:"XDM.ASSET_GROUP.CREATION_TIME"`
	CreatedBy           string      `json:"XDM.ASSET_GROUP.CREATED_BY"`
	CreatedByPretty     string      `json:"XDM.ASSET_GROUP.CREATED_BY_PRETTY"`
	LastUpdateTime      int64       `json:"XDM.ASSET_GROUP.LAST_UPDATE_TIME"`
	ModifiedBy          string      `json:"XDM.ASSET_GROUP.MODIFIED_BY"`
	ModifiedByPretty    string      `json:"XDM.ASSET_GROUP.MODIFIED_BY_PRETTY"`
	MembershipPredicate types.CriteriaFilter `json:"XDM.ASSET_GROUP.MEMBERSHIP_PREDICATE"`
	IsUsedBySBAC        bool        `json:"IS_USED_BY_SBAC"`
}

// AssetGroupFilter represents a filter component in the asset group list response.
type AssetGroupFilter struct {
	PrettyName string  `json:"pretty_name"`
	DataType   string `json:"data_type"`
	RenderType string  `json:"render_type"`
	EntityMap  any     `json:"entity_map"`
	DMLType    any     `json:"dml_type"`
}

type createAssetGroupRequestWrapper struct {
	AssetGroup CreateAssetGroupRequest `json:"asset_group"`
}

// CreateAssetGroupRequest is the request for creating an asset group.
type CreateAssetGroupRequest struct {
	GroupName           string `json:"group_name"`
	GroupType           string `json:"group_type"`
	GroupDescription    string `json:"group_description,omitempty"`
	MembershipPredicate any    `json:"membership_predicate,omitempty"`
}

// UpdateAssetGroupRequest is the request for updating an asset group.
type UpdateAssetGroupRequest struct {
	GroupName           string `json:"group_name"`
	GroupType           string `json:"group_type"`
	GroupDescription    string `json:"group_description,omitempty"`
	MembershipPredicate types.CriteriaFilter `json:"membership_predicate"`
}

// CreateAssetGroupResponseData is the response for creating an asset group.
type CreateAssetGroupResponseData struct {
	Success      bool `json:"success"`
	AssetGroupID int  `json:"asset_group_id"`
}

type CreateAssetGroupResponseWrapper struct {
	Data CreateAssetGroupResponseData `json:"data"`
}

// GenericAssetGroupsResponseData is a generic response for asset group operations.
type GenericAssetGroupsResponseData struct {
	Success bool `json:"success"`
}

type GenericAssetGroupsResponseWrapper struct {
	Data GenericAssetGroupsResponseData `json:"data"`
}

// ListAssetGroupsRequest is the request for listing asset groups.
type ListAssetGroupsRequest struct {
	Filters    types.CriteriaFilter `json:"filters"` // Can be AndFilterForGroups or OrFilterForGroups
	//Sort       []types.SortFilter   `json:"-"`
	Sort       []types.SortFilter   `json:"sort,omitempty"`
	SearchFrom int                  `json:"search_from,omitempty"`
	SearchTo   int                  `json:"search_to,omitempty"`
}

// ListAssetGroupsResponse is the response for listing asset groups.
type ListAssetGroupsResponse struct {
	Data []AssetGroup `json:"data"`
}

// CreateAssetGroup creates a new asset group.
//
// TODO: make sure that the SEARCH_FIELD value is forced to uppercase
func (c *Client) CreateAssetGroup(ctx context.Context, req CreateAssetGroupRequest) (CreateAssetGroupResponseData, error) {
	var respWrapper CreateAssetGroupResponseWrapper
	_, err := c.internalClient.Do(ctx, http.MethodPost, CreateAssetGroupEndpoint, nil, nil, createAssetGroupRequestWrapper{ AssetGroup: req, }, &respWrapper, &app.DoOptions{
		RequestWrapperKey:  "request_data",
		ResponseWrapperKey: "reply",
	})
	return respWrapper.Data, err
}

// ListAssetGroups retrieves a list of asset groups.
func (c *Client) ListAssetGroups(ctx context.Context, req ListAssetGroupsRequest) (ListAssetGroupsResponse, error) {
	var resp ListAssetGroupsResponse
	_, err := c.internalClient.Do(ctx, http.MethodPost, ListAssetGroupsEndpoint, nil, nil, req, &resp, &app.DoOptions{
		RequestWrapperKey: "request_data",
		ResponseWrapperKey: "reply",
	})
	return resp, err
}

// UpdateAssetGroup updates an existing asset group.
func (c *Client) UpdateAssetGroup(ctx context.Context, groupID string, req UpdateAssetGroupRequest) (GenericAssetGroupsResponseData, error) {
	var respWrapper GenericAssetGroupsResponseWrapper
	pathParams := &[]string{groupID}
	_, err := c.internalClient.Do(ctx, http.MethodPost, UpdateAssetGroupEndpoint, pathParams, nil, req, &respWrapper, &app.DoOptions{
		RequestWrapperKey:  "asset_group",
		ResponseWrapperKey: "reply",
	})
	return respWrapper.Data, err
}

// DeleteAssetGroup deletes an asset group.
func (c *Client) DeleteAssetGroup(ctx context.Context, groupID int) (GenericAssetGroupsResponseData, error) {
	var respWrapper GenericAssetGroupsResponseWrapper
	pathParams := &[]string{strconv.Itoa(groupID)}
	_, err := c.internalClient.Do(ctx, http.MethodPost, DeleteAssetGroupEndpoint, pathParams, nil, nil, &respWrapper, &app.DoOptions{
		ResponseWrapperKey: "reply",
	})
	return respWrapper.Data, err
}
