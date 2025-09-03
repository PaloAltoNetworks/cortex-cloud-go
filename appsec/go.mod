module github.com/PaloAltoNetworks/cortex-cloud-go/appsec

go 1.24.3

require (
	github.com/PaloAltoNetworks/cortex-cloud-go/api v0.0.3-beta
	github.com/PaloAltoNetworks/cortex-cloud-go/enums v0.0.3-beta
	github.com/PaloAltoNetworks/cortex-cloud-go/internal/app v0.0.3-beta3
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/PaloAltoNetworks/cortex-cloud-go v0.0.3-beta // indirect
	github.com/PaloAltoNetworks/cortex-cloud-go/errors v0.0.3-beta // indirect
	github.com/PaloAltoNetworks/cortex-cloud-go/log v0.0.3-beta // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/hashicorp/go-hclog v1.5.0 // indirect
	github.com/hashicorp/terraform-plugin-log v0.9.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.3-beta14 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
)

replace github.com/PaloAltoNetworks/cortex-cloud-go/api => ../api

replace github.com/PaloAltoNetworks/cortex-cloud-go/enums => ../enums

replace github.com/PaloAltoNetworks/cortex-cloud-go/internal/app => ../internal/app

require (
	dario.cat/mergo v1.0.2
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
