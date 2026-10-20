package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"school-management/internal/domain/school"
	"school-management/internal/domain/shared"
)

type SchoolRepository struct {
	q Querier
}

func (r *SchoolRepository) Get(ctx context.Context) (*school.School, error) {
	row, err := r.q.GetSchool(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("school")
		}
		return nil, shared.Internal(err)
	}

	createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)

	var logoID *shared.ID
	if row.LogoDocumentID.Valid {
		id := shared.ID(row.LogoDocumentID.String)
		logoID = &id
	}

	return &school.School{
		ID:             shared.ID(row.ID),
		Name:           row.Name,
		LegalName:      row.LegalName.String,
		Code:           row.Code,
		LogoDocumentID: logoID,
		Address: school.Address{
			Line1:      row.AddressLine1.String,
			Line2:      row.AddressLine2.String,
			City:       row.City.String,
			Region:     row.Region.String,
			PostalCode: row.PostalCode.String,
			Country:    row.Country.String,
		},
		Phone:              row.Phone.String,
		Email:              row.Email.String,
		Website:            row.Website.String,
		Timezone:           row.Timezone,
		CurrencyCode:       row.CurrencyCode,
		CurrencyMinorUnits: int(row.CurrencyMinorUnits),
		Locale:             row.Locale,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
		Revision:           row.RowVersion,
	}, nil
}

func (r *SchoolRepository) Insert(ctx context.Context, s *school.School) error {
	now := s.CreatedAt.UTC().Format(time.RFC3339Nano)
	var logoID sql.NullString
	if s.LogoDocumentID != nil {
		logoID = sql.NullString{String: s.LogoDocumentID.String(), Valid: true}
	}

	err := r.q.InsertSchool(ctx, InsertSchoolParams{
		ID:                 s.ID.String(),
		Name:               s.Name,
		LegalName:          sql.NullString{String: s.LegalName, Valid: s.LegalName != ""},
		Code:               s.Code,
		LogoDocumentID:     logoID,
		AddressLine1:       sql.NullString{String: s.Address.Line1, Valid: s.Address.Line1 != ""},
		AddressLine2:       sql.NullString{String: s.Address.Line2, Valid: s.Address.Line2 != ""},
		City:               sql.NullString{String: s.Address.City, Valid: s.Address.City != ""},
		Region:             sql.NullString{String: s.Address.Region, Valid: s.Address.Region != ""},
		PostalCode:         sql.NullString{String: s.Address.PostalCode, Valid: s.Address.PostalCode != ""},
		Country:            sql.NullString{String: s.Address.Country, Valid: s.Address.Country != ""},
		Phone:              sql.NullString{String: s.Phone, Valid: s.Phone != ""},
		Email:              sql.NullString{String: s.Email, Valid: s.Email != ""},
		Website:            sql.NullString{String: s.Website, Valid: s.Website != ""},
		Timezone:           s.Timezone,
		CurrencyCode:       s.CurrencyCode,
		CurrencyMinorUnits: int64(s.CurrencyMinorUnits),
		Locale:             s.Locale,
		CreatedAt:          now,
		UpdatedAt:          now,
		RowVersion:         1,
	})
	if err != nil {
		return shared.Internal(err)
	}
	s.Revision = 1
	return nil
}

