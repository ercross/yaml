package testdata

import "github.com/ercross/yaml/token"

// ExpectedMultilineTokens is the output expected after running tokenizer on multiline_strings.yaml file
var ExpectedMultilineTokens = [][]token.Token{

	// multiline1: |
	{
		token.New(token.TypeData, "multiline1", 1, 1),
		token.New(token.TypeColon, "", 1, 11),
		token.New(token.TypePipe, "", 1, 13),
		token.New(token.TypeNewline, "", 1, 14),
	},

	// This is a
	{

		token.New(token.TypeIndentation, "  ", 2, 1),
		token.New(token.TypeData, "This is a", 2, 3),
		token.New(token.TypeNewline, "", 2, 12),
	},

	// multiline string
	{
		token.New(token.TypeIndentation, "  ", 3, 1),
		token.New(token.TypeData, "multiline string", 3, 3),
		token.New(token.TypeNewline, "", 3, 19),
	},

	// that spans several lines.
	{
		token.New(token.TypeIndentation, "  ", 4, 1),
		token.New(token.TypeData, "that spans several lines.", 4, 3),
		token.New(token.TypeNewline, "", 4, 28),
	},

	{
		token.New(token.TypeNewline, "", 5, 1),
	},

	// multiline2: |
	{
		token.New(token.TypeData, "multiline2", 6, 1),
		token.New(token.TypeColon, "", 6, 11),
		token.New(token.TypePipe, "", 6, 13),
		token.New(token.TypeNewline, "", 6, 14),
	},

	// Another example of
	{
		token.New(token.TypeIndentation, "  ", 7, 1),
		token.New(token.TypeData, "Another example of", 7, 3),
		token.New(token.TypeNewline, "", 7, 21),
	},

	// a multiline string
	{
		token.New(token.TypeIndentation, "  ", 8, 1),
		token.New(token.TypeData, "a multiline string", 8, 3),
		token.New(token.TypeNewline, "", 8, 21),
	},

	// with newline retained
	{
		token.New(token.TypeIndentation, "  ", 9, 1),
		token.New(token.TypeData, "with newlines retained.", 9, 3),
		token.New(token.TypeNewline, "", 9, 26),
	},
}
