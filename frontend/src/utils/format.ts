export function formatMinor(minor: number, currency = 'PKR'): string {
	const amount = (minor || 0) / 100
	return new Intl.NumberFormat('en-PK', {
		style: 'currency',
		currency: currency,
		minimumFractionDigits: 2,
		maximumFractionDigits: 2
	}).format(amount)
}

export function formatBasisPoints(bp: number): string {
	const pct = (bp || 0) / 100
	return `${pct.toFixed(2)}%`
}

export function formatCivilDate(dateVal: any): string {
	if (!dateVal) return '—'
	if (typeof dateVal === 'object' && dateVal.Year) {
		const y = String(dateVal.Year).padStart(4, '0')
		const m = String(dateVal.Month).padStart(2, '0')
		const d = String(dateVal.Day).padStart(2, '0')
		return `${y}-${m}-${d}`
	}
	if (typeof dateVal === 'string') {
		return dateVal.split('T')[0]
	}
	return String(dateVal)
}

