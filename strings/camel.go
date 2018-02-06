package strings

/**
	"ABcDef" => "a_bc_def"
**/
func StringFromCamelCase(str string) string {
	c := len(str)
	s := ""
	var b byte
	for i := 0; i < c; i++ {
		b = str[i]
		if i == 0 {
			if b >= 'A' && b <= 'Z' {
				s += string(b - 'a' + 'A')
				continue
			}
		}
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
	"ab_c_d_ef" => "AbCDef"
**/
func StringToCamelCase(str string) string {
	c := len(str)
	s := ""
	var b byte
	for i := 0; i < c; i++ {
		b = str[i]
		if b >= 'A' && b <= 'Z' {
			if i > 0 {
				s += "_"
			}
			s += string((b - 'A' + 'a'))
		} else {
			if i == 0 {
				if b >= 'A' && b <= 'Z' {
					b = b - 'A' + 'a'
				}
			}
			s += string(b)
		}
	}
	return s
}