func (r *SchoolRepository) Update(ctx context.Context, s *school.School) error {
	now := s.UpdatedAt.UTC().Format(time.RFC3339Nano)
	var logoID sql.NullString
	if s.LogoDocumentID != nil {
		logoID = sql.NullString{String: s.LogoDocumentID.String(), Valid: true}
	}

	err := r.q.UpdateSchool(ctx, UpdateSchoolParams{
		Name:               s.Name,
		LegalName:          sql.NullString{String: s.LegalName, Valid: s.LegalName != ""},
		Code:               s.Code,
		LogoDocumentID:     logoID,
		AddressLine1:       sql.NullString{String: s.Address.Line1, Valid: s.Address.Line1 != ""},
		AddressLine2:       sql.NullString{String: s.Address.Line2, Valid: s.Address.Line2 != ""},
		City:               sql.NullString{String: s.Address.City, Valid: s.Address.City != ""},
		Region:             sql.NullString{String: s.Address.Region, Valid: s.Address.Region != ""},
		PostalCode:         sql.NullString{String: s.Address.PostalCode, Valid: s.Address.PostalCode != ""},
		Country:            sql.NullString{String: s.Address.Country, Valid: s.Address.Country != ""},
		Phone:              sql.NullString{String: s.Phone, Valid: s.Phone != ""},
		Email:              sql.NullString{String: s.Email, Valid: s.Email != ""},
		Website:            sql.NullString{String: s.Website, Valid: s.Website != ""},
		Timezone:           s.Timezone,
		CurrencyCode:       s.CurrencyCode,
		CurrencyMinorUnits: int64(s.CurrencyMinorUnits),
		Locale:             s.Locale,
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

func (r *SchoolRepository) GetSettings(ctx context.Context, schoolID shared.ID) (*school.Settings, error) {
	row, err := r.q.GetSchoolSettings(ctx, schoolID.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("settings")
		}
		return nil, shared.Internal(err)
	}
	updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
	return &school.Settings{
		SchoolID:              shared.ID(row.SchoolID),
		AdmissionNoFormat:     row.AdmissionNoFormat,
		InvoicePrefix:         row.InvoicePrefix,
		ReceiptPrefix:         row.ReceiptPrefix,
		RoundingUnitMinor:     row.RoundingUnitMinor,
		RoundingMode:          shared.RoundingMode(row.RoundingMode),
		DefaultDueDay:         int(row.DefaultDueDay),
		LateFeeGraceDays:      int(row.LateFeeGraceDays),
		DiscountCombination:   school.DiscountCombination(row.DiscountCombination),
		RequireSecondApprover: row.RequireSecondApprover == 1,
		DateFormat:            row.DateFormat,
		UpdatedAt:             updatedAt,
		Revision:              row.RowVersion,
	}, nil
}

func (r *SchoolRepository) InsertSettings(ctx context.Context, s *school.Settings, at time.Time) error {
	now := at.UTC().Format(time.RFC3339Nano)
	secondApprover := int64(0)
	if s.RequireSecondApprover {
		secondApprover = 1
	}
	err := r.q.InsertSchoolSettings(ctx, InsertSchoolSettingsParams{
		SchoolID:              s.SchoolID.String(),
		AdmissionNoFormat:     s.AdmissionNoFormat,
		InvoicePrefix:         s.InvoicePrefix,
		ReceiptPrefix:         s.ReceiptPrefix,
		RoundingUnitMinor:     s.RoundingUnitMinor,
		RoundingMode:          string(s.RoundingMode),
		DefaultDueDay:         int64(s.DefaultDueDay),
		LateFeeGraceDays:      int64(s.LateFeeGraceDays),
		DiscountCombination:   string(s.DiscountCombination),
		RequireSecondApprover: secondApprover,
		DateFormat:            s.DateFormat,
		CreatedAt:             now,
		UpdatedAt:             now,
		RowVersion:            1,
	})
	if err != nil {
		return shared.Internal(err)
	}
	s.Revision = 1
	return nil
}

func (r *SchoolRepository) UpdateSettings(ctx context.Context, s *school.Settings, at time.Time) error {
	now := at.UTC().Format(time.RFC3339Nano)
	secondApprover := int64(0)
	if s.RequireSecondApprover {
		secondApprover = 1
	}
	err := r.q.UpdateSchoolSettings(ctx, UpdateSchoolSettingsParams{
		AdmissionNoFormat:     s.AdmissionNoFormat,
		InvoicePrefix:         s.InvoicePrefix,
		ReceiptPrefix:         s.ReceiptPrefix,
		RoundingUnitMinor:     s.RoundingUnitMinor,
		RoundingMode:          string(s.RoundingMode),
		DefaultDueDay:         int64(s.DefaultDueDay),
		LateFeeGraceDays:      int64(s.LateFeeGraceDays),
		DiscountCombination:   string(s.DiscountCombination),
		RequireSecondApprover: secondApprover,
		DateFormat:            s.DateFormat,
		UpdatedAt:             now,
		SchoolID:              s.SchoolID.String(),
		RowVersion:            s.Revision,
	})
	if err != nil {
		return shared.Internal(err)
	}
	s.Revision++
	return nil
}

