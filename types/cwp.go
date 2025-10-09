package types

// ----------------------------------------------------------------------------
// CWP Policy
// ----------------------------------------------------------------------------

// Policy defines the structure for a CWP policy.
type Policy struct {
	Id                  string   `json:"id"`
	Revision            int      `json:"revision"`
	CreatedAt           string   `json:"created_at"`
	ModifiedAt          string   `json:"modified_at"`
	Type                string   `json:"type"`
	CreatedBy           string   `json:"created_by"`
	Disabled            bool     `json:"disabled"`
	Name                string   `json:"name"`
	Description         string   `json:"description"`
	EvaluationModes     []string `json:"evaluation_modes"`
	EvaluationStage     string   `json:"evaluation_stage"`
	RulesIDs            []string `json:"rules_ids"`
	Condition           string   `json:"condition"`
	Exception           string   `json:"exception"`
	AssetScope          string   `json:"asset_scope"`
	AssetGroupIDs       []int    `json:"asset_group_ids"`
	AssetGroups         []string `json:"asset_groups"`
	PolicyAction        string   `json:"policy_action"`
	PolicySeverity      string   `json:"policy_severity"`
	RemediationGuidance string   `json:"remediation_guidance"`
}

// CreatePolicyRequest is the request for creating a CWP policy.
type CreatePolicyRequest struct {
	Type                string   `json:"type"`
	Name                string   `json:"name"`
	Description         string   `json:"description"`
	EvaluationStage     string   `json:"evaluationStage"`
	RulesIDs            []string `json:"rulesIds"`
	AssetGroupIDs       []int    `json:"assetGroupsIDs"`
	PolicyAction        string   `json:"action"`
	PolicySeverity      string   `json:"severity"`
	RemediationGuidance string   `json:"remediationGuidance"`
}

// CreatePolicyResponse is the response for creating a CWP policy.
type CreatePolicyResponse struct {
	Id string `json:"id"`
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
