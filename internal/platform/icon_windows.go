//go:build windows

package platform

import (
	"encoding/base64"
	"errors"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

var (
	user32                       = syscall.NewLazyDLL("user32.dll")
	procEnumWindows              = user32.NewProc("EnumWindows")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procIsWindowVisible          = user32.NewProc("IsWindowVisible")
	procSendMessageW             = user32.NewProc("SendMessageW")
	procCreateIconFromResourceEx = user32.NewProc("CreateIconFromResourceEx")
	procDestroyIcon              = user32.NewProc("DestroyIcon")
)

const (
	wmSetIcon = 0x0080
	iconSmall = 0
	iconBig   = 1
	lrDefault = 0x00000000
)

var currentHIcon syscall.Handle

// SetWindowIconFromPNG updates the Windows taskbar and titlebar icon at runtime from PNG data.
func SetWindowIconFromPNG(pngBytes []byte) error {
	if len(pngBytes) == 0 {
		return errors.New("empty icon data")
	}

	hwnd := findAppWindow()
	if hwnd == 0 {
		return errors.New("could not find application window handle")
	}

	// Create HICON from PNG bytes
	pBits := &pngBytes[0]
	ret, _, err := procCreateIconFromResourceEx.Call(
		uintptr(unsafe.Pointer(pBits)),
		uintptr(len(pngBytes)),
		uintptr(1),          // fIcon = TRUE
		uintptr(0x00030000), // dwVersion
		uintptr(0),          // cxDesired (0 = use image width)
		uintptr(0),          // cyDesired (0 = use image height)
		uintptr(lrDefault),
	)

	if ret == 0 {
		return err
	}

	hIcon := syscall.Handle(ret)

	// Send WM_SETICON to window
	// ICON_BIG (1) updates Windows Taskbar icon
	procSendMessageW.Call(uintptr(hwnd), uintptr(wmSetIcon), uintptr(iconBig), uintptr(hIcon))
	// ICON_SMALL (0) updates Titlebar / Alt-Tab icon
	procSendMessageW.Call(uintptr(hwnd), uintptr(wmSetIcon), uintptr(iconSmall), uintptr(hIcon))

	// Clean up previous icon to prevent resource leaks
	if currentHIcon != 0 {
		procDestroyIcon.Call(uintptr(currentHIcon))
	}
	currentHIcon = hIcon

	return nil
}

// SetWindowIconFromBase64 parses data URL or base64 string and updates the icon.
func SetWindowIconFromBase64(base64Data string) error {
	if base64Data == "" {
		return nil
	}
	// Strip data:image/...;base64, prefix if present
	idx := strings.Index(base64Data, ",")
	raw := base64Data
	if idx != -1 {
		raw = base64Data[idx+1:]
	}
	raw = strings.TrimSpace(raw)
	data, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return err
	}
	return SetWindowIconFromPNG(data)
}

// findAppWindow finds the visible top-level window for the current process
func findAppWindow() syscall.Handle {
	currentPID := uint32(os.Getpid())
	var targetHWND syscall.Handle

	cb := syscall.NewCallback(func(hwnd syscall.Handle, lparam uintptr) uintptr {
		var windowPID uint32
		procGetWindowThreadProcessId.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&windowPID)))

		if windowPID == currentPID {
			isVisible, _, _ := procIsWindowVisible.Call(uintptr(hwnd))
			if isVisible != 0 {
				targetHWND = hwnd
				return 0 // Stop enumerating
			}
		}
		return 1 // Continue enumerating
	})

	procEnumWindows.Call(cb, 0)
	return targetHWND
}
