package students

import (
	"context"
	"fmt"
	"strings"

	"school-management/internal/application/services/auditlog"
	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/uow"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
	"school-management/internal/domain/students"
)

type RegisterStudentInput struct {
	AdmissionNo            string             `json:"admissionNo"`
	FirstName              string             `json:"firstName"`
	MiddleName             string             `json:"middleName"`
	LastName               string             `json:"lastName"`
	DateOfBirth            shared.Date        `json:"dateOfBirth"`
	Gender                 students.Gender    `json:"gender"`
	NationalID             string             `json:"nationalId"`
	Address                shared.Address     `json:"address"`
	PhotoDocID            *shared.ID         `json:"photoDocId,omitempty"`
	AdmissionDate          shared.Date        `json:"admissionDate"`
	Notes                  string             `json:"notes"`
	AcademicYearID         shared.ID          `json:"academicYearId"`
	ClassID                shared.ID          `json:"classId"`
	SectionID              *shared.ID         `json:"sectionId,omitempty"`
	RollNo                 *string            `json:"rollNo,omitempty"`
	AdmissionApplicationID *shared.ID         `json:"admissionApplicationId,omitempty"`
}

type UpdateStudentInput struct {
	ID          shared.ID       `json:"id"`
	FirstName   string          `json:"firstName"`
	MiddleName  string          `json:"middleName"`
	LastName    string          `json:"lastName"`
	DateOfBirth shared.Date     `json:"dateOfBirth"`
	Gender      students.Gender `json:"gender"`
	NationalID  string          `json:"nationalId"`
	Address     shared.Address  `json:"address"`
	PhotoDocID  *shared.ID      `json:"photoDocId,omitempty"`
	Notes       string          `json:"notes"`
}

type ChangeStudentStatusInput struct {
	StudentID   shared.ID       `json:"studentId"`
	NewStatus   students.Status `json:"newStatus"`
	EffectiveOn shared.Date     `json:"effectiveOn"`
	Reason      string          `json:"reason"`
}

type StudentCommandsHandler struct {
	uow      uow.UnitOfWork
	clock    shared.Clock
	ids      shared.IDGenerator
	recorder auditlog.Recorder
}

func NewStudentCommandsHandler(
	uow uow.UnitOfWork,
	clock shared.Clock,
	ids shared.IDGenerator,
) *StudentCommandsHandler {
	return &StudentCommandsHandler{
		uow:      uow,
		clock:    clock,
		ids:      ids,
		recorder: auditlog.Recorder{IDs: ids, Clock: clock},
	}
}

