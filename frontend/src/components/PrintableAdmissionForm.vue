<template>
	<div v-if="isOpen && application" class="fixed inset-0 bg-slate-900/60 backdrop-blur-xs flex items-center justify-center p-4 z-50 overflow-y-auto animate-in fade-in duration-150 print:static print:bg-transparent print:p-0 print:overflow-visible print:block">
		<div class="bg-white rounded-2xl shadow-2xl max-w-3xl w-full p-6 border border-slate-200 my-8 max-h-[95vh] overflow-y-auto print:max-h-none print:overflow-visible print:p-0 print:border-none print:shadow-none print:my-0 print:w-full print:max-w-none">
			<!-- Top Control Bar (Hidden on print) -->
			<div class="flex items-center justify-between pb-4 mb-4 border-b border-slate-200 no-print">
				<div class="flex items-center gap-2">
					<span class="px-2.5 py-1 rounded-full text-xs font-bold uppercase tracking-wider bg-indigo-50 text-indigo-700 border border-indigo-200">
						Official Admission Form
					</span>
					<span class="text-xs text-slate-500 font-mono">
						Ref: {{ application.ApplicationNo }}
					</span>
				</div>

				<div class="flex items-center gap-2">
					<label class="py-1.5 px-3 bg-white border border-slate-300 hover:bg-slate-50 text-slate-700 font-medium rounded-lg text-xs transition shadow-2xs flex items-center gap-1.5 cursor-pointer">
						<Upload class="w-3.5 h-3.5 text-slate-500" />
						<span>{{ candidatePhoto ? 'Change Candidate Photo' : 'Upload Candidate Photo' }}</span>
						<input type="file" accept="image/png,image/jpeg,image/webp" class="hidden" @change="onPhotoSelected" />
					</label>
					<button
						@click="printForm"
						class="py-1.5 px-3.5 bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg text-xs transition shadow-xs flex items-center gap-1.5 cursor-pointer"
					>
						<Printer class="w-3.5 h-3.5" />
						<span>Print Admission Form</span>
					</button>
					<button
						@click="$emit('close')"
						class="p-1.5 text-slate-400 hover:text-slate-600 rounded-lg hover:bg-slate-100 transition cursor-pointer"
						title="Close"
					>
						<X class="w-4 h-4" />
					</button>
				</div>
			</div>

			<!-- Printable Paper Form -->
			<div class="admission-paper bg-white p-6 md:p-8 space-y-4 text-slate-900 border border-slate-300 rounded-xl print:border-none print:space-y-3 relative overflow-hidden">
				<!-- Background Institutional Watermark (opacity 3%) -->
				<div class="absolute inset-0 flex items-center justify-center pointer-events-none opacity-[0.03] select-none">
					<img v-if="authStore.schoolLogo" :src="authStore.schoolLogo" alt="watermark" class="w-96 h-96 object-contain" />
					<div v-else class="text-9xl font-black text-slate-900 tracking-widest">{{ schoolInitials }}</div>
				</div>

				<!-- HEADER SECTION: Logo + Title + Candidate Photo -->
				<div class="flex items-start justify-between gap-4 border-b-2 border-slate-900 pb-3 relative z-10">
					<!-- Left: School Logo / Crest -->
					<div class="w-20 h-20 rounded-full border-2 border-slate-800 bg-white flex items-center justify-center font-black text-2xl shrink-0 overflow-hidden shadow-2xs">
						<img v-if="authStore.schoolLogo" :src="authStore.schoolLogo" alt="School Logo" class="w-full h-full object-contain p-1" />
						<div v-else class="text-center">
							<span class="text-indigo-950 font-black text-xl">{{ schoolInitials }}</span>
						</div>
					</div>

					<!-- Center: School Details & Document Title -->
					<div class="flex-1 text-center px-2">
						<h1 class="text-xl md:text-2xl font-black uppercase text-slate-900 tracking-tight leading-tight">
							{{ currentSchoolName }}
						</h1>
						<p class="text-[11px] font-semibold text-slate-600 uppercase tracking-wide mt-0.5">
							{{ currentSchoolLegalName || 'Center of Academic Excellence & Holistic Learning' }}
						</p>
						<p class="text-xs font-bold text-slate-800 uppercase tracking-widest mt-0.5">
							{{ currentSchoolCampus }}
						</p>
						<div class="inline-block mt-1.5 border-b-2 border-slate-900 px-4 pb-0.5">
							<h2 class="text-sm md:text-base font-black tracking-widest uppercase text-slate-900">
								ADMISSION FORM
							</h2>
						</div>
					</div>

					<!-- Right: Candidate Passport Photo Box -->
					<div class="w-24 h-28 border-2 border-slate-800 bg-slate-50 flex items-center justify-center shrink-0 relative overflow-hidden group shadow-2xs">
						<img v-if="candidatePhoto" :src="candidatePhoto" alt="Candidate Photo" class="w-full h-full object-cover" />
						<div v-else class="text-center p-1.5 text-slate-400">
							<User class="w-6 h-6 mx-auto mb-1 text-slate-300" />
							<p class="text-[8px] font-bold uppercase leading-tight text-slate-400">
								Candidate Photo
							</p>
							<p class="text-[7px] text-slate-400">(Passport Size)</p>
						</div>

						<!-- Hover trigger to upload photo (Hidden on print) -->
						<label class="absolute inset-0 bg-slate-900/60 opacity-0 group-hover:opacity-100 flex flex-col items-center justify-center text-white text-[9px] font-bold transition cursor-pointer no-print p-1 text-center">
							<Camera class="w-4 h-4 mb-1" />
							<span>{{ candidatePhoto ? 'Change' : 'Upload' }}</span>
							<input type="file" accept="image/png,image/jpeg,image/webp" class="hidden" @change="onPhotoSelected" />
						</label>
					</div>
				</div>

				<!-- SECTION 1: OFFICE USE & FINANCIALS -->
				<div class="space-y-1 relative z-10">
					<div class="bg-slate-900 text-white px-3 py-1 font-black text-[10px] uppercase tracking-wider flex items-center justify-between">
						<span>SECTION 1: OFFICE USE &amp; ACADEMIC CLASSIFICATION</span>
						<span class="font-mono text-[9px] font-bold text-slate-300">App No: {{ application.ApplicationNo }}</span>
					</div>
					<table class="w-full border-collapse border border-slate-400 text-xs">
						<tbody>
							<tr class="border-b border-slate-300">
								<td class="w-1/4 p-1.5 bg-slate-50 font-bold text-slate-700 border-r border-slate-300 text-[11px]">Admission Date:</td>
								<td class="w-1/4 p-1.5 font-medium border-r border-slate-300 text-[11px]">{{ formattedAdmissionDate }}</td>
								<td class="w-1/4 p-1.5 bg-slate-50 font-bold text-slate-700 border-r border-slate-300 text-[11px]">Admission / Roll #:</td>
								<td class="w-1/4 p-1.5 font-mono font-bold text-[11px]">{{ application.ApplicationNo }}</td>
							</tr>
							<tr class="border-b border-slate-300">
								<td class="p-1.5 bg-slate-50 font-bold text-slate-700 border-r border-slate-300 text-[11px]">Admission For Class:</td>
								<td class="p-1.5 font-bold text-indigo-900 border-r border-slate-300 text-[11px]">{{ application.ClassName || 'Unassigned' }}</td>
								<td class="p-1.5 bg-slate-50 font-bold text-slate-700 border-r border-slate-300 text-[11px]">Campus:</td>
								<td class="p-1.5 font-semibold text-[11px]">{{ currentSchoolCampus }}</td>
							</tr>
							<tr>
								<td class="p-1.5 bg-slate-50 font-bold text-slate-700 border-r border-slate-300 text-[10px]">Academic Session:</td>
								<td class="p-1.5 font-medium border-r border-slate-300 text-[11px]">{{ application.YearCode || '2026-2027' }}</td>
								<td class="p-1.5 bg-slate-50 font-bold text-slate-700 border-r border-slate-300 text-[10px]">Application Status:</td>
								<td class="p-1.5 font-bold uppercase text-[10px]" :class="application.Status === 'enrolled' ? 'text-emerald-700' : 'text-slate-700'">
									{{ application.Status.replace('_', ' ') }}
								</td>
							</tr>
						</tbody>
					</table>
				</div>

				<!-- SECTION 2: STUDENT PERSONAL INFO -->
				<div class="space-y-1 relative z-10">
					<div class="bg-slate-900 text-white px-3 py-1 font-black text-[10px] uppercase tracking-wider">
						SECTION 2: STUDENT PERSONAL INFO
					</div>
					<table class="w-full border-collapse border border-slate-400 text-xs">
						<tbody>
							<tr class="border-b border-slate-300">
								<td class="w-1/4 p-1.5 bg-slate-50 font-bold text-slate-700 border-r border-slate-300 text-[11px]">Student Full Name:</td>
								<td colspan="3" class="p-1.5 font-black text-slate-900 text-sm tracking-wide">
									{{ application.ApplicantFirstName }} {{ application.ApplicantLastName }}
								</td>
							</tr>
							<tr class="border-b border-slate-300">
								<td class="p-1.5 bg-slate-50 font-bold text-slate-700 border-r border-slate-300 text-[11px]">Date of Birth:</td>
								<td class="p-1.5 font-medium border-r border-slate-300 text-[11px]">{{ formattedDOB }}</td>
								<td class="p-1.5 bg-slate-50 font-bold text-slate-700 border-r border-slate-300 text-[11px]">Gender:</td>
								<td class="p-1.5 font-bold uppercase border-r border-slate-300 text-[11px]">{{ application.ApplicantGender || 'Unspecified' }}</td>
							</tr>
							<tr class="border-b border-slate-300">
								<td class="p-1.5 bg-slate-50 font-bold text-slate-700 border-r border-slate-300 text-[11px]">Residential Address:</td>
								<td colspan="3" class="p-1.5 font-medium text-[11px]">
									{{ application.Response?.address || application.Response?.homeAddress || 'As per guardian records' }}
								</td>
							</tr>
							<tr>
								<td class="p-1.5 bg-slate-50 font-bold text-slate-700 border-r border-slate-300 text-[11px]">Previous School:</td>
								<td class="p-1.5 font-medium border-r border-slate-300 text-[11px]">
									{{ application.Response?.previousSchool || application.Response?.notes || '—' }}
								</td>
								<td class="p-1.5 bg-slate-50 font-bold text-slate-700 border-r border-slate-300 text-[11px]">Medical Notes:</td>
								<td class="p-1.5 font-medium text-[11px]">
									{{ application.Response?.medicalHistory || 'NO' }}
								</td>
							</tr>
						</tbody>
					</table>
				</div>

				<!-- SECTION 3: PARENTS & GUARDIAN DETAILS -->
				<div class="space-y-1 relative z-10">
					<div class="bg-slate-900 text-white px-3 py-1 font-black text-[10px] uppercase tracking-wider">
						SECTION 3: PARENTS &amp; GUARDIAN DETAILS
					</div>
					<table class="w-full border-collapse border border-slate-400 text-xs">
						<tbody>
							<tr class="border-b border-slate-300">
								<td class="w-1/4 p-1.5 bg-slate-50 font-bold text-slate-700 border-r border-slate-300 text-[11px]">Father / Guardian Name:</td>
								<td class="w-1/4 p-1.5 font-bold text-slate-900 border-r border-slate-300 text-[11px]">
									{{ application.Response?.guardianName || application.Response?.fatherName || 'Guardian' }}
								</td>
								<td class="w-1/4 p-1.5 bg-slate-50 font-bold text-slate-700 border-r border-slate-300 text-[11px]">Mobile (WhatsApp):</td>
								<td class="w-1/4 p-1.5 font-mono font-bold text-[11px]">
									{{ application.GuardianPhoneNormalized || application.Response?.guardianPhone || '—' }}
								</td>
							</tr>
							<tr class="border-b border-slate-300">
								<td class="p-1.5 bg-slate-50 font-bold text-slate-700 border-r border-slate-300 text-[11px]">Guardian Email:</td>
								<td class="p-1.5 font-medium border-r border-slate-300 text-[11px] text-slate-700 truncate">
									{{ application.Response?.guardianEmail || '—' }}
								</td>
								<td class="p-1.5 bg-slate-50 font-bold text-slate-700 border-r border-slate-300 text-[11px]">Occupation / Business:</td>
								<td class="p-1.5 font-medium text-[11px]">
									{{ application.Response?.guardianOccupation || application.Response?.occupation || '—' }}
								</td>
							</tr>
							<tr class="border-b border-slate-300">
								<td class="p-1.5 bg-slate-50 font-bold text-slate-700 border-r border-slate-300 text-[11px]">Mother's Name:</td>
								<td class="p-1.5 font-semibold text-slate-800 border-r border-slate-300 text-[11px]">
									{{ application.Response?.motherName || '—' }}
								</td>
								<td class="p-1.5 bg-slate-50 font-bold text-slate-700 border-r border-slate-300 text-[11px]">Emergency Phone:</td>
								<td class="p-1.5 font-mono text-[11px]">
									{{ application.Response?.motherPhone || application.Response?.emergencyPhone || '—' }}
								</td>
							</tr>
						</tbody>
					</table>
				</div>

				<!-- SECTION 4: SIBLINGS & EMERGENCY REFERENCE -->
				<div class="space-y-1 relative z-10">
					<div class="bg-slate-900 text-white px-3 py-1 font-black text-[10px] uppercase tracking-wider">
						SECTION 4: SIBLINGS &amp; EMERGENCY REFERENCE
					</div>
					<table class="w-full border-collapse border border-slate-400 text-xs">
						<tbody>
							<tr class="border-b border-slate-300">
								<td class="w-1/4 p-1.5 bg-slate-50 font-bold text-slate-700 border-r border-slate-300 text-[11px]">Reference Contact:</td>
								<td class="w-1/4 p-1.5 font-medium border-r border-slate-300 text-[11px]">
									{{ application.Response?.referenceName || 'Direct Admission' }}
								</td>
								<td class="w-1/4 p-1.5 bg-slate-50 font-bold text-slate-700 border-r border-slate-300 text-[11px]">Reference Phone:</td>
								<td class="w-1/4 p-1.5 font-mono text-[11px]">
									{{ application.Response?.referencePhone || '—' }}
								</td>
							</tr>
							<tr>
								<td class="p-1.5 bg-slate-50 font-bold text-slate-700 border-r border-slate-300 text-[11px]">Sibling Enrolled in School:</td>
								<td class="p-1.5 font-medium border-r border-slate-300 text-[11px]">
									{{ application.Response?.siblingInSchool || 'NO' }}
								</td>
								<td class="p-1.5 bg-slate-50 font-bold text-slate-700 border-r border-slate-300 text-[11px]">Sibling Class / Grade:</td>
								<td class="p-1.5 font-medium text-[11px]">
									{{ application.Response?.siblingClass || '—' }}
								</td>
							</tr>
						</tbody>
					</table>
				</div>

				<!-- FORMAL DECLARATION -->
				<div class="pt-2 text-xs space-y-1 relative z-10">
					<p class="font-bold uppercase tracking-wider text-[10px] text-slate-900">DECLARATION:</p>
					<p class="text-[11px] leading-relaxed text-slate-700">
						I wish to enroll my child at <strong class="text-slate-900">{{ currentSchoolName }}</strong>.
						I have read all the rules carefully and agree to abide by the School Regulations and fee policies. All information provided above is correct to the best of my knowledge.
					</p>
				</div>

				<!-- SIGNATURES -->
				<div class="pt-8 flex justify-between items-end text-xs text-slate-700 relative z-10">
					<div class="text-center w-56">
						<div class="border-b-2 border-slate-900 pb-1 mb-1 font-medium text-slate-800 text-[11px]">
							{{ application.Response?.guardianName || application.ApplicantLastName + ' Guardian' }}
						</div>
						<div class="text-[10px] text-slate-500 uppercase tracking-wider font-bold">
							Parent / Guardian Signature
						</div>
					</div>

					<div class="text-center w-60">
						<div class="border-b-2 border-slate-900 pb-1 mb-1 font-bold text-slate-900 text-[11px]">
							Office of Admissions
						</div>
						<div class="text-[10px] text-slate-500 uppercase tracking-wider font-bold">
							Chairman / Principal Signature
						</div>
					</div>
				</div>

				<!-- FOOTER -->
				<div class="pt-4 border-t border-slate-200 text-center text-[10px] text-slate-500 space-y-0.5 relative z-10">
					<p class="font-semibold text-slate-700">{{ currentSchoolAddress }}</p>
					<p v-if="currentSchoolContact">Contact: {{ currentSchoolContact }}</p>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import type { ApplicationDTO } from '@/types/common'
