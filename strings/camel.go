package strings

/**
	"aB" = StringFromCamelCase("a_b")
**/
func StringFromCamelCase(str string) string {
	c := len(str)
	s := ""
	for i := 0; i < c; i++ {
		c := str[i]
		if c == '_' {
			i++
			if i >= c {
				continue
			}
			c = str[i]
			if c >= 'a' && c <= 'z' {
				s += c - 'a' + 'A'
			}
		} else {
			s += c
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
	for i := 0; i < c; i++ {
		c := str[i]
		if c >= 'A' && c <= 'Z' {
			s += "_" + (c - 'A' + 'a')
		} else {
			s += c
		}
	}
	return s
}
