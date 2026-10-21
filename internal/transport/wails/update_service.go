package wails

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"school-management/internal/platform/version"
	"school-management/internal/transport/dto"
)

type UpdateInfoDTO struct {
	HasUpdate        bool   `json:"hasUpdate"`
	CurrentVersion   string `json:"currentVersion"`
	LatestVersion    string `json:"latestVersion"`
	ReleaseTitle     string `json:"releaseTitle"`
	ReleaseNotes     string `json:"releaseNotes"`
	PublishedAt      string `json:"publishedAt"`
	DownloadURL      string `json:"downloadURL"`
	DownloadFileName string `json:"downloadFileName"`
	AssetSizeBytes   int64  `json:"assetSizeBytes"`
	Repository       string `json:"repository"`
}

type UpdateService struct {
	httpClient *http.Client
}

func NewUpdateService() *UpdateService {
	return &UpdateService{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

type githubRelease struct {
	TagName     string        `json:"tag_name"`
	Name        string        `json:"name"`
	Body        string        `json:"body"`
	PublishedAt string        `json:"published_at"`
	Assets      []githubAsset `json:"assets"`
}

type githubAsset struct {
	Name               string `json:"name"`
	Size               int64  `json:"size"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// GetVersion returns the current running application version.
func (s *UpdateService) GetVersion() dto.Response[string] {
	return dto.Ok(version.Current)
}

// CheckForUpdates queries the GitHub Releases API for the specified repository (or default).
func (s *UpdateService) CheckForUpdates(customRepo string) dto.Response[UpdateInfoDTO] {
	repo := strings.TrimSpace(customRepo)
	if repo == "" {
		repo = version.DefaultGitHubRepo
	}

	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return dto.Fail[UpdateInfoDTO]("internal", fmt.Sprintf("Failed to create update request: %v", err), nil, "")
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", fmt.Sprintf("SchoolManagement-DesktopApp/%s", version.Current))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return dto.Fail[UpdateInfoDTO]("network_error", fmt.Sprintf("Could not connect to GitHub: %v", err), nil, "")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return dto.Fail[UpdateInfoDTO]("not_found", fmt.Sprintf("No releases found on GitHub repository '%s'. Publish a release on GitHub first.", repo), nil, "")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return dto.Fail[UpdateInfoDTO]("github_error", fmt.Sprintf("GitHub API returned status %d: %s", resp.StatusCode, string(body)), nil, "")
	}

	var rel githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return dto.Fail[UpdateInfoDTO]("internal", fmt.Sprintf("Failed to parse release information: %v", err), nil, "")
	}

	hasUpdate := version.IsNewerVersion(version.Current, rel.TagName)

	// Pick appropriate binary asset based on operating system
	var chosenAsset *githubAsset
	osTarget := runtime.GOOS

	switch osTarget {
	case "darwin":
		for i := range rel.Assets {
			name := strings.ToLower(rel.Assets[i].Name)
			if strings.Contains(name, "darwin") || strings.HasSuffix(name, ".dmg") || strings.HasSuffix(name, ".zip") {
				chosenAsset = &rel.Assets[i]
				break
			}
		}
	case "linux":
		for i := range rel.Assets {
			name := strings.ToLower(rel.Assets[i].Name)
			if strings.Contains(name, "linux") {
				chosenAsset = &rel.Assets[i]
				break
			}
		}
	default: // Windows and default
		// First prefer standalone executable for in-place replacement (retains all data safely)
		for i := range rel.Assets {
			name := strings.ToLower(rel.Assets[i].Name)
			if strings.HasSuffix(name, ".exe") && !strings.Contains(name, "setup") && !strings.Contains(name, "installer") {
				chosenAsset = &rel.Assets[i]
				break
			}
		}
		// Fallback to installer if standalone binary is not found
		if chosenAsset == nil {
			for i := range rel.Assets {
				name := strings.ToLower(rel.Assets[i].Name)
				if strings.HasSuffix(name, ".exe") {
					chosenAsset = &rel.Assets[i]
					break
				}
			}
		}
	}

	info := UpdateInfoDTO{
		HasUpdate:      hasUpdate,
		CurrentVersion: version.Current,
		LatestVersion:  strings.TrimPrefix(rel.TagName, "v"),
		ReleaseTitle:   rel.Name,
		ReleaseNotes:   rel.Body,
		PublishedAt:    rel.PublishedAt,
		Repository:     repo,
	}

	if chosenAsset != nil {
		info.DownloadURL = chosenAsset.BrowserDownloadURL
		info.DownloadFileName = chosenAsset.Name
		info.AssetSizeBytes = chosenAsset.Size
	}

	return dto.Ok(info)
}

// DownloadAndApplyUpdate downloads the update and applies the in-place patch or installer.
func (s *UpdateService) DownloadAndApplyUpdate(downloadURL, downloadFileName string) dto.Response[bool] {
	if downloadURL == "" {
		return dto.Fail[bool]("validation", "Download URL cannot be empty", nil, "")
	}

	exePath, err := os.Executable()
	if err != nil {
		return dto.Fail[bool]("internal", fmt.Sprintf("Could not locate running executable: %v", err), nil, "")
	}
	appDir := filepath.Dir(exePath)

	updatesDir := filepath.Join(appDir, "updates")
	if err := os.MkdirAll(updatesDir, 0o755); err != nil {
		return dto.Fail[bool]("internal", fmt.Sprintf("Failed to create updates directory: %v", err), nil, "")
	}

	fileName := downloadFileName
	if fileName == "" {
		fileName = "update_patch.exe"
	}
	targetFile := filepath.Join(updatesDir, fileName)

	// Download file
	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return dto.Fail[bool]("internal", err.Error(), nil, "")
	}
	req.Header.Set("User-Agent", fmt.Sprintf("SchoolManagement-DesktopApp/%s", version.Current))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return dto.Fail[bool]("network_error", fmt.Sprintf("Failed to download update: %v", err), nil, "")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dto.Fail[bool]("download_error", fmt.Sprintf("Download failed with HTTP %d", resp.StatusCode), nil, "")
	}

	tmpFile := targetFile + ".tmp"
	out, err := os.Create(tmpFile)
	if err != nil {
		return dto.Fail[bool]("internal", fmt.Sprintf("Failed to create temporary file: %v", err), nil, "")
	}

	if _, err := io.Copy(out, resp.Body); err != nil {
		out.Close()
		_ = os.Remove(tmpFile)
		return dto.Fail[bool]("internal", fmt.Sprintf("Failed to save update bytes: %v", err), nil, "")
	}
	out.Close()

	if err := os.Rename(tmpFile, targetFile); err != nil {
		_ = os.Remove(tmpFile)
		return dto.Fail[bool]("internal", fmt.Sprintf("Failed to finalize update file: %v", err), nil, "")
	}

	// Hot-patch or run installer
	lowerName := strings.ToLower(fileName)
	isInstaller := strings.Contains(lowerName, "setup") || strings.Contains(lowerName, "installer")

	go func() {
		time.Sleep(500 * time.Millisecond)

		if isInstaller {
			// Launch downloaded installer
			cmd := exec.Command(targetFile)
			_ = cmd.Start()
			os.Exit(0)
		} else {
			// In-place binary replacement (Windows allows renaming running .exe)
			oldExe := exePath + ".old"
			_ = os.Remove(oldExe)

			if err := os.Rename(exePath, oldExe); err == nil {
				// Copy new binary into place
				if copyErr := copyExecutableFile(targetFile, exePath); copyErr == nil {
					// Launch newly patched application
					cmd := exec.Command(exePath)
					cmd.Dir = appDir
					_ = cmd.Start()
					os.Exit(0)
				}
			}
		}
	}()

	return dto.Ok(true)
}

func copyExecutableFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
