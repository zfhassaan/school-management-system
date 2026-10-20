package school

import (
	"context"
	"strings"

	"school-management/internal/application/services/auditlog"
	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/uow"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/school"
	"school-management/internal/domain/shared"
)

type UpdateSchoolProfileInput struct {
	Name        string         `json:"name"`
	LegalName   string         `json:"legalName"`
	Code        string         `json:"code"`
	Address     school.Address `json:"address"`
	Phone       string         `json:"phone"`
	Email       string         `json:"email"`
	Website     string         `json:"website"`
	Timezone    string         `json:"timezone"`
	Locale      string         `json:"locale"`
}

type UpdateSchoolSettingsInput struct {
	AdmissionNoFormat     string `json:"admissionNoFormat"`
	InvoicePrefix         string `json:"invoicePrefix"`
	ReceiptPrefix         string `json:"receiptPrefix"`
	RoundingUnitMinor     int64  `json:"roundingUnitMinor"`
	RoundingMode          string `json:"roundingMode"`
	DefaultDueDay         int    `json:"defaultDueDay"`
	LateFeeGraceDays      int    `json:"lateFeeGraceDays"`
	DiscountCombination   string `json:"discountCombination"`
	RequireSecondApprover bool   `json:"requireSecondApprover"`
	DateFormat            string `json:"dateFormat"`
}

type SchoolCommandsHandler struct {
	uow      uow.UnitOfWork
	clock    shared.Clock
	ids      shared.IDGenerator
	recorder auditlog.Recorder
}

func NewSchoolCommandsHandler(
	uow uow.UnitOfWork,
	clock shared.Clock,
	ids shared.IDGenerator,
) *SchoolCommandsHandler {
	return &SchoolCommandsHandler{
		uow:      uow,
		clock:    clock,
		ids:      ids,
		recorder: auditlog.Recorder{IDs: ids, Clock: clock},
	}
}

func (h *SchoolCommandsHandler) UpdateProfile(ctx context.Context, in UpdateSchoolProfileInput) error {
	if _, err := authz.Require(ctx, identity.PermSettingsManage); err != nil {
		return err
	}

	now := h.clock.Now()

	return h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}

		before := map[string]any{
			"name":    sch.Name,
			"phone":   sch.Phone,
			"email":   sch.Email,
			"website": sch.Website,
		}

		sch.Name = strings.TrimSpace(in.Name)
		sch.LegalName = strings.TrimSpace(in.LegalName)
		sch.Code = strings.ToUpper(strings.TrimSpace(in.Code))
		sch.Address = in.Address
		sch.Phone = strings.TrimSpace(in.Phone)
		sch.Email = strings.TrimSpace(in.Email)
		sch.Website = strings.TrimSpace(in.Website)
		sch.Timezone = strings.TrimSpace(in.Timezone)
		if in.Locale != "" {
			sch.Locale = in.Locale
		}
		sch.UpdatedAt = now

		if err := sch.Validate(); err != nil {
			return err
		}

		if err := r.School.Update(ctx, sch); err != nil {
			return err
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "school.profile_updated",
			EntityType: "school",
			EntityID:   sch.ID.String(),
			Reason:     "school profile updated by administrator",
			Before:     before,
			After: map[string]any{
				"name":    sch.Name,
				"phone":   sch.Phone,
				"email":   sch.Email,
				"website": sch.Website,
			},
		})
	})
}

func (h *SchoolCommandsHandler) UpdateSettings(ctx context.Context, in UpdateSchoolSettingsInput) error {
	if _, err := authz.Require(ctx, identity.PermSettingsManage); err != nil {
		return err
	}

	now := h.clock.Now()

	return h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}

		settings, err := r.School.GetSettings(ctx, sch.ID)
		if err != nil {
			return err
		}

		settings.AdmissionNoFormat = in.AdmissionNoFormat
		settings.InvoicePrefix = strings.ToUpper(strings.TrimSpace(in.InvoicePrefix))
		settings.ReceiptPrefix = strings.ToUpper(strings.TrimSpace(in.ReceiptPrefix))
		settings.RoundingUnitMinor = in.RoundingUnitMinor
		settings.RoundingMode = shared.RoundingMode(in.RoundingMode)
		settings.DefaultDueDay = in.DefaultDueDay
		settings.LateFeeGraceDays = in.LateFeeGraceDays
		settings.DiscountCombination = school.DiscountCombination(in.DiscountCombination)
		settings.RequireSecondApprover = in.RequireSecondApprover
		settings.DateFormat = in.DateFormat
		settings.UpdatedAt = now

		if err := settings.Validate(); err != nil {
			return err
		}

		if err := r.School.UpdateSettings(ctx, settings, now); err != nil {
			return err
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "school.settings_updated",
			EntityType: "school_settings",
			EntityID:   sch.ID.String(),
			Reason:     "school settings updated by administrator",
		})
	})
}

