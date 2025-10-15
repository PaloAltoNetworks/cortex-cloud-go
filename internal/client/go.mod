module github.com/PaloAltoNetworks/cortex-cloud-go/internal/client

go 1.25.0

require (
	github.com/PaloAltoNetworks/cortex-cloud-go/errors v0.0.0
	github.com/PaloAltoNetworks/cortex-cloud-go/types v0.0.0
	github.com/stretchr/testify v1.11.1
)

replace (
	github.com/PaloAltoNetworks/cortex-cloud-go/config => ../config
	github.com/PaloAltoNetworks/cortex-cloud-go/errors => ../../errors
	github.com/PaloAltoNetworks/cortex-cloud-go/types => ../../types
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
