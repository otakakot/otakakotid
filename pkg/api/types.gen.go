// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Message string `json:"message"`
}

// FinalizeAssertionRequest Finalize Assertion Request
type FinalizeAssertionRequest = map[string]interface{}

// FinalizeAssertionResponse Finalize Assertion Response
type FinalizeAssertionResponse = map[string]interface{}

// FinalizeAttestationRequest Finalize Attestation Request
type FinalizeAttestationRequest = map[string]interface{}

// InitializeAssertionRequest Initialize Assertion Response
type InitializeAssertionRequest = map[string]interface{}

// InitializeAssertionResponse Initialize Assertion Response
type InitializeAssertionResponse = map[string]interface{}

// InitializeAttestationResponse Initialize Attestation Response
type InitializeAttestationResponse = map[string]interface{}

// OpenIDConfigurationResponseSchema https://openid.net/specs/openid-connect-discovery-1_0.html#ProviderMetadata
type OpenIDConfigurationResponseSchema struct {
	// AuthorizationEndpoint http://localhost:8787/authorize
	AuthorizationEndpoint            string   `json:"authorization_endpoint"`
	IdTokenSigningAlgValuesSupported []string `json:"id_token_signing_alg_values_supported"`

	// Issuer http://localhost:8787
	Issuer string `json:"issuer"`

	// JwksUri http://localhost:8787/certs
	JwksUri string `json:"jwks_uri"`

	// RevocationEndpoint http://localhost:8787/revoke
	RevocationEndpoint    string   `json:"revocation_endpoint"`
	SubjectTypesSupported []string `json:"subject_types_supported"`

	// TokenEndpoint http://localhost:8787/token
	TokenEndpoint string `json:"token_endpoint"`

	// UserinfoEndpoint http://localhost:8787/userinfo
	UserinfoEndpoint string `json:"userinfo_endpoint"`
}

// RegistrationRequest Registration Request
type RegistrationRequest struct {
	Email string `json:"email"`
}

// FinalizeAssertionParams defines parameters for FinalizeAssertion.
type FinalizeAssertionParams struct {
	// Assertion session
	Assertion string `form:"__assertion__" json:"__assertion__"`
}

// FinalizeAttestationParams defines parameters for FinalizeAttestation.
type FinalizeAttestationParams struct {
	// Attestation session
	Attestation string `form:"__attestation__" json:"__attestation__"`
}

// FinalizeRegistrationParams defines parameters for FinalizeRegistration.
type FinalizeRegistrationParams struct {
	// Code code
	Code string `form:"code" json:"code"`
}

// InitializeAssertionJSONRequestBody defines body for InitializeAssertion for application/json ContentType.
type InitializeAssertionJSONRequestBody = InitializeAssertionRequest

// FinalizeAssertionJSONRequestBody defines body for FinalizeAssertion for application/json ContentType.
type FinalizeAssertionJSONRequestBody = FinalizeAssertionRequest

// FinalizeAttestationJSONRequestBody defines body for FinalizeAttestation for application/json ContentType.
type FinalizeAttestationJSONRequestBody = FinalizeAttestationRequest

// FinalizeRegistrationJSONRequestBody defines body for FinalizeRegistration for application/json ContentType.
type FinalizeRegistrationJSONRequestBody = FinalizeAttestationRequest

// InitialiseRegistrationJSONRequestBody defines body for InitialiseRegistration for application/json ContentType.
type InitialiseRegistrationJSONRequestBody = RegistrationRequest
