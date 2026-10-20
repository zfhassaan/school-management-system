<template>
	<div class="space-y-6 max-w-6xl mx-auto">
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-xl font-bold text-slate-900">Audit Trail</h1>
				<p class="text-xs text-slate-500 mt-0.5">Immutable record of security events, administrative updates, and financial actions.</p>
			</div>
			<button
				@click="loadLogs"
				class="py-2 px-3 bg-white hover:bg-slate-50 border border-slate-200 text-slate-700 font-medium rounded-lg text-xs transition shadow-2xs flex items-center gap-1.5"
			>
				<RefreshCw class="w-3.5 h-3.5" :class="isLoading ? 'animate-spin' : ''" />
				<span>Refresh</span>
			</button>
		</div>

		<!-- Filters Bar -->
		<div class="bg-white p-4 rounded-xl border border-slate-200 shadow-xs flex flex-wrap gap-4 items-center">
			<div class="w-48">
				<label class="block text-[10px] font-bold uppercase tracking-wider text-slate-400 mb-1">Entity Type</label>
				<select
					v-model="filterType"
					@change="loadLogs"
					class="w-full px-2.5 py-1.5 text-xs bg-slate-50 border border-slate-200 rounded-md outline-hidden focus:bg-white focus:ring-1 focus:ring-indigo-500"
				>
					<option value="">All Entities</option>
					<option value="user">User</option>
					<option value="role">Role</option>
					<option value="school">School</option>
					<option value="school_settings">School Settings</option>
				</select>
			</div>

			<div class="w-64">
				<label class="block text-[10px] font-bold uppercase tracking-wider text-slate-400 mb-1">Search Entity ID</label>
				<input
					v-model="filterEntityId"
					@keyup.enter="loadLogs"
					type="text"
					placeholder="Filter by Entity ID..."
					class="w-full px-2.5 py-1.5 text-xs bg-slate-50 border border-slate-200 rounded-md outline-hidden focus:bg-white focus:ring-1 focus:ring-indigo-500 font-mono"
				/>
			</div>
		</div>

		<!-- Audit Logs Table -->
		<div class="bg-white rounded-xl border border-slate-200 overflow-hidden shadow-xs">
			<table class="w-full text-left border-collapse">
				<thead>
					<tr class="bg-slate-50 border-b border-slate-200 text-[11px] font-bold uppercase tracking-wider text-slate-500">
						<th class="py-3 px-4">Timestamp</th>
						<th class="py-3 px-4">Action</th>
						<th class="py-3 px-4">Entity</th>
						<th class="py-3 px-4">Reason / Details</th>
						<th class="py-3 px-4">Source</th>
						<th class="py-3 px-4 text-right">Data Diff</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-slate-100 text-xs text-slate-700">
					<tr v-for="log in logs" :key="log.id" class="hover:bg-slate-50/70 transition">
						<td class="py-3 px-4 font-mono text-[11px] text-slate-500">
							{{ new Date(log.occurredAt).toLocaleString() }}
						</td>
						<td class="py-3 px-4">
							<span class="px-2 py-0.5 rounded text-[10px] font-semibold bg-indigo-50 text-indigo-700 border border-indigo-100 font-mono">
								{{ log.action }}
							</span>
						</td>
						<td class="py-3 px-4">
							<span class="font-medium text-slate-800 capitalize">{{ log.entityType }}</span>
							<span v-if="log.entityId" class="text-[10px] text-slate-400 font-mono ml-1.5">({{ log.entityId.substring(0, 8) }}...)</span>
						</td>
						<td class="py-3 px-4 text-slate-600">{{ log.reason || '—' }}</td>
						<td class="py-3 px-4">
							<span class="text-[10px] uppercase font-bold text-slate-400">{{ log.source }}</span>
						</td>
						<td class="py-3 px-4 text-right">
							<button
								v-if="log.before || log.after"
								@click="viewDiff(log)"
								class="text-indigo-600 hover:text-indigo-800 font-medium text-xs inline-flex items-center gap-1"
							>
								<Eye class="w-3.5 h-3.5" />
								<span>Inspect</span>
							</button>
							<span v-else class="text-slate-400 text-[11px]">—</span>
						</td>
					</tr>
					<tr v-if="logs.length === 0">
						<td colspan="6" class="py-8 text-center text-slate-400 text-xs">
							No audit logs match the current filter.
						</td>
					</tr>
				</tbody>
			</table>
		</div>

		<!-- Diff Modal -->
		<div v-if="selectedLog" class="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/50 backdrop-blur-xs p-4">
			<div class="bg-white rounded-2xl shadow-xl max-w-2xl w-full p-6 border border-slate-100 space-y-4">
				<div class="flex items-center justify-between border-b pb-3">
					<div>
						<h3 class="text-sm font-bold text-slate-900">Event Details: {{ selectedLog.action }}</h3>
						<p class="text-[11px] text-slate-500 font-mono">Correlation ID: {{ selectedLog.correlationId }}</p>
					</div>
					<button @click="selectedLog = null" class="p-1 text-slate-400 hover:text-slate-600 rounded">
						<X class="w-4 h-4" />
					</button>
				</div>

				<div class="grid grid-cols-2 gap-4">
					<div>
						<h4 class="text-xs font-semibold text-slate-600 mb-1">State Before</h4>
						<pre class="p-3 bg-slate-50 rounded-lg border border-slate-200 text-[11px] font-mono overflow-x-auto max-h-64">{{ selectedLog.before ? JSON.stringify(selectedLog.before, null, 2) : 'null' }}</pre>
					</div>
					<div>
						<h4 class="text-xs font-semibold text-slate-600 mb-1">State After</h4>
						<pre class="p-3 bg-slate-50 rounded-lg border border-slate-200 text-[11px] font-mono overflow-x-auto max-h-64">{{ selectedLog.after ? JSON.stringify(selectedLog.after, null, 2) : 'null' }}</pre>
					</div>
				</div>

				<div class="flex justify-end pt-2">
					<button
						@click="selectedLog = null"
						class="px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 font-medium rounded-lg text-xs"
					>
						Close
					</button>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getServices } from '@/services/factory'
import { useToast } from '@/composables/useToast'
import type { AuditLogDTO } from '@/types/common'
import { RefreshCw, Eye, X } from 'lucide-vue-next'

const services = getServices()
const toast = useToast()

const logs = ref<AuditLogDTO[]>([])
const isLoading = ref(false)
const filterType = ref('')
const filterEntityId = ref('')
const selectedLog = ref<AuditLogDTO | null>(null)

onMounted(async () => {
	await loadLogs()
})

async function loadLogs() {
	isLoading.value = true
	try {
		const res = await services.audit.listLogs({
			entityType: filterType.value || undefined,
			entityId: filterEntityId.value || undefined,
			limit: 50
		})
		if (res.success && res.data) {
			logs.value = res.data
		}
	} catch (e: any) {
		toast.error(e.message || 'Failed to load audit logs')
	} finally {
		isLoading.value = false
	}
}

function viewDiff(log: AuditLogDTO) {
	selectedLog.value = log
}
</script>

