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
	studCmd "school-management/internal/application/commands/students"
	feeQuery "school-management/internal/application/queries/fees"
	"school-management/internal/application/services/authz"
	"school-management/internal/domain/academics"
	"school-management/internal/domain/fees"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
	"school-management/internal/infrastructure/database/sqlite"
	"school-management/internal/platform/config"
	"school-management/internal/platform/security"
)

func TestFeeEngineAndInvoicingLifecycle(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sms-phase4-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := config.DefaultConfig(tempDir)
	cfg.Database.Path = filepath.Join(tempDir, "phase4_test.db")
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

	// Step 1: School setup & login
	_, err = setupHandler.Execute(ctx, authCmd.FirstRunSetupInput{
		SchoolName:         "Oakridge International School",
		SchoolCode:         "OIS",
		Timezone:           "UTC",
		CurrencyCode:       "USD",
		CurrencyMinorUnits: 2,
		AdminUsername:      "superadmin",
		AdminDisplayName:   "System Admin",
		AdminEmail:         "admin@oakridge.edu",
		AdminPassword:      "Password123!",
	})
	if err != nil {
		t.Fatalf("first run setup failed: %v", err)
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

	// Step 2: Create Academic Year, Class, Section
	yearID, err := acadCmdHandler.CreateAcademicYear(authCtx, acadCmd.CreateAcademicYearInput{
		Name:     "Academic Year 2026-2027",
		Code:     "2026-27",
		StartsOn: shared.CivilDate(2026, 9, 1),
		EndsOn:   shared.CivilDate(2027, 6, 30),
		Status:   academics.YearActive,
	})
	if err != nil {
		t.Fatalf("failed to create academic year: %v", err)
	}

	classID, err := acadCmdHandler.CreateClass(authCtx, acadCmd.CreateClassInput{
		Name:       "Grade 1",
		Code:       "G1",
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

	// Step 3: Register 2 students (Alice & Bob)
	aliceID, err := studCmdHandler.RegisterStudent(authCtx, studCmd.RegisterStudentInput{
		AdmissionNo:    "ADM-2026-0001",
		FirstName:      "Alice",
		LastName:       "Wonderland",
		DateOfBirth:    shared.CivilDate(2020, 5, 10),
		Gender:         "female",
		AdmissionDate:  shared.CivilDate(2026, 9, 1),
		AcademicYearID: yearID,
		ClassID:        classID,
		SectionID:      &sectionID,
	})
	if err != nil {
		t.Fatalf("failed to register student Alice: %v", err)
	}

	bobID, err := studCmdHandler.RegisterStudent(authCtx, studCmd.RegisterStudentInput{
		AdmissionNo:    "ADM-2026-0002",
		FirstName:      "Bob",
		LastName:       "Builder",
		DateOfBirth:    shared.CivilDate(2020, 8, 20),
		Gender:         "male",
		AdmissionDate:  shared.CivilDate(2026, 9, 1),
		AcademicYearID: yearID,
		ClassID:        classID,
		SectionID:      &sectionID,
	})
	if err != nil {
		t.Fatalf("failed to register student Bob: %v", err)
	}

	// Step 4: Create Fee Heads
	tuitionHeadID, err := feeCmdHandler.CreateFeeHead(authCtx, feeCmd.CreateFeeHeadInput{
		Code:               "TUITION",
		Name:               "Tuition Fee",
		Description:        "Monthly core instruction fee",
		IsRecurring:        true,
		DefaultAmountMinor: 50000, // $500.00
	})
	if err != nil {
		t.Fatalf("failed to create tuition fee head: %v", err)
	}

	labHeadID, err := feeCmdHandler.CreateFeeHead(authCtx, feeCmd.CreateFeeHeadInput{
		Code:               "COMPUTER_LAB",
		Name:               "Computer Lab Fee",
		Description:        "Monthly technology fee",
		IsRecurring:        true,
		DefaultAmountMinor: 10000, // $100.00
	})
	if err != nil {
		t.Fatalf("failed to create lab fee head: %v", err)
	}

	// Step 5: Create Fee Structure for Grade 1
	feeStructureID, err := feeCmdHandler.CreateFeeStructure(authCtx, feeCmd.CreateFeeStructureInput{
		AcademicYearID: yearID,
		ClassID:        classID,
		Code:           "FS-G1-2026",
		Name:           "Grade 1 Fee Structure 2026-27",
		Description:    "Standard fees for Grade 1 students",
	})
	if err != nil {
		t.Fatalf("failed to create fee structure: %v", err)
	}

	// Create Version with line items
	versionID, err := feeCmdHandler.CreateFeeStructureVersion(authCtx, feeCmd.CreateFeeStructureVersionInput{
		FeeStructureID: feeStructureID,
		Changelog:      "Initial 2026-27 structure",
		Items: []feeCmd.FeeStructureItemInput{
			{
				FeeHeadID:   tuitionHeadID,
				AmountMinor: 50000,
				Frequency:   fees.FrequencyMonthly,
				DueDay:      10,
			},
			{
				FeeHeadID:   labHeadID,
				AmountMinor: 10000,
				Frequency:   fees.FrequencyMonthly,
				DueDay:      10,
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create fee structure version: %v", err)
	}

	// Publish Version
	err = feeCmdHandler.PublishFeeStructureVersion(authCtx, versionID)
	if err != nil {
		t.Fatalf("failed to publish fee structure version: %v", err)
	}

	// Verify published version detail
	detail, err := feeQueries.GetFeeStructureVersionDetail(authCtx, versionID)
	if err != nil {
		t.Fatalf("failed to get version detail: %v", err)
	}
	if detail.Version.Status != fees.VersionPublished {
		t.Errorf("expected version status published, got %s", detail.Version.Status)
	}
	if len(detail.Items) != 2 {
		t.Errorf("expected 2 fee structure items, got %d", len(detail.Items))
	}

	// Step 6: Create Sibling Discount Policy (20% off Tuition)
	discountPolicyID, err := feeCmdHandler.CreateDiscountPolicy(authCtx, feeCmd.CreateDiscountPolicyInput{
		Code:             "SIBLING_20",
		Name:             "Sibling Concession 20%",
		Description:      "20% discount on tuition for second child",
		DiscountType:     fees.DiscountPercentage,
		RateBasisPoints:  2000, // 20.00%
		FeeHeadID:        &tuitionHeadID,
		RequiresApproval: true,
	})
	if err != nil {
		t.Fatalf("failed to create discount policy: %v", err)
	}

	// Apply discount to Bob
	studentDiscountID, err := feeCmdHandler.ApplyStudentDiscount(authCtx, feeCmd.ApplyStudentDiscountInput{
		StudentID:        bobID,
		DiscountPolicyID: discountPolicyID,
		AcademicYearID:   yearID,
		ValidFrom:        shared.CivilDate(2026, 9, 1),
	})
	if err != nil {
		t.Fatalf("failed to apply student discount: %v", err)
	}

	// Approve discount
	err = feeCmdHandler.DecideStudentDiscount(authCtx, feeCmd.DecideStudentDiscountInput{
		DiscountID: studentDiscountID,
		Approved:   true,
	})
	if err != nil {
		t.Fatalf("failed to approve student discount: %v", err)
	}

	// Step 7: Batch Generate Invoices for Grade 1 (September 2026)
	invoiceIDs, err := feeCmdHandler.GenerateInvoicesBatch(authCtx, feeCmd.GenerateInvoicesBatchInput{
		AcademicYearID: yearID,
		ClassID:        &classID,
		BillingPeriod:  "2026-09",
		IssueDate:      shared.CivilDate(2026, 9, 1),
		DueDate:        shared.CivilDate(2026, 9, 10),
	})
	if err != nil {
		t.Fatalf("failed to generate invoices batch: %v", err)
	}
	if len(invoiceIDs) != 2 {
		t.Fatalf("expected 2 invoices generated, got %d", len(invoiceIDs))
	}

	// Verify Alice's invoice (No discount)
	var aliceInvoice *feeQuery.InvoiceDetail
	var bobInvoice *feeQuery.InvoiceDetail

	for _, invID := range invoiceIDs {
		inv, err := feeQueries.GetInvoice(authCtx, invID)
		if err != nil {
			t.Fatalf("failed to get invoice %s: %v", invID, err)
		}
		if inv.Invoice.StudentID == aliceID {
			aliceInvoice = inv
		} else if inv.Invoice.StudentID == bobID {
			bobInvoice = inv
		}
	}

	if aliceInvoice == nil || bobInvoice == nil {
		t.Fatalf("could not find both Alice and Bob invoices")
	}

	// Alice: Subtotal=60000, Discount=0, Total=60000, Balance=60000
	if aliceInvoice.Invoice.SubtotalAmountMinor != 60000 {
		t.Errorf("Alice subtotal mismatch: want 60000, got %d", aliceInvoice.Invoice.SubtotalAmountMinor)
	}
	if aliceInvoice.Invoice.DiscountAmountMinor != 0 {
		t.Errorf("Alice discount mismatch: want 0, got %d", aliceInvoice.Invoice.DiscountAmountMinor)
	}
	if aliceInvoice.Invoice.TotalAmountMinor != 60000 {
		t.Errorf("Alice total mismatch: want 60000, got %d", aliceInvoice.Invoice.TotalAmountMinor)
	}
	if aliceInvoice.Invoice.BalanceMinor != 60000 {
		t.Errorf("Alice balance mismatch: want 60000, got %d", aliceInvoice.Invoice.BalanceMinor)
	}

	// Bob: Tuition 50000 with 20% discount (10000) -> 40000 net; Lab 10000 -> 10000 net
	// Bob Subtotal=60000, Discount=10000, Total=50000, Balance=50000
	if bobInvoice.Invoice.SubtotalAmountMinor != 60000 {
		t.Errorf("Bob subtotal mismatch: want 60000, got %d", bobInvoice.Invoice.SubtotalAmountMinor)
	}
	if bobInvoice.Invoice.DiscountAmountMinor != 10000 {
		t.Errorf("Bob discount mismatch: want 10000, got %d", bobInvoice.Invoice.DiscountAmountMinor)
	}
	if bobInvoice.Invoice.TotalAmountMinor != 50000 {
		t.Errorf("Bob total mismatch: want 50000, got %d", bobInvoice.Invoice.TotalAmountMinor)
	}
	if bobInvoice.Invoice.BalanceMinor != 50000 {
		t.Errorf("Bob balance mismatch: want 50000, got %d", bobInvoice.Invoice.BalanceMinor)
	}

	// Verify line items on Bob's invoice
	for _, item := range bobInvoice.Items {
		if item.FeeHeadID == tuitionHeadID {
			if item.OriginalAmountMinor != 50000 {
				t.Errorf("Bob tuition original amount want 50000, got %d", item.OriginalAmountMinor)
			}
			if item.DiscountAmountMinor != 10000 {
				t.Errorf("Bob tuition discount want 10000, got %d", item.DiscountAmountMinor)
			}
			if item.NetAmountMinor != 40000 {
				t.Errorf("Bob tuition net amount want 40000, got %d", item.NetAmountMinor)
			}
		}
	}

	// Step 8: Void Bob's invoice
	err = feeCmdHandler.VoidInvoice(authCtx, feeCmd.VoidInvoiceInput{
		InvoiceID: bobInvoice.Invoice.ID,
		Reason:    "Regenerating due to revised scholarship policy",
	})
	if err != nil {
		t.Fatalf("failed to void invoice: %v", err)
	}

	voidedBobInv, err := feeQueries.GetInvoice(authCtx, bobInvoice.Invoice.ID)
	if err != nil {
		t.Fatalf("failed to get voided invoice: %v", err)
	}
	if voidedBobInv.Invoice.Status != fees.InvoiceVoid {
		t.Errorf("expected invoice status void, got %s", voidedBobInv.Invoice.Status)
	}
	if voidedBobInv.Invoice.VoidReason != "Regenerating due to revised scholarship policy" {
		t.Errorf("unexpected void reason: %s", voidedBobInv.Invoice.VoidReason)
	}

	t.Logf("Phase 4 fee engine integration test passed! Alice Total=%d, Bob Total=%d (Discount=%d), Invoices=%d",
		aliceInvoice.Invoice.TotalAmountMinor, bobInvoice.Invoice.TotalAmountMinor, bobInvoice.Invoice.DiscountAmountMinor, len(invoiceIDs))
}
