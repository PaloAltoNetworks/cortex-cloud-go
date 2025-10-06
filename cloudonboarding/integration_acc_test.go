// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build acceptance

package cloudonboarding

//import (
//	"context"
//	"fmt"
//	//"os"
//	"strconv"
//	"testing"
//	"time"
//
//	"github.com/PaloAltoNetworks/cortex-cloud-go/enums"
//	acctest "github.com/PaloAltoNetworks/cortex-cloud-go/internal/test/acceptance"
//	//"github.com/aws/aws-sdk-go-v2/aws"
//	//"github.com/aws/aws-sdk-go-v2/config"
//	//"github.com/aws/aws-sdk-go-v2/service/cloudformation"
//	//"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
//	"github.com/stretchr/testify/assert"
//)
//
//
//const (
//	AWSIntegrationTemplateAutomatedLinkRegexp = `^https:\/\/([a-z0-9-]+\.)?console\.aws\.amazon\.com\/cloudformation\/home#\/stacks\/quickcreate\?templateURL=https%3A%2F%2F.+$`
//	AWSIntegrationTemplateManualLinkRegexp = `^\/.*\.ya?ml(\?.*)?$`
//)
//
//const (
//	GUIDRegexp = `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`
//	TrackingGUIDRegexp = `^[0-9a-fA-F]{32}$`
//)
//func setupAcceptanceTest(t *testing.T) *Client {
//	//apiUrl := os.Getenv("CORTEX_CLOUD_TEST_API_URL")
//	//apiKey := os.Getenv("CORTEX_CLOUD_TEST_API_KEY")
//	//apiKeyIDStr := os.Getenv("CORTEX_CLOUD_TEST_API_KEY_ID")
//	//apiKeyID, err := strconv.Atoi(apiKeyIDStr)
//	//if err != nil {
//	//	t.Fatalf("failed to convert API key ID \"%s\" to int: %s", apiKeyIDStr, err.Error())
//	//}
//
//	//config := &api.Config{
//	//	ApiUrl:   apiUrl,
//	//	ApiKey:   apiKey,
//	//	ApiKeyId: apiKeyID,
//	//}
//	//config, err := api.NewConfigFromFile("../susan-polgar-config.json", false)
//	//config, err := api.NewConfigFromFile("../pcs-lab-config.json", false)
//	config, err := api.NewConfigFromFile(acctest.ConfigFile, false)
//	if err != nil {
//		t.Fatalf("failed to read config file: %s", err.Error())
//	}
//
//	client, err := NewClient(config)
//	assert.NoError(t, err)
//	assert.NotNil(t, client)
//	return client
//}
//
//func TestAccAwsOrganizationIntegrationTemplateLifecycle(t *testing.T) {
//	client := setupAcceptanceTest(t)
//	ctx := context.Background()
//	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
//
//	// Test values
//	instanceName := fmt.Sprintf("tf-acctest-aws-account-%s", timestamp)
//	cloudProvider := enums.CloudProviderAWS.String()
//	//scope := enums.ScopeOrganization.String()
//	scope := enums.ScopeAccount.String()
//	scanMode := enums.ScanModeManaged.String()
//
//	additionalCapabilities := AdditionalCapabilities{
//		DataSecurityPostureManagement: true,
//		RegistryScanning:              true,
//		RegistryScanningOptions: RegistryScanningOptions{
//			Type: enums.RegistryScanningTypeAll.String(),
//		},
//		XsiamAnalytics: true,
//	}
//	collectionConfiguration := CollectionConfiguration{
//		AuditLogs: AuditLogsConfiguration{
//			Enabled: true,
//		},
//	}
//	customResourcesTags := []Tag{
//		{
//			Key:   "managed_by",
//			Value: "paloaltonetworks",
//		},
//		{
//			Key:   "test_tag",
//			Value: timestamp,
//		},
//	}
//	scopeModifications := ScopeModifications{
//		//Accounts: &ScopeModificationsOptionsGeneric{
//		//	Enabled: false,
//		//	//Enabled: true,
//		//	//Type: "INCLUDE",
//		//	//AccountIDs: []string{"123456789012"},
//		//},
//		Regions: &ScopeModificationsOptionsRegions{
//			//Enabled: false,
//			Enabled: true,
//			Type:    enums.ScopeModificationTypeInclude.String(),
//			Regions: []string{"us-east-1"},
//		},
//	}
//
//	// Execute create request
//	createReq := CreateIntegrationTemplateRequest{
//		Data: CreateIntegrationTemplateRequestData{
//			InstanceName:            instanceName,
//			CloudProvider:           cloudProvider,
//			Scope:                   scope,
//			ScanMode:                scanMode,
//			AdditionalCapabilities:  additionalCapabilities,
//			CollectionConfiguration: collectionConfiguration,
//			CustomResourcesTags:     customResourcesTags,
//			ScopeModifications:      scopeModifications,
//		},
//	}
//	createResp, err := client.CreateIntegrationTemplate(ctx, createReq)
//	if err != nil {
//		t.Fatalf("failed to create integration template: %s", err.Error())
//	}
//	createRespData := createResp.Reply
//
//	// Check response
//	assert.NotNil(t, createResp)
//	assert.NotNil(t, createRespData)
//	assert.NotNil(t, createRespData.Automated)
//	assert.Regexp(t, acctest.AWSIntegrationTemplateAutomatedLinkRegexp, createRespData.Automated.Link)
//	assert.Regexp(t, test.TrackingGUIDRegexp, createRespData.Automated.TrackingGuid)
//	// TODO: fix this
//	//assert.NotNil(t, response.Manual)
//	//assert.Regexp(t, test.AWSIntegrationTemplateManualLinkRegexp, response.Manual.CF)
//
//	// Execute get request
//	instanceID := createRespData.Automated.TrackingGuid
//	getReq := ListIntegrationInstancesRequest{
//		RequestData: ListIntegrationInstancesRequestData{
//			FilterData: FilterData{
//				Paging: PagingFilter{
//					From: 0,
//					To:   1000,
//				},
//				Filter: CriteriaFilter{
//					And: []Criteria{
//						{
//							SearchField: "ID",
//							SearchType:  "WILDCARD",
//							SearchValue: instanceID,
//						},
//					},
//				},
//			},
//		},
//	}
//	getResp, err := client.ListIntegrationInstances(ctx, getReq)
//	if err != nil {
//		t.Fatalf("failed to retrieve integration template: %s", err.Error())
//	}
//	assert.NotNil(t, getResp)
//	assert.NotNil(t, getResp.Reply)
//	assert.NotNil(t, getResp.Reply.Data)
//	assert.Len(t, getResp.Reply.Data, 1)
//	getRespData := getResp.Reply.Data[0]
//
//	// Check response
//	assert.NotNil(t, getRespData)
//	assert.Equal(t, instanceID, getRespData.InstanceID)
//	assert.Equal(t, cloudProvider, getRespData.CloudProvider)
//	assert.Equal(t, instanceName, getRespData.InstanceName)
//	assert.Equal(t, scope, getRespData.Scope)
//	assert.Equal(t, scanMode, getRespData.ScanMode)
//	assert.Equal(t, "PENDING", getRespData.Status)
//
//	marshalledGetResp, err := getResp.Marshal()
//	assert.NoError(t, err)
//	assert.NotNil(t, marshalledGetResp)
//	assert.Len(t, marshalledGetResp, 1)
//	marshalledGetRespData := marshalledGetResp[0]
//
//	assert.NotNil(t, marshalledGetRespData)
//	assert.NotNil(t, marshalledGetRespData.CustomResourcesTags)
//	assert.Equal(t, marshalledGetRespData.AdditionalCapabilities, additionalCapabilities)
//	assert.Equal(t, marshalledGetRespData.CustomResourcesTags, customResourcesTags)
//
//	// Deploy CloudFormation stack
//	templateURL, err := createResp.GetTemplateUrl()
//	if err != nil {
//		t.Fatalf("failed to parse CloudFormation template URL: %s", err.Error())
//	}
//	t.Logf("CloudFormation template URL: %s", templateURL)
//	acctest.DeployCloudFormationStack(t, ctx, "us-east-1", timestamp, templateURL)
//
//	//// Update test values
//	//updatedInstanceName := fmt.Sprintf("%s-updated", instanceName)
//	//// TODO: test Days Modified setting for registry scanning options
//	//updatedAdditionalCapabilities := AdditionalCapabilities{
//	//	DataSecurityPostureManagement: false,
//	//	RegistryScanning: false,
//	//	RegistryScanningOptions: RegistryScanningOptions{
//	//		Type: enums.RegistryScanningTypeLatestTag.String(),
//	//	},
//	//	XsiamAnalytics: false,
//	//}
//	//updatedCustomResourcesTags := []Tag{
//	//	{
//	//		Key: "managed_by",
//	//		Value: "paloaltonetworks",
//	//	},
//	//	{
//	//		Key: "test_tag_updated",
//	//		Value: fmt.Sprintf("%s-updated", timestamp),
//	//	},
//	//}
//	//updatedCollectionConfiguration := CollectionConfiguration{
//	//	AuditLogs: AuditLogsConfiguration{
//	//		Enabled: false,
//	//	},
//	//}
//	////scopeModifications := ScopeModifications{
//	////	Accounts: &ScopeModificationsOptionsGeneric{
//	////		Enabled: false,
//	////		//Enabled: true,
//	////		//Type: "INCLUDE",
//	////		//AccountIDs: []string{"123456789012"},
//	////	},
//	////	Regions: &ScopeModificationsOptionsRegions{
//	////		Enabled: false,
//	////		//Type: regionScopeType,
//	////		//Regions: regions,
//	////	},
//	////}
//	//
//	//// Execute update request
//	//updateReq := EditIntegrationInstanceRequest{
//	//	RequestData: EditIntegrationInstanceRequestData{
//	//		InstanceID: instanceID,
//	//		CloudProvider: cloudProvider,
//	//		// TODO: replace with logic to populate this
//	//		ScanEnvID: "43083abe03a648e7b029b9b1b5403b13",
//	//		InstanceName: updatedInstanceName,
//	//		AdditionalCapabilities: updatedAdditionalCapabilities,
//	//		CollectionConfiguration: updatedCollectionConfiguration,
//	//		CustomResourcesTags: updatedCustomResourcesTags,
//	//		// TODO: updated scope modifications
//	//		ScopeModifications: scopeModifications,
//	//	},
//	//}
//	//updateResp, err := client.EditIntegrationInstance(ctx, updateReq)
//	//assert.NoError(t, err)
//	//assert.NotNil(t, updateResp)
//	//assert.NotNil(t, updateResp.Reply)
//	//updateRespData := updateResp.Reply
//
//	//// Check response
//	//assert.NotNil(t, updateRespData)
//	//assert.NotNil(t, updateRespData.Automated)
//	//assert.Regexp(t, acctest.AWSIntegrationTemplateAutomatedLinkRegexp, updateRespData.Automated.Link)
//	//assert.Regexp(t, test.TrackingGUIDRegexp, updateRespData.Automated.TrackingGuid)
//	// TODO: fix this
//	//assert.NotNil(t, response.Manual)
//	//assert.Regexp(t, acctest.AWSIntegrationTemplateManualLinkRegexp, response.Manual.CF)
//
//	//// Check response
//	//assert.NotNil(t, updateRespData)
//	//assert.Equal(t, instanceID, updateRespData.InstanceID)
//	//assert.Equal(t, cloudProvider, updateRespData.CloudProvider)
//	//assert.Equal(t, instanceName, updateRespData.InstanceName)
//	//assert.Equal(t, scope, updateRespData.Scope)
//	//assert.Equal(t, scanMode, updateRespData.ScanMode)
//	//assert.Equal(t, "PENDING", updateRespData.Status)
//	//
//	//marshalledupdateResp, err := updateResp.Marshal()
//	//assert.NoError(t, err)
//	//assert.NotNil(t, marshalledupdateResp)
//	//assert.Len(t, marshalledupdateResp, 1)
//	//marshalledupdateRespData := marshalledupdateResp[0]
//
//	//assert.NotNil(t, marshalledGetRespData)
//	//assert.NotNil(t, marshalledGetRespData.CustomResourcesTags)
//	//assert.Equal(t, marshalledGetRespData.AdditionalCapabilities, additionalCapabilities)
//	//assert.Equal(t, marshalledGetRespData.CustomResourcesTags, customResourcesTags)
//}

