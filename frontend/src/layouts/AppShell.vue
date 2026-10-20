<template>
	<div class="flex h-screen w-screen overflow-hidden bg-slate-50 text-slate-800">
		<!-- Sidebar -->
		<aside
			class="relative flex flex-col border-r border-slate-200 bg-white transition-all duration-200 z-20 shrink-0 no-print"
			:class="isCollapsed ? 'w-16' : 'w-56'"
		>
			<!-- School Brand Header -->
			<div class="flex h-14 items-center justify-between px-3 border-b border-slate-100">
				<div v-if="!isCollapsed" class="flex items-center gap-2 overflow-hidden">
					<div class="w-8 h-8 rounded-lg bg-indigo-600 text-white flex items-center justify-center font-bold text-sm shrink-0 shadow-xs">
						{{ schoolInitials }}
					</div>
					<div class="truncate">
						<h1 class="text-xs font-bold text-slate-900 truncate">{{ authStore.school?.name || 'School Management' }}</h1>
						<span class="text-[10px] font-semibold text-indigo-600 uppercase tracking-wider">{{ authStore.school?.code || 'OFFLINE' }}</span>
					</div>
				</div>
				<div v-else class="mx-auto">
					<div class="w-8 h-8 rounded-lg bg-indigo-600 text-white flex items-center justify-center font-bold text-sm shadow-xs">
						{{ schoolInitials }}
					</div>
				</div>

				<button
					v-if="!isCollapsed"
					@click="isCollapsed = true"
					class="p-1 rounded-md text-slate-400 hover:text-slate-600 hover:bg-slate-100 transition"
					title="Collapse Sidebar"
				>
					<PanelLeftClose class="w-4 h-4" />
				</button>
			</div>

			<!-- Expand button when collapsed -->
			<div v-if="isCollapsed" class="py-2 flex justify-center border-b border-slate-100">
				<button
					@click="isCollapsed = false"
					class="p-1 rounded-md text-slate-400 hover:text-slate-600 hover:bg-slate-100 transition"
					title="Expand Sidebar"
				>
					<PanelLeftOpen class="w-4 h-4" />
				</button>
			</div>

			<!-- Navigation Links -->
			<nav class="flex-1 overflow-y-auto px-2 py-3 space-y-1">
				<router-link
					to="/"
					class="flex items-center gap-3 px-2.5 py-2 rounded-lg text-xs font-medium transition"
					:class="$route.path === '/' ? 'bg-indigo-50 text-indigo-700 font-semibold' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'"
				>
					<LayoutDashboard class="w-4 h-4 shrink-0" />
					<span v-if="!isCollapsed" class="truncate">Dashboard</span>
				</router-link>

				<div v-if="!isCollapsed" class="pt-2 pb-1 px-2.5 text-[10px] font-bold text-slate-400 uppercase tracking-wider">
					Students & Academics
				</div>

				<router-link
					v-if="authStore.hasPermission('students.view')"
					to="/students"
					class="flex items-center gap-3 px-2.5 py-2 rounded-lg text-xs font-medium transition"
					:class="$route.path.startsWith('/students') ? 'bg-indigo-50 text-indigo-700 font-semibold' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'"
				>
					<GraduationCap class="w-4 h-4 shrink-0" />
					<span v-if="!isCollapsed" class="truncate">Students</span>
				</router-link>

				<router-link
					v-if="authStore.hasPermission('admissions.view')"
					to="/admissions"
					class="flex items-center gap-3 px-2.5 py-2 rounded-lg text-xs font-medium transition"
					:class="$route.path.startsWith('/admissions') ? 'bg-indigo-50 text-indigo-700 font-semibold' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'"
				>
					<FileText class="w-4 h-4 shrink-0" />
					<span v-if="!isCollapsed" class="truncate">Admissions</span>
				</router-link>

				<router-link
					v-if="authStore.hasPermission('academics.view')"
					to="/academics"
					class="flex items-center gap-3 px-2.5 py-2 rounded-lg text-xs font-medium transition"
					:class="$route.path.startsWith('/academics') ? 'bg-indigo-50 text-indigo-700 font-semibold' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'"
				>
					<BookOpen class="w-4 h-4 shrink-0" />
					<span v-if="!isCollapsed" class="truncate">Academics</span>
				</router-link>

				<router-link
					v-if="authStore.hasPermission('academics.promote')"
					to="/promotions"
					class="flex items-center gap-3 px-2.5 py-2 rounded-lg text-xs font-medium transition"
					:class="$route.path.startsWith('/promotions') ? 'bg-indigo-50 text-indigo-700 font-semibold' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'"
				>
					<ArrowRightLeft class="w-4 h-4 shrink-0" />
					<span v-if="!isCollapsed" class="truncate">Promotions</span>
				</router-link>

				<div v-if="!isCollapsed" class="pt-2 pb-1 px-2.5 text-[10px] font-bold text-slate-400 uppercase tracking-wider">
					Billing & Finance
				</div>

				<router-link
					v-if="authStore.hasPermission('fees.view')"
					to="/fees/structure"
					class="flex items-center gap-3 px-2.5 py-2 rounded-lg text-xs font-medium transition"
					:class="$route.path.startsWith('/fees/structure') ? 'bg-indigo-50 text-indigo-700 font-semibold' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'"
				>
					<DollarSign class="w-4 h-4 shrink-0" />
					<span v-if="!isCollapsed" class="truncate">Fee Structures</span>
				</router-link>

				<router-link
					v-if="authStore.hasPermission('fees.view')"
					to="/fees/discounts"
					class="flex items-center gap-3 px-2.5 py-2 rounded-lg text-xs font-medium transition"
					:class="$route.path.startsWith('/fees/discounts') ? 'bg-indigo-50 text-indigo-700 font-semibold' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'"
				>
					<Percent class="w-4 h-4 shrink-0" />
					<span v-if="!isCollapsed" class="truncate">Concessions</span>
				</router-link>

				<router-link
					v-if="authStore.hasPermission('invoices.view')"
					to="/invoices"
					class="flex items-center gap-3 px-2.5 py-2 rounded-lg text-xs font-medium transition"
					:class="$route.path.startsWith('/invoices') ? 'bg-indigo-50 text-indigo-700 font-semibold' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'"
				>
					<Receipt class="w-4 h-4 shrink-0" />
					<span v-if="!isCollapsed" class="truncate">Invoices</span>
				</router-link>

				<router-link
					v-if="authStore.hasPermission('payments.view')"
					to="/payments"
					class="flex items-center gap-3 px-2.5 py-2 rounded-lg text-xs font-medium transition"
					:class="$route.path.startsWith('/payments') ? 'bg-indigo-50 text-indigo-700 font-semibold' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'"
				>
					<CreditCard class="w-4 h-4 shrink-0" />
					<span v-if="!isCollapsed" class="truncate">Payments</span>
				</router-link>

				<div v-if="!isCollapsed" class="pt-2 pb-1 px-2.5 text-[10px] font-bold text-slate-400 uppercase tracking-wider">
					Administration
				</div>

				<router-link
					v-if="authStore.hasPermission('settings.manage') || authStore.hasPermission('school.view')"
					to="/settings/school"
					class="flex items-center gap-3 px-2.5 py-2 rounded-lg text-xs font-medium transition"
					:class="$route.path.startsWith('/settings/school') ? 'bg-indigo-50 text-indigo-700 font-semibold' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'"
				>
					<Building2 class="w-4 h-4 shrink-0" />
					<span v-if="!isCollapsed" class="truncate">School Settings</span>
				</router-link>

				<router-link
					v-if="authStore.hasPermission('users.view')"
					to="/settings/users"
					class="flex items-center gap-3 px-2.5 py-2 rounded-lg text-xs font-medium transition"
					:class="$route.path.startsWith('/settings/users') ? 'bg-indigo-50 text-indigo-700 font-semibold' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'"
				>
					<Users class="w-4 h-4 shrink-0" />
					<span v-if="!isCollapsed" class="truncate">User Accounts</span>
				</router-link>

				<router-link
					v-if="authStore.hasPermission('roles.manage') || authStore.hasPermission('users.view')"
					to="/settings/roles"
					class="flex items-center gap-3 px-2.5 py-2 rounded-lg text-xs font-medium transition"
					:class="$route.path.startsWith('/settings/roles') ? 'bg-indigo-50 text-indigo-700 font-semibold' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'"
				>
					<ShieldCheck class="w-4 h-4 shrink-0" />
					<span v-if="!isCollapsed" class="truncate">Roles & Access</span>
				</router-link>

				<router-link
					v-if="authStore.hasPermission('audit.view')"
					to="/audit"
					class="flex items-center gap-3 px-2.5 py-2 rounded-lg text-xs font-medium transition"
					:class="$route.path === '/audit' ? 'bg-indigo-50 text-indigo-700 font-semibold' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'"
				>
					<History class="w-4 h-4 shrink-0" />
					<span v-if="!isCollapsed" class="truncate">Audit Logs</span>
				</router-link>
			</nav>

			<!-- Bottom Connection Badge & User Footer -->
			<div class="border-t border-slate-100 p-2">
				<div
					class="flex items-center gap-2 px-2 py-1.5 rounded-md text-[11px] font-medium text-emerald-700 bg-emerald-50 mb-2"
					:class="isCollapsed ? 'justify-center' : ''"
				>
					<span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
					<span v-if="!isCollapsed" class="truncate">Local SQLite (Active)</span>
				</div>

				<div class="flex items-center justify-between p-1.5 rounded-lg bg-slate-50">
					<div v-if="!isCollapsed" class="flex items-center gap-2 overflow-hidden">
						<div class="w-7 h-7 rounded-full bg-slate-200 text-slate-700 font-bold text-xs flex items-center justify-center shrink-0">
							{{ userInitials }}
						</div>
						<div class="truncate">
							<p class="text-xs font-semibold text-slate-800 truncate">{{ authStore.user?.displayName }}</p>
							<p class="text-[10px] text-slate-500 truncate capitalize">{{ authStore.user?.roles?.[0] || 'Staff' }}</p>
						</div>
					</div>
					<button
						@click="handleLogout"
						class="p-1.5 rounded-md text-slate-400 hover:text-rose-600 hover:bg-white transition"
						title="Sign Out"
					>
						<LogOut class="w-4 h-4" />
					</button>
				</div>

				<!-- Developer Branding Trigger -->
				<div class="mt-2 pt-2 border-t border-slate-100">
					<button
						@click="showDeveloperModal = true"
						class="w-full flex items-center gap-2 p-1.5 rounded-lg text-slate-600 hover:text-indigo-600 hover:bg-slate-50 transition text-left group"
						:class="isCollapsed ? 'justify-center' : ''"
						title="Architect & Developer: Hassaan Ali (@zfhassaan)"
					>
						<DeveloperLogo sizeClass="w-5 h-5 group-hover:scale-110 transition-transform" />
						<div v-if="!isCollapsed" class="truncate flex-1">
							<p class="text-[10px] font-bold text-slate-700 group-hover:text-indigo-600 transition truncate">zfhassaan</p>
							<p class="text-[9px] text-slate-400 truncate">Architect & Developer</p>
						</div>
						<Info v-if="!isCollapsed" class="w-3.5 h-3.5 text-slate-400 group-hover:text-indigo-500 transition shrink-0" />
					</button>
				</div>
			</div>
		</aside>

		<!-- Main Workspace Area -->
		<div class="flex-1 flex flex-col min-w-0 overflow-hidden">
			<!-- Header Bar -->
			<header class="h-14 bg-white border-b border-slate-200 flex items-center justify-between px-6 shrink-0 z-10 no-print">
				<div class="flex items-center gap-3">
					<h2 class="text-sm font-semibold text-slate-800 capitalize">{{ currentRouteTitle }}</h2>
				</div>
				<div class="flex items-center gap-4">
					<span class="text-xs text-slate-500 font-medium">
						{{ formattedDate }}
					</span>
				</div>
			</header>

			<!-- Page Content View -->
			<main class="flex-1 overflow-y-auto p-6 pb-16 print:p-0 print:overflow-visible">
				<router-view />
			</main>
		</div>

		<!-- Floating Bottom-Right Developer Badge -->
		<DeveloperBadge class="no-print" @open-modal="showDeveloperModal = true" />

		<!-- Mandatory Password Change Modal -->
		<ChangePasswordModal />

		<!-- Developer System Info Modal -->
		<DeveloperAboutModal :isOpen="showDeveloperModal" @close="showDeveloperModal = false" />
	</div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import ChangePasswordModal from '@/components/ChangePasswordModal.vue'
