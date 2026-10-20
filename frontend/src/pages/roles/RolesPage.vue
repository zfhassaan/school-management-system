<template>
	<div class="space-y-6 max-w-6xl mx-auto">
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-xl font-bold text-slate-900">Roles & Permissions</h1>
				<p class="text-xs text-slate-500 mt-0.5">Control access rights, capabilities, and sensitive operation authorization.</p>
			</div>
			<button
				@click="openCreateModal"
				class="py-2 px-3.5 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5"
			>
				<Plus class="w-4 h-4" />
				<span>Create Custom Role</span>
			</button>
		</div>

		<!-- Roles Cards Grid -->
		<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
			<div
				v-for="role in roles"
				:key="role.id"
				class="bg-white rounded-xl border p-5 shadow-xs transition hover:shadow-sm"
				:class="selectedRole?.id === role.id ? 'border-indigo-500 ring-2 ring-indigo-100' : 'border-slate-200'"
			>
				<div class="flex items-start justify-between">
					<div>
						<div class="flex items-center gap-2">
							<h3 class="text-sm font-bold text-slate-900">{{ role.name }}</h3>
							<span
								v-if="role.isSystem"
								class="px-1.5 py-0.5 rounded text-[9px] font-semibold bg-slate-100 text-slate-600 border border-slate-200"
							>
								SYSTEM
							</span>
						</div>
						<p class="text-xs text-slate-500 mt-1 line-clamp-2">{{ role.description || 'No description' }}</p>
					</div>
				</div>

				<div class="mt-4 pt-3 border-t border-slate-100 flex items-center justify-between text-xs text-slate-500">
					<span>{{ role.permissions.length }} permissions</span>
					<span>{{ role.userCount }} users assigned</span>
				</div>

				<div class="mt-3 flex gap-2">
					<button
						@click="selectRole(role)"
						class="flex-1 py-1.5 px-3 bg-slate-50 hover:bg-slate-100 text-slate-700 rounded-md text-xs font-medium border border-slate-200 transition"
					>
						Inspect Permissions
					</button>
					<button
						v-if="!role.isSystem"
						@click="deleteRole(role)"
						class="py-1.5 px-2 bg-rose-50 hover:bg-rose-100 text-rose-600 rounded-md text-xs font-medium border border-rose-200 transition"
						title="Delete Role"
					>
						<Trash2 class="w-3.5 h-3.5" />
					</button>
				</div>
			</div>
		</div>

		<!-- Permission Matrix for Selected Role -->
		<div v-if="selectedRole" class="bg-white rounded-2xl border border-slate-200 p-6 shadow-xs space-y-4">
			<div class="flex items-center justify-between border-b pb-3">
				<div>
					<h3 class="text-base font-bold text-slate-900">
						Permissions for <span class="text-indigo-600">{{ selectedRole.name }}</span>
					</h3>
					<p class="text-xs text-slate-500">
						{{ selectedRole.isSystem && selectedRole.code === 'super_admin' ? 'Super Administrator always holds all catalog capabilities.' : 'Fine-grained permissions assigned to this role.' }}
					</p>
				</div>
				<button
					v-if="!selectedRole.isSystem || selectedRole.code !== 'super_admin'"
					@click="saveRolePermissions"
					:disabled="isSaving"
					class="py-1.5 px-4 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5"
				>
					<Save class="w-4 h-4" />
					<span>{{ isSaving ? 'Saving...' : 'Save Permissions' }}</span>
				</button>
			</div>

			<!-- Permissions grouped by module -->
			<div class="space-y-6">
				<div v-for="(perms, moduleName) in groupedPermissions" :key="moduleName" class="space-y-2">
					<h4 class="text-xs font-bold uppercase tracking-wider text-slate-500 capitalize">{{ moduleName }} Module</h4>
					<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2">
						<label
							v-for="p in perms"
							:key="p.code"
							class="flex items-start gap-2.5 p-2.5 rounded-lg border transition text-xs cursor-pointer"
							:class="activePermissions.includes(p.code) ? 'bg-indigo-50/60 border-indigo-200' : 'bg-slate-50/40 border-slate-200 hover:bg-slate-50'"
						>
							<input
								type="checkbox"
								:value="p.code"
								v-model="activePermissions"
								:disabled="selectedRole.isSystem && selectedRole.code === 'super_admin'"
								class="rounded text-indigo-600 focus:ring-indigo-500 mt-0.5"
							/>
							<div class="flex-1">
								<div class="font-semibold text-slate-800 flex items-center gap-1.5">
									<span>{{ p.code }}</span>
									<span v-if="p.sensitive" class="px-1 py-0.2 rounded text-[8px] font-bold bg-amber-100 text-amber-800">
										SENSITIVE
									</span>
								</div>
								<div class="text-[11px] text-slate-500 mt-0.5">{{ p.description }}</div>
							</div>
						</label>
					</div>
				</div>
			</div>
		</div>

		<!-- Modal: Create Role -->
		<div v-if="showCreateModal" class="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/50 backdrop-blur-xs p-4">
			<div class="bg-white rounded-2xl shadow-xl max-w-md w-full p-6 border border-slate-100">
				<h3 class="text-base font-bold text-slate-900 mb-4">Create Custom Role</h3>
				<form @submit.prevent="submitCreateRole" class="space-y-4">
					<div>
						<label class="block text-xs font-semibold text-slate-700 mb-1">Role Code * (lowercase_underscore)</label>
						<input
							v-model="createRoleForm.code"
							type="text"
							required
							placeholder="e.g. fee_auditor"
							class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white outline-hidden border-slate-300 font-mono text-xs"
						/>
					</div>

					<div>
						<label class="block text-xs font-semibold text-slate-700 mb-1">Display Name *</label>
						<input
							v-model="createRoleForm.name"
							type="text"
							required
							placeholder="e.g. Fee Auditor"
							class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white outline-hidden border-slate-300"
						/>
					</div>

					<div>
						<label class="block text-xs font-semibold text-slate-700 mb-1">Description</label>
						<textarea
							v-model="createRoleForm.description"
							rows="2"
							placeholder="Responsibilities and access scope..."
							class="w-full px-3 py-2 border rounded-lg text-sm bg-slate-50 focus:bg-white outline-hidden border-slate-300"
						></textarea>
					</div>

					<div class="flex justify-end gap-2 pt-2">
						<button
							type="button"
							@click="showCreateModal = false"
							class="px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 font-medium rounded-lg text-xs"
						>
							Cancel
						</button>
						<button
							type="submit"
							class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg text-xs shadow-xs"
						>
							Create Role
						</button>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { getServices } from '@/services/factory'
