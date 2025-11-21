// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/config"
	"github.com/PaloAltoNetworks/cortex-cloud-go/platform"
	"github.com/PaloAltoNetworks/cortex-cloud-go/version"
)

func main() {
	// Example 1: Custom User-Agent
	client1, err := platform.NewClient(
		config.WithAgent("my-security-app/2.0.0"),
		config.WithCortexAPIURL("https://api-tenant.xdr.us.paloaltonetworks.com"),
		config.WithCortexAPIKey("your-api-key"),
		config.WithCortexAPIKeyID(1),
		config.WithCortexAPIKeyType("standard"),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	groups, err := client1.ListAssetGroups(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to list asset groups: %v", err)
	}

	fmt.Printf("Example 1: Found %d asset groups with custom User-Agent\n", len(groups))

	// Example 2: Append to default User-Agent
	customUA := version.UserAgentWithCustom(
		platform.ModuleName,
		platform.Version,
		"my-app/2.0.0",
	)

	client2, err := platform.NewClient(
		config.WithAgent(customUA),
		config.WithCortexAPIURL("https://api-tenant.xdr.us.paloaltonetworks.com"),
		config.WithCortexAPIKey("your-api-key"),
		config.WithCortexAPIKeyID(1),
		config.WithCortexAPIKeyType("standard"),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	groups2, err := client2.ListAssetGroups(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to list asset groups: %v", err)
	}

	fmt.Printf("Example 2: Found %d asset groups with appended User-Agent\n", len(groups2))
	fmt.Printf("User-Agent: %s\n", customUA)

	// Example 3: Get version information
	info := version.Info()
	fmt.Println("\nSDK Version Information:")
	for key, value := range info {
		fmt.Printf("  %s: %s\n", key, value)
	}
}
