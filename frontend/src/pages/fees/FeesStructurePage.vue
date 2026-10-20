<template>
	<div class="space-y-6 max-w-6xl mx-auto">
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-xl font-bold text-slate-900">Fee Management</h1>
				<p class="text-xs text-slate-500 mt-0.5">Define fee heads, create class fee structures, and manage immutable versioned fee schedules.</p>
			</div>
			<div class="flex gap-2">
				<button
					v-if="currentTab === 'heads'"
					@click="openAddHeadModal"
					class="py-2 px-3.5 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5 cursor-pointer"
				>
					<Plus class="w-4 h-4" />
					<span>Add Fee Head</span>
				</button>
				<button
					v-if="currentTab === 'structures'"
					@click="showCreateStructureModal = true"
					class="py-2 px-3.5 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5 cursor-pointer"
				>
					<Plus class="w-4 h-4" />
					<span>Create Fee Structure</span>
				</button>
			</div>
		</div>

		<!-- Tabs -->
		<div class="flex border-b border-slate-200 space-x-6 text-xs font-medium">
			<button
				v-for="tab in ['heads', 'structures']"
				:key="tab"
				@click="currentTab = tab"
				class="pb-3 capitalize cursor-pointer transition border-b-2"
				:class="currentTab === tab ? 'border-indigo-600 text-indigo-600 font-semibold' : 'border-transparent text-slate-500 hover:text-slate-800'"
			>
				{{ tab === 'heads' ? 'Fee Heads (Catalog)' : 'Class Fee Structures' }}
			</button>
		</div>

		<!-- TAB 1: Fee Heads -->
		<div v-if="currentTab === 'heads'" class="space-y-4">
			<div class="bg-white rounded-xl border border-slate-200 overflow-hidden shadow-xs">
				<table class="w-full text-left border-collapse text-xs">
					<thead>
						<tr class="bg-slate-50 border-b border-slate-200 text-[11px] font-bold uppercase tracking-wider text-slate-500">
							<th class="py-3 px-4">Code</th>
							<th class="py-3 px-4">Name</th>
							<th class="py-3 px-4">Type</th>
							<th class="py-3 px-4">Default Amount</th>
							<th class="py-3 px-4">Status</th>
							<th class="py-3 px-4 text-right">Actions</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-slate-100 text-slate-700">
						<tr v-for="h in feeHeads" :key="h.ID" class="hover:bg-slate-50/50">
							<td class="py-3 px-4 font-mono font-medium text-slate-900">{{ h.Code }}</td>
							<td class="py-3 px-4">
								<div class="font-semibold text-slate-900">{{ h.Name }}</div>
								<div v-if="h.Description" class="text-[11px] text-slate-400 truncate max-w-xs">{{ h.Description }}</div>
							</td>
							<td class="py-3 px-4">
								<span
									class="px-2 py-0.5 rounded-full text-[10px] font-semibold"
									:class="h.IsRecurring ? 'bg-blue-50 text-blue-700 border border-blue-200' : 'bg-amber-50 text-amber-700 border border-amber-200'"
								>
									{{ h.IsRecurring ? 'Recurring' : 'One-time' }}
								</span>
							</td>
							<td class="py-3 px-4 font-medium text-slate-900">
								{{ formatMinor(h.DefaultAmountMinor) }}
							</td>
							<td class="py-3 px-4">
								<span
									class="px-2 py-0.5 rounded-full text-[10px] font-semibold"
									:class="h.IsActive ? 'bg-emerald-50 text-emerald-700 border border-emerald-200' : 'bg-slate-100 text-slate-500 border border-slate-200'"
								>
									{{ h.IsActive ? 'Active' : 'Inactive' }}
								</span>
							</td>
							<td class="py-3 px-4 text-right">
								<button
									@click="openEditHeadModal(h)"
									class="text-indigo-600 hover:text-indigo-800 font-medium text-xs cursor-pointer"
								>
									Edit
								</button>
							</td>
						</tr>
						<tr v-if="feeHeads.length === 0">
							<td colspan="6" class="py-8 text-center text-slate-400 italic">No fee heads configured yet. Click "Add Fee Head" to create one.</td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>

		<!-- TAB 2: Fee Structures & Versions -->
		<div v-if="currentTab === 'structures'" class="space-y-5">
			<!-- Academic Year Selector -->
			<div class="bg-white p-4 rounded-xl border border-slate-200 shadow-xs flex items-center justify-between">
				<div class="flex items-center gap-3">
					<label class="text-xs font-semibold text-slate-700">Academic Year:</label>
					<select
						v-model="selectedYearId"
						@change="loadStructures"
						class="text-xs border border-slate-300 rounded-lg py-1.5 px-3 bg-white text-slate-700 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
					>
						<option v-for="y in academicYears" :key="y.ID" :value="y.ID">{{ y.Name }} ({{ y.Code }})</option>
					</select>
				</div>
				<button
					@click="loadStructures"
					class="p-2 border border-slate-200 hover:bg-slate-50 rounded-lg text-slate-600 transition"
					title="Refresh"
				>
					<RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': loading }" />
				</button>
			</div>

			<!-- Structures Grid & Version Details Master-Detail Layout -->
			<div class="grid grid-cols-12 gap-5">
				<!-- Left column: Structures List -->
				<div class="col-span-4 space-y-3">
					<div class="text-xs font-bold uppercase tracking-wider text-slate-500">Fee Structures</div>
					<div
						v-for="s in structures"
						:key="s.ID"
						@click="selectStructure(s)"
						class="p-4 rounded-xl border transition cursor-pointer shadow-xs"
						:class="selectedStructure?.ID === s.ID ? 'bg-indigo-50/60 border-indigo-300' : 'bg-white border-slate-200 hover:border-slate-300'"
					>
						<div class="flex items-center justify-between">
							<div class="font-bold text-xs text-slate-900">{{ s.Name }}</div>
							<span class="font-mono text-[10px] text-slate-400">{{ s.Code }}</span>
						</div>
						<div class="text-[11px] text-slate-500 mt-1">Class: <span class="font-medium text-slate-700">{{ getClassName(s.ClassID) }}</span></div>
						<div v-if="s.Description" class="text-[11px] text-slate-400 mt-0.5 truncate">{{ s.Description }}</div>
					</div>
					<div v-if="structures.length === 0" class="p-6 bg-white rounded-xl border border-dashed border-slate-200 text-center text-xs text-slate-400 italic">
						No fee structures found for this academic year.
					</div>
				</div>

				<!-- Right column: Versions & Version Items -->
				<div class="col-span-8">
					<div v-if="selectedStructure" class="space-y-4">
						<div class="bg-white p-5 rounded-xl border border-slate-200 shadow-xs space-y-4">
							<div class="flex items-center justify-between border-b border-slate-100 pb-3">
								<div>
									<h2 class="text-sm font-bold text-slate-900">{{ selectedStructure.Name }}</h2>
									<p class="text-[11px] text-slate-500">Class: {{ getClassName(selectedStructure.ClassID) }} • Structure Code: {{ selectedStructure.Code }}</p>
								</div>
								<button
									@click="openAddVersionModal"
									class="py-1.5 px-3 bg-indigo-50 hover:bg-indigo-100 text-indigo-700 font-medium rounded-lg text-xs transition flex items-center gap-1 cursor-pointer"
								>
									<Plus class="w-3.5 h-3.5" />
									<span>Create Version</span>
								</button>
							</div>

							<!-- Version selector tabs -->
							<div class="flex gap-2 items-center flex-wrap">
								<span class="text-xs font-semibold text-slate-600 mr-1">Versions:</span>
								<button
									v-for="v in versions"
									:key="v.ID"
									@click="selectVersion(v.ID)"
									class="px-2.5 py-1 rounded-lg text-xs font-medium cursor-pointer transition flex items-center gap-1.5"
									:class="selectedVersion?.ID === v.ID ? 'bg-indigo-600 text-white shadow-xs' : 'bg-slate-100 hover:bg-slate-200 text-slate-700'"
								>
									<span>v{{ v.VersionNo }}</span>
									<span
										class="px-1.5 py-0.2 rounded-full text-[9px] uppercase tracking-wide font-bold"
										:class="v.Status === 'published' ? 'bg-emerald-100 text-emerald-800' : v.Status === 'draft' ? 'bg-amber-100 text-amber-800' : 'bg-slate-200 text-slate-600'"
									>
										{{ v.Status }}
									</span>
								</button>
								<span v-if="versions.length === 0" class="text-xs text-slate-400 italic">No versions created yet</span>
							</div>

							<!-- Selected Version Items Table -->
							<div v-if="selectedVersionDetail" class="space-y-3 pt-2">
								<div class="flex items-center justify-between bg-slate-50 p-3 rounded-lg border border-slate-100">
									<div>
										<div class="text-xs font-medium text-slate-700">
											Status:
											<span class="font-bold capitalize" :class="selectedVersionDetail.version.Status === 'published' ? 'text-emerald-600' : 'text-amber-600'">
												{{ selectedVersionDetail.version.Status }}
											</span>
											<span v-if="selectedVersionDetail.version.PublishedAt" class="text-[11px] text-slate-400 ml-2">
												(Published: {{ formatCivilDate(selectedVersionDetail.version.PublishedAt) }})
											</span>
										</div>
										<div v-if="selectedVersionDetail.version.Changelog" class="text-[11px] text-slate-500 mt-0.5">
											Notes: {{ selectedVersionDetail.version.Changelog }}
										</div>
									</div>

									<div v-if="selectedVersionDetail.version.Status === 'draft'">
										<button
											@click="publishVersion(selectedVersionDetail.version.ID)"
											class="py-1.5 px-3 bg-emerald-600 hover:bg-emerald-700 text-white font-medium rounded-lg text-xs transition flex items-center gap-1 cursor-pointer"
										>
											<Lock class="w-3.5 h-3.5" />
											<span>Publish & Lock Version</span>
										</button>
									</div>
									<div v-else class="text-[11px] font-semibold text-emerald-700 flex items-center gap-1">
										<CheckCircle2 class="w-3.5 h-3.5" />
										<span>Immutable (Locked)</span>
									</div>
								</div>

								<!-- Items Table -->
								<div class="border border-slate-200 rounded-lg overflow-hidden">
									<table class="w-full text-left text-xs">
										<thead>
											<tr class="bg-slate-50 border-b border-slate-200 text-[10px] font-bold uppercase tracking-wider text-slate-500">
												<th class="py-2.5 px-3">Fee Head</th>
												<th class="py-2.5 px-3">Amount</th>
												<th class="py-2.5 px-3">Frequency</th>
												<th class="py-2.5 px-3">Due Day</th>
												<th class="py-2.5 px-3">Optional</th>
											</tr>
										</thead>
										<tbody class="divide-y divide-slate-100">
											<tr v-for="item in selectedVersionDetail.items" :key="item.ID" class="hover:bg-slate-50/50">
												<td class="py-2.5 px-3 font-medium text-slate-800">
													{{ item.FeeHeadName || getFeeHeadName(item.FeeHeadID) }}
													<span class="font-mono text-[10px] text-slate-400 ml-1">({{ item.FeeHeadCode || getFeeHeadCode(item.FeeHeadID) }})</span>
												</td>
												<td class="py-2.5 px-3 font-bold text-slate-900">{{ formatMinor(item.AmountMinor) }}</td>
												<td class="py-2.5 px-3 capitalize text-slate-600">{{ item.Frequency }}</td>
												<td class="py-2.5 px-3 text-slate-600">Day {{ item.DueDay }}</td>
												<td class="py-2.5 px-3">
													<span
														class="px-1.5 py-0.5 rounded text-[10px] font-semibold"
														:class="item.IsOptional ? 'bg-amber-50 text-amber-700' : 'bg-slate-100 text-slate-600'"
													>
														{{ item.IsOptional ? 'Optional' : 'Mandatory' }}
													</span>
												</td>
											</tr>
											<tr v-if="selectedVersionDetail.items.length === 0">
												<td colspan="5" class="py-6 text-center text-slate-400 italic">No line items in this version schedule.</td>
											</tr>
										</tbody>
										<tfoot v-if="selectedVersionDetail.items.length > 0" class="bg-slate-50 border-t border-slate-200">
											<tr>
												<td class="py-2.5 px-3 font-bold text-slate-700">Total Mandatory Fee Schedule:</td>
												<td colspan="4" class="py-2.5 px-3 font-bold text-indigo-700 text-sm">
													{{ formatMinor(calculateVersionTotal(selectedVersionDetail.items)) }}
												</td>
											</tr>
										</tfoot>
									</table>
								</div>
							</div>
						</div>
					</div>

					<div v-else class="bg-white p-12 rounded-xl border border-dashed border-slate-200 text-center text-xs text-slate-400">
						Select a fee structure on the left to inspect and manage versions.
					</div>
				</div>
			</div>
		</div>

		<!-- MODAL: Add / Edit Fee Head -->
		<div v-if="showHeadModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
			<div class="bg-white rounded-xl shadow-xl border border-slate-200 max-w-md w-full overflow-hidden">
				<div class="p-4 border-b border-slate-100 flex items-center justify-between">
					<h3 class="font-bold text-sm text-slate-900">{{ isEditingHead ? 'Edit Fee Head' : 'Add Fee Head' }}</h3>
					<button @click="showHeadModal = false" class="text-slate-400 hover:text-slate-600 text-lg cursor-pointer">&times;</button>
				</div>
				<form @submit.prevent="submitHeadForm" class="p-4 space-y-3">
					<div v-if="headModalError" class="p-2.5 bg-rose-50 text-rose-700 border border-rose-200 rounded-lg text-xs">
						{{ headModalError }}
					</div>

					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">Fee Code *</label>
						<input
							v-model="headForm.code"
							:disabled="isEditingHead"
							type="text"
							required
							placeholder="e.g. TUITION, LAB, LIBRARY"
							class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 font-mono focus:outline-none focus:ring-2 focus:ring-indigo-500/20 disabled:bg-slate-100"
						/>
					</div>

					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">Name *</label>
						<input
							v-model="headForm.name"
							type="text"
							required
							placeholder="e.g. Monthly Tuition Fee"
							class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
						/>
					</div>

					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">Description</label>
						<textarea
							v-model="headForm.description"
							rows="2"
							placeholder="Short details..."
							class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
						></textarea>
					</div>

					<div class="grid grid-cols-2 gap-3">
						<div>
							<label class="block text-xs font-medium text-slate-700 mb-1">Default Amount ($) *</label>
							<input
								v-model.number="headForm.amountDecimal"
								type="number"
								step="0.01"
								min="0"
								required
								class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
							/>
						</div>
						<div class="pt-5 flex items-center gap-2">
							<input id="isRecurring" v-model="headForm.isRecurring" type="checkbox" class="rounded border-slate-300 text-indigo-600 focus:ring-indigo-500/20" />
							<label for="isRecurring" class="text-xs text-slate-700 select-none">Recurring Fee</label>
						</div>
					</div>

					<div v-if="isEditingHead" class="pt-2 flex items-center gap-2">
						<input id="isActive" v-model="headForm.isActive" type="checkbox" class="rounded border-slate-300 text-indigo-600 focus:ring-indigo-500/20" />
						<label for="isActive" class="text-xs text-slate-700 select-none">Active</label>
					</div>

					<div class="flex justify-end gap-2 pt-3 border-t border-slate-100">
						<button
							type="button"
							@click="showHeadModal = false"
							class="px-3 py-1.5 border border-slate-200 text-slate-600 rounded-lg text-xs hover:bg-slate-50 cursor-pointer"
						>
							Cancel
						</button>
						<button
							type="submit"
							:disabled="savingHead"
							class="px-4 py-1.5 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg text-xs font-medium transition cursor-pointer"
						>
							{{ savingHead ? 'Saving...' : 'Save Fee Head' }}
						</button>
					</div>
				</form>
			</div>
		</div>

		<!-- MODAL: Create Fee Structure -->
		<div v-if="showCreateStructureModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
			<div class="bg-white rounded-xl shadow-xl border border-slate-200 max-w-md w-full overflow-hidden">
				<div class="p-4 border-b border-slate-100 flex items-center justify-between">
					<h3 class="font-bold text-sm text-slate-900">Create Fee Structure</h3>
					<button @click="showCreateStructureModal = false" class="text-slate-400 hover:text-slate-600 text-lg cursor-pointer">&times;</button>
				</div>
				<form @submit.prevent="submitStructureForm" class="p-4 space-y-3">
					<div v-if="structureModalError" class="p-2.5 bg-rose-50 text-rose-700 border border-rose-200 rounded-lg text-xs">
						{{ structureModalError }}
					</div>

					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">Class *</label>
						<select
							v-model="structureForm.classId"
							required
							class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
						>
							<option value="" disabled>Select Class</option>
							<option v-for="c in classes" :key="c.ID" :value="c.ID">{{ c.Name }} ({{ c.Code }})</option>
						</select>
					</div>

					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">Structure Code *</label>
						<input
							v-model="structureForm.code"
							type="text"
							required
							placeholder="e.g. FS-G1-2026"
							class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 font-mono focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
						/>
					</div>

					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">Structure Name *</label>
						<input
							v-model="structureForm.name"
							type="text"
							required
							placeholder="e.g. Grade 1 Standard Fee Structure"
							class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
						/>
					</div>

					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">Description</label>
						<textarea
							v-model="structureForm.description"
							rows="2"
							placeholder="Notes or applicable criteria..."
							class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
						></textarea>
					</div>

					<div class="flex justify-end gap-2 pt-3 border-t border-slate-100">
						<button
							type="button"
							@click="showCreateStructureModal = false"
							class="px-3 py-1.5 border border-slate-200 text-slate-600 rounded-lg text-xs hover:bg-slate-50 cursor-pointer"
						>
							Cancel
						</button>
						<button
							type="submit"
							:disabled="savingStructure"
							class="px-4 py-1.5 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg text-xs font-medium transition cursor-pointer"
						>
							{{ savingStructure ? 'Creating...' : 'Create Structure' }}
						</button>
					</div>
				</form>
			</div>
		</div>

		<!-- MODAL: Create New Version with Items -->
		<div v-if="showCreateVersionModal" class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs flex items-center justify-center p-4 z-50">
			<div class="bg-white rounded-xl shadow-xl border border-slate-200 max-w-2xl w-full overflow-hidden">
				<div class="p-4 border-b border-slate-100 flex items-center justify-between">
					<h3 class="font-bold text-sm text-slate-900">New Fee Schedule Version ({{ selectedStructure?.Name }})</h3>
					<button @click="showCreateVersionModal = false" class="text-slate-400 hover:text-slate-600 text-lg cursor-pointer">&times;</button>
				</div>
				<form @submit.prevent="submitVersionForm" class="p-4 space-y-4">
					<div v-if="versionModalError" class="p-2.5 bg-rose-50 text-rose-700 border border-rose-200 rounded-lg text-xs">
						{{ versionModalError }}
					</div>

					<div>
						<label class="block text-xs font-medium text-slate-700 mb-1">Changelog / Notes *</label>
						<input
							v-model="versionForm.changelog"
							type="text"
							required
							placeholder="e.g. Initial schedule for 2026/2027 academic year"
							class="w-full px-3 py-1.5 text-xs rounded-lg border border-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20"
						/>
					</div>

					<div class="space-y-2">
						<div class="flex items-center justify-between">
							<span class="text-xs font-bold text-slate-700">Fee Items in Schedule</span>
							<button
								type="button"
								@click="addVersionItemRow"
								class="text-xs text-indigo-600 hover:text-indigo-800 font-medium flex items-center gap-1 cursor-pointer"
							>
								<Plus class="w-3.5 h-3.5" />
								<span>Add Item</span>
							</button>
						</div>

						<div class="space-y-2 max-h-64 overflow-y-auto pr-1">
							<div
								v-for="(row, idx) in versionForm.items"
								:key="idx"
								class="flex items-center gap-2 p-2.5 bg-slate-50 border border-slate-200 rounded-lg"
							>
								<div class="flex-1">
									<select
										v-model="row.feeHeadId"
										required
										class="w-full px-2 py-1 text-xs rounded border border-slate-300 bg-white"
									>
										<option value="" disabled>Select Fee Head</option>
										<option v-for="h in feeHeads" :key="h.ID" :value="h.ID">{{ h.Name }} ({{ h.Code }})</option>
									</select>
								</div>
								<div class="w-24">
									<input
										v-model.number="row.amountDecimal"
										type="number"
										step="0.01"
										min="0"
										placeholder="Amount"
										required
										class="w-full px-2 py-1 text-xs rounded border border-slate-300 bg-white"
									/>
								</div>
								<div class="w-28">
									<select v-model="row.frequency" class="w-full px-2 py-1 text-xs rounded border border-slate-300 bg-white">
										<option value="monthly">Monthly</option>
										<option value="termly">Termly</option>
										<option value="annual">Annual</option>
										<option value="one_time">One-time</option>
									</select>
								</div>
								<div class="w-20">
									<input
										v-model.number="row.dueDay"
										type="number"
										min="1"
										max="31"
										placeholder="Due day"
										class="w-full px-2 py-1 text-xs rounded border border-slate-300 bg-white"
									/>
								</div>
								<div class="flex items-center gap-1">
									<input :id="'opt-' + idx" v-model="row.isOptional" type="checkbox" class="rounded text-indigo-600" />
									<label :for="'opt-' + idx" class="text-[10px] text-slate-600 select-none">Opt</label>
								</div>
								<button
									type="button"
									@click="removeVersionItemRow(idx)"
									class="text-rose-500 hover:text-rose-700 text-xs px-1 cursor-pointer"
								>
									&times;
								</button>
							</div>
							<div v-if="versionForm.items.length === 0" class="text-center text-xs text-slate-400 py-4 italic">
								At least one fee item is required in the version schedule.
							</div>
						</div>
					</div>

					<div class="flex justify-end gap-2 pt-3 border-t border-slate-100">
						<button
							type="button"
							@click="showCreateVersionModal = false"
							class="px-3 py-1.5 border border-slate-200 text-slate-600 rounded-lg text-xs hover:bg-slate-50 cursor-pointer"
						>
							Cancel
						</button>
						<button
							type="submit"
							:disabled="savingVersion || versionForm.items.length === 0"
							class="px-4 py-1.5 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg text-xs font-medium transition cursor-pointer"
						>
							{{ savingVersion ? 'Saving...' : 'Save Draft Version' }}
						</button>
					</div>
				</form>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getServices } from '@/services/factory'
