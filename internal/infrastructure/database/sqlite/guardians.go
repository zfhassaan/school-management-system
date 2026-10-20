package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"school-management/internal/domain/guardians"
	"school-management/internal/domain/shared"
)

type GuardiansRepository struct {
	q Querier
}

func (r *GuardiansRepository) Get(ctx context.Context, id shared.ID) (*guardians.Guardian, error) {
	row, err := r.q.GetGuardian(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("guardian")
		}
		return nil, shared.Internal(err)
	}
	createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
	return &guardians.Guardian{
		ID:                     shared.ID(row.ID),
		SchoolID:               shared.ID(row.SchoolID),
		FirstName:              row.FirstName,
		LastName:               row.LastName.String,
		PhonePrimary:           row.PhonePrimary,
		PhonePrimaryNormalized: row.PhonePrimaryNormalized,
		PhoneSecondary:         row.PhoneSecondary.String,
		Email:                  row.Email.String,
		Occupation:             row.Occupation.String,
		NationalID:             row.NationalID.String,
		Address: shared.Address{
			Line1:      row.AddressLine1.String,
			Line2:      row.AddressLine2.String,
			City:       row.City.String,
			Region:     row.Region.String,
			PostalCode: row.PostalCode.String,
			Country:    row.Country.String,
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Revision:  row.RowVersion,
	}, nil
}

