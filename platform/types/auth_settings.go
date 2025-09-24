package types

type AuthSettings struct {
	TenantID           string           `json:"tenant_id"`
	Name               string           `json:"name"`
	Domain             string           `json:"domain"`
	IDPEnabled         bool             `json:"idp_enabled"`
	DefaultRole        string           `json:"default_role"`
	IsAccountRole      bool             `json:"is_account_role"`
	IDPCertificate     string           `json:"idp_certificate"`
	IDPIssuer          string           `json:"idp_issuer"`
	IDPSingleSignOnURL string           `json:"idp_sso_url"`
	MetadataURL        string           `json:"metadata_url"`
	Mappings           Mappings         `json:"mappings"`
	AdvancedSettings   AdvancedSettings `json:"advanced_settings"`
	SpEntityID         string           `json:"sp_entity_id"`
	SpLogoutURL        string           `json:"sp_logout_url"`
	SpURL              string           `json:"sp_url"`
}

type Mappings struct {
	Email     string `json:"email"`
	FirstName string `json:"firstname"`
	GroupName string `json:"group_name"`
	LastName  string `json:"lastname"`
}

type AdvancedSettings struct {
	AuthnContextEnabled bool `json:"authn_context_enabled"`
	//ForceAuthn bool `json:"authn_context_enabled"`
	IDPSingleLogoutURL        string `json:"idp_single_logout_url"`
	RelayState                string `json:"relay_state"`
	ServiceProviderPrivateKey string `json:"service_provider_private_key"`
	ServiceProviderPublicCert string `json:"service_provider_public_cert"`
}

// --------------------------- 
// Request/Response structs
// ---------------------------

// ListIDPMetadata

type ListIDPMetadataRequestData struct{}

type ListIDPMetadataResponse struct {
	TenantID    string `json:"tenant_id"`
	SpEntityID  string `json:"sp_entity_id"`
	SpLogoutURL string `json:"sp_logout_url"`
	SpURL       string `json:"sp_url"`
}

// ListAuthSettings

// TODO: This endpoint currently doesn't have any parameters defined. Populate
// or remove this before release, depending on whether or not the endpoint 
// has been updated before then.
type ListAuthSettingsRequestData struct{}

// CreateAuthSettings

type CreateAuthSettingsRequestData struct {
	Name               string           `json:"name"`
	DefaultRole        string           `json:"default_role"`
	IsAccountRole      bool             `json:"is_account_role"`
	Domain             string           `json:"domain"`
	Mappings           Mappings         `json:"mappings"`
	AdvancedSettings   AdvancedSettings `json:"advanced_settings"`
	IDPSingleSignOnURL string           `json:"idp_sso_url"`
	IDPCertificate     string           `json:"idp_certificate"`
	IDPIssuer          string           `json:"idp_issuer"`
	MetadataURL        string           `json:"metadata_url"`
}

// UpdateAuthSettings

type UpdateAuthSettingsRequestData struct {
	Name               string           `json:"name"`
	DefaultRole        string           `json:"default_role"`
	IsAccountRole      bool             `json:"is_account_role"`
	CurrentDomain      string           `json:"current_domain_value"`
	NewDomain          string           `json:"new_domain_value"`
	Mappings           Mappings         `json:"mappings"`
	AdvancedSettings   AdvancedSettings `json:"advanced_settings"`
	IDPSingleSignOnURL string           `json:"idp_sso_url"`
	IDPCertificate     string           `json:"idp_certificate"`
	IDPIssuer          string           `json:"idp_issuer"`
	MetadataURL        string           `json:"metadata_url"`
}

// DeleteAuthSettings

type DeleteAuthSettingsRequestData struct {
	Domain string `json:"domain"`
}
