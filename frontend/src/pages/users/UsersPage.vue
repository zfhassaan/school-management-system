<template>
	<div class="space-y-6 max-w-6xl mx-auto">
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-xl font-bold text-slate-900">User Accounts</h1>
				<p class="text-xs text-slate-500 mt-0.5">Manage administrative and staff accounts with offline authentication.</p>
			</div>
			<button
				@click="openCreateModal"
				class="py-2 px-3.5 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5"
			>
				<UserPlus class="w-4 h-4" />
				<span>Add New User</span>
			</button>
		</div>

		<!-- Users Table -->
		<div class="bg-white rounded-xl border border-slate-200 overflow-hidden shadow-xs">
			<table class="w-full text-left border-collapse">
				<thead>
					<tr class="bg-slate-50 border-b border-slate-200 text-[11px] font-bold uppercase tracking-wider text-slate-500">
						<th class="py-3 px-4">User</th>
						<th class="py-3 px-4">Email</th>
						<th class="py-3 px-4">Assigned Roles</th>
						<th class="py-3 px-4">Status</th>
						<th class="py-3 px-4">Last Login</th>
						<th class="py-3 px-4 text-right">Actions</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-slate-100 text-xs text-slate-700">
					<tr v-for="user in users" :key="user.id" class="hover:bg-slate-50/70 transition">
						<td class="py-3 px-4">
							<div class="font-semibold text-slate-900">{{ user.displayName }}</div>
							<div class="text-[11px] text-slate-500 font-mono">@{{ user.username }}</div>
						</td>
						<td class="py-3 px-4 text-slate-600">{{ user.email || '—' }}</td>
						<td class="py-3 px-4">
							<div class="flex flex-wrap gap-1">
								<span
									v-for="role in user.roles"
									:key="role"
									class="px-2 py-0.5 rounded-full text-[10px] font-medium bg-slate-100 text-slate-700 capitalize border border-slate-200"
								>
									{{ role.replace('_', ' ') }}
								</span>
								<span v-if="user.roles.length === 0" class="text-slate-400 text-[11px]">No roles</span>
							</div>
						</td>
						<td class="py-3 px-4">
							<span
								class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-semibold"
								:class="user.status === 'active' ? 'bg-emerald-50 text-emerald-700 border border-emerald-200' : 'bg-slate-100 text-slate-500 border border-slate-200'"
							>
								<span class="w-1.5 h-1.5 rounded-full" :class="user.status === 'active' ? 'bg-emerald-500' : 'bg-slate-400'"></span>
								{{ user.status }}
							</span>
						</td>
						<td class="py-3 px-4 text-slate-500 text-[11px]">
							{{ user.lastLoginAt ? new Date(user.lastLoginAt).toLocaleString() : 'Never' }}
						</td>
						<td class="py-3 px-4 text-right">
							<div class="flex items-center justify-end gap-1">
								<button
									@click="openEditModal(user)"
									class="p-1.5 text-slate-500 hover:text-indigo-600 rounded hover:bg-slate-100 transition"
									title="Edit User"
								>
									<Edit class="w-3.5 h-3.5" />
								</button>
								<button
									@click="openResetPwModal(user)"
									class="p-1.5 text-slate-500 hover:text-amber-600 rounded hover:bg-slate-100 transition"
									title="Reset Password"
								>
									<KeyRound class="w-3.5 h-3.5" />
								</button>
								<button
									@click="toggleUserStatus(user)"
									class="p-1.5 rounded hover:bg-slate-100 transition"
									:class="user.status === 'active' ? 'text-slate-500 hover:text-rose-600' : 'text-slate-500 hover:text-emerald-600'"
									:title="user.status === 'active' ? 'Disable Account' : 'Enable Account'"
								>
									<UserX v-if="user.status === 'active'" class="w-3.5 h-3.5" />
									<UserCheck v-else class="w-3.5 h-3.5" />
								</button>
							</div>
						</td>
					</tr>
				</tbody>
			</table>
		</div>

		<!-- Modal: Create / Edit User -->
		<div v-if="showUserModal" class="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/50 backdrop-blur-xs p-4">
			<div class="bg-white rounded-2xl shadow-xl max-w-md w-full p-6 border border-slate-100">
				<h3 class="text-base font-bold text-slate-900 mb-4">
					{{ isEditing ? 'Edit User Profile' : 'Create New User Account' }}
				</h3>
				<form @submit.prevent="saveUser" class="space-y-4">
					<div v-if="!isEditing">
						<label class="block text-xs font-semibold text-slate-700 mb-1">Username *</label>
						<input
							v-model="userForm.username"
							type="text"
							required
							placeholder="e.g. accountant1"
							class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white outline-hidden border-slate-300"
						/>
					</div>

					<div>
						<label class="block text-xs font-semibold text-slate-700 mb-1">Display Name *</label>
						<input
							v-model="userForm.displayName"
							type="text"
							required
							placeholder="e.g. Sarah Khan"
							class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white outline-hidden border-slate-300"
						/>
					</div>

					<div>
						<label class="block text-xs font-semibold text-slate-700 mb-1">Email Address</label>
						<input
							v-model="userForm.email"
							type="email"
							placeholder="sarah@school.internal"
							class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white outline-hidden border-slate-300"
						/>
					</div>

					<div v-if="!isEditing">
						<label class="block text-xs font-semibold text-slate-700 mb-1">Initial Password * (Min 10 chars)</label>
						<input
							v-model="userForm.password"
							type="password"
							required
							minlength="10"
							placeholder="••••••••••••"
							class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white outline-hidden border-slate-300"
						/>
					</div>

					<div>
						<label class="block text-xs font-semibold text-slate-700 mb-1">Assign Roles</label>
						<div class="space-y-1.5 max-h-36 overflow-y-auto p-2 bg-slate-50 rounded-lg border border-slate-200">
							<label
								v-for="r in availableRoles"
								:key="r.id"
								class="flex items-center gap-2 text-xs text-slate-800 cursor-pointer"
							>
								<input
									type="checkbox"
									:value="r.id"
									v-model="userForm.roleIds"
									class="rounded text-indigo-600 focus:ring-indigo-500"
								/>
								<span class="font-medium">{{ r.name }}</span>
								<span class="text-[10px] text-slate-400">({{ r.code }})</span>
							</label>
						</div>
					</div>

					<div class="flex justify-end gap-2 pt-3">
						<button
							type="button"
							@click="showUserModal = false"
							class="px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 font-medium rounded-lg text-xs"
						>
							Cancel
						</button>
						<button
							type="submit"
							class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg text-xs shadow-xs"
						>
							Save Account
						</button>
					</div>
				</form>
			</div>
		</div>

		<!-- Modal: Admin Reset Password -->
		<div v-if="showResetModal" class="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/50 backdrop-blur-xs p-4">
			<div class="bg-white rounded-2xl shadow-xl max-w-sm w-full p-6 border border-slate-100">
				<h3 class="text-base font-bold text-slate-900 mb-1">Reset Password</h3>
				<p class="text-xs text-slate-500 mb-4">Set a temporary password for <strong>{{ targetUser?.displayName }}</strong>. They will be forced to change it on their next sign-in.</p>

				<form @submit.prevent="submitResetPw" class="space-y-4">
					<div>
						<label class="block text-xs font-semibold text-slate-700 mb-1">New Temporary Password</label>
						<input
							v-model="resetPasswordVal"
							type="password"
							required
							minlength="10"
							placeholder="••••••••••••"
							class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white outline-hidden border-slate-300"
						/>
					</div>

					<div class="flex justify-end gap-2 pt-2">
						<button
							type="button"
							@click="showResetModal = false"
							class="px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 font-medium rounded-lg text-xs"
						>
							Cancel
						</button>
						<button
							type="submit"
							class="px-4 py-2 bg-amber-600 hover:bg-amber-700 text-white font-medium rounded-lg text-xs shadow-xs"
						>
							Update Password
						</button>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getServices } from '@/services/factory'
