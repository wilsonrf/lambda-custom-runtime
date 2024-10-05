package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFile(url string, destPath string) error {
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file at %s: %w", destPath, err)
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file from %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status when downloading file: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy file contents to %s: %w", destPath, err)
	}

	if err := out.Chmod(0755); err != nil {
		return fmt.Errorf("failed to set executable permissions for %s: %w", destPath, err)
	}

	return nil
}

func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, srcInfo.Mode())
}
