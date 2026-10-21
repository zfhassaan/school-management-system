export function updateFavicon(iconUrl: string) {
	try {
		let link = document.querySelector<HTMLLinkElement>("link[rel~='icon']")
		if (!link) {
			link = document.createElement('link')
			link.rel = 'icon'
			document.head.appendChild(link)
		}
		link.href = iconUrl
	} catch (e) {
		console.warn('Failed to update favicon', e)
	}
}

export function resetFavicon() {
	try {
		const link = document.querySelector<HTMLLinkElement>("link[rel~='icon']")
		if (link) {
			link.href = '/favicon.ico'
		}
	} catch (e) {
		console.warn('Failed to reset favicon', e)
	}
}
