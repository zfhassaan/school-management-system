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
			<div id="admission-printable-area" class="admission-paper bg-white text-slate-900 border border-slate-200 shadow-xs print:border-none print:shadow-none p-6 md:p-8 space-y-3 relative">
				
				<!-- 1. FORMAL INSTITUTIONAL HEADER -->
				<div class="flex items-center justify-between gap-4 border-b-2 border-slate-900 pb-2.5">
					<!-- Left: School Crest / Logo -->
					<div class="w-20 h-20 md:w-24 md:h-24 rounded-xl border border-slate-300 bg-white p-1 flex items-center justify-center shrink-0 shadow-2xs">
						<img v-if="authStore.schoolLogo" :src="authStore.schoolLogo" alt="School Logo" class="max-w-full max-h-full object-contain" />
						<div v-else class="w-full h-full bg-slate-900 text-white rounded-lg flex items-center justify-center font-black text-2xl tracking-widest font-serif">
							{{ schoolInitials }}
						</div>
					</div>

					<!-- Center: School Name, Campus & Form Title -->
					<div class="flex-1 text-center px-2">
						<h1 class="text-2xl md:text-3xl font-black uppercase text-slate-900 tracking-tight leading-tight font-serif">
							{{ currentSchoolName }}
						</h1>
						<p v-if="currentSchoolLegalName" class="text-xs font-semibold text-slate-700 tracking-wide mt-0.5">
							{{ currentSchoolLegalName }}
						</p>
						<p class="text-xs font-black uppercase text-slate-800 tracking-wider mt-0.5">
							{{ campusName }}
						</p>
						<div class="inline-block mt-1.5 px-6 py-0.5 bg-slate-900 text-white rounded shadow-2xs">
							<h2 class="text-xs md:text-sm font-black tracking-widest uppercase">
								ADMISSION FORM
							</h2>
						</div>
					</div>

					<!-- Right: Candidate Passport Photo Frame -->
					<div class="w-24 h-32 md:w-28 md:h-36 border-2 border-slate-900 bg-slate-50 rounded flex flex-col items-center justify-center shrink-0 relative overflow-hidden group shadow-2xs">
						<img v-if="candidatePhoto" :src="candidatePhoto" alt="Candidate" class="w-full h-full object-cover" />
						<div v-else class="text-center p-2 flex flex-col items-center justify-center h-full">
							<Camera class="w-6 h-6 mx-auto mb-1 text-slate-400" />
							<span class="text-[9px] font-bold uppercase tracking-wider text-slate-700 leading-tight">
								Affix Photo
							</span>
							<span class="text-[8px] text-slate-500 mt-0.5">Passport Size</span>
						</div>
						<!-- Non-printable hover to upload -->
						<label class="absolute inset-0 bg-slate-900/60 opacity-0 group-hover:opacity-100 flex flex-col items-center justify-center text-white text-[9px] font-bold transition cursor-pointer no-print p-1 text-center">
							<Camera class="w-4 h-4 mb-1" />
							<span>{{ candidatePhoto ? 'Change' : 'Upload' }}</span>
							<input type="file" accept="image/png,image/jpeg,image/webp" class="hidden" @change="onPhotoSelected" />
						</label>
					</div>
				</div>

				<!-- SECTION 1: OFFICE USE & FINANCIALS -->
				<div class="space-y-0.5">
					<div class="bg-[#002060] text-white px-2.5 py-0.5 text-[11px] font-black uppercase tracking-wider flex items-center justify-between">
						<span>SECTION 1: OFFICE USE &amp; FINANCIALS</span>
						<span class="text-[9px] font-normal tracking-normal text-slate-200 font-mono">Form: {{ application.ApplicationNo }}</span>
					</div>
					<table class="w-full border-collapse border border-black text-[11px]">
						<tbody>
							<tr class="border-b border-black">
								<td class="w-[20%] p-1.5 font-bold border-r border-black bg-slate-50">Admission Date:</td>
								<td class="w-[30%] p-1.5 border-r border-black font-semibold">{{ formattedAdmissionDate }}</td>
								<td class="w-[20%] p-1.5 font-bold border-r border-black bg-slate-50">Admission / Roll #:</td>
								<td class="w-[30%] p-1.5 font-bold font-mono text-slate-900">{{ admissionOrRollNo }}</td>
							</tr>
							<tr class="border-b border-black">
								<td class="p-1.5 font-bold border-r border-black bg-slate-50">Admission For Class:</td>
								<td class="p-1.5 border-r border-black font-semibold">{{ application.ClassName || 'Unassigned' }}</td>
								<td class="p-1.5 font-bold border-r border-black bg-slate-50">Campus:</td>
								<td class="p-1.5 font-semibold">{{ campusName }}</td>
							</tr>
							<tr>
								<td colspan="4" class="p-0">
									<div class="grid grid-cols-5 divide-x divide-black text-center text-[10.5px]">
										<div class="p-1.5">
											<span class="font-bold text-[10px] text-slate-600 block">Adm Fee:</span>
											<span class="font-bold text-slate-900">{{ admFee }}</span>
										</div>
										<div class="p-1.5">
											<span class="font-bold text-[10px] text-slate-600 block">Monthly:</span>
											<span class="font-bold text-slate-900">{{ monthlyFee }}</span>
										</div>
										<div class="p-1.5">
											<span class="font-bold text-[10px] text-slate-600 block">Annual:</span>
											<span class="font-bold text-slate-900">{{ annualFee }}</span>
										</div>
										<div class="p-1.5">
											<span class="font-bold text-[10px] text-slate-600 block">Paper:</span>
											<span class="font-bold text-slate-900">{{ paperFund }}</span>
										</div>
										<div class="p-1.5">
											<span class="font-bold text-[10px] text-slate-600 block">ID Card:</span>
											<span class="font-bold text-slate-900">{{ idCardFee }}</span>
										</div>
									</div>
								</td>
							</tr>
						</tbody>
					</table>
				</div>

				<!-- SECTION 2: STUDENT PERSONAL INFO -->
				<div class="space-y-0.5">
					<div class="bg-[#002060] text-white px-2.5 py-0.5 text-[11px] font-black uppercase tracking-wider">
						SECTION 2: STUDENT PERSONAL INFO
					</div>
					<table class="w-full border-collapse border border-black text-[11px]">
						<tbody>
							<tr class="border-b border-black">
								<td class="w-[20%] p-1.5 font-bold border-r border-black bg-slate-50">Student Name:</td>
								<td class="w-[30%] p-1.5 border-r border-black font-black text-slate-900 uppercase">
									{{ application.ApplicantFirstName }} {{ application.ApplicantLastName }}
								</td>
								<td class="w-[20%] p-1.5 font-bold border-r border-black bg-slate-50">B-Form / CNIC:</td>
								<td class="w-[30%] p-1.5 font-bold font-mono text-slate-900">
									{{ bFormCnic }}
								</td>
							</tr>
							<tr class="border-b border-black">
								<td class="p-1.5 font-bold border-r border-black bg-slate-50">Date of Birth:</td>
								<td class="p-1.5 border-r border-black font-semibold">{{ formattedDOB }}</td>
								<td class="p-1.5 font-bold border-r border-black bg-slate-50">Medical History:</td>
								<td class="p-1.5 font-semibold">{{ medicalHistory }}</td>
							</tr>
							<tr class="border-b border-black">
								<td class="p-1.5 font-bold border-r border-black bg-slate-50">Home Address:</td>
								<td colspan="3" class="p-1.5 font-medium text-slate-800">{{ homeAddress }}</td>
							</tr>
							<tr>
								<td class="p-1.5 font-bold border-r border-black bg-slate-50">Previous School:</td>
								<td class="p-1.5 border-r border-black font-semibold">{{ previousSchool }}</td>
								<td class="p-1.5 font-bold border-r border-black bg-slate-50">Last Grade/Result:</td>
								<td class="p-1.5 font-semibold">{{ lastGradeResult }}</td>
							</tr>
						</tbody>
					</table>
				</div>

				<!-- SECTION 3: PARENTS DETAILS -->
				<div class="space-y-0.5">
					<div class="bg-[#002060] text-white px-2.5 py-0.5 text-[11px] font-black uppercase tracking-wider">
						SECTION 3: PARENTS DETAILS
					</div>
					<table class="w-full border-collapse border border-black text-[11px]">
						<tbody>
							<tr class="border-b border-black">
								<td class="w-[20%] p-1.5 font-bold border-r border-black bg-slate-50">Father's Name:</td>
								<td class="w-[30%] p-1.5 border-r border-black font-bold text-slate-900">{{ fatherName }}</td>
								<td class="w-[20%] p-1.5 font-bold border-r border-black bg-slate-50">Mobile (WhatsApp):</td>
								<td class="w-[30%] p-1.5 font-bold font-mono text-slate-900">{{ fatherWhatsApp }}</td>
							</tr>
							<tr class="border-b border-black">
								<td class="p-1.5 font-bold border-r border-black bg-slate-50">Qualification:</td>
								<td class="p-1.5 border-r border-black font-medium">{{ fatherQualification }}</td>
								<td class="p-1.5 font-bold border-r border-black bg-slate-50">Occupation:</td>
								<td class="p-1.5 font-medium">{{ fatherOccupation }}</td>
							</tr>
							<tr class="border-b border-black">
								<td class="p-1.5 font-bold border-r border-black bg-slate-50">Office Address:</td>
								<td class="p-1.5 border-r border-black font-medium">{{ fatherOfficeAddress }}</td>
								<td class="p-1.5 font-bold border-r border-black bg-slate-50">Office Phone:</td>
								<td class="p-1.5 font-mono font-medium">{{ fatherOfficePhone }}</td>
							</tr>
							<tr class="border-b border-black">
								<td class="p-1.5 font-bold border-r border-black bg-slate-50">Mother's Name:</td>
								<td class="p-1.5 border-r border-black font-medium">{{ motherName }}</td>
								<td class="p-1.5 font-bold border-r border-black bg-slate-50">Mobile Number:</td>
								<td class="p-1.5 font-mono font-medium">{{ motherPhone }}</td>
							</tr>
							<tr>
								<td class="p-1.5 font-bold border-r border-black bg-slate-50">Qualification:</td>
								<td class="p-1.5 border-r border-black font-medium">{{ motherQualification }}</td>
								<td class="p-1.5 font-bold border-r border-black bg-slate-50">Occupation:</td>
								<td class="p-1.5 font-medium">{{ motherOccupation }}</td>
							</tr>
						</tbody>
					</table>
				</div>

				<!-- SECTION 4: SIBLINGS & REFERENCE -->
				<div class="space-y-0.5">
					<div class="bg-[#002060] text-white px-2.5 py-0.5 text-[11px] font-black uppercase tracking-wider">
						SECTION 4: SIBLINGS &amp; REFERENCE
					</div>
					<table class="w-full border-collapse border border-black text-[11px]">
						<tbody>
							<tr class="border-b border-black">
								<td class="w-[20%] p-1.5 font-bold border-r border-black bg-slate-50">Reference Name:</td>
								<td class="w-[30%] p-1.5 border-r border-black font-medium">{{ referenceName }}</td>
								<td class="w-[20%] p-1.5 font-bold border-r border-black bg-slate-50">Reference Phone:</td>
								<td class="w-[30%] p-1.5 font-mono font-medium">{{ referencePhone }}</td>
							</tr>
							<tr>
								<td class="p-1.5 font-bold border-r border-black bg-slate-50">Sibling in School:</td>
								<td class="p-1.5 border-r border-black font-medium">{{ siblingInSchool }}</td>
								<td colspan="2" class="p-0">
									<div class="grid grid-cols-2 divide-x divide-black h-full">
										<div class="p-1.5 flex items-center gap-2">
											<span class="font-bold text-slate-700">Class:</span>
											<span class="font-medium">{{ siblingClass }}</span>
										</div>
										<div class="p-1.5 flex items-center gap-2">
											<span class="font-bold text-slate-700">Since:</span>
											<span class="font-medium">{{ siblingSince }}</span>
										</div>
									</div>
								</td>
							</tr>
						</tbody>
					</table>
				</div>

				<!-- DECLARATION -->
				<div class="pt-2 text-[10.5px] text-slate-800 leading-relaxed">
					<p class="font-black uppercase tracking-wider text-slate-900 mb-0.5">
						DECLARATION:
					</p>
					<p>
						I Wish to enroll my child at <strong>{{ currentSchoolName }}</strong>. I have read all the rules carefully and agree to abide by the School Regulations. All information provided above is correct to the best of my knowledge.
					</p>
				</div>

				<!-- FORMAL SIGNATURES -->
				<div class="flex justify-between items-end pt-12 pb-2 px-8 text-xs font-bold text-slate-900">
					<div class="text-center w-56">
						<div class="border-t-2 border-slate-900 pt-1.5">
							Parent/Guardian Signature
						</div>
					</div>
					<div class="text-center w-56">
						<div class="border-t-2 border-slate-900 pt-1.5">
							Chairman / Principal
						</div>
					</div>
				</div>

				<!-- OFFICIAL FOOTER -->
				<div class="text-center text-[10.5px] text-slate-800 pt-3 border-t border-slate-300 space-y-0.5">
					<p class="font-bold">{{ currentSchoolAddress }}</p>
					<p v-if="currentSchoolContact" class="font-medium text-slate-600">{{ currentSchoolContact }}</p>
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

