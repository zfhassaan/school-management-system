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
		Title:       "Official Admission Application Form",
		Description: "Official admission intake form for student registration, parents details and institutional financials.",
		Sections: []FormSection{
			{
				ID:          "sec_office_financials",
				Title:       "SECTION 1: OFFICE USE & FINANCIALS",
				Description: "Institutional enrollment details, student identification and fee schedule",
				Fields: []FormField{
					{
						ID:          "b_form_cnic",
						Label:       "B-Form / CNIC No.",
						Type:        FieldText,
						Placeholder: "e.g. 34101-1234567-1",
						Required:    false,
					},
					{
						ID:          "adm_fee",
						Label:       "Admission Fee (PKR)",
						Type:        FieldNumber,
						Placeholder: "e.g. 2000",
						Required:    false,
					},
					{
						ID:          "monthly_fee",
						Label:       "Monthly Tuition Fee (PKR)",
						Type:        FieldNumber,
						Placeholder: "e.g. 2500",
						Required:    false,
					},
					{
						ID:          "annual_fee",
						Label:       "Annual Charges (PKR)",
						Type:        FieldNumber,
						Placeholder: "e.g. 3500",
						Required:    false,
					},
					{
						ID:          "paper_fund",
						Label:       "Paper / Exam Fund (PKR)",
						Type:        FieldNumber,
						Placeholder: "e.g. 1000",
						Required:    false,
					},
					{
						ID:          "id_card_fee",
						Label:       "ID Card Fee (PKR)",
						Type:        FieldNumber,
						Placeholder: "e.g. 500",
						Required:    false,
					},
				},
			},
			{
				ID:          "sec_student_personal",
				Title:       "SECTION 2: STUDENT PERSONAL INFO",
				Description: "Student personal particulars, medical history and previous academic record",
				Fields: []FormField{
					{
						ID:          "medical_history",
						Label:       "Medical History / Special Needs",
						Type:        FieldText,
						Placeholder: "e.g. None / Asthma / Allergies",
						Required:    false,
					},
					{
						ID:          "home_address",
						Label:       "Home / Residential Address",
						Type:        FieldText,
						Placeholder: "House #, Street #, Area / City",
						Required:    false,
					},
					{
						ID:          "previous_school",
						Label:       "Previous School Attended",
						Type:        FieldText,
						Placeholder: "Name of previous school",
						Required:    false,
					},
					{
						ID:          "last_grade_result",
						Label:       "Last Grade / Result",
						Type:        FieldText,
						Placeholder: "e.g. Grade 1 - Passed",
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
				ID:          "sec_parents_details",
				Title:       "SECTION 3: PARENTS DETAILS",
				Description: "Father and mother contact details, occupations and qualifications",
				Fields: []FormField{
					{
						ID:          "father_name",
						Label:       "Father's Name",
						Type:        FieldText,
						Placeholder: "Father's full name",
						Required:    false,
					},
					{
						ID:          "father_whatsapp",
						Label:       "Mobile (WhatsApp)",
						Type:        FieldPhone,
						Placeholder: "0321 7404703",
						Required:    false,
					},
					{
						ID:          "father_qualification",
						Label:       "Father Qualification",
						Type:        FieldText,
						Placeholder: "e.g. Matric / Graduate / Master",
						Required:    false,
					},
					{
						ID:          "father_occupation",
						Label:       "Father Occupation",
						Type:        FieldText,
						Placeholder: "e.g. Business / Government / Private",
						Required:    false,
					},
					{
						ID:          "father_office_address",
						Label:       "Office Address",
						Type:        FieldText,
						Placeholder: "Office or business address",
						Required:    false,
					},
					{
						ID:          "father_office_phone",
						Label:       "Office Phone",
						Type:        FieldPhone,
						Placeholder: "e.g. 042-1234567",
						Required:    false,
					},
					{
						ID:          "mother_name",
						Label:       "Mother's Name",
						Type:        FieldText,
						Placeholder: "Mother's full name",
						Required:    false,
					},
					{
						ID:          "mother_phone",
						Label:       "Mother Mobile Number",
						Type:        FieldPhone,
						Placeholder: "0320 7721922",
						Required:    false,
					},
					{
						ID:          "mother_qualification",
						Label:       "Mother Qualification",
						Type:        FieldText,
						Placeholder: "e.g. Matric / Graduate",
						Required:    false,
					},
					{
						ID:          "mother_occupation",
						Label:       "Mother Occupation",
						Type:        FieldText,
						Placeholder: "e.g. Housewife / Teacher",
						Required:    false,
					},
				},
			},
			{
				ID:          "sec_siblings_reference",
				Title:       "SECTION 4: SIBLINGS & REFERENCE",
				Description: "Siblings currently enrolled and emergency contact reference",
				Fields: []FormField{
					{
						ID:          "reference_name",
						Label:       "Reference Name",
						Type:        FieldText,
						Placeholder: "Name of emergency reference",
						Required:    false,
					},
					{
						ID:          "reference_phone",
						Label:       "Reference Phone",
						Type:        FieldPhone,
						Placeholder: "0300 1234567",
						Required:    false,
					},
					{
						ID:          "sibling_in_school",
						Label:       "Sibling in School?",
						Type:        FieldText,
						Placeholder: "Yes / No",
						Required:    false,
					},
					{
						ID:          "sibling_class",
						Label:       "Sibling Class",
						Type:        FieldText,
						Placeholder: "e.g. Class 2",
						Required:    false,
					},
					{
						ID:          "sibling_since",
						Label:       "Enrolled Since",
						Type:        FieldText,
						Placeholder: "e.g. 2024",
						Required:    false,
					},
				},
			},
		},
	}
}

type FormDefinition struct {
	ID          shared.ID `json:"id"`
	SchoolID    shared.ID `json:"schoolId"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Revision    int64     `json:"revision"`
}

type FormVersion struct {
	ID               shared.ID   `json:"id"`
	FormDefinitionID shared.ID   `json:"formDefinitionId"`
	VersionNo        int         `json:"versionNo"`
	Status           FormStatus  `json:"status"`
	Schema           FormSchema  `json:"schema"`
	SchemaHash       string      `json:"schemaHash"`
	PublishedAt      *time.Time  `json:"publishedAt,omitempty"`
	PublishedBy      *shared.ID  `json:"publishedBy,omitempty"`
	CreatedAt        time.Time   `json:"createdAt"`
	UpdatedAt        time.Time   `json:"updatedAt"`
	Revision         int64       `json:"revision"`
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

