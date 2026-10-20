package enrollments

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"school-management/internal/application/services/auditlog"
	"school-management/internal/application/services/authz"
	"school-management/internal/application/services/uow"
	"school-management/internal/domain/identity"
	"school-management/internal/domain/shared"
	"school-management/internal/domain/students"
)

type EnrollStudentInput struct {
	StudentID      shared.ID   `json:"studentId"`
	AcademicYearID shared.ID   `json:"academicYearId"`
	ClassID        shared.ID   `json:"classId"`
	SectionID      *shared.ID  `json:"sectionId,omitempty"`
	RollNo         *string     `json:"rollNo,omitempty"`
	EnrolledOn     shared.Date `json:"enrolledOn"`
}

type PromotionItem struct {
	StudentID       shared.ID  `json:"studentId"`
	TargetClassID   shared.ID  `json:"targetClassId"`
	TargetSectionID *shared.ID `json:"targetSectionId,omitempty"`
	RollNo          *string    `json:"rollNo,omitempty"`
	Action          string     `json:"action"` // "promoted", "retained", "graduated"
}

type ExecutePromotionBatchInput struct {
	FromAcademicYearID shared.ID       `json:"fromAcademicYearId"`
	ToAcademicYearID   shared.ID       `json:"toAcademicYearId"`
	EffectiveDate      shared.Date     `json:"effectiveDate"`
	Notes              string          `json:"notes"`
	Items              []PromotionItem `json:"items"`
}

type EnrollmentCommandsHandler struct {
	uow      uow.UnitOfWork
	clock    shared.Clock
	ids      shared.IDGenerator
	recorder auditlog.Recorder
}

func NewEnrollmentCommandsHandler(
	uow uow.UnitOfWork,
	clock shared.Clock,
	ids shared.IDGenerator,
) *EnrollmentCommandsHandler {
	return &EnrollmentCommandsHandler{
		uow:      uow,
		clock:    clock,
		ids:      ids,
		recorder: auditlog.Recorder{IDs: ids, Clock: clock},
	}
}

func (h *EnrollmentCommandsHandler) EnrollStudent(ctx context.Context, in EnrollStudentInput) (shared.ID, error) {
	if _, err := authz.Require(ctx, identity.PermStudentsEdit); err != nil {
		return "", err
	}
	id := h.ids.NewID()
	now := h.clock.Now()

	err := h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		// Close any existing active enrollment
		if active, err := r.Enrollments.GetActive(ctx, in.StudentID); err == nil && active != nil {
			_ = r.Enrollments.Close(ctx, active.ID, students.EnrollmentTransferred, in.EnrolledOn, active.Revision)
		}

		enr := &students.Enrollment{
			ID:             id,
			StudentID:      in.StudentID,
			AcademicYearID: in.AcademicYearID,
			ClassID:        in.ClassID,
			SectionID:      in.SectionID,
			RollNo:         in.RollNo,
			Status:         students.EnrollmentActive,
			EnrolledOn:     in.EnrolledOn,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		if err := r.Enrollments.Insert(ctx, enr); err != nil {
			return err
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "enrollment.created",
			EntityType: "enrollment",
			EntityID:   id.String(),
			Reason:     "student enrolled",
			After: map[string]any{
				"studentId":      in.StudentID.String(),
				"academicYearId": in.AcademicYearID.String(),
				"classId":        in.ClassID.String(),
			},
		})
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func (h *EnrollmentCommandsHandler) ExecutePromotionBatch(ctx context.Context, in ExecutePromotionBatchInput) (shared.ID, error) {
	actor, err := authz.Require(ctx, identity.PermAcademicsPromote)
	if err != nil {
		return "", err
	}
	if len(in.Items) == 0 {
		return "", shared.Validation("no students provided for promotion", nil)
	}

	batchID := h.ids.NewID()
	now := h.clock.Now()

	summaryBytes, _ := json.Marshal(map[string]any{
		"notes":         strings.TrimSpace(in.Notes),
		"totalStudents": len(in.Items),
	})

	err = h.uow.Write(ctx, func(ctx context.Context, r uow.Repos) error {
		batch := &students.PromotionBatch{
			ID:                 batchID,
			FromAcademicYearID: in.FromAcademicYearID,
			ToAcademicYearID:   in.ToAcademicYearID,
			Status:             students.PromotionExecuted,
			SummaryJSON:        string(summaryBytes),
			ExecutedBy:         &actor.UserID,
			ExecutedAt:         &now,
			CreatedAt:          now,
			UpdatedAt:          now,
		}

		if err := r.Enrollments.InsertPromotionBatch(ctx, batch); err != nil {
			return fmt.Errorf("failed to insert promotion batch: %w", err)
		}

		for _, item := range in.Items {
			active, err := r.Enrollments.GetActive(ctx, item.StudentID)
			if err == nil && active != nil {
				_ = r.Enrollments.Close(ctx, active.ID, students.EnrollmentPromoted, in.EffectiveDate, active.Revision)
			}

			if item.Action == "graduated" {
				s, err := r.Students.Get(ctx, item.StudentID)
				if err == nil {
					s.Status = students.StatusGraduated
					s.UpdatedAt = now
					_ = r.Students.Update(ctx, s)
					_ = r.Students.AppendStatusHistory(ctx, students.StatusHistoryEntry{
						ID:          h.ids.NewID(),
						StudentID:   item.StudentID,
						FromStatus:  string(students.StatusActive),
						ToStatus:    string(students.StatusGraduated),
						EffectiveOn: in.EffectiveDate,
						Reason:      "Graduation via promotion batch",
						ChangedBy:   actor.UserID,
						ChangedAt:   now,
					})
				}
				continue
			}

			newEnrollmentID := h.ids.NewID()
			var prevID *shared.ID
			if active != nil {
				prevID = &active.ID
			}
			pBatchID := batchID

			newEnr := &students.Enrollment{
				ID:                   newEnrollmentID,
				StudentID:            item.StudentID,
				AcademicYearID:       in.ToAcademicYearID,
				ClassID:              item.TargetClassID,
				SectionID:            item.TargetSectionID,
				RollNo:               item.RollNo,
				Status:               students.EnrollmentActive,
				EnrolledOn:           in.EffectiveDate,
				PreviousEnrollmentID: prevID,
				PromotionBatchID:     &pBatchID,
				CreatedAt:            now,
				UpdatedAt:            now,
			}
			if err := r.Enrollments.Insert(ctx, newEnr); err != nil {
				return fmt.Errorf("failed to insert promoted enrollment: %w", err)
			}
		}

		return h.recorder.Record(ctx, r.Audit, auditlog.Change{
			Action:     "academics.promoted_batch",
			EntityType: "promotion_batch",
			EntityID:   batchID.String(),
			Reason:     "promoted students to next academic year",
			After: map[string]any{
				"totalStudents": len(in.Items),
				"fromYear":      in.FromAcademicYearID.String(),
				"toYear":        in.ToAcademicYearID.String(),
			},
		})
	})
	if err != nil {
		return "", err
	}
	return batchID, nil
}

