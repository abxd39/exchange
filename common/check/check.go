package check

import "regexp"

func CheckPhone(phone string) bool {
	reg := `^1([38][0-9]|14[57]|5[^4])\d{8}$`
	rgx := regexp.MustCompile(reg)
	if ok := rgx.MatchString(phone); !ok {
		return false
	}
	return true
}

func CheckEmail(email string) bool {
	reg := `^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`
	rgx := regexp.MustCompile(reg)
	if ok := rgx.MatchString(email); !ok {
		return false
	}
	return true
}
