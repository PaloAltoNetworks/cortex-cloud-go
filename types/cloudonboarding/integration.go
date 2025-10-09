package types

import (
	"encoding/json"
	"fmt"
	"net/url"

	filterTypes "github.com/PaloAltoNetworks/cortex-cloud-go/types/filter"
)

// ----------------------------------------------------------------------------
// Cloud Integration Instance Management
// ----------------------------------------------------------------------------

type IntegrationInstance struct {
	ID                      string                  `json:"id"`
	Collector               string                  `json:"collector"`
	InstanceName            string                  `json:"instance_name"`
	AccountName             string                  `json:"account_name,omitempty"`
	Accounts                int                     `json:"accounts,omitempty"`
	Scope                   string                  `json:"scope"`
	CustomResourcesTags     []Tag                   `json:"tags"`
	Scan                    Scan                    `json:"scan"`
	Status                  string                  `json:"status"`
	CloudProvider           string                  `json:"cloud_provider"`
	SecurityCapabilities    []SecurityCapability    `json:"security_capabilities"`
	CollectionConfiguration CollectionConfiguration `json:"collection_configuration"`
	AdditionalCapabilities  AdditionalCapabilities  `json:"additional_capabilities"`
	CreationTime            int                     `json:"creation_time,omitempty"`
	ProvisioningMethod      string                  `json:"provisioning_method,omitempty"`
	UpdateStatus            string                  `json:"update_status,omitempty"`
	UpgradeAvailable        bool                    `json:"upgrade_available,omitempty"`
	IsPendingChanges        int                     `json:"is_pending_changes,omitempty"`
	OutpostID               string                  `json:"outpost_id,omitempty"`
}

type Tag struct {
	Key   string `json:"key" tfsdk:"key"`
	Value string `json:"value" tfsdk:"value"`
}

type Scan struct {
	StatusUI   int    `json:"status_ui,omitempty" tfsdk:"status_ui"`
	OutpostID  string `json:"outpost_id,omitempty" tfsdk:"outpost_id"`
	ScanMethod string `json:"scan_method" tfsdk:"scan_method"`
}

type SecurityCapability struct {
	Name             string            `json:"name" tfsdk:"name"`
	Description      string            `json:"description" tfsdk:"description"`
	Status           int               `json:"status" tfsdk:"status"`
	LastScanCoverage *LastScanCoverage `json:"last_scan_coverage,omitempty" tfsdk:"last_scan_coverage"`
}

type LastScanCoverage struct {
	Excluded    int `json:"excluded" tfsdk:"excluded"`
	Issues      int `json:"issues" tfsdk:"issues"`
	Pending     int `json:"pending" tfsdk:"pending"`
	Success     int `json:"success" tfsdk:"success"`
	Unsupported int `json:"unsupported" tfsdk:"unsupported"`
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
	//ProjectIDs      []string `json:"project_ids,omitempty" tfsdk:"project_ids"`
	//SubscriptionIDs []string `json:"subscription_ids,omitempty" tfsdk:"subscription_ids"`
}

type ScopeModificationsOptionsRegions struct {
	Enabled bool     `json:"enabled" tfsdk:"enabled"`
	Type    string   `json:"type,omitempty" tfsdk:"type"`
	Regions []string `json:"regions,omitempty" tfsdk:"regions"`
}

