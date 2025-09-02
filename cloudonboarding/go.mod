module github.com/PaloAltoNetworks/cortex-cloud-go/cloudonboarding

go 1.24.3

require (
	github.com/PaloAltoNetworks/cortex-cloud-go/api v0.0.1-beta
	github.com/PaloAltoNetworks/cortex-cloud-go/internal/app v0.0.0
)

replace github.com/PaloAltoNetworks/cortex-cloud-go/api => ../api
replace github.com/PaloAltoNetworks/cortex-cloud-go/internal/app => ../internal/app
