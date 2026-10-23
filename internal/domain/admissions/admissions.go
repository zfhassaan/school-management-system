package admissions

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"school-management/internal/domain/shared"
)

type FormStatus string

const (
	FormDraft     FormStatus = "draft"
	FormPublished FormStatus = "published"
	FormRetired   FormStatus = "retired"
)

type FieldType string

const (
	FieldText       FieldType = "text"
	FieldNumber     FieldType = "number"
	FieldDate       FieldType = "date"
	FieldDropdown   FieldType = "dropdown"
	FieldRadio      FieldType = "radio"
	FieldCheckbox   FieldType = "checkbox"
	FieldPhone      FieldType = "phone"
	FieldEmail      FieldType = "email"
	FieldAddress    FieldType = "address"
	FieldFileUpload FieldType = "file_upload"
	FieldTextarea   FieldType = "textarea"
)

type FormField struct {
	ID          string      `json:"id"`
	Label       string      `json:"label"`
	Type        FieldType   `json:"type"`
	Placeholder string      `json:"placeholder,omitempty"`
	Required    bool        `json:"required"`
	Options     []string    `json:"options,omitempty"`
	DependsOn   *Visibility `json:"dependsOn,omitempty"`
}

type Visibility struct {
	FieldID string `json:"fieldId"`
	Equals  string `json:"equals"`
}

type FormSection struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description,omitempty"`
	Fields      []FormField `json:"fields"`
}

type FormSchema struct {
	Title             string        `json:"title"`
	Description       string        `json:"description,omitempty"`
	Sections          []FormSection `json:"sections"`
	RequiredDocuments []string      `json:"requiredDocuments,omitempty"`
}

// DefaultFormSchema returns the standard institutional admission form schema.
func DefaultFormSchema() FormSchema {
	return FormSchema{
		Title:       "General Admission Application Form",
		Description: "Official admission intake form for student registration and enrollment.",
		Sections: []FormSection{
			{
				ID:          "sec_student_background",
				Title:       "Academic Background",
				Description: "Previous school background and student personal details",
				Fields: []FormField{
					{
						ID:          "prev_school",
						Label:       "Previous School Attended",
						Type:        FieldText,
						Placeholder: "Name of previous school",
						Required:    false,
					},
					{
						ID:          "last_grade",
						Label:       "Last Grade / Class Passed",
						Type:        FieldText,
						Placeholder: "e.g. Grade 4",
						Required:    false,
					},
					{
						ID:          "blood_group",
						Label:       "Blood Group",
						Type:        FieldDropdown,
						Placeholder: "Select Blood Group",
						Required:    false,
						Options:     []string{"A+", "A-", "B+", "B-", "O+", "O-", "AB+", "AB-"},
					},
				},
			},
			{
				ID:          "sec_emergency",
				Title:       "Emergency & Medical Contact",
				Description: "Emergency contact numbers and special health requirements",
				Fields: []FormField{
					{
						ID:          "emergency_contact",
						Label:       "Emergency Contact Person",
						Type:        FieldText,
						Placeholder: "Name of emergency contact",
						Required:    false,
					},
					{
						ID:          "emergency_phone",
						Label:       "Emergency Phone Number",
						Type:        FieldPhone,
						Placeholder: "+92-300-1234567",
						Required:    false,
					},
				},
			},
			{
				ID:          "sec_notes",
				Title:       "General Notes",
				Description: "Optional additional notes or remarks",
				Fields: []FormField{
					{
						ID:          "notes",
						Label:       "Notes / Remarks",
						Type:        FieldTextarea,
						Placeholder: "Optional background or remarks...",
						Required:    false,
					},
				},
			},
		},
	}
}

type FormDefinition struct {
	ID          shared.ID
	SchoolID    shared.ID
	Code        string
	Name        string
	Description string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Revision    int64
}

type FormVersion struct {
	ID               shared.ID
	FormDefinitionID shared.ID
	VersionNo        int
	Status           FormStatus
	Schema           FormSchema
	SchemaHash       string
	PublishedAt      *time.Time
	PublishedBy      *shared.ID
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Revision         int64
}

