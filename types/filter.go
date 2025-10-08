package types

import (
	"encoding/json"
	"fmt"
)

type FilterData struct {
	Sort   []SortFilter `json:"sort"`
	Paging PagingFilter `json:"paging"`
	Filter Filter       `json:"filter"`
}

// Filter is a recursive structure that can be used to build complex search filters.
// It can represent both logical groupings (AND/OR) and individual search criteria.
type Filter interface {
	ToJSON() (string, error)
}

type FilterGeneric struct {
	And         []Filter `json:"AND,omitempty"`
	Or          []Filter `json:"OR,omitempty"`
	SearchField string   `json:"SEARCH_FIELD,omitempty"`
	SearchType  string   `json:"SEARCH_TYPE,omitempty"`
	SearchValue string   `json:"SEARCH_VALUE,omitempty"`
}

func (f FilterGeneric) ToJSON() (jsonString string, err error) {
	data, err := json.Marshal(f)
	if err != nil {
		return "", fmt.Errorf("failed to marshal FilterGeneric to JSON: %w", err)
	}
	return string(data), nil
}

type FilterTimespan struct {
	And         []Filter            `json:"AND,omitempty"`
	Or          []Filter            `json:"OR,omitempty"`
	SearchField string              `json:"SEARCH_FIELD,omitempty"`
	SearchType  string              `json:"SEARCH_TYPE,omitempty"`
	SearchValue SearchValueTimespan `json:"SEARCH_VALUE"`
}

func (f FilterTimespan) ToJSON() (jsonString string, err error) {
	data, err := json.Marshal(f)
	if err != nil {
		return "", fmt.Errorf("failed to marshal FilterTimespan to JSON: %w", err)
	}
	return string(data), nil
}

type SearchValueTimespan struct {
	From int `json:"from,omitempty"`
	To   int `json:"to,omitempty"`
}

type SortFilter struct {
	Field string `json:"FIELD"`
	Order string `json:"ORDER"`
}

type PagingFilter struct {
	From int `json:"from"`
	To   int `json:"to"`
}

// NewAndFilter returns a new Filter that represents a logical AND of the provided filters.
func NewAndFilter(filters ...Filter) Filter {
	return FilterGeneric{
		And: filters,
	}
}

// NewOrFilter returns a new Filter that represents a logical OR of the provided filters.
func NewOrFilter(filters ...Filter) Filter {
	return FilterGeneric{
		Or: filters,
	}
}

// NewSearchFilter returns a new search filter criterion.
func NewSearchFilter(field, searchType, value string) Filter {
	return FilterGeneric{
		SearchField: field,
		SearchType:  searchType,
		SearchValue: value,
	}
}

// NewTimespanFilter returns a new timespan filter criterion.
func NewTimespanFilter(field, searchType string, from, to int) Filter {
	return FilterTimespan{
		SearchField: field,
		SearchType:  searchType,
		SearchValue: SearchValueTimespan{
			From: from,
			To:   to,
		},
	}
}
