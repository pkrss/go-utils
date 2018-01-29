package strings

/**
	"aB" = StringFromCamelCase("a_b")
**/
func StringFromCamelCase(str string) string {
	c := len(str)
	s := ""
	var b byte
	for i := 0; i < c; i++ {
		b = str[i]
		if b == '_' {
			i++
			if i >= c {
				continue
			}
			b = str[i]
			if b >= 'a' && b <= 'z' {
				s += string(b - 'a' + 'A')
			}
		} else {
			s += string(b)
		}
	}
	return s
}

/**
	"a_b" = StringToCamelCase("aB")
**/
func StringToCamelCase(str string) string {
	c := len(str)
	s := ""
	var b byte
	for i := 0; i < c; i++ {
		b = str[i]
		if b >= 'A' && b <= 'Z' {
			s += "_"
			s += string((b - 'A' + 'a'))
		} else {
			s += string(b)
		}
	}
	return s
}
