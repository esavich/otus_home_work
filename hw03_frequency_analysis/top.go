package hw03frequencyanalysis

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type CountedWord struct {
	word string
	freq int
}

var (
	charRegexp = regexp.MustCompile(`[a-zA-Zа-яА-Я]`)
	// Не уверен, что правильно понял условие, сделал исключение для дефисов (2 и больше).
	defisRegexp = regexp.MustCompile(`^-{2,}$`)
)

func Top10(text string) []string {
	if text == "" {
		return nil
	}

	words := strings.Fields(text)

	freq := make(map[string]int)
	for _, v := range words {
		if v == "-" {
			continue
		}

		lowerString := strings.ToLower(v)
		if defisRegexp.MatchString(lowerString) {
			freq[lowerString]++
			continue
		}

		trimmedString := strings.TrimFunc(lowerString, func(r rune) bool {
			return !charRegexp.MatchString(string(r))
		})

		if trimmedString != "" {
			freq[trimmedString]++
		}
	}
	wordsWithFreq := make([]CountedWord, 0, len(freq))
	for w, f := range freq {
		wordsWithFreq = append(wordsWithFreq, CountedWord{w, f})
	}

	sort.Slice(wordsWithFreq, func(i, j int) bool {
		if wordsWithFreq[i].freq == wordsWithFreq[j].freq {
			return wordsWithFreq[i].word < wordsWithFreq[j].word
		}
		return wordsWithFreq[i].freq > wordsWithFreq[j].freq
	})
	fmt.Println(freq)
	fmt.Println(wordsWithFreq)

	res := make([]string, 0, 10)
	for i := 0; i < 10; i++ {
		if i < len(wordsWithFreq) {
			res = append(res, wordsWithFreq[i].word)
		}
	}
	fmt.Println(res)

	return res
}
