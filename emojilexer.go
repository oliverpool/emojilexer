package twemoji

// NewLexer splits the text and the provided emojis from s.
func NewLexer(emojis []string) func(s string, text func(string), emoji func(string)) {
	root := newTree(emojis)

	return func(s string, text func(string), emoji func(string)) {
		textStart := 0 // leftmost byte. Everything before has been sent
		emojiStart := 0

		var j int
		var r rune
		current := root
		for j, r = range s {
			// is (still) matching an emoji
			if c := current.Children[r]; c != nil {
				// emoji starts now
				if current == root {
					emojiStart = j
				}
				current = c
				continue
			}

			// was matching an emoji
			if current != root {
				if current.IsLeaf {
					if textStart < emojiStart {
						text(s[textStart:emojiStart])
					}
					emoji(s[emojiStart:j])
					textStart = j
				}
				current = root
				// next emoji starts now
				if c := current.Children[r]; c != nil {
					emojiStart = j
					current = c
				}
			}
		}

		if j >= textStart {
			// currently matching an emoji
			if current != root && current.IsLeaf {
				if textStart < emojiStart {
					text(s[textStart:emojiStart])
				}
				emoji(s[emojiStart:])
			} else {
				text(s[textStart:])
			}
		}
	}
}

type node struct {
	IsLeaf   bool
	Children map[rune]*node
}

func (n *node) insert(s []rune) {
	if len(s) == 0 {
		n.IsLeaf = true
		return
	}
	if n.Children == nil {
		n.Children = make(map[rune]*node)
	}
	child := n.Children[s[0]]
	if child == nil {
		child = &node{}
		n.Children[s[0]] = child
	}
	child.insert(s[1:])
}

func newTree(emojis []string) *node {
	root := &node{
		IsLeaf:   false,
		Children: make(map[rune]*node),
	}
	for _, s := range emojis {
		root.insert([]rune(s))
	}
	return root
}
