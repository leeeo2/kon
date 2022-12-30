package str

import "strings"

// ToUpper :to upper
// Aaaa => AAAA
func ToUpper(src string) string {
	return strings.ToUpper(src)
}

// ToLower: to lower
// AAaa => aaaa
func ToLower(src string) string {
	return strings.ToLower(src)
}

// UnderscoreToCamelCase :Underscore name to camel case
// aa_bb_cc =>  AaBbCc
func UnderscoreToCamelCase(str string) string {
	res := make([]byte, 0, len(str))
	nextCharShouldBeUpper := true // The first char upper.
	for i := 0; i < len(str); i++ {
		c := str[i]
		if nextCharShouldBeUpper {
			if c >= 'a' && c <= 'z' {
				c = c - 32
			}
		}
		if c == '_' {
			nextCharShouldBeUpper = true
			continue
		} else {
			nextCharShouldBeUpper = false
		}
		res = append(res, c)
	}
	return string(res)
}

// CamelCaseToUnderscore Camel case to underscore
func CamelCaseToUnderscore(str string) string {
	res := make([]byte, 0, len(str))
	for i := 0; i < len(str); i++ {
		c := str[i]
		if c >= 'A' && c <= 'Z' {
			c = c + 32
			if i != 0 {
				res = append(res, '_')
			}
		}
		res = append(res, c)
	}
	return string(res)
}
