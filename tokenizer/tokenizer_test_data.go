package tokenizer

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

var expectedScalarTokens = [][]Token{
	// strings: "Hello\, World\n"
	{
		NewToken(TokenTypeData, "string", 1, 1),
		NewToken(TokenTypeColon, "", 1, 7),
		NewToken(TokenTypeData, `Hello\, World`, 1, 9),
		NewToken(TokenTypeNewline, "", 1, 23),
	},

	{
		// integer: 12345
		NewToken(TokenTypeData, "integer", 2, 1),
		NewToken(TokenTypeColon, "", 2, 8),
		NewToken(TokenTypeData, "12345", 2, 10),
		NewToken(TokenTypeNewline, "", 2, 15),
	},

	{
		// float: 3.14159
		NewToken(TokenTypeData, "float", 3, 1),
		NewToken(TokenTypeColon, "", 3, 6),
		NewToken(TokenTypeData, "3.14159", 3, 8),
		NewToken(TokenTypeNewline, "", 3, 15),
	},

	{
		// boolean_true: true
		NewToken(TokenTypeData, "boolean_true", 4, 1),
		NewToken(TokenTypeColon, "", 4, 13),
		NewToken(TokenTypeData, "true", 4, 15),
		NewToken(TokenTypeNewline, "", 4, 19),
	},

	{
		// boolean_false: false
		NewToken(TokenTypeData, "boolean_false", 5, 1),
		NewToken(TokenTypeColon, "", 5, 14),
		NewToken(TokenTypeData, "false", 5, 16),
		NewToken(TokenTypeNewline, "", 5, 21),
	},
	{
		// null_value: null
		NewToken(TokenTypeData, "null_value", 6, 1),
		NewToken(TokenTypeColon, "", 6, 11),
		NewToken(TokenTypeData, "null", 6, 13),
		NewToken(TokenTypeNewline, "", 6, 17),
	},

	{
		// single_quote_string: 'This is YAML!'
		NewToken(TokenTypeData, "single_quote_string", 7, 1),
		NewToken(TokenTypeColon, "", 7, 20),
		NewToken(TokenTypeData, "This is YAML!", 7, 22),
		NewToken(TokenTypeNewline, "", 7, 36),
	},

	{
		// escaped_chars: "Line with a \"quote\" inside"
		NewToken(TokenTypeData, "escaped_chars", 8, 1),
		NewToken(TokenTypeColon, "", 8, 14),
		NewToken(TokenTypeData, `Line with a "quote" inside`, 8, 16),
		NewToken(TokenTypeNewline, "", 8, 41),
	},

	{
		// scientific: 1.23e4
		NewToken(TokenTypeData, "scientific", 9, 1),
		NewToken(TokenTypeColon, "", 9, 11),
		NewToken(TokenTypeData, "1.23e4", 9, 13),
		NewToken(TokenTypeNewline, "", 9, 19),
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
