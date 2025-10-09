package types

import (
	"encoding/json"
	"fmt"
)

// Filter is a recursive structure that can be used to build complex search filters.
// It can represent both logical groupings (AND/OR) and individual search criteria.
type Filter interface {
	isFilter() // Marker method
}

// ----------------------------------------------------------------------------
// Filter Implementations
// ----------------------------------------------------------------------------

// FilterRoot represents the root of a filter tree.
// Its fields are unexported to enforce creation via constructors.
type FilterRoot struct {
	and []Filter
	or  []Filter
}

// FilterGeneric represents a generic search criterion or a nested logical block.
// Its fields are unexported to enforce creation via constructors.
type FilterGeneric struct {
	and         []Filter
	or          []Filter
	searchField string
	searchType  string
	searchValue string
}

// FilterTimespan represents a time-based search criterion.
// Its fields are unexported to enforce creation via constructors.
type FilterTimespan struct {
	and         []Filter
	or          []Filter
	searchField string
	searchType  string
	searchValue SearchValueTimespan
}

// Marker methods for Filter interface compliance.
func (FilterRoot) isFilter()     {}
func (FilterGeneric) isFilter()  {}
func (FilterTimespan) isFilter() {}

// ----------------------------------------------------------------------------
// Constructors
// ----------------------------------------------------------------------------

// NewRootFilter returns a new Filter that represents the root of a nested
// collection of filters.
func NewRootFilter(and []Filter, or []Filter) FilterRoot {
	fr := FilterRoot{
		and: and,
		or:  or,
	}
	if len(fr.and) == 0 {
		fr.and = nil
	}
	if len(fr.or) == 0 {
		fr.or = nil
	}
	return fr
}

// NewAndFilter returns a new Filter that represents a logical AND of the provided filters.
func NewAndFilter(filters ...Filter) FilterGeneric {
	fg := FilterGeneric{
		and: filters,
	}
	if len(fg.and) == 0 {
		fg.and = nil
	}
	return fg
}

// NewOrFilter returns a new Filter that represents a logical OR of the provided filters.
func NewOrFilter(filters ...Filter) FilterGeneric {
	fg := FilterGeneric{
		or: filters,
	}
	if len(fg.or) == 0 {
		fg.or = nil
	}
	return fg
}

// NewSearchFilter returns a new search filter criterion.
func NewSearchFilter(field, searchType, value string) Filter {
	return FilterGeneric{
		searchField: field,
		searchType:  searchType,
		searchValue: value,
	}
}

// NewTimespanFilter returns a new timespan filter criterion.
func NewTimespanFilter(field, searchType string, from, to int) Filter {
	return FilterTimespan{
		searchField: field,
		searchType:  searchType,
		searchValue: SearchValueTimespan{
			From: from,
			To:   to,
		},
	}
}

// ----------------------------------------------------------------------------
// Helper Functions
// ----------------------------------------------------------------------------

// AddAnd appends filters to the And slice of a FilterRoot.
func (f *FilterRoot) AddAnd(filters ...Filter) {
	f.and = append(f.and, filters...)
}

// AddOr appends filters to the Or slice of a FilterRoot.
func (f *FilterRoot) AddOr(filters ...Filter) {
	f.or = append(f.or, filters...)
}

// AddAnd appends filters to the And slice of a FilterGeneric.
func (f *FilterGeneric) AddAnd(filters ...Filter) {
	f.and = append(f.and, filters...)
}

// AddOr appends filters to the Or slice of a FilterGeneric.
func (f *FilterGeneric) AddOr(filters ...Filter) {
	f.or = append(f.or, filters...)
}

// AddAnd appends filters to the And slice of a FilterTimespan.
func (f *FilterTimespan) AddAnd(filters ...Filter) {
	f.and = append(f.and, filters...)
}

// AddOr appends filters to the Or slice of a FilterTimespan.
func (f *FilterTimespan) AddOr(filters ...Filter) {
	f.or = append(f.or, filters...)
}

// ----------------------------------------------------------------------------
// JSON Marshaling and Unmarshaling
// ----------------------------------------------------------------------------

// unmarshalFilter determines the concrete type of a Filter from its JSON representation and unmarshals it.
func unmarshalFilter(b []byte) (Filter, error) {
	var probe map[string]json.RawMessage
	if err := json.Unmarshal(b, &probe); err != nil {
		return nil, fmt.Errorf("failed to probe filter type: %w", err)
	}

	if _, ok := probe["SEARCH_FIELD"]; ok {
		if val, ok := probe["SEARCH_VALUE"]; ok && len(val) > 0 && val[0] == '{' {
			var f FilterTimespan
			if err := json.Unmarshal(b, &f); err != nil {
				return nil, err
			}
			return f, nil
		}
		var f FilterGeneric
		if err := json.Unmarshal(b, &f); err != nil {
			return nil, err
		}
		return f, nil
	}

	_, andOk := probe["AND"]
	_, orOk := probe["OR"]
	if andOk || orOk {
		var f FilterGeneric
		if err := json.Unmarshal(b, &f); err != nil {
			return nil, err
		}
		return f, nil
	}

	// Default to an empty generic filter.
	var f FilterGeneric
	if err := json.Unmarshal(b, &f); err != nil {
		return nil, err
	}
	return f, nil
}

