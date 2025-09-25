package util

import (
	"net/url"
)

func StringToQuery(key string, value string) url.Values {
	result := url.Values{}
	result.Add(key, value)
	return result
}

func StringSliceToQuery(key string, values []string) url.Values {
	result := url.Values{}
	for _, value := range values {
		result.Add(key, value)
	}
	return result
}
