package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"school-management/internal/domain/fees"
	"school-management/internal/domain/shared"
)

type FeesRepository struct {
	q Querier
}

// --- Fee Heads ---

func (r *FeesRepository) ListFeeHeads(ctx context.Context, schoolID shared.ID) ([]fees.FeeHead, error) {
	rows, err := r.q.ListFeeHeads(ctx, schoolID.String())
	if err != nil {
		return nil, shared.Internal(err)
	}
	result := make([]fees.FeeHead, len(rows))
	for i, row := range rows {
		createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
		result[i] = fees.FeeHead{
			ID:                 shared.ID(row.ID),
			SchoolID:           shared.ID(row.SchoolID),
			Code:               row.Code,
			Name:               row.Name,
			Description:        row.Description,
			IsRecurring:        row.IsRecurring == 1,
			DefaultAmountMinor: row.DefaultAmountMinor,
			IsActive:           row.IsActive == 1,
			CreatedAt:          createdAt,
			UpdatedAt:          updatedAt,
			Revision:           row.RowVersion,
		}
	}
	return result, nil
}

func (r *FeesRepository) GetFeeHead(ctx context.Context, id shared.ID) (*fees.FeeHead, error) {
	row, err := r.q.GetFeeHead(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("fee head")
		}
		return nil, shared.Internal(err)
	}
	createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
	return &fees.FeeHead{
		ID:                 shared.ID(row.ID),
		SchoolID:           shared.ID(row.SchoolID),
		Code:               row.Code,
		Name:               row.Name,
		Description:        row.Description,
		IsRecurring:        row.IsRecurring == 1,
		DefaultAmountMinor: row.DefaultAmountMinor,
		IsActive:           row.IsActive == 1,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
		Revision:           row.RowVersion,
	}, nil
}

func (r *FeesRepository) InsertFeeHead(ctx context.Context, h *fees.FeeHead) error {
	now := h.CreatedAt.UTC().Format(time.RFC3339Nano)
	var isRec, isAct int64
	if h.IsRecurring {
		isRec = 1
	}
	if h.IsActive {
		isAct = 1
	}
	err := r.q.InsertFeeHead(ctx, InsertFeeHeadParams{
		ID:                 h.ID.String(),
		SchoolID:           h.SchoolID.String(),
		Code:               h.Code,
		Name:               h.Name,
		Description:        h.Description,
		IsRecurring:        isRec,
		DefaultAmountMinor: h.DefaultAmountMinor,
		IsActive:           isAct,
		CreatedAt:          now,
		UpdatedAt:          now,
	})
	if err != nil {
		return shared.Internal(err)
	}
	return nil
}

func (r *FeesRepository) UpdateFeeHead(ctx context.Context, h *fees.FeeHead) error {
	now := h.UpdatedAt.UTC().Format(time.RFC3339Nano)
	var isRec, isAct int64
	if h.IsRecurring {
		isRec = 1
	}
	if h.IsActive {
		isAct = 1
	}
	err := r.q.UpdateFeeHead(ctx, UpdateFeeHeadParams{
		Name:               h.Name,
		Description:        h.Description,
		IsRecurring:        isRec,
		DefaultAmountMinor: h.DefaultAmountMinor,
		IsActive:           isAct,
		UpdatedAt:          now,
		ID:                 h.ID.String(),
	})
	if err != nil {
		return shared.Internal(err)
	}
	return nil
}

// --- Fee Structures ---

func (r *FeesRepository) ListFeeStructures(ctx context.Context, schoolID shared.ID, academicYearID shared.ID) ([]fees.FeeStructure, error) {
	rows, err := r.q.ListFeeStructures(ctx, ListFeeStructuresParams{
		SchoolID:       schoolID.String(),
		AcademicYearID: academicYearID.String(),
	})
	if err != nil {
		return nil, shared.Internal(err)
	}
	result := make([]fees.FeeStructure, len(rows))
	for i, row := range rows {
		var campusID *shared.ID
		if row.CampusID.Valid {
			id := shared.ID(row.CampusID.String)
			campusID = &id
		}
		createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
		result[i] = fees.FeeStructure{
			ID:             shared.ID(row.ID),
			SchoolID:       shared.ID(row.SchoolID),
			AcademicYearID: shared.ID(row.AcademicYearID),
			ClassID:        shared.ID(row.ClassID),
			CampusID:       campusID,
			Code:           row.Code,
			Name:           row.Name,
			Description:    row.Description,
			IsActive:       row.IsActive == 1,
			CreatedAt:      createdAt,
			UpdatedAt:      updatedAt,
			Revision:       row.RowVersion,
		}
	}
	return result, nil
}

