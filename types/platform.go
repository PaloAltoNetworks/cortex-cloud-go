package types

// ----------------------------------------------------------------------------
// Asset Group
// ----------------------------------------------------------------------------

// AssetGroup defines the structure for an asset group as returned by the API.
type AssetGroup struct {
	ID                  int                `json:"XDM.ASSET_GROUP.ID"`
	Name                string             `json:"XDM.ASSET_GROUP.NAME"`
	Type                string             `json:"XDM.ASSET_GROUP.TYPE"`
	Description         string             `json:"XDM.ASSET_GROUP.DESCRIPTION"`
	Filter              []AssetGroupFilter `json:"XDM.ASSET_GROUP.FILTER"`
	CreationTime        int64              `json:"XDM.ASSET_GROUP.CREATION_TIME"`
	CreatedBy           string             `json:"XDM.ASSET_GROUP.CREATED_BY"`
	CreatedByPretty     string             `json:"XDM.ASSET_GROUP.CREATED_BY_PRETTY"`
	LastUpdateTime      int64              `json:"XDM.ASSET_GROUP.LAST_UPDATE_TIME"`
	ModifiedBy          string             `json:"XDM.ASSET_GROUP.MODIFIED_BY"`
	ModifiedByPretty    string             `json:"XDM.ASSET_GROUP.MODIFIED_BY_PRETTY"`
	MembershipPredicate CriteriaFilter     `json:"XDM.ASSET_GROUP.MEMBERSHIP_PREDICATE"`
	IsUsedBySBAC        bool               `json:"IS_USED_BY_SBAC"`
}

// AssetGroupFilter represents a filter component in the asset group list response.
type AssetGroupFilter struct {
	PrettyName string `json:"pretty_name"`
	DataType   string `json:"data_type"`
	RenderType string `json:"render_type"`
	EntityMap  any    `json:"entity_map"`
	DMLType    any    `json:"dml_type"`
}

type CreateOrUpdateAssetGroupRequestWrapper struct {
	AssetGroup CreateOrUpdateAssetGroupRequest `json:"asset_group"`
}

// CreateOrUpdateAssetGroupRequest is the request for creating or updating an
// asset group.
type CreateOrUpdateAssetGroupRequest struct {
	GroupName           string         `json:"group_name"`
	GroupType           string         `json:"group_type"`
	GroupDescription    string         `json:"group_description,omitempty"`
	MembershipPredicate CriteriaFilter `json:"membership_predicate"`
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
	Filters    CriteriaFilter `json:"filters"`
	Sort       []SortFilter   `json:"sort,omitempty"`
	SearchFrom int            `json:"search_from,omitempty"`
	SearchTo   int            `json:"search_to,omitempty"`
}

// ListAssetGroupsResponse is the response for listing asset groups.
type ListAssetGroupsResponse struct {
	Data []AssetGroup `json:"data"`
}

// ----------------------------------------------------------------------------
// Authentication Settings
// ----------------------------------------------------------------------------

type AuthSettings struct {
	TenantID           string           `json:"tenant_id"`
	Name               string           `json:"name"`
	Domain             string           `json:"domain"`
	IDPEnabled         bool             `json:"idp_enabled"`
	DefaultRole        string           `json:"default_role"`
	IsAccountRole      bool             `json:"is_account_role"`
	IDPCertificate     string           `json:"idp_certificate"`
	IDPIssuer          string           `json:"idp_issuer"`
	IDPSingleSignOnURL string           `json:"idp_sso_url"`
	MetadataURL        string           `json:"metadata_url"`
	Mappings           Mappings         `json:"mappings"`
	AdvancedSettings   AdvancedSettings `json:"advanced_settings"`
	SpEntityID         string           `json:"sp_entity_id"`
	SpLogoutURL        string           `json:"sp_logout_url"`
	SpURL              string           `json:"sp_url"`
}

type Mappings struct {
	Email     string `json:"email"`
	FirstName string `json:"firstname"`
	GroupName string `json:"group_name"`
	LastName  string `json:"lastname"`
}

type AdvancedSettings struct {
	AuthnContextEnabled       bool `json:"authn_context_enabled"`
	IDPSingleLogoutURL        string `json:"idp_single_logout_url"`
	RelayState                string `json:"relay_state"`
	ServiceProviderPrivateKey string `json:"service_provider_private_key"`
	ServiceProviderPublicCert string `json:"service_provider_public_cert"`
}

// ListIDPMetadataRequest is the request for listing IDP metadata.
type ListIDPMetadataRequest struct{}

// ListIDPMetadataResponse is the response for listing IDP metadata.
type ListIDPMetadataResponse struct {
	TenantID    string `json:"tenant_id"`
	SpEntityID  string `json:"sp_entity_id"`
	SpLogoutURL string `json:"sp_logout_url"`
	SpURL       string `json:"sp_url"`
}

// ListAuthSettingsRequest is the request for listing auth settings.
type ListAuthSettingsRequest struct{}