import DeveloperAboutModal from '@/components/DeveloperAboutModal.vue'
import DeveloperBadge from '@/components/DeveloperBadge.vue'
import DeveloperLogo from '@/components/DeveloperLogo.vue'
import {
	LayoutDashboard,
	GraduationCap,
	FileText,
	BookOpen,
	ArrowRightLeft,
	Building2,
	Users,
	ShieldCheck,
	History,
	PanelLeftClose,
	PanelLeftOpen,
	LogOut,
	DollarSign,
	Percent,
	Receipt,
	CreditCard,
	Info
} from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const isCollapsed = ref(false)
const showDeveloperModal = ref(false)

const schoolInitials = computed(() => {
	const code = authStore.school?.code
	if (code && code !== 'OFFLINE') {
		return code.length <= 3 ? code.toUpperCase() : code.substring(0, 2).toUpperCase()
	}
	const name = authStore.school?.name || 'SM'
	const words = name.trim().split(/\s+/)
	if (words.length >= 2) {
		return (words[0][0] + words[1][0]).toUpperCase()
	}
	return name.substring(0, 2).toUpperCase()
})

const userInitials = computed(() => {
	const name = authStore.user?.displayName || 'User'
	return name.substring(0, 2).toUpperCase()
})

const currentRouteTitle = computed(() => {
	return (route.meta?.title as string) || 'Dashboard'
})

const formattedDate = computed(() => {
	const d = new Date()
	return d.toLocaleDateString(undefined, { weekday: 'short', month: 'short', day: 'numeric', year: 'numeric' })
})

async function handleLogout() {
	await authStore.logout()
	router.push('/login')
}

onMounted(() => {
	authStore.refreshSchool()
})
</script>

