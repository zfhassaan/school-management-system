package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"school-management/internal/domain/shared"
	"school-management/internal/domain/students"
)

type StudentsRepository struct {
	q Querier
}

func (r *StudentsRepository) Get(ctx context.Context, id shared.ID) (*students.Student, error) {
	row, err := r.q.GetStudent(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("student")
		}
		return nil, shared.Internal(err)
	}
	return toDomainStudent(row), nil
}

func (r *StudentsRepository) GetByAdmissionNo(ctx context.Context, schoolID shared.ID, no string) (*students.Student, error) {
	row, err := r.q.GetStudentByAdmissionNo(ctx, GetStudentByAdmissionNoParams{
		SchoolID:    schoolID.String(),
		AdmissionNo: no,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("student")
		}
		return nil, shared.Internal(err)
	}
	return toDomainStudent(row), nil
}

func (r *StudentsRepository) Insert(ctx context.Context, s *students.Student) error {
	now := s.CreatedAt.UTC().Format(time.RFC3339Nano)
	var photoID sql.NullString
	if s.PhotoDocID != nil {
		photoID = sql.NullString{String: s.PhotoDocID.String(), Valid: true}
	}

	err := r.q.InsertStudent(ctx, InsertStudentParams{
		ID:              s.ID.String(),
		SchoolID:        s.SchoolID.String(),
		AdmissionNo:     s.AdmissionNo,
		FirstName:       s.FirstName,
		MiddleName:      sql.NullString{String: s.MiddleName, Valid: s.MiddleName != ""},
		LastName:        s.LastName,
		SearchName:      s.SearchName,
		DateOfBirth:     s.DateOfBirth.String(),
		Gender:          string(s.Gender),
		NationalID:      sql.NullString{String: s.NationalID, Valid: s.NationalID != ""},
		AddressLine1:    sql.NullString{String: s.Address.Line1, Valid: s.Address.Line1 != ""},
		AddressLine2:    sql.NullString{String: s.Address.Line2, Valid: s.Address.Line2 != ""},
		City:            sql.NullString{String: s.Address.City, Valid: s.Address.City != ""},
		Region:          sql.NullString{String: s.Address.Region, Valid: s.Address.Region != ""},
		PostalCode:      sql.NullString{String: s.Address.PostalCode, Valid: s.Address.PostalCode != ""},
		Country:         sql.NullString{String: s.Address.Country, Valid: s.Address.Country != ""},
		PhotoDocumentID: photoID,
		AdmissionDate:   s.AdmissionDate.String(),
		Status:          string(s.Status),
		Notes:           sql.NullString{String: s.Notes, Valid: s.Notes != ""},
		CreatedAt:       now,
		UpdatedAt:       now,
	})
	if err != nil {
		return shared.Internal(err)
	}
	s.Revision = 1
	return nil
}

func (r *StudentsRepository) Update(ctx context.Context, s *students.Student) error {
	now := s.UpdatedAt.UTC().Format(time.RFC3339Nano)
	var photoID sql.NullString
	if s.PhotoDocID != nil {
		photoID = sql.NullString{String: s.PhotoDocID.String(), Valid: true}
	}

	err := r.q.UpdateStudent(ctx, UpdateStudentParams{
		FirstName:       s.FirstName,
		MiddleName:      sql.NullString{String: s.MiddleName, Valid: s.MiddleName != ""},
		LastName:        s.LastName,
		SearchName:      s.SearchName,
		DateOfBirth:     s.DateOfBirth.String(),
		Gender:          string(s.Gender),
		NationalID:      sql.NullString{String: s.NationalID, Valid: s.NationalID != ""},
		AddressLine1:    sql.NullString{String: s.Address.Line1, Valid: s.Address.Line1 != ""},
		AddressLine2:    sql.NullString{String: s.Address.Line2, Valid: s.Address.Line2 != ""},
		City:            sql.NullString{String: s.Address.City, Valid: s.Address.City != ""},
		Region:          sql.NullString{String: s.Address.Region, Valid: s.Address.Region != ""},
		PostalCode:      sql.NullString{String: s.Address.PostalCode, Valid: s.Address.PostalCode != ""},
		Country:         sql.NullString{String: s.Address.Country, Valid: s.Address.Country != ""},
		PhotoDocumentID: photoID,
		Status:          string(s.Status),
		Notes:           sql.NullString{String: s.Notes, Valid: s.Notes != ""},
		UpdatedAt:       now,
		ID:              s.ID.String(),
		RowVersion:      s.Revision,
	})
	if err != nil {
		return shared.Internal(err)
	}
	s.Revision++
	return nil
}

