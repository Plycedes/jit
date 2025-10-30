package repo

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Add(rootPath, targetPath string) error {
	jitDir := filepath.Join(rootPath, ".jit")
	indexPath := filepath.Join(jitDir, "index")

	if _, err := os.Stat(jitDir); os.IsNotExist(err) {
		return fmt.Errorf("fatal: not a jit repository")
	}

	indexEntries, _ := readIndex(indexPath)

	info, err := os.Stat(targetPath)
	if err != nil {
		return fmt.Errorf("cannot access %s: %v", targetPath, err)
	}

	if info.IsDir() {
		err = filepath.Walk(targetPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				if info.Name() == ".jit" {
					return filepath.SkipDir
				}
				return nil
			}
			return addSingleFile(rootPath, path, indexEntries)
		})

		if err != nil {
			return err
		}
	} else {
		if err := addSingleFile(rootPath, targetPath, indexEntries); err != nil {
			return err
		}
	}

	if err := writeIndex(indexPath, indexEntries); err != nil {
		return err
	}

	return nil
}

func addSingleFile(rootPath, filePath string, indexEntries map[string]string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", filePath, err)
	}

	// Creating jitDir again
	jitDir := filepath.Join(rootPath, ".jit")
	sha, err := HashObject(jitDir, data, "blob", true)
	if err != nil {
		return fmt.Errorf("failed to hash file %s: %v", filePath, err)
	}

	relPath, err := filepath.Rel(rootPath, filePath)
	if err != nil {
		return err
	}

	indexEntries[relPath] = sha
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
		return nil, err
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
