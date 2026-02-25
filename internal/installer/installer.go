package installer

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func getLatestRelease(repo string) (*githubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api returned %d", resp.StatusCode)
	}

	var rel githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, err
	}

	return &rel, nil
}

func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download %s: status %d", url, resp.StatusCode)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func InstallSops(destDir string) error {
	color.Cyan("Fetching latest sops release...")
	rel, err := getLatestRelease("getsops/sops")
	if err != nil {
		return fmt.Errorf("failed to get sops release: %w", err)
	}

	// sops-v3.8.1.linux.amd64
	osName := runtime.GOOS
	arch := runtime.GOARCH

	// sops mapping
	if osName == "darwin" && arch == "amd64" {
		// macos amd64 is often just 'darwin' or 'darwin.amd64'
	}

	var assetURL string
	expectedSuffix := fmt.Sprintf("%s.%s", osName, arch)
	if osName == "windows" {
		expectedSuffix = fmt.Sprintf("%s.exe", arch) // e.g. amd64.exe
	}

	for _, asset := range rel.Assets {
		if strings.HasSuffix(asset.Name, expectedSuffix) {
			assetURL = asset.BrowserDownloadURL
			break
		}
	}

	if assetURL == "" {
		return fmt.Errorf("no sops asset found for %s/%s", osName, arch)
	}

	binaryName := "sops"
	if osName == "windows" {
		binaryName = "sops.exe"
	}

	destPath := filepath.Join(destDir, binaryName)
	color.Cyan("Downloading sops %s to %s...", rel.TagName, destPath)

	if err := downloadFile(assetURL, destPath); err != nil {
		return fmt.Errorf("failed to download sops: %w", err)
	}

	// Make executable
	if err := os.Chmod(destPath, 0755); err != nil {
		return fmt.Errorf("failed to chmod sops: %w", err)
	}

	color.Green("sops installed successfully!")
	return nil
}

func InstallAge(destDir string) error {
	color.Cyan("Fetching latest age release...")
	rel, err := getLatestRelease("FiloSottile/age")
	if err != nil {
		return fmt.Errorf("failed to get age release: %w", err)
	}

	osName := runtime.GOOS
	arch := runtime.GOARCH

	// freebsd, linux, darwin -> .tar.gz
	// windows -> .zip

	var expectedSubstring string
	if osName == "windows" {
		expectedSubstring = fmt.Sprintf("windows-%s.zip", arch)
	} else {
		expectedSubstring = fmt.Sprintf("%s-%s.tar.gz", osName, arch)
	}

	var assetURL string
	var assetName string
	for _, asset := range rel.Assets {
		if strings.Contains(asset.Name, expectedSubstring) && !strings.HasSuffix(asset.Name, ".proof") {
			assetURL = asset.BrowserDownloadURL
			assetName = asset.Name
			break
		}
	}

	if assetURL == "" {
		return fmt.Errorf("no age asset found for %s/%s", osName, arch)
	}

	tempDir, err := os.MkdirTemp("", "age-install-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	archivePath := filepath.Join(tempDir, assetName)
	color.Cyan("Downloading age %s...", rel.TagName)
	if err := downloadFile(assetURL, archivePath); err != nil {
		return fmt.Errorf("failed to download age archive: %w", err)
	}

	color.Cyan("Extracting age...")

	binaryNames := []string{"age", "age-keygen"}
	if osName == "windows" {
		binaryNames = []string{"age.exe", "age-keygen.exe"}
	}

	if strings.HasSuffix(assetName, ".zip") {
		r, err := zip.OpenReader(archivePath)
		if err != nil {
			return err
		}
		defer r.Close()

		for _, f := range r.File {
			for _, bin := range binaryNames {
				if strings.HasSuffix(f.Name, bin) {
					rc, err := f.Open()
					if err != nil {
						return err
					}
					outPath := filepath.Join(destDir, bin)
					out, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
					if err != nil {
						rc.Close()
						return err
					}
					_, err = io.Copy(out, rc)
					out.Close()
					rc.Close()
					if err != nil {
						return err
					}
				}
			}
		}
	} else if strings.HasSuffix(assetName, ".tar.gz") {
		f, err := os.Open(archivePath)
		if err != nil {
			return err
		}
		defer f.Close()

		gzr, err := gzip.NewReader(f)
		if err != nil {
			return err
		}
		defer gzr.Close()

		tr := tar.NewReader(gzr)
		for {
			header, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			if header.Typeflag == tar.TypeReg {
				for _, bin := range binaryNames {
					if strings.HasSuffix(header.Name, "/"+bin) || header.Name == bin {
						outPath := filepath.Join(destDir, bin)
						out, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
						if err != nil {
							return err
						}
						_, err = io.Copy(out, tr)
						out.Close()
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}

	color.Green("age installed successfully!")
	return nil
}