func (r *FeesRepository) GetFeeStructure(ctx context.Context, id shared.ID) (*fees.FeeStructure, error) {
	row, err := r.q.GetFeeStructure(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("fee structure")
		}
		return nil, shared.Internal(err)
	}
	var campusID *shared.ID
	if row.CampusID.Valid {
		cid := shared.ID(row.CampusID.String)
		campusID = &cid
	}
	createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
	return &fees.FeeStructure{
		ID:             shared.ID(row.ID),
		SchoolID:       shared.ID(row.SchoolID),
		AcademicYearID: shared.ID(row.AcademicYearID),
		ClassID:        shared.ID(row.ClassID),
		CampusID:       campusID,
		Code:           row.Code,
		Name:           row.Name,
		Description:    row.Description,
		IsActive:       row.IsActive == 1,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		Revision:       row.RowVersion,
	}, nil
}

func (r *FeesRepository) GetFeeStructureForClass(ctx context.Context, academicYearID shared.ID, classID shared.ID) (*fees.FeeStructure, error) {
	row, err := r.q.GetFeeStructureForClass(ctx, GetFeeStructureForClassParams{
		AcademicYearID: academicYearID.String(),
		ClassID:        classID.String(),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("fee structure for class")
		}
		return nil, shared.Internal(err)
	}
	var campusID *shared.ID
	if row.CampusID.Valid {
		cid := shared.ID(row.CampusID.String)
		campusID = &cid
	}
	createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
	return &fees.FeeStructure{
		ID:             shared.ID(row.ID),
		SchoolID:       shared.ID(row.SchoolID),
		AcademicYearID: shared.ID(row.AcademicYearID),
		ClassID:        shared.ID(row.ClassID),
		CampusID:       campusID,
		Code:           row.Code,
		Name:           row.Name,
		Description:    row.Description,
		IsActive:       row.IsActive == 1,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		Revision:       row.RowVersion,
	}, nil
}

func (r *FeesRepository) InsertFeeStructure(ctx context.Context, s *fees.FeeStructure) error {
	now := s.CreatedAt.UTC().Format(time.RFC3339Nano)
	var campusID sql.NullString
	if s.CampusID != nil {
		campusID = sql.NullString{String: s.CampusID.String(), Valid: true}
	}
	var isAct int64 = 1
	if !s.IsActive {
		isAct = 0
	}
	err := r.q.InsertFeeStructure(ctx, InsertFeeStructureParams{
		ID:             s.ID.String(),
		SchoolID:       s.SchoolID.String(),
		AcademicYearID: s.AcademicYearID.String(),
		ClassID:        s.ClassID.String(),
		CampusID:       campusID,
		Code:           s.Code,
		Name:           s.Name,
		Description:    s.Description,
		IsActive:       isAct,
		CreatedAt:      now,
		UpdatedAt:      now,
	})
	if err != nil {
		return shared.Internal(err)
	}
	return nil
}

// --- Fee Structure Versions & Items ---

func (r *FeesRepository) ListFeeStructureVersions(ctx context.Context, feeStructureID shared.ID) ([]fees.FeeStructureVersion, error) {
	rows, err := r.q.ListFeeStructureVersions(ctx, feeStructureID.String())
	if err != nil {
		return nil, shared.Internal(err)
	}
	result := make([]fees.FeeStructureVersion, len(rows))
	for i, row := range rows {
		var pubByID *shared.ID
		if row.PublishedBy.Valid {
			id := shared.ID(row.PublishedBy.String)
			pubByID = &id
		}
		var pubAt *time.Time
		if row.PublishedAt.Valid {
			t, _ := time.Parse(time.RFC3339Nano, row.PublishedAt.String)
			pubAt = &t
		}
		createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
		result[i] = fees.FeeStructureVersion{
			ID:             shared.ID(row.ID),
			FeeStructureID: shared.ID(row.FeeStructureID),
			VersionNo:      int(row.VersionNo),
			Status:         fees.StructureVersionStatus(row.Status),
			Changelog:      row.Changelog,
			CreatedByID:    shared.ID(row.CreatedBy),
			PublishedByID:  pubByID,
			PublishedAt:    pubAt,
			CreatedAt:      createdAt,
			UpdatedAt:      updatedAt,
			Revision:       row.RowVersion,
		}
	}
	return result, nil
}