//func TestAccAwsAccountIntegrationTemplateLifecycle(t *testing.T) {
//	client := setupAcceptanceTest(t)
//	ctx := context.Background()
//	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
//
//	// Create a new AWS account integration template
//	instanceName := fmt.Sprintf("tf-acctest-aws-account-%s", timestamp)
//	//regionScopeType := enums.ScopeModificationTypeInclude.String()
//	//regions := []string{ "us-east-1", "us-east-2" }
//
//	createReq := CreateIntegrationTemplateRequest{
//		Data: CreateIntegrationTemplateRequestData{
//			InstanceName: instanceName,
//			CloudProvider: enums.CloudProviderAWS.String(),
//			Scope: enums.ScopeOrganization.String(),
//			//Scope: enums.ScopeAccount.String(),
//			ScanMode: enums.ScanModeManaged.String(),
//			AdditionalCapabilities: AdditionalCapabilities{
//				DataSecurityPostureManagement: true,
//				RegistryScanning: true,
//				RegistryScanningOptions: RegistryScanningOptions{
//					Type: enums.RegistryScanningTypeAll.String(),
//				},
//				XsiamAnalytics: true,
//			},
//			CollectionConfiguration: CollectionConfiguration{
//				AuditLogs: AuditLogsConfiguration{
//					Enabled: true,
//				},
//			},
//			CustomResourcesTags: []Tag{
//				{
//					Key: "managed_by",
//					Value: "paloaltonetworks",
//				},
//				{
//					Key: "test_tag",
//					Value: timestamp,
//				},
//			},
//			ScopeModifications: ScopeModifications{
//				// TODO: Accounts
//				Accounts: &ScopeModificationsOptionsGeneric{
//					Enabled: false,
//					//Enabled: true,
//					//Type: "INCLUDE",
//					//AccountIDs: []string{"123456789012"},
//				},
//				Regions: &ScopeModificationsOptionsRegions{
//					Enabled: false,
//					//Type: regionScopeType,
//					//Regions: regions,
//				},
//			},
//		},
//	}
//
//	response, err := client.CreateIntegrationTemplate(ctx, createReq)
//	assert.NoError(t, err)
//	assert.NotNil(t, response)
//
//	responseData := response.Reply
//	assert.NotNil(t, responseData)
//	// TODO: check response fields
//	//assert.Equal(t, instanceName, createdTemplate.InstanceName)
//	//assert.Equal(t, enums.CloudProviderAWS.String(), createdTemplate.CloudProvider)
//
//	//// Get the template
//	//getReq := GetIntegrationInstanceRequest{
//
//	//}
//	//createdTemplate, err := client.GetInstanceDetails(ctx,
//	//assert.NoError(t, err)
//	//assert.NotNil(t, gotRule)
//	//assert.Equal(t, ruleID, gotRule.Id)
//	//assert.Equal(t, instanceName, gotRule.Name)
//
//	// Update the rule
//	//updatedName := fmt.Sprintf("updated-test-rule-%d", time.Now().Unix())
//	//updateReq := UpdateRequest{
//	//	Name: updatedName,
//	//}
//	//updatedResp, err := client.Update(ctx, ruleID, updateReq)
//	//assert.NoError(t, err)
//	//assert.NotNil(t, updatedResp)
//	//assert.Equal(t, updatedName, updatedResp.Rule.Name)
//
//	//// Delete the rule
//	//err = client.Delete(ctx, ruleID)
//	//assert.NoError(t, err)
//
//	//// Verify the rule is deleted
//	//_, err = client.Get(ctx, ruleID)
//	//assert.Error(t, err) // Expect an error when getting a deleted rule
//}

