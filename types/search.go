package types

type FilterData struct {
	Sort   []SortFilter `json:"sort"`
	Paging PagingFilter `json:"paging"`
	Filter Filter       `json:"filter"`
}

// Filter is a recursive structure that can be used to build complex search filters.
// It can represent both logical groupings (AND/OR) and individual search criteria.
type Filter struct {
	And         []*Filter `json:"AND,omitempty"`
	Or          []*Filter `json:"OR,omitempty"`
	SearchField string    `json:"SEARCH_FIELD,omitempty"`
	SearchType  string    `json:"SEARCH_TYPE,omitempty"`
	SearchValue string    `json:"SEARCH_VALUE,omitempty"`
}

type SortFilter struct {
	Field string `json:"FIELD"`
	Order string `json:"ORDER"`
}

type PagingFilter struct {
	From int `json:"from"`
	To   int `json:"to"`
}