const campusName = computed(() => {
	const code = authStore.school?.code
	if (code && code !== 'OFFLINE') {
		return `${code} Campus`
	}
	return 'Main Campus'
})

const resp = computed<Record<string, any>>(() => {
	return (props.application?.Response as Record<string, any>) || {}
})

const admissionOrRollNo = computed(() => {
	return props.application?.AdmissionNo || props.application?.RollNo || props.application?.ApplicationNo || '—'
})

// Financials
const admFee = computed(() => resp.value.adm_fee ? `Rs ${resp.value.adm_fee}` : (resp.value.admission_fee ? `Rs ${resp.value.admission_fee}` : '—'))
const monthlyFee = computed(() => resp.value.monthly_fee ? `Rs ${resp.value.monthly_fee}` : (resp.value.monthly_tuition_fee ? `Rs ${resp.value.monthly_tuition_fee}` : '—'))
const annualFee = computed(() => resp.value.annual_fee ? `Rs ${resp.value.annual_fee}` : (resp.value.annual_charges ? `Rs ${resp.value.annual_charges}` : '—'))
const paperFund = computed(() => resp.value.paper_fund ? `Rs ${resp.value.paper_fund}` : (resp.value.paper_exam_fund ? `Rs ${resp.value.paper_exam_fund}` : '—'))
const idCardFee = computed(() => resp.value.id_card_fee ? `Rs ${resp.value.id_card_fee}` : '—')