func (r *FeesRepository) GetFeeStructureVersion(ctx context.Context, id shared.ID) (*fees.FeeStructureVersion, error) {
	row, err := r.q.GetFeeStructureVersion(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("fee structure version")
		}
		return nil, shared.Internal(err)
	}
	var pubByID *shared.ID
	if row.PublishedBy.Valid {
		pid := shared.ID(row.PublishedBy.String)
		pubByID = &pid
	}
	var pubAt *time.Time
	if row.PublishedAt.Valid {
		t, _ := time.Parse(time.RFC3339Nano, row.PublishedAt.String)
		pubAt = &t
	}
	createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
	return &fees.FeeStructureVersion{
		ID:             shared.ID(row.ID),
		FeeStructureID: shared.ID(row.FeeStructureID),
		VersionNo:      int(row.VersionNo),
		Status:         fees.StructureVersionStatus(row.Status),
		Changelog:      row.Changelog,
		CreatedByID:    shared.ID(row.CreatedBy),
		PublishedByID:  pubByID,
		PublishedAt:    pubAt,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		Revision:       row.RowVersion,
	}, nil
}

func (r *FeesRepository) GetPublishedFeeStructureVersion(ctx context.Context, feeStructureID shared.ID) (*fees.FeeStructureVersion, error) {
	row, err := r.q.GetPublishedFeeStructureVersion(ctx, feeStructureID.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("published fee structure version")
		}
		return nil, shared.Internal(err)
	}
	var pubByID *shared.ID
	if row.PublishedBy.Valid {
		pid := shared.ID(row.PublishedBy.String)
		pubByID = &pid
	}
	var pubAt *time.Time
	if row.PublishedAt.Valid {
		t, _ := time.Parse(time.RFC3339Nano, row.PublishedAt.String)
		pubAt = &t
	}
	createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
	return &fees.FeeStructureVersion{
		ID:             shared.ID(row.ID),
		FeeStructureID: shared.ID(row.FeeStructureID),
		VersionNo:      int(row.VersionNo),
		Status:         fees.StructureVersionStatus(row.Status),
		Changelog:      row.Changelog,
		CreatedByID:    shared.ID(row.CreatedBy),
		PublishedByID:  pubByID,
		PublishedAt:    pubAt,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		Revision:       row.RowVersion,
	}, nil
}

func (r *FeesRepository) InsertFeeStructureVersion(ctx context.Context, v *fees.FeeStructureVersion) error {
	now := v.CreatedAt.UTC().Format(time.RFC3339Nano)
	err := r.q.InsertFeeStructureVersion(ctx, InsertFeeStructureVersionParams{
		ID:             v.ID.String(),
		FeeStructureID: v.FeeStructureID.String(),
		VersionNo:      int64(v.VersionNo),
		Status:         string(v.Status),
		Changelog:      v.Changelog,
		CreatedBy:      v.CreatedByID.String(),
		CreatedAt:      now,
		UpdatedAt:      now,
	})
	if err != nil {
		return shared.Internal(err)
	}
	return nil
}

func (r *FeesRepository) PublishFeeStructureVersion(ctx context.Context, id shared.ID, publishedByID shared.ID, publishedAt time.Time) error {
	now := publishedAt.UTC().Format(time.RFC3339Nano)
	err := r.q.PublishFeeStructureVersion(ctx, PublishFeeStructureVersionParams{
		PublishedBy: sql.NullString{String: publishedByID.String(), Valid: true},
		PublishedAt: sql.NullString{String: now, Valid: true},
		UpdatedAt:   now,
		ID:          id.String(),
	})
	if err != nil {
		return shared.Internal(err)
	}
	return nil
}

func (r *FeesRepository) ListFeeStructureItems(ctx context.Context, versionID shared.ID) ([]fees.FeeStructureItem, error) {
	rows, err := r.q.ListFeeStructureItems(ctx, versionID.String())
	if err != nil {
		return nil, shared.Internal(err)
	}
	result := make([]fees.FeeStructureItem, len(rows))
	for i, row := range rows {
		createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
		result[i] = fees.FeeStructureItem{
			ID:                    shared.ID(row.ID),
			FeeStructureVersionID: shared.ID(row.FeeStructureVersionID),
			FeeHeadID:             shared.ID(row.FeeHeadID),
			AmountMinor:           row.AmountMinor,
			Frequency:             fees.FeeFrequency(row.Frequency),
			DueDay:                int(row.DueDay),
			IsOptional:            row.IsOptional == 1,
			CreatedAt:             createdAt,
			UpdatedAt:             updatedAt,
			FeeHeadCode:           row.FeeHeadCode,
			FeeHeadName:           row.FeeHeadName,
		}
	}
	return result, nil
}

