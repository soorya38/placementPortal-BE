package companyPresenter

type CreateCompany struct {
	CompanyName     string `json:"companyName"`
	CompanyAddress  string `json:"companyAddress"`
	Drive           string `json:"drive"`
	TypeOfDrive     string `json:"typeOfDrive"`
	FollowUp        string `json:"followUp"`
	IsContacted     bool   `json:"isContacted"`
	Remarks         string `json:"remarks"`
	ContactDetails  string `json:"contactDetails"`
	Hr1Details      string `json:"hr1Details"`
	Hr2Details      string `json:"hr2Details"`
	Package         string `json:"package"`
	AssignedOfficer []string `json:"assignedOfficer"`
}
