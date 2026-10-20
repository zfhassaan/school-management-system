package academics

import (
	"context"
	"strings"

	"school-management/internal/application/services/auditlog"
	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/uow"
	"school-management/internal/domain/academics"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
)

type CreateAcademicYearInput struct {
	Name     string               `json:"name"`
	Code     string               `json:"code"`
	StartsOn shared.Date          `json:"startsOn"`
	EndsOn   shared.Date          `json:"endsOn"`
	Status   academics.YearStatus `json:"status"`
}

type UpdateAcademicYearInput struct {
	ID       shared.ID            `json:"id"`
	Name     string               `json:"name"`
	StartsOn shared.Date          `json:"startsOn"`
	EndsOn   shared.Date          `json:"endsOn"`
	Status   academics.YearStatus `json:"status"`
}

type AcademicTermInput struct {
	Sequence int         `json:"sequence"`
	Name     string      `json:"name"`
	StartsOn shared.Date `json:"startsOn"`
	EndsOn   shared.Date `json:"endsOn"`
}

type SetTermsInput struct {
	AcademicYearID shared.ID           `json:"academicYearId"`
	Terms          []AcademicTermInput `json:"terms"`
}

type CreateCampusInput struct {
	Name    string `json:"name"`
	Code    string `json:"code"`
	Address string `json:"address"`
}

type CreateClassInput struct {
	Name        string     `json:"name"`
	Code        string     `json:"code"`
	LevelOrder  int        `json:"levelOrder"`
	NextClassID *shared.ID `json:"nextClassId,omitempty"`
}

type UpdateClassInput struct {
	ID          shared.ID  `json:"id"`
	Name        string     `json:"name"`
	Code        string     `json:"code"`
	LevelOrder  int        `json:"levelOrder"`
	NextClassID *shared.ID `json:"nextClassId,omitempty"`
	IsActive    bool       `json:"isActive"`
}

type CreateSectionInput struct {
	ClassID            shared.ID  `json:"classId"`
	AcademicYearID     shared.ID  `json:"academicYearId"`
	CampusID           *shared.ID `json:"campusId,omitempty"`
	Name               string     `json:"name"`
	Capacity           *int       `json:"capacity,omitempty"`
	ClassTeacherUserID *shared.ID `json:"classTeacherUserId,omitempty"`
}

type UpdateSectionInput struct {
	ID                 shared.ID  `json:"id"`
	Name               string     `json:"name"`
	Capacity           *int       `json:"capacity,omitempty"`
	ClassTeacherUserID *shared.ID `json:"classTeacherUserId,omitempty"`
	IsActive           bool       `json:"isActive"`
}

type AcademicCommandsHandler struct {
	uow      uow.UnitOfWork
	clock    shared.Clock
	ids      shared.IDGenerator
	recorder auditlog.Recorder
}

func NewAcademicCommandsHandler(
	uow uow.UnitOfWork,
	clock shared.Clock,
	ids shared.IDGenerator,
) *AcademicCommandsHandler {
	return &AcademicCommandsHandler{
		uow:      uow,
		clock:    clock,
		ids:      ids,
		recorder: auditlog.Recorder{IDs: ids, Clock: clock},
	}
}

func (h *AcademicCommandsHandler) CreateAcademicYear(ctx context.Context, in CreateAcademicYearInput) (shared.ID, error) {
	if _, err := authz.Require(ctx, identity.PermAcademicsManage); err != nil {
		return "", err
	}
	id := h.ids.NewID()
	now := h.clock.Now()

	status := in.Status
	if status == "" {
		status = academics.YearPlanned
	}

	err := h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}
		ay := &academics.AcademicYear{
			ID:        id,
			SchoolID:  sch.ID,
			Code:      strings.TrimSpace(in.Code),
			Name:      strings.TrimSpace(in.Name),
			StartsOn:  in.StartsOn,
			EndsOn:    in.EndsOn,
			Status:    status,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := ay.Validate(); err != nil {
			return err
		}
		if err := r.Academics.InsertAcademicYear(ctx, ay); err != nil {
			return err
		}
		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "academic_year.created",
			EntityType: "academic_year",
			EntityID:   id.String(),
			Reason:     "academic year created",
			After:      map[string]any{"name": ay.Name, "code": ay.Code, "status": ay.Status},
		})
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func (h *AcademicCommandsHandler) UpdateAcademicYear(ctx context.Context, in UpdateAcademicYearInput) error {
	if _, err := authz.Require(ctx, identity.PermAcademicsManage); err != nil {
		return err
	}
	now := h.clock.Now()

	return h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		ay, err := r.Academics.GetAcademicYear(ctx, in.ID)
		if err != nil {
			return err
		}
		ay.Name = strings.TrimSpace(in.Name)
		ay.StartsOn = in.StartsOn
		ay.EndsOn = in.EndsOn
		if in.Status != "" {
			ay.Status = in.Status
		}
		ay.UpdatedAt = now

		if err := ay.Validate(); err != nil {
			return err
		}
		if err := r.Academics.UpdateAcademicYear(ctx, ay); err != nil {
			return err
		}
		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "academic_year.updated",
			EntityType: "academic_year",
			EntityID:   in.ID.String(),
			Reason:     "academic year updated",
			After:      map[string]any{"name": ay.Name, "status": ay.Status},
		})
	})
}

