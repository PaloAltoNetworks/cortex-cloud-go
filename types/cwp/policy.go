package types

// Policy defines the structure for a CWP policy.
type Policy struct {
	ID                  string   `json:"id"`
	Revision            int      `json:"revision"`
	CreatedAt           string   `json:"createdAt"`
	ModifiedAt          string   `json:"modifiedAt"`
	Type                string   `json:"type"`
	CreatedBy           string   `json:"createdBy"`
	Disabled            bool     `json:"disabled"`
	Name                string   `json:"name"`
	Description         string   `json:"description"`
	EvaluationModes     []string `json:"evaluationModes"`
	EvaluationStage     string   `json:"evaluationStage"`
	RulesIDs            []string `json:"rulesIds"` // lowercase "i"
	Condition           string   `json:"condition"`
	Exception           string   `json:"exception"`
	AssetScope          string   `json:"assetScope"`
	AssetGroupIDs       []int    `json:"assetGroupsIDs"` // "s" before IDs
	AssetGroups         []string `json:"assetGroups"`
	PolicyAction        string   `json:"action"`   // not "policy_action"
	PolicySeverity      string   `json:"severity"` // not "policy_severity"
	RemediationGuidance string   `json:"remediationGuidance"`
}

// CreatePolicyRequest is the request for creating a CWP policy.
type CreatePolicyRequest struct {
	Type                string   `json:"type"`
	Name                string   `json:"name"`
	Description         string   `json:"description"`
	EvaluationModes     []string `json:"evaluationModes,omitempty"`
	EvaluationStage     string   `json:"evaluationStage"`
	RulesIDs            []string `json:"rulesIds,omitempty"` // note lowercase "i"
	Condition           string   `json:"condition,omitempty"`
	Exception           string   `json:"exception,omitempty"`
	AssetScope          string   `json:"assetScope,omitempty"`
	AssetGroupIDs       []int    `json:"assetGroupsIDs,omitempty"` // note "s" before "IDs"
	AssetGroups         []string `json:"assetGroups,omitempty"`
	PolicyAction        string   `json:"action"`   // maps to "action"
	PolicySeverity      string   `json:"severity"` // maps to "severity"
	RemediationGuidance string   `json:"remediationGuidance,omitempty"`
}

// CreatedModifiedAt represents the datetime value that the rule was created
// or updated.
type CreatedModifiedAt struct {
	Value string `json:"value,omitempty"`
}

// CreatePolicyResponse is the response for creating a CWP policy.
type CreatePolicyResponse struct {
	ID string `json:"id"`
}

// ListPoliciesRequest is the request for listing CWP policies.
type ListPoliciesRequest struct {
	PolicyTypes []string `json:"policy_types,omitempty"`
}

// DeletePolicyRequest is the request for deleting a CWP policy.
type DeletePolicyRequest struct {
	ID          string `json:"id"`
	CloseIssues bool   `json:"close_issues,omitempty"`
}

// UpdatePolicyRequest is the request for updating a CWP policy.
type UpdatePolicyRequest struct {
	Policy
}