import { useToast } from '@/composables/useToast'
import type { UserListItemDTO, RoleDTO } from '@/types/common'
import { UserPlus, Edit, KeyRound, UserX, UserCheck } from 'lucide-vue-next'

const services = getServices()
const toast = useToast()

const users = ref<UserListItemDTO[]>([])
const availableRoles = ref<RoleDTO[]>([])

const showUserModal = ref(false)
const isEditing = ref(false)
const userForm = reactive({
	id: '',
	username: '',
	displayName: '',
	email: '',
	password: '',
	roleIds: [] as string[]
})

const showResetModal = ref(false)
const targetUser = ref<UserListItemDTO | null>(null)
const resetPasswordVal = ref('')

onMounted(async () => {
	await loadUsersAndRoles()
})

async function loadUsersAndRoles() {
	try {
		const [uRes, rRes] = await Promise.all([
			services.users.listUsers(),
			services.roles.listRoles()
		])
		if (uRes.success && uRes.data) users.value = uRes.data
		if (rRes.success && rRes.data) availableRoles.value = rRes.data
	} catch (e: any) {
		toast.error(e.message || 'Failed to load users')
	}
}

function openCreateModal() {
	isEditing.value = false
	userForm.id = ''
	userForm.username = ''
	userForm.displayName = ''
	userForm.email = ''
	userForm.password = ''
	userForm.roleIds = []
	showUserModal.value = true
}

