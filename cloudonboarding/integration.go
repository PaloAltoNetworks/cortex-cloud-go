// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package cloudonboarding

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/app"
)

// ----------------------------------------------------------------------------
// Core structs
// ----------------------------------------------------------------------------

type IntegrationInstance struct {
	ID                      string                  `json:"id" tfsdk:"id"`
	Collector               string                  `json:"collector" tfsdk:"collector"`
	InstanceName            string                  `json:"instance_name" tfsdk:"instance_name"`
	Scope                   string                  `json:"scope" tfsdk:"scope"`
	CustomResourcesTags     []Tag                   `json:"tags" tfsdk:"custom_resource_tags"`
	Scan                    Scan                    `json:"scan" tfsdk:"scan"`
	Status                  string                  `json:"status" tfsdk:"status"`
	CloudProvider           string                  `json:"cloud_provider" tfsdk:"cloud_provider"`
	SecurityCapabilities    []SecurityCapability    `json:"security_capabilities" tfsdk:"security_capabilities"`
	CollectionConfiguration CollectionConfiguration `json:"collection_configuration"`
	AdditionalCapabilities  AdditionalCapabilities  `json:"additional_capabilities"`
}

type Tag struct {
	Key   string `json:"key" tfsdk:"key"`
	Value string `json:"value" tfsdk:"value"`
}

type Scan struct {
	ScanMethod string `json:"scan_method" tfsdk:"scan_method"`
}

type SecurityCapability struct {
	Name        string `json:"name" tfsdk:"name"`
	Description string `json:"description" tfsdk:"description"`
	Status      int    `json:"status" tfsdk:"status"`
}

type AccountDetails struct {
	OrganizationID *string `json:"organization_id,omitempty"`
}

type CollectionConfiguration struct {
	AuditLogs AuditLogsConfiguration `json:"audit_logs" tfsdk:"audit_logs"`
}

type AuditLogsConfiguration struct {
	Enabled bool `json:"enabled" tfsdk:"enabled"`
}

type ScopeModifications struct {
	Accounts      *ScopeModificationsOptionsGeneric `json:"accounts,omitempty" tfsdk:"accounts"`
	Projects      *ScopeModificationsOptionsGeneric `json:"projects,omitempty" tfsdk:"projects"`
	Subscriptions *ScopeModificationsOptionsGeneric `json:"subscriptions,omitempty" tfsdk:"subscriptions"`
	Regions       *ScopeModificationsOptionsRegions `json:"regions,omitempty" tfsdk:"regions"`
}

type ScopeModificationsOptionsGeneric struct {
	Enabled         bool     `json:"enabled" tfsdk:"enabled"`
	Type            string   `json:"type,omitempty" tfsdk:"type"`
	AccountIDs      []string `json:"account_ids,omitempty" tfsdk:"account_ids"`
	ProjectIDs      []string `json:"project_ids,omitempty" tfsdk:"project_ids"`
	SubscriptionIDs []string `json:"subscription_ids,omitempty" tfsdk:"subscription_ids"`
}

type ScopeModificationsOptionsRegions struct {
	Enabled bool     `json:"enabled" tfsdk:"enabled"`
	Type    string   `json:"type,omitempty" tfsdk:"type"`
	Regions []string `json:"regions,omitempty" tfsdk:"regions"`
}

type DefaultScanningScope struct {
	RegistryScanningScope      RegistryScanningScope      `json:"registry_scanning_scope"`
	AgentlessDiskScanningScope AgentlessDiskScanningScope `json:"agentless_disk_scanning_scope"`
	// TODO: DataAssetsClassificationOptions
}

type RegistryScanningScope struct {
	Enabled bool `json:"enabled"`
}

type AgentlessDiskScanningScope struct {
	Enabled bool `json:"enabled"`
}

type AdditionalCapabilities struct {
	XsiamAnalytics                bool                    `json:"xsiam_analytics" tfsdk:"xsiam_analytics"`
	DataSecurityPostureManagement bool                    `json:"data_security_posture_management" tfsdk:"data_security_posture_management"`
	RegistryScanning              bool                    `json:"registry_scanning" tfsdk:"registry_scanning"`
	RegistryScanningOptions       RegistryScanningOptions `json:"registry_scanning_options" tfsdk:"registry_scanning_options"`
}