import { useToast } from '@/composables/useToast'
import type {
	FeeHeadDTO,
	FeeStructureDTO,
	FeeStructureVersionDTO,
	FeeStructureVersionDetailDTO,
	FeeStructureItemDTO,
	AcademicYearDTO,
	ClassDTO
} from '@/types/common'
import { formatMinor, formatCivilDate } from '@/utils/format'
import { Plus, RefreshCw, Lock, CheckCircle2 } from 'lucide-vue-next'

const services = getServices()
const toast = useToast()

const currentTab = ref('heads')
const loading = ref(false)

// Data states
const feeHeads = ref<FeeHeadDTO[]>([])
const academicYears = ref<AcademicYearDTO[]>([])
const classes = ref<ClassDTO[]>([])
const structures = ref<FeeStructureDTO[]>([])
const selectedYearId = ref('')
const selectedStructure = ref<FeeStructureDTO | null>(null)
const versions = ref<FeeStructureVersionDTO[]>([])
const selectedVersion = ref<FeeStructureVersionDTO | null>(null)
const selectedVersionDetail = ref<FeeStructureVersionDetailDTO | null>(null)

// Fee Head Modal
const showHeadModal = ref(false)
const isEditingHead = ref(false)
const savingHead = ref(false)
const headModalError = ref('')
const headForm = ref({
	id: '',
	code: '',
	name: '',
	description: '',
	isRecurring: true,
	amountDecimal: 0,
	isActive: true
})