type ApplicationStatus string

const (
	AppDraft       ApplicationStatus = "draft"
	AppSubmitted   ApplicationStatus = "submitted"
	AppUnderReview ApplicationStatus = "under_review"
	AppApproved    ApplicationStatus = "approved"
	AppRejected    ApplicationStatus = "rejected"
	AppEnrolled    ApplicationStatus = "enrolled"
	AppWithdrawn   ApplicationStatus = "withdrawn"
)

type Application struct {
	ID                      shared.ID
	SchoolID                shared.ID
	ApplicationNo           string
	FormVersionID           shared.ID
	AcademicYearID          shared.ID
	AppliedClassID          shared.ID
	Status                  ApplicationStatus
	ApplicantFirstName      string
	ApplicantLastName       string
	ApplicantDOB            *shared.Date
	ApplicantGender         string
	GuardianPhoneNormalized string
	StudentID               *shared.ID
	Response                map[string]any
	SubmittedAt             *time.Time
	DecidedAt               *time.Time
	DecidedBy               *shared.ID
	DecisionReason          string
	CreatedBy               shared.ID
	CreatedAt               time.Time
	UpdatedAt               time.Time
	Revision                int64
	// Expanded view fields
	ClassName string
	YearCode  string
}

// ComputeSchemaHash returns SHA256 of normalized canonical JSON schema.
func ComputeSchemaHash(s FormSchema) (string, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(b)
	return hex.EncodeToString(hash[:]), nil
}

// ValidateResponse checks submitted response data against the published form schema.
func ValidateResponse(schema FormSchema, values map[string]any) error {
	var v shared.Validator
	for _, sec := range schema.Sections {
		for _, f := range sec.Fields {
			// Check conditional visibility
			if f.DependsOn != nil {
				parentVal, ok := values[f.DependsOn.FieldID]
				if !ok || fmt.Sprint(parentVal) != f.DependsOn.Equals {
					// Field is hidden due to condition; skip requirement check
					continue
				}
			}

			val, exists := values[f.ID]
			if f.Required {
				if !exists || val == nil || strings.TrimSpace(fmt.Sprint(val)) == "" {
					v.Add(f.ID, fmt.Sprintf("%s is required", f.Label))
					continue
				}
			}

			if exists && val != nil && strings.TrimSpace(fmt.Sprint(val)) != "" {
				sval := fmt.Sprint(val)
				switch f.Type {
				case FieldDate:
					if _, err := shared.ParseDate(sval); err != nil {
						v.Add(f.ID, fmt.Sprintf("%s must be a valid date (YYYY-MM-DD)", f.Label))
					}
				case FieldEmail:
					if !strings.Contains(sval, "@") {
						v.Add(f.ID, fmt.Sprintf("%s must be a valid email", f.Label))
					}
				}
			}
		}
	}
	return v.Err()
}

type Repository interface {
	ListFormDefinitions(ctx context.Context, schoolID shared.ID) ([]FormDefinition, error)
	GetFormDefinition(ctx context.Context, id shared.ID) (*FormDefinition, error)
	InsertFormDefinition(ctx context.Context, d *FormDefinition) error

	ListFormVersions(ctx context.Context, formDefID shared.ID) ([]FormVersion, error)
	GetFormVersion(ctx context.Context, id shared.ID) (*FormVersion, error)
	GetPublishedVersion(ctx context.Context, formDefID shared.ID) (*FormVersion, error)
	InsertFormVersion(ctx context.Context, v *FormVersion) error
	PublishFormVersion(ctx context.Context, id shared.ID, publishedBy shared.ID, at time.Time, hash string) error

	ListApplications(ctx context.Context, schoolID shared.ID, status string, yearID string, limit int) ([]Application, error)
	GetApplication(ctx context.Context, id shared.ID) (*Application, error)
	InsertApplication(ctx context.Context, a *Application) error
	UpdateStatus(ctx context.Context, a *Application, reason string, actorID shared.ID, at time.Time) error
}

