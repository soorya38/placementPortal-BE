package entity

type CompanyTemp struct {
	ID              string   `json:"id"`
	CompanyID       string   `json:"companyId"`
	CompanyName     string   `json:"companyName"`
	CompanyAddress  string   `json:"companyAddress"`
	Drive           string   `json:"drive"`
	TypeOfDrive     string   `json:"typeOfDrive"`
	FollowUp        string   `json:"followUp"`
	IsContacted     bool     `json:"isContacted"`
	Remarks         string   `json:"remarks"`
	ContactDetails  string   `json:"contactDetails"`
	HR1Details      string   `json:"hr1Details"`
	HR2Details      string   `json:"hr2Details"`
	Package         string   `json:"package"`
	AssignedOfficer []string `json:"assignedOfficer"`
	Status          string   `json:"status"`
	CreatedBy       string   `json:"createdBy"`
	CreatedAt       string   `json:"createdAt"`
	UpdatedAt       string   `json:"updatedAt"`
}