// Create Structure Modal
const showCreateStructureModal = ref(false)
const savingStructure = ref(false)
const structureModalError = ref('')
const structureForm = ref({
	classId: '',
	code: '',
	name: '',
	description: ''
})

// Create Version Modal
const showCreateVersionModal = ref(false)
const savingVersion = ref(false)
const versionModalError = ref('')
const versionForm = ref({
	changelog: '',
	items: [] as {
		feeHeadId: string
		amountDecimal: number
		frequency: string
		dueDay: number
		isOptional: boolean
	}[]
})

onMounted(async () => {
	await Promise.all([loadFeeHeads(), loadAcademicData()])
})

async function loadFeeHeads() {
	try {
		const res = await services.fees.listFeeHeads()
		if (res.success && res.data) {
			feeHeads.value = res.data
		}
	} catch (e: any) {
		console.error('Failed to load fee heads', e)
	}
}

async function loadAcademicData() {
	try {
		const [yearsRes, classesRes] = await Promise.all([
			services.academics.listAcademicYears(),
			services.academics.listClasses()
		])
		if (yearsRes.success && yearsRes.data) {
			academicYears.value = yearsRes.data
			const active = yearsRes.data.find(y => y.Status === 'active')
			if (active) selectedYearId.value = active.ID
			else if (yearsRes.data.length > 0) selectedYearId.value = yearsRes.data[0].ID
		}
		if (classesRes.success && classesRes.data) {
			classes.value = classesRes.data
		}
		if (selectedYearId.value) {
			await loadStructures()
		}
	} catch (e: any) {
		console.error('Failed to load academic data', e)
	}
}