import { formatCivilDate } from '@/utils/format'
import { Printer, Upload, X, User, Camera } from 'lucide-vue-next'

const props = defineProps<{
	application: ApplicationDTO | null
	isOpen: boolean
}>()

const emit = defineEmits<{
	(e: 'close'): void
	(e: 'photoUpdated', photoBase64: string): void
}>()

const authStore = useAuthStore()
const candidatePhoto = ref<string | null>(null)

const currentSchoolName = computed(() => {
	return authStore.school?.name || 'School Management System'
})

const currentSchoolLegalName = computed(() => {
	return authStore.school?.legalName || ''
})

const currentSchoolCode = computed(() => {
	return authStore.school?.code || ''
})

const currentSchoolCampus = computed(() => {
	if (authStore.school?.code) {
		return `${authStore.school.code} CAMPUS`
	}
	return 'MAIN CAMPUS'
})

const schoolInitials = computed(() => {
	const code = currentSchoolCode.value
	if (code && code !== 'OFFLINE' && code.length >= 2) {
		return code.length <= 3 ? code.toUpperCase() : code.substring(0, 2).toUpperCase()
	}
	const name = currentSchoolName.value
	const words = name.trim().split(/\s+/)
	if (words.length >= 2) {
		return (words[0][0] + words[1][0]).toUpperCase()
	}
	return name.substring(0, 2).toUpperCase()
})

