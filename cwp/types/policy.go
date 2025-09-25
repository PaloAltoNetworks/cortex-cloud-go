package types

import (
	////"context"
	//"encoding/json"
	////"net/http"
	//"net/url"
	//"strconv"

	//"dario.cat/mergo"
	//"github.com/PaloAltoNetworks/cortex-cloud-go/enums"
)

// ----------------------------
// Shared Structs
// ----------------------------

type Policy struct {
	Id                  string                 `json:"id" tfsdk:"id"`
	Revision            int                    `json:"revision" tfsdk:"revision"`
	CreatedAt           string                 `json:"created_at" tfsdk:"created_at"`
	ModifiedAt          string                 `json:"modified_at" tfsdk:"modified_at"`
	Type                string `json:"type" tfsdk:"type"`
	CreatedBy           string                 `json:"created_by" tfsdk:"created_by"`
	Disabled            bool                   `json:"disabled" tfsdk:"disabled"`
	Name                string                 `json:"name" tfsdk:"name"`
	Description         string                 `json:"description" tfsdk:"description"`
	EvaluationModes     []string `json:"evaluation_modes" tfsdk:"evaluation_modes"`
	EvaluationStage     string                 `json:"evaluation_stage" tfsdk:"evaluation_stage"`
	RulesIDs            []string               `json:"rules_ids" tfsdk:"rules_ids"`
	Condition           string                 `json:"condition" tfsdk:"condition"`
	Exception           string                 `json:"exception" tfsdk:"exception"`
	AssetScope          string                 `json:"asset_scope" tfsdk:"asset_scope"`
	AssetGroupIDs       []int                  `json:"asset_group_ids" tfsdk:"asset_group_ids"`
	AssetGroups         []string               `json:"asset_groups" tfsdk:"asset_groups"`
	PolicyAction        string `json:"policy_action" tfsdk:"policy_action"`
	PolicySeverity      string `json:"policy_severity" tfsdk:"policy_severity"`
	RemediationGuidance string                 `json:"remediation_guidance" tfsdk:"remediation_guidance"`
}
// ---------------------------
// Request/Response structs
// ---------------------------

// Add CWP policies

type CreatePolicyRequest struct {
	//Data Policy `json:"data"`
	//Policy
	//Id                  string                 `json:"id" tfsdk:"id"`
	//Revision            int                    `json:"revision" tfsdk:"revision"`
	//CreatedAt           string                 `json:"created_at" tfsdk:"created_at"`
	//ModifiedAt          string                 `json:"modified_at" tfsdk:"modified_at"`
	//Type                enums.PolicyType       `json:"type" tfsdk:"type"`
	Type                string `json:"type" tfsdk:"type"`
	//CreatedBy           string                 `json:"created_by" tfsdk:"created_by"`
	//Disabled            bool                   `json:"disabled" tfsdk:"disabled"`
	Name                string                 `json:"name" tfsdk:"name"`
	Description         string                 `json:"description" tfsdk:"description"`
	//EvaluationModes     []enums.EvaluationMode `json:"evaluation_modes" tfsdk:"evaluation_modes"`
	//EvaluationModes     []string `json:"evaluation_modes" tfsdk:"evaluation_modes"`
	EvaluationStage     string                 `json:"evaluationStage" tfsdk:"evaluation_stage"`
	RulesIDs            []string               `json:"rulesIds" tfsdk:"rules_ids"`
	//Condition           string                 `json:"condition" tfsdk:"condition"`
	//Condition           []byte `json:"condition" tfsdk:"condition"`
	//Exception           string                 `json:"exception" tfsdk:"exception"`
	//Exception           []byte `json:"exception" tfsdk:"exception"`
	//AssetScope          string                 `json:"asset_scope" tfsdk:"asset_scope"`
	//AssetScope          []byte `json:"asset_scope" tfsdk:"asset_scope"`
	AssetGroupIDs       []int                  `json:"assetGroupsIDs" tfsdk:"asset_group_ids"`
	//AssetGroups         []string               `json:"asset_groups" tfsdk:"asset_groups"`
	//PolicyAction        enums.PolicyAction     `json:"policy_action" tfsdk:"policy_action"`
	PolicyAction        string `json:"action" tfsdk:"policy_action"`
	//PolicySeverity      enums.PolicySeverity   `json:"policy_severity" tfsdk:"policy_severity"`
	PolicySeverity      string `json:"severity" tfsdk:"policy_severity"`
	RemediationGuidance string                 `json:"remediationGuidance" tfsdk:"remediation_guidance"`
}

type CreatePolicyResponse struct {
	Id string `json:"id"`
}

// Get a CWP policy by ID

//type GetPolicyByIDRequest struct {
//	Id string `json:"id"`
//}

type GetPolicyByIDResponse struct {
	//Policy Policy `json:"policy"`
	Policy
}

// Get CWP policies

type ListPoliciesRequest struct {
	//PolicyTypes []enums.PolicyType `json:"policy_types,omitempty"`
	PolicyTypes []string `json:"policy_types,omitempty"`
}