async function loadStructures() {
	if (!selectedYearId.value) return
	loading.value = true
	try {
		const res = await services.fees.listFeeStructures(selectedYearId.value)
		if (res.success && res.data) {
			structures.value = res.data
			if (structures.value.length > 0) {
				await selectStructure(structures.value[0])
			} else {
				selectedStructure.value = null
				versions.value = []
				selectedVersionDetail.value = null
			}
		}
	} catch (e: any) {
		console.error('Failed to load fee structures', e)
	} finally {
		loading.value = false
	}
}

async function selectStructure(st: FeeStructureDTO) {
	selectedStructure.value = st
	try {
		const res = await services.fees.listFeeStructureVersions(st.ID)
		if (res.success && res.data) {
			versions.value = res.data
			if (res.data.length > 0) {
				await selectVersion(res.data[0].ID)
			} else {
				selectedVersion.value = null
				selectedVersionDetail.value = null
			}
		}
	} catch (e: any) {
		console.error('Failed to load versions', e)
	}
}

async function selectVersion(versionId: string) {
	try {
		const res = await services.fees.getFeeStructureVersionDetail(versionId)
		if (res.success && res.data) {
			selectedVersionDetail.value = res.data
			selectedVersion.value = res.data.version
		}
	} catch (e: any) {
		console.error('Failed to get version detail', e)
	}
}

