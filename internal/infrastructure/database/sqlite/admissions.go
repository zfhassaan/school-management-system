package sqlite

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"school-management/internal/domain/admissions"
	"school-management/internal/domain/shared"
)

type AdmissionsRepository struct {
	q Querier
}

func (r *AdmissionsRepository) ListFormDefinitions(ctx context.Context, schoolID shared.ID) ([]admissions.FormDefinition, error) {
	rows, err := r.q.ListAdmissionFormDefinitions(ctx, schoolID.String())
	if err != nil {
		return nil, shared.Internal(err)
	}
	result := make([]admissions.FormDefinition, len(rows))
	for i, row := range rows {
		createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
		result[i] = admissions.FormDefinition{
			ID:          shared.ID(row.ID),
			SchoolID:    shared.ID(row.SchoolID),
			Code:        row.Code,
			Name:        row.Name,
			Description: row.Description.String,
			IsActive:    row.IsActive == 1,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			Revision:    row.RowVersion,
		}
	}
	return result, nil
}

func (r *AdmissionsRepository) GetFormDefinition(ctx context.Context, id shared.ID) (*admissions.FormDefinition, error) {
	row, err := r.q.GetAdmissionFormDefinition(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("admission form definition")
		}
		return nil, shared.Internal(err)
	}
	createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
	return &admissions.FormDefinition{
		ID:          shared.ID(row.ID),
		SchoolID:    shared.ID(row.SchoolID),
		Code:        row.Code,
		Name:        row.Name,
		Description: row.Description.String,
		IsActive:    row.IsActive == 1,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Revision:    row.RowVersion,
	}, nil
}

func (r *AdmissionsRepository) InsertFormDefinition(ctx context.Context, d *admissions.FormDefinition) error {
	now := d.CreatedAt.UTC().Format(time.RFC3339Nano)
	err := r.q.InsertAdmissionFormDefinition(ctx, InsertAdmissionFormDefinitionParams{
		ID:          d.ID.String(),
		SchoolID:    d.SchoolID.String(),
		Code:        d.Code,
		Name:        d.Name,
		Description: sql.NullString{String: d.Description, Valid: d.Description != ""},
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		return shared.Internal(err)
	}
	d.Revision = 1
	return nil
}

func (r *AdmissionsRepository) ListFormVersions(ctx context.Context, formDefID shared.ID) ([]admissions.FormVersion, error) {
	rows, err := r.q.ListAdmissionFormVersions(ctx, formDefID.String())
	if err != nil {
		return nil, shared.Internal(err)
	}
	result := make([]admissions.FormVersion, len(rows))
	for i, row := range rows {
		var schema admissions.FormSchema
		_ = json.Unmarshal([]byte(row.SchemaJson), &schema)
		var pubAt *time.Time
		if row.PublishedAt.Valid {
			t, _ := time.Parse(time.RFC3339Nano, row.PublishedAt.String)
			pubAt = &t
		}
		var pubBy *shared.ID
		if row.PublishedBy.Valid {
			id := shared.ID(row.PublishedBy.String)
			pubBy = &id
		}
		createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)

		result[i] = admissions.FormVersion{
			ID:               shared.ID(row.ID),
			FormDefinitionID: shared.ID(row.FormDefinitionID),
			VersionNo:        int(row.VersionNo),
			Status:           admissions.FormStatus(row.Status),
			Schema:           schema,
			SchemaHash:       row.SchemaHash,
			PublishedAt:      pubAt,
			PublishedBy:      pubBy,
			CreatedAt:        createdAt,
			UpdatedAt:        updatedAt,
			Revision:         row.RowVersion,
		}
	}
	return result, nil
}

