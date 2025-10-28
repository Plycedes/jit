package repo

import (
	"fmt"
	"os"
	"path/filepath"
)

func InitRepo(path string) error {
	jitDir := filepath.Join(path, ".jit")

	if _, err := os.Stat(jitDir); !os.IsNotExist(err) {
		return fmt.Errorf("repository already exists")
	}

	dirs := []string{
		filepath.Join(jitDir, "objects"),
		filepath.Join(jitDir, "refs", "heads"),
		filepath.Join(jitDir, "refs", "tags"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create %s: %v", dir, err)
		}
	}

	// Create HEAD file
	headPath := filepath.Join(jitDir, "HEAD")
	headContent := []byte("ref: refs/heads/main\n")

	if err := os.WriteFile(headPath, headContent, 0644); err != nil {
		return fmt.Errorf("failed to create HEAD: %v", err)
	}

	fmt.Println("Intialized empty Jit repository in", jitDir)
	return nil
}
