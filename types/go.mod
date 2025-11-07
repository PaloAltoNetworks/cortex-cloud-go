module github.com/PaloAltoNetworks/cortex-cloud-go/types

go 1.25.0

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
	github.com/PaloAltoNetworks/cortex-cloud-go/log => ../log
	github.com/PaloAltoNetworks/cortex-cloud-go/platform => ../platform
)

require github.com/PaloAltoNetworks/cortex-cloud-go/log v0.0.0-00010101000000-000000000000

require (
	github.com/fatih/color v1.18.0 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/terraform-plugin-log v0.9.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.36.0 // indirect
)
