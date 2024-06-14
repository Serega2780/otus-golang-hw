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

	words := strings.Fields(str)
	if len(words) == 0 {
		return []string{}
	}

	for _, w := range words {
		if val, ok := m[w]; ok {
			m[w] = val + 1
		} else {
			m[w] = 1
		}
	}
	for k, v := range m {
		stats = append(stats, stat{k, v})
	}
	sortSlice(&stats)
	for i := 0; i < ResultCount; i++ {
		results = append(results, stats[i].Word)
	}

	return results
}

func sortSlice(stats *[]stat) {
	sort.Slice(*stats, func(p, n int) bool {
		s := *stats
		if s[p].Count == s[n].Count {
			return strings.Compare(s[p].Word, s[n].Word) < 0
		}

		return s[n].Count < s[p].Count
	})
}