//func TestAppsecRuleLifecycle(t *testing.T) {
//	client := setupAcceptanceTest(t)
//	ctx := context.Background()
//
//	// Create a new rule
//	instanceName := fmt.Sprintf("test-rule-%d", time.Now().Unix())
//	createReq := CreateOrCloneRequest{
//		Name:        instanceName,
//		Description: "test description",
//		Category:    string(enums.IacCategoryCompute),
//		SubCategory: string(enums.IacSubCategoryComputeOverprovisioned),
//		Scanner:     string(enums.ScannerIAC),
//		Severity:    string(enums.SeverityLow),
//		Frameworks: []FrameworkData{
//			{
//				Name:       string(enums.FrameworkNameTerraform),
//				Definition: "scope:\n  provider: \"aws\"\ndefinition:\n  or:\n    - cond_type: \"attribute\"\n      resource_types:\n        - \"aws_instance\"\n      attribute: \"instance_type\"\n      operator: \"equals\"\n      value: \"t3.micro\"",
//			},
//		},
//		Labels: []string{"test-label"},
//	}
//	createdRule, err := client.CreateOrClone(ctx, createReq)
//	assert.NoError(t, err)
//	assert.NotNil(t, createdRule)
//	assert.Equal(t, instanceName, createdRule.Name)
//	ruleID := createdRule.Id
//
//	// Get the rule
//	gotRule, err := client.Get(ctx, ruleID)
//	assert.NoError(t, err)
//	assert.NotNil(t, gotRule)
//	assert.Equal(t, ruleID, gotRule.Id)
//	assert.Equal(t, instanceName, gotRule.Name)
//
//	// Update the rule
//	updatedName := fmt.Sprintf("updated-test-rule-%d", time.Now().Unix())
//	updateReq := UpdateRequest{
//		Name: updatedName,
//	}
//	updatedResp, err := client.Update(ctx, ruleID, updateReq)
//	assert.NoError(t, err)
//	assert.NotNil(t, updatedResp)
//	assert.Equal(t, updatedName, updatedResp.Rule.Name)
//
//	// Delete the rule
//	err = client.Delete(ctx, ruleID)
//	assert.NoError(t, err)
//
//	// Verify the rule is deleted
//	_, err = client.Get(ctx, ruleID)
//	assert.Error(t, err) // Expect an error when getting a deleted rule
//}
//
//func TestClient_List_Acceptance(t *testing.T) {
//	t.Skip("Skipping test due to persistent failures")
//	client := setupAcceptanceTest(t)
//	listReq := ListRequest{
//		Limit: 1,
//	}
//	resp, err := client.List(context.Background(), listReq)
//	assert.NoError(t, err)
//	assert.NotNil(t, resp)
//}

