package form

import (
	"strings"
)

type UITag string

func (tag UITag) Lookup(key string) (string, bool) {
	input := string(tag)
	k := ""
	i := 0
	start := i
	for i < len(input) {
		for i < len(input) && input[i] == ' ' {
			i++
		}
		start = i
		for i < len(input) && input[i] != '=' && input[i] != ',' && input[i] != ' ' {
			i++
		}
		k = input[start:i]
		for i < len(input) && input[i] == ' ' {
			i++
		}
		if i < len(input) && input[i] == '=' {
			i++
			for i < len(input) && input[i] == ' ' {
				i++
			}
			var v string
			if i < len(input) && input[i] == '"' {
				i++
				start := i
				for i < len(input) {
					if input[i] == '\\' {
						i++
					} else if input[i] == '"' {
						break
					}
					i++
				}
				v = input[start:i]
				v = strings.ReplaceAll(v, `\"`, `"`)
				i++
			} else {
				start := i
				for i < len(input) && input[i] != ',' {
					i++
				}
				v = input[start:i]
				v = strings.TrimSpace(v)
			}
			if k == key {
				return v, true
			}
		} else if k == key {
			return "", true
		}
		if i < len(input) && input[i] == ',' {
			i++
		}
	}
	return "", false
}
