package guardians

import (
	"context"
	"strings"

	"school-management/internal/application/services/auditlog"
	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/uow"
	"school-management/internal/domain/guardians"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
)

type CreateGuardianInput struct {
	FirstName      string         `json:"firstName"`
	LastName       string         `json:"lastName"`
	PhonePrimary   string         `json:"phonePrimary"`
	PhoneSecondary string         `json:"phoneSecondary"`
	Email          string         `json:"email"`
	Occupation     string         `json:"occupation"`
	NationalID     string         `json:"nationalId"`
	Address        shared.Address `json:"address"`
}

type UpdateGuardianInput struct {
	ID             shared.ID      `json:"id"`
	FirstName      string         `json:"firstName"`
	LastName       string         `json:"lastName"`
	PhonePrimary   string         `json:"phonePrimary"`
	PhoneSecondary string         `json:"phoneSecondary"`
	Email          string         `json:"email"`
	Occupation     string         `json:"occupation"`
	NationalID     string         `json:"nationalId"`
	Address        shared.Address `json:"address"`
}

type StudentGuardianLinkInput struct {
	GuardianID         shared.ID              `json:"guardianId"`
	Relationship       guardians.Relationship `json:"relationship"`
	IsPrimary          bool                   `json:"isPrimary"`
	IsEmergencyContact bool                   `json:"isEmergencyContact"`
	ReceivesBilling    bool                   `json:"receivesBilling"`
	CanPickup          bool                   `json:"canPickup"`
	Priority           int                    `json:"priority"`
}

type SetStudentGuardiansInput struct {
	StudentID shared.ID                  `json:"studentId"`
	Guardians []StudentGuardianLinkInput `json:"guardians"`
}

type GuardianCommandsHandler struct {
	uow      uow.UnitOfWork
	clock    shared.Clock
	ids      shared.IDGenerator
	recorder auditlog.Recorder
}

func NewGuardianCommandsHandler(
	uow uow.UnitOfWork,
	clock shared.Clock,
	ids shared.IDGenerator,
) *GuardianCommandsHandler {
	return &GuardianCommandsHandler{
		uow:      uow,
		clock:    clock,
		ids:      ids,
		recorder: auditlog.Recorder{IDs: ids, Clock: clock},
	}
}

func (h *GuardianCommandsHandler) CreateGuardian(ctx context.Context, in CreateGuardianInput) (shared.ID, error) {
	if _, err := authz.Require(ctx, identity.PermGuardiansManage); err != nil {
		return "", err
	}
	id := h.ids.NewID()
	now := h.clock.Now()

	err := h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}
		g := &guardians.Guardian{
			ID:             id,
			SchoolID:       sch.ID,
			FirstName:      strings.TrimSpace(in.FirstName),
			LastName:       strings.TrimSpace(in.LastName),
			PhonePrimary:   strings.TrimSpace(in.PhonePrimary),
			PhoneSecondary: strings.TrimSpace(in.PhoneSecondary),
			Email:          strings.TrimSpace(in.Email),
			Occupation:     strings.TrimSpace(in.Occupation),
			NationalID:     strings.TrimSpace(in.NationalID),
			Address:        in.Address,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		if err := g.Validate(); err != nil {
			return err
		}

		if err := r.Guardians.Insert(ctx, g); err != nil {
			return err
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "guardian.created",
			EntityType: "guardian",
			EntityID:   id.String(),
			Reason:     "guardian created",
			After:      map[string]any{"name": g.FirstName + " " + g.LastName, "phone": g.PhonePrimary},
		})
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func (h *GuardianCommandsHandler) UpdateGuardian(ctx context.Context, in UpdateGuardianInput) error {
	if _, err := authz.Require(ctx, identity.PermGuardiansManage); err != nil {
		return err
	}
	now := h.clock.Now()

	return h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		g, err := r.Guardians.Get(ctx, in.ID)
		if err != nil {
			return err
		}
		g.FirstName = strings.TrimSpace(in.FirstName)
		g.LastName = strings.TrimSpace(in.LastName)
		g.PhonePrimary = strings.TrimSpace(in.PhonePrimary)
		g.PhoneSecondary = strings.TrimSpace(in.PhoneSecondary)
		g.Email = strings.TrimSpace(in.Email)
		g.Occupation = strings.TrimSpace(in.Occupation)
		g.NationalID = strings.TrimSpace(in.NationalID)
		g.Address = in.Address
		g.UpdatedAt = now

		if err := g.Validate(); err != nil {
			return err
		}

		if err := r.Guardians.Update(ctx, g); err != nil {
			return err
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "guardian.updated",
			EntityType: "guardian",
			EntityID:   in.ID.String(),
			Reason:     "guardian updated",
			After:      map[string]any{"name": g.FirstName + " " + g.LastName, "phone": g.PhonePrimary},
		})
	})
}

func (h *GuardianCommandsHandler) SetStudentGuardians(ctx context.Context, in SetStudentGuardiansInput) error {
	if _, err := authz.Require(ctx, identity.PermGuardiansManage); err != nil {
		return err
	}
	links := make([]guardians.StudentGuardianLink, len(in.Guardians))
	for i, item := range in.Guardians {
		links[i] = guardians.StudentGuardianLink{
			StudentID:          in.StudentID,
			GuardianID:         item.GuardianID,
			Relationship:       item.Relationship,
			IsPrimary:          item.IsPrimary,
			IsEmergencyContact: item.IsEmergencyContact,
			ReceivesBilling:    item.ReceivesBilling,
			CanPickup:          item.CanPickup,
			Priority:           item.Priority,
		}
	}

	return h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		if err := r.Guardians.SetStudentGuardians(ctx, in.StudentID, links); err != nil {
			return err
		}
		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "student_guardians.updated",
			EntityType: "student",
			EntityID:   in.StudentID.String(),
			Reason:     "student guardians updated",
			After:      map[string]any{"count": len(links)},
		})
	})
}

