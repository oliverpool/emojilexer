package twemoji

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func FromHexStringMust(t *testing.T, s string) string {
	out, err := FromHexString(s)
	require.NoError(t, err)
	return out
}
func TestToHexString(t *testing.T) {
	cc := []struct {
		name string
		in   string
		out  string
	}{
		{
			"empty in",
			"",
			"",
		},
		{
			"one letter",
			"a",
			"61",
		},
		{
			"some text",
			"ab",
			"61-62",
		},
		{
			"heart",
			"❤",
			"2764",
		},
		{
			"running woman",
			"🏃‍♀️",
			"1f3c3-200d-2640-fe0f",
		},
	}
	for _, c := range cc {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.out, ToHexString(c.in))
		})
	}
}
func TestFromHexString(t *testing.T) {
	cc := []struct {
		name string
		in   string
		out  string
	}{
		{
			"empty in",
			"",
			"",
		},
		{
			"one letter",
			"61",
			"a",
		},
		{
			"some text",
			"61-62",
			"ab",
		},
		{
			"39-20e3",
			"39-20e3",
			"9⃣",
		},
		{
			"a9",
			"a9",
			"©",
		},
		{
			"heart",
			"2764",
			"❤",
		},
		{
			"running woman",
			"1f3c3-200d-2640-fe0f",
			"🏃‍♀️",
		},
	}
	for _, c := range cc {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.out, FromHexStringMust(t, c.in))
		})
	}
}
