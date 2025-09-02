module github.com/PaloAltoNetworks/cortex-cloud-go/internal/app

go 1.24.3

require (
	github.com/PaloAltoNetworks/cortex-cloud-go/api v0.0.1-beta
	github.com/PaloAltoNetworks/cortex-cloud-go/errors v0.0.1-beta
	github.com/PaloAltoNetworks/cortex-cloud-go/log v0.0.1-beta
	github.com/stretchr/testify v1.10.0
)

replace github.com/PaloAltoNetworks/cortex-cloud-go/api => ../../api
replace github.com/PaloAltoNetworks/cortex-cloud-go/errors => ../../errors
replace github.com/PaloAltoNetworks/cortex-cloud-go/log => ../../log

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
