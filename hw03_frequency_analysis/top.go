package hw03frequencyanalysis

import (
	"container/heap"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

type Word struct {
	word      string
	frequency int
}

type WordHeap []Word

func (h WordHeap) Len() int { 
	return len(h) 
}

func (h WordHeap) Less(i, j int) bool {
	dis := h[i].frequency - h[j].frequency

	if dis == 0 {
		words := []string{h[i].word, h[j].word}
		slices.Sort(words)
		return words[0] == h[i].word
	}

	return dis > 0
}

func (h WordHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *WordHeap) Push(x any) {
	*h = append(*h, x.(Word))
}

func (h *WordHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

var regex = regexp.MustCompile(`[\p{L}\d_]+`)

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

	h := &WordHeap{}
	heap.Init(h)

	for k, v := range dict {
		heap.Push(h, Word{word: k, frequency: v})
	}

	res := make([]string, 0)
	for len(res) < 10 && h.Len() != 0 {
		word := heap.Pop(h)
		res = append(res, word.(Word).word)
	}

	return res
}
