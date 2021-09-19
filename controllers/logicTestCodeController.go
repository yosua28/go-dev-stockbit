package controllers

import (
	"sort"
	"strings"
)

func LogisTest(arrAnagram []string) []interface{} {

	var result []interface{}

	list := make(map[string][]string)
	for _, word := range arrAnagram {
		key := sortStr(word)
		list[key] = append(list[key], word)
	}

	for _, words := range list {
		data := []string{}
		for _, w := range words {
			data = append(data, w)
		}
		result = append(result, data)
	}
	return result
}

func sortStr(word string) string {
	s := strings.Split(word, "")
	sort.Strings(s)

	return strings.Join(s, "")
}
