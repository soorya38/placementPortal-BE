package companyPresenter

type CreateCompanyTemp struct {
	CompanyID       string   `json:"company_id"`
	CompanyName     string   `json:"company_name"`
	CompanyAddress  string   `json:"company_address"`
	Drive           string   `json:"drive"`
	TypeOfDrive     string   `json:"type_of_drive"`
	FollowUp        string   `json:"follow_up"`
	IsContacted     bool     `json:"is_contacted"`
	Remarks         string   `json:"remarks"`
	ContactDetails  string   `json:"contact_details"`
	Hr1Details      string   `json:"hr1_details"`
	Hr2Details      string   `json:"hr2_details"`
	Package         string   `json:"package"`
	AssignedOfficer []string `json:"assigned_officer"`
	CreatedBy       string   `json:"created_by"`
}

type CompanyTempResponse struct {
	ID              string   `json:"id"`
	CompanyID       string   `json:"company_id"`
	CompanyName     string   `json:"company_name"`
	CompanyAddress  string   `json:"company_address"`
	Drive           string   `json:"drive"`
	TypeOfDrive     string   `json:"type_of_drive"`
	FollowUp        string   `json:"follow_up"`
	IsContacted     bool     `json:"is_contacted"`
	Remarks         string   `json:"remarks"`
	ContactDetails  string   `json:"contact_details"`
	Hr1Details      string   `json:"hr1_details"`
	Hr2Details      string   `json:"hr2_details"`
	Package         string   `json:"package"`
	AssignedOfficer []string `json:"assigned_officer"`
	CreatedBy       string   `json:"created_by"`
	Status          string   `json:"status"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
}
