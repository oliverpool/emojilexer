package emojilexer

import (
	"fmt"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/require"
)

func Example() {
	// List the supported emojis
	lexer := NewLexer([]string{
		"🤩",
		"🥳",
		"🏃‍♀️",
	})

	// And lex your input
	input := `Some text with emojis:🤩🥳...` +
		`Some unsupported emojis as well: 🤠💩🏃‍♀️(the last one is supported)`
	lexer(input, func(s string) {
		fmt.Println("TEXT:", s)
	}, func(s string) {
		fmt.Println("EMOJI:", s)
	})

	// Output:
	// TEXT: Some text with emojis:
	// EMOJI: 🤩
	// EMOJI: 🥳
	// TEXT: ...Some unsupported emojis as well: 🤠💩
	// EMOJI: 🏃‍♀️
	// TEXT: (the last one is supported)
}

func TestEmptyQuick(t *testing.T) {
	f := func(s string) bool {
		var strs []string
		var emojis []string
		NewLexer(nil)(s, func(sub string) {
			strs = append(strs, sub)
		}, func(sub string) {
			emojis = append(emojis, sub)
		})
		require.Empty(t, emojis, "@"+s+"@")
		require.Len(t, strs, 1, "@"+s+"@")
		require.Equal(t, s, strs[0])
		return len(strs) == 1 && strs[0] == s && len(emojis) == 0
	}
	require.NoError(t, quick.Check(f, nil))
}

func TestConcatenationInvariantQuick(t *testing.T) {
	f := func(s string, emojis []string) bool {
		var out string
		NewLexer(emojis)(s, func(text string) {
			out += text
		}, func(emoji string) {
			out += emoji
		})
		require.Equal(t, s, out)
		return s == out
	}
	require.NoError(t, quick.Check(f, nil))
}

func TestEmptySet(t *testing.T) {
	cc := []struct {
		name string
		in   string
		out  []string
	}{
		{
			"empty in",
			"",
			[]string{""},
		},
		{
			"one letter",
			"a",
			[]string{"a"},
		},
		{
			"some text",
			"some text with some emoji",
			[]string{"some text with some emoji"},
		},
	}
	for _, c := range cc {
		t.Run(c.name, func(t *testing.T) {
			var out []string
			NewLexer(nil)(c.in, func(s string) {
				out = append(out, s)
			}, func(s string) {
				require.Fail(t, "found emoji: %s", s)
			})
			require.Equal(t, c.out, out)
		})
	}
}

func FromHexStringsMust(t *testing.T, ss ...string) []string {
	var out []string
	for _, s := range ss {
		out = append(out, FromHexStringMust(t, s))
	}
	return out
}

func TestMatch(t *testing.T) {
	cc := []struct {
		name      string
		filenames []string
		in        string
		out       []string
	}{
		{
			"only man emoji",
			FromHexStringsMust(t, "1f3c3"), // running man
			"Here is a woman: 🏃‍♀️, she is running with a man: 🏃",
			[]string{"Here is a woman: ", ":)🏃", "\u200d♀️, she is running with a man: ", ":)🏃"},
		},
		{
			"only woman emoji",
			FromHexStringsMust(t, "1f3c3-200d-2640-fe0f"), // running woman
			"Here is a woman: 🏃‍♀️, she is running with a man: 🏃",
			[]string{"Here is a woman: ", ":)🏃‍♀️", ", she is running with a man: 🏃"},
		},
		{
			"both emojis",
			FromHexStringsMust(t, "1f3c3", "1f3c3-200d-2640-fe0f"), // running man + woman
			"Here is a woman: 🏃‍♀️, she is running with a man: 🏃",
			[]string{"Here is a woman: ", ":)🏃‍♀️", ", she is running with a man: ", ":)🏃"},
		},
		{
			"emojis near another",
			FromHexStringsMust(t, "1f3c3", "1f3c3-200d-2640-fe0f"), // running man + woman
			"🏃🏃‍♀️🏃🏃",
			[]string{":)🏃", ":)🏃‍♀️", ":)🏃", ":)🏃"},
		},
	}
	for _, c := range cc {
		t.Run(c.name, func(t *testing.T) {
			var out []string
			NewLexer(c.filenames)(c.in, func(s string) {
				out = append(out, s)
			}, func(s string) {
				out = append(out, ":)"+s)
			})
			require.Equal(t, c.out, out)

		})
	}
}
