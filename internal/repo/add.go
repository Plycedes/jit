package repo

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func AddFile(rootPath, filePath string) error {
	jitDir := filepath.Join(rootPath, ".jit")
	indexPath := filepath.Join(jitDir, "index")

	if _, err := os.Stat(jitDir); os.IsNotExist(err) {
		return fmt.Errorf("fatal: not a jit repository")
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", filePath, err)
	}

	sha, err := HashObject(rootPath, data, "blob", true)
	if err != nil {
		return fmt.Errorf("failed to hash file %s: %v", filePath, err)
	}

	relPath, err := filepath.Rel(rootPath, filePath)
	if err != nil {
		return err
	}

	indexEntries, _ := readIndex(indexPath)
	indexEntries[relPath] = sha
	if err := writeIndex(indexPath, indexEntries); err != nil {
		return err
	}

	fmt.Printf("added '%s'\n", relPath)
	return nil
}

func readIndex(indexPath string) (map[string]string, error) {
	index := make(map[string]string)

	f, err := os.Open(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return index, nil
		}
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}

		index[parts[1]] = parts[0]
	}

	return index, scanner.Err()
}

func writeIndex(indexPath string, entries map[string]string) error {
	f, err := os.Create(indexPath)
	if err != nil {
		return fmt.Errorf("failed to write index: %v", err)
	}
	defer f.Close()

	for path, sha := range entries {
		line := fmt.Sprintf("%s %s\n", sha, path)
		if _, err := f.WriteString(line); err != nil {
			return err
		}
	}

	return nil
}