function openAddHeadModal() {
	isEditingHead.value = false
	headModalError.value = ''
	headForm.value = {
		id: '',
		code: '',
		name: '',
		description: '',
		isRecurring: true,
		amountDecimal: 100,
		isActive: true
	}
	showHeadModal.value = true
}

function openEditHeadModal(h: FeeHeadDTO) {
	isEditingHead.value = true
	headModalError.value = ''
	headForm.value = {
		id: h.ID,
		code: h.Code,
		name: h.Name,
		description: h.Description,
		isRecurring: h.IsRecurring,
		amountDecimal: (h.DefaultAmountMinor || 0) / 100,
		isActive: h.IsActive
	}
	showHeadModal.value = true
}

async function submitHeadForm() {
	savingHead.value = true
	headModalError.value = ''
	try {
		const amountMinor = Math.round(headForm.value.amountDecimal * 100)
		if (isEditingHead.value) {
			const res = await services.fees.updateFeeHead({
				id: headForm.value.id,
				name: headForm.value.name,
				description: headForm.value.description,
				isRecurring: headForm.value.isRecurring,
				defaultAmountMinor: amountMinor,
				isActive: headForm.value.isActive
			})
			if (!res.success) {
				headModalError.value = res.error?.message || 'Failed to update fee head'
				return
			}
		} else {
			const res = await services.fees.createFeeHead({
				code: headForm.value.code.toUpperCase(),
				name: headForm.value.name,
				description: headForm.value.description,
				isRecurring: headForm.value.isRecurring,
				defaultAmountMinor: amountMinor
			})
			if (!res.success) {
				headModalError.value = res.error?.message || 'Failed to create fee head'
				return
			}
		}
		showHeadModal.value = false
		await loadFeeHeads()
	} catch (e: any) {
		headModalError.value = e.message || 'An error occurred'
	} finally {
		savingHead.value = false
	}
}

