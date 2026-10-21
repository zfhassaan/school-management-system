<template>
	<div v-if="isOpen && application" class="fixed inset-0 bg-slate-900/60 backdrop-blur-xs flex items-center justify-center p-4 z-50 overflow-y-auto animate-in fade-in duration-150 print:static print:bg-transparent print:p-0 print:overflow-visible print:block">
		<div class="bg-white rounded-2xl shadow-2xl max-w-4xl w-full p-6 border border-slate-200 my-6 max-h-[95vh] overflow-y-auto print:max-h-none print:overflow-visible print:p-0 print:border-none print:shadow-none print:my-0 print:w-full print:max-w-none">
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
						<span>{{ candidatePhoto ? 'Change Photo' : 'Upload Candidate Photo' }}</span>
						<input type="file" accept="image/png,image/jpeg,image/webp" class="hidden" @change="onPhotoSelected" />
					</label>
					<button
						v-if="candidatePhoto"
						@click="removeCandidatePhoto"
						type="button"
						class="py-1.5 px-3 bg-rose-50 border border-rose-200 hover:bg-rose-100 text-rose-700 font-medium rounded-lg text-xs transition cursor-pointer"
					>
						Remove Photo
					</button>
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
			<div class="admission-paper bg-white text-slate-900 border-2 border-slate-900 p-6 md:p-8 space-y-3.5 relative">
				
				<!-- 1. FORMAL INSTITUTIONAL HEADER -->
				<div class="flex items-center justify-between gap-4 border-b-2 border-slate-900 pb-3">
					<!-- Left: School Crest / Logo -->
					<div class="w-24 h-24 rounded-xl border border-slate-300 bg-white p-1.5 flex items-center justify-center shrink-0 shadow-2xs">
						<img v-if="authStore.schoolLogo" :src="authStore.schoolLogo" alt="School Logo" class="max-w-full max-h-full object-contain" />
						<div v-else class="w-full h-full bg-slate-900 text-white rounded-lg flex items-center justify-center font-black text-2xl tracking-widest font-serif">
							{{ schoolInitials }}
						</div>
					</div>

					<!-- Center: School Name & Document Title -->
					<div class="flex-1 text-center px-2">
						<h1 class="text-2xl md:text-3xl font-black uppercase text-slate-900 tracking-tight leading-none font-serif">
							{{ currentSchoolName }}
						</h1>
						<p v-if="currentSchoolLegalName" class="text-xs font-semibold text-slate-600 uppercase tracking-wider mt-1">
							{{ currentSchoolLegalName }}
						</p>
						<p class="text-xs font-bold text-slate-700 tracking-wide mt-1">
							{{ currentSchoolAddress }}
						</p>
						<p v-if="currentSchoolContact" class="text-[11px] text-slate-500 font-medium mt-0.5">
							{{ currentSchoolContact }}
						</p>
						<div class="inline-block mt-2 px-6 py-1 bg-slate-900 text-white rounded shadow-2xs">
							<h2 class="text-xs md:text-sm font-black tracking-widest uppercase">
								OFFICIAL ADMISSION APPLICATION FORM
							</h2>
						</div>
					</div>

					<!-- Right: Candidate Passport Photo Frame -->
					<div class="w-28 h-36 border-2 border-dashed border-slate-400 bg-slate-50 rounded flex flex-col items-center justify-center shrink-0 relative overflow-hidden group shadow-2xs">
						<img v-if="candidatePhoto" :src="candidatePhoto" alt="Candidate" class="w-full h-full object-cover" />
						<div v-else class="text-center p-2 flex flex-col items-center justify-center h-full">
							<Camera class="w-7 h-7 mx-auto mb-1 text-slate-400" />
							<span class="text-[9px] font-bold uppercase tracking-wider text-slate-700 leading-tight">
								Affix Student Photo
							</span>
							<span class="text-[8px] text-slate-500 mt-1">Passport Size</span>
							<span class="text-[7px] text-slate-400 italic mt-0.5">(Recent &amp; Attested)</span>
						</div>
						<!-- Non-printable hover to upload -->
						<label class="absolute inset-0 bg-slate-900/60 opacity-0 group-hover:opacity-100 flex flex-col items-center justify-center text-white text-[9px] font-bold transition cursor-pointer no-print p-1 text-center">
							<Camera class="w-4 h-4 mb-1" />
							<span>{{ candidatePhoto ? 'Change' : 'Upload' }}</span>
							<input type="file" accept="image/png,image/jpeg,image/webp" class="hidden" @change="onPhotoSelected" />
						</label>
					</div>
				</div>

				<!-- 2. METADATA STAT STRIP -->
				<div class="grid grid-cols-4 border-2 border-slate-900 bg-slate-100 divide-x-2 divide-slate-900 text-xs text-center font-medium">
					<div class="p-1.5">
						<span class="text-[9px] uppercase tracking-wider text-slate-600 block font-bold">Application No</span>
						<span class="font-mono font-black text-slate-900 text-sm">{{ application.ApplicationNo }}</span>
					</div>
					<div class="p-1.5">
						<span class="text-[9px] uppercase tracking-wider text-slate-600 block font-bold">Admission Date</span>
						<span class="font-bold text-slate-900 text-xs">{{ formattedAdmissionDate }}</span>
					</div>
					<div class="p-1.5">
						<span class="text-[9px] uppercase tracking-wider text-slate-600 block font-bold">Enrolling Class</span>
						<span class="font-black text-indigo-950 text-xs">{{ application.ClassName || 'Unassigned' }}</span>
					</div>
					<div class="p-1.5">
						<span class="text-[9px] uppercase tracking-wider text-slate-600 block font-bold">Admission Status</span>
						<span class="font-black uppercase text-xs" :class="application.Status === 'enrolled' ? 'text-emerald-700' : 'text-slate-900'">
							{{ application.Status.replace('_', ' ') }}
						</span>
					</div>
				</div>

				<!-- 3. SECTION 1: APPLICANT PERSONAL DETAILS -->
				<div class="space-y-1">
					<div class="bg-slate-900 text-white px-3 py-1 text-[11px] font-black uppercase tracking-wider flex items-center justify-between">
						<span>SECTION 1: APPLICANT PERSONAL INFORMATION</span>
						<span class="text-[9px] font-normal tracking-normal text-slate-300 font-mono">App No: {{ application.ApplicationNo }}</span>
					</div>
					<table class="w-full border-collapse border-2 border-slate-900 text-xs">
						<tbody>
							<tr class="border-b border-slate-300">
								<td class="w-[20%] p-1.5 bg-slate-100 font-bold text-slate-800 border-r border-slate-300 text-[11px]">Full Student Name:</td>
								<td class="w-[30%] p-1.5 font-black text-slate-900 text-xs border-r border-slate-300">
									{{ application.ApplicantFirstName }} {{ application.ApplicantLastName }}
								</td>
								<td class="w-[20%] p-1.5 bg-slate-100 font-bold text-slate-800 border-r border-slate-300 text-[11px]">B-Form / CNIC No:</td>
								<td class="w-[30%] p-1.5 font-mono font-bold text-slate-900 text-[11px]">
									{{ application.Response?.bFormNo || application.Response?.cnic || '—' }}
								</td>
							</tr>
							<tr class="border-b border-slate-300">
								<td class="p-1.5 bg-slate-100 font-bold text-slate-800 border-r border-slate-300 text-[11px]">Date of Birth:</td>
								<td class="p-1.5 font-medium border-r border-slate-300 text-[11px]">
									{{ formattedDOB }}
								</td>
								<td class="p-1.5 bg-slate-100 font-bold text-slate-800 border-r border-slate-300 text-[11px]">Gender / Blood Group:</td>
								<td class="p-1.5 font-bold uppercase text-[11px]">
									{{ application.ApplicantGender || 'Not Specified' }} {{ application.Response?.bloodGroup ? `(${application.Response.bloodGroup})` : '' }}
								</td>
							</tr>
							<tr class="border-b border-slate-300">
								<td class="p-1.5 bg-slate-100 font-bold text-slate-800 border-r border-slate-300 text-[11px]">Residential Address:</td>
								<td colspan="3" class="p-1.5 font-medium text-[11px] text-slate-800">
									{{ application.Response?.address || application.Response?.homeAddress || 'As per guardian records' }}
								</td>
							</tr>
							<tr>
								<td class="p-1.5 bg-slate-100 font-bold text-slate-800 border-r border-slate-300 text-[11px]">Previous School Passed:</td>
								<td class="p-1.5 font-medium border-r border-slate-300 text-[11px]">
									{{ application.Response?.previousSchool || '— None / Fresh Admission —' }}
								</td>
								<td class="p-1.5 bg-slate-100 font-bold text-slate-800 border-r border-slate-300 text-[11px]">Medical / Special Notes:</td>
								<td class="p-1.5 font-medium text-[11px]">
									{{ application.Response?.medicalHistory || 'None Reported' }}
								</td>
							</tr>
						</tbody>
					</table>
				</div>

				<!-- 4. SECTION 2: PARENTS & GUARDIAN DETAILS -->
				<div class="space-y-1">
					<div class="bg-slate-900 text-white px-3 py-1 text-[11px] font-black uppercase tracking-wider flex items-center justify-between">
						<span>SECTION 2: PARENTS &amp; GUARDIAN DETAILS</span>
					</div>
					<table class="w-full border-collapse border-2 border-slate-900 text-xs">
						<tbody>
							<tr class="border-b border-slate-300">
								<td class="w-[20%] p-1.5 bg-slate-100 font-bold text-slate-800 border-r border-slate-300 text-[11px]">Father / Guardian Name:</td>
								<td class="w-[30%] p-1.5 font-bold text-slate-900 border-r border-slate-300 text-[11px]">
									{{ application.Response?.guardianName || application.ApplicantLastName + ' Guardian' }}
								</td>
								<td class="w-[20%] p-1.5 bg-slate-100 font-bold text-slate-800 border-r border-slate-300 text-[11px]">Father / Guardian CNIC:</td>
								<td class="w-[30%] p-1.5 font-mono font-medium text-slate-900 text-[11px]">
									{{ application.Response?.guardianCNIC || '—' }}
								</td>
							</tr>
							<tr class="border-b border-slate-300">
								<td class="p-1.5 bg-slate-100 font-bold text-slate-800 border-r border-slate-300 text-[11px]">Primary Mobile (WhatsApp):</td>
								<td class="p-1.5 font-mono font-bold text-slate-900 border-r border-slate-300 text-[11px]">
									{{ application.GuardianPhoneNormalized || '—' }}
								</td>
								<td class="p-1.5 bg-slate-100 font-bold text-slate-800 border-r border-slate-300 text-[11px]">Occupation / Business:</td>
								<td class="p-1.5 font-medium text-slate-800 text-[11px]">
									{{ application.Response?.guardianOccupation || '—' }}
								</td>
							</tr>
							<tr>
								<td class="p-1.5 bg-slate-100 font-bold text-slate-800 border-r border-slate-300 text-[11px]">Mother's Name:</td>
								<td class="p-1.5 font-medium border-r border-slate-300 text-[11px]">
									{{ application.Response?.motherName || '—' }}
								</td>
								<td class="p-1.5 bg-slate-100 font-bold text-slate-800 border-r border-slate-300 text-[11px]">Emergency Contact Phone:</td>
								<td class="p-1.5 font-mono font-medium text-slate-900 text-[11px]">
									{{ application.Response?.emergencyPhone || application.GuardianPhoneNormalized || '—' }}
								</td>
							</tr>
						</tbody>
					</table>
				</div>

				<!-- 5. SECTION 3: SIBLINGS & EMERGENCY REFERENCE -->
				<div class="space-y-1">
					<div class="bg-slate-900 text-white px-3 py-1 text-[11px] font-black uppercase tracking-wider flex items-center justify-between">
						<span>SECTION 3: SIBLINGS &amp; EMERGENCY REFERENCE</span>
					</div>
					<table class="w-full border-collapse border-2 border-slate-900 text-xs">
						<tbody>
							<tr class="border-b border-slate-300">
								<td class="w-[20%] p-1.5 bg-slate-100 font-bold text-slate-800 border-r border-slate-300 text-[11px]">Sibling Enrolled in School:</td>
								<td class="w-[30%] p-1.5 font-medium border-r border-slate-300 text-[11px]">
									{{ application.Response?.siblingInSchool || 'No' }}
								</td>
								<td class="w-[20%] p-1.5 bg-slate-100 font-bold text-slate-800 border-r border-slate-300 text-[11px]">Sibling Class / Roll No:</td>
								<td class="w-[30%] p-1.5 font-medium text-[11px]">
									{{ application.Response?.siblingClass || '—' }}
								</td>
							</tr>
							<tr>
								<td class="p-1.5 bg-slate-100 font-bold text-slate-800 border-r border-slate-300 text-[11px]">Emergency Reference Person:</td>
								<td class="p-1.5 font-medium border-r border-slate-300 text-[11px]">
									{{ application.Response?.referenceContact || 'Family Guardian' }}
								</td>
								<td class="p-1.5 bg-slate-100 font-bold text-slate-800 border-r border-slate-300 text-[11px]">Contact Person Phone:</td>
								<td class="p-1.5 font-mono font-medium text-[11px]">
									{{ application.Response?.referencePhone || application.GuardianPhoneNormalized || '—' }}
								</td>
							</tr>
						</tbody>
					</table>
				</div>

				<!-- 6. FORMAL UNDERTAKING & DECLARATION -->
				<div class="border border-slate-400 p-2.5 bg-slate-50 text-[10.5px] space-y-1 text-slate-700 leading-relaxed">
					<p class="font-black uppercase tracking-wider text-[10px] text-slate-900">
						OFFICIAL UNDERTAKING &amp; DECLARATION BY PARENT / GUARDIAN:
					</p>
					<p>
						I/We hereby declare that all particulars entered in this admission application are authentic and true to the best of my knowledge. I understand that admission of my child/ward is subject to credential verification and approval by the school authority. I agree to abide by all the rules, regulations, fee schedule, and discipline policies of <strong>{{ currentSchoolName }}</strong> in force now or hereafter.
					</p>
				</div>

				<!-- 7. FORMAL SIGNATURES -->
				<div class="grid grid-cols-3 gap-6 pt-5 pb-2 text-center text-xs">
					<div>
						<div class="border-b-2 border-slate-900 pb-1 mb-1 font-bold text-slate-800 text-[11px] min-h-[20px]">
							{{ application.Response?.guardianName || application.ApplicantLastName + ' Guardian' }}
						</div>
						<p class="text-[10px] font-bold uppercase tracking-wider text-slate-700">Parent / Guardian Signature</p>
						<p class="text-[9px] text-slate-400">Date: {{ formattedAdmissionDate }}</p>
					</div>
					<div>
						<div class="border-b-2 border-slate-900 pb-1 mb-1 font-bold text-slate-800 text-[11px] min-h-[20px]">
							Admissions Committee
						</div>
						<p class="text-[10px] font-bold uppercase tracking-wider text-slate-700">Verified &amp; Checked By</p>
						<p class="text-[9px] text-slate-400">Credential Verification</p>
					</div>
					<div>
						<div class="border-b-2 border-slate-900 pb-1 mb-1 font-bold text-slate-900 text-[11px] min-h-[20px]">
							Principal / Director
						</div>
						<p class="text-[10px] font-bold uppercase tracking-wider text-slate-700">Principal / Head Seal</p>
						<p class="text-[9px] text-slate-400">Official Sanction</p>
					</div>
				</div>

				<!-- 8. OFFICIAL TEAR-OFF OFFICE RECORD SLIP -->
				<div class="border-t-2 border-dashed border-slate-400 pt-2.5 space-y-2">
					<div class="flex items-center justify-between text-[10px] font-bold uppercase tracking-wider text-slate-600">
						<span>ADMISSION ACKNOWLEDGEMENT SLIP (OFFICE / PARENT COPY)</span>
						<span class="font-mono">{{ application.ApplicationNo }}</span>
					</div>
					<div class="border border-slate-300 p-2 rounded bg-slate-50 grid grid-cols-4 gap-2 text-[10.5px]">
						<div>
							<span class="text-[9px] text-slate-500 block uppercase font-bold">Candidate Name</span>
							<span class="font-bold text-slate-900">{{ application.ApplicantFirstName }} {{ application.ApplicantLastName }}</span>
						</div>
						<div>
							<span class="text-[9px] text-slate-500 block uppercase font-bold">Enrolled Class</span>
							<span class="font-bold text-slate-900">{{ application.ClassName || 'Grade 1' }}</span>
						</div>
						<div>
							<span class="text-[9px] text-slate-500 block uppercase font-bold">Roll / Form No</span>
							<span class="font-bold text-slate-900 font-mono">{{ application.ApplicationNo }}</span>
						</div>
						<div>
							<span class="text-[9px] text-slate-500 block uppercase font-bold">Admission Status</span>
							<span class="font-bold uppercase text-emerald-700">{{ application.Status.replace('_', ' ') }}</span>
						</div>
					</div>
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
import { Printer, Upload, X, Camera } from 'lucide-vue-next'

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
	const legal = authStore.school?.legalName || ''
	if (legal && legal.trim().toLowerCase() === currentSchoolName.value.trim().toLowerCase()) {
		return ''
	}
	return legal
})

const currentSchoolCode = computed(() => {
	return authStore.school?.code || ''
})

const currentSchoolAddress = computed(() => {
	const addr = authStore.school?.address
	if (addr) {
		const parts = [addr.line1, addr.city, addr.country].filter(Boolean)
		if (parts.length > 0) return parts.join(', ')
	}
	if (authStore.school?.code) {
		return `${authStore.school.code} Campus`
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

function removeCandidatePhoto() {
	candidatePhoto.value = null
	if (props.application) {
		localStorage.removeItem(`candidate_photo_${props.application.ID}`)
		if (props.application.Response) {
			delete props.application.Response.candidatePhoto
		}
		emit('photoUpdated', '')
	}
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
	* {
		-webkit-print-color-adjust: exact !important;
		print-color-adjust: exact !important;
	}

	@page {
		margin: 0 !important;
		size: A4 portrait;
	}

	.no-print {
		display: none !important;
	}

	.admission-paper {
		border: 2px solid #0f172a !important;
		box-shadow: none !important;
		padding: 8mm 10mm !important;
		margin: 0 auto !important;
		max-width: 100% !important;
		width: 100% !important;
		box-sizing: border-box !important;
		page-break-after: avoid;
		page-break-inside: avoid;
	}
}
</style>
