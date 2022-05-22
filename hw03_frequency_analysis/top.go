package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var re = regexp.MustCompile(`[^0-9А-Яа-яA-za-z_]`)

func Top10(str string) []string {
	var result []string
	if str == "" {
		return result
	}
	sl := strings.Fields(str)

	dict := make(map[string]int)

	for _, word := range sl {
		word = strings.ToLower(word)
		word := re.ReplaceAllString(word, "")
		if word == "" {
			continue
		}
		dict[word]++
	}
	keys := make([]string, 0, len(dict))
	values := make([]int, 0, len(dict))
	for k, v := range dict {
		keys = append(keys, k)
		values = append(values, v)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	sort.Slice(values, func(i, j int) bool {
		return values[i] > values[j]
	})

	prev := values[0]
	for i, val := range values {
		if i != 0 && prev == val {
			prev = val
			continue
		}
		for _, key := range keys {
			if dict[key] == val {
				result = append(result, key)
				delete(dict, key)
			}
		}
	}

	if len(result) < 10 {
		return result
	}

	return result[:10]
}
