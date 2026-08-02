package assets

import (
	"os"
	"strconv"
	"strings"
)

// SpelloCorrecter computes weighted Levenshtein distance to find closest matching names.
type SpelloCorrecter struct {
	weights map[rune]map[rune]float64
}

// NewSpelloCorrecter initializes a SpelloCorrecter with substitution weight mappings.
func NewSpelloCorrecter() *SpelloCorrecter {
	weights := map[rune]map[rune]float64{
		'k': {'c': 0.25, 'g': 0.75, 'q': 0.125},
		'c': {'k': 0.25, 'g': 0.75, 's': 0.5, 'z': 0.5, 'q': 0.125},
		's': {'z': 0.25, 'c': 0.5},
		'z': {'s': 0.25, 'c': 0.5},
		'g': {'k': 0.75, 'c': 0.75, 'q': 0.9},
		'o': {'u': 0.5},
		'u': {'o': 0.5, 'v': 0.75, 'w': 0.5},
		'b': {'v': 0.75},
		'v': {'b': 0.75, 'w': 0.5, 'u': 0.7},
		'w': {'v': 0.5, 'u': 0.5},
		'q': {'c': 0.125, 'k': 0.125, 'g': 0.9},
	}
	return &SpelloCorrecter{weights: weights}
}

// Distance computes the weighted Levenshtein distance between target and candidate.
func (sc *SpelloCorrecter) Distance(target, candidate string) float64 {
	targetRunes := []rune(strings.ToLower(target))
	candRunes := []rune(strings.ToLower(candidate))

	m := len(targetRunes)
	n := len(candRunes)

	if m == 0 {
		return float64(n)
	}
	if n == 0 {
		return float64(m)
	}

	dp := make([][]float64, m+1)
	for i := range dp {
		dp[i] = make([]float64, n+1)
		dp[i][0] = float64(i)
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = float64(j)
	}

	for i := 1; i <= m; i++ {
		t := targetRunes[i-1]
		for j := 1; j <= n; j++ {
			c := candRunes[j-1]

			cost := 1.0
			if t == c {
				cost = 0.0
			} else if wMap, ok := sc.weights[t]; ok {
				if w, ok2 := wMap[c]; ok2 {
					cost = w
				}
			}

			del := dp[i-1][j] + 1.0
			ins := dp[i][j-1] + 1.0
			sub := dp[i-1][j-1] + cost

			minVal := del
			if ins < minVal {
				minVal = ins
			}
			if sub < minVal {
				minVal = sub
			}
			dp[i][j] = minVal
		}
	}

	return dp[m][n]
}

// Correct finds the closest matching candidates for target string.
// Returns candidate strings with minimum distance, and that distance.
func (sc *SpelloCorrecter) Correct(target string, candidates []string) ([]string, float64) {
	if len(candidates) == 0 {
		return nil, 1e9
	}

	limit := 5.0
	if envLimit := os.Getenv("PONYSAY_TYPO_LIMIT"); envLimit != "" {
		if parsed, err := strconv.ParseFloat(envLimit, 64); err == nil {
			limit = parsed
		}
	}

	var bestMatches []string
	bestDist := 1e9

	for _, cand := range candidates {
		dist := sc.Distance(target, cand)
		if dist < bestDist {
			bestDist = dist
			bestMatches = []string{cand}
		} else if dist == bestDist {
			bestMatches = append(bestMatches, cand)
		}
	}

	if bestDist <= limit {
		return bestMatches, bestDist
	}

	return nil, bestDist
}
