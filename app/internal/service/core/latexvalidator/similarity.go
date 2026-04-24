package latexvalidator

import "math"

func similarity(a, b string) float64 {
	dist := distance(a, b)
	maxLen := math.Max(float64(len(a)), float64(len(b)))

	if maxLen == 0 {
		return 1
	}

	return 1 - float64(dist)/maxLen
}

func distance(a, b string) int {
	n, m := len(a), len(b)

	if n == 0 {
		return m
	}
	if m == 0 {
		return n
	}

	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}

	for i := 0; i <= n; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= m; j++ {
		dp[0][j] = j
	}

	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}

			dp[i][j] = min(
				dp[i-1][j]+1,
				dp[i][j-1]+1,
				dp[i-1][j-1]+cost,
			)
		}
	}

	return dp[n][m]
}

func min(a, b, c int) int {
	if a < b && a < c {
		return a
	}
	if b < c {
		return b
	}
	return c
}