const currentSchoolAddress = computed(() => {
	const addr = authStore.school?.address
	if (addr) {
		const parts = [addr.line1, addr.city, addr.country].filter(Boolean)
		if (parts.length > 0) return parts.join(', ')
	}
	return 'Academic Campus & Central Office'
})

const currentSchoolContact = computed(() => {
	const s = authStore.school
	const parts = []
	if (s?.phone) parts.push(`Mob / Tel: ${s.phone}`)
	if (s?.email) parts.push(`Email: ${s.email}`)
	return parts.join(' • ')
})

const formattedAdmissionDate = computed(() => {
	if (props.application?.SubmittedAt) {
		try {
			const d = new Date(props.application.SubmittedAt)
			return d.toLocaleDateString(undefined, { year: 'numeric', month: '2-digit', day: '2-digit' })
		} catch (e) {
			return props.application.SubmittedAt.split('T')[0]
		}
	}
	return new Date().toLocaleDateString(undefined, { year: 'numeric', month: '2-digit', day: '2-digit' })
})

const formattedDOB = computed(() => {
	if (props.application?.ApplicantDOB) {
		return formatCivilDate(props.application.ApplicantDOB)
	}
	return '—'
})

watch(() => props.application, (app) => {
	if (app) {
		loadCandidatePhoto(app.ID)
	}
}, { immediate: true })