func (r *StudentsRepository) AppendStatusHistory(ctx context.Context, e students.StatusHistoryEntry) error {
	return r.q.AppendStudentStatusHistory(ctx, AppendStudentStatusHistoryParams{
		ID:          e.ID.String(),
		StudentID:   e.StudentID.String(),
		FromStatus:  sql.NullString{String: e.FromStatus, Valid: e.FromStatus != ""},
		ToStatus:    e.ToStatus,
		EffectiveOn: e.EffectiveOn.String(),
		Reason:      e.Reason,
		ChangedBy:   e.ChangedBy.String(),
		ChangedAt:   e.ChangedAt.UTC().Format(time.RFC3339Nano),
	})
}

func (r *StudentsRepository) ListStatusHistory(ctx context.Context, studentID shared.ID) ([]students.StatusHistoryEntry, error) {
	rows, err := r.q.ListStudentStatusHistory(ctx, studentID.String())
	if err != nil {
		return nil, shared.Internal(err)
	}
	entries := make([]students.StatusHistoryEntry, len(rows))
	for i, row := range rows {
		effOn, _ := shared.ParseDate(row.EffectiveOn)
		changedAt, _ := time.Parse(time.RFC3339Nano, row.ChangedAt)
		entries[i] = students.StatusHistoryEntry{
			ID:          shared.ID(row.ID),
			StudentID:   shared.ID(row.StudentID),
			FromStatus:  row.FromStatus.String,
			ToStatus:    row.ToStatus,
			EffectiveOn: effOn,
			Reason:      row.Reason,
			ChangedBy:   shared.ID(row.ChangedBy),
			ChangedAt:   changedAt,
		}
	}
	return entries, nil
}

func (r *StudentsRepository) ListCursor(ctx context.Context, f students.StudentFilter) ([]students.Student, error) {
	limit := int64(f.Limit)
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	rows, err := r.q.ListStudentsCursor(ctx, ListStudentsCursorParams{
		SchoolID:     f.SchoolID.String(),
		Column2:      f.Status,
		Status:       f.Status,
		Column4:      f.ClassID,
		ClassID:      f.ClassID,
		Column6:      f.SectionID,
		SectionID:    sql.NullString{String: f.SectionID, Valid: f.SectionID != ""},
		SearchName:   f.CursorName,
		SearchName_2: f.CursorName,
		ID:           f.CursorID,
		Limit:        limit,
	})
	if err != nil {
		return nil, shared.Internal(err)
	}

	result := make([]students.Student, len(rows))
	for i, row := range rows {
		dob, _ := shared.ParseDate(row.DateOfBirth)
		admDate, _ := shared.ParseDate(row.AdmissionDate)
		createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
		var photoID *shared.ID
		if row.PhotoDocumentID.Valid {
			id := shared.ID(row.PhotoDocumentID.String)
			photoID = &id
		}
		result[i] = students.Student{
			ID:            shared.ID(row.ID),
			SchoolID:      shared.ID(row.SchoolID),
			AdmissionNo:   row.AdmissionNo,
			FirstName:     row.FirstName,
			MiddleName:    row.MiddleName.String,
			LastName:      row.LastName,
			SearchName:    row.SearchName,
			DateOfBirth:   dob,
			Gender:        students.Gender(row.Gender),
			NationalID:    row.NationalID.String,
			Address: shared.Address{
				Line1:      row.AddressLine1.String,
				Line2:      row.AddressLine2.String,
				City:       row.City.String,
				Region:     row.Region.String,
				PostalCode: row.PostalCode.String,
				Country:    row.Country.String,
			},
			PhotoDocID:    photoID,
			AdmissionDate: admDate,
			Status:        students.Status(row.Status),
			Notes:         row.Notes.String,
			CreatedAt:     createdAt,
			UpdatedAt:     updatedAt,
			Revision:      row.RowVersion,
		}
	}
	return result, nil
}

