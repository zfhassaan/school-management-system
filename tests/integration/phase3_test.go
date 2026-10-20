package integration_test

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	acadCmd "school-management/internal/application/commands/academics"
	admCmd "school-management/internal/application/commands/admissions"
	authCmd "school-management/internal/application/commands/auth"
	docCmd "school-management/internal/application/commands/documents"
	enrCmd "school-management/internal/application/commands/enrollments"
	guardCmd "school-management/internal/application/commands/guardians"
	studCmd "school-management/internal/application/commands/students"
	acadQuery "school-management/internal/application/queries/academics"
	admQuery "school-management/internal/application/queries/admissions"
	authQuery "school-management/internal/application/queries/auth"
	docQuery "school-management/internal/application/queries/documents"
	enrQuery "school-management/internal/application/queries/enrollments"
	guardQuery "school-management/internal/application/queries/guardians"
	studQuery "school-management/internal/application/queries/students"
	"school-management/internal/application/services/authz"
	"school-management/internal/domain/academics"
	"school-management/internal/domain/admissions"
	"school-management/internal/domain/guardians"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
	"school-management/internal/domain/students"
	"school-management/internal/infrastructure/database/sqlite"
	"school-management/internal/infrastructure/storage"
	"school-management/internal/platform/config"
	"school-management/internal/platform/security"
)

