# YAML
Fully backward-compatible YAML parser built from scratch in Golang
Yaml 1.2 syntax is specified in the official specification release at https://yaml.org/spec

## Backward Compatibility
This parser will start out to be fully compliant with YAML 1.2 specification (revision 2021)
https://yaml.org/spec/1.2.2/  
Compatibility with specifications 1.1.x will come in future releases.

## Key Features
- Supports all YAML data types: scalars, sequences, and mappings.
- Handles block style and flow style YAML formats.
- Supports advanced YAML features like anchors, aliases, and tags.
- Full support for multi-line and folded strings.
- Proper error handling for invalid YAML syntax, including indentation errors and invalid characters.

## Components
- **Tokenizer**: Breaks the input stream into YAML tokens such as scalars, mappings, sequences, comments, and indentation tokens.
- **Parser**: Constructs an Abstract Syntax Tree (AST) based on tokenized input.
- **Emitter**: Converts the parsed structure back into human-readable YAML if needed.

## Use cases
Ideal for configuration management, data serialization, and parsing of structured data in YAML format.