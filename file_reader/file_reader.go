package reader

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

const maxFileSize = 100 * 1024 * 1024

func ReadFromYamlFile(filename string, out chan<- string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filename, err)
	}

	defer file.Close()
	if err = validateFile(file); err != nil {
		return fmt.Errorf("file failed validation %s: %w", filename, err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		out <- scanner.Text()
	}
	if err = scanner.Err(); err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	close(out)
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