func TestAcademicStudentAdmissionLifecycle(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sms-phase3-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := config.DefaultConfig(tempDir)
	cfg.Database.Path = filepath.Join(tempDir, "phase3_test.db")
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	ctx := context.Background()
	db, err := sqlite.Open(ctx, &cfg, logger)
	if err != nil {
		t.Fatalf("failed to open sqlite database: %v", err)
	}
	defer db.Close()

	blobDir := filepath.Join(tempDir, "blobs")
	blobStore, err := storage.NewFilesystemBlobStore(blobDir)
	if err != nil {
		t.Fatalf("failed to open blob store: %v", err)
	}

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

	// Handlers
	setupHandler := authCmd.NewFirstRunSetupHandler(uow, ids, clock, hasher)
	loginHandler := authCmd.NewLoginHandler(uow, sessionStore, ids, clock, hasher, db.DeviceID())
	authQueries := authQuery.NewAuthQueriesHandler(uow)

	acadCmdHandler := acadCmd.NewAcademicCommandsHandler(uow, clock, ids)
	acadQueries := acadQuery.NewAcademicQueriesHandler(uow)

	studCmdHandler := studCmd.NewStudentCommandsHandler(uow, clock, ids)
	studQueries := studQuery.NewStudentQueriesHandler(uow)

	guardCmdHandler := guardCmd.NewGuardianCommandsHandler(uow, clock, ids)
	guardQueries := guardQuery.NewGuardianQueriesHandler(uow)

	enrCmdHandler := enrCmd.NewEnrollmentCommandsHandler(uow, clock, ids)
	enrQueries := enrQuery.NewEnrollmentQueriesHandler(uow)

	admCmdHandler := admCmd.NewAdmissionCommandsHandler(uow, clock, ids)
	admQueries := admQuery.NewAdmissionQueriesHandler(uow)

	docCmdHandler := docCmd.NewDocumentCommandsHandler(uow, blobStore, clock, ids)
	docQueries := docQuery.NewDocumentQueriesHandler(uow, blobStore)

	_ = authQueries
	_ = acadQueries
	_ = studCmdHandler
	_ = guardQueries
	_ = admQueries

	// Step 1: Initial school setup
	_, err = setupHandler.Execute(ctx, authCmd.FirstRunSetupInput{
		SchoolName:         "Greenwood Academy",
		SchoolCode:         "GWA",
		Timezone:           "UTC",
		CurrencyCode:       "USD",
		CurrencyMinorUnits: 2,
		AdminUsername:      "superadmin",
		AdminDisplayName:   "System Administrator",
		AdminEmail:         "admin@greenwood.edu",
		AdminPassword:      "Password123!",
	})
	if err != nil {
		t.Fatalf("first run setup failed: %v", err)
	}

	// Login
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

	// Step 2: Create Academic Year 2026-2027
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

	// Set Terms
	err = acadCmdHandler.SetTerms(authCtx, acadCmd.SetTermsInput{
		AcademicYearID: yearID,
		Terms: []acadCmd.AcademicTermInput{
			{Sequence: 1, Name: "Fall Term", StartsOn: shared.CivilDate(2026, 9, 1), EndsOn: shared.CivilDate(2026, 12, 20)},
			{Sequence: 2, Name: "Spring Term", StartsOn: shared.CivilDate(2027, 1, 10), EndsOn: shared.CivilDate(2027, 6, 30)},
		},
	})
	if err != nil {
		t.Fatalf("failed to set terms: %v", err)
	}

	// Create Campus
	campusID, err := acadCmdHandler.CreateCampus(authCtx, acadCmd.CreateCampusInput{
		Name:    "Main Campus",
		Code:    "MAIN",
		Address: "123 School Way",
	})
	if err != nil {
		t.Fatalf("failed to create campus: %v", err)
	}

	// Create Class Grade 1
	class1ID, err := acadCmdHandler.CreateClass(authCtx, acadCmd.CreateClassInput{
		Name:       "Grade 1",
		Code:       "G1",
		LevelOrder: 1,
	})
	if err != nil {
		t.Fatalf("failed to create class: %v", err)
	}

	// Create Section Grade 1-A
	cap40 := 40
	section1ID, err := acadCmdHandler.CreateSection(authCtx, acadCmd.CreateSectionInput{
		ClassID:        class1ID,
		AcademicYearID: yearID,
		CampusID:       &campusID,
		Name:           "Section A",
		Capacity:       &cap40,
	})
	if err != nil {
		t.Fatalf("failed to create section: %v", err)
	}

	// Step 3: Configure Admission Form Engine
	formDefID, err := admCmdHandler.CreateFormDefinition(authCtx, admCmd.CreateFormDefInput{
		Code:        "STD-2026",
		Name:        "Standard Admission Form 2026",
		Description: "General admission form for Greenwood Academy",
	})
	if err != nil {
		t.Fatalf("failed to create admission form def: %v", err)
	}

	schema := admissions.FormSchema{
		Title:       "Greenwood Admission Form",
		Description: "Please complete all required fields",
		Sections: []admissions.FormSection{
			{
				ID:    "student_info",
				Title: "Student Information",
				Fields: []admissions.FormField{
					{
						ID:       "prev_school",
						Label:    "Previous School",
						Type:     admissions.FieldText,
						Required: false,
					},
					{
						ID:       "has_allergies",
						Label:    "Has Allergies",
						Type:     admissions.FieldDropdown,
						Required: true,
						Options:  []string{"Yes", "No"},
					},
					{
						ID:          "allergy_details",
						Label:       "Allergy Details",
						Type:        admissions.FieldText,
						Required:    true,
						DependsOn:   &admissions.Visibility{FieldID: "has_allergies", Equals: "Yes"},
						Placeholder: "Describe allergies and emergency treatments",
					},
				},
			},
		},
	}

	formVerID, err := admCmdHandler.CreateFormVersion(authCtx, admCmd.CreateFormVersionInput{
		FormDefinitionID: formDefID,
		Schema:           schema,
	})
	if err != nil {
		t.Fatalf("failed to create form version: %v", err)
	}

	// Publish Version
	err = admCmdHandler.PublishFormVersion(authCtx, admCmd.PublishFormVersionInput{
		VersionID: formVerID,
	})
	if err != nil {
		t.Fatalf("failed to publish form version: %v", err)
	}

	// Step 4: Submit Admission Application
	dob := shared.CivilDate(2020, 5, 15)
	appID, err := admCmdHandler.SubmitApplication(authCtx, admCmd.SubmitApplicationInput{
		FormVersionID:           formVerID,
		AcademicYearID:          yearID,
		AppliedClassID:          class1ID,
		ApplicantFirstName:      "Alice",
		ApplicantLastName:       "Wonderland",
		ApplicantDOB:            &dob,
		ApplicantGender:         "female",
		GuardianPhoneNormalized: "+1234567890",
		GuardianName:            "Charles Wonderland",
		GuardianEmail:           "charles@wonderland.org",
		Response: map[string]any{
			"has_allergies":   "Yes",
			"allergy_details": "Peanut allergy",
			"prev_school":     "Wonderland Nursery",
		},
	})
	if err != nil {
		t.Fatalf("failed to submit application: %v", err)
	}

	// Step 5: Review and Approve Application
	err = admCmdHandler.UpdateStatus(authCtx, admCmd.UpdateApplicationStatusInput{
		ApplicationID: appID,
		NewStatus:     admissions.AppApproved,
		Reason:        "Met admission interview criteria",
	})
	if err != nil {
		t.Fatalf("failed to approve application: %v", err)
	}

	// Step 6: Enroll Application
	studentID, err := admCmdHandler.EnrollApplication(authCtx, admCmd.EnrollApplicationInput{
		ApplicationID: appID,
		SectionID:     &section1ID,
		AdmissionNo:   "ADM-2026-0001",
		AdmissionDate: shared.CivilDate(2026, 9, 1),
	})
	if err != nil {
		t.Fatalf("failed to enroll application: %v", err)
	}

	// Verify student created and active
	student, err := studQueries.GetStudent(authCtx, studentID)
	if err != nil {
		t.Fatalf("failed to get student: %v", err)
	}
	if student.FirstName != "Alice" || student.LastName != "Wonderland" {
		t.Errorf("unexpected student name: %s %s", student.FirstName, student.LastName)
	}
	if student.Status != students.StatusActive {
		t.Errorf("expected student status active, got %s", student.Status)
	}

	// Verify active enrollment
	activeEnr, err := enrQueries.GetActiveEnrollment(authCtx, studentID)
	if err != nil {
		t.Fatalf("failed to get active enrollment: %v", err)
	}
	if activeEnr.ClassID != class1ID || *activeEnr.SectionID != section1ID {
		t.Errorf("unexpected class/section in enrollment")
	}

	// Step 7: Create Guardian and Link
	guardianID, err := guardCmdHandler.CreateGuardian(authCtx, guardCmd.CreateGuardianInput{
		FirstName:    "Charles",
		LastName:     "Wonderland",
		PhonePrimary: "+1234567890",
		Email:        "charles@wonderland.org",
	})
	if err != nil {
		t.Fatalf("failed to create guardian: %v", err)
	}

	err = guardCmdHandler.SetStudentGuardians(authCtx, guardCmd.SetStudentGuardiansInput{
		StudentID: studentID,
		Guardians: []guardCmd.StudentGuardianLinkInput{
			{
				GuardianID:         guardianID,
				Relationship:       guardians.RelFather,
				IsPrimary:          true,
				IsEmergencyContact: true,
				ReceivesBilling:    true,
				CanPickup:          true,
				Priority:           1,
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to link guardian: %v", err)
	}

	// Step 8: Full Text Search (FTS5) test
	searchRes, err := studQueries.ListStudents(authCtx, students.StudentFilter{
		SearchQuery: "Alice",
		Limit:       10,
	})
	if err != nil {
		t.Fatalf("failed to search students: %v", err)
	}
	if len(searchRes) == 0 {
		t.Fatalf("expected FTS search to find student Alice, got 0 results")
	}
	if searchRes[0].ID != studentID {
		t.Errorf("expected student ID %s, got %s", studentID, searchRes[0].ID)
	}

	// Step 9: Document upload to BlobStore
	docContent := []byte("%PDF-1.4 Mock Birth Certificate for Alice")
	docID, err := docCmdHandler.UploadDocument(authCtx, docCmd.UploadDocumentInput{
		StudentID:        &studentID,
		DocumentTypeCode: "BIRTH_CERT",
		FileName:         "alice_birth_cert.pdf",
		MimeType:         "application/pdf",
		Data:             docContent,
	})
	if err != nil {
		t.Fatalf("failed to upload document: %v", err)
	}

	readContent, docMeta, err := docQueries.GetDocumentContent(authCtx, docID)
	if err != nil {
		t.Fatalf("failed to get document content: %v", err)
	}
	if string(readContent) != string(docContent) {
		t.Errorf("retrieved document content mismatch")
	}
	if docMeta.FileName != "alice_birth_cert.pdf" {
		t.Errorf("unexpected document filename: %s", docMeta.FileName)
	}

	// Step 10: Academic Promotion to Next Year
	year2ID, err := acadCmdHandler.CreateAcademicYear(authCtx, acadCmd.CreateAcademicYearInput{
		Name:     "Academic Year 2027-2028",
		Code:     "2027-28",
		StartsOn: shared.CivilDate(2027, 9, 1),
		EndsOn:   shared.CivilDate(2028, 6, 30),
		Status:   academics.YearPlanned,
	})
	if err != nil {
		t.Fatalf("failed to create year 2: %v", err)
	}

	class2ID, err := acadCmdHandler.CreateClass(authCtx, acadCmd.CreateClassInput{
		Name:       "Grade 2",
		Code:       "G2",
		LevelOrder: 2,
	})
	if err != nil {
		t.Fatalf("failed to create class 2: %v", err)
	}

	section2ID, err := acadCmdHandler.CreateSection(authCtx, acadCmd.CreateSectionInput{
		ClassID:        class2ID,
		AcademicYearID: year2ID,
		CampusID:       &campusID,
		Name:           "Section A",
		Capacity:       &cap40,
	})
	if err != nil {
		t.Fatalf("failed to create section 2: %v", err)
	}

	rollNo := "01"
	batchID, err := enrCmdHandler.ExecutePromotionBatch(authCtx, enrCmd.ExecutePromotionBatchInput{
		FromAcademicYearID: yearID,
		ToAcademicYearID:   year2ID,
		EffectiveDate:      shared.CivilDate(2027, 9, 1),
		Notes:              "Annual promotion to Grade 2",
		Items: []enrCmd.PromotionItem{
			{
				StudentID:       studentID,
				TargetClassID:   class2ID,
				TargetSectionID: &section2ID,
				RollNo:          &rollNo,
				Action:          "promoted",
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to execute promotion batch: %v", err)
	}
	if batchID == "" {
		t.Fatalf("expected non-empty batch ID")
	}

	// Verify new active enrollment is in Grade 2
	newActiveEnr, err := enrQueries.GetActiveEnrollment(authCtx, studentID)
	if err != nil {
		t.Fatalf("failed to get active enrollment after promotion: %v", err)
	}
	if newActiveEnr.ClassID != class2ID || *newActiveEnr.SectionID != section2ID {
		t.Errorf("promoted enrollment does not point to Grade 2")
	}

	// Verify enrollment history has 2 records
	history, err := enrQueries.ListEnrollmentHistory(authCtx, studentID)
	if err != nil {
		t.Fatalf("failed to list enrollment history: %v", err)
	}
	if len(history) != 2 {
		t.Errorf("expected 2 enrollment records in history, got %d", len(history))
	}

	t.Logf("Phase 3 complete end-to-end integration test passed successfully! Admin=%s, Student=%s, Year=%s", loginRes.User.Username, studentID, yearID)
}
