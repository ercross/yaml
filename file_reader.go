package yaml

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// maxFileSize
const maxFileSize = 10 * 1024 * 1024

func parseSingleFile(path string, maxNestingLevel int) (*abstractSyntaxTree, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", path, err)
	}

	defer file.Close()
	if err = validateFile(file); err != nil {
		return nil, fmt.Errorf("file failed validation %s: %w", path, err)
	}

	reader := bufio.NewReader(file)
	var line string
	lineNumber := 0
	p := newParser(maxNestingLevel)
	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("failed to read line %s: %w", line, err)
		}
		lineNumber++
		tokens, err := tokenizeLine(line, lineNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to tokenize line %d: %w", lineNumber, err)
		}
		if err = p.parseTokens(tokens, nil); err != nil {
			return nil, fmt.Errorf("token parser error: %w", err)
		}
	}

	return p.ast, nil
}

func validateFile(file *os.File) error {
	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file %s: %w", file.Name(), err)
	}

	if stat.Size() > maxFileSize {
		return fmt.Errorf("file %s is too big", file.Name())
	}

	ext := filepath.Ext(file.Name())
	if !(ext == "yaml" || ext == "yml") {
		return fmt.Errorf("invalid file extension: %s", ext)
	}

	return nil
}