type RegistryScanningOptions struct {
	Type string `json:"type" tfsdk:"type"`
}

type Automated struct {
	Link         string `json:"link" tfsdk:"automated_deployment_link"`
	TrackingGuid string `json:"tracking_guid" tfsdk:"tracking_guid"`
}

type Manual struct {
	CF string `json:"CF" tfsdk:"manual_deployment_link"`
}

// ----------------------------------------------------------------------------
// Create Integration Template
// ----------------------------------------------------------------------------

type CreateIntegrationTemplateRequestData struct {
	AccountDetails          *AccountDetails         `json:"account_details,omitempty"`
	AdditionalCapabilities  AdditionalCapabilities  `json:"additional_capabilities"`
	CloudProvider           string                  `json:"cloud_provider"`
	CollectionConfiguration CollectionConfiguration `json:"collection_configuration"`
	CustomResourcesTags     []Tag                   `json:"custom_resources_tags"`
	InstanceName            string                  `json:"instance_name"`
	ScanMode                string                  `json:"scan_mode"`
	Scope                   string                  `json:"scope"`
	ScopeModifications      ScopeModifications      `json:"scope_modifications"`
}

type CreateTemplateOrEditIntegrationInstanceResponseReply struct {
	Automated Automated `json:"automated"`
	Manual    Manual    `json:"manual"`
}

func (r CreateTemplateOrEditIntegrationInstanceResponseReply) GetTemplateUrl() (string, error) {
	if r.Automated.Link == "" {
		return "", fmt.Errorf("Failed to retrieve template URL: reply.automated.link is empty string")
	}

	parsedUrl, err := url.Parse(r.Automated.Link)
	if err != nil {
		return "", err
	}

	queryValues, err := url.ParseQuery(parsedUrl.RawFragment)
	if err != nil {
		return "", err
	}

	templateUrl := queryValues.Get("/stacks/quickcreate?templateURL")

	return templateUrl, nil
}

// CreateTemplate creates a new Cloud Onboarding Integration Template.
//
// TODO: details
func (c *Client) CreateIntegrationTemplate(ctx context.Context, input CreateIntegrationTemplateRequestData) (CreateTemplateOrEditIntegrationInstanceResponseReply, error) {
	var ans CreateTemplateOrEditIntegrationInstanceResponseReply
	_, err := c.internalClient.Do(ctx, http.MethodPost, CreateIntegrationTemplateEndpoint, nil, nil, input, &ans, &app.DoOptions{
		RequestWrapperKey:  "request_data",
		ResponseWrapperKey: "reply",
	})

	return ans, err
}

// ----------------------------------------------------------------------------
// Get Integration Instance Details
// ----------------------------------------------------------------------------

type GetIntegrationInstanceRequestData struct {
	InstanceID string `json:"id"`
}

type GetIntegrationInstanceResponseReply struct {
	ID                      string               `json:"id"`
	Collector               string               `json:"collector"`
	InstanceName            string               `json:"instance_name"`
	Scope                   string               `json:"scope"`
	CustomResourcesTags     []Tag                `json:"tags"`
	Scan                    Scan                 `json:"scan"`
	Status                  string               `json:"status"`
	CloudProvider           string               `json:"cloud_provider"`
	SecurityCapabilities    []SecurityCapability `json:"security_capabilities"`
	CollectionConfiguration string               `json:"collection_configuration"`
	AdditionalCapabilities  string               `json:"additional_capabilities"`
}

