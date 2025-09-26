module github.com/PaloAltoNetworks/cortex-cloud-go/api

go 1.25.0

require github.com/PaloAltoNetworks/cortex-cloud-go/log v0.0.4

require (
	github.com/fatih/color v1.18.0 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/terraform-plugin-log v0.9.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.36.0 // indirect
)

replace github.com/PaloAltoNetworks/cortex-cloud-go/appsec => ../appsec

replace github.com/PaloAltoNetworks/cortex-cloud-go/client => ../client

replace github.com/PaloAltoNetworks/cortex-cloud-go/cloudonboarding => ../cloudonboarding

replace github.com/PaloAltoNetworks/cortex-cloud-go/cwp => ../cwp

replace github.com/PaloAltoNetworks/cortex-cloud-go/enums => ../enums

replace github.com/PaloAltoNetworks/cortex-cloud-go/errors => ../errors

replace github.com/PaloAltoNetworks/cortex-cloud-go/log => ../log

replace github.com/PaloAltoNetworks/cortex-cloud-go/platform => ../platform

replace github.com/PaloAltoNetworks/cortex-cloud-go/types => ../types
