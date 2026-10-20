import { createRouter, createWebHashHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

import LoginPage from '@/pages/auth/LoginPage.vue'
import SetupPage from '@/pages/auth/SetupPage.vue'
import AppShell from '@/layouts/AppShell.vue'
import DashboardPage from '@/pages/dashboard/DashboardPage.vue'
import SchoolSettingsPage from '@/pages/settings/SchoolSettingsPage.vue'
import UsersPage from '@/pages/users/UsersPage.vue'
import RolesPage from '@/pages/roles/RolesPage.vue'
import AuditPage from '@/pages/audit/AuditPage.vue'
import StudentsListPage from '@/pages/students/StudentsListPage.vue'
import StudentProfilePage from '@/pages/students/StudentProfilePage.vue'
import AcademicsPage from '@/pages/academics/AcademicsPage.vue'
import AdmissionsPage from '@/pages/admissions/AdmissionsPage.vue'
import AdmissionFormBuilderPage from '@/pages/admissions/AdmissionFormBuilderPage.vue'
import PromotionsPage from '@/pages/promotions/PromotionsPage.vue'
import FeesStructurePage from '@/pages/fees/FeesStructurePage.vue'
import DiscountPoliciesPage from '@/pages/fees/DiscountPoliciesPage.vue'
import InvoicesListPage from '@/pages/invoices/InvoicesListPage.vue'
import InvoiceDetailPage from '@/pages/invoices/InvoiceDetailPage.vue'
import PaymentsListPage from '@/pages/payments/PaymentsListPage.vue'
import ReceiptDetailPage from '@/pages/payments/ReceiptDetailPage.vue'
import FinancialAdjustmentsPage from '@/pages/payments/FinancialAdjustmentsPage.vue'

const router = createRouter({
	history: createWebHashHistory(),
	routes: [
		{
			path: '/login',
			name: 'login',
			component: LoginPage,
			meta: { public: true, title: 'Sign In' }
		},
		{
			path: '/setup',
			name: 'setup',
			component: SetupPage,
			meta: { public: true, title: 'First-Run Setup' }
		},
		{
			path: '/',
			component: AppShell,
			children: [
				{
					path: '',
					name: 'dashboard',
					component: DashboardPage,
					meta: { title: 'Dashboard' }
				},
				{
					path: 'students',
					name: 'students',
					component: StudentsListPage,
					meta: { title: 'Students Directory', permission: 'students.view' }
				},
				{
					path: 'students/:id',
					name: 'student-profile',
					component: StudentProfilePage,
					meta: { title: 'Student Profile', permission: 'students.view' }
				},
				{
					path: 'academics',
					name: 'academics',
					component: AcademicsPage,
					meta: { title: 'Academic Structure', permission: 'academics.view' }
				},
				{
					path: 'admissions',
					name: 'admissions',
					component: AdmissionsPage,
					meta: { title: 'Admissions', permission: 'admissions.view' }
				},
				{
					path: 'admissions/builder',
					name: 'admission-builder',
					component: AdmissionFormBuilderPage,
					meta: { title: 'Form Builder', permission: 'admissions.forms.manage' }
				},
				{
					path: 'promotions',
					name: 'promotions',
					component: PromotionsPage,
					meta: { title: 'Class Promotions', permission: 'academics.promote' }
				},
				{
					path: 'fees/structure',
					name: 'fees-structure',
					component: FeesStructurePage,
					meta: { title: 'Fee Structures', permission: 'fees.view' }
				},
				{
					path: 'fees/discounts',
					name: 'fees-discounts',
					component: DiscountPoliciesPage,
					meta: { title: 'Discounts & Concessions', permission: 'fees.view' }
				},
				{
					path: 'invoices',
					name: 'invoices',
					component: InvoicesListPage,
					meta: { title: 'Student Invoices', permission: 'invoices.view' }
				},
				{
					path: 'invoices/:id',
					name: 'invoice-detail',
					component: InvoiceDetailPage,
					meta: { title: 'Invoice Details', permission: 'invoices.view' }
				},
				{
					path: 'payments',
					name: 'payments',
					component: PaymentsListPage,
					meta: { title: 'Payments & Receipts', permission: 'payments.view' }
				},
				{
					path: 'payments/:id',
					name: 'receipt-detail',
					component: ReceiptDetailPage,
					meta: { title: 'Payment Receipt', permission: 'payments.view' }
				},
				{
					path: 'payments/adjustments',
					name: 'financial-adjustments',
					component: FinancialAdjustmentsPage,
					meta: { title: 'Financial Adjustments', permission: 'adjustments.request' }
				},
				{
					path: 'settings/school',
					name: 'settings-school',
					component: SchoolSettingsPage,
					meta: { title: 'School Settings', permission: 'settings.manage' }
				},
				{
					path: 'settings/users',
					name: 'settings-users',
					component: UsersPage,
					meta: { title: 'User Accounts', permission: 'users.view' }
				},
				{
					path: 'settings/roles',
					name: 'settings-roles',
					component: RolesPage,
					meta: { title: 'Roles & Permissions', permission: 'roles.manage' }
				},
				{
					path: 'audit',
					name: 'audit',
					component: AuditPage,
					meta: { title: 'Audit Trail', permission: 'audit.view' }
				}
			]
		}
	]
})

router.beforeEach(async (to, _from, next) => {
	const authStore = useAuthStore()

	// 1. Ensure setup status is known
	if (authStore.isInitialized === null) {
		await authStore.checkSetup()
	}

	// 2. If system is not initialized, only /setup is accessible
	if (!authStore.isInitialized) {
		if (to.path !== '/setup') {
			return next('/setup')
		}
		return next()
	}

	// 3. If system is initialized, /setup cannot be accessed again
	if (to.path === '/setup') {
		return next('/login')
	}

	// 4. Check Authentication
	if (!to.meta.public && !authStore.isAuthenticated) {
		return next('/login')
	}

	// 5. If already authenticated and visiting /login, redirect to /
	if (to.path === '/login' && authStore.isAuthenticated) {
		return next('/')
	}

	// 6. Check specific permission
	const requiredPerm = to.meta.permission as string | undefined
	if (requiredPerm && !authStore.hasPermission(requiredPerm)) {
		return next('/')
	}

	next()
})

export default router

