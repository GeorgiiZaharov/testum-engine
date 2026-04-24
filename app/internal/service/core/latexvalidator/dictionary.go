package latexvalidator

import (
	"math"
	"sort"
	"strings"
)

type Dictionary struct {
	functions []string
}

func NewDictionary() *Dictionary {
	fn := []string{
		"frac", "matrix", "sum", "lim", "begin", "end", "limits",
		"langle", "rangle", "infty", "ne", "to", "right", "left",
		"sqrt", "cdots", "ldots", "partial", "leq", "geq", "le", "ge",
		"int", "iint", "iiint", "log", "min", "prod",
		"lambda", "alpha", "beta", "gamma", "delta", "epsilon",
		"eta", "zeta", "theta", "kappa", "mu", "nu", "xi",
		"pi", "rho", "sigma", "tau", "phi", "psi", "chi", "omega",
		"div", "parallel", "sim", "simeq",
	}

	sort.Strings(fn)

	return &Dictionary{
		functions: fn,
	}
}

// =========================
// VALIDATION
// =========================

func (d *Dictionary) IsValid(cmd string) bool {
	i := sort.SearchStrings(d.functions, cmd)
	return i < len(d.functions) && d.functions[i] == cmd
}

// =========================
// SUGGESTIONS
// =========================

func (d *Dictionary) Suggest(word string) string {
	bestScore := 0.0
	var best []string

	for _, fn := range d.functions {
		score := similarity(word, fn)

		if score > bestScore {
			bestScore = score
			best = []string{fn}
		} else if math.Abs(score-bestScore) < 0.001 {
			best = append(best, fn)
		}
	}

	if bestScore < 0.5 {
		return ""
	}

	return strings.Join(best, ", ")
}
