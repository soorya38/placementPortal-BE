package company

import (
	"backend/companyd/entity"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) Usecase {
	return &Service{repo: repo}
}

func (s *Service) CreateCompany(companyName,
	companyAddress,
	drive,
	typeOfDrive,
	followUp,
	isContacted,
	remarks,
	contactDetails,
	hr1Details,
	hr2Details,
	pkg string,
	assignedOfficer []string,
) (*entity.Company, error) {
	company, err := s.repo.CreateCompany(companyName,
		companyAddress,
		drive,
		typeOfDrive,
		followUp,
		isContacted,
		remarks,
		contactDetails,
		hr1Details,
		hr2Details,
		pkg,
		assignedOfficer)
	if err != nil {
		return nil, err
	}

	return company, nil
}

func (s *Service) DeleteCompany(id string) error {
	return s.repo.DeleteCompany(id)
}

func (s *Service) ListCompanies() ([]*entity.Company, error) {
	return s.repo.ListCompanies()
}

func (s *Service) UpdateCompany(id, companyName, companyAddress, drive, typeOfDrive, followUp, isContacted, remarks, contactDetails, hr1Details, hr2Details, pkg string, assignedOfficer []string) (*entity.Company, error) {
	company, err := s.repo.UpdateCompany(id, companyName, companyAddress, drive, typeOfDrive, followUp, isContacted, remarks, contactDetails, hr1Details, hr2Details, pkg, assignedOfficer)
	if err != nil {
		return nil, err
	}
	return company, nil
}

func (s *Service) ListCompaniesByUsername(username string) ([]*entity.Company, error) {
	return s.repo.ListCompaniesByUsername(username)
}

func (s *Service) CreateCompanyTemp(companyId, companyName, companyAddress, drive, typeOfDrive, followUp, isContacted, remarks, contactDetails, hr1Details, hr2Details, pkg string, assignedOfficer []string, createdBy string) (*entity.CompanyTemp, error) {
	return s.repo.CreateCompanyTemp(companyId, companyName, companyAddress, drive, typeOfDrive, followUp, isContacted, remarks, contactDetails, hr1Details, hr2Details, pkg, assignedOfficer, createdBy)
}

func (s *Service) ListCompanyTemps() ([]*entity.CompanyTemp, error) {
	return s.repo.ListCompanyTemps()
}

func (s *Service) UpdateCompanyTempStatus(id string, status string) error {
	return s.repo.UpdateCompanyTempStatus(id, status)
}

func (s *Service) ApproveCompanyTemp(id string) error {
	return s.repo.ApproveCompanyTemp(id)
}

func (s *Service) CreateEvent(date, eventType, title, description, createdBy string) (*entity.Event, error) {
	return s.repo.CreateEvent(date, eventType, title, description, createdBy)
}

func (s *Service) ListEvents() ([]*entity.Event, error) {
	return s.repo.ListEvents()
}
