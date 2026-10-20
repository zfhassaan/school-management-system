package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"school-management/internal/domain/academics"
	"school-management/internal/domain/shared"
)

type AcademicsRepository struct {
	q Querier
}

func (r *AcademicsRepository) ListAcademicYears(ctx context.Context, schoolID shared.ID) ([]academics.AcademicYear, error) {
	rows, err := r.q.ListAcademicYears(ctx, schoolID.String())
	if err != nil {
		return nil, shared.Internal(err)
	}
	years := make([]academics.AcademicYear, len(rows))
	for i, row := range rows {
		startsOn, _ := shared.ParseDate(row.StartsOn)
		endsOn, _ := shared.ParseDate(row.EndsOn)
		createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
		years[i] = academics.AcademicYear{
			ID:        shared.ID(row.ID),
			SchoolID:  shared.ID(row.SchoolID),
			Code:      row.Code,
			Name:      row.Name,
			StartsOn:  startsOn,
			EndsOn:    endsOn,
			Status:    academics.YearStatus(row.Status),
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			Revision:  row.RowVersion,
		}
	}
	return years, nil
}

func (r *AcademicsRepository) GetAcademicYear(ctx context.Context, id shared.ID) (*academics.AcademicYear, error) {
	row, err := r.q.GetAcademicYear(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("academic year")
		}
		return nil, shared.Internal(err)
	}
	startsOn, _ := shared.ParseDate(row.StartsOn)
	endsOn, _ := shared.ParseDate(row.EndsOn)
	createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
	return &academics.AcademicYear{
		ID:        shared.ID(row.ID),
		SchoolID:  shared.ID(row.SchoolID),
		Code:      row.Code,
		Name:      row.Name,
		StartsOn:  startsOn,
		EndsOn:    endsOn,
		Status:    academics.YearStatus(row.Status),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Revision:  row.RowVersion,
	}, nil
}

func (r *AcademicsRepository) GetActiveAcademicYear(ctx context.Context, schoolID shared.ID) (*academics.AcademicYear, error) {
	row, err := r.q.GetActiveAcademicYear(ctx, schoolID.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("active academic year")
		}
		return nil, shared.Internal(err)
	}
	startsOn, _ := shared.ParseDate(row.StartsOn)
	endsOn, _ := shared.ParseDate(row.EndsOn)
	createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
	return &academics.AcademicYear{
		ID:        shared.ID(row.ID),
		SchoolID:  shared.ID(row.SchoolID),
		Code:      row.Code,
		Name:      row.Name,
		StartsOn:  startsOn,
		EndsOn:    endsOn,
		Status:    academics.YearStatus(row.Status),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Revision:  row.RowVersion,
	}, nil
}

