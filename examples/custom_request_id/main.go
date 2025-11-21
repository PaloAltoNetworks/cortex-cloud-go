// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/client"
	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/config"
	"github.com/PaloAltoNetworks/cortex-cloud-go/platform"
)

func main() {
	// Create client
	platformClient, err := platform.NewClient(
		config.WithCortexAPIURL("https://api-tenant.xdr.us.paloaltonetworks.com"),
		config.WithCortexAPIKey("your-api-key"),
		config.WithCortexAPIKeyID(1),
		config.WithCortexAPIKeyType("standard"),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Example 1: Custom request ID
	ctx := context.Background()
	ctx = client.WithRequestID(ctx, "my-workflow-12345")

	groups, err := platformClient.ListAssetGroups(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to list asset groups: %v", err)
	}

	fmt.Printf("Found %d asset groups (request ID: my-workflow-12345)\n", len(groups))

	// Example 2: Get or generate request ID
	ctx2 := context.Background()
	ctx2, requestID := client.GetOrGenerateRequestID(ctx2)

	fmt.Printf("Using request ID: %s\n", requestID)

	groups2, err := platformClient.ListAssetGroups(ctx2, nil)
	if err != nil {
		log.Fatalf("Failed to list asset groups: %v", err)
	}

	fmt.Printf("Found %d asset groups\n", len(groups2))

	// Example 3: Request ID propagation through functions
	if err := processWorkflow(context.Background(), platformClient); err != nil {
		log.Fatalf("Workflow failed: %v", err)
	}
}

func processWorkflow(ctx context.Context, client *platform.Client) error {
	// Get or generate request ID for the entire workflow
	ctx, requestID := client.GetOrGenerateRequestID(ctx)
	log.Printf("[%s] Starting workflow", requestID)

	// All API calls in this workflow share the same request ID
	if err := step1(ctx, client); err != nil {
		return fmt.Errorf("[%s] step1 failed: %w", requestID, err)
	}

	if err := step2(ctx, client); err != nil {
		return fmt.Errorf("[%s] step2 failed: %w", requestID, err)
	}

	log.Printf("[%s] Workflow completed", requestID)
	return nil
}

func step1(ctx context.Context, client *platform.Client) error {
	requestID := client.GetRequestID(ctx)
	log.Printf("[%s] Executing step1", requestID)

	groups, err := client.ListAssetGroups(ctx, nil)
	if err != nil {
		return err
	}

	log.Printf("[%s] Step1 found %d groups", requestID, len(groups))
	return nil
}

func step2(ctx context.Context, client *platform.Client) error {
	requestID := client.GetRequestID(ctx)
	log.Printf("[%s] Executing step2", requestID)

	// Additional processing here

	log.Printf("[%s] Step2 completed", requestID)
	return nil
}
