// Copyright (c) Palo Alto Networks, Inc.
// SPDX-License-Identifier: MPL-2.0

package cortexcloud

/*

Package cortexcloud is the root of the packages used to access the Cortex Cloud
API.

TODO: add pkg.go.dev link to list of sub-modules

# Client Configuration

All clients are configured by passing the Config struct from the
cortex-cloud-go/api package with the appropriate values. See the below section
for examples.

# Authentication

All clients support authentication via the basic API key.

TODO: add link/description for API key creation in Cortex UI

The following example authenticates to the Cloud Onboarding Public API and
returns the respective API client struct:

TODO: update examples and add to examples/ dir

	import (
		"github.com/PaloAltoNetworks/cortex-cloud-go/api"
		"github.com/PaloAltoNetworks/cortex-cloud-go/cloudonboarding"
	)

	func main() {
		config := api.Config{
			ApiUrl:   "https://api-cortexcloud.xdr.us.paloaltonetworks.com/",
			ApiKey:   "api-key-value-here",
			ApiKeyId: 123,
		}

		cloudOnboardingClient, err := cloudonboarding.NewClient(config)
		if err != nil {
			log.Fatal(err)
		}
	}

*/
