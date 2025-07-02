package company

import (
	"backend/companyd/entity"
)

type Repository interface {
	CreateCompany(companyName, companyAddress, drive, typeOfDrive, followUp, isContacted, remarks, contactDetails, hr1Details, hr2Details, pkg string, assignedOfficer []string) (*entity.Company, error)
	ListCompanies() ([]*entity.Company, error)
	DeleteCompany(id string) error
	UpdateCompany(id, companyName, companyAddress, drive, typeOfDrive, followUp, isContacted, remarks, contactDetails, hr1Details, hr2Details, pkg string, assignedOfficer []string) (*entity.Company, error)
	ListCompaniesByUsername(username string) ([]*entity.Company, error)
	CreateCompanyTemp(companyId, companyName, companyAddress, drive, typeOfDrive, followUp, isContacted, remarks, contactDetails, hr1Details, hr2Details, pkg string, assignedOfficer []string, createdBy string) (*entity.CompanyTemp, error)
	ListCompanyTemps() ([]*entity.CompanyTemp, error)
	UpdateCompanyTempStatus(id string, status string) error
	ApproveCompanyTemp(id string) error
	CreateEvent(date, eventType, title, description, createdBy string) (*entity.Event, error)
	ListEvents() ([]*entity.Event, error)
}

type Writer interface {
	CreateCompany(companyName string,
		companyAddress string,
		drive string,
		typeOfDrive string,
		followUp string,
		isContacted string,
		remarks string,
		contactDetails string,
		hr1Details string,
		hr2Details string,
		pkg string,
		assignedOfficer []string,
	) (*entity.Company, error)
	DeleteCompany(id string) error
	ApproveCompanyTemp(id string, status string) error
	UpdateCompany(id, companyName, companyAddress, drive, typeOfDrive, followUp, isContacted, remarks, contactDetails, hr1Details, hr2Details, pkg string, assignedOfficer []string) (*entity.Company, error)
	CreateEvent(date, eventType, title, description, createdBy string) (*entity.Event, error)
}

type Reader interface {
	ListCompanies() ([]*entity.Company, error)
	ListCompaniesByUsername(username string) ([]*entity.Company, error)
	ListEvents() ([]*entity.Event, error)
}

type Usecase interface {
	CreateCompany(companyName string,
		companyAddress string,
		drive string,
		typeOfDrive string,
		followUp string,
		isContacted string,
		remarks string,
		contactDetails string,
		hr1Details string,
		hr2Details string,
		pkg string,
		assignedOfficer []string,
	) (*entity.Company, error)
	ListCompanies() ([]*entity.Company, error)
	DeleteCompany(id string) error
	UpdateCompany(id, companyName, companyAddress, drive, typeOfDrive, followUp, isContacted, remarks, contactDetails, hr1Details, hr2Details, pkg string, assignedOfficer []string) (*entity.Company, error)
	ListCompaniesByUsername(username string) ([]*entity.Company, error)
	CreateCompanyTemp(companyId, companyName, companyAddress, drive, typeOfDrive, followUp, isContacted, remarks, contactDetails, hr1Details, hr2Details, pkg string, assignedOfficer []string, createdBy string) (*entity.CompanyTemp, error)
	ListCompanyTemps() ([]*entity.CompanyTemp, error)
	UpdateCompanyTempStatus(id string, status string) error
	ApproveCompanyTemp(id string) error
	CreateEvent(date, eventType, title, description, createdBy string) (*entity.Event, error)
	ListEvents() ([]*entity.Event, error)
}