func (r *FeesRepository) InsertFeeStructureItem(ctx context.Context, item *fees.FeeStructureItem) error {
	now := item.CreatedAt.UTC().Format(time.RFC3339Nano)
	var isOpt int64 = 0
	if item.IsOptional {
		isOpt = 1
	}
	err := r.q.InsertFeeStructureItem(ctx, InsertFeeStructureItemParams{
		ID:                    item.ID.String(),
		FeeStructureVersionID: item.FeeStructureVersionID.String(),
		FeeHeadID:             item.FeeHeadID.String(),
		AmountMinor:           item.AmountMinor,
		Frequency:             string(item.Frequency),
		DueDay:                int64(item.DueDay),
		IsOptional:            isOpt,
		CreatedAt:             now,
		UpdatedAt:             now,
	})
	if err != nil {
		return shared.Internal(err)
	}
	return nil
}

// --- Discount Policies ---

func (r *FeesRepository) ListDiscountPolicies(ctx context.Context, schoolID shared.ID) ([]fees.DiscountPolicy, error) {
	rows, err := r.q.ListDiscountPolicies(ctx, schoolID.String())
	if err != nil {
		return nil, shared.Internal(err)
	}
	result := make([]fees.DiscountPolicy, len(rows))
	for i, row := range rows {
		var fhID *shared.ID
		if row.FeeHeadID.Valid {
			id := shared.ID(row.FeeHeadID.String)
			fhID = &id
		}
		createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
		result[i] = fees.DiscountPolicy{
			ID:               shared.ID(row.ID),
			SchoolID:         shared.ID(row.SchoolID),
			Code:             row.Code,
			Name:             row.Name,
			Description:      row.Description,
			DiscountType:     fees.DiscountType(row.DiscountType),
			RateBasisPoints:  int(row.RateBasisPoints),
			AmountMinor:      row.AmountMinor,
			FeeHeadID:        fhID,
			RequiresApproval: row.RequiresApproval == 1,
			IsActive:         row.IsActive == 1,
			CreatedAt:        createdAt,
			UpdatedAt:        updatedAt,
			Revision:         row.RowVersion,
		}
	}
	return result, nil
}

func (r *FeesRepository) GetDiscountPolicy(ctx context.Context, id shared.ID) (*fees.DiscountPolicy, error) {
	row, err := r.q.GetDiscountPolicy(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("discount policy")
		}
		return nil, shared.Internal(err)
	}
	var fhID *shared.ID
	if row.FeeHeadID.Valid {
		fid := shared.ID(row.FeeHeadID.String)
		fhID = &fid
	}
	createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
	return &fees.DiscountPolicy{
		ID:               shared.ID(row.ID),
		SchoolID:         shared.ID(row.SchoolID),
		Code:             row.Code,
		Name:             row.Name,
		Description:      row.Description,
		DiscountType:     fees.DiscountType(row.DiscountType),
		RateBasisPoints:  int(row.RateBasisPoints),
		AmountMinor:      row.AmountMinor,
		FeeHeadID:        fhID,
		RequiresApproval: row.RequiresApproval == 1,
		IsActive:         row.IsActive == 1,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
		Revision:         row.RowVersion,
	}, nil
}

func (r *FeesRepository) InsertDiscountPolicy(ctx context.Context, p *fees.DiscountPolicy) error {
	now := p.CreatedAt.UTC().Format(time.RFC3339Nano)
	var fhID sql.NullString
	if p.FeeHeadID != nil {
		fhID = sql.NullString{String: p.FeeHeadID.String(), Valid: true}
	}
	var reqApp, isAct int64
	if p.RequiresApproval {
		reqApp = 1
	}
	if p.IsActive {
		isAct = 1
	}
	err := r.q.InsertDiscountPolicy(ctx, InsertDiscountPolicyParams{
		ID:               p.ID.String(),
		SchoolID:         p.SchoolID.String(),
		Code:             p.Code,
		Name:             p.Name,
		Description:      p.Description,
		DiscountType:     string(p.DiscountType),
		RateBasisPoints:  int64(p.RateBasisPoints),
		AmountMinor:      p.AmountMinor,
		FeeHeadID:        fhID,
		RequiresApproval: reqApp,
		IsActive:         isAct,
		CreatedAt:        now,
		UpdatedAt:        now,
	})
	if err != nil {
		return shared.Internal(err)
	}
	return nil
}

