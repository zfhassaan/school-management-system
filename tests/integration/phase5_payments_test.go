package integration_test

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	acadCmd "school-management/internal/application/commands/academics"
	authCmd "school-management/internal/application/commands/auth"
	feeCmd "school-management/internal/application/commands/fees"
	paymentCmd "school-management/internal/application/commands/payments"
	studCmd "school-management/internal/application/commands/students"
	feeQuery "school-management/internal/application/queries/fees"
	paymentQuery "school-management/internal/application/queries/payments"
	"school-management/internal/application/services/authz"
	"school-management/internal/domain/academics"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/payments"
	"school-management/internal/domain/shared"
	"school-management/internal/domain/students"
	"school-management/internal/infrastructure/database/sqlite"
	"school-management/internal/platform/config"
	"school-management/internal/platform/security"
)

func TestPaymentLifecycleAndCreditEngine(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sms-phase5-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := config.DefaultConfig(tempDir)
	cfg.Database.Path = filepath.Join(tempDir, "phase5_test.db")
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	ctx := context.Background()
	db, err := sqlite.Open(ctx, &cfg, logger)
	if err != nil {
		t.Fatalf("failed to open sqlite database: %v", err)
	}
	defer db.Close()

	uow := sqlite.NewUnitOfWork(db)
	sessionStore := sqlite.NewSessionStore(db)
	clock := shared.SystemClock{}
	ids := shared.UUIDv7Generator{}
	hasher := security.NewArgon2Hasher(security.Argon2Params{
		Memory:      1024,
		Iterations:  1,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	})

	setupHandler := authCmd.NewFirstRunSetupHandler(uow, ids, clock, hasher)
	loginHandler := authCmd.NewLoginHandler(uow, sessionStore, ids, clock, hasher, db.DeviceID())
	acadCmdHandler := acadCmd.NewAcademicCommandsHandler(uow, clock, ids)
	studCmdHandler := studCmd.NewStudentCommandsHandler(uow, clock, ids)
	feeCmdHandler := feeCmd.NewFeeCommandsHandler(uow, clock, ids)
	feeQueries := feeQuery.NewFeeQueriesHandler(uow)
	paymentCmdHandler := paymentCmd.NewPaymentCommandsHandler(uow, clock, ids)
	paymentQueries := paymentQuery.NewPaymentQueriesHandler(uow)

	// Step 1: School setup & login
	_, err = setupHandler.Execute(ctx, authCmd.FirstRunSetupInput{
		SchoolName:         "Financial Academy",
		SchoolCode:         "FA-01",
		Timezone:           "UTC",
		CurrencyCode:       "USD",
		CurrencyMinorUnits: 2,
		AdminUsername:      "superadmin",
		AdminDisplayName:   "Chief Bursar",
		AdminEmail:         "bursar@fa.edu",
		AdminPassword:      "Password123!",
	})
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	loginRes, err := loginHandler.Execute(ctx, authCmd.LoginInput{
		Username: "superadmin",
		Password: "Password123!",
	})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	adminUserID, _ := shared.ParseID(loginRes.User.ID)
	adminPerms := make([]identity.Permission, len(loginRes.User.Permissions))
	for i, p := range loginRes.User.Permissions {
		adminPerms[i] = identity.Permission(p)
	}
	authCtx := authz.WithActor(ctx, authz.Actor{
		UserID:      adminUserID,
		Username:    loginRes.User.Username,
		DisplayName: loginRes.User.DisplayName,
		DeviceID:    db.DeviceID(),
		SessionID:   loginRes.SessionToken,
		Roles:       loginRes.User.Roles,
		Permissions: identity.NewPermissionSet(adminPerms),
	})

	// Step 2: Setup Academic Year & Class
	yearID, err := acadCmdHandler.CreateAcademicYear(authCtx, acadCmd.CreateAcademicYearInput{
		Name:     "Academic Year 2026-2027",
		Code:     "AY2026",
		StartsOn: shared.CivilDate(2026, 9, 1),
		EndsOn:   shared.CivilDate(2027, 6, 30),
		Status:   academics.YearActive,
	})
	if err != nil {
		t.Fatalf("failed to create academic year: %v", err)
	}

	classID, err := acadCmdHandler.CreateClass(authCtx, acadCmd.CreateClassInput{
		Code:       "G1",
		Name:       "Grade 1",
		LevelOrder: 1,
	})
	if err != nil {
		t.Fatalf("failed to create class: %v", err)
	}

	cap40 := 40
	sectionID, err := acadCmdHandler.CreateSection(authCtx, acadCmd.CreateSectionInput{
		ClassID:        classID,
		AcademicYearID: yearID,
		Name:           "Section A",
		Capacity:       &cap40,
	})
	if err != nil {
		t.Fatalf("failed to create section: %v", err)
	}

	// Step 3: Register Student Alice
	aliceID, err := studCmdHandler.RegisterStudent(authCtx, studCmd.RegisterStudentInput{
		AdmissionNo:    "ADM-P5-001",
		FirstName:      "Alice",
		LastName:       "Wonderland",
		DateOfBirth:    shared.CivilDate(2018, 5, 10),
		Gender:         students.GenderFemale,
		AdmissionDate:  shared.CivilDate(2026, 9, 1),
		AcademicYearID: yearID,
		ClassID:        classID,
		SectionID:      &sectionID,
	})
	if err != nil {
		t.Fatalf("failed to register Alice: %v", err)
	}

	// Step 4: Create Fee Head, Structure & Publish Version
	headID, err := feeCmdHandler.CreateFeeHead(authCtx, feeCmd.CreateFeeHeadInput{
		Code:               "TUITION",
		Name:               "Monthly Tuition",
		IsRecurring:        true,
		DefaultAmountMinor: 60000, // $600.00
	})
	if err != nil {
		t.Fatalf("failed to create fee head: %v", err)
	}

	fsID, err := feeCmdHandler.CreateFeeStructure(authCtx, feeCmd.CreateFeeStructureInput{
		AcademicYearID: yearID,
		ClassID:        classID,
		Code:           "FS-G1-2026",
		Name:           "Grade 1 Fees",
	})
	if err != nil {
		t.Fatalf("failed to create fee structure: %v", err)
	}

	versionID, err := feeCmdHandler.CreateFeeStructureVersion(authCtx, feeCmd.CreateFeeStructureVersionInput{
		FeeStructureID: fsID,
		Changelog:      "Initial schedule",
		Items: []feeCmd.FeeStructureItemInput{
			{
				FeeHeadID:   headID,
				AmountMinor: 60000,
				Frequency:   "monthly",
				DueDay:      10,
				IsOptional:  false,
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create version: %v", err)
	}

	if err := feeCmdHandler.PublishFeeStructureVersion(authCtx, versionID); err != nil {
		t.Fatalf("failed to publish version: %v", err)
	}

	// Step 5: Batch Generate Invoices
	invIDs, err := feeCmdHandler.GenerateInvoicesBatch(authCtx, feeCmd.GenerateInvoicesBatchInput{
		AcademicYearID: yearID,
		ClassID:        &classID,
		BillingPeriod:  "2026-09",
		IssueDate:      shared.CivilDate(2026, 9, 1),
		DueDate:        shared.CivilDate(2026, 9, 15),
	})
	if err != nil {
		t.Fatalf("failed to generate batch invoices: %v", err)
	}
	if len(invIDs) != 1 {
		t.Fatalf("expected 1 invoice, got %d", len(invIDs))
	}
	aliceInvID := invIDs[0]

	// Verify Alice's initial invoice: Balance = $600.00, Paid = $0, Status = issued
	invDetail, err := feeQueries.GetInvoice(authCtx, aliceInvID)
	if err != nil {
		t.Fatalf("failed to get invoice: %v", err)
	}
	if invDetail.Invoice.TotalAmountMinor != 60000 || invDetail.Invoice.BalanceMinor != 60000 || invDetail.Invoice.Status != "issued" {
		t.Fatalf("unexpected invoice state: %+v", invDetail.Invoice)
	}

	// Step 6: Partial Payment: Alice pays $200.00 (20000 minor units)
	payRes1, err := paymentCmdHandler.RecordPayment(authCtx, paymentCmd.RecordPaymentInput{
		StudentID:     aliceID,
		PaymentDate:   "2026-09-05",
		PaymentMethod: payments.MethodCash,
		AmountMinor:   20000,
		PayerName:     "Alice Guardian",
		Notes:         "First instalment",
	})
	if err != nil {
		t.Fatalf("failed to record payment 1: %v", err)
	}
	if payRes1.ReceiptNo != "RCP-2026-000001" {
		t.Fatalf("expected receipt number RCP-2026-000001, got %s", payRes1.ReceiptNo)
	}
	if payRes1.Allocated != 20000 || payRes1.Unallocated != 0 {
		t.Fatalf("expected allocated 20000, unallocated 0, got allocated=%d, unallocated=%d", payRes1.Allocated, payRes1.Unallocated)
	}

	// Verify invoice after payment 1: Paid = 20000, Balance = 40000, Status = partially_paid
	invDetail, _ = feeQueries.GetInvoice(authCtx, aliceInvID)
	if invDetail.Invoice.PaidAmountMinor != 20000 || invDetail.Invoice.BalanceMinor != 40000 || invDetail.Invoice.Status != "partially_paid" {
		t.Fatalf("expected partially_paid with balance 40000, got status=%s, balance=%d", invDetail.Invoice.Status, invDetail.Invoice.BalanceMinor)
	}

	// Step 7: Overpayment: Alice pays $500.00 when only $400.00 is due
	payRes2, err := paymentCmdHandler.RecordPayment(authCtx, paymentCmd.RecordPaymentInput{
		StudentID:     aliceID,
		PaymentDate:   "2026-09-08",
		PaymentMethod: payments.MethodBankTransfer,
		AmountMinor:   50000,
		ReferenceNo:   "TXN-998877",
		PayerName:     "Alice Guardian",
	})
	if err != nil {
		t.Fatalf("failed to record payment 2: %v", err)
	}
	if payRes2.ReceiptNo != "RCP-2026-000002" {
		t.Fatalf("expected receipt number RCP-2026-000002, got %s", payRes2.ReceiptNo)
	}
	if payRes2.Allocated != 40000 || payRes2.Unallocated != 10000 {
		t.Fatalf("expected allocated 40000, unallocated 10000, got allocated=%d, unallocated=%d", payRes2.Allocated, payRes2.Unallocated)
	}

	// Verify invoice after payment 2: Paid = 60000, Balance = 0, Status = paid
	invDetail, _ = feeQueries.GetInvoice(authCtx, aliceInvID)
	if invDetail.Invoice.PaidAmountMinor != 60000 || invDetail.Invoice.BalanceMinor != 0 || invDetail.Invoice.Status != "paid" {
		t.Fatalf("expected paid with balance 0, got status=%s, balance=%d", invDetail.Invoice.Status, invDetail.Invoice.BalanceMinor)
	}

	// Verify student credit balance received the $100.00 (10000 minor units) overpayment
	creditInfo, err := paymentQueries.GetStudentCreditBalance(authCtx, aliceID)
	if err != nil {
		t.Fatalf("failed to query credit balance: %v", err)
	}
	if creditInfo.BalanceMinor != 10000 {
		t.Fatalf("expected credit balance 10000, got %d", creditInfo.BalanceMinor)
	}
	if len(creditInfo.Ledger) != 1 || creditInfo.Ledger[0].AmountMinor != 10000 {
		t.Fatalf("unexpected credit ledger entry: %+v", creditInfo.Ledger)
	}

	// Step 8: Payment Reversal: Reverse payment 1 ($200.00)
	err = paymentCmdHandler.ReversePayment(authCtx, paymentCmd.ReversePaymentInput{
		PaymentID: payRes1.PaymentID,
		Reason:    "Bounced cheque / transaction cancellation",
	})
	if err != nil {
		t.Fatalf("failed to reverse payment 1: %v", err)
	}

	// Verify invoice is reinstated: Paid = 40000, Balance = 20000, Status = partially_paid
	invDetail, _ = feeQueries.GetInvoice(authCtx, aliceInvID)
	if invDetail.Invoice.PaidAmountMinor != 40000 || invDetail.Invoice.BalanceMinor != 20000 || invDetail.Invoice.Status != "partially_paid" {
		t.Fatalf("expected invoice reinstated to partially_paid balance 20000, got status=%s, balance=%d", invDetail.Invoice.Status, invDetail.Invoice.BalanceMinor)
	}

	// Verify payment 1 detail reflects reversal
	pDetail1, err := paymentQueries.GetPayment(authCtx, payRes1.PaymentID)
	if err != nil {
		t.Fatalf("failed to get payment 1 detail: %v", err)
	}
	if pDetail1.Payment.Status != "reversed" || pDetail1.Reversal == nil || pDetail1.Reversal.Reason != "Bounced cheque / transaction cancellation" {
		t.Fatalf("expected payment reversal record, got: %+v", pDetail1)
	}

	// Step 9: Re-reversing the same payment should be prevented
	err = paymentCmdHandler.ReversePayment(authCtx, paymentCmd.ReversePaymentInput{
		PaymentID: payRes1.PaymentID,
		Reason:    "Duplicate reversal attempt",
	})
	if err == nil {
		t.Fatalf("expected conflict on re-reversing, got nil")
	}

	// Step 10: Financial Adjustment (Waiver): Waive remaining $200.00 balance
	adjID, err := paymentCmdHandler.RequestFinancialAdjustment(authCtx, paymentCmd.RequestFinancialAdjustmentInput{
		StudentID:      aliceID,
		InvoiceID:      &aliceInvID,
		AdjustmentType: payments.AdjWaiver,
		AmountMinor:    20000,
		Reason:         "Hardship concession granted by Principal",
	})
	if err != nil {
		t.Fatalf("failed to request financial adjustment: %v", err)
	}

	// Approve the adjustment
	err = paymentCmdHandler.DecideFinancialAdjustment(authCtx, paymentCmd.DecideFinancialAdjustmentInput{
		AdjustmentID:  adjID,
		Approved:      true,
		DecisionNotes: "Approved under discretionary hardship policy",
	})
	if err != nil {
		t.Fatalf("failed to decide financial adjustment: %v", err)
	}

	// Verify invoice is now fully settled (balance 0, status paid)
	invDetail, _ = feeQueries.GetInvoice(authCtx, aliceInvID)
	if invDetail.Invoice.BalanceMinor != 0 || invDetail.Invoice.Status != "paid" {
		t.Fatalf("expected waived invoice to be paid with balance 0, got status=%s, balance=%d", invDetail.Invoice.Status, invDetail.Invoice.BalanceMinor)
	}

	t.Logf("Phase 5 payment, allocation, credit engine, and adjustment integration test passed successfully!")
}