// CreateAuthSettingsRequest is the request for creating auth settings.
type CreateAuthSettingsRequest struct {
	Name               string           `json:"name"`
	DefaultRole        string           `json:"default_role"`
	IsAccountRole      bool             `json:"is_account_role"`
	Domain             string           `json:"domain"`
	Mappings           Mappings         `json:"mappings"`
	AdvancedSettings   AdvancedSettings `json:"advanced_settings"`
	IDPSingleSignOnURL string           `json:"idp_sso_url"`
	IDPCertificate     string           `json:"idp_certificate"`
	IDPIssuer          string           `json:"idp_issuer"`
	MetadataURL        string           `json:"metadata_url"`
}

// UpdateAuthSettingsRequest is the request for updating auth settings.
type UpdateAuthSettingsRequest struct {
	Name               string           `json:"name"`
	DefaultRole        string           `json:"default_role"`
	IsAccountRole      bool             `json:"is_account_role"`
	CurrentDomain      string           `json:"current_domain_value"`
	NewDomain          string           `json:"new_domain_value"`
	Mappings           Mappings         `json:"mappings"`
	AdvancedSettings   AdvancedSettings `json:"advanced_settings"`
	IDPSingleSignOnURL string           `json:"idp_sso_url"`
	IDPCertificate     string           `json:"idp_certificate"`
	IDPIssuer          string           `json:"idp_issuer"`
	MetadataURL        string           `json:"metadata_url"`
}

// DeleteAuthSettingsRequest is the request for deleting auth settings.
type DeleteAuthSettingsRequest struct {
	Domain string `json:"domain"`
}

// ----------------------------------------------------------------------------
// System Management
// ----------------------------------------------------------------------------

type User struct {
	UserEmail     string   `json:"user_email"`
	UserFirstName string   `json:"user_first_name"`
	UserLastName  string   `json:"user_last_name"`
	RoleName      string   `json:"role_name"`
	LastLoggedIn  int      `json:"last_logged_in"`
	UserType      string   `json:"user_type"`
	Groups        []string `json:"groups"`
	Scope         Scope    `json:"scope"`
}

type Scope struct {
	Endpoints   Endpoints   `json:"endpoints"`
	CasesIssues CasesIssues `json:"cases_issues"`
}

type Endpoints struct {
	EndpointGroups EndpointGroups `json:"endpoint_groups"`
	EndpointTags   EndpointTags   `json:"endpoint_tags"`
	Mode           string         `json:"mode"`
}

type EndpointGroups struct {
	IDs  []string `json:"ids"`
	Mode string   `json:"mode"`
}

type EndpointTags struct {
	IDs  []string `json:"ids"`
	Mode string   `json:"mode"`
}

type CasesIssues struct {
	IDs  []string `json:"ids"`
	Mode string   `json:"mode"`
}

type Reason struct {
	DateCreated string `json:"date_created"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Status      string `json:"status"`
	Points      int    `json:"points"`
}

// GetUserRequest is the request for getting a user.
type GetUserRequest struct {
	Email string `json:"email"`
}

// ListRolesRequest is the request for listing roles.
type ListRolesRequest struct {
	RoleNames []string `json:"role_names" validate:"required,min=1"`
}

// ListRolesResponse is the response for listing roles.
type ListRolesResponse struct {
	PrettyName  string   `json:"pretty_name"`
	Permissions []string `json:"permissions"`
	InsertTime  int      `json:"insert_time"`
	UpdateTime  int      `json:"update_time"`
	CreatedBy   string   `json:"created_by"`
	Description string   `json:"description"`
	Tags        string   `json:"tags"`
	Groups      []string `json:"groups"`
	Users       []string `json:"users"`
}

// SetRoleRequest is the request for setting a role.
type SetRoleRequest struct {
	UserEmails []string `json:"user_emails" validate:"required,min=1,dive,required,email"`
	RoleName   string   `json:"role_name"`
}

// SetRoleResponse is the response for setting a role.
type SetRoleResponse struct {
	UpdateCount string `json:"update_count"`
}

// GetRiskScoreRequest is the request for getting a risk score.
type GetRiskScoreRequest struct {
	ID string `json:"id" validate:"required,sysmgmtID"`
}

// GetRiskScoreResponse is the response for getting a risk score.
type GetRiskScoreResponse struct {
	Type          string   `json:"type"`
	ID            string   `json:"id"`
	Score         int      `json:"score"`
	NormRiskScore int      `json:"norm_risk_score"`
	RiskLevel     string   `json:"risk_level"`
	Reasons       []Reason `json:"reasons"`
	Email         string   `json:"email"`
}

// ListRiskyUsersResponse is the response for listing risky users.
type ListRiskyUsersResponse struct {
	Type          string   `json:"type"`
	ID            string   `json:"id"`
	Score         int      `json:"score"`
	NormRiskScore int      `json:"norm_risk_score"`
	RiskLevel     string   `json:"risk_level"`
	Reasons       []Reason `json:"reasons"`
	Email         string   `json:"email"`
}

// ListRiskyHostsResponse is the response for listing risky hosts.
type ListRiskyHostsResponse struct {
	Type          string   `json:"type"`
	ID            string   `json:"id"`
	Score         int      `json:"score"`
	NormRiskScore int      `json:"norm_risk_score"`
	RiskLevel     string   `json:"risk_level"`
	Reasons       []Reason `json:"reasons"`
}
