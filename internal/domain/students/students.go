package students

import (
	"context"
	"strings"
	"time"

	"school-management/internal/domain/shared"
)

type Status string

const (
	StatusActive    Status = "active"
	StatusInactive  Status = "inactive"
	StatusGraduated Status = "graduated"
	StatusWithdrawn Status = "withdrawn"
	StatusAlumni    Status = "alumni"
)

type Gender string

const (
	GenderMale        Gender = "male"
	GenderFemale      Gender = "female"
	GenderOther       Gender = "other"
	GenderUnspecified Gender = "unspecified"
)

type Student struct {
	ID            shared.ID
	SchoolID      shared.ID
	AdmissionNo   string
	FirstName     string
	MiddleName    string
	LastName      string
	SearchName    string // lower unaccented full name
	DateOfBirth   shared.Date
	Gender        Gender
	NationalID    string
	Address       shared.Address
	PhotoDocID    *shared.ID
	AdmissionDate shared.Date
	Status        Status
	Notes         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Revision      int64
}

// ComputeSearchName generates a normalized lowercase search string.
func ComputeSearchName(first, middle, last string) string {
	parts := []string{}
	if f := strings.TrimSpace(first); f != "" {
		parts = append(parts, f)
	}
	if m := strings.TrimSpace(middle); m != "" {
		parts = append(parts, m)
	}
	if l := strings.TrimSpace(last); l != "" {
		parts = append(parts, l)
	}
	return strings.ToLower(strings.Join(parts, " "))
}

func (s *Student) FullName() string {
	parts := []string{s.FirstName}
	if s.MiddleName != "" {
		parts = append(parts, s.MiddleName)
	}
	parts = append(parts, s.LastName)
	return strings.Join(parts, " ")
}

func (s *Student) Validate() error {
	var v shared.Validator
	v.Check(strings.TrimSpace(s.AdmissionNo) != "", "admissionNo", "admission number is required")
	v.Check(strings.TrimSpace(s.FirstName) != "", "firstName", "first name is required")
	v.Check(strings.TrimSpace(s.LastName) != "", "lastName", "last name is required")
	v.Check(!s.DateOfBirth.IsZero(), "dateOfBirth", "date of birth is required")
	v.Check(!s.AdmissionDate.IsZero(), "admissionDate", "admission date is required")
	v.Check(s.Gender == GenderMale || s.Gender == GenderFemale || s.Gender == GenderOther || s.Gender == GenderUnspecified, "gender", "invalid gender")
	s.SearchName = ComputeSearchName(s.FirstName, s.MiddleName, s.LastName)
	return v.Err()
}

type StatusHistoryEntry struct {
	ID          shared.ID
	StudentID   shared.ID
	FromStatus  string
	ToStatus    string
	EffectiveOn shared.Date
	Reason      string
	ChangedBy   shared.ID
	ChangedAt   time.Time
}

type EnrollmentStatus string

const (
	EnrollmentActive      EnrollmentStatus = "active"
	EnrollmentPromoted    EnrollmentStatus = "promoted"
	EnrollmentRepeated    EnrollmentStatus = "repeated"
	EnrollmentTransferred EnrollmentStatus = "transferred"
	EnrollmentWithdrawn   EnrollmentStatus = "withdrawn"
	EnrollmentCompleted   EnrollmentStatus = "completed"
	EnrollmentCancelled   EnrollmentStatus = "cancelled"
)

type Enrollment struct {
	ID                     shared.ID
	StudentID              shared.ID
	AcademicYearID         shared.ID
	ClassID                shared.ID
	SectionID              *shared.ID
	RollNo                 *string
	Status                 EnrollmentStatus
	EnrolledOn             shared.Date
	EndedOn                *shared.Date
	PreviousEnrollmentID   *shared.ID
	PromotionBatchID       *shared.ID
	AdmissionApplicationID *shared.ID
	CreatedAt              time.Time
	UpdatedAt              time.Time
	Revision               int64
	// Expanded view fields
	ClassName   string
	SectionName string
	YearCode    string
}

type PromotionBatchStatus string

const (
	PromotionPreviewed  PromotionBatchStatus = "previewed"
	PromotionExecuted   PromotionBatchStatus = "executed"
	PromotionRolledBack PromotionBatchStatus = "rolled_back"
)

type PromotionBatch struct {
	ID                 shared.ID
	FromAcademicYearID shared.ID
	ToAcademicYearID   shared.ID
	Status             PromotionBatchStatus
	SummaryJSON        string
	ExecutedBy         *shared.ID
	ExecutedAt         *time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Revision           int64
}

type StudentFilter struct {
	SchoolID       shared.ID
	Status         string
	ClassID        string
	SectionID      string
	SearchQuery    string
	CursorName     string
	CursorID       string
	Limit          int
}

type Repository interface {
	Get(ctx context.Context, id shared.ID) (*Student, error)
	GetByAdmissionNo(ctx context.Context, schoolID shared.ID, no string) (*Student, error)
	Insert(ctx context.Context, s *Student) error
	Update(ctx context.Context, s *Student) error
	AppendStatusHistory(ctx context.Context, e StatusHistoryEntry) error
	ListStatusHistory(ctx context.Context, studentID shared.ID) ([]StatusHistoryEntry, error)
	ListCursor(ctx context.Context, f StudentFilter) ([]Student, error)
	SearchFTS(ctx context.Context, schoolID shared.ID, query string, limit int) ([]Student, error)
	Count(ctx context.Context, schoolID shared.ID) (int64, error)
	FindSiblings(ctx context.Context, studentID shared.ID) ([]Student, error)
}

type EnrollmentRepository interface {
	GetActive(ctx context.Context, studentID shared.ID) (*Enrollment, error)
	ListHistory(ctx context.Context, studentID shared.ID) ([]Enrollment, error)
	Insert(ctx context.Context, e *Enrollment) error
	Close(ctx context.Context, id shared.ID, status EnrollmentStatus, endedOn shared.Date, revision int64) error
	ListSectionRoster(ctx context.Context, sectionID shared.ID) ([]Enrollment, error)
	InsertPromotionBatch(ctx context.Context, batch *PromotionBatch) error
	GetPromotionBatch(ctx context.Context, id shared.ID) (*PromotionBatch, error)
}
