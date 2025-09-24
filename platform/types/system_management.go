package types


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

// ---------------------------
// Request/Response structs
// ---------------------------

// Get User

type GetUserRequest struct {
	Email string `json:"email"`
}

// List Roles

type ListRolesRequestData struct {
	// TODO: add validation tag/function for role names?
	RoleNames []string `json:"role_names" validate:"required,min=1"`
}

type ListRolesResponseReply struct {
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

// Set Role

type SetRoleRequestData struct {
	UserEmails []string `json:"user_emails" validate:"required,min=1,dive,required,email"`
	RoleName   string   `json:"role_name"`
}

type SetRoleResponseReply struct {
	UpdateCount string `json:"update_count"`
}

// Get Risk Score

type GetRiskScoreRequestData struct {
	ID string `json:"id" validate:"required,sysmgmtID"`
}

type GetRiskScoreResponseReply struct {
	Type          string   `json:"type"`
	ID            string   `json:"id"`
	Score         int      `json:"score"`
	NormRiskScore int      `json:"norm_risk_score"`
	RiskLevel     string   `json:"risk_level"`
	Reasons       []Reason `json:"reasons"`
	Email         string   `json:"email"`
}

// List Risky Users

type ListRiskyUsersResponseReply struct {
	Type          string   `json:"type"`
	ID            string   `json:"id"`
	Score         int      `json:"score"`
	NormRiskScore int      `json:"norm_risk_score"`
	RiskLevel     string   `json:"risk_level"`
	Reasons       []Reason `json:"reasons"`
	Email         string   `json:"email"`
}

// List Risky Hosts

type ListRiskyHostsResponseReply struct {
	Type          string   `json:"type"`
	ID            string   `json:"id"`
	Score         int      `json:"score"`
	NormRiskScore int      `json:"norm_risk_score"`
	RiskLevel     string   `json:"risk_level"`
	Reasons       []Reason `json:"reasons"`
}
