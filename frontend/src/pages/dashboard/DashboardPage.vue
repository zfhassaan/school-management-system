<template>
	<div class="space-y-6 max-w-7xl mx-auto">
		<!-- Welcome Banner -->
		<div class="bg-gradient-to-r from-indigo-700 to-indigo-900 rounded-2xl p-6 text-white shadow-md relative overflow-hidden">
			<div class="relative z-10">
				<div class="inline-flex items-center gap-2 px-2.5 py-1 rounded-full bg-indigo-500/30 text-indigo-100 text-xs font-medium mb-3 backdrop-blur-xs">
					<span class="w-2 h-2 rounded-full bg-emerald-400"></span>
					Offline-First Desktop System Active
				</div>
				<h1 class="text-2xl font-bold tracking-tight">Welcome, {{ authStore.user?.displayName }}</h1>
				<p class="text-xs text-indigo-200 mt-1 max-w-xl">
					Managing <span class="font-semibold text-white">{{ authStore.school?.name }}</span> (Code: {{ authStore.school?.code }}).
					All financial and operational data is safely stored locally on this computer.
				</p>
			</div>
		</div>

		<!-- Status & Metrics Cards -->
		<div class="grid grid-cols-1 md:grid-cols-4 gap-4">
			<div class="bg-white p-5 rounded-xl border border-slate-200 shadow-xs">
				<div class="flex items-center justify-between text-slate-500 mb-2">
					<span class="text-xs font-semibold uppercase tracking-wider">Database Mode</span>
					<Database class="w-4 h-4 text-emerald-600" />
				</div>
				<div class="text-lg font-bold text-slate-900">SQLite (WAL)</div>
				<p class="text-[11px] text-emerald-700 mt-1 flex items-center gap-1">
					<CheckCircle2 class="w-3.5 h-3.5" /> High-concurrency WAL mode
				</p>
			</div>

			<div class="bg-white p-5 rounded-xl border border-slate-200 shadow-xs">
				<div class="flex items-center justify-between text-slate-500 mb-2">
					<span class="text-xs font-semibold uppercase tracking-wider">Active Staff Users</span>
					<Users class="w-4 h-4 text-indigo-600" />
				</div>
				<div class="text-lg font-bold text-slate-900">{{ usersCount }} Accounts</div>
				<p class="text-[11px] text-slate-500 mt-1">Argon2id offline authentication</p>
			</div>

			<div class="bg-white p-5 rounded-xl border border-slate-200 shadow-xs">
				<div class="flex items-center justify-between text-slate-500 mb-2">
					<span class="text-xs font-semibold uppercase tracking-wider">Access Roles</span>
					<ShieldCheck class="w-4 h-4 text-sky-600" />
				</div>
				<div class="text-lg font-bold text-slate-900">{{ rolesCount }} Roles</div>
				<p class="text-[11px] text-slate-500 mt-1">Fine-grained RBAC capabilities</p>
			</div>

			<div class="bg-white p-5 rounded-xl border border-slate-200 shadow-xs">
				<div class="flex items-center justify-between text-slate-500 mb-2">
					<span class="text-xs font-semibold uppercase tracking-wider">Audit Trail</span>
					<History class="w-4 h-4 text-amber-600" />
				</div>
				<div class="text-lg font-bold text-slate-900">{{ auditCount }} Events</div>
				<p class="text-[11px] text-slate-500 mt-1">Immutable trigger-enforced logs</p>
			</div>
		</div>

		<!-- Quick Navigation Grid -->
		<div class="space-y-3">
			<h3 class="text-xs font-bold uppercase tracking-wider text-slate-500">Quick Administration Actions</h3>
			<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
				<router-link
					to="/settings/users"
					class="p-4 bg-white rounded-xl border border-slate-200 hover:border-indigo-400 hover:shadow-xs transition group flex items-start gap-4"
				>
					<div class="w-10 h-10 rounded-lg bg-indigo-50 text-indigo-600 flex items-center justify-center shrink-0 group-hover:scale-105 transition">
						<Users class="w-5 h-5" />
					</div>
					<div>
						<h4 class="text-sm font-semibold text-slate-900 group-hover:text-indigo-600 transition">User Accounts</h4>
						<p class="text-xs text-slate-500 mt-0.5">Manage administrative logins, assign roles, and handle password resets.</p>
					</div>
				</router-link>

				<router-link
					to="/settings/roles"
					class="p-4 bg-white rounded-xl border border-slate-200 hover:border-indigo-400 hover:shadow-xs transition group flex items-start gap-4"
				>
					<div class="w-10 h-10 rounded-lg bg-sky-50 text-sky-600 flex items-center justify-center shrink-0 group-hover:scale-105 transition">
						<ShieldCheck class="w-5 h-5" />
					</div>
					<div>
						<h4 class="text-sm font-semibold text-slate-900 group-hover:text-indigo-600 transition">Roles & Permissions</h4>
						<p class="text-xs text-slate-500 mt-0.5">Configure role permissions, restrict sensitive actions, and inspect policies.</p>
					</div>
				</router-link>

				<router-link
					to="/settings/school"
					class="p-4 bg-white rounded-xl border border-slate-200 hover:border-indigo-400 hover:shadow-xs transition group flex items-start gap-4"
				>
					<div class="w-10 h-10 rounded-lg bg-emerald-50 text-emerald-600 flex items-center justify-center shrink-0 group-hover:scale-105 transition">
						<Building2 class="w-5 h-5" />
					</div>
					<div>
						<h4 class="text-sm font-semibold text-slate-900 group-hover:text-indigo-600 transition">School Settings</h4>
						<p class="text-xs text-slate-500 mt-0.5">Configure invoice prefixes, receipt sequences, and financial rounding rules.</p>
					</div>
				</router-link>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { getServices } from '@/services/factory'
import {
	Database,
	Users,
	ShieldCheck,
	History,
	CheckCircle2,
	Building2
} from 'lucide-vue-next'

const authStore = useAuthStore()
const services = getServices()

const usersCount = ref(0)
const rolesCount = ref(0)
const auditCount = ref(0)

onMounted(async () => {
	try {
		const [usersRes, rolesRes, auditRes] = await Promise.all([
			services.users.listUsers(),
			services.roles.listRoles(),
			services.audit.listLogs({ limit: 10 })
		])
		if (usersRes.success && usersRes.data) usersCount.value = usersRes.data.length
		if (rolesRes.success && rolesRes.data) rolesCount.value = rolesRes.data.length
		if (auditRes.success && auditRes.data) auditCount.value = auditRes.data.length
	} catch (e) {
		console.warn('Dashboard metrics fetch:', e)
	}
})
</script>

