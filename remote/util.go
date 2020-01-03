/*
 * Copyright The Titan Project Contributors.
 */
package remote

import (
	"errors"
	"fmt"
	"net/url"
)

/*
 * Wrap remote URL parsing in an easier-to use function that will handle converting to the intermediate URL format,
 * processing any query parameters (for tags) and fragment (for commit IDs).
 */
func ParseURL(input string, properties map[string]string) (string, map[string]interface{}, []string, string, error) {
	u, err := url.Parse(input)
	if err != nil {
		return "", nil, nil, "", err
	}

	var provider string
	if u.Scheme != "" {
		provider = u.Scheme
	} else {
		provider = u.Path
	}

	var r = Get(provider)
	if r == nil {
		return "", nil, nil, "", errors.New(fmt.Sprintf("unknown remote provider '%s'", provider))
	}

	commit := u.Fragment
	tags := []string{}
	for k := range u.Query() {
		if k != "tag" {
			return "", nil, nil, "", errors.New(fmt.Sprintf("invalid query parameter '%s'", k))
		}
	}
	if u.Query()["tag"] != nil {
		tags = u.Query()["tag"]
	}

	props, err := r.FromURL(u, properties)
	if err != nil {
		return "", nil, nil, "", err
	}

	return provider, props, tags, commit, nil
}