func (r *GuardiansRepository) GetByPhone(ctx context.Context, schoolID shared.ID, phoneNormalized string) (*guardians.Guardian, error) {
	row, err := r.q.GetGuardianByPhone(ctx, GetGuardianByPhoneParams{
		SchoolID:                schoolID.String(),
		PhonePrimaryNormalized: phoneNormalized,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.NotFound("guardian")
		}
		return nil, shared.Internal(err)
	}
	createdAt, _ := time.Parse(time.RFC3339Nano, row.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339Nano, row.UpdatedAt)
	return &guardians.Guardian{
		ID:                     shared.ID(row.ID),
		SchoolID:               shared.ID(row.SchoolID),
		FirstName:              row.FirstName,
		LastName:               row.LastName.String,
		PhonePrimary:           row.PhonePrimary,
		PhonePrimaryNormalized: row.PhonePrimaryNormalized,
		PhoneSecondary:         row.PhoneSecondary.String,
		Email:                  row.Email.String,
		Occupation:             row.Occupation.String,
		NationalID:             row.NationalID.String,
		Address: shared.Address{
			Line1:      row.AddressLine1.String,
			Line2:      row.AddressLine2.String,
			City:       row.City.String,
			Region:     row.Region.String,
			PostalCode: row.PostalCode.String,
			Country:    row.Country.String,
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Revision:  row.RowVersion,
	}, nil
}

func (r *GuardiansRepository) Insert(ctx context.Context, g *guardians.Guardian) error {
	now := g.CreatedAt.UTC().Format(time.RFC3339Nano)
	err := r.q.InsertGuardian(ctx, InsertGuardianParams{
		ID:                     g.ID.String(),
		SchoolID:               g.SchoolID.String(),
		FirstName:              g.FirstName,
		LastName:               sql.NullString{String: g.LastName, Valid: g.LastName != ""},
		PhonePrimary:           g.PhonePrimary,
		PhonePrimaryNormalized: g.PhonePrimaryNormalized,
		PhoneSecondary:         sql.NullString{String: g.PhoneSecondary, Valid: g.PhoneSecondary != ""},
		Email:                  sql.NullString{String: g.Email, Valid: g.Email != ""},
		Occupation:             sql.NullString{String: g.Occupation, Valid: g.Occupation != ""},
		NationalID:             sql.NullString{String: g.NationalID, Valid: g.NationalID != ""},
		AddressLine1:           sql.NullString{String: g.Address.Line1, Valid: g.Address.Line1 != ""},
		AddressLine2:           sql.NullString{String: g.Address.Line2, Valid: g.Address.Line2 != ""},
		City:                   sql.NullString{String: g.Address.City, Valid: g.Address.City != ""},
		Region:                 sql.NullString{String: g.Address.Region, Valid: g.Address.Region != ""},
		PostalCode:             sql.NullString{String: g.Address.PostalCode, Valid: g.Address.PostalCode != ""},
		Country:                sql.NullString{String: g.Address.Country, Valid: g.Address.Country != ""},
		CreatedAt:              now,
		UpdatedAt:              now,
	})
	if err != nil {
		return shared.Internal(err)
	}
	g.Revision = 1
	return nil
}

func (r *GuardiansRepository) Update(ctx context.Context, g *guardians.Guardian) error {
	now := g.UpdatedAt.UTC().Format(time.RFC3339Nano)
	err := r.q.UpdateGuardian(ctx, UpdateGuardianParams{
		FirstName:              g.FirstName,
		LastName:               sql.NullString{String: g.LastName, Valid: g.LastName != ""},
		PhonePrimary:           g.PhonePrimary,
		PhonePrimaryNormalized: g.PhonePrimaryNormalized,
		PhoneSecondary:         sql.NullString{String: g.PhoneSecondary, Valid: g.PhoneSecondary != ""},
		Email:                  sql.NullString{String: g.Email, Valid: g.Email != ""},
		Occupation:             sql.NullString{String: g.Occupation, Valid: g.Occupation != ""},
		NationalID:             sql.NullString{String: g.NationalID, Valid: g.NationalID != ""},
		AddressLine1:           sql.NullString{String: g.Address.Line1, Valid: g.Address.Line1 != ""},
		AddressLine2:           sql.NullString{String: g.Address.Line2, Valid: g.Address.Line2 != ""},
		City:                   sql.NullString{String: g.Address.City, Valid: g.Address.City != ""},
		Region:                 sql.NullString{String: g.Address.Region, Valid: g.Address.Region != ""},
		PostalCode:             sql.NullString{String: g.Address.PostalCode, Valid: g.Address.PostalCode != ""},
		Country:                sql.NullString{String: g.Address.Country, Valid: g.Address.Country != ""},
		UpdatedAt:              now,
		ID:                     g.ID.String(),
		RowVersion:             g.Revision,
	})
	if err != nil {
		return shared.Internal(err)
	}
	g.Revision++
	return nil
}

func (r *GuardiansRepository) ListForStudent(ctx context.Context, studentID shared.ID) ([]guardians.StudentGuardianLink, error) {
	rows, err := r.q.ListStudentGuardians(ctx, studentID.String())
	if err != nil {
		return nil, shared.Internal(err)
	}
	links := make([]guardians.StudentGuardianLink, len(rows))
	for i, row := range rows {
		links[i] = guardians.StudentGuardianLink{
			StudentID:          shared.ID(row.StudentID),
			GuardianID:         shared.ID(row.GuardianID),
			Relationship:       guardians.Relationship(row.Relationship),
			IsPrimary:          row.IsPrimary == 1,
			IsEmergencyContact: row.IsEmergencyContact == 1,
			ReceivesBilling:    row.ReceivesBilling == 1,
			CanPickup:          row.CanPickup == 1,
			Priority:           int(row.Priority),
			FirstName:          row.FirstName,
			LastName:           row.LastName.String,
			PhonePrimary:       row.PhonePrimary,
			PhoneSecondary:     row.PhoneSecondary.String,
			Email:              row.Email.String,
			Occupation:         row.Occupation.String,
		}
	}
	return links, nil
}

func (r *GuardiansRepository) SetStudentGuardians(ctx context.Context, studentID shared.ID, links []guardians.StudentGuardianLink) error {
	if err := r.q.DeleteStudentGuardians(ctx, studentID.String()); err != nil {
		return shared.Internal(err)
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	for _, l := range links {
		isPrimary := int64(0)
		if l.IsPrimary {
			isPrimary = 1
		}
		isEmergency := int64(0)
		if l.IsEmergencyContact {
			isEmergency = 1
		}
		receivesBilling := int64(0)
		if l.ReceivesBilling {
			receivesBilling = 1
		}
		canPickup := int64(0)
		if l.CanPickup {
			canPickup = 1
		}

		err := r.q.InsertStudentGuardian(ctx, InsertStudentGuardianParams{
			StudentID:          studentID.String(),
			GuardianID:         l.GuardianID.String(),
			Relationship:       string(l.Relationship),
			IsPrimary:          isPrimary,
			IsEmergencyContact: isEmergency,
			ReceivesBilling:    receivesBilling,
			CanPickup:          canPickup,
			Priority:           int64(l.Priority),
			CreatedAt:          now,
			UpdatedAt:          now,
		})
		if err != nil {
			return shared.Internal(err)
		}
	}
	return nil
}