async function submitStructureForm() {
	savingStructure.value = true
	structureModalError.value = ''
	try {
		const res = await services.fees.createFeeStructure({
			academicYearId: selectedYearId.value,
			classId: structureForm.value.classId,
			code: structureForm.value.code.toUpperCase(),
			name: structureForm.value.name,
			description: structureForm.value.description
		})
		if (!res.success) {
			structureModalError.value = res.error?.message || 'Failed to create fee structure'
			return
		}
		showCreateStructureModal.value = false
		structureForm.value = { classId: '', code: '', name: '', description: '' }
		await loadStructures()
	} catch (e: any) {
		structureModalError.value = e.message || 'An error occurred'
	} finally {
		savingStructure.value = false
	}
}

function openAddVersionModal() {
	versionModalError.value = ''
	versionForm.value = {
		changelog: 'Revision schedule',
		items: feeHeads.value.slice(0, 2).map(h => ({
			feeHeadId: h.ID,
			amountDecimal: (h.DefaultAmountMinor || 0) / 100,
			frequency: h.IsRecurring ? 'monthly' : 'one_time',
			dueDay: 5,
			isOptional: false
		}))
	}
	showCreateVersionModal.value = true
}

function addVersionItemRow() {
	const firstHead = feeHeads.value[0]
	versionForm.value.items.push({
		feeHeadId: firstHead ? firstHead.ID : '',
		amountDecimal: firstHead ? (firstHead.DefaultAmountMinor || 0) / 100 : 0,
		frequency: 'monthly',
		dueDay: 5,
		isOptional: false
	})
}