func (r *StudentsRepository) SearchFTS(ctx context.Context, schoolID shared.ID, query string, limit int) ([]students.Student, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	rows, err := r.q.SearchStudentsFTS(ctx, SearchStudentsFTSParams{
		SearchName: query,
		SchoolID:   schoolID.String(),
		Limit:      int64(limit),
	})
	if err != nil {
		return nil, shared.Internal(err)
	}
	result := make([]students.Student, len(rows))
	for i, row := range rows {
		dob, _ := shared.ParseDate(row.DateOfBirth)
		admDate, _ := shared.ParseDate(row.AdmissionDate)
		createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
		var photoID *shared.ID
		if row.PhotoDocumentID.Valid {
			id := shared.ID(row.PhotoDocumentID.String)
			photoID = &id
		}
		result[i] = students.Student{
			ID:            shared.ID(row.ID),
			SchoolID:      shared.ID(row.SchoolID),
			AdmissionNo:   row.AdmissionNo,
			FirstName:     row.FirstName,
			MiddleName:    row.MiddleName.String,
			LastName:      row.LastName,
			SearchName:    row.SearchName,
			DateOfBirth:   dob,
			Gender:        students.Gender(row.Gender),
			NationalID:    row.NationalID.String,
			Address: shared.Address{
				Line1:      row.AddressLine1.String,
				Line2:      row.AddressLine2.String,
				City:       row.City.String,
				Region:     row.Region.String,
				PostalCode: row.PostalCode.String,
				Country:    row.Country.String,
			},
			PhotoDocID:    photoID,
			AdmissionDate: admDate,
			Status:        students.Status(row.Status),
			Notes:         row.Notes.String,
			CreatedAt:     createdAt,
			UpdatedAt:     updatedAt,
			Revision:      row.RowVersion,
		}
	}
	return result, nil
}

func (r *StudentsRepository) Count(ctx context.Context, schoolID shared.ID) (int64, error) {
	count, err := r.q.CountStudents(ctx, schoolID.String())
	if err != nil {
		return 0, shared.Internal(err)
	}
	return count, nil
}

func (r *StudentsRepository) FindSiblings(ctx context.Context, studentID shared.ID) ([]students.Student, error) {
	rows, err := r.q.FindSiblingsByGuardian(ctx, FindSiblingsByGuardianParams{
		StudentID: studentID.String(),
		ID:        studentID.String(),
	})
	if err != nil {
		return nil, shared.Internal(err)
	}
	result := make([]students.Student, len(rows))
	for i, row := range rows {
		dob, _ := shared.ParseDate(row.DateOfBirth)
		admDate, _ := shared.ParseDate(row.AdmissionDate)
		createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
		var photoID *shared.ID
		if row.PhotoDocumentID.Valid {
			id := shared.ID(row.PhotoDocumentID.String)
			photoID = &id
		}
		result[i] = students.Student{
			ID:            shared.ID(row.ID),
			SchoolID:      shared.ID(row.SchoolID),
			AdmissionNo:   row.AdmissionNo,
			FirstName:     row.FirstName,
			MiddleName:    row.MiddleName.String,
			LastName:      row.LastName,
			SearchName:    row.SearchName,
			DateOfBirth:   dob,
			Gender:        students.Gender(row.Gender),
			NationalID:    row.NationalID.String,
			Address: shared.Address{
				Line1:      row.AddressLine1.String,
				Line2:      row.AddressLine2.String,
				City:       row.City.String,
				Region:     row.Region.String,
				PostalCode: row.PostalCode.String,
				Country:    row.Country.String,
			},
			PhotoDocID:    photoID,
			AdmissionDate: admDate,
			Status:        students.Status(row.Status),
			Notes:         row.Notes.String,
			CreatedAt:     createdAt,
			UpdatedAt:     updatedAt,
			Revision:      row.RowVersion,
		}
	}
	return result, nil
}

func toDomainStudent(row Student) *students.Student {
	dob, _ := shared.ParseDate(row.DateOfBirth)
	admDate, _ := shared.ParseDate(row.AdmissionDate)
	createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
	var photoID *shared.ID
	if row.PhotoDocumentID.Valid {
		id := shared.ID(row.PhotoDocumentID.String)
		photoID = &id
	}
	return &students.Student{
		ID:            shared.ID(row.ID),
		SchoolID:      shared.ID(row.SchoolID),
		AdmissionNo:   row.AdmissionNo,
		FirstName:     row.FirstName,
		MiddleName:    row.MiddleName.String,
		LastName:      row.LastName,
		SearchName:    row.SearchName,
		DateOfBirth:   dob,
		Gender:        students.Gender(row.Gender),
		NationalID:    row.NationalID.String,
		Address: shared.Address{
			Line1:      row.AddressLine1.String,
			Line2:      row.AddressLine2.String,
			City:       row.City.String,
			Region:     row.Region.String,
			PostalCode: row.PostalCode.String,
			Country:    row.Country.String,
		},
		PhotoDocID:    photoID,
		AdmissionDate: admDate,
		Status:        students.Status(row.Status),
		Notes:         row.Notes.String,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		Revision:      row.RowVersion,
	}
}
