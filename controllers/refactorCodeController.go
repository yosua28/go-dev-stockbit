package controllers

import (
	"strings"
)

func FindFirstStringInBracket(str string) string {
	ret := ""
	if len(str) > 0 {
		indexFirstBracketFound := strings.Index(str, "(")
		indexClosingBracketFound := strings.Index(str, ")")
		if indexClosingBracketFound >= 0 && indexFirstBracketFound >= 0 && indexClosingBracketFound > indexFirstBracketFound {
			ret = string(str[indexFirstBracketFound+1 : indexClosingBracketFound])
		}
	}
	return ret
}