function removeVersionItemRow(idx: number) {
	versionForm.value.items.splice(idx, 1)
}

async function submitVersionForm() {
	if (!selectedStructure.value) return
	savingVersion.value = true
	versionModalError.value = ''
	try {
		const items = versionForm.value.items.map(it => ({
			feeHeadId: it.feeHeadId,
			amountMinor: Math.round(it.amountDecimal * 100),
			frequency: it.frequency,
			dueDay: it.dueDay,
			isOptional: it.isOptional
		}))

		const res = await services.fees.createFeeStructureVersion({
			feeStructureId: selectedStructure.value.ID,
			changelog: versionForm.value.changelog,
			items: items
		})
		if (!res.success) {
			versionModalError.value = res.error?.message || 'Failed to create version'
			return
		}
		showCreateVersionModal.value = false
		if (selectedStructure.value) {
			await selectStructure(selectedStructure.value)
		}
	} catch (e: any) {
		versionModalError.value = e.message || 'An error occurred'
	} finally {
		savingVersion.value = false
	}
}

async function publishVersion(versionId: string) {
	if (!confirm('Are you sure you want to publish this fee structure version? Once published, it becomes permanently immutable and will be used for student billing.')) {
		return
	}
	try {
		const res = await services.fees.publishFeeStructureVersion(versionId)
		if (res.success) {
			toast.success('Fee structure version published successfully and locked for billing.')
			await selectVersion(versionId)
		} else {
			toast.error(res.error?.message || 'Failed to publish version')
		}
	} catch (e: any) {
		toast.error(e.message || 'An error occurred while publishing version')
	}
}

function getClassName(classId: string) {
	const c = classes.value.find(cl => cl.ID === classId)
	return c ? c.Name : classId
}

function getFeeHeadName(headId: string) {
	const h = feeHeads.value.find(hd => hd.ID === headId)
	return h ? h.Name : headId
}

function getFeeHeadCode(headId: string) {
	const h = feeHeads.value.find(hd => hd.ID === headId)
	return h ? h.Code : ''
}

function calculateVersionTotal(items: FeeStructureItemDTO[]): number {
	return items
		.filter(it => !it.IsOptional)
		.reduce((acc, it) => acc + (it.AmountMinor || 0), 0)
}
</script>

