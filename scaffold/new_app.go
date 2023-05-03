package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func NewApp(dirName string) error {
	err := download(dirName)
	if err != nil {
		return fmt.Errorf("failed to download repository: %w", err)
	}

	zipFile := filepath.Join(dirName, "plum-starter-app.zip")
	destDir := dirName
	err = unzip(zipFile, destDir)
	if err != nil {
		return fmt.Errorf("failed to unzip file: %w", err)
	}

	err = os.Remove(zipFile)
	if err != nil {
		return fmt.Errorf("failed to remove zip file: %w", err)
	}

	return nil
}

func download(dirName string) error {
	repo := "plum-starter-app"
	if dirName == "" {
		return fmt.Errorf("missing directory name argument")
	}

	url := fmt.Sprintf("https://github.com/scottraio/%s/archive/refs/heads/main.zip", repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "GoApp")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download repository: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to download repository: %s", resp.Status)
	}
	if err := os.MkdirAll(dirName, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	zipFile := filepath.Join(dirName, "plum-starter-app.zip")

	file, err := os.Create(zipFile)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save zip file: %w", err)
	}
	fmt.Printf("downloaded %s to %s\n", zipFile, dirName)
	return nil
}

func unzip(src, dest string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		path := filepath.Join(dest, strings.Replace(file.Name, "plum-starter-app-main/", "", 1))
		if file.FileInfo().IsDir() {
			err := os.MkdirAll(path, file.Mode())
			if err != nil {
				return err
			}
			continue
		}
		err := os.MkdirAll(filepath.Dir(path), file.Mode())
		if err != nil {
			return err
		}
		writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer writer.Close()

		reader, err := file.Open()
		if err != nil {
			return err
		}
		defer reader.Close()

		if _, err = io.Copy(writer, reader); err != nil {
			return err
		}
	}
	return nil
}
