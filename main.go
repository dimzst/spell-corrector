package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

var reWords = regexp.MustCompile("\\w+")

var words map[string]int
var totalWords int

func main() {
	readSourceFile()
	fmt.Println(correction(os.Args[1]))
}

func readSourceFile() {
	bs, err := ioutil.ReadFile("big.txt")
	if err != nil {
		log.Fatal(err)
	}
	words = make(map[string]int)
	result := reWords.FindAll(bs, -1)
	for _, v := range result {
		words[strings.ToLower(string(v))]++
		totalWords++
	}
}

func correction(word string) string {
	var maxProbabilty float64
	var result string
	list := candidates(word)
	for _, v := range list {
		probability := P(v)
		if probability > maxProbabilty {
			maxProbabilty = probability
			result = v
		}
	}
	if result == "" {
		result = word
	}
	return result
}

func candidates(word string) []string {
	var result []string
	wordSet := set([]string{word})
	result = known(wordSet)
	if len(result) > 0 {
		return result
	}
	result = known(edits1(word))
	if len(result) > 0 {
		return result
	}
	result = known(edits2(word))
	if len(result) > 0 {
		return result
	}
	return []string{word}
}

func P(word string) float64 {
	return float64(words[word]) / float64(totalWords)
}

func known(input map[string]bool) []string {
	result := make([]string, 0)
	for k := range input {
		if _, ok := words[k]; ok {
			result = append(result, k)
		}
	}
	return result
}

func edits2(input string) map[string]bool {
	e1 := edits1(input)
	result := make([]string, 0)
	for k := range e1 {
		e2 := edits1(k)
		for word := range e2 {
			result = append(result, word)
		}
	}
	return set(result)
}

func edits1(input string) map[string]bool {
	letters := "abcdefghijklmnopqrstuvwxyz"
	splits := make([][]string, 0)
	for i := range input {
		splits = append(splits, []string{input[:i], input[i:]})
	}
	deletes := make([]string, 0)
	for _, val := range splits {
		if len(val[1]) > 0 {
			deletes = append(deletes, val[0]+val[1][1:])
		}
	}
	transposes := make([]string, 0)
	for _, val := range splits {
		if len(val[1]) > 1 {
			transposes = append(transposes, val[0]+val[1][1:2]+val[1][:1]+val[1][2:])
		}
	}
	replaces := make([]string, 0)
	for _, val := range splits {
		if len(val[1]) > 0 {
			for _, c := range letters {
				replaces = append(replaces, val[0]+string(c)+val[1][1:])
			}
		}
	}
	inserts := make([]string, 0)
	for _, val := range splits {
		for _, c := range letters {
			inserts = append(inserts, val[0]+string(c)+val[1])
		}
	}
	result := make([]string, 0)
	result = append(result, deletes...)
	result = append(result, transposes...)
	result = append(result, replaces...)
	result = append(result, inserts...)
	return set(result)
}

func set(input []string) map[string]bool {
	result := make(map[string]bool)
	for _, v := range input {
		result[v] = true
	}
	return result
}
