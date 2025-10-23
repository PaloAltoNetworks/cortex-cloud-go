package types

// Filter represents a search filter for compliance queries.
// This is a shared type used across multiple compliance endpoints.
type Filter struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}
