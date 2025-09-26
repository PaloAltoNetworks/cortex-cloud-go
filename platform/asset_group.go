package platform

import (
	"context"
	"net/http"
	"strconv"

	"github.com/PaloAltoNetworks/cortex-cloud-go/client"
	"github.com/PaloAltoNetworks/cortex-cloud-go/types"
)

// CreateAssetGroup creates a new asset group.
//
// TODO: make sure that the SEARCH_FIELD value is forced to uppercase
func (c *Client) CreateAssetGroup(ctx context.Context, req types.CreateOrUpdateAssetGroupRequest) (types.CreateAssetGroupResponseData, error) {
	var respWrapper types.CreateAssetGroupResponseWrapper
	_, err := c.internalClient.Do(ctx, http.MethodPost, CreateAssetGroupEndpoint, nil, nil, types.CreateOrUpdateAssetGroupRequestWrapper{AssetGroup: req}, &respWrapper, &client.DoOptions{
		RequestWrapperKey:  "request_data",
		ResponseWrapperKey: "reply",
	})
	return respWrapper.Data, err
}

// ListAssetGroups retrieves a list of asset groups.
func (c *Client) ListAssetGroups(ctx context.Context, req types.ListAssetGroupsRequest) (types.ListAssetGroupsResponse, error) {
	var resp types.ListAssetGroupsResponse
	_, err := c.internalClient.Do(ctx, http.MethodPost, ListAssetGroupsEndpoint, nil, nil, req, &resp, &client.DoOptions{
		RequestWrapperKey:  "request_data",
		ResponseWrapperKey: "reply",
	})
	return resp, err
}

// UpdateAssetGroup updates an existing asset group.
func (c *Client) UpdateAssetGroup(ctx context.Context, groupID int, req types.CreateOrUpdateAssetGroupRequest) (types.GenericAssetGroupsResponseData, error) {
	var respWrapper types.GenericAssetGroupsResponseWrapper
	pathParams := &[]string{strconv.Itoa(groupID)}
	_, err := c.internalClient.Do(ctx, http.MethodPost, UpdateAssetGroupEndpoint, pathParams, nil, types.CreateOrUpdateAssetGroupRequestWrapper{AssetGroup: req}, &respWrapper, &client.DoOptions{
		RequestWrapperKey:  "request_data",
		ResponseWrapperKey: "reply",
	})
	return respWrapper.Data, err
}

// DeleteAssetGroup deletes an asset group.
func (c *Client) DeleteAssetGroup(ctx context.Context, groupID int) (types.GenericAssetGroupsResponseData, error) {
	var respWrapper types.GenericAssetGroupsResponseWrapper
	pathParams := &[]string{strconv.Itoa(groupID)}
	_, err := c.internalClient.Do(ctx, http.MethodPost, DeleteAssetGroupEndpoint, pathParams, nil, nil, &respWrapper, &client.DoOptions{
		ResponseWrapperKey: "reply",
	})
	return respWrapper.Data, err
}