func (f FilterRoot) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		And []Filter `json:"AND,omitempty"`
		Or  []Filter `json:"OR,omitempty"`
	}{
		And: f.and,
		Or:  f.or,
	})
}

func (f *FilterRoot) UnmarshalJSON(b []byte) error {
	var raw struct {
		And []json.RawMessage `json:"AND,omitempty"`
		Or  []json.RawMessage `json:"OR,omitempty"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal raw filter root: %w", err)
	}

	if len(raw.And) > 0 {
		f.and = make([]Filter, len(raw.And))
		for i, filterJSON := range raw.And {
			filter, err := unmarshalFilter(filterJSON)
			if err != nil {
				return err
			}
			f.and[i] = filter
		}
	} else {
		f.and = nil
	}

	if len(raw.Or) > 0 {
		f.or = make([]Filter, len(raw.Or))
		for i, filterJSON := range raw.Or {
			filter, err := unmarshalFilter(filterJSON)
			if err != nil {
				return err
			}
			f.or[i] = filter
		}
	} else {
		f.or = nil
	}

	return nil
}

func (f FilterGeneric) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		And         []Filter `json:"AND,omitempty"`
		Or          []Filter `json:"OR,omitempty"`
		SearchField string   `json:"SEARCH_FIELD,omitempty"`
		SearchType  string   `json:"SEARCH_TYPE,omitempty"`
		SearchValue string   `json:"SEARCH_VALUE,omitempty"`
	}{
		And:         f.and,
		Or:          f.or,
		SearchField: f.searchField,
		SearchType:  f.searchType,
		SearchValue: f.searchValue,
	})
}

func (f *FilterGeneric) UnmarshalJSON(b []byte) error {
	var raw struct {
		And         []json.RawMessage `json:"AND,omitempty"`
		Or          []json.RawMessage `json:"OR,omitempty"`
		SearchField string            `json:"SEARCH_FIELD,omitempty"`
		SearchType  string            `json:"SEARCH_TYPE,omitempty"`
		SearchValue string            `json:"SEARCH_VALUE,omitempty"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal raw generic filter: %w", err)
	}

	f.searchField = raw.SearchField
	f.searchType = raw.SearchType
	f.searchValue = raw.SearchValue

	if len(raw.And) > 0 {
		f.and = make([]Filter, len(raw.And))
		for i, filterJSON := range raw.And {
			filter, err := unmarshalFilter(filterJSON)
			if err != nil {
				return err
			}
			f.and[i] = filter
		}
	} else {
		f.and = nil
	}

	if len(raw.Or) > 0 {
		f.or = make([]Filter, len(raw.Or))
		for i, filterJSON := range raw.Or {
			filter, err := unmarshalFilter(filterJSON)
			if err != nil {
				return err
			}
			f.or[i] = filter
		}
	} else {
		f.or = nil
	}

	return nil
}

func (f FilterTimespan) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		And         []Filter            `json:"AND,omitempty"`
		Or          []Filter            `json:"OR,omitempty"`
		SearchField string              `json:"SEARCH_FIELD,omitempty"`
		SearchType  string              `json:"SEARCH_TYPE,omitempty"`
		SearchValue SearchValueTimespan `json:"SEARCH_VALUE"`
	}{
		And:         f.and,
		Or:          f.or,
		SearchField: f.searchField,
		SearchType:  f.searchType,
		SearchValue: f.searchValue,
	})
}

func (f *FilterTimespan) UnmarshalJSON(b []byte) error {
	var raw struct {
		And         []json.RawMessage   `json:"AND,omitempty"`
		Or          []json.RawMessage   `json:"OR,omitempty"`
		SearchField string              `json:"SEARCH_FIELD,omitempty"`
		SearchType  string              `json:"SEARCH_TYPE,omitempty"`
		SearchValue SearchValueTimespan `json:"SEARCH_VALUE"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal raw timespan filter: %w", err)
	}

	f.searchField = raw.SearchField
	f.searchType = raw.SearchType
	f.searchValue = raw.SearchValue

	if len(raw.And) > 0 {
		f.and = make([]Filter, len(raw.And))
		for i, filterJSON := range raw.And {
			filter, err := unmarshalFilter(filterJSON)
			if err != nil {
				return err
			}
			f.and[i] = filter
		}
	} else {
		f.and = nil
	}

	if len(raw.Or) > 0 {
		f.or = make([]Filter, len(raw.Or))
		for i, filterJSON := range raw.Or {
			filter, err := unmarshalFilter(filterJSON)
			if err != nil {
				return err
			}
			f.or[i] = filter
		}
	} else {
		f.or = nil
	}

	return nil
}

// ----------------------------------------------------------------------------
// Other Types
// ----------------------------------------------------------------------------

type FilterData struct {
	Sort   []SortFilter `json:"sort"`
	Paging PagingFilter `json:"paging"`
	Filter Filter       `json:"filter"`
}

func (fd *FilterData) UnmarshalJSON(b []byte) error {
	type alias FilterData
	var raw struct {
		Filter json.RawMessage `json:"filter"`
		alias
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal raw filter data: %w", err)
	}

	*fd = FilterData(raw.alias)
	if raw.Filter != nil {
		var err error
		fd.Filter, err = unmarshalFilter(raw.Filter)
		if err != nil {
			return fmt.Errorf("failed to unmarshal nested filter: %w", err)
		}
	}
	return nil
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
