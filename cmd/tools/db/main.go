package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"

	acadCmd "school-management/internal/application/commands/academics"
	admCmd "school-management/internal/application/commands/admissions"
	authCmd "school-management/internal/application/commands/auth"
	feeCmd "school-management/internal/application/commands/fees"
	guardCmd "school-management/internal/application/commands/guardians"
	paymentCmd "school-management/internal/application/commands/payments"
	studCmd "school-management/internal/application/commands/students"
	"school-management/internal/application/services/authz"
	"school-management/internal/domain/academics"
	"school-management/internal/domain/admissions"
	"school-management/internal/domain/guardians"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/payments"
	"school-management/internal/domain/shared"
	"school-management/internal/domain/students"
	"school-management/internal/infrastructure/database/sqlite"
	"school-management/internal/platform/config"
	"school-management/internal/platform/security"
)

func main() {
	freshFlag := flag.Bool("fresh", false, "Drop and re-run all migrations from scratch")
	seedFlag := flag.Bool("seed", false, "Seed realistic demo school data")
	dbPathFlag := flag.String("db", "", "Custom path to SQLite database (optional)")
	flag.Parse()

	if !*freshFlag && !*seedFlag {
		fmt.Println("Usage: go run ./cmd/tools/db [-fresh] [-seed] [-db <path>]")
		fmt.Println("  -fresh: Deletes existing SQLite database and re-runs all Goose migrations")
		fmt.Println("  -seed:  Populates school, academic structures, students, fees, and sample payments")
		os.Exit(0)
	}

	cfg, err := config.Load("")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	if *dbPathFlag != "" {
		cfg.Database.Path = *dbPathFlag
	}

	dbPath := cfg.Database.Path
	ctx := context.Background()

	// 1. Fresh: Remove existing database files or drop all tables
	if *freshFlag {
		fmt.Printf("[RESET] Performing fresh migration on: %s\n", dbPath)
		err1 := os.Remove(dbPath)
		_ = os.Remove(dbPath + "-wal")
		_ = os.Remove(dbPath + "-shm")

		if err1 != nil && !os.IsNotExist(err1) {
			fmt.Println("[WARN] Database file is currently open by running app. Resetting tables via SQL connection...")
			cleanCfg := cfg
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			rawDB, openErr := sqlite.Open(ctx, cleanCfg, logger)
			if openErr == nil {
				_, _ = rawDB.WriterDB().ExecContext(ctx, "PRAGMA foreign_keys = OFF;")
				rows, qErr := rawDB.WriterDB().QueryContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%';")
				if qErr == nil {
					var tables []string
					for rows.Next() {
						var tbl string
						if err := rows.Scan(&tbl); err == nil {
							tables = append(tables, tbl)
						}
					}
					rows.Close()
					for _, tbl := range tables {
						_, _ = rawDB.WriterDB().ExecContext(ctx, fmt.Sprintf("DROP TABLE IF EXISTS %s;", tbl))
					}
				}
				_, _ = rawDB.WriterDB().ExecContext(ctx, "PRAGMA foreign_keys = ON;")
				rawDB.Close()
			}
		}
		fmt.Println("[OK] Existing database cleared.")
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	db, err := sqlite.Open(ctx, cfg, logger)
	if err != nil {
		fmt.Printf("[ERROR] Failed to open database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	fmt.Println("[OK] Embedded Goose migrations applied up to latest version.")

	// 2. Seed: Populate demo data
	if *seedFlag {
		fmt.Println("[SEED] Seeding comprehensive demo dataset...")
		if err := seedDatabase(ctx, db); err != nil {
			fmt.Printf("[ERROR] Seeding failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("[SUCCESS] Database seeded successfully!")
	}
}

func seedDatabase(ctx context.Context, db *sqlite.DB) error {
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
	guardCmdHandler := guardCmd.NewGuardianCommandsHandler(uow, clock, ids)
	feeCmdHandler := feeCmd.NewFeeCommandsHandler(uow, clock, ids)
	paymentCmdHandler := paymentCmd.NewPaymentCommandsHandler(uow, clock, ids)

	// Step 1: First-Run Setup
	fmt.Println("   • Creating school profile and superadmin account...")
	setupRes, err := setupHandler.Execute(ctx, authCmd.FirstRunSetupInput{
		SchoolName:         "Oakridge International Academy",
		SchoolCode:         "OAK-01",
		Timezone:           "Asia/Karachi",
		CurrencyCode:       "PKR",
		CurrencyMinorUnits: 2,
		AdminUsername:      "admin",
		AdminDisplayName:   "Principal Admin",
		AdminEmail:         "admin@oakridge.edu",
		AdminPassword:      "Password123!",
	})
	if err != nil {
		return fmt.Errorf("setup failed: %w", err)
	}

	// Login as admin
	loginRes, err := loginHandler.Execute(ctx, authCmd.LoginInput{
		Username: "admin",
		Password: "Password123!",
	})
	if err != nil {
		return fmt.Errorf("admin login failed: %w", err)
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

	// Step 2: Academic Year & Campuses
	fmt.Println("   • Creating academic years, terms, and campuses...")
	yearID, err := acadCmdHandler.CreateAcademicYear(authCtx, acadCmd.CreateAcademicYearInput{
		Name:     "Academic Year 2026-2027",
		Code:     "2026-27",
		StartsOn: shared.CivilDate(2026, 9, 1),
		EndsOn:   shared.CivilDate(2027, 6, 30),
		Status:   academics.YearActive,
	})
	if err != nil {
		return err
	}

	_ = acadCmdHandler.SetTerms(authCtx, acadCmd.SetTermsInput{
		AcademicYearID: yearID,
		Terms: []acadCmd.AcademicTermInput{
			{Name: "Term 1 (Fall)", Sequence: 1, StartsOn: shared.CivilDate(2026, 9, 1), EndsOn: shared.CivilDate(2026, 12, 15)},
			{Name: "Term 2 (Winter)", Sequence: 2, StartsOn: shared.CivilDate(2027, 1, 5), EndsOn: shared.CivilDate(2027, 3, 26)},
			{Name: "Term 3 (Spring)", Sequence: 3, StartsOn: shared.CivilDate(2027, 4, 6), EndsOn: shared.CivilDate(2027, 6, 30)},
		},
	})

	mainCampusID, err := acadCmdHandler.CreateCampus(authCtx, acadCmd.CreateCampusInput{
		Code: "MAIN",
		Name: "Main Campus",
	})
	if err != nil {
		return err
	}

	// Step 3: Classes & Sections
	fmt.Println("   • Creating classes (Grade 1, Grade 2, Grade 3) and sections...")
	g1ID, err := acadCmdHandler.CreateClass(authCtx, acadCmd.CreateClassInput{
		Code:       "G1",
		Name:       "Grade 1",
		LevelOrder: 1,
	})
	if err != nil {
		return err
	}

	g2ID, err := acadCmdHandler.CreateClass(authCtx, acadCmd.CreateClassInput{
		Code:       "G2",
		Name:       "Grade 2",
		LevelOrder: 2,
	})
	if err != nil {
		return err
	}

	g3ID, err := acadCmdHandler.CreateClass(authCtx, acadCmd.CreateClassInput{
		Code:       "G3",
		Name:       "Grade 3",
		LevelOrder: 3,
	})
	if err != nil {
		return err
	}

	cap30 := 30
	sec1AID, err := acadCmdHandler.CreateSection(authCtx, acadCmd.CreateSectionInput{
		ClassID:        g1ID,
		AcademicYearID: yearID,
		CampusID:       &mainCampusID,
		Name:           "Section A",
		Capacity:       &cap30,
	})
	if err != nil {
		return err
	}

	sec1BID, err := acadCmdHandler.CreateSection(authCtx, acadCmd.CreateSectionInput{
		ClassID:        g1ID,
		AcademicYearID: yearID,
		CampusID:       &mainCampusID,
		Name:           "Section B",
		Capacity:       &cap30,
	})
	if err != nil {
		return err
	}

	sec2AID, err := acadCmdHandler.CreateSection(authCtx, acadCmd.CreateSectionInput{
		ClassID:        g2ID,
		AcademicYearID: yearID,
		CampusID:       &mainCampusID,
		Name:           "Section A",
		Capacity:       &cap30,
	})
	if err != nil {
		return err
	}

	// Step 4: Students & Guardians
	fmt.Println("   • Registering students (Alice, Bob, Charlie, Diana)...")
	aliceID, err := studCmdHandler.RegisterStudent(authCtx, studCmd.RegisterStudentInput{
		AdmissionNo:    "ADM-2026-0001",
		FirstName:      "Alice",
		LastName:       "Wonderland",
		DateOfBirth:    shared.CivilDate(2020, 5, 10),
		Gender:         students.GenderFemale,
		AdmissionDate:  shared.CivilDate(2026, 9, 1),
		AcademicYearID: yearID,
		ClassID:        g1ID,
		SectionID:      &sec1AID,
	})
	if err != nil {
		return err
	}

	bobID, err := studCmdHandler.RegisterStudent(authCtx, studCmd.RegisterStudentInput{
		AdmissionNo:    "ADM-2026-0002",
		FirstName:      "Bob",
		LastName:       "Builder",
		DateOfBirth:    shared.CivilDate(2020, 8, 20),
		Gender:         students.GenderMale,
		AdmissionDate:  shared.CivilDate(2026, 9, 1),
		AcademicYearID: yearID,
		ClassID:        g1ID,
		SectionID:      &sec1AID,
	})
	if err != nil {
		return err
	}

	_, err = studCmdHandler.RegisterStudent(authCtx, studCmd.RegisterStudentInput{
		AdmissionNo:    "ADM-2026-0003",
		FirstName:      "Charlie",
		LastName:       "Brown",
		DateOfBirth:    shared.CivilDate(2020, 3, 14),
		Gender:         students.GenderMale,
		AdmissionDate:  shared.CivilDate(2026, 9, 1),
		AcademicYearID: yearID,
		ClassID:        g1ID,
		SectionID:      &sec1BID,
	})
	if err != nil {
		return err
	}

	_, err = studCmdHandler.RegisterStudent(authCtx, studCmd.RegisterStudentInput{
		AdmissionNo:    "ADM-2026-0004",
		FirstName:      "Diana",
		LastName:       "Prince",
		DateOfBirth:    shared.CivilDate(2019, 11, 2),
		Gender:         students.GenderFemale,
		AdmissionDate:  shared.CivilDate(2026, 9, 1),
		AcademicYearID: yearID,
		ClassID:        g2ID,
		SectionID:      &sec2AID,
	})
	if err != nil {
		return err
	}

	// Guardians
	gAlice, _ := guardCmdHandler.CreateGuardian(authCtx, guardCmd.CreateGuardianInput{
		FirstName:    "Arthur",
		LastName:     "Wonderland",
		PhonePrimary: "+1-555-0101",
		Email:        "arthur@wonderland.org",
	})
	_ = guardCmdHandler.SetStudentGuardians(authCtx, guardCmd.SetStudentGuardiansInput{
		StudentID: aliceID,
		Guardians: []guardCmd.StudentGuardianLinkInput{
			{GuardianID: gAlice, Relationship: guardians.RelFather, IsPrimary: true, ReceivesBilling: true, CanPickup: true, Priority: 1},
		},
	})

	// Step 4.5: Admissions Form & Sample Applications
	fmt.Println("   • Creating admission form and sample applications...")
	admCmdHandler := admCmd.NewAdmissionCommandsHandler(uow, clock, ids)
	formDefID, err := admCmdHandler.CreateFormDefinition(authCtx, admCmd.CreateFormDefInput{
		Code:        "STD-2026",
		Name:        "General Admission Form",
		Description: "General admission registration form for upcoming academic session",
	})
	if err == nil {
		schema := admissions.FormSchema{
			Title: "General Admission Form",
			Sections: []admissions.FormSection{
				{
					ID:    "sec_student",
					Title: "Student Information",
					Fields: []admissions.FormField{
						{ID: "f_prev_school", Label: "Previous School Attended", Type: admissions.FieldText, Required: false},
					},
				},
			},
		}
		verID, err := admCmdHandler.CreateFormVersion(authCtx, admCmd.CreateFormVersionInput{
			FormDefinitionID: formDefID,
			Schema:           schema,
		})
		if err == nil {
			_ = admCmdHandler.PublishFormVersion(authCtx, admCmd.PublishFormVersionInput{
				VersionID: verID,
			})

			// 1. Submitted Application
			app1DOB := shared.CivilDate(2020, 6, 15)
			_, _ = admCmdHandler.SubmitApplication(authCtx, admCmd.SubmitApplicationInput{
				FormVersionID:           verID,
				AcademicYearID:          yearID,
				AppliedClassID:          g1ID,
				ApplicantFirstName:      "Hamza",
				ApplicantLastName:       "Tariq",
				ApplicantDOB:            &app1DOB,
				ApplicantGender:         "male",
				GuardianName:            "Tariq Mahmood",
				GuardianPhoneNormalized: "+923001112233",
				GuardianEmail:           "tariq.m@example.com",
			})

			// 2. Under Review Application
			app2DOB := shared.CivilDate(2019, 9, 25)
			app2ID, err := admCmdHandler.SubmitApplication(authCtx, admCmd.SubmitApplicationInput{
				FormVersionID:           verID,
				AcademicYearID:          yearID,
				AppliedClassID:          g2ID,
				ApplicantFirstName:      "Zainab",
				ApplicantLastName:       "Fatima",
				ApplicantDOB:            &app2DOB,
				ApplicantGender:         "female",
				GuardianName:            "Fatima Noor",
				GuardianPhoneNormalized: "+923214445566",
				GuardianEmail:           "fatima.noor@example.com",
			})
			if err == nil {
				_ = admCmdHandler.UpdateStatus(authCtx, admCmd.UpdateApplicationStatusInput{
					ApplicationID: app2ID,
					NewStatus:     admissions.AppUnderReview,
					Reason:        "Interview scheduled for next Tuesday",
				})
			}

			// 3. Approved Application (Ready to be enrolled)
			app3DOB := shared.CivilDate(2018, 4, 10)
			app3ID, err := admCmdHandler.SubmitApplication(authCtx, admCmd.SubmitApplicationInput{
				FormVersionID:           verID,
				AcademicYearID:          yearID,
				AppliedClassID:          g3ID,
				ApplicantFirstName:      "Bilal",
				ApplicantLastName:       "Ahmed",
				ApplicantDOB:            &app3DOB,
				ApplicantGender:         "male",
				GuardianName:            "Ahmed Hassan",
				GuardianPhoneNormalized: "+923337778899",
				GuardianEmail:           "ahmed.hassan@example.com",
			})
			if err == nil {
				_ = admCmdHandler.UpdateStatus(authCtx, admCmd.UpdateApplicationStatusInput{
					ApplicationID: app3ID,
					NewStatus:     admissions.AppApproved,
					Reason:        "Passed entrance test with 92% marks",
				})
			}
		}
	}

	// Step 5: Fee Heads & Grade 1 Fee Structure
	fmt.Println("   • Configuring fee heads, published schedules, and concessions...")
	tuitionHead, err := feeCmdHandler.CreateFeeHead(authCtx, feeCmd.CreateFeeHeadInput{
		Code:               "TUITION",
		Name:               "Monthly Tuition",
		Description:        "Core academic tuition instruction",
		IsRecurring:        true,
		DefaultAmountMinor: 60000, // $600.00
	})
	if err != nil {
		return err
	}

	labHead, err := feeCmdHandler.CreateFeeHead(authCtx, feeCmd.CreateFeeHeadInput{
		Code:               "LAB",
		Name:               "Computer & STEM Lab",
		Description:        "Technology and laboratory facility access",
		IsRecurring:        true,
		DefaultAmountMinor: 5000, // $50.00
	})
	if err != nil {
		return err
	}

	libraryHead, err := feeCmdHandler.CreateFeeHead(authCtx, feeCmd.CreateFeeHeadInput{
		Code:               "LIBRARY",
		Name:               "Digital Library Fee",
		Description:        "Online research portal and library borrowing",
		IsRecurring:        true,
		DefaultAmountMinor: 3000, // $30.00
	})
	if err != nil {
		return err
	}

	fsG1, err := feeCmdHandler.CreateFeeStructure(authCtx, feeCmd.CreateFeeStructureInput{
		AcademicYearID: yearID,
		ClassID:        g1ID,
		Code:           "FS-G1-2026",
		Name:           "Grade 1 Standard Fee Structure",
		Description:    "Standard fee schedule for all Grade 1 enrollments",
	})
	if err != nil {
		return err
	}

	vG1, err := feeCmdHandler.CreateFeeStructureVersion(authCtx, feeCmd.CreateFeeStructureVersionInput{
		FeeStructureID: fsG1,
		Changelog:      "Initial approved schedule 2026/27",
		Items: []feeCmd.FeeStructureItemInput{
			{FeeHeadID: tuitionHead, AmountMinor: 60000, Frequency: "monthly", DueDay: 10, IsOptional: false},
			{FeeHeadID: labHead, AmountMinor: 5000, Frequency: "monthly", DueDay: 10, IsOptional: false},
			{FeeHeadID: libraryHead, AmountMinor: 3000, Frequency: "monthly", DueDay: 10, IsOptional: false},
		},
	})
	if err != nil {
		return err
	}

	if err := feeCmdHandler.PublishFeeStructureVersion(authCtx, vG1); err != nil {
		return err
	}

	// Step 6: Discount Policies & Assignment
	sibPolicyID, err := feeCmdHandler.CreateDiscountPolicy(authCtx, feeCmd.CreateDiscountPolicyInput{
		Code:             "SIBLING_20",
		Name:             "Sibling Concession 20%",
		Description:      "20% discount on tuition for second and subsequent children",
		DiscountType:     "percentage",
		RateBasisPoints:  2000, // 20.00%
		FeeHeadID:        &tuitionHead,
		RequiresApproval: true,
	})
	if err != nil {
		return err
	}

	_, _ = feeCmdHandler.CreateDiscountPolicy(authCtx, feeCmd.CreateDiscountPolicyInput{
		Code:             "MERIT_50",
		Name:             "Academic Merit Scholarship 50%",
		Description:      "50% discount on tuition for top academic achievers",
		DiscountType:     "percentage",
		RateBasisPoints:  5000, // 50.00%
		FeeHeadID:        &tuitionHead,
		RequiresApproval: true,
	})

	// Assign sibling discount to Bob and approve it
	discBobID, err := feeCmdHandler.ApplyStudentDiscount(authCtx, feeCmd.ApplyStudentDiscountInput{
		StudentID:        bobID,
		DiscountPolicyID: sibPolicyID,
		AcademicYearID:   yearID,
		ValidFrom:        shared.CivilDate(2026, 9, 1),
	})
	if err != nil {
		return err
	}

	_ = feeCmdHandler.DecideStudentDiscount(authCtx, feeCmd.DecideStudentDiscountInput{
		DiscountID: discBobID,
		Approved:   true,
	})

	// Step 7: Batch Invoicing for Grade 1
	fmt.Println("   • Generating Grade 1 invoices for billing period 2026-09...")
	invIDs, err := feeCmdHandler.GenerateInvoicesBatch(authCtx, feeCmd.GenerateInvoicesBatchInput{
		AcademicYearID: yearID,
		ClassID:        &g1ID,
		BillingPeriod:  "2026-09",
		IssueDate:      shared.CivilDate(2026, 9, 1),
		DueDate:        shared.CivilDate(2026, 9, 15),
	})
	if err != nil {
		return err
	}
	fmt.Printf("     -> Generated %d invoices\n", len(invIDs))

	// Step 8: Sample Payments
	fmt.Println("   • Posting sample payments and FIFO allocations...")
	// Alice pays $300.00 (partial payment against $680.00 total)
	payRes1, err := paymentCmdHandler.RecordPayment(authCtx, paymentCmd.RecordPaymentInput{
		StudentID:     aliceID,
		PaymentDate:   "2026-09-05",
		PaymentMethod: payments.MethodBankTransfer,
		AmountMinor:   30000, // Rs 300.00
		ReferenceNo:   "WIRE-889900",
		PayerName:     "Arthur Wonderland",
		Notes:         "First term tuition instalment",
	})
	if err != nil {
		return err
	}
	fmt.Printf("     -> Receipt %s issued for Alice (Rs 300.00 allocated)\n", payRes1.ReceiptNo)

	// Bob pays Rs 600.00 (Bob's bill is Rs 560.00 due to 20% discount on tuition: Rs 480 + Rs 50 + Rs 30 = Rs 560.00)
	// Excess Rs 40.00 goes to Bob's credit balance advance!
	payRes2, err := paymentCmdHandler.RecordPayment(authCtx, paymentCmd.RecordPaymentInput{
		StudentID:     bobID,
		PaymentDate:   "2026-09-08",
		PaymentMethod: payments.MethodCash,
		AmountMinor:   60000, // Rs 600.00
		PayerName:     "Bob Parent",
		Notes:         "Cash payment at bursar counter",
	})
	if err != nil {
		return err
	}
	fmt.Printf("     -> Receipt %s issued for Bob (Fully paid + Rs 40.00 credit advance stored)\n", payRes2.ReceiptNo)

	fmt.Println("\n==================================================================")
	fmt.Println("SCHOOL DEMO ENVIRONMENT READY")
	fmt.Println("==================================================================")
	fmt.Printf("  School Name:      Oakridge International Academy (%s)\n", setupRes.SchoolID)
	fmt.Printf("  Admin Username:   admin\n")
	fmt.Printf("  Admin Password:   Password123!\n")
	fmt.Printf("  Students:         Alice Wonderland, Bob Builder, Charlie Brown, Diana Prince\n")
	fmt.Printf("  Invoices:         Period 2026-09 issued with automated sibling discount\n")
	fmt.Printf("  Receipts:         RCP-2026-000001 (Alice), RCP-2026-000002 (Bob)\n")
	fmt.Println("==================================================================")
	return nil
}
