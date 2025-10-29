package repo

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func HashObject(jitPath string, data []byte, typ string, write bool) (string, error) {
	header := fmt.Sprintf("%s %d\x00", typ, len(data))
	store := append([]byte(header), data...)

	h := sha1.Sum(store)
	sha := hex.EncodeToString(h[:])

	if write {
		if err := writeObjectFile(jitPath, sha, store); err != nil {
			return "", err
		}
	}

	return sha, nil
}

func writeObjectFile(jitPath, sha string, store []byte) error {
	objectsDir := filepath.Join(jitPath, "objects")
	if err := os.MkdirAll(objectsDir, 0755); err != nil {
		return fmt.Errorf("failed to create objects dir: %v", err)
	}

	dir := filepath.Join(objectsDir, sha[:2])
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create object subdir: %v", err)
	}

	objPath := filepath.Join(dir, sha[2:])
	if _, err := os.Stat(objPath); err == nil {
		return nil
	}

	f, err := os.OpenFile(objPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0444)
	if err != nil {
		return fmt.Errorf("failed to compress object: %v", err)
	}
	defer f.Close()

	zw := zlib.NewWriter(f)
	if _, err := zw.Write(store); err != nil {
		return fmt.Errorf("failed to compress object: %v", err)
	}
	if err := zw.Close(); err != nil {
		return fmt.Errorf("failed to finish compression: %v", err)
	}

	return nil
}

func ReadObject(jitPath, sha string) (string, []byte, error) {
	if len(sha) < 3 {
		return "", nil, fmt.Errorf("invalid sha")
	}

	objPath := filepath.Join(jitPath, "objects", sha[:2], sha[2:])
	f, err := os.Open(objPath)
	if err != nil {
		return "", nil, fmt.Errorf("failed to open object: %v", err)
	}
	defer f.Close()

	zr, err := zlib.NewReader(f)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create zlib reader: %v", err)
	}
	defer zr.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, zr); err != nil {
		return "", nil, fmt.Errorf("failed to decompress object: %v", err)
	}
	raw := buf.Bytes()

	nullIdx := bytes.IndexByte(raw, 0)
	if nullIdx < 0 {
		return "", nil, fmt.Errorf("malformed object (no header)")
	}

	header := string(raw[:nullIdx])
	var typ string
	_, err = fmt.Sscanf(header, "%s", &typ)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse header: %v", err)
	}

	content := raw[nullIdx+1:]
	return typ, content, nil
}
