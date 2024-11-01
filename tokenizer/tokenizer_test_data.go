package tokenizer

import "github.com/ercross/yaml/token"

var ScalarLines = scalarLines

var scalarLines = []string{
	`string: "Hello\, World"` + "\n",
	"integer: 12345\n",
	"float: 3.14159\n",
	"boolean_true: true\n",
	"boolean_false: false\n",
	"null_value: null\n",
	"single_quote_string: 'This is YAML!'\n",
	`escaped_chars: "Line with a "quote" inside"` + "\n",
	"scientific: 1.23e4\n",
}

var expectedScalarTokens = [][]token.Token{
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
		token.New(token.TypeData, `Line with a "quote" inside`, 8, 16),
		token.New(token.TypeNewline, "", 8, 41),
	},

	{
		// scientific: 1.23e4
		token.New(token.TypeData, "scientific", 9, 1),
		token.New(token.TypeColon, "", 9, 11),
		token.New(token.TypeData, "1.23e4", 9, 13),
		token.New(token.TypeNewline, "", 9, 19),
	},
}

var AllNodes = []string{
	"# Example YAML 1.2 file",
	"# --- Document start",
	"---",
	"# Scalar types",
	"string: \"Hello, World\"",
	"integer: 12345",
	"float: 3.14159",
	"boolean_true: true",
	"boolean_false: false",
	"null_value: null",
	"multiline_string: |",
	"  This is a",
	"  multiline string.",
	"folded_string: >",
	"  This is a folded",
	"  string with",
	"  newlines replaced by spaces.",
	"# Sequences (lists)",
	"shopping_list:",
	"  - apples",
	"  - oranges",
	"  - bananas",
	"# Sequences with inline notation",
	"numbers: [1, 2, 3, 4, 5]",
	"# Mappings (dictionaries)",
	"person:",
	"  name: John Doe",
	"  age: 30",
	"  email: john.doe@example.com",
	"  address:",
	"    street: 123 Main St",
	"    city: Anytown",
	"    zip: 12345",
	"# Nested mappings with anchors and references",
	"defaults: &defaults",
	"  user: guest",
	"  permissions: read-only",
	"user1:",
	"  <<: *defaults",
	"  user: admin",
	"  permissions: read-write",
	"user2:",
	"  <<: *defaults",
	"# Complex keys",
	"? [home, work]",
	"  : phone_number: 555-1234",
	"# Multiple documents in a single file",
	"---",
	"document2:",
	"  content: \"This is another document\"",
	"---",
	"document3:",
	"  content: \"This is yet another document\"",
	"# ... Document end",
	"# YAML tags (explicit typing)",
	"binary_data: !!binary |",
	"  R0lGODlhAQABAIAAAAUEBAg=",
	"timestamp: !!timestamp 2024-09-08T12:34:56Z",
	"set: !!set",
	"  ? item1",
	"  ? item2",
}
