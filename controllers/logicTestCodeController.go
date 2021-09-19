package controllers

import (
	"log"
	"sort"
	"strings"
)

func logisTest() []interface{} {

	var result []interface{}

	arrAnagram := []string{"kita", "atik", "tika", "aku", "kia", "makan", "kua"}
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
	log.Println(result)
	return result
}

func sortStr(word string) string {
	s := strings.Split(word, "")
	sort.Strings(s)

	return strings.Join(s, "")
}
