package testdata

import "github.com/ercross/yaml/token"

var ScalarLineTokens = [][]token.Token{
	// strings: "Hello\, World\n"
	{
		token.New(token.TypeData, "string", 1, 1),
		token.New(token.TypeColon, "", 1, 7),
		token.New(token.TypeData, `Hello\, World`, 1, 9),
		token.New(token.TypeNewline, "", 1, 23),
	},

	{
		// integer: 12345
		token.New(token.TypeData, "integer", 2, 1),
		token.New(token.TypeColon, "", 2, 8),
		token.New(token.TypeData, "12345", 2, 10),
		token.New(token.TypeNewline, "", 2, 15),
	},

	{
		// float: 3.14159
		token.New(token.TypeData, "float", 3, 1),
		token.New(token.TypeColon, "", 3, 6),
		token.New(token.TypeData, "3.14159", 3, 8),
		token.New(token.TypeNewline, "", 3, 15),
	},

	{
		// boolean_true: true
		token.New(token.TypeData, "boolean_true", 4, 1),
		token.New(token.TypeColon, "", 4, 13),
		token.New(token.TypeData, "true", 4, 15),
		token.New(token.TypeNewline, "", 4, 19),
	},

	{
		// boolean_false: false
		token.New(token.TypeData, "boolean_false", 5, 1),
		token.New(token.TypeColon, "", 5, 14),
		token.New(token.TypeData, "false", 5, 16),
		token.New(token.TypeNewline, "", 5, 21),
	},
	{
		// null_value: null
		token.New(token.TypeData, "null_value", 6, 1),
		token.New(token.TypeColon, "", 6, 11),
		token.New(token.TypeData, "null", 6, 13),
		token.New(token.TypeNewline, "", 6, 17),
	},

	{
		// single_quote_string: 'This is YAML!'
		token.New(token.TypeData, "single_quote_string", 7, 1),
		token.New(token.TypeColon, "", 7, 20),
		token.New(token.TypeData, "This is YAML!", 7, 22),
		token.New(token.TypeNewline, "", 7, 36),
	},

	{
		// escaped_chars: "Line with a \"quote\" inside"
		token.New(token.TypeData, "escaped_chars", 8, 1),
		token.New(token.TypeColon, "", 8, 14),
		token.New(token.TypeData, `Line with a \"quote\" inside`, 8, 16),
		token.New(token.TypeNewline, "", 8, 45),
	},

	{
		// scientific: 1.23e4
		token.New(token.TypeData, "scientific", 9, 1),
		token.New(token.TypeColon, "", 9, 11),
		token.New(token.TypeData, "1.23e4", 9, 13),
	},
}
