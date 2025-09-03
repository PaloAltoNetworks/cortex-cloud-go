module github.com/PaloAltoNetworks/cortex-cloud-go/cloudonboarding

go 1.24.3

require (
	github.com/PaloAltoNetworks/cortex-cloud-go/api v0.0.3-beta
	github.com/PaloAltoNetworks/cortex-cloud-go/internal/app v0.0.3-beta3
)

require (
	github.com/PaloAltoNetworks/cortex-cloud-go v0.0.3-beta // indirect
	github.com/PaloAltoNetworks/cortex-cloud-go/errors v0.0.3-beta // indirect
	github.com/PaloAltoNetworks/cortex-cloud-go/log v0.0.3-beta // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/hashicorp/go-hclog v1.5.0 // indirect
	github.com/hashicorp/terraform-plugin-log v0.9.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.3-beta14 // indirect
	golang.org/x/sys v0.30.0 // indirect
)

replace github.com/PaloAltoNetworks/cortex-cloud-go/api => ../api

replace github.com/PaloAltoNetworks/cortex-cloud-go/internal/app => ../internal/app
