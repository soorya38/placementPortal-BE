package entity

type Company struct {
	ID              string   `json:"id"`
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
	CreatedAt       string   `json:"createdAt"`
	UpdatedAt       string   `json:"updatedAt"`
}
