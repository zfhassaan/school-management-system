import { ref } from 'vue'

export interface Toast {
	id: string
	type: 'success' | 'error' | 'info' | 'warning'
	title?: string
	message: string
	duration?: number
}

const toasts = ref<Toast[]>([])

export function useToast() {
	function show(toast: Omit<Toast, 'id'>) {
		const id = Math.random().toString(36).substring(2, 9)
		const duration = toast.duration ?? 4000
		const newToast: Toast = { ...toast, id }
		toasts.value.push(newToast)

		if (duration > 0) {
			setTimeout(() => {
				remove(id)
			}, duration)
		}
	}

	function success(message: string, title = 'Success') {
		show({ type: 'success', title, message })
	}

	function error(message: string, title = 'Error') {
		show({ type: 'error', title, message, duration: 6000 })
	}

	function info(message: string, title = 'Information') {
		show({ type: 'info', title, message })
	}

	function remove(id: string) {
		toasts.value = toasts.value.filter(t => t.id !== id)
	}

	return {
		toasts,
		show,
		success,
		error,
		info,
		remove
	}
}