func (r *AdmissionsRepository) GetFormVersion(ctx context.Context, id shared.ID) (*admissions.FormVersion, error) {
	row, err := r.q.GetAdmissionFormVersion(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("admission form version")
		}
		return nil, shared.Internal(err)
	}
	var schema admissions.FormSchema
	_ = json.Unmarshal([]byte(row.SchemaJson), &schema)
	var pubAt *time.Time
	if row.PublishedAt.Valid {
		t, _ := time.Parse(time.RFC3339Nano, row.PublishedAt.String)
		pubAt = &t
	}
	var pubBy *shared.ID
	if row.PublishedBy.Valid {
		id := shared.ID(row.PublishedBy.String)
		pubBy = &id
	}
	createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)

	return &admissions.FormVersion{
		ID:               shared.ID(row.ID),
		FormDefinitionID: shared.ID(row.FormDefinitionID),
		VersionNo:        int(row.VersionNo),
		Status:           admissions.FormStatus(row.Status),
		Schema:           schema,
		SchemaHash:       row.SchemaHash,
		PublishedAt:      pubAt,
		PublishedBy:      pubBy,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
		Revision:         row.RowVersion,
	}, nil
}

func (r *AdmissionsRepository) GetPublishedVersion(ctx context.Context, formDefID shared.ID) (*admissions.FormVersion, error) {
	row, err := r.q.GetPublishedFormVersion(ctx, formDefID.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("published form version")
		}
		return nil, shared.Internal(err)
	}
	var schema admissions.FormSchema
	_ = json.Unmarshal([]byte(row.SchemaJson), &schema)
	var pubAt *time.Time
	if row.PublishedAt.Valid {
		t, _ := time.Parse(time.RFC3339Nano, row.PublishedAt.String)
		pubAt = &t
	}
	var pubBy *shared.ID
	if row.PublishedBy.Valid {
		id := shared.ID(row.PublishedBy.String)
		pubBy = &id
	}
	createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)

	return &admissions.FormVersion{
		ID:               shared.ID(row.ID),
		FormDefinitionID: shared.ID(row.FormDefinitionID),
		VersionNo:        int(row.VersionNo),
		Status:           admissions.FormStatus(row.Status),
		Schema:           schema,
		SchemaHash:       row.SchemaHash,
		PublishedAt:      pubAt,
		PublishedBy:      pubBy,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
		Revision:         row.RowVersion,
	}, nil
}

func (r *AdmissionsRepository) InsertFormVersion(ctx context.Context, v *admissions.FormVersion) error {
	schemaJSON, err := json.Marshal(v.Schema)
	if err != nil {
		return shared.Internal(err)
	}
	now := v.CreatedAt.UTC().Format(time.RFC3339Nano)
	err = r.q.InsertAdmissionFormVersion(ctx, InsertAdmissionFormVersionParams{
		ID:               v.ID.String(),
		FormDefinitionID: v.FormDefinitionID.String(),
		VersionNo:        int64(v.VersionNo),
		SchemaJson:       string(schemaJSON),
		SchemaHash:       v.SchemaHash,
		CreatedAt:        now,
		UpdatedAt:        now,
	})
	if err != nil {
		return shared.Internal(err)
	}
	v.Revision = 1
	return nil
}

func (r *AdmissionsRepository) PublishFormVersion(ctx context.Context, id shared.ID, publishedBy shared.ID, at time.Time, hash string) error {
	now := at.UTC().Format(time.RFC3339Nano)
	return r.q.PublishAdmissionFormVersion(ctx, PublishAdmissionFormVersionParams{
		PublishedAt: sql.NullString{String: now, Valid: true},
		PublishedBy: sql.NullString{String: publishedBy.String(), Valid: publishedBy != ""},
		SchemaHash:  hash,
		UpdatedAt:   now,
		ID:          id.String(),
	})
}

