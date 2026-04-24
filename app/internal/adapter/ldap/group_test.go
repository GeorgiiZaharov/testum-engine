package ldap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractAcademicGroup_NilOrEmpty(t *testing.T) {
	assert.Nil(t, extractAcademicGroup(""))
}

func TestExtractAcademicGroup_NoOU(t *testing.T) {
	dn := "cn=user,dc=example,dc=com"

	assert.Nil(t, extractAcademicGroup(dn))
}

func TestExtractAcademicGroup_NonStudentDN_Ignored(t *testing.T) {
	dn := "cn=user,ou=22307,dc=example,dc=com"

	assert.Nil(t, extractAcademicGroup(dn))
}

// student + invalid ou
func TestExtractAcademicGroup_Student_NoNumericOU(t *testing.T) {
	dn := "cn=user,ou=abc123,dc=example,dc=com,ou=student"

	assert.Nil(t, extractAcademicGroup(dn))
}

// student + valid ou
func TestExtractAcademicGroup_Student_ValidOU(t *testing.T) {
	dn := "cn=user,ou=22307,dc=example,dc=com,ou=student"

	res := extractAcademicGroup(dn)

	require.NotNil(t, res)
	assert.Equal(t, "22307", *res)
}

// student + multiple OU (first valid numeric wins)
func TestExtractAcademicGroup_Student_MultipleOU_TakesFirstValid(t *testing.T) {
	dn := "ou=abc,ou=12345,ou=999,ou=student"

	res := extractAcademicGroup(dn)

	require.NotNil(t, res)
	assert.Equal(t, "12345", *res)
}

// whitespace + student
func TestExtractAcademicGroup_Student_WhitespaceHandling(t *testing.T) {
	dn := " cn=user , ou= 7777 , dc=example , ou=student "

	res := extractAcademicGroup(dn)

	require.NotNil(t, res)
	assert.Equal(t, "7777", *res)
}

func TestIsDigits_True(t *testing.T) {
	assert.True(t, isDigits("123456"))
}

func TestIsDigits_FalseLetters(t *testing.T) {
	assert.False(t, isDigits("12a34"))
}

func TestIsDigits_FalseEmpty(t *testing.T) {
	assert.False(t, isDigits(""))
}