func (r *FeesRepository) UpdateDiscountPolicy(ctx context.Context, p *fees.DiscountPolicy) error {
	now := p.UpdatedAt.UTC().Format(time.RFC3339Nano)
	var fhID sql.NullString
	if p.FeeHeadID != nil {
		fhID = sql.NullString{String: p.FeeHeadID.String(), Valid: true}
	}
	var reqApp, isAct int64
	if p.RequiresApproval {
		reqApp = 1
	}
	if p.IsActive {
		isAct = 1
	}
	err := r.q.UpdateDiscountPolicy(ctx, UpdateDiscountPolicyParams{
		Name:             p.Name,
		Description:      p.Description,
		DiscountType:     string(p.DiscountType),
		RateBasisPoints:  int64(p.RateBasisPoints),
		AmountMinor:      p.AmountMinor,
		FeeHeadID:        fhID,
		RequiresApproval: reqApp,
		IsActive:         isAct,
		UpdatedAt:        now,
		ID:               p.ID.String(),
	})
	if err != nil {
		return shared.Internal(err)
	}
	return nil
}

// --- Student Discounts ---

func (r *FeesRepository) ListStudentDiscounts(ctx context.Context, studentID shared.ID, academicYearID shared.ID) ([]fees.StudentDiscount, error) {
	rows, err := r.q.ListStudentDiscounts(ctx, ListStudentDiscountsParams{
		StudentID:      studentID.String(),
		AcademicYearID: academicYearID.String(),
	})
	if err != nil {
		return nil, shared.Internal(err)
	}
	result := make([]fees.StudentDiscount, len(rows))
	for i, row := range rows {
		var appByID *shared.ID
		if row.ApprovedBy.Valid {
			id := shared.ID(row.ApprovedBy.String)
			appByID = &id
		}
		var appAt *time.Time
		if row.ApprovedAt.Valid {
			t, _ := time.Parse(time.RFC3339Nano, row.ApprovedAt.String)
			appAt = &t
		}
		validFrom, _ := shared.ParseDate(row.ValidFrom)
		var validTo *shared.Date
		if row.ValidTo.Valid {
			d, _ := shared.ParseDate(row.ValidTo.String)
			validTo = &d
		}
		var policyFeeHeadID *shared.ID
		if row.PolicyFeeHeadID.Valid {
			id := shared.ID(row.PolicyFeeHeadID.String)
			policyFeeHeadID = &id
		}
		createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
		result[i] = fees.StudentDiscount{
			ID:                shared.ID(row.ID),
			StudentID:         shared.ID(row.StudentID),
			DiscountPolicyID:  shared.ID(row.DiscountPolicyID),
			AcademicYearID:    shared.ID(row.AcademicYearID),
			Status:            fees.StudentDiscountStatus(row.Status),
			ApprovedByID:      appByID,
			ApprovedAt:        appAt,
			RejectionReason:   row.RejectionReason,
			ValidFrom:         validFrom,
			ValidTo:           validTo,
			CreatedByID:       shared.ID(row.CreatedBy),
			CreatedAt:         createdAt,
			UpdatedAt:         updatedAt,
			Revision:          row.RowVersion,
			PolicyCode:        row.PolicyCode,
			PolicyName:        row.PolicyName,
			DiscountType:      fees.DiscountType(row.DiscountType),
			RateBasisPoints:   int(row.RateBasisPoints),
			PolicyAmountMinor: row.PolicyAmountMinor,
			PolicyFeeHeadID:   policyFeeHeadID,
		}
	}
	return result, nil
}

func (r *FeesRepository) ListPendingStudentDiscounts(ctx context.Context, schoolID shared.ID) ([]fees.StudentDiscount, error) {
	rows, err := r.q.ListPendingStudentDiscounts(ctx, schoolID.String())
	if err != nil {
		return nil, shared.Internal(err)
	}
	result := make([]fees.StudentDiscount, len(rows))
	for i, row := range rows {
		validFrom, _ := shared.ParseDate(row.ValidFrom)
		createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
		result[i] = fees.StudentDiscount{
			ID:               shared.ID(row.ID),
			StudentID:        shared.ID(row.StudentID),
			DiscountPolicyID: shared.ID(row.DiscountPolicyID),
			AcademicYearID:   shared.ID(row.AcademicYearID),
			Status:           fees.StudentDiscountStatus(row.Status),
			ValidFrom:        validFrom,
			CreatedByID:      shared.ID(row.CreatedBy),
			CreatedAt:        createdAt,
			PolicyName:       row.PolicyName,
			StudentFirstName: row.FirstName,
			StudentLastName:  row.LastName,
			AdmissionNo:      row.AdmissionNo,
		}
	}
	return result, nil
}

