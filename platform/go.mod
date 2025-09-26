module github.com/PaloAltoNetworks/cortex-cloud-go/platform

go 1.25.0

require (
	github.com/PaloAltoNetworks/cortex-cloud-go/api v0.0.4
	github.com/PaloAltoNetworks/cortex-cloud-go/client v0.0.4
	github.com/PaloAltoNetworks/cortex-cloud-go/types v0.0.4
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/PaloAltoNetworks/cortex-cloud-go/errors v0.0.4 // indirect
	github.com/PaloAltoNetworks/cortex-cloud-go/log v0.0.4 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/terraform-plugin-log v0.9.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/PaloAltoNetworks/cortex-cloud-go/api => ../api

replace github.com/PaloAltoNetworks/cortex-cloud-go/appsec => ../appsec

replace github.com/PaloAltoNetworks/cortex-cloud-go/client => ../client

replace github.com/PaloAltoNetworks/cortex-cloud-go/cloudonboarding => ../cloudonboarding

replace github.com/PaloAltoNetworks/cortex-cloud-go/cwp => ../cwp

replace github.com/PaloAltoNetworks/cortex-cloud-go/enums => ../enums

replace github.com/PaloAltoNetworks/cortex-cloud-go/errors => ../errors

replace github.com/PaloAltoNetworks/cortex-cloud-go/log => ../log

replace github.com/PaloAltoNetworks/cortex-cloud-go/types => ../types
