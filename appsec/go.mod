module github.com/PaloAltoNetworks/cortex-cloud-go/appsec

go 1.25.0

require github.com/PaloAltoNetworks/cortex-cloud-go/enums v0.0.8

replace (
	github.com/PaloAltoNetworks/cortex-cloud-go/client => ../internal/client
	github.com/PaloAltoNetworks/cortex-cloud-go/config => ../internal/config
)
