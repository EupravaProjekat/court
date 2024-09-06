package Models

type User struct {
	Uuid     string    `bson:"uuid,omitempty" json:"uuid,omitempty"`
	Email    string    `bson:"email,omitempty" json:"email,omitempty"`
	Role     string    `bson:"role,omitempty" json:"role,omitempty"`
	Requests []Request ` bson:"requests,omitempty" json:"requests,omitempty"`
}
type Case struct {
	ID           string `bson:"ID,omitempty" json:"id,omitempty"`               // Unique identifier for the case
	Type         string `bson:"type,omitempty" json:"type"`                     // Type of case (e.g., Civil, Criminal)
	Status       string `bson:"status,omitempty" json:"status"`                 // Current status of the case (e.g., Open, Closed)
	FilingDate   string `bson:"filingDate,omitempty" json:"filing_date"`        // Date when the case was filed
	HearingDates string `bson:"hearingDates,omitempty" json:"hearing_dates"`    // List of scheduled hearing dates
	Judge        string `bson:"judge,omitempty" json:"judge,omitempty"`         // Name of the judge
	Plaintiff    string `bson:"plaintiff,omitempty" json:"plaintiff,omitempty"` // Name of the plaintiff
	Defendant    string `bson:"defendant,omitempty" json:"defendant,omitempty"` // Name of the defendant
	Lawyers      string `bson:"lawyers,omitempty" json:"lawyers,omitempty"`     // List of lawyers involved
}
type Request struct {
	ID          string `bson:"id,omitempty" json:"id,omitempty"`         // Unique identifier for the request
	Type        string `bson:"type,omitempty" json:"type,omitempty"`     // Type of request (e.g., access, support)
	Status      string `bson:"status,omitempty" json:"status,omitempty"` // Current status of the request (e.g., pending, resolved)
	Case        string `bson:"case,omitempty" json:"case,omitempty"`
	Description string `bson:"description,omitempty" json:"description,omitempty"` // Description of the request
	CreatedAt   string `bson:"created_at,omitempty" json:"created_at,omitempty"`   // Timestamp when the request was created
}
type GetRequest struct {
	Uuid string `bson:"uuid,omitempty" json:"uuid,omitempty"`
}
type Response struct {
	CaseStatus bool `json:"case_status"`
}
