package str

import "testing"

func TestConvert(t *testing.T) {
	// test low case to upper case
	lowCase := "AaaaBbbbCccc"
	want := "AAAABBBBCCCC"
	if got := ToUpper(lowCase); got != want {
		t.Errorf("convert low case to upper case failed,want %s,but got %s", want, got)
	}

	// test upper case to low case
	upperCase := "AAAaBbbbCCcc"
	want = "aaaabbbbcccc"
	if got := ToLower(upperCase); got != want {
		t.Errorf("convert upper case to low case failed,want %s,but got %s", want, got)
	}

	// test underscore to camel case
	src := "custom_config_path"
	want = "CustomConfigPath"
	if got := UnderscoreToCamelCase(src); got != want {
		t.Errorf("convert underscore to camel case failed,want %s,but got %s", want, got)
	}

	// test camel case to underscore
	src = "CustomConfigPath"
	want = "custom_config_path"
	if got := CamelCaseToUnderscore(src); got != want {
		t.Errorf("convert camel case to underscore failed,want %s,but got %s", want, got)
	}
}
