package models

import "strings"

//acronym color example color White convert to WH
func AcronymColorProduct(text string) string {
	words := strings.Split(text, " ")

	res := ""

	if len(words) > 1 {
		res = res + string(words[0][0]) + string(words[1][0]) + string(words[1][1])
	} else {
		res = res + string(words[0][0]) + string(words[0][1]) + string(words[0][2])
	}

	return strings.ToUpper(res)
}

//acronym size example size S convert to SS
func AcronymSizeProduct(text string) string {
	res := ""

	switch text {
	case "S":
		res = "SS"
	case "M":
		res = "MM"
	case "L":
		res = "LL"
	case "XL":
		res = "XL"
	case "XXL":
		res = "XX"
	case "XXXL":
		res = "X3"
	}

	return res
}
