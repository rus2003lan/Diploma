package elastic

import (
	"bytes"
	"strings"
)

const (
	queryStringReservedChars = `+-&|!(){}[]^"~*?:\/`
)

func QueryStringEscape(s string) string {
	i := strings.IndexAny(s, queryStringReservedChars)
	if i == -1 {
		return s
	}

	var buf bytes.Buffer
	for i != -1 {
		if _, err := buf.WriteString(s[:i]); err != nil {
			panic(err)
		}

		var esc string

		switch s[i] {
		case '+':
			esc = "\\+"

		case '-':
			esc = "\\-"

		case '&':
			esc = "\\&"

		case '|':
			esc = "\\|"

		case '!':
			esc = "\\!"

		case '(':
			esc = "\\("

		case ')':
			esc = "\\)"

		case '{':
			esc = "\\{"

		case '}':
			esc = "\\}"

		case '[':
			esc = "\\["

		case ']':
			esc = "\\]"

		case '^':
			esc = "\\^"

		case '"':
			esc = "\\\""

		case '~':
			esc = "\\~"

		case '*':
			esc = "\\*"

		case '?':
			esc = "\\?"

		case ':':
			esc = "\\:"

		case '\\':
			esc = "\\\\"

		case '/':
			esc = "\\/"

		default:
			panic("unrecognized escape character")
		}

		s = s[i+1:]

		if _, err := buf.WriteString(esc); err != nil {
			panic(err)
		}

		i = strings.IndexAny(s, queryStringReservedChars)
	}

	buf.WriteString(s)

	return buf.String()
}
