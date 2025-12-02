package types

// CloudWorkloadPolicy defines a Cloud Workload Policy.
//
// Cloud Workload Policies leverage actionable findings or enforce preventive measures at defined stages of the Software Development Life Cycle (SDLC).
//
// Required license: Cortex Cloud Runtime Security or Cortex Cloud Posture Management
type CloudWorkloadPolicy struct {
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
	RulesIDs            []string `json:"rulesIds"`
	Condition           string   `json:"condition"`
	Exception           string   `json:"exception"`
	AssetScope          string   `json:"assetScope"`
	AssetGroupIDs       []int    `json:"assetGroupsIDs"`
	AssetGroups         []string `json:"assetGroups"`
	PolicyAction        string   `json:"action"`
	PolicySeverity      string   `json:"severity"`
	RemediationGuidance string   `json:"remediationGuidance"`
}

// CreateCloudWorkloadPolicyResponse is the request body for the Create Cloud Workload Policy endpoint.
type CreateCloudWorkloadPolicyRequest struct {
	Type                string   `json:"type"`
	Name                string   `json:"name"`
	Description         string   `json:"description"`
	EvaluationModes     []string `json:"evaluationModes,omitempty"`
	EvaluationStage     string   `json:"evaluationStage"`
	RulesIDs            []string `json:"rulesIds,omitempty"`
	Condition           string   `json:"condition,omitempty"`
	Exception           string   `json:"exception,omitempty"`
	AssetScope          string   `json:"assetScope,omitempty"`
	AssetGroupIDs       []int    `json:"assetGroupsIDs,omitempty"`
	AssetGroups         []string `json:"assetGroups,omitempty"`
	PolicyAction        string   `json:"action"`
	PolicySeverity      string   `json:"severity"`
	RemediationGuidance string   `json:"remediationGuidance,omitempty"`
}

// CreateCloudWorkloadPolicyResponse is the response body for the Create Cloud Workload Policy endpoint.
type CreateCloudWorkloadPolicyResponse struct {
	ID string `json:"id"`
}

// ListCloudWorkloadPoliciesRequest is the request body for the List Cloud Workload Policies endpoint.
type ListCloudWorkloadPoliciesRequest struct {
	PolicyTypes []string `json:"policy_types,omitempty"`
}

// DeleteCloudWorkloadPolicyRequest is the request body for the Delete Cloud Workload Policy endpoint.
type DeleteCloudWorkloadPolicyRequest struct {
	ID          string `json:"id"`
	CloseIssues bool   `json:"close_issues,omitempty"`
}

// UpdateCloudWorkloadPolicyRequest is the request body for the Update Cloud Workload Policy Endpoint
type UpdateCloudWorkloadPolicyRequest struct {
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
	RulesIDs            []string `json:"rulesIds"`
	Condition           string   `json:"condition"`
	Exception           string   `json:"exception"`
	AssetScope          string   `json:"assetScope"`
	AssetGroupIDs       []int    `json:"assetGroupsIDs"`
	AssetGroups         []string `json:"assetGroups"`
	Action              string   `json:"action"`
	Severity            string   `json:"severity"`
	RemediationGuidance string   `json:"remediationGuidance"`
}
