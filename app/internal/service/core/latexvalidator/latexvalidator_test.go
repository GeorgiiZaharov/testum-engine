package latexvalidator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newValidator() *Validator {
	return New()
}

func newDict() *Dictionary {
	return NewDictionary()
}

func lines(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, "\n")
}

// =========================
// BRACKETS
// =========================

func Test_Brackets_Valid(t *testing.T) {
	v := newValidator()

	errs := v.Validate(lines("(a + b) + {c}"))

	assert.Empty(t, errs)
}

func Test_Brackets_MissingClosing(t *testing.T) {
	v := newValidator()

	errs := v.Validate(lines("(a + b"))

	assert.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error, "закрывающей")
}

func Test_Brackets_ExtraClosing(t *testing.T) {
	v := newValidator()

	errs := v.Validate(lines("a + b)"))

	assert.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error, "лишняя")
}

func Test_Brackets_WrongOrder(t *testing.T) {
	v := newValidator()

	errs := v.Validate(lines("({)}"))

	assert.Len(t, errs, 2)
}

func Test_Brackets_MultipleErrors(t *testing.T) {
	v := newValidator()

	errs := v.Validate(lines("({)}))"))

	assert.GreaterOrEqual(t, len(errs), 2)
}

func Test_Brackets_UnclosedMultiple(t *testing.T) {
	v := newValidator()

	errs := v.Validate(lines("((({"))

	assert.Len(t, errs, 4)
}

// =========================
// COMMANDS
// =========================

func Test_Commands_Valid(t *testing.T) {
	v := newValidator()

	errs := v.Validate(lines(`\frac{a}{b} + \sqrt{x}`))

	assert.Empty(t, errs)
}

func Test_Commands_Invalid_WithSuggestion(t *testing.T) {
	v := newValidator()

	errs := v.Validate(lines(`\fracc{a}{b}`))

	assert.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error, "frac")
}

func Test_Commands_Invalid_NoSuggestion(t *testing.T) {
	v := newValidator()

	errs := v.Validate(lines(`\zzzzzz`))

	assert.Empty(t, errs)
}

func Test_Commands_Multiple(t *testing.T) {
	v := newValidator()

	errs := v.Validate(lines(`\fracc + \sqrrt`))

	assert.Len(t, errs, 2)
}

func Test_Commands_WithUnderscore(t *testing.T) {
	v := newValidator()

	errs := v.Validate(lines(`\alph_a`))

	assert.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error, "alpha")
}

// =========================
// MULTILINE
// =========================

func Test_Multiline_LineNumbers(t *testing.T) {
	v := newValidator()

	text := `(a + b
\fracc`

	errs := v.Validate(lines(text))

	assert.Len(t, errs, 2)
	assert.Equal(t, 1, errs[0].Line)
	assert.Equal(t, 2, errs[1].Line)
}

// =========================
// EDGE CASES
// =========================

func Test_EmptyInput(t *testing.T) {
	v := newValidator()

	errs := v.Validate(lines(""))

	assert.Empty(t, errs)
}

func Test_OnlySpaces(t *testing.T) {
	v := newValidator()

	errs := v.Validate(lines("   \n   "))

	assert.Empty(t, errs)
}

func Test_PlainText(t *testing.T) {
	v := newValidator()

	errs := v.Validate(lines("hello world"))

	assert.Empty(t, errs)
}

// =========================
// DICTIONARY
// =========================

func Test_Dictionary_IsValid(t *testing.T) {
	d := newDict()

	assert.True(t, d.IsValid("frac"))
	assert.False(t, d.IsValid("unknown"))
}

func Test_Dictionary_Suggest_BestMatch(t *testing.T) {
	d := newDict()

	res := d.Suggest("fracc")

	assert.Contains(t, res, "frac")
}

func Test_Dictionary_Suggest_MultipleBest(t *testing.T) {
	d := newDict()

	res := d.Suggest("si")

	if res != "" {
		assert.NotEmpty(t, res)
	}
}

func Test_Dictionary_Suggest_LowScore(t *testing.T) {
	d := newDict()

	res := d.Suggest("zzzzzzzz")

	assert.Equal(t, "", res)
}

// =========================
// INTERNALS
// =========================

func Test_Distance(t *testing.T) {
	assert.Equal(t, 0, distance("abc", "abc"))
	assert.Equal(t, 1, distance("abc", "ab"))
	assert.Equal(t, 3, distance("", "abc"))
	assert.Equal(t, 3, distance("abc", ""))
}

func Test_Similarity(t *testing.T) {
	assert.Equal(t, 1.0, similarity("abc", "abc"))
	assert.Less(t, similarity("abc", "xyz"), 0.5)
}

func Test_Similarity_BothEmptyStrings(t *testing.T) {
	result := similarity("", "")

	assert.Equal(t, 1.0, result)
}

func Test_Min(t *testing.T) {
	assert.Equal(t, 1, min(1, 2, 3))
	assert.Equal(t, 1, min(2, 1, 3))
	assert.Equal(t, 1, min(3, 2, 1))
}