function openEditModal(u: UserListItemDTO) {
	isEditing.value = true
	userForm.id = u.id
	userForm.username = u.username
	userForm.displayName = u.displayName
	userForm.email = u.email
	userForm.password = ''
	userForm.roleIds = [...u.roleIds]
	showUserModal.value = true
}

async function saveUser() {
	try {
		if (isEditing.value) {
			const res = await services.users.updateUser({
				id: userForm.id,
				displayName: userForm.displayName,
				email: userForm.email,
				roleIds: userForm.roleIds
			})
			if (res.success) {
				toast.success('User updated successfully')
				showUserModal.value = false
				await loadUsersAndRoles()
			} else {
				toast.error(res.error?.message || 'Update failed')
			}
		} else {
			const res = await services.users.createUser({
				username: userForm.username,
				displayName: userForm.displayName,
				email: userForm.email,
				password: userForm.password,
				roleIds: userForm.roleIds
			})
			if (res.success) {
				toast.success('User created successfully')
				showUserModal.value = false
				await loadUsersAndRoles()
			} else {
				toast.error(res.error?.message || 'Creation failed')
			}
		}
	} catch (e: any) {
		toast.error(e.message || 'Operation failed')
	}
}

async function toggleUserStatus(u: UserListItemDTO) {
	const newActive = u.status !== 'active'
	try {
		const res = await services.users.setUserStatus(u.id, newActive)
		if (res.success) {
			toast.success(`User ${newActive ? 'enabled' : 'disabled'} successfully`)
			await loadUsersAndRoles()
		} else {
			toast.error(res.error?.message || 'Failed to update user status')
		}
	} catch (e: any) {
		toast.error(e.message || 'Operation failed')
	}
}

function openResetPwModal(u: UserListItemDTO) {
	targetUser.value = u
	resetPasswordVal.value = ''
	showResetModal.value = true
}

async function submitResetPw() {
	if (!targetUser.value) return
	try {
		const res = await services.users.adminResetPassword({
			userId: targetUser.value.id,
			newPassword: resetPasswordVal.value
		})
		if (res.success) {
			toast.success('Temporary password set. User must change it at next login.')
			showResetModal.value = false
		} else {
			toast.error(res.error?.message || 'Reset password failed')
		}
	} catch (e: any) {
		toast.error(e.message || 'Operation failed')
	}
}
</script>