func (r *AdmissionsRepository) ListApplications(ctx context.Context, schoolID shared.ID, status string, yearID string, limit int) ([]admissions.Application, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	rows, err := r.q.ListAdmissionApplications(ctx, ListAdmissionApplicationsParams{
		SchoolID:       schoolID.String(),
		Column2:        status,
		Status:         status,
		Column4:        yearID,
		AcademicYearID: yearID,
		Limit:          int64(limit),
	})
	if err != nil {
		return nil, shared.Internal(err)
	}
	result := make([]admissions.Application, len(rows))
	for i, row := range rows {
		var dob *shared.Date
		if row.ApplicantDob.Valid {
			d, _ := shared.ParseDate(row.ApplicantDob.String)
			dob = &d
		}
		var studentID *shared.ID
		if row.StudentID.Valid {
			id := shared.ID(row.StudentID.String)
			studentID = &id
		}
		var subAt *time.Time
		if row.SubmittedAt.Valid {
			t, _ := time.Parse(time.RFC3339Nano, row.SubmittedAt.String)
			subAt = &t
		}
		createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)

		result[i] = admissions.Application{
			ID:                      shared.ID(row.ID),
			SchoolID:                shared.ID(row.SchoolID),
			ApplicationNo:           row.ApplicationNo,
			FormVersionID:           shared.ID(row.FormVersionID),
			AcademicYearID:          shared.ID(row.AcademicYearID),
			AppliedClassID:          shared.ID(row.AppliedClassID),
			Status:                  admissions.ApplicationStatus(row.Status),
			ApplicantFirstName:      row.ApplicantFirstName,
			ApplicantLastName:       row.ApplicantLastName,
			ApplicantDOB:            dob,
			ApplicantGender:         row.ApplicantGender.String,
			GuardianPhoneNormalized: row.GuardianPhoneNormalized.String,
			StudentID:               studentID,
			SubmittedAt:             subAt,
			CreatedBy:               shared.ID(row.CreatedBy),
			ClassName:               row.ClassName,
			YearCode:                row.YearCode,
			CreatedAt:               createdAt,
			UpdatedAt:               updatedAt,
			Revision:                row.RowVersion,
		}
	}
	return result, nil
}

func (r *AdmissionsRepository) GetApplication(ctx context.Context, id shared.ID) (*admissions.Application, error) {
	row, err := r.q.GetAdmissionApplication(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("admission application")
		}
		return nil, shared.Internal(err)
	}
	var dob *shared.Date
	if row.ApplicantDob.Valid {
		d, _ := shared.ParseDate(row.ApplicantDob.String)
		dob = &d
	}
	var studentID *shared.ID
	if row.StudentID.Valid {
		id := shared.ID(row.StudentID.String)
		studentID = &id
	}
	var subAt *time.Time
	if row.SubmittedAt.Valid {
		t, _ := time.Parse(time.RFC3339Nano, row.SubmittedAt.String)
		subAt = &t
	}
	var resp map[string]any
	if row.ResponseJson.Valid {
		_ = json.Unmarshal([]byte(row.ResponseJson.String), &resp)
	}
	createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)

	return &admissions.Application{
		ID:                      shared.ID(row.ID),
		SchoolID:                shared.ID(row.SchoolID),
		ApplicationNo:           row.ApplicationNo,
		FormVersionID:           shared.ID(row.FormVersionID),
		AcademicYearID:          shared.ID(row.AcademicYearID),
		AppliedClassID:          shared.ID(row.AppliedClassID),
		Status:                  admissions.ApplicationStatus(row.Status),
		ApplicantFirstName:      row.ApplicantFirstName,
		ApplicantLastName:       row.ApplicantLastName,
		ApplicantDOB:            dob,
		ApplicantGender:         row.ApplicantGender.String,
		GuardianPhoneNormalized: row.GuardianPhoneNormalized.String,
		StudentID:               studentID,
		Response:                resp,
		SubmittedAt:             subAt,
		CreatedBy:               shared.ID(row.CreatedBy),
		ClassName:               row.ClassName,
		YearCode:                row.YearCode,
		CreatedAt:               createdAt,
		UpdatedAt:               updatedAt,
		Revision:                row.RowVersion,
	}, nil
}

