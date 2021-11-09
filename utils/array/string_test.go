package array

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContainsString(t *testing.T) {
	list := []string{
		"big",
		"brother",
		"is",
		"watching",
		"you",
	}
	s := "big"
	containsString := ContainsString(list, s)
	assert.True(t, containsString)

	s = "brother"
	containsString = ContainsString(list, s)
	assert.True(t, containsString)

	s = "is not"
	containsString = ContainsString(list, s)
	assert.False(t, containsString)

	s = "watching your"
	containsString = ContainsString(list, s)
	assert.False(t, containsString)
}

func TestRemoveString(t *testing.T) {
	list := []string{
		"big",
		"brother",
		"is",
		"watching",
		"you",
	}
	s := "big"
	result := RemoveString(list, s)
	assert.NotContains(t, result, s)

	s = "is not you"
	result = RemoveString(list, s)
	assert.Equal(t, list, result)
}
