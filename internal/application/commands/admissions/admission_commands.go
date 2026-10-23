package admissions

import (
	"context"
	"fmt"
	"strings"

	"school-management/internal/application/services/auditlog"
	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/uow"
	"school-management/internal/domain/admissions"
	"school-management/internal/domain/guardians"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
	"school-management/internal/domain/students"
)

type CreateFormDefInput struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateFormVersionInput struct {
	FormDefinitionID shared.ID             `json:"formDefinitionId"`
	Schema           admissions.FormSchema `json:"schema"`
}

type PublishFormVersionInput struct {
	VersionID shared.ID `json:"versionId"`
}

type SubmitApplicationInput struct {
	FormVersionID           shared.ID              `json:"formVersionId"`
	AcademicYearID          shared.ID              `json:"academicYearId"`
	AppliedClassID          shared.ID              `json:"appliedClassId"`
	ApplicantFirstName      string                 `json:"applicantFirstName"`
	ApplicantLastName       string                 `json:"applicantLastName"`
	ApplicantDOB            *shared.Date           `json:"applicantDob,omitempty"`
	ApplicantGender         string                 `json:"applicantGender"`
	GuardianPhoneNormalized string                 `json:"guardianPhoneNormalized"`
	GuardianName            string                 `json:"guardianName"`
	GuardianEmail           string                 `json:"guardianEmail"`
	Response                map[string]any         `json:"response"`
}

type UpdateApplicationStatusInput struct {
	ApplicationID shared.ID                    `json:"applicationId"`
	NewStatus     admissions.ApplicationStatus `json:"newStatus"`
	Reason        string                       `json:"reason"`
}

type EnrollApplicationInput struct {
	ApplicationID shared.ID   `json:"applicationId"`
	ClassID       *shared.ID  `json:"classId,omitempty"`
	SectionID     *shared.ID  `json:"sectionId,omitempty"`
	RollNo        *string     `json:"rollNo,omitempty"`
	AdmissionNo   string      `json:"admissionNo"`
	AdmissionDate shared.Date `json:"admissionDate"`
}

type AdmissionCommandsHandler struct {
	uow      uow.UnitOfWork
	clock    shared.Clock
	ids      shared.IDGenerator
	recorder auditlog.Recorder
}

func NewAdmissionCommandsHandler(
	uow uow.UnitOfWork,
	clock shared.Clock,
	ids shared.IDGenerator,
) *AdmissionCommandsHandler {
	return &AdmissionCommandsHandler{
		uow:      uow,
		clock:    clock,
		ids:      ids,
		recorder: auditlog.Recorder{IDs: ids, Clock: clock},
	}
}

func (h *AdmissionCommandsHandler) CreateFormDefinition(ctx context.Context, in CreateFormDefInput) (shared.ID, error) {
	if _, err := authz.Require(ctx, identity.PermAdmissionFormsMgt); err != nil {
		return "", err
	}
	id := h.ids.NewID()
	now := h.clock.Now()

	err := h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}
		fd := &admissions.FormDefinition{
			ID:          id,
			SchoolID:    sch.ID,
			Code:        strings.TrimSpace(in.Code),
			Name:        strings.TrimSpace(in.Name),
			Description: strings.TrimSpace(in.Description),
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if err := r.Admissions.InsertFormDefinition(ctx, fd); err != nil {
			return err
		}
		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "admissions.form_def_created",
			EntityType: "admission_form_definition",
			EntityID:   id.String(),
			Reason:     "form definition created",
			After:      map[string]any{"name": fd.Name, "code": fd.Code},
		})
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func (h *AdmissionCommandsHandler) CreateFormVersion(ctx context.Context, in CreateFormVersionInput) (shared.ID, error) {
	if _, err := authz.Require(ctx, identity.PermAdmissionFormsMgt); err != nil {
		return "", err
	}
	id := h.ids.NewID()
	now := h.clock.Now()

	hash, err := admissions.ComputeSchemaHash(in.Schema)
	if err != nil {
		return "", shared.Validation("invalid schema structure", nil)
	}

	err = h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		versions, err := r.Admissions.ListFormVersions(ctx, in.FormDefinitionID)
		if err != nil {
			return err
		}
		nextVerNo := len(versions) + 1

		v := &admissions.FormVersion{
			ID:               id,
			FormDefinitionID: in.FormDefinitionID,
			VersionNo:        nextVerNo,
			Status:           admissions.FormDraft,
			Schema:           in.Schema,
			SchemaHash:       hash,
			CreatedAt:        now,
			UpdatedAt:        now,
		}

		if err := r.Admissions.InsertFormVersion(ctx, v); err != nil {
			return err
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "admissions.form_version_created",
			EntityType: "admission_form_version",
			EntityID:   id.String(),
			Reason:     "new draft form version created",
			After:      map[string]any{"versionNo": nextVerNo, "schemaHash": hash},
		})
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func (h *AdmissionCommandsHandler) PublishFormVersion(ctx context.Context, in PublishFormVersionInput) error {
	actor, err := authz.Require(ctx, identity.PermAdmissionFormsMgt)
	if err != nil {
		return err
	}
	now := h.clock.Now()

	return h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		v, err := r.Admissions.GetFormVersion(ctx, in.VersionID)
		if err != nil {
			return err
		}
		hash, err := admissions.ComputeSchemaHash(v.Schema)
		if err != nil {
			return err
		}

		if err := r.Admissions.PublishFormVersion(ctx, in.VersionID, actor.UserID, now, hash); err != nil {
			return err
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "admissions.form_version_published",
			EntityType: "admission_form_version",
			EntityID:   in.VersionID.String(),
			Reason:     "form version published",
			After:      map[string]any{"versionNo": v.VersionNo, "schemaHash": hash},
		})
	})
}

