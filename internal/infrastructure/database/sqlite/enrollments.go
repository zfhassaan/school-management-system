package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"school-management/internal/domain/shared"
	"school-management/internal/domain/students"
)

type EnrollmentsRepository struct {
	q Querier
}

func (r *EnrollmentsRepository) GetActive(ctx context.Context, studentID shared.ID) (*students.Enrollment, error) {
	row, err := r.q.GetActiveEnrollment(ctx, studentID.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("active enrollment")
		}
		return nil, shared.Internal(err)
	}

	enrolledOn, _ := shared.ParseDate(row.EnrolledOn)
	var endedOn *shared.Date
	if row.EndedOn.Valid {
		d, _ := shared.ParseDate(row.EndedOn.String)
		endedOn = &d
	}
	var secID *shared.ID
	if row.SectionID.Valid {
		id := shared.ID(row.SectionID.String)
		secID = &id
	}
	var rollNo *string
	if row.RollNo.Valid {
		rollNo = &row.RollNo.String
	}
	var prevID *shared.ID
	if row.PreviousEnrollmentID.Valid {
		id := shared.ID(row.PreviousEnrollmentID.String)
		prevID = &id
	}
	var batchID *shared.ID
	if row.PromotionBatchID.Valid {
		id := shared.ID(row.PromotionBatchID.String)
		batchID = &id
	}
	var appID *shared.ID
	if row.AdmissionApplicationID.Valid {
		id := shared.ID(row.AdmissionApplicationID.String)
		appID = &id
	}
	createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)

	return &students.Enrollment{
		ID:                     shared.ID(row.ID),
		StudentID:              shared.ID(row.StudentID),
		AcademicYearID:         shared.ID(row.AcademicYearID),
		ClassID:                shared.ID(row.ClassID),
		SectionID:              secID,
		RollNo:                 rollNo,
		Status:                 students.EnrollmentStatus(row.Status),
		EnrolledOn:             enrolledOn,
		EndedOn:                endedOn,
		PreviousEnrollmentID:   prevID,
		PromotionBatchID:       batchID,
		AdmissionApplicationID: appID,
		ClassName:              row.ClassName,
		SectionName:            row.SectionName,
		YearCode:               row.YearCode,
		CreatedAt:              createdAt,
		UpdatedAt:              updatedAt,
		Revision:               row.RowVersion,
	}, nil
}

func (r *EnrollmentsRepository) ListHistory(ctx context.Context, studentID shared.ID) ([]students.Enrollment, error) {
	rows, err := r.q.ListEnrollmentHistory(ctx, studentID.String())
	if err != nil {
		return nil, shared.Internal(err)
	}
	result := make([]students.Enrollment, len(rows))
	for i, row := range rows {
		enrolledOn, _ := shared.ParseDate(row.EnrolledOn)
		var endedOn *shared.Date
		if row.EndedOn.Valid {
			d, _ := shared.ParseDate(row.EndedOn.String)
			endedOn = &d
		}
		var secID *shared.ID
		if row.SectionID.Valid {
			id := shared.ID(row.SectionID.String)
			secID = &id
		}
		var rollNo *string
		if row.RollNo.Valid {
			rollNo = &row.RollNo.String
		}
		createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)

		result[i] = students.Enrollment{
			ID:             shared.ID(row.ID),
			StudentID:      shared.ID(row.StudentID),
			AcademicYearID: shared.ID(row.AcademicYearID),
			ClassID:        shared.ID(row.ClassID),
			SectionID:      secID,
			RollNo:         rollNo,
			Status:         students.EnrollmentStatus(row.Status),
			EnrolledOn:     enrolledOn,
			EndedOn:        endedOn,
			ClassName:      row.ClassName,
			SectionName:    row.SectionName,
			YearCode:       row.YearCode,
			CreatedAt:      createdAt,
			UpdatedAt:      updatedAt,
			Revision:       row.RowVersion,
		}
	}
	return result, nil
}

