package types

// Filter is a recursive structure that can be used to build complex search filters.
// It can represent both logical groupings (AND/OR) and individual search criteria.
type Filter interface {
	isFilter()
}