function loadCandidatePhoto(appId: string) {
	const stored = localStorage.getItem(`candidate_photo_${appId}`)
	if (stored) {
		candidatePhoto.value = stored
	} else if (props.application?.Response?.candidatePhoto) {
		candidatePhoto.value = props.application.Response.candidatePhoto
	} else {
		candidatePhoto.value = null
	}
}

function onPhotoSelected(e: Event) {
	const file = (e.target as HTMLInputElement).files?.[0]
	if (!file || !props.application) return
	const reader = new FileReader()
	reader.onload = (event) => {
		const base64 = event.target?.result as string
		candidatePhoto.value = base64
		localStorage.setItem(`candidate_photo_${props.application!.ID}`, base64)
		emit('photoUpdated', base64)
	}
	reader.readAsDataURL(file)
}

function printForm() {
	const originalTitle = document.title
	document.title = ''
	window.print()
	window.addEventListener('afterprint', () => {
		document.title = originalTitle
	}, { once: true })
	setTimeout(() => {
		document.title = originalTitle
	}, 1500)
}

onMounted(async () => {
	if (!authStore.school?.name) {
		await authStore.refreshSchool()
	}
})
</script>

<style scoped>
@media print {
	@page {
		margin: 0 !important;
		size: auto;
	}

	.no-print {
		display: none !important;
	}

	.admission-paper {
		border: none !important;
		box-shadow: none !important;
		padding: 10mm 12mm !important;
		margin: 0 auto !important;
		max-width: 100% !important;
		width: 100% !important;
		box-sizing: border-box !important;
		page-break-after: avoid;
	}
}
</style>

