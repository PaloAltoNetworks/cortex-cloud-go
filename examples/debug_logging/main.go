// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"log"

	"github.com/PaloAltoNetworks/cortex-cloud-go/internal/config"
	cortexlog "github.com/PaloAltoNetworks/cortex-cloud-go/log"
	"github.com/PaloAltoNetworks/cortex-cloud-go/platform"
)

func main() {
	// Create logger with debug level
	logger := cortexlog.NewDefaultLogger()

	// Create client with debug logging enabled
	client, err := platform.NewClient(
		config.WithLogger(logger),
		config.WithLogLevel("debug"),
		config.WithCortexAPIURL("https://api-tenant.xdr.us.paloaltonetworks.com"),
		config.WithCortexAPIKey("your-api-key"),
		config.WithCortexAPIKeyID(1),
		config.WithCortexAPIKeyType("standard"),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Make API call - will see detailed request/response logs with request ID
	ctx := context.Background()
	groups, err := client.ListAssetGroups(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to list asset groups: %v", err)
	}

	log.Printf("Found %d asset groups", len(groups))

	// Debug logs will show:
	// ---[ REQUEST req_abc123 ]-----------------------------
	// GET /asset_groups HTTP/1.1
	// Host: api-tenant.xdr.us.paloaltonetworks.com
	// User-Agent: cortex-cloud-go/1.0.0 (platform/1.0.0; go1.25.1; darwin/arm64)
	// X-Request-ID: req_abc123
	// ...
	// -----------------------------------------------------
	//
	// ---[ RESPONSE req_abc123 ]----------------------------
	// HTTP/1.1 200 OK
	// Content-Type: application/json
	// ...
	// -----------------------------------------------------
}
