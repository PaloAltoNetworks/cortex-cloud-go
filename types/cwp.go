package types

// ----------------------------------------------------------------------------
// CWP Policy
// ----------------------------------------------------------------------------

// Policy defines the structure for a CWP policy.
type Policy struct {
	Id                  string   `json:"id" tfsdk:"id"`
	Revision            int      `json:"revision" tfsdk:"revision"`
	CreatedAt           string   `json:"created_at" tfsdk:"created_at"`
	ModifiedAt          string   `json:"modified_at" tfsdk:"modified_at"`
	Type                string   `json:"type" tfsdk:"type"`
	CreatedBy           string   `json:"created_by" tfsdk:"created_by"`
	Disabled            bool     `json:"disabled" tfsdk:"disabled"`
	Name                string   `json:"name" tfsdk:"name"`
	Description         string   `json:"description" tfsdk:"description"`
	EvaluationModes     []string `json:"evaluation_modes" tfsdk:"evaluation_modes"`
	EvaluationStage     string   `json:"evaluation_stage" tfsdk:"evaluation_stage"`
	RulesIDs            []string `json:"rules_ids" tfsdk:"rules_ids"`
	Condition           string   `json:"condition" tfsdk:"condition"`
	Exception           string   `json:"exception" tfsdk:"exception"`
	AssetScope          string   `json:"asset_scope" tfsdk:"asset_scope"`
	AssetGroupIDs       []int    `json:"asset_group_ids" tfsdk:"asset_group_ids"`
	AssetGroups         []string `json:"asset_groups" tfsdk:"asset_groups"`
	PolicyAction        string   `json:"policy_action" tfsdk:"policy_action"`
	PolicySeverity      string   `json:"policy_severity" tfsdk:"policy_severity"`
	RemediationGuidance string   `json:"remediation_guidance" tfsdk:"remediation_guidance"`
}

// CreatePolicyRequest is the request for creating a CWP policy.
type CreatePolicyRequest struct {
	Type                string   `json:"type" tfsdk:"type"`
	Name                string   `json:"name" tfsdk:"name"`
	Description         string   `json:"description" tfsdk:"description"`
	EvaluationStage     string   `json:"evaluationStage" tfsdk:"evaluation_stage"`
	RulesIDs            []string `json:"rulesIds" tfsdk:"rules_ids"`
	AssetGroupIDs       []int    `json:"assetGroupsIDs" tfsdk:"asset_group_ids"`
	PolicyAction        string   `json:"action" tfsdk:"policy_action"`
	PolicySeverity      string   `json:"severity" tfsdk:"policy_severity"`
	RemediationGuidance string   `json:"remediationGuidance" tfsdk:"remediation_guidance"`
}

// CreatePolicyResponse is the response for creating a CWP policy.
type CreatePolicyResponse struct {
	Id string `json:"id"`
}

// GetPolicyByIDResponse is the response for getting a CWP policy by ID.
type GetPolicyByIDResponse struct {
	Policy
}

// ListPoliciesRequest is the request for listing CWP policies.
type ListPoliciesRequest struct {
	PolicyTypes []string `json:"policy_types,omitempty"`
}

// DeletePolicyRequest is the request for deleting a CWP policy.
type DeletePolicyRequest struct {
	Id          string `json:"id"`
	CloseIssues bool   `json:"close_issues,omitempty"`
}

// UpdatePolicyRequest is the request for updating a CWP policy.
type UpdatePolicyRequest struct {
	Policy
}
