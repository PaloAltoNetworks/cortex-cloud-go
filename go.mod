module github.com/PaloAltoNetworks/cortex-cloud-go

go 1.24.3

replace github.com/PaloAltoNetworks/cortex-cloud-go => .

replace github.com/PaloAltoNetworks/cortex-cloud-go/internal/app => ./internal/app

replace github.com/PaloAltoNetworks/cortex-cloud-go/errors => ./errors

replace github.com/PaloAltoNetworks/cortex-cloud-go/enums => ./enums

replace github.com/PaloAltoNetworks/cortex-cloud-go/log => ./log

replace github.com/PaloAltoNetworks/cortex-cloud-go/api => ./api

replace github.com/PaloAltoNetworks/cortex-cloud-go/appsec => ./appsec

replace github.com/PaloAltoNetworks/cortex-cloud-go/cloudonboarding => ./cloudonboarding

replace github.com/PaloAltoNetworks/cortex-cloud-go/platform => ./platform

require github.com/go-playground/validator/v10 v10.27.0

require (
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	golang.org/x/crypto v0.33.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
)