func (h *StudentCommandsHandler) RegisterStudent(ctx context.Context, in RegisterStudentInput) (shared.ID, error) {
	actor, err := authz.Require(ctx, identity.PermStudentsCreate)
	if err != nil {
		return "", err
	}
	id := h.ids.NewID()
	enrollmentID := h.ids.NewID()
	now := h.clock.Now()

	err = h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		sch, err := r.School.Get(ctx)
		if err != nil {
			return err
		}

		admissionNo := strings.TrimSpace(in.AdmissionNo)
		if admissionNo == "" {
			count, err := r.Students.Count(ctx, sch.ID)
			if err != nil {
				return err
			}
			admissionNo = fmt.Sprintf("ADM-%d-%04d", now.Year(), count+1)
		}

		s := &students.Student{
			ID:            id,
			SchoolID:      sch.ID,
			AdmissionNo:   admissionNo,
			FirstName:     strings.TrimSpace(in.FirstName),
			MiddleName:    strings.TrimSpace(in.MiddleName),
			LastName:      strings.TrimSpace(in.LastName),
			DateOfBirth:   in.DateOfBirth,
			Gender:        in.Gender,
			NationalID:    strings.TrimSpace(in.NationalID),
			Address:       in.Address,
			PhotoDocID:    in.PhotoDocID,
			AdmissionDate: in.AdmissionDate,
			Status:        students.StatusActive,
			Notes:         strings.TrimSpace(in.Notes),
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		if err := s.Validate(); err != nil {
			return err
		}

		if err := r.Students.Insert(ctx, s); err != nil {
			return err
		}

		// Initial enrollment
		enr := &students.Enrollment{
			ID:                     enrollmentID,
			StudentID:              id,
			AcademicYearID:         in.AcademicYearID,
			ClassID:                in.ClassID,
			SectionID:              in.SectionID,
			RollNo:                 in.RollNo,
			Status:                 students.EnrollmentActive,
			EnrolledOn:             in.AdmissionDate,
			AdmissionApplicationID: in.AdmissionApplicationID,
			CreatedAt:              now,
			UpdatedAt:              now,
		}
		if err := r.Enrollments.Insert(ctx, enr); err != nil {
			return err
		}

		// Record initial status history
		sh := students.StatusHistoryEntry{
			ID:          h.ids.NewID(),
			StudentID:   id,
			FromStatus:  "",
			ToStatus:    string(students.StatusActive),
			EffectiveOn: in.AdmissionDate,
			Reason:      "Initial registration",
			ChangedBy:   actor.UserID,
			ChangedAt:   now,
		}
		_ = r.Students.AppendStatusHistory(ctx, sh)

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "student.registered",
			EntityType: "student",
			EntityID:   id.String(),
			Reason:     "student registered",
			After:      map[string]any{"admissionNo": admissionNo, "name": s.FullName()},
		})
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func (h *StudentCommandsHandler) UpdateStudent(ctx context.Context, in UpdateStudentInput) error {
	if _, err := authz.Require(ctx, identity.PermStudentsEdit); err != nil {
		return err
	}
	now := h.clock.Now()

	return h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		s, err := r.Students.Get(ctx, in.ID)
		if err != nil {
			return err
		}

		s.FirstName = strings.TrimSpace(in.FirstName)
		s.MiddleName = strings.TrimSpace(in.MiddleName)
		s.LastName = strings.TrimSpace(in.LastName)
		s.DateOfBirth = in.DateOfBirth
		s.Gender = in.Gender
		s.NationalID = strings.TrimSpace(in.NationalID)
		s.Address = in.Address
		s.PhotoDocID = in.PhotoDocID
		s.Notes = strings.TrimSpace(in.Notes)
		s.UpdatedAt = now

		if err := s.Validate(); err != nil {
			return err
		}

		if err := r.Students.Update(ctx, s); err != nil {
			return err
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "student.updated",
			EntityType: "student",
			EntityID:   in.ID.String(),
			Reason:     "student profile updated",
			After:      map[string]any{"name": s.FullName()},
		})
	})
}

func (h *StudentCommandsHandler) ChangeStudentStatus(ctx context.Context, in ChangeStudentStatusInput) error {
	actor, err := authz.Require(ctx, identity.PermStudentsStatus)
	if err != nil {
		return err
	}
	now := h.clock.Now()

	return h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		s, err := r.Students.Get(ctx, in.StudentID)
		if err != nil {
			return err
		}
		if s.Status == in.NewStatus {
			return nil
		}

		fromStatus := s.Status
		s.Status = in.NewStatus
		s.UpdatedAt = now

		if err := r.Students.Update(ctx, s); err != nil {
			return err
		}

		sh := students.StatusHistoryEntry{
			ID:          h.ids.NewID(),
			StudentID:   in.StudentID,
			FromStatus:  string(fromStatus),
			ToStatus:    string(in.NewStatus),
			EffectiveOn: in.EffectiveOn,
			Reason:      strings.TrimSpace(in.Reason),
			ChangedBy:   actor.UserID,
			ChangedAt:   now,
		}
		if err := r.Students.AppendStatusHistory(ctx, sh); err != nil {
			return err
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "student.status_changed",
			EntityType: "student",
			EntityID:   in.StudentID.String(),
			Reason:     in.Reason,
			Before:     map[string]any{"status": fromStatus},
			After:      map[string]any{"status": in.NewStatus},
		})
	})
}