func (r *FeesRepository) GetStudentDiscount(ctx context.Context, id shared.ID) (*fees.StudentDiscount, error) {
	row, err := r.q.GetStudentDiscount(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("student discount")
		}
		return nil, shared.Internal(err)
	}
	var appByID *shared.ID
	if row.ApprovedBy.Valid {
		aid := shared.ID(row.ApprovedBy.String)
		appByID = &aid
	}
	var appAt *time.Time
	if row.ApprovedAt.Valid {
		t, _ := time.Parse(time.RFC3339Nano, row.ApprovedAt.String)
		appAt = &t
	}
	validFrom, _ := shared.ParseDate(row.ValidFrom)
	var validTo *shared.Date
	if row.ValidTo.Valid {
		d, _ := shared.ParseDate(row.ValidTo.String)
		validTo = &d
	}
	createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
	return &fees.StudentDiscount{
		ID:               shared.ID(row.ID),
		StudentID:        shared.ID(row.StudentID),
		DiscountPolicyID: shared.ID(row.DiscountPolicyID),
		AcademicYearID:   shared.ID(row.AcademicYearID),
		Status:           fees.StudentDiscountStatus(row.Status),
		ApprovedByID:     appByID,
		ApprovedAt:       appAt,
		RejectionReason:  row.RejectionReason,
		ValidFrom:        validFrom,
		ValidTo:          validTo,
		CreatedByID:      shared.ID(row.CreatedBy),
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
		Revision:         row.RowVersion,
	}, nil
}

func (r *FeesRepository) InsertStudentDiscount(ctx context.Context, d *fees.StudentDiscount) error {
	now := d.CreatedAt.UTC().Format(time.RFC3339Nano)
	var validTo sql.NullString
	if d.ValidTo != nil {
		validTo = sql.NullString{String: d.ValidTo.String(), Valid: true}
	}
	err := r.q.InsertStudentDiscount(ctx, InsertStudentDiscountParams{
		ID:               d.ID.String(),
		StudentID:        d.StudentID.String(),
		DiscountPolicyID: d.DiscountPolicyID.String(),
		AcademicYearID:   d.AcademicYearID.String(),
		Status:           string(d.Status),
		ApprovedBy:       sql.NullString{},
		ApprovedAt:       sql.NullString{},
		RejectionReason:  d.RejectionReason,
		ValidFrom:        d.ValidFrom.String(),
		ValidTo:          validTo,
		CreatedBy:        d.CreatedByID.String(),
		CreatedAt:        now,
		UpdatedAt:        now,
	})
	if err != nil {
		return shared.Internal(err)
	}
	return nil
}

func (r *FeesRepository) UpdateStudentDiscountStatus(ctx context.Context, id shared.ID, status fees.StudentDiscountStatus, approvedByID *shared.ID, approvedAt *time.Time, rejectionReason string) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	var appBy, appAtStr sql.NullString
	if approvedByID != nil {
		appBy = sql.NullString{String: approvedByID.String(), Valid: true}
	}
	if approvedAt != nil {
		appAtStr = sql.NullString{String: approvedAt.UTC().Format(time.RFC3339Nano), Valid: true}
	}
	err := r.q.UpdateStudentDiscountStatus(ctx, UpdateStudentDiscountStatusParams{
		Status:          string(status),
		ApprovedBy:      appBy,
		ApprovedAt:      appAtStr,
		RejectionReason: rejectionReason,
		UpdatedAt:       now,
		ID:              id.String(),
	})
	if err != nil {
		return shared.Internal(err)
	}
	return nil
}

// --- Invoices ---

func (r *FeesRepository) ListInvoices(ctx context.Context, filter fees.InvoiceFilter) ([]fees.Invoice, error) {
	limit := int64(filter.Limit)
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	rows, err := r.q.ListInvoices(ctx, ListInvoicesParams{
		SchoolID:       filter.SchoolID.String(),
		Column2:        filter.Status,
		Status:         filter.Status,
		Column4:        filter.AcademicYearID,
		AcademicYearID: filter.AcademicYearID,
		Column6:        filter.StudentID,
		StudentID:      filter.StudentID,
		Limit:          limit,
	})
	if err != nil {
		return nil, shared.Internal(err)
	}
	result := make([]fees.Invoice, len(rows))
	for i, row := range rows {
		issueDate, _ := shared.ParseDate(row.IssueDate)
		dueDate, _ := shared.ParseDate(row.DueDate)
		var voidByID *shared.ID
		if row.VoidedBy.Valid {
			id := shared.ID(row.VoidedBy.String)
			voidByID = &id
		}
		var voidAt *time.Time
		if row.VoidedAt.Valid {
			t, _ := time.Parse(time.RFC3339Nano, row.VoidedAt.String)
			voidAt = &t
		}
		createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
		updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
		result[i] = fees.Invoice{
			ID:                  shared.ID(row.ID),
			SchoolID:            shared.ID(row.SchoolID),
			StudentID:           shared.ID(row.StudentID),
			AcademicYearID:      shared.ID(row.AcademicYearID),
			InvoiceNo:           row.InvoiceNo,
			IssueDate:           issueDate,
			DueDate:             dueDate,
			BillingPeriod:       row.BillingPeriod,
			SubtotalAmountMinor: row.SubtotalAmountMinor,
			DiscountAmountMinor: row.DiscountAmountMinor,
			TotalAmountMinor:    row.TotalAmountMinor,
			PaidAmountMinor:     row.PaidAmountMinor,
			BalanceMinor:        row.BalanceMinor,
			Status:              fees.InvoiceStatus(row.Status),
			VoidReason:          row.VoidReason,
			VoidedByID:          voidByID,
			VoidedAt:            voidAt,
			CreatedByID:         shared.ID(row.CreatedBy),
			CreatedAt:           createdAt,
			UpdatedAt:           updatedAt,
			Revision:            row.RowVersion,
			StudentFirstName:    row.FirstName,
			StudentLastName:     row.LastName,
			AdmissionNo:         row.AdmissionNo,
		}
	}
	return result, nil
}

