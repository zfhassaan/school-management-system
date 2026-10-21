package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

//go:embed school-management.exe
var appBinary []byte

const (
	appName      = "School Management System"
	appDirName   = "SchoolManagement"
	exeName      = "school-management.exe"
	appVersion   = "1.0.0"
	appPublisher = "School Management System"
)

func main() {
	fmt.Println("==========================================================")
	fmt.Printf("   %s - Setup & Installer v%s\n", appName, appVersion)
	fmt.Println("==========================================================")
	fmt.Println()

	// Default installation directory: %LOCALAPPDATA%\Programs\SchoolManagement
	defaultInstallDir := getDefaultInstallDir()

	fmt.Printf("Default installation folder:\n  %s\n\n", defaultInstallDir)
	fmt.Print("Press ENTER to install to default location, or type custom path: ")

	var customPath string
	fmt.Scanln(&customPath)
	customPath = strings.TrimSpace(customPath)

	installDir := defaultInstallDir
	if customPath != "" {
		installDir = customPath
	}

	installDir, err := filepath.Abs(installDir)
	if err != nil {
		fmt.Printf("Error resolving path: %v\n", err)
		pauseAndExit(1)
	}

	fmt.Printf("\nInstalling %s to: %s\n", appName, installDir)

	// 1. Create target directories (including data/ folder for encrypted DB)
	dirs := []string{
		installDir,
		filepath.Join(installDir, "data"),
		filepath.Join(installDir, "logs"),
		filepath.Join(installDir, "documents"),
		filepath.Join(installDir, "backups"),
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", d, err)
			pauseAndExit(1)
		}
	}

	// 2. Extract application executable
	targetExe := filepath.Join(installDir, exeName)
	fmt.Printf("Extracting application executable (%d MB)... ", len(appBinary)/(1024*1024))
	if err := os.WriteFile(targetExe, appBinary, 0o755); err != nil {
		fmt.Printf("Failed: %v\n", err)
		pauseAndExit(1)
	}
	fmt.Println("Done.")

	// 3. Create Desktop Shortcut
	desktopDir, err := getDesktopDir()
	if err == nil {
		shortcutPath := filepath.Join(desktopDir, appName+".lnk")
		if err := createWindowsShortcut(targetExe, shortcutPath, installDir, "School Management Desktop Application"); err == nil {
			fmt.Println("Created Desktop shortcut.")
		}
	}

	// 4. Create Start Menu Shortcut
	startMenuDir, err := getStartMenuDir()
	if err == nil {
		smsStartDir := filepath.Join(startMenuDir, appName)
		_ = os.MkdirAll(smsStartDir, 0o755)
		startShortcut := filepath.Join(smsStartDir, appName+".lnk")
		if err := createWindowsShortcut(targetExe, startShortcut, installDir, "School Management Desktop Application"); err == nil {
			fmt.Println("Created Start Menu shortcut.")
		}
	}

	// 5. Create Uninstaller script
	uninstallerBat := filepath.Join(installDir, "uninstall.bat")
	uninstallScript := fmt.Sprintf(`@echo off
echo Uninstalling %s...
taskkill /F /IM %s 2>nul
timeout /t 1 /nobreak >nul
del "%s" 2>nul
del "%s" 2>nul
echo %s has been uninstalled.
echo Your data in "%s\data" has been preserved for safety.
pause
`, appName, exeName,
		filepath.Join(desktopDir, appName+".lnk"),
		filepath.Join(startMenuDir, appName, appName+".lnk"),
		appName, installDir)
	_ = os.WriteFile(uninstallerBat, []byte(uninstallScript), 0o755)

	// 6. Register with Windows Add/Remove Programs (Registry)
	registerUninstallInRegistry(installDir, targetExe, uninstallerBat)

	fmt.Println()
	fmt.Println("==========================================================")
	fmt.Println("   Installation Completed Successfully!")
	fmt.Println("==========================================================")
	fmt.Println("Database and files will reside in the same installed folder:")
	fmt.Printf("  Database Vault: %s\n", filepath.Join(installDir, "data", "school.db.enc"))
	fmt.Printf("  Executable:     %s\n", targetExe)
	fmt.Println()
	fmt.Print("Launch School Management System now? (Y/n): ")

	var launchChoice string
	fmt.Scanln(&launchChoice)
	launchChoice = strings.TrimSpace(strings.ToLower(launchChoice))

	if launchChoice == "" || launchChoice == "y" || launchChoice == "yes" {
		cmd := exec.Command(targetExe)
		cmd.Dir = installDir
		_ = cmd.Start()
	}
}

func getDefaultInstallDir() string {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData != "" {
		return filepath.Join(localAppData, "Programs", appDirName)
	}
	return filepath.Join("C:\\", appDirName)
}

func getDesktopDir() (string, error) {
	userProfile := os.Getenv("USERPROFILE")
	if userProfile != "" {
		return filepath.Join(userProfile, "Desktop"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Desktop"), nil
}

func getStartMenuDir() (string, error) {
	appData := os.Getenv("APPDATA")
	if appData != "" {
		return filepath.Join(appData, "Microsoft", "Windows", "Start Menu", "Programs"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "AppData", "Roaming", "Microsoft", "Windows", "Start Menu", "Programs"), nil
}

// createWindowsShortcut creates a .lnk file using PowerShell WScript.Shell
func createWindowsShortcut(targetExe, shortcutPath, workDir, description string) error {
	psScript := fmt.Sprintf(
		`$ws = New-Object -ComObject WScript.Shell; $s = $ws.CreateShortcut('%s'); $s.TargetPath = '%s'; $s.WorkingDirectory = '%s'; $s.Description = '%s'; $s.IconLocation = '%s,0'; $s.Save()`,
		escapeForPowerShell(shortcutPath),
		escapeForPowerShell(targetExe),
		escapeForPowerShell(workDir),
		escapeForPowerShell(description),
		escapeForPowerShell(targetExe),
	)

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psScript)
	return cmd.Run()
}

func escapeForPowerShell(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

func registerUninstallInRegistry(installDir, targetExe, uninstallerBat string) {
	regCmd := fmt.Sprintf(
		`reg add "HKCU\Software\Microsoft\Windows\CurrentVersion\Uninstall\%s" /v DisplayName /d "%s" /f && `+
			`reg add "HKCU\Software\Microsoft\Windows\CurrentVersion\Uninstall\%s" /v DisplayVersion /d "%s" /f && `+
			`reg add "HKCU\Software\Microsoft\Windows\CurrentVersion\Uninstall\%s" /v Publisher /d "%s" /f && `+
			`reg add "HKCU\Software\Microsoft\Windows\CurrentVersion\Uninstall\%s" /v InstallLocation /d "%s" /f && `+
			`reg add "HKCU\Software\Microsoft\Windows\CurrentVersion\Uninstall\%s" /v DisplayIcon /d "%s" /f && `+
			`reg add "HKCU\Software\Microsoft\Windows\CurrentVersion\Uninstall\%s" /v UninstallString /d "\"%s\"" /f`,
		appDirName, appName,
		appDirName, appVersion,
		appDirName, appPublisher,
		appDirName, installDir,
		appDirName, targetExe,
		appDirName, uninstallerBat,
	)

	_ = exec.Command("cmd", "/c", regCmd).Run()
}

func pauseAndExit(code int) {
	fmt.Println("\nPress ENTER to exit...")
	var dummy string
	fmt.Scanln(&dummy)
	os.Exit(code)
}

// Suppress unused imports on Windows if needed
var _ = unsafe.Pointer(nil)
var _ = syscall.Handle(0)

