/**
 * Utility to safely open external links in the default system browser
 * regardless of whether running inside Wails Desktop or a regular web browser.
 */
export function openExternalUrl(url: string): void {
	try {
		// Wails runtime injection
		if (typeof (window as any).runtime?.BrowserOpenURL === 'function') {
			(window as any).runtime.BrowserOpenURL(url)
			return
		}
	} catch {
		// Fallback to standard browser window
	}
	window.open(url, '_blank', 'noopener,noreferrer')
}
