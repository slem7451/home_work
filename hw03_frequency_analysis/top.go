package hw03frequencyanalysis

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/biningo/bstree" //nolint:depguard
)

type Word struct {
	word      string
	frequency int
}

var regex = regexp.MustCompile(`[\p{L}\d_]+`)

func treeComp(a, b interface{}) int {
	ia, ib := a.(Word), b.(Word)

	if ia.word == ib.word {
		return 0
	}

	dis := ia.frequency - ib.frequency

	if dis == 0 {
		words := []string{ia.word, ib.word}
		slices.Sort(words)

		if words[0] == ia.word {
			return 1
		}

		return -1
	}

	return dis
}

func Top10(text string) []string {
	words := strings.Fields(text)
	dict := make(map[string]int)

	for i := range words {
		if !regex.MatchString(words[i]) {
			if utf8.RuneCountInString(words[i]) > 1 {
				dict[words[i]]++
			}
			continue
		}

		words[i] = strings.Trim(words[i], "'\".,:;`!?/\\+=*-")

		if words[i] != "" {
			dict[strings.ToLower(words[i])]++
		}
	}

	tree := bstree.NewBSTree(treeComp)

	for k, v := range dict {
		tree.Set(Word{word: k, frequency: v})
	}

	res := make([]string, 0)
	for len(res) < 10 && tree.Len() != 0 {
		word := tree.Max().(Word)
		res = append(res, word.word)
		tree.Del(word)
	}

	return res
}