func (r GetIntegrationInstanceResponseReply) Marshal() (IntegrationInstance, error) {
	var collectionConfiguration CollectionConfiguration
	err := json.Unmarshal([]byte(r.CollectionConfiguration), &collectionConfiguration)
	if err != nil {
		return IntegrationInstance{}, err
	}

	var additionalCapabilities AdditionalCapabilities
	err = json.Unmarshal([]byte(r.AdditionalCapabilities), &additionalCapabilities)
	if err != nil {
		return IntegrationInstance{}, err
	}

	marshalledResponse := IntegrationInstance{
		ID:                      r.ID,
		Collector:               r.Collector,
		InstanceName:            r.InstanceName,
		Scope:                   r.Scope,
		CustomResourcesTags:     r.CustomResourcesTags,
		Scan:                    r.Scan,
		Status:                  r.Status,
		CloudProvider:           r.CloudProvider,
		SecurityCapabilities:    r.SecurityCapabilities,
		CollectionConfiguration: collectionConfiguration,
		AdditionalCapabilities:  additionalCapabilities,
	}

	return marshalledResponse, nil
}

// GetDetails returns the configuration details of the specified integration instance.
func (c *Client) GetIntegrationInstanceDetails(ctx context.Context, input GetIntegrationInstanceRequestData) (GetIntegrationInstanceResponseReply, error) {
	var ans GetIntegrationInstanceResponseReply
	_, err := c.internalClient.Do(ctx, http.MethodPost, GetIntegrationInstanceDetailsEndpoint, nil, nil, input, &ans, &app.DoOptions{
		RequestWrapperKey:  "request_data",
		ResponseWrapperKey: "reply",
	})

	return ans, err
}

// ----------------------------------------------------------------------------
// List Integration Instances
// ----------------------------------------------------------------------------

type ListIntegrationInstancesRequestData struct {
	FilterData FilterData `json:"filter_data"`
}

type ListIntegrationInstancesResponseReply struct {
	Data []ListIntegrationInstancesResponseData `json:"DATA"`
}

type ListIntegrationInstancesResponseData struct {
	InstanceName            string               `json:"instance_name"`
	CloudProvider           string               `json:"cloud_provider"`
	Scope                   string               `json:"scope"`
	ScanMode                string               `json:"scan_mode"`
	CustomResourcesTags     string               `json:"custom_resources_tags"`
	ProvisioningMethod      string               `json:"provisioning_method"`
	AccountDetails          AccountDetails       `json:"account_details"`
	ScopeModifications      ScopeModifications   `json:"scope_modifications"`
	CollectionConfiguration string               `json:"collection_configuration"`
	AdditionalCapabilities  string               `json:"additional_capabilities"`
	InstanceID              string               `json:"instance_id"`
	Status                  string               `json:"status"`
	CloudPartition          string               `json:"cloud_partition"`
	CreatedAt               int                  `json:"created_at"`
	ModifiedAt              int                  `json:"modified_at"`
	DeletedAt               int                  `json:"deleted_at"`
	DefaultScanningScope    DefaultScanningScope `json:"default_scanning_scope"`
	OutpostID               string               `json:"outpost_id"`
}

func (r ListIntegrationInstancesResponseReply) Marshal() ([]IntegrationInstance, error) {
	// TODO: make sure Paging.To is set to 1000 (the max accepted value)
	// if not configured.

	// TODO: Where is outpost_id populated? is there a static list of
	// outpost IDs for managed integrations?

	marshalledResponse := []IntegrationInstance{}

	for _, data := range r.Data {
		var customResourcesTags []Tag
		if data.CustomResourcesTags != "" {
			err := json.Unmarshal([]byte(data.CustomResourcesTags), &customResourcesTags)
			if err != nil {
				return []IntegrationInstance{}, err
			}
		} else {
			customResourcesTags = []Tag{}
		}

		// TODO: verify this behavior
		var collectionConfiguration CollectionConfiguration
		if data.CollectionConfiguration != "" {
			err := json.Unmarshal([]byte(data.CollectionConfiguration), &collectionConfiguration)
			if err != nil {
				return []IntegrationInstance{}, err
			}
		} else {
			collectionConfiguration = CollectionConfiguration{}
		}

		var additionalCapabilities AdditionalCapabilities
		if data.AdditionalCapabilities != "" {
			err := json.Unmarshal([]byte(data.AdditionalCapabilities), &additionalCapabilities)
			if err != nil {
				return []IntegrationInstance{}, err
			}
		} else {
			additionalCapabilities = AdditionalCapabilities{}
		}

		marshalledData := IntegrationInstance{
			ID:                      data.InstanceID,
			InstanceName:            data.InstanceName,
			Scope:                   data.Scope,
			CustomResourcesTags:     customResourcesTags,
			Scan:                    Scan{ScanMethod: data.ScanMode},
			Status:                  data.Status,
			CloudProvider:           data.CloudProvider,
			CollectionConfiguration: collectionConfiguration,
			AdditionalCapabilities:  additionalCapabilities,
		}

		marshalledResponse = append(marshalledResponse, marshalledData)
	}

	return marshalledResponse, nil
}