func (r *FeesRepository) GetInvoice(ctx context.Context, id shared.ID) (*fees.Invoice, error) {
	row, err := r.q.GetInvoice(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("invoice")
		}
		return nil, shared.Internal(err)
	}
	issueDate, _ := shared.ParseDate(row.IssueDate)
	dueDate, _ := shared.ParseDate(row.DueDate)
	var voidByID *shared.ID
	if row.VoidedBy.Valid {
		vid := shared.ID(row.VoidedBy.String)
		voidByID = &vid
	}
	var voidAt *time.Time
	if row.VoidedAt.Valid {
		t, _ := time.Parse(time.RFC3339Nano, row.VoidedAt.String)
		voidAt = &t
	}
	createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
	return &fees.Invoice{
		ID:                  shared.ID(row.ID),
		SchoolID:            shared.ID(row.SchoolID),
		StudentID:           shared.ID(row.StudentID),
		AcademicYearID:      shared.ID(row.AcademicYearID),
		InvoiceNo:           row.InvoiceNo,
		IssueDate:           issueDate,
		DueDate:             dueDate,
		BillingPeriod:       row.BillingPeriod,
		SubtotalAmountMinor: row.SubtotalAmountMinor,
		DiscountAmountMinor: row.DiscountAmountMinor,
		TotalAmountMinor:    row.TotalAmountMinor,
		PaidAmountMinor:     row.PaidAmountMinor,
		BalanceMinor:        row.BalanceMinor,
		Status:              fees.InvoiceStatus(row.Status),
		VoidReason:          row.VoidReason,
		VoidedByID:          voidByID,
		VoidedAt:            voidAt,
		CreatedByID:         shared.ID(row.CreatedBy),
		CreatedAt:           createdAt,
		UpdatedAt:           updatedAt,
		Revision:            row.RowVersion,
		StudentFirstName:    row.FirstName,
		StudentLastName:     row.LastName,
		AdmissionNo:         row.AdmissionNo,
	}, nil
}

func (r *FeesRepository) GetInvoiceByNo(ctx context.Context, schoolID shared.ID, invoiceNo string) (*fees.Invoice, error) {
	row, err := r.q.GetInvoiceByNo(ctx, GetInvoiceByNoParams{
		SchoolID:  schoolID.String(),
		InvoiceNo: invoiceNo,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("invoice")
		}
		return nil, shared.Internal(err)
	}
	issueDate, _ := shared.ParseDate(row.IssueDate)
	dueDate, _ := shared.ParseDate(row.DueDate)
	createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
	return &fees.Invoice{
		ID:                  shared.ID(row.ID),
		SchoolID:            shared.ID(row.SchoolID),
		StudentID:           shared.ID(row.StudentID),
		AcademicYearID:      shared.ID(row.AcademicYearID),
		InvoiceNo:           row.InvoiceNo,
		IssueDate:           issueDate,
		DueDate:             dueDate,
		BillingPeriod:       row.BillingPeriod,
		SubtotalAmountMinor: row.SubtotalAmountMinor,
		DiscountAmountMinor: row.DiscountAmountMinor,
		TotalAmountMinor:    row.TotalAmountMinor,
		PaidAmountMinor:     row.PaidAmountMinor,
		BalanceMinor:        row.BalanceMinor,
		Status:              fees.InvoiceStatus(row.Status),
		CreatedAt:           createdAt,
		UpdatedAt:           updatedAt,
		Revision:            row.RowVersion,
	}, nil
}