//func DeployCloudFormationStack(t *testing.T, ctx context.Context, region string, timestamp string, templateURL string) {
//	// Load default AWS configuration from environment
//	//cfnConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
//	cfnConfig, err := config.LoadDefaultConfig(ctx)
//	if err != nil {
//		t.Fatalf("unable to load AWS SDK config: %s", err.Error())
//	}
//
//	// Create CloudFormation client
//	cfnClient := cloudformation.NewFromConfig(cfnConfig)
//
//	// Deploy CloudFormation stack
//	stackName := fmt.Sprintf("go-sdk-acctest-%s", timestamp)
//	createStackInput := &cloudformation.CreateStackInput{
//		StackName: aws.String(stackName),
//		TemplateURL: aws.String(templateURL),
//		// Specify OrganizationalUnitId for Organization deployment
//		//Parameters:
//		Capabilities: []types.Capability{
//			types.CapabilityCapabilityIam,
//			types.CapabilityCapabilityNamedIam,
//		},
//	}
//
//	startTime := time.Now()
//	t.Logf("Creating CloudFormation stack: %s", stackName)
//	_, err = cfnClient.CreateStack(ctx, createStackInput)
//	if err != nil {
//		t.Fatalf("failed to create CloudFormation stack: %v", err)
//	}
//
//	// Register cleanup func
//	t.Cleanup(func() {
//		t.Logf("Tearing down CloudFormation stack: %s", stackName)
//		deleteInput := &cloudformation.DeleteStackInput{
//			StackName: aws.String(stackName),
//		}
//		_, err := cfnClient.DeleteStack(ctx, deleteInput)
//		if err != nil {
//			t.Logf("failed to initiate CloudFormation stack deletion, you may need to delete it manually: %v", err)
//		}
//
//		// Optional: Wait for deletion to complete
//		deleteWaiter := cloudformation.NewStackDeleteCompleteWaiter(cfnClient)
//		err = deleteWaiter.Wait(ctx, &cloudformation.DescribeStacksInput{StackName: aws.String(stackName)}, 5*time.Minute)
//		if err != nil {
//			t.Logf("error waiting for CloudFormation stack deletion: %v", err)
//		}
//		t.Logf("Stack %s deleted successfully", stackName)
//	})
//
//	// Wait for the stack to be created successfully
//	waiter := cloudformation.NewStackCreateCompleteWaiter(cfnClient)
//	describeInput := &cloudformation.DescribeStacksInput{StackName: aws.String(stackName)}
//
//	// Wait for up to 10 minutes for create operation to finish
//	err = waiter.Wait(ctx, describeInput, 10*time.Minute)
//	if err != nil {
//		t.Fatalf("timed out or failed waiting for stack creation: %v", err)
//	}
//	endTime := time.Now()
//	deployTime := endTime.Sub(startTime)
//
//	t.Logf("Stack %s successfully created in %s", deployTime.String())
//}