func (c *Client) ListIntegrationInstances(ctx context.Context, input ListIntegrationInstancesRequestData) (ListIntegrationInstancesResponseReply, error) {
	var ans ListIntegrationInstancesResponseReply
	_, err := c.internalClient.Do(ctx, http.MethodPost, ListIntegrationInstancesEndpoint, nil, nil, input, &ans, &app.DoOptions{
		RequestWrapperKey:  "request_data",
		ResponseWrapperKey: "reply",
	})

	return ans, err
}

// ----------------------------------------------------------------------------
// Edit Integration Instance
// ----------------------------------------------------------------------------

type EditIntegrationInstanceRequestData struct {
	InstanceID              string                  `json:"id"`
	ScanEnvID               string                  `json:"scan_env_id"`
	InstanceName            string                  `json:"instance_name"`
	AdditionalCapabilities  AdditionalCapabilities  `json:"additional_capabilities"`
	CloudProvider           string                  `json:"cloud_provider"`
	CustomResourcesTags     []Tag                   `json:"custom_resources_tags"`
	CollectionConfiguration CollectionConfiguration `json:"collection_configuration"`
	ScopeModifications      ScopeModifications      `json:"scope_modifications"`
}

func (c *Client) EditIntegrationInstance(ctx context.Context, input EditIntegrationInstanceRequestData) (CreateTemplateOrEditIntegrationInstanceResponseReply, error) {
	var ans CreateTemplateOrEditIntegrationInstanceResponseReply
	_, err := c.internalClient.Do(ctx, http.MethodPost, EditIntegrationInstanceEndpoint, nil, nil, input, &ans, &app.DoOptions{
		RequestWrapperKey:  "request_data",
		ResponseWrapperKey: "reply",
	})

	return ans, err
}

// ----------------------------------------------------------------------------
// Enable or Disable Instances
// ----------------------------------------------------------------------------

type EnableOrDisableIntegrationInstancesRequestData struct {
	IDs    []string `json:"ids"`
	Enable bool     `json:"enable"`
}

func (c *Client) EnableIntegrationInstances(ctx context.Context, instanceIDs []string) error {
	body := EnableOrDisableIntegrationInstancesRequestData{
		IDs:    instanceIDs,
		Enable: true,
	}

	_, err := c.internalClient.Do(ctx, http.MethodPost, EnableOrDisableIntegrationInstancesEndpoint, nil, nil, body, nil, &app.DoOptions{
		RequestWrapperKey: "request_data",
	})

	return err
}

func (c *Client) DisableIntegrationInstances(ctx context.Context, instanceIDs []string) error {
	body := EnableOrDisableIntegrationInstancesRequestData{
		IDs:    instanceIDs,
		Enable: false,
	}

	_, err := c.internalClient.Do(ctx, http.MethodPost, EnableOrDisableIntegrationInstancesEndpoint, nil, nil, body, nil, &app.DoOptions{
		RequestWrapperKey: "request_data",
	})

	return err
}

// ----------------------------------------------------------------------------
// Delete Integration Instances
// ----------------------------------------------------------------------------

type DeleteIntegrationInstanceRequestData struct {
	IDs []string `json:"ids"`
}

func (c *Client) DeleteIntegrationInstances(ctx context.Context, instanceIDs []string) error {
	body := DeleteIntegrationInstanceRequestData{
		IDs: instanceIDs,
	}

	_, err := c.internalClient.Do(ctx, http.MethodPost, DeleteIntegrationInstancesEndpoint, nil, nil, body, nil, &app.DoOptions{
		RequestWrapperKey: "request_data",
	})

	return err
}