func (r *AdmissionsRepository) InsertApplication(ctx context.Context, a *admissions.Application) error {
	now := a.CreatedAt.UTC().Format(time.RFC3339Nano)
	var dob sql.NullString
	if a.ApplicantDOB != nil {
		dob = sql.NullString{String: a.ApplicantDOB.String(), Valid: true}
	}
	var subAt sql.NullString
	if a.SubmittedAt != nil {
		subAt = sql.NullString{String: a.SubmittedAt.UTC().Format(time.RFC3339Nano), Valid: true}
	}
	var studentID sql.NullString
	if a.StudentID != nil {
		studentID = sql.NullString{String: a.StudentID.String(), Valid: true}
	}

	err := r.q.InsertAdmissionApplication(ctx, InsertAdmissionApplicationParams{
		ID:                      a.ID.String(),
		SchoolID:                a.SchoolID.String(),
		ApplicationNo:           a.ApplicationNo,
		FormVersionID:           a.FormVersionID.String(),
		AcademicYearID:          a.AcademicYearID.String(),
		AppliedClassID:          a.AppliedClassID.String(),
		Status:                  string(a.Status),
		ApplicantFirstName:      a.ApplicantFirstName,
		ApplicantLastName:       a.ApplicantLastName,
		ApplicantDob:            dob,
		ApplicantGender:         sql.NullString{String: a.ApplicantGender, Valid: a.ApplicantGender != ""},
		GuardianPhoneNormalized: sql.NullString{String: a.GuardianPhoneNormalized, Valid: a.GuardianPhoneNormalized != ""},
		StudentID:               studentID,
		SubmittedAt:             subAt,
		CreatedBy:               a.CreatedBy.String(),
		CreatedAt:               now,
		UpdatedAt:               now,
	})
	if err != nil {
		return shared.Internal(err)
	}

	if a.Response != nil {
		respJSON, err := json.Marshal(a.Response)
		if err == nil {
			hash := sha256.Sum256(respJSON)
			_ = r.q.InsertAdmissionFormResponse(ctx, InsertAdmissionFormResponseParams{
				ID:            shared.UUIDv7Generator{}.NewID().String(),
				ApplicationID: a.ID.String(),
				FormVersionID: a.FormVersionID.String(),
				ResponseJson:  string(respJSON),
				ResponseHash:  hex.EncodeToString(hash[:]),
				CreatedAt:     now,
				UpdatedAt:     now,
			})
		}
	}

	a.Revision = 1
	return nil
}

func (r *AdmissionsRepository) UpdateStatus(ctx context.Context, a *admissions.Application, reason string, actorID shared.ID, at time.Time) error {
	now := at.UTC().Format(time.RFC3339Nano)
	var studentID sql.NullString
	if a.StudentID != nil {
		studentID = sql.NullString{String: a.StudentID.String(), Valid: true}
	}

	err := r.q.UpdateAdmissionApplicationStatus(ctx, UpdateAdmissionApplicationStatusParams{
		Status:         string(a.Status),
		DecidedAt:      sql.NullString{String: now, Valid: true},
		DecidedBy:      sql.NullString{String: actorID.String(), Valid: true},
		DecisionReason: sql.NullString{String: reason, Valid: reason != ""},
		StudentID:      studentID,
		UpdatedAt:      now,
		ID:             a.ID.String(),
		RowVersion:     a.Revision,
	})
	if err != nil {
		return shared.Internal(err)
	}

	_ = r.q.AppendAdmissionStatusHistory(ctx, AppendAdmissionStatusHistoryParams{
		ID:            shared.UUIDv7Generator{}.NewID().String(),
		ApplicationID: a.ID.String(),
		FromStatus:    sql.NullString{},
		ToStatus:      string(a.Status),
		Reason:        sql.NullString{String: reason, Valid: reason != ""},
		ChangedBy:     actorID.String(),
		ChangedAt:     now,
	})

	a.Revision++
	return nil
}
