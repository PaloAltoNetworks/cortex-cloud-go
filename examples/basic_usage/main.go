// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/config"
	"github.com/PaloAltoNetworks/cortex-cloud-go/platform"
)

func main() {
	// Create client with environment variables
	// Set: CORTEX_API_URL, CORTEX_API_KEY, CORTEX_API_KEY_ID, CORTEX_API_KEY_TYPE
	client, err := platform.NewClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Or create client with explicit configuration
	client, err = platform.NewClient(
		config.WithCortexAPIURL("https://api-tenant.xdr.us.paloaltonetworks.com"),
		config.WithCortexAPIKey("your-api-key"),
		config.WithCortexAPIKeyID(1),
		config.WithCortexAPIKeyType("standard"),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// List asset groups
	ctx := context.Background()
	groups, err := client.ListAssetGroups(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to list asset groups: %v", err)
	}

	fmt.Printf("Found %d asset groups\n", len(groups))
	for _, group := range groups {
		fmt.Printf("  - %s: %s\n", group.ID, group.Name)
	}
}
