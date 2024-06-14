package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

const ResultCount int = 10

type stat struct {
	Word  string
	Count int
}

func Top10(str string) []string {
	m := make(map[string]int)
	stats := []stat{}
	var results []string
	counter := ResultCount

	words := strings.Fields(str)
	if len(words) == 0 {
		return []string{}
	}
	for _, w := range words {
		m[w]++
	}
	for k, v := range m {
		stats = append(stats, stat{k, v})
	}
	sortSlice(&stats)

	if len(stats) < ResultCount {
		counter = len(stats)
	}
	for i := 0; i < counter; i++ {
		results = append(results, stats[i].Word)
	}

	return results
}

func sortSlice(stats *[]stat) {
	sort.Slice(*stats, func(p, n int) bool {
		s := *stats
		if s[p].Count == s[n].Count {
			return s[p].Word < s[n].Word
		}

		return s[n].Count < s[p].Count
	})
}
