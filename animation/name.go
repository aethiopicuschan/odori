package animation

import "regexp"

func IsValidName(name string) bool {
	r, _ := regexp.Compile("^[0-9a-zA-Z]+$")
	return r.Match([]byte(name))
}