// Personal Info
const bFormCnic = computed(() => resp.value.b_form_cnic || resp.value.b_form || resp.value.bFormNo || resp.value.cnic || '—')
const medicalHistory = computed(() => resp.value.medical_history || resp.value.medicalHistory || resp.value.medical_special_needs || 'NO')
const homeAddress = computed(() => resp.value.home_address || resp.value.homeAddress || resp.value.residential_address || resp.value.address || '—')
const previousSchool = computed(() => resp.value.previous_school || resp.value.previousSchool || resp.value.prev_school || resp.value.notes || '—')
const lastGradeResult = computed(() => resp.value.last_grade_result || resp.value.lastGradeResult || resp.value.last_grade || '—')

// Parents Details
const fatherName = computed(() => resp.value.father_name || resp.value.fatherName || props.application?.GuardianName || resp.value.guardianName || (props.application?.ApplicantLastName ? `${props.application.ApplicantLastName} Guardian` : '—'))
const fatherWhatsApp = computed(() => resp.value.father_whatsapp || resp.value.fatherWhatsApp || props.application?.GuardianPhoneNormalized || resp.value.guardianPhone || '—')
const fatherQualification = computed(() => resp.value.father_qualification || resp.value.fatherQualification || resp.value.guardianQualification || '—')
const fatherOccupation = computed(() => resp.value.father_occupation || resp.value.fatherOccupation || resp.value.guardianOccupation || '—')
const fatherOfficeAddress = computed(() => resp.value.father_office_address || resp.value.fatherOfficeAddress || '—')
const fatherOfficePhone = computed(() => resp.value.father_office_phone || resp.value.fatherOfficePhone || '—')

