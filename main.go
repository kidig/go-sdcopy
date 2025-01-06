package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/djherbis/times"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: sdcopy <src_path> <dst_path>")
		os.Exit(1)
	}

	sourcePath := os.Args[1]
	destinationPath := os.Args[2]

	sem := make(chan string, 4)
	var wg sync.WaitGroup

	err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		sem <- info.Name()
		wg.Add(1)

		go func() {
			defer wg.Done()

			err := copyFile(path, destinationPath, info)
			if err != nil {
				fmt.Printf("Error copying file %s: %v\n", path, err)
			}

			<-sem
		}()

		return nil
	})

	if err != nil {
		fmt.Printf("Error scanning source path: %v\n", err)
		os.Exit(1)
	}

	wg.Wait()

	fmt.Println("Media files successfully copied.")
}

type placeholders struct {
	year  string
	month string
	day   string
}

func resolveDestinationPath(destinationPath string, d time.Time) string {
	pl := placeholders{
		year:  d.Format("2006"),
		month: d.Format("01"),
		day:   d.Format("02"),
	}

	resolved := destinationPath
	resolved = strings.ReplaceAll(resolved, "{year}", pl.year)
	resolved = strings.ReplaceAll(resolved, "{month}", pl.month)
	resolved = strings.ReplaceAll(resolved, "{day}", pl.day)

	re := regexp.MustCompile(`\{[^}]*\}`)
	resolved = re.ReplaceAllString(resolved, "")

	return resolved
}

func copyFile(sourcePath, destinationPath string, info os.FileInfo) error {
	ts := times.Get(info)
	mtime := info.ModTime()
	atime := ts.AccessTime()

	destDir := resolveDestinationPath(destinationPath, mtime)

	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", destDir, err)
	}
	destPath := filepath.Join(destDir, info.Name())

	fmt.Printf("Copying file: %s -> %s\n", sourcePath, destPath)

	if fileExists(destPath) {
		fmt.Printf("File already exists. skipping: %s\n", destPath)
	}

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %v", sourcePath, err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %v", destPath, err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copu file from %s to %s: %v", sourcePath, destPath, err)
	}

	if err := os.Chtimes(destPath, atime, mtime); err != nil {
		return fmt.Errorf("failed to preserve timestamps for %s: %v", destPath, err)
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