func (h *AcademicCommandsHandler) SetTerms(ctx context.Context, in SetTermsInput) error {
	if _, err := authz.Require(ctx, identity.PermAcademicsManage); err != nil {
		return err
	}
	terms := make([]academics.AcademicTerm, len(in.Terms))
	for i, t := range in.Terms {
		terms[i] = academics.AcademicTerm{
			ID:             h.ids.NewID(),
			AcademicYearID: in.AcademicYearID,
			Sequence:       t.Sequence,
			Name:           strings.TrimSpace(t.Name),
			StartsOn:       t.StartsOn,
			EndsOn:         t.EndsOn,
		}
	}

	return h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		if err := r.Academics.ReplaceAcademicTerms(ctx, in.AcademicYearID, terms); err != nil {
			return err
		}
		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "academic_year.terms_replaced",
			EntityType: "academic_year",
			EntityID:   in.AcademicYearID.String(),
			Reason:     "academic terms updated",
			After:      map[string]any{"count": len(terms)},
		})
	})
}

func (h *AcademicCommandsHandler) CreateCampus(ctx context.Context, in CreateCampusInput) (shared.ID, error) {
	if _, err := authz.Require(ctx, identity.PermAcademicsManage); err != nil {
		return "", err
	}
	id := h.ids.NewID()
	now := h.clock.Now()

	err := h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}
		c := &academics.Campus{
			ID:        id,
			SchoolID:  sch.ID,
			Code:      strings.TrimSpace(in.Code),
			Name:      strings.TrimSpace(in.Name),
			Address:   strings.TrimSpace(in.Address),
			IsActive:  true,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := r.Academics.InsertCampus(ctx, c); err != nil {
			return err
		}
		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "campus.created",
			EntityType: "campus",
			EntityID:   id.String(),
			Reason:     "campus created",
			After:      map[string]any{"name": c.Name, "code": c.Code},
		})
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func (h *AcademicCommandsHandler) CreateClass(ctx context.Context, in CreateClassInput) (shared.ID, error) {
	if _, err := authz.Require(ctx, identity.PermAcademicsManage); err != nil {
		return "", err
	}
	id := h.ids.NewID()
	now := h.clock.Now()

	err := h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}
		cls := &academics.Class{
			ID:          id,
			SchoolID:    sch.ID,
			Code:        strings.TrimSpace(in.Code),
			Name:        strings.TrimSpace(in.Name),
			LevelOrder:  in.LevelOrder,
			NextClassID: in.NextClassID,
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if err := cls.Validate(); err != nil {
			return err
		}
		if err := r.Academics.InsertClass(ctx, cls); err != nil {
			return err
		}
		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "class.created",
			EntityType: "class",
			EntityID:   id.String(),
			Reason:     "class created",
			After:      map[string]any{"name": cls.Name, "code": cls.Code, "levelOrder": cls.LevelOrder},
		})
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func (h *AcademicCommandsHandler) UpdateClass(ctx context.Context, in UpdateClassInput) error {
	if _, err := authz.Require(ctx, identity.PermAcademicsManage); err != nil {
		return err
	}
	now := h.clock.Now()

	return h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		cls, err := r.Academics.GetClass(ctx, in.ID)
		if err != nil {
			return err
		}
		cls.Name = strings.TrimSpace(in.Name)
		cls.Code = strings.TrimSpace(in.Code)
		cls.LevelOrder = in.LevelOrder
		cls.NextClassID = in.NextClassID
		cls.IsActive = in.IsActive
		cls.UpdatedAt = now

		if err := cls.Validate(); err != nil {
			return err
		}
		if err := r.Academics.UpdateClass(ctx, cls); err != nil {
			return err
		}
		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "class.updated",
			EntityType: "class",
			EntityID:   in.ID.String(),
			Reason:     "class updated",
			After:      map[string]any{"name": cls.Name, "code": cls.Code, "levelOrder": cls.LevelOrder},
		})
	})
}

func (h *AcademicCommandsHandler) CreateSection(ctx context.Context, in CreateSectionInput) (shared.ID, error) {
	if _, err := authz.Require(ctx, identity.PermAcademicsManage); err != nil {
		return "", err
	}
	id := h.ids.NewID()
	now := h.clock.Now()

	err := h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		sec := &academics.Section{
			ID:                 id,
			ClassID:            in.ClassID,
			AcademicYearID:     in.AcademicYearID,
			CampusID:           in.CampusID,
			Name:               strings.TrimSpace(in.Name),
			Capacity:           in.Capacity,
			ClassTeacherUserID: in.ClassTeacherUserID,
			IsActive:           true,
			CreatedAt:          now,
			UpdatedAt:          now,
		}
		if err := sec.Validate(); err != nil {
			return err
		}
		if err := r.Academics.InsertSection(ctx, sec); err != nil {
			return err
		}
		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "section.created",
			EntityType: "section",
			EntityID:   id.String(),
			Reason:     "section created",
			After:      map[string]any{"classId": in.ClassID.String(), "name": sec.Name},
		})
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func (h *AcademicCommandsHandler) UpdateSection(ctx context.Context, in UpdateSectionInput) error {
	if _, err := authz.Require(ctx, identity.PermAcademicsManage); err != nil {
		return err
	}
	now := h.clock.Now()

	return h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		sec, err := r.Academics.GetSection(ctx, in.ID)
		if err != nil {
			return err
		}
		sec.Name = strings.TrimSpace(in.Name)
		sec.Capacity = in.Capacity
		sec.ClassTeacherUserID = in.ClassTeacherUserID
		sec.IsActive = in.IsActive
		sec.UpdatedAt = now

		if err := sec.Validate(); err != nil {
			return err
		}
		if err := r.Academics.UpdateSection(ctx, sec); err != nil {
			return err
		}
		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "section.updated",
			EntityType: "section",
			EntityID:   in.ID.String(),
			Reason:     "section updated",
			After:      map[string]any{"name": sec.Name},
		})
	})
}