// EnsureDefaultAdmissionForm verifies that a published admission form exists.
// If none exists, it automatically creates and publishes the default institutional form.
func (h *AdmissionCommandsHandler) EnsureDefaultAdmissionForm(ctx context.Context) (shared.ID, error) {
	var actorID shared.ID
	if a, err := authz.ActorFrom(ctx); err == nil {
		actorID = a.UserID
	}
	now := h.clock.Now()
	var publishedVerID shared.ID

	err := h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}

		defs, err := r.Admissions.ListFormDefinitions(ctx, sch.ID)
		if err == nil {
			for _, def := range defs {
				pub, err := r.Admissions.GetPublishedVersion(ctx, def.ID)
				if err == nil && pub != nil && pub.Status == admissions.FormPublished {
					publishedVerID = pub.ID
					return nil
				}
			}
		}

		var targetDefID shared.ID
		var targetDefName string

		if len(defs) > 0 {
			targetDefID = defs[0].ID
			targetDefName = defs[0].Name
		} else {
			defID := h.ids.NewID()
			def := &admissions.FormDefinition{
				ID:          defID,
				SchoolID:    sch.ID,
				Code:        "STD",
				Name:        "General Admission Form",
				Description: "Standard admission application form",
				IsActive:    true,
				CreatedAt:   now,
				UpdatedAt:   now,
			}
			if err := r.Admissions.InsertFormDefinition(ctx, def); err != nil {
				return err
			}
			targetDefID = defID
			targetDefName = def.Name
		}

		versions, _ := r.Admissions.ListFormVersions(ctx, targetDefID)
		if len(versions) > 0 {
			latest := versions[len(versions)-1]
			hash, _ := admissions.ComputeSchemaHash(latest.Schema)
			_ = r.Admissions.PublishFormVersion(ctx, latest.ID, actorID, now, hash)
			publishedVerID = latest.ID
			return nil
		}

		verID := h.ids.NewID()
		schema := admissions.DefaultFormSchema()
		schema.Title = targetDefName
		hash, _ := admissions.ComputeSchemaHash(schema)
		ver := &admissions.FormVersion{
			ID:               verID,
			FormDefinitionID: targetDefID,
			VersionNo:        1,
			Status:           admissions.FormDraft,
			Schema:           schema,
			SchemaHash:       hash,
			CreatedAt:        now,
			UpdatedAt:        now,
		}
		if err := r.Admissions.InsertFormVersion(ctx, ver); err != nil {
			return err
		}
		if err := r.Admissions.PublishFormVersion(ctx, verID, actorID, now, hash); err != nil {
			return err
		}
		publishedVerID = verID
		return nil
	})

	if err != nil {
		return "", err
	}
	return publishedVerID, nil
}