func (r *EnrollmentsRepository) Insert(ctx context.Context, e *students.Enrollment) error {
	now := e.CreatedAt.UTC().Format(time.RFC3339Nano)
	var secID sql.NullString
	if e.SectionID != nil {
		secID = sql.NullString{String: e.SectionID.String(), Valid: true}
	}
	var rollNo sql.NullString
	if e.RollNo != nil {
		rollNo = sql.NullString{String: *e.RollNo, Valid: true}
	}
	var endedOn sql.NullString
	if e.EndedOn != nil {
		endedOn = sql.NullString{String: e.EndedOn.String(), Valid: true}
	}
	var prevID sql.NullString
	if e.PreviousEnrollmentID != nil {
		prevID = sql.NullString{String: e.PreviousEnrollmentID.String(), Valid: true}
	}
	var batchID sql.NullString
	if e.PromotionBatchID != nil {
		batchID = sql.NullString{String: e.PromotionBatchID.String(), Valid: true}
	}
	var appID sql.NullString
	if e.AdmissionApplicationID != nil {
		appID = sql.NullString{String: e.AdmissionApplicationID.String(), Valid: true}
	}

	err := r.q.InsertEnrollment(ctx, InsertEnrollmentParams{
		ID:                     e.ID.String(),
		StudentID:              e.StudentID.String(),
		AcademicYearID:         e.AcademicYearID.String(),
		ClassID:                e.ClassID.String(),
		SectionID:              secID,
		RollNo:                 rollNo,
		Status:                 string(e.Status),
		EnrolledOn:             e.EnrolledOn.String(),
		EndedOn:                endedOn,
		PreviousEnrollmentID:   prevID,
		PromotionBatchID:       batchID,
		AdmissionApplicationID: appID,
		CreatedAt:              now,
		UpdatedAt:              now,
	})
	if err != nil {
		return shared.Internal(err)
	}
	e.Revision = 1
	return nil
}

func (r *EnrollmentsRepository) Close(ctx context.Context, id shared.ID, status students.EnrollmentStatus, endedOn shared.Date, revision int64) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	return r.q.CloseEnrollment(ctx, CloseEnrollmentParams{
		Status:     string(status),
		EndedOn:    sql.NullString{String: endedOn.String(), Valid: true},
		UpdatedAt:  now,
		ID:         id.String(),
		RowVersion: revision,
	})
}

func (r *EnrollmentsRepository) ListSectionRoster(ctx context.Context, sectionID shared.ID) ([]students.Enrollment, error) {
	rows, err := r.q.ListSectionRoster(ctx, sql.NullString{String: sectionID.String(), Valid: true})
	if err != nil {
		return nil, shared.Internal(err)
	}
	result := make([]students.Enrollment, len(rows))
	for i, row := range rows {
		var rollNo *string
		if row.RollNo.Valid {
			rollNo = &row.RollNo.String
		}
		result[i] = students.Enrollment{
			ID:        shared.ID(row.EnrollmentID),
			StudentID: shared.ID(row.ID),
			RollNo:    rollNo,
			Status:    students.EnrollmentStatus(row.EnrollmentStatus),
		}
	}
	return result, nil
}

func (r *EnrollmentsRepository) InsertPromotionBatch(ctx context.Context, batch *students.PromotionBatch) error {
	now := batch.CreatedAt.UTC().Format(time.RFC3339Nano)
	var execBy sql.NullString
	if batch.ExecutedBy != nil {
		execBy = sql.NullString{String: batch.ExecutedBy.String(), Valid: true}
	}
	var execAt sql.NullString
	if batch.ExecutedAt != nil {
		execAt = sql.NullString{String: batch.ExecutedAt.UTC().Format(time.RFC3339Nano), Valid: true}
	}

	err := r.q.InsertPromotionBatch(ctx, InsertPromotionBatchParams{
		ID:                 batch.ID.String(),
		FromAcademicYearID: batch.FromAcademicYearID.String(),
		ToAcademicYearID:   batch.ToAcademicYearID.String(),
		Status:             string(batch.Status),
		SummaryJson:        batch.SummaryJSON,
		ExecutedBy:         execBy,
		ExecutedAt:         execAt,
		CreatedAt:          now,
		UpdatedAt:          now,
	})
	if err != nil {
		return shared.Internal(err)
	}
	batch.Revision = 1
	return nil
}

func (r *EnrollmentsRepository) GetPromotionBatch(ctx context.Context, id shared.ID) (*students.PromotionBatch, error) {
	row, err := r.q.GetPromotionBatch(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("promotion batch")
		}
		return nil, shared.Internal(err)
	}
	var execBy *shared.ID
	if row.ExecutedBy.Valid {
		id := shared.ID(row.ExecutedBy.String)
		execBy = &id
	}
	var execAt *time.Time
	if row.ExecutedAt.Valid {
		t, _ := time.Parse(time.RFC3339Nano, row.ExecutedAt.String)
		execAt = &t
	}
	createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
	return &students.PromotionBatch{
		ID:                 shared.ID(row.ID),
		FromAcademicYearID: shared.ID(row.FromAcademicYearID),
		ToAcademicYearID:   shared.ID(row.ToAcademicYearID),
		Status:             students.PromotionBatchStatus(row.Status),
		SummaryJSON:        row.SummaryJson,
		ExecutedBy:         execBy,
		ExecutedAt:         execAt,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
		Revision:           row.RowVersion,
	}, nil
}
