# Cortex Cloud Go SDK

Go SDK for [Cortex Cloud](https://www.paloaltonetworks.com/cortex/cloud), the industry-leading unified cloud security platform by Palo Alto Networks&reg;.

## Getting Started

### Prerequisites

Along with access to your Cortex Cloud tenant, you will need:
* An API key and the corresponding API key ID
* Your Cortex Cloud tenant's API URL

Refer to the [official API documentation](https://docs-cortex.paloaltonetworks.com/r/Cortex-Cloud-Platform-APIs/Create-a-new-API-key) for guidance on both of these.

It is recommended that you use an Advanced API key with the Cortex Cloud Go SDK, but either type will work as long as you specify the correct key type (`"standard"` or `"advanced"`) in the client configuration.

### Installation

* Create a module file in your project directory by running `go mod init` if there is not already a `go.mod` file present.
* Install your desired Cortex Cloud Go SDK package by running the following command (replacing `platform` with your desired package):
```bash
go get github.com/PaloAltoNetworks/cortex-cloud-go/platform@latest
```
* Import the package(s) into your project with `import "github.com/paloaltonetworks/cortex-cloud-go/platform"`

### Example Usage

This example shows how to create a new role by making a request to the [Create a new role](https://docs-cortex.paloaltonetworks.com/r/Cortex-Cloud-Platform-APIs/Create-a-new-role) endpoint:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/PaloAltoNetworks/cortex-cloud-go/platform"
	platformTypes "github.com/PaloAltoNetworks/cortex-cloud-go/types/platform"
)

func main() {
	// Create an instance of the Cortex Cloud Platform API client
	platformClient, err := platform.NewClient(
		platform.WithCortexAPIURL("https://{tenant-name}.xdr.{region}.paloaltonetworks.com"),
		platform.WithCortexAPIKey("{your-api-key}"),
		platform.WithCortexAPIKeyID(100),
		platform.WithCortexAPIKeyType("advanced"), // Defaults to "advanced" if not specified
	)

	// Build the request input
	input := platformTypes.RoleCreateRequest{
		RequestData: platformTypes.RoleCreateRequestData{
			PrettyName: "ExampleRole",
			ComponentPermissions: []string{"rules_action"},
		},
	}
	
	// Execute the request
	ctx := context.Background()
	resp, err := platformClient.CreateRole(ctx, input)
	if err != nil {
		log.Fatalf("failed to create role: %s", err.Error())
	}

	// Handle the response
	fmt.Printf("new role ID: %s", resp.RoleID)
}
```

## Resources

* [Cortex Cloud API Documentation](https://docs-cortex.paloaltonetworks.com/r/Cortex-Cloud-Platform-APIs/Create-a-new-API-key)
* [Cortex Cloud Documentation Portal](https://docs-cortex.paloaltonetworks.com/p/Cortex+CLOUD)