func (h *AdmissionCommandsHandler) SubmitApplication(ctx context.Context, in SubmitApplicationInput) (shared.ID, error) {
	actor, err := authz.Require(ctx, identity.PermAdmissionsCreate)
	if err != nil {
		return "", err
	}
	id := h.ids.NewID()
	now := h.clock.Now()

	err = h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}

		formVerID := in.FormVersionID
		if formVerID == "" {
			defs, err := r.Admissions.ListFormDefinitions(ctx, sch.ID)
			if err == nil {
				for _, def := range defs {
					pub, err := r.Admissions.GetPublishedVersion(ctx, def.ID)
					if err == nil && pub != nil && pub.Status == admissions.FormPublished {
						formVerID = pub.ID
						break
					}
				}
			}

			if formVerID == "" {
				var targetDefID shared.ID
				var targetDefName string

				if len(defs) > 0 {
					targetDefID = defs[0].ID
					targetDefName = defs[0].Name
				} else {
					defID := h.ids.NewID()
					def := &admissions.FormDefinition{
						ID:          defID,
						SchoolID:    sch.ID,
						Code:        "STD",
						Name:        "General Admission Form",
						Description: "Standard admission application form",
						IsActive:    true,
						CreatedAt:   now,
						UpdatedAt:   now,
					}
					if err := r.Admissions.InsertFormDefinition(ctx, def); err != nil {
						return err
					}
					targetDefID = defID
					targetDefName = def.Name
				}

				// Check if there are existing versions for this definition
				versions, _ := r.Admissions.ListFormVersions(ctx, targetDefID)
				if len(versions) > 0 {
					latest := versions[len(versions)-1]
					hash, _ := admissions.ComputeSchemaHash(latest.Schema)
					_ = r.Admissions.PublishFormVersion(ctx, latest.ID, actor.UserID, now, hash)
					formVerID = latest.ID
				} else {
					verID := h.ids.NewID()
					schema := admissions.FormSchema{
						Title:    targetDefName,
						Sections: []admissions.FormSection{},
					}
					hash, _ := admissions.ComputeSchemaHash(schema)
					ver := &admissions.FormVersion{
						ID:               verID,
						FormDefinitionID: targetDefID,
						VersionNo:        1,
						Status:           admissions.FormDraft,
						Schema:           schema,
						SchemaHash:       hash,
						CreatedAt:        now,
						UpdatedAt:        now,
					}
					if err := r.Admissions.InsertFormVersion(ctx, ver); err != nil {
						return err
					}
					if err := r.Admissions.PublishFormVersion(ctx, verID, actor.UserID, now, hash); err != nil {
						return err
					}
					formVerID = verID
				}
			}
		}

		// Validate against published version schema
		v, err := r.Admissions.GetFormVersion(ctx, formVerID)
		if err != nil {
			return err
		}
		if v.Status == admissions.FormDraft {
			hash, _ := admissions.ComputeSchemaHash(v.Schema)
			if pubErr := r.Admissions.PublishFormVersion(ctx, formVerID, actor.UserID, now, hash); pubErr == nil {
				v.Status = admissions.FormPublished
			}
		}
		if v.Status != admissions.FormPublished {
			return shared.Validation("can only submit against published form versions", nil)
		}
		if in.Response != nil {
			if err := admissions.ValidateResponse(v.Schema, in.Response); err != nil {
				return err
			}
		}

		idSuffix := strings.ToUpper(strings.ReplaceAll(id.String(), "-", ""))
		if len(idSuffix) > 4 {
			idSuffix = idSuffix[len(idSuffix)-4:]
		}
		appNo := fmt.Sprintf("APP-%d-%04d-%s", now.Year(), (now.UnixNano()/1e6)%10000, idSuffix)

		app := &admissions.Application{
			ID:                      id,
			SchoolID:                sch.ID,
			ApplicationNo:           appNo,
			FormVersionID:           formVerID,
			AcademicYearID:          in.AcademicYearID,
			AppliedClassID:          in.AppliedClassID,
			Status:                  admissions.AppSubmitted,
			ApplicantFirstName:      strings.TrimSpace(in.ApplicantFirstName),
			ApplicantLastName:       strings.TrimSpace(in.ApplicantLastName),
			ApplicantDOB:            in.ApplicantDOB,
			ApplicantGender:         strings.TrimSpace(in.ApplicantGender),
			GuardianPhoneNormalized: guardians.NormalizePhone(in.GuardianPhoneNormalized),
			Response:                in.Response,
			SubmittedAt:             &now,
			CreatedBy:               actor.UserID,
			CreatedAt:               now,
			UpdatedAt:               now,
		}

		if err := r.Admissions.InsertApplication(ctx, app); err != nil {
			return err
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "admissions.application_submitted",
			EntityType: "admission_application",
			EntityID:   id.String(),
			Reason:     "application submitted",
			After:      map[string]any{"applicationNo": appNo, "applicant": app.ApplicantFirstName + " " + app.ApplicantLastName},
		})
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func (h *AdmissionCommandsHandler) UpdateStatus(ctx context.Context, in UpdateApplicationStatusInput) error {
	actor, err := authz.Require(ctx, identity.PermAdmissionsReview)
	if err != nil {
		return err
	}
	now := h.clock.Now()

	return h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		app, err := r.Admissions.GetApplication(ctx, in.ApplicationID)
		if err != nil {
			return err
		}
		oldStatus := app.Status
		app.Status = in.NewStatus
		if in.NewStatus == admissions.AppApproved || in.NewStatus == admissions.AppRejected {
			app.DecidedAt = &now
			app.DecidedBy = &actor.UserID
			app.DecisionReason = strings.TrimSpace(in.Reason)
		}

		if err := r.Admissions.UpdateStatus(ctx, app, in.Reason, actor.UserID, now); err != nil {
			return err
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "admissions.status_updated",
			EntityType: "admission_application",
			EntityID:   in.ApplicationID.String(),
			Reason:     in.Reason,
			Before:     map[string]any{"status": oldStatus},
			After:      map[string]any{"status": in.NewStatus},
		})
	})
}