func (r *FeesRepository) InsertInvoice(ctx context.Context, inv *fees.Invoice) error {
	now := inv.CreatedAt.UTC().Format(time.RFC3339Nano)
	err := r.q.InsertInvoice(ctx, InsertInvoiceParams{
		ID:                  inv.ID.String(),
		SchoolID:            inv.SchoolID.String(),
		StudentID:           inv.StudentID.String(),
		AcademicYearID:      inv.AcademicYearID.String(),
		InvoiceNo:           inv.InvoiceNo,
		IssueDate:           inv.IssueDate.String(),
		DueDate:             inv.DueDate.String(),
		BillingPeriod:       inv.BillingPeriod,
		SubtotalAmountMinor: inv.SubtotalAmountMinor,
		DiscountAmountMinor: inv.DiscountAmountMinor,
		TotalAmountMinor:    inv.TotalAmountMinor,
		PaidAmountMinor:     inv.PaidAmountMinor,
		BalanceMinor:        inv.BalanceMinor,
		Status:              string(inv.Status),
		CreatedBy:           inv.CreatedByID.String(),
		CreatedAt:           now,
		UpdatedAt:           now,
	})
	if err != nil {
		return shared.Internal(err)
	}
	return nil
}

func (r *FeesRepository) UpdateInvoiceStatus(ctx context.Context, id shared.ID, status fees.InvoiceStatus, paidAmountMinor int64, balanceMinor int64) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	err := r.q.UpdateInvoiceStatus(ctx, UpdateInvoiceStatusParams{
		Status:          string(status),
		PaidAmountMinor: paidAmountMinor,
		BalanceMinor:    balanceMinor,
		UpdatedAt:       now,
		ID:              id.String(),
	})
	if err != nil {
		return shared.Internal(err)
	}
	return nil
}

func (r *FeesRepository) VoidInvoice(ctx context.Context, id shared.ID, reason string, voidedByID shared.ID, voidedAt time.Time) error {
	now := voidedAt.UTC().Format(time.RFC3339Nano)
	err := r.q.VoidInvoice(ctx, VoidInvoiceParams{
		VoidReason: reason,
		VoidedBy:   sql.NullString{String: voidedByID.String(), Valid: true},
		VoidedAt:   sql.NullString{String: now, Valid: true},
		UpdatedAt:  now,
		ID:         id.String(),
	})
	if err != nil {
		return shared.Internal(err)
	}
	return nil
}

func (r *FeesRepository) ListInvoiceItems(ctx context.Context, invoiceID shared.ID) ([]fees.InvoiceItem, error) {
	rows, err := r.q.ListInvoiceItems(ctx, invoiceID.String())
	if err != nil {
		return nil, shared.Internal(err)
	}
	result := make([]fees.InvoiceItem, len(rows))
	for i, row := range rows {
		createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
		result[i] = fees.InvoiceItem{
			ID:                  shared.ID(row.ID),
			InvoiceID:           shared.ID(row.InvoiceID),
			FeeHeadID:           shared.ID(row.FeeHeadID),
			FeeHeadName:         row.FeeHeadName,
			Description:         row.Description,
			OriginalAmountMinor: row.OriginalAmountMinor,
			DiscountAmountMinor: row.DiscountAmountMinor,
			NetAmountMinor:      row.NetAmountMinor,
			CreatedAt:           createdAt,
		}
	}
	return result, nil
}

func (r *FeesRepository) InsertInvoiceItem(ctx context.Context, item *fees.InvoiceItem) error {
	now := item.CreatedAt.UTC().Format(time.RFC3339Nano)
	err := r.q.InsertInvoiceItem(ctx, InsertInvoiceItemParams{
		ID:                  item.ID.String(),
		InvoiceID:           item.InvoiceID.String(),
		FeeHeadID:           item.FeeHeadID.String(),
		FeeHeadName:         item.FeeHeadName,
		Description:         item.Description,
		OriginalAmountMinor: item.OriginalAmountMinor,
		DiscountAmountMinor: item.DiscountAmountMinor,
		NetAmountMinor:      item.NetAmountMinor,
		CreatedAt:           now,
	})
	if err != nil {
		return shared.Internal(err)
	}
	return nil
}

func (r *FeesRepository) CountInvoicesInPeriod(ctx context.Context, schoolID shared.ID, billingPeriod string) (int64, error) {
	count, err := r.q.CountInvoicesInPeriod(ctx, CountInvoicesInPeriodParams{
		SchoolID:      schoolID.String(),
		BillingPeriod: billingPeriod,
	})
	if err != nil {
		return 0, shared.Internal(err)
	}
	return count, nil
}

