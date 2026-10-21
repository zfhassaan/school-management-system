//go:build !windows

package platform

func SetWindowIconFromPNG(pngBytes []byte) error {
	return nil
}

func SetWindowIconFromBase64(base64Data string) error {
	return nil
}