func (h *AdmissionCommandsHandler) EnrollApplication(ctx context.Context, in EnrollApplicationInput) (shared.ID, error) {
	actor, err := authz.Require(ctx, identity.PermAdmissionsEnroll)
	if err != nil {
		return "", err
	}
	studentID := h.ids.NewID()
	enrollmentID := h.ids.NewID()
	now := h.clock.Now()

	err = h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		app, err := r.Admissions.GetApplication(ctx, in.ApplicationID)
		if err != nil {
			return err
		}
		if app.Status != admissions.AppApproved {
			return shared.Validation("only approved applications can be enrolled", nil)
		}

		dob := shared.Date{}
		if app.ApplicantDOB != nil {
			dob = *app.ApplicantDOB
		}

		gender := students.Gender(app.ApplicantGender)
		if gender == "" {
			gender = students.GenderUnspecified
		}

		targetClassID := app.AppliedClassID
		if in.ClassID != nil && !in.ClassID.IsZero() {
			targetClassID = *in.ClassID
		}

		admNo := strings.TrimSpace(in.AdmissionNo)
		if admNo == "" {
			admNo = strings.Replace(app.ApplicationNo, "APP-", "ADM-", 1)
			if admNo == app.ApplicationNo {
				admNo = "ADM-" + app.ApplicationNo
			}
		}

		student := &students.Student{
			ID:            studentID,
			SchoolID:      app.SchoolID,
			AdmissionNo:   admNo,
			FirstName:     app.ApplicantFirstName,
			LastName:      app.ApplicantLastName,
			DateOfBirth:   dob,
			Gender:        gender,
			Status:        students.StatusActive,
			AdmissionDate: in.AdmissionDate,
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		if err := student.Validate(); err != nil {
			return err
		}

		if err := r.Students.Insert(ctx, student); err != nil {
			return fmt.Errorf("failed to insert student: %w", err)
		}

		appID := app.ID
		enr := &students.Enrollment{
			ID:                     enrollmentID,
			StudentID:              studentID,
			AcademicYearID:         app.AcademicYearID,
			ClassID:                targetClassID,
			SectionID:              in.SectionID,
			RollNo:                 in.RollNo,
			Status:                 students.EnrollmentActive,
			EnrolledOn:             in.AdmissionDate,
			AdmissionApplicationID: &appID,
			CreatedAt:              now,
			UpdatedAt:              now,
		}
		if err := r.Enrollments.Insert(ctx, enr); err != nil {
			return fmt.Errorf("failed to insert enrollment: %w", err)
		}

		// Update application to enrolled
		app.Status = admissions.AppEnrolled
		app.StudentID = &studentID
		app.AppliedClassID = targetClassID
		if err := r.Admissions.UpdateStatus(ctx, app, "Enrolled as student", actor.UserID, now); err != nil {
			return err
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "admissions.application_enrolled",
			EntityType: "admission_application",
			EntityID:   in.ApplicationID.String(),
			Reason:     "application enrolled as student",
			After: map[string]any{
				"studentId":   studentID.String(),
				"admissionNo": student.AdmissionNo,
			},
		})
	})
	if err != nil {
		return "", err
	}
	return studentID, nil
}

