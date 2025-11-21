module github.com/PaloAltoNetworks/cortex-cloud-go/compliance

go 1.25.0

require (
	github.com/PaloAltoNetworks/cortex-cloud-go/internal/client v0.0.0
	github.com/PaloAltoNetworks/cortex-cloud-go/internal/config v0.0.0
	github.com/PaloAltoNetworks/cortex-cloud-go/log v0.0.0
	github.com/PaloAltoNetworks/cortex-cloud-go/types v0.0.0
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/PaloAltoNetworks/cortex-cloud-go/errors v0.0.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/terraform-plugin-log v0.9.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	golang.org/x/sys v0.36.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/PaloAltoNetworks/cortex-cloud-go/errors => ../errors
	github.com/PaloAltoNetworks/cortex-cloud-go/internal/client => ../internal/client
	github.com/PaloAltoNetworks/cortex-cloud-go/internal/config => ../internal/config
	github.com/PaloAltoNetworks/cortex-cloud-go/log => ../log
	github.com/PaloAltoNetworks/cortex-cloud-go/types => ../types
)