const motherName = computed(() => resp.value.mother_name || resp.value.motherName || '—')
const motherPhone = computed(() => resp.value.mother_phone || resp.value.motherPhone || resp.value.emergencyPhone || '—')
const motherQualification = computed(() => resp.value.mother_qualification || resp.value.motherQualification || '—')
const motherOccupation = computed(() => resp.value.mother_occupation || resp.value.motherOccupation || '—')

// Siblings & Reference
const referenceName = computed(() => resp.value.reference_name || resp.value.referenceName || resp.value.referenceContact || 'NO')
const referencePhone = computed(() => resp.value.reference_phone || resp.value.referencePhone || 'NO')
const siblingInSchool = computed(() => resp.value.sibling_in_school || resp.value.siblingInSchool || 'NO')
const siblingClass = computed(() => resp.value.sibling_class || resp.value.siblingClass || '—')
const siblingSince = computed(() => resp.value.sibling_since || resp.value.siblingSince || '—')

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

<style>
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

	/* Hide everything on the page */
	body * {
		visibility: hidden !important;
	}

	/* Exclusively reveal and position the admission paper container */
	#admission-printable-area,
	#admission-printable-area * {
		visibility: visible !important;
	}

	#admission-printable-area {
		position: fixed !important;
		left: 0 !important;
		top: 0 !important;
		right: 0 !important;
		width: 100% !important;
		max-width: 100% !important;
		margin: 0 !important;
		padding: 6mm 8mm !important;
		background: #ffffff !important;
		box-sizing: border-box !important;
		border: none !important;
		box-shadow: none !important;
		z-index: 999999 !important;
		page-break-after: avoid !important;
		page-break-inside: avoid !important;
	}
}
</style>
