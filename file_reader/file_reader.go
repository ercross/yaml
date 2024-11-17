package reader

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const maxFileSize = 100 * 1024 * 1024

func ReadFromYamlFile(filename string, out chan<- string) error {
	defer close(out)

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filename, err)
	}

	defer file.Close()
	if err = validateFile(file); err != nil {
		return fmt.Errorf("file failed validation %s: %w", filename, err)
	}

	reader := bufio.NewReader(file)
	var line string
	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				// this bufio.Reader rely on finding a delimiter ('\n' in this case) to determine the end of a line.
				// If the file's last line does not end with a newline character, it will not be recognized as
				// a complete line and may not return it.
				// This behavior is consistent with POSIX standards, where lines are expected to be
				// terminated with a newline.
				//
				// To ensure the last line is read even if it does not end with a newline,
				// explicitly check for EOF and process any remaining data
				if len(line) > 0 {
					out <- line
				}
				break
			}
			return fmt.Errorf("failed to read line: %w", err)
		}
		out <- line
	}

	return nil
}

func validateFile(file *os.File) error {
	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file %s: %w", file.Name(), err)
	}

	if stat.Size() > maxFileSize {
		return fmt.Errorf("file (%s) is too big", file.Name())
	}

	ext := filepath.Ext(file.Name())
	if !(ext == ".yaml" || ext == ".yml") {
		return fmt.Errorf("invalid file extension: %s", ext)
	}

	return nil
}
