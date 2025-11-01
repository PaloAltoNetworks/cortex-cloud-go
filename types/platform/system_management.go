package types

type User struct {
	Email        string   `json:"user_email"`
	FirstName    string   `json:"user_first_name"`
	LastName     string   `json:"user_last_name"`
	RoleName     string   `json:"role_name"`
	LastLoggedIn int      `json:"last_logged_in"`
	UserType     string   `json:"user_type"`
	Groups       []string `json:"groups"`
	Scope        Scope    `json:"scope"`
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

type Role struct {
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

// GetUserRequest is the request for getting a user.
type GetUserRequest struct {
	Email string `json:"email"`
}

// SetRoleRequest is the request for setting a role.
type SetRoleRequest struct {
	UserEmails []string `json:"user_emails"`
	RoleName   string   `json:"role_name"`
}

// SetRoleResponse is the response for setting a role.
type SetRoleResponse struct {
	UpdateCount string `json:"update_count"`
}

// GetRiskScoreRequest is the request for getting a risk score.
type GetRiskScoreRequest struct {
	ID string `json:"id"`
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

// HealthCheckResponse defines the response for the health check endpoint.
type HealthCheckResponse struct {
	Service   string `json:"service"`
	Status    string `json:"status"`
	Reason    string `json:"reason"`
	Timestamp int64  `json:"timestamp"`
}

// GetTenantInfoRequest defines the request for the get_tenant_info endpoint.
type GetTenantInfoRequest struct {
	Tenants []string `json:"tenants,omitempty"`
}

// TenantInfoLicense defines a license associated with a tenant.
type TenantInfoLicense struct {
	LicenseID      string `json:"license_id"`
	LicenseType    string `json:"license_type"`
	LicenseName    string `json:"license_name"`
	ExpirationDate int64  `json:"expiration_date"`
	IsExpired      bool   `json:"is_expired"`
}

// TenantInfo defines information about a tenant.
type TenantInfo struct {
	TenantID       string              `json:"tenant_id"`
	TenantName     string              `json:"tenant_name"`
	DisplayName    string              `json:"display_name"`
	CustomerName   string              `json:"customer_name"`
	ParentTenantID string              `json:"parent_tenant_id"`
	TenantType     string              `json:"tenant_type"`
	IsActive       bool                `json:"is_active"`
	Licenses       []TenantInfoLicense `json:"licenses"`
}

// GetUserGroupRequest defines the request for the get_user_group endpoint.
type GetUserGroupRequest struct {
	GroupNames []string `json:"group_names"`
}

// NestedGroup represents a child group within a user group.
type NestedGroup struct {
	GroupID   string `json:"group_id"`
	GroupName string `json:"group_name"`
}

// UserGroup defines the structure for a single user group.
type UserGroup struct {
	GroupID        string        `json:"group_id"`
	GroupName      string        `json:"group_name"`
	Description    string        `json:"description"`
	RoleName       string        `json:"role_name"`
	PrettyRoleName string        `json:"pretty_role_name"`
	CreatedBy      string        `json:"created_by"`
	UpdatedBy      string        `json:"updated_by"`
	CreatedTS      int64         `json:"created_ts"`
	UpdatedTS      int64         `json:"updated_ts"`
	Users          []string      `json:"users"`
	GroupType      string        `json:"group_type"`
	NestedGroups   []NestedGroup `json:"nested_groups"`
	IDPGroups      []string      `json:"idp_groups"`
}
