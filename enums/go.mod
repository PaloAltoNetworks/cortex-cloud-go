module github.com/PaloAltoNetworks/cortex-cloud-go/enums

go 1.25.0

replace (
	github.com/PaloAltoNetworks/cortex-cloud-go/appsec => ../appsec
	github.com/PaloAltoNetworks/cortex-cloud-go/client => ../internal/client
	github.com/PaloAltoNetworks/cortex-cloud-go/cloudonboarding => ../cloudonboarding
	github.com/PaloAltoNetworks/cortex-cloud-go/config => ../internal/config
	github.com/PaloAltoNetworks/cortex-cloud-go/cwp => ../cwp
	github.com/PaloAltoNetworks/cortex-cloud-go/errors => ../errors
	github.com/PaloAltoNetworks/cortex-cloud-go/internal/client => ../internal/client
	github.com/PaloAltoNetworks/cortex-cloud-go/internal/config => ../internal/config
	github.com/PaloAltoNetworks/cortex-cloud-go/log => ../log
	github.com/PaloAltoNetworks/cortex-cloud-go/platform => ../platform
	github.com/PaloAltoNetworks/cortex-cloud-go/types => ../types
)