type DefaultScanningScope struct {
	RegistryScanningScope      RegistryScanningScope      `json:"registry_scanning_scope"`
	AgentlessDiskScanningScope AgentlessDiskScanningScope `json:"agentless_disk_scanning_scope"`
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
	AgentlessDiskScanning         bool                    `json:"agentless_disk_scanning" tfsdk:"agentless_disk_scanning"`
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

// CreateIntegrationTemplateRequest is the request for creating an integration template.
type CreateIntegrationTemplateRequest struct {
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

// CreateTemplateOrEditIntegrationInstanceResponse is the response for creating or editing an integration instance.
type CreateTemplateOrEditIntegrationInstanceResponse struct {
	Automated Automated `json:"automated"`
	Manual    Manual    `json:"manual"`
}

func (r CreateTemplateOrEditIntegrationInstanceResponse) GetTemplateUrl() (string, error) {
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

// GetIntegrationInstanceRequest is the request for getting integration instance details.
type GetIntegrationInstanceRequest struct {
	InstanceID string `json:"id"`
}

// GetIntegrationInstanceResponse is the response for getting integration instance details.
type GetIntegrationInstanceResponse struct {
	ID                      string               `json:"id"`
	Collector               string               `json:"collector"`
	InstanceName            string               `json:"instance_name"`
	Scope                   string               `json:"scope"`
	Tags                    []Tag                `json:"tags"`
	Scan                    Scan                 `json:"scan"`
	Status                  string               `json:"status"`
	CloudProvider           string               `json:"cloud_provider"`
	SecurityCapabilities    []SecurityCapability `json:"security_capabilities"`
	CollectionConfiguration string               `json:"collection_configuration"`
	AdditionalCapabilities  string               `json:"additional_capabilities"`
	UpgradeAvailable        bool                 `json:"upgrade_available,omitempty"`
}

func (r GetIntegrationInstanceResponse) Marshal() (IntegrationInstance, error) {
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
		CustomResourcesTags:     r.Tags,
		Scan:                    r.Scan,
		UpgradeAvailable:        r.UpgradeAvailable,
		Status:                  r.Status,
		CloudProvider:           r.CloudProvider,
		SecurityCapabilities:    r.SecurityCapabilities,
		CollectionConfiguration: collectionConfiguration,
		AdditionalCapabilities:  additionalCapabilities,
	}

	return marshalledResponse, nil
}

// ListIntegrationInstancesRequest is the request for listing integration instances.
type ListIntegrationInstancesRequest struct {
	FilterData filterTypes.FilterData `json:"filter_data"`
}

// ListIntegrationInstancesResponseWrapper is the response wrapper for listing integration instances.
type ListIntegrationInstancesResponseWrapper struct {
	Data []ListIntegrationInstancesResponse `json:"DATA"`
}

// ListIntegrationInstancesResponse is the response for listing integration instances.
type ListIntegrationInstancesResponse struct {
	InstanceName        string `json:"instance_name"`
	CloudProvider       string `json:"cloud_provider"`
	Accounts            int    `json:"accounts,omitempty"`
	AccountName         string `json:"account_name,omitempty"`
	Scope               string `json:"scope"`
	ScanMode            string `json:"scan_mode"`
	CustomResourcesTags string `json:"custom_resources_tags"`
	ProvisioningMethod  string `json:"provisioning_method"`
	//AccountDetails          AccountDetails     `json:"account_details"`
	//ScopeModifications      ScopeModifications `json:"scope_modifications"`
	CollectionConfiguration string `json:"collection_configuration"`
	AdditionalCapabilities  string `json:"additional_capabilities"`
	InstanceID              string `json:"instance_id"`
	Status                  string `json:"status"`
	//CloudPartition          string               `json:"cloud_partition"`
	//ModifiedAt              int                  `json:"modified_at"`
	DeletedAt            int                  `json:"deleted_at"`
	DefaultScanningScope DefaultScanningScope `json:"default_scanning_scope"`
	OutpostID            string               `json:"outpost_id"`
	CreationTime         int                  `json:"creation_time,omitempty"`
	UpdateStatus         string               `json:"update_status,omitempty"`
	IsPendingChanges     int                  `json:"is_pending_changes,omitempty"`
}

func (r ListIntegrationInstancesResponseWrapper) Marshal() ([]IntegrationInstance, error) {
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
			AccountName:             data.AccountName,
			Accounts:                data.Accounts,
			CreationTime:            data.CreationTime,
			ProvisioningMethod:      data.ProvisioningMethod,
			UpdateStatus:            data.UpdateStatus,
			IsPendingChanges:        data.IsPendingChanges,
			OutpostID:               data.OutpostID,
		}

		marshalledResponse = append(marshalledResponse, marshalledData)
	}

	return marshalledResponse, nil
}

// EditIntegrationInstanceRequest is the request for editing an integration instance.
type EditIntegrationInstanceRequest struct {
	InstanceID              string                  `json:"id"`
	ScanEnvID               string                  `json:"scan_env_id"`
	InstanceName            string                  `json:"instance_name"`
	AdditionalCapabilities  AdditionalCapabilities  `json:"additional_capabilities"`
	CloudProvider           string                  `json:"cloud_provider"`
	CustomResourcesTags     []Tag                   `json:"custom_resources_tags"`
	CollectionConfiguration CollectionConfiguration `json:"collection_configuration"`
	ScopeModifications      ScopeModifications      `json:"scope_modifications"`
}

// EnableOrDisableIntegrationInstancesRequest is the request for enabling or disabling integration instances.
type EnableOrDisableIntegrationInstancesRequest struct {
	IDs    []string `json:"ids"`
	Enable bool     `json:"enable"`
}

// DeleteIntegrationInstanceRequest is the request for deleting integration instances.
type DeleteIntegrationInstanceRequest struct {
	IDs []string `json:"ids"`
}