import { useToast } from '@/composables/useToast'
import type { RoleDTO, PermissionInfoDTO } from '@/types/common'
import { Plus, Trash2, Save } from 'lucide-vue-next'

const services = getServices()
const toast = useToast()

const roles = ref<RoleDTO[]>([])
const catalog = ref<PermissionInfoDTO[]>([])
const selectedRole = ref<RoleDTO | null>(null)
const activePermissions = ref<string[]>([])
const isSaving = ref(false)

const showCreateModal = ref(false)
const createRoleForm = reactive({
	code: '',
	name: '',
	description: ''
})

const groupedPermissions = computed(() => {
	const map: Record<string, PermissionInfoDTO[]> = {}
	for (const p of catalog.value) {
		if (!map[p.module]) map[p.module] = []
		map[p.module].push(p)
	}
	return map
})

onMounted(async () => {
	await loadData()
})

async function loadData() {
	try {
		const [rRes, pRes] = await Promise.all([
			services.roles.listRoles(),
			services.roles.listPermissions()
		])
		if (rRes.success && rRes.data) {
			roles.value = rRes.data
			if (!selectedRole.value && roles.value.length > 0) {
				selectRole(roles.value[0])
			} else if (selectedRole.value) {
				const updated = roles.value.find(r => r.id === selectedRole.value?.id)
				if (updated) selectRole(updated)
			}
		}
		if (pRes.success && pRes.data) catalog.value = pRes.data
	} catch (e: any) {
		toast.error(e.message || 'Failed to load roles')
	}
}

function selectRole(role: RoleDTO) {
	selectedRole.value = role
	activePermissions.value = [...role.permissions]
}

async function saveRolePermissions() {
	if (!selectedRole.value) return
	isSaving.value = true
	try {
		const res = await services.roles.updateRole({
			id: selectedRole.value.id,
			name: selectedRole.value.name,
			description: selectedRole.value.description,
			permissions: activePermissions.value
		})
		if (res.success) {
			toast.success(`Permissions updated for role ${selectedRole.value.name}`)
			await loadData()
		} else {
			toast.error(res.error?.message || 'Update failed')
		}
	} catch (e: any) {
		toast.error(e.message || 'An error occurred')
	} finally {
		isSaving.value = false
	}
}

function openCreateModal() {
	createRoleForm.code = ''
	createRoleForm.name = ''
	createRoleForm.description = ''
	showCreateModal.value = true
}

async function submitCreateRole() {
	try {
		const res = await services.roles.createRole({
			code: createRoleForm.code,
			name: createRoleForm.name,
			description: createRoleForm.description,
			permissions: []
		})
		if (res.success) {
			toast.success('Custom role created')
			showCreateModal.value = false
			await loadData()
		} else {
			toast.error(res.error?.message || 'Role creation failed')
		}
	} catch (e: any) {
		toast.error(e.message || 'An error occurred')
	}
}

async function deleteRole(role: RoleDTO) {
	if (!confirm(`Are you sure you want to delete role "${role.name}"?`)) return
	try {
		const res = await services.roles.deleteRole(role.id)
		if (res.success) {
			toast.success('Role deleted')
			if (selectedRole.value?.id === role.id) selectedRole.value = null
			await loadData()
		} else {
			toast.error(res.error?.message || 'Delete failed')
		}
	} catch (e: any) {
		toast.error(e.message || 'An error occurred')
	}
}
</script>