func (r *AcademicsRepository) InsertAcademicYear(ctx context.Context, y *academics.AcademicYear) error {
	now := y.CreatedAt.UTC().Format(time.RFC3339Nano)
	err := r.q.InsertAcademicYear(ctx, InsertAcademicYearParams{
		ID:        y.ID.String(),
		SchoolID:  y.SchoolID.String(),
		Code:      y.Code,
		Name:      y.Name,
		StartsOn:  y.StartsOn.String(),
		EndsOn:    y.EndsOn.String(),
		Status:    string(y.Status),
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return shared.Internal(err)
	}
	y.Revision = 1
	return nil
}

func (r *AcademicsRepository) UpdateAcademicYear(ctx context.Context, y *academics.AcademicYear) error {
	now := y.UpdatedAt.UTC().Format(time.RFC3339Nano)
	err := r.q.UpdateAcademicYear(ctx, UpdateAcademicYearParams{
		Name:       y.Name,
		StartsOn:   y.StartsOn.String(),
		EndsOn:     y.EndsOn.String(),
		Status:     string(y.Status),
		UpdatedAt:  now,
		ID:         y.ID.String(),
		RowVersion: y.Revision,
	})
	if err != nil {
		return shared.Internal(err)
	}
	y.Revision++
	return nil
}

func (r *AcademicsRepository) ListAcademicTerms(ctx context.Context, yearID shared.ID) ([]academics.AcademicTerm, error) {
	rows, err := r.q.ListAcademicTerms(ctx, yearID.String())
	if err != nil {
		return nil, shared.Internal(err)
	}
	terms := make([]academics.AcademicTerm, len(rows))
	for i, row := range rows {
		startsOn, _ := shared.ParseDate(row.StartsOn)
		endsOn, _ := shared.ParseDate(row.EndsOn)
		terms[i] = academics.AcademicTerm{
			ID:             shared.ID(row.ID),
			AcademicYearID: shared.ID(row.AcademicYearID),
			Sequence:       int(row.Sequence),
			Name:           row.Name,
			StartsOn:       startsOn,
			EndsOn:         endsOn,
		}
	}
	return terms, nil
}

func (r *AcademicsRepository) ReplaceAcademicTerms(ctx context.Context, yearID shared.ID, terms []academics.AcademicTerm) error {
	if err := r.q.DeleteAcademicTermsByYear(ctx, yearID.String()); err != nil {
		return shared.Internal(err)
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	for _, t := range terms {
		err := r.q.InsertAcademicTerm(ctx, InsertAcademicTermParams{
			ID:             t.ID.String(),
			AcademicYearID: yearID.String(),
			Sequence:       int64(t.Sequence),
			Name:           t.Name,
			StartsOn:       t.StartsOn.String(),
			EndsOn:         t.EndsOn.String(),
			CreatedAt:      now,
			UpdatedAt:      now,
		})
		if err != nil {
			return shared.Internal(err)
		}
	}
	return nil
}

func (r *AcademicsRepository) ListClasses(ctx context.Context, schoolID shared.ID) ([]academics.Class, error) {
	rows, err := r.q.ListClasses(ctx, schoolID.String())
	if err != nil {
		return nil, shared.Internal(err)
	}
	classes := make([]academics.Class, len(rows))
	for i, row := range rows {
		var nextID *shared.ID
		if row.NextClassID.Valid {
			id := shared.ID(row.NextClassID.String)
			nextID = &id
		}
		createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
		classes[i] = academics.Class{
			ID:          shared.ID(row.ID),
			SchoolID:    shared.ID(row.SchoolID),
			Code:        row.Code,
			Name:        row.Name,
			LevelOrder:  int(row.LevelOrder),
			NextClassID: nextID,
			IsActive:    row.IsActive == 1,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			Revision:    row.RowVersion,
		}
	}
	return classes, nil
}

func (r *AcademicsRepository) GetClass(ctx context.Context, id shared.ID) (*academics.Class, error) {
	row, err := r.q.GetClass(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("class")
		}
		return nil, shared.Internal(err)
	}
	var nextID *shared.ID
	if row.NextClassID.Valid {
		id := shared.ID(row.NextClassID.String)
		nextID = &id
	}
	createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
	return &academics.Class{
		ID:          shared.ID(row.ID),
		SchoolID:    shared.ID(row.SchoolID),
		Code:        row.Code,
		Name:        row.Name,
		LevelOrder:  int(row.LevelOrder),
		NextClassID: nextID,
		IsActive:    row.IsActive == 1,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Revision:    row.RowVersion,
	}, nil
}

func (r *AcademicsRepository) InsertClass(ctx context.Context, c *academics.Class) error {
	now := c.CreatedAt.UTC().Format(time.RFC3339Nano)
	var nextID sql.NullString
	if c.NextClassID != nil {
		nextID = sql.NullString{String: c.NextClassID.String(), Valid: true}
	}
	isActive := int64(0)
	if c.IsActive {
		isActive = 1
	}
	err := r.q.InsertClass(ctx, InsertClassParams{
		ID:          c.ID.String(),
		SchoolID:    c.SchoolID.String(),
		Code:        c.Code,
		Name:        c.Name,
		LevelOrder:  int64(c.LevelOrder),
		NextClassID: nextID,
		IsActive:    isActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		return shared.Internal(err)
	}
	c.Revision = 1
	return nil
}

func (r *AcademicsRepository) UpdateClass(ctx context.Context, c *academics.Class) error {
	now := c.UpdatedAt.UTC().Format(time.RFC3339Nano)
	var nextID sql.NullString
	if c.NextClassID != nil {
		nextID = sql.NullString{String: c.NextClassID.String(), Valid: true}
	}
	isActive := int64(0)
	if c.IsActive {
		isActive = 1
	}
	err := r.q.UpdateClass(ctx, UpdateClassParams{
		Code:        c.Code,
		Name:        c.Name,
		LevelOrder:  int64(c.LevelOrder),
		NextClassID: nextID,
		IsActive:    isActive,
		UpdatedAt:   now,
		ID:         c.ID.String(),
		RowVersion: c.Revision,
	})
	if err != nil {
		return shared.Internal(err)
	}
	c.Revision++
	return nil
}

func (r *AcademicsRepository) ListSectionsByYear(ctx context.Context, yearID shared.ID) ([]academics.Section, error) {
	rows, err := r.q.ListSectionsByYear(ctx, yearID.String())
	if err != nil {
		return nil, shared.Internal(err)
	}
	sections := make([]academics.Section, len(rows))
	for i, row := range rows {
		var campusID *shared.ID
		if row.CampusID.Valid {
			id := shared.ID(row.CampusID.String)
			campusID = &id
		}
		var teacherID *shared.ID
		if row.ClassTeacherUserID.Valid {
			id := shared.ID(row.ClassTeacherUserID.String)
			teacherID = &id
		}
		var capacity *int
		if row.Capacity.Valid {
			cap := int(row.Capacity.Int64)
			capacity = &cap
		}
		createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
		sections[i] = academics.Section{
			ID:                 shared.ID(row.ID),
			ClassID:            shared.ID(row.ClassID),
			AcademicYearID:     shared.ID(row.AcademicYearID),
			CampusID:           campusID,
			Name:               row.Name,
			Capacity:           capacity,
			ClassTeacherUserID: teacherID,
			IsActive:           row.IsActive == 1,
			CreatedAt:          createdAt,
			UpdatedAt:          updatedAt,
			Revision:           row.RowVersion,
		}
	}
	return sections, nil
}

func (r *AcademicsRepository) GetSection(ctx context.Context, id shared.ID) (*academics.Section, error) {
	row, err := r.q.GetSection(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("section")
		}
		return nil, shared.Internal(err)
	}
	var campusID *shared.ID
	if row.CampusID.Valid {
		id := shared.ID(row.CampusID.String)
		campusID = &id
	}
	var teacherID *shared.ID
	if row.ClassTeacherUserID.Valid {
		id := shared.ID(row.ClassTeacherUserID.String)
		teacherID = &id
	}
	var capacity *int
	if row.Capacity.Valid {
		cap := int(row.Capacity.Int64)
		capacity = &cap
	}
	createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
	return &academics.Section{
		ID:                 shared.ID(row.ID),
		ClassID:            shared.ID(row.ClassID),
		AcademicYearID:     shared.ID(row.AcademicYearID),
		CampusID:           campusID,
		Name:               row.Name,
		Capacity:           capacity,
		ClassTeacherUserID: teacherID,
		IsActive:           row.IsActive == 1,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
		Revision:           row.RowVersion,
	}, nil
}

func (r *AcademicsRepository) InsertSection(ctx context.Context, s *academics.Section) error {
	now := s.CreatedAt.UTC().Format(time.RFC3339Nano)
	var campusID sql.NullString
	if s.CampusID != nil {
		campusID = sql.NullString{String: s.CampusID.String(), Valid: true}
	}
	var teacherID sql.NullString
	if s.ClassTeacherUserID != nil {
		teacherID = sql.NullString{String: s.ClassTeacherUserID.String(), Valid: true}
	}
	var capacity sql.NullInt64
	if s.Capacity != nil {
		capacity = sql.NullInt64{Int64: int64(*s.Capacity), Valid: true}
	}
	isActive := int64(0)
	if s.IsActive {
		isActive = 1
	}
	err := r.q.InsertSection(ctx, InsertSectionParams{
		ID:                 s.ID.String(),
		ClassID:            s.ClassID.String(),
		AcademicYearID:     s.AcademicYearID.String(),
		CampusID:           campusID,
		Name:               s.Name,
		Capacity:           capacity,
		ClassTeacherUserID: teacherID,
		IsActive:           isActive,
		CreatedAt:          now,
		UpdatedAt:          now,
	})
	if err != nil {
		return shared.Internal(err)
	}
	s.Revision = 1
	return nil
}

func (r *AcademicsRepository) UpdateSection(ctx context.Context, s *academics.Section) error {
	now := s.UpdatedAt.UTC().Format(time.RFC3339Nano)
	var teacherID sql.NullString
	if s.ClassTeacherUserID != nil {
		teacherID = sql.NullString{String: s.ClassTeacherUserID.String(), Valid: true}
	}
	var capacity sql.NullInt64
	if s.Capacity != nil {
		capacity = sql.NullInt64{Int64: int64(*s.Capacity), Valid: true}
	}
	isActive := int64(0)
	if s.IsActive {
		isActive = 1
	}
	err := r.q.UpdateSection(ctx, UpdateSectionParams{
		Name:               s.Name,
		Capacity:           capacity,
		ClassTeacherUserID: teacherID,
		IsActive:           isActive,
		UpdatedAt:          now,
		ID:                 s.ID.String(),
		RowVersion:         s.Revision,
	})
	if err != nil {
		return shared.Internal(err)
	}
	s.Revision++
	return nil
}

func (r *AcademicsRepository) ListCampuses(ctx context.Context, schoolID shared.ID) ([]academics.Campus, error) {
	rows, err := r.q.ListCampuses(ctx, schoolID.String())
	if err != nil {
		return nil, shared.Internal(err)
	}
	campuses := make([]academics.Campus, len(rows))
	for i, row := range rows {
		createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
		campuses[i] = academics.Campus{
			ID:        shared.ID(row.ID),
			SchoolID:  shared.ID(row.SchoolID),
			Code:      row.Code,
			Name:      row.Name,
			Address:   row.Address.String,
			IsActive:  row.IsActive == 1,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			Revision:  row.RowVersion,
		}
	}
	return campuses, nil
}

func (r *AcademicsRepository) InsertCampus(ctx context.Context, c *academics.Campus) error {
	now := c.CreatedAt.UTC().Format(time.RFC3339Nano)
	isActive := int64(0)
	if c.IsActive {
		isActive = 1
	}
	err := r.q.InsertCampus(ctx, InsertCampusParams{
		ID:        c.ID.String(),
		SchoolID:  c.SchoolID.String(),
		Code:      c.Code,
		Name:      c.Name,
		Address:   sql.NullString{String: c.Address, Valid: c.Address != ""},
		IsActive:  isActive,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return shared.Internal(err)
	}
	c.Revision = 1
	return nil
}