//func (r ListPoliciesRequest) ToQueryValues() url.Values {
//	result := url.Values{}
//
//	for _, policyType := range r.PolicyTypes {
//		result.Add("policy_types", policyType)
//	}
//
//	return result
//}

//type ListPoliciesResponse struct {
//	Policies []PolicyData `json:"policies"`
//}

//func (r ListPoliciesResponse) Marshal() ([]PolicyData, error) {
//
//	marshalledResponse := []PolicyData{}
//	for _, pd := range r.Policies {
//		var condition string
//		var exception string
//		var assetScope string
//		err := json.Unmarshal([]byte(pd.Condition), &condition)
//		if err != nil {
//			return []PolicyData{}, err
//		}
//
//		err = json.Unmarshal([]byte(pd.Exception), &exception)
//		if err != nil {
//			return []PolicyData{}, err
//		}
//		err = json.Unmarshal([]byte(pd.AssetScope), &assetScope)
//		if err != nil {
//			return []PolicyData{}, err
//		}
//
//		marshalledResponseItem := PolicyData{
//			Id:                  pd.Id,
//			Revision:            pd.Revision,
//			CreatedAt:           pd.CreatedAt,
//			ModifiedAt:          pd.ModifiedAt,
//			Type:                pd.Type,
//			CreatedBy:           pd.CreatedBy,
//			Disabled:            pd.Disabled,
//			Name:                pd.Name,
//			Description:         pd.Description,
//			EvaluationModes:     pd.EvaluationModes,
//			EvaluationStage:     pd.EvaluationStage,
//			RulesIDs:            pd.RulesIDs,
//			Condition:           condition,
//			Exception:           exception,
//			AssetScope:          assetScope,
//			AssetGroupIDs:       pd.AssetGroupIDs,
//			AssetGroups:         pd.AssetGroups,
//			PolicyAction:        pd.PolicyAction,
//			PolicySeverity:      pd.PolicySeverity,
//			RemediationGuidance: pd.RemediationGuidance,
//		}
//
//		marshalledResponse = append(marshalledResponse, marshalledResponseItem)
//	}
//
//	return marshalledResponse, nil
//}

//func (r GetPolicyByIDResponse) Marshal() (PolicyData, error) {
//	var condition string
//	var exception string
//	var assetScope string
//
//	err := json.Unmarshal([]byte(r.Policy.Condition), &condition)
//	if err != nil {
//		return PolicyData{}, err
//	}
//
//	err = json.Unmarshal([]byte(r.Policy.Exception), &exception)
//	if err != nil {
//		return PolicyData{}, err
//	}
//
//	err = json.Unmarshal([]byte(r.Policy.AssetScope), &assetScope)
//	if err != nil {
//		return PolicyData{}, err
//	}
//
//	marshalledPolicy := PolicyData{
//		Id:                  r.Policy.Id,
//		Revision:            r.Policy.Revision,
//		CreatedAt:           r.Policy.CreatedAt,
//		ModifiedAt:          r.Policy.ModifiedAt,
//		Type:                r.Policy.Type,
//		CreatedBy:           r.Policy.CreatedBy,
//		Disabled:            r.Policy.Disabled,
//		Name:                r.Policy.Name,
//		Description:         r.Policy.Description,
//		EvaluationModes:     r.Policy.EvaluationModes,
//		EvaluationStage:     r.Policy.EvaluationStage,
//		RulesIDs:            r.Policy.RulesIDs,
//		Condition:           condition,
//		Exception:           exception,
//		AssetScope:          assetScope,
//		AssetGroupIDs:       r.Policy.AssetGroupIDs,
//		AssetGroups:         r.Policy.AssetGroups,
//		PolicyAction:        r.Policy.PolicyAction,
//		PolicySeverity:      r.Policy.PolicySeverity,
//		RemediationGuidance: r.Policy.RemediationGuidance,
//	}
//
//	return marshalledPolicy, nil
//}

type DeletePolicyRequest struct {
	Id          string `json:"id"`
	CloseIssues bool   `json:"close_issues,omitempty"`
}

//func (r DeletePolicyRequest) ToQueryValues() url.Values {
//	result := url.Values{}
//
//	if r.CloseIssues {
//		result.Add("close_issues", strconv.FormatBool(r.CloseIssues))
//	}
//
//	return result
//}

type UpdatePolicyRequest struct {
	//Id   string     `json:"id"`
	//Data Policy `json:"data"`
	Policy
}

//func (r PolicyData) ToUpdateRequest() UpdatePolicyRequest {
//	return UpdatePolicyRequest{
//		Id: r.Id,
//		Data: PolicyData{
//			Name:                r.Name,
//			Description:         r.Description,
//			Type:                r.Type,
//			EvaluationModes:     r.EvaluationModes,
//			EvaluationStage:     r.EvaluationStage,
//			RulesIDs:            r.RulesIDs,
//			Condition:           r.Condition,
//			Exception:           r.Exception,
//			AssetScope:          r.AssetScope,
//			AssetGroupIDs:       r.AssetGroupIDs,
//			PolicyAction:        r.PolicyAction,
//			PolicySeverity:      r.PolicySeverity,
//			RemediationGuidance: r.RemediationGuidance,
//		},
//	}
//}
