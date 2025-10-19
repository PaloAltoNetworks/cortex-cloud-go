module github.com/PaloAltoNetworks/cortex-cloud-go/log

go 1.25.0

require github.com/hashicorp/terraform-plugin-log v0.9.0

require (
	github.com/fatih/color v1.13.0 // indirect
	github.com/hashicorp/go-hclog v1.5.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	golang.org/x/sys v0.0.0-20220503163025-988cb79eb6c6 // indirect
)

replace (
	github.com/PaloAltoNetworks/cortex-cloud-go/appsec => ../appsec
	github.com/PaloAltoNetworks/cortex-cloud-go/client => ../internal/client
	github.com/PaloAltoNetworks/cortex-cloud-go/cloudonboarding => ../cloudonboarding
	github.com/PaloAltoNetworks/cortex-cloud-go/config => ../internal/config
	github.com/PaloAltoNetworks/cortex-cloud-go/cwp => ../cwp
	github.com/PaloAltoNetworks/cortex-cloud-go/enums => ../enums
	github.com/PaloAltoNetworks/cortex-cloud-go/errors => ../errors
	github.com/PaloAltoNetworks/cortex-cloud-go/internal/client => ../internal/client
	github.com/PaloAltoNetworks/cortex-cloud-go/internal/config => ../internal/config
	github.com/PaloAltoNetworks/cortex-cloud-go/platform => ../platform
	github.com/PaloAltoNetworks/cortex-cloud-go/types => ../types
)
