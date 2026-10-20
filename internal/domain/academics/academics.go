package academics

import (
	"context"
	"strings"
	"time"

	"school-management/internal/domain/shared"
)

type YearStatus string

const (
	YearPlanned YearStatus = "planned"
	YearActive  YearStatus = "active"
	YearClosed  YearStatus = "closed"
)

type AcademicYear struct {
	ID        shared.ID
	SchoolID  shared.ID
	Code      string // e.g. "2026-27"
	Name      string // e.g. "Academic Year 2026-2027"
	StartsOn  shared.Date
	EndsOn    shared.Date
	Status    YearStatus
	CreatedAt time.Time
	UpdatedAt time.Time
	Revision  int64
}

func (y *AcademicYear) Validate() error {
	var v shared.Validator
	v.Check(y.Code != "", "code", "code is required")
	v.Check(y.Name != "", "name", "name is required")
	v.Check(!y.StartsOn.IsZero(), "startsOn", "startsOn is required")
	v.Check(!y.EndsOn.IsZero(), "endsOn", "endsOn is required")
	v.Check(y.EndsOn.After(y.StartsOn), "endsOn", "endsOn must be strictly after startsOn")
	v.Check(y.Status == YearPlanned || y.Status == YearActive || y.Status == YearClosed, "status", "invalid year status")
	return v.Err()
}

type AcademicTerm struct {
	ID             shared.ID
	AcademicYearID shared.ID
	Sequence       int
	Name           string
	StartsOn       shared.Date
	EndsOn         shared.Date
}

type Class struct {
	ID          shared.ID
	SchoolID    shared.ID
	Code        string // e.g. "G1", "NUR", "KG"
	Name        string // e.g. "Grade 1"
	LevelOrder  int    // e.g. 1, 2, 3 ...
	NextClassID *shared.ID
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Revision    int64
}

func (c *Class) Validate() error {
	var v shared.Validator
	v.Check(c.Code != "", "code", "code is required")
	v.Check(c.Name != "", "name", "name is required")
	v.Check(c.LevelOrder >= 0, "levelOrder", "levelOrder must be non-negative")
	return v.Err()
}

type Section struct {
	ID                 shared.ID
	ClassID            shared.ID
	AcademicYearID     shared.ID
	CampusID           *shared.ID
	Name               string // e.g. "Red", "A", "Rose"
	Capacity           *int
	ClassTeacherUserID *shared.ID
	IsActive           bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Revision           int64
}

func (s *Section) Validate() error {
	var v shared.Validator
	v.Check(strings.TrimSpace(s.Name) != "", "name", "section name is required")
	if s.Capacity != nil {
		v.Check(*s.Capacity > 0, "capacity", "capacity must be positive")
	}
	return v.Err()
}

type Campus struct {
	ID        shared.ID
	SchoolID  shared.ID
	Code      string
	Name      string
	Address   string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	Revision  int64
}

type Repository interface {
	ListAcademicYears(ctx context.Context, schoolID shared.ID) ([]AcademicYear, error)
	GetAcademicYear(ctx context.Context, id shared.ID) (*AcademicYear, error)
	GetActiveAcademicYear(ctx context.Context, schoolID shared.ID) (*AcademicYear, error)
	InsertAcademicYear(ctx context.Context, y *AcademicYear) error
	UpdateAcademicYear(ctx context.Context, y *AcademicYear) error

	ListAcademicTerms(ctx context.Context, yearID shared.ID) ([]AcademicTerm, error)
	ReplaceAcademicTerms(ctx context.Context, yearID shared.ID, terms []AcademicTerm) error

	ListClasses(ctx context.Context, schoolID shared.ID) ([]Class, error)
	GetClass(ctx context.Context, id shared.ID) (*Class, error)
	InsertClass(ctx context.Context, c *Class) error
	UpdateClass(ctx context.Context, c *Class) error

	ListSectionsByYear(ctx context.Context, yearID shared.ID) ([]Section, error)
	GetSection(ctx context.Context, id shared.ID) (*Section, error)
	InsertSection(ctx context.Context, s *Section) error
	UpdateSection(ctx context.Context, s *Section) error

	ListCampuses(ctx context.Context, schoolID shared.ID) ([]Campus, error)
	InsertCampus(ctx context.Context, c *Campus) error
}
