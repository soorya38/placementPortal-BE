package repository

import (
	"backend/companyd/entity"
	"database/sql"

	"github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewCompanyRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateCompany(companyName, companyAddress, drive, typeOfDrive, followUp, isContacted, remarks, contactDetails, hr1Details, hr2Details, pkg string, assignedOfficer []string) (*entity.Company, error) {
	query := `
		INSERT INTO companies (company_name, company_address, drive, type_of_drive, follow_up, is_contacted, remarks, contact_details, hr1_details, hr2_details, package, assigned_officer) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, company_name, company_address, drive, type_of_drive, follow_up, is_contacted, remarks, contact_details, hr1_details, hr2_details, package, assigned_officer, created_at, updated_at`

	var company entity.Company
	var assignedOfficerResult []string
	err := r.db.QueryRow(query, companyName, companyAddress, drive, typeOfDrive, followUp, isContacted, remarks, contactDetails, hr1Details, hr2Details, pkg, pq.Array(assignedOfficer)).Scan(
		&company.ID, &company.CompanyName, &company.CompanyAddress, &company.Drive, &company.TypeOfDrive, &company.FollowUp, &company.IsContacted, &company.Remarks, &company.ContactDetails, &company.HR1Details, &company.HR2Details, &company.Package, pq.Array(&assignedOfficerResult), &company.CreatedAt, &company.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	company.AssignedOfficer = assignedOfficerResult
	return &company, nil
}

func (r *Repository) DeleteCompany(id string) error {
	query := `DELETE FROM companies WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *Repository) ListCompanies() ([]*entity.Company, error) {
	query := `
		SELECT id, company_name, company_address, drive, type_of_drive, follow_up, is_contacted, remarks, contact_details, hr1_details, hr2_details, package, assigned_officer, created_at, updated_at 
		FROM companies`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []*entity.Company
	for rows.Next() {
		var company entity.Company
		var assignedOfficer []string
		err := rows.Scan(
			&company.ID, &company.CompanyName, &company.CompanyAddress, &company.Drive, &company.TypeOfDrive, &company.FollowUp, &company.IsContacted, &company.Remarks, &company.ContactDetails, &company.HR1Details, &company.HR2Details, &company.Package, pq.Array(&assignedOfficer), &company.CreatedAt, &company.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		company.AssignedOfficer = assignedOfficer
		companies = append(companies, &company)
	}
	return companies, nil
}

func (r *Repository) UpdateCompany(id, companyName, companyAddress, drive, typeOfDrive, followUp, isContacted, remarks, contactDetails, hr1Details, hr2Details, pkg string, assignedOfficer []string) (*entity.Company, error) {
	query := `
		UPDATE companies 
		SET company_name = $1, 
			company_address = $2, 
			drive = $3, 
			type_of_drive = $4, 
			follow_up = $5, 
			is_contacted = $6, 
			remarks = $7, 
			contact_details = $8, 
			hr1_details = $9, 
			hr2_details = $10, 
			package = $11, 
			assigned_officer = $12,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $13
		RETURNING id, company_name, company_address, drive, type_of_drive, follow_up, is_contacted, remarks, contact_details, hr1_details, hr2_details, package, assigned_officer, created_at, updated_at`

	var company entity.Company
	var assignedOfficerResult []string
	err := r.db.QueryRow(query,
		companyName, companyAddress, drive, typeOfDrive, followUp, isContacted, remarks,
		contactDetails, hr1Details, hr2Details, pkg, pq.Array(assignedOfficer), id).Scan(
		&company.ID, &company.CompanyName, &company.CompanyAddress, &company.Drive,
		&company.TypeOfDrive, &company.FollowUp, &company.IsContacted, &company.Remarks,
		&company.ContactDetails, &company.HR1Details, &company.HR2Details, &company.Package,
		pq.Array(&assignedOfficerResult), &company.CreatedAt, &company.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	company.AssignedOfficer = assignedOfficerResult
	return &company, nil
}

func (r *Repository) ListCompaniesByUsername(username string) ([]*entity.Company, error) {
	query := `
		SELECT id, company_name, company_address, drive, type_of_drive, follow_up, is_contacted, remarks, contact_details, hr1_details, hr2_details, package, assigned_officer, created_at, updated_at 
		FROM companies 
		WHERE $1 = ANY(assigned_officer)`

	rows, err := r.db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []*entity.Company
	for rows.Next() {
		var company entity.Company
		var assignedOfficer []string
		err := rows.Scan(
			&company.ID, &company.CompanyName, &company.CompanyAddress, &company.Drive, &company.TypeOfDrive, &company.FollowUp, &company.IsContacted, &company.Remarks, &company.ContactDetails, &company.HR1Details, &company.HR2Details, &company.Package, pq.Array(&assignedOfficer), &company.CreatedAt, &company.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		company.AssignedOfficer = assignedOfficer
		companies = append(companies, &company)
	}
	return companies, nil
}

func (r *Repository) CreateCompanyTemp(companyId, companyName, companyAddress, drive, typeOfDrive, followUp, isContacted, remarks, contactDetails, hr1Details, hr2Details, pkg string, assignedOfficer []string, createdBy string) (*entity.CompanyTemp, error) {
	query := `
		INSERT INTO companies_temp (company_id, company_name, company_address, drive, type_of_drive, follow_up, is_contacted, remarks, contact_details, hr1_details, hr2_details, package, assigned_officer, created_by) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, company_id, company_name, company_address, drive, type_of_drive, follow_up, is_contacted, remarks, contact_details, hr1_details, hr2_details, package, assigned_officer, status, created_by, created_at, updated_at`

	var companyTemp entity.CompanyTemp
	var assignedOfficerResult []string
	err := r.db.QueryRow(query, companyId, companyName, companyAddress, drive, typeOfDrive, followUp, isContacted, remarks, contactDetails, hr1Details, hr2Details, pkg, pq.Array(assignedOfficer), createdBy).Scan(
		&companyTemp.ID, &companyTemp.CompanyID, &companyTemp.CompanyName, &companyTemp.CompanyAddress, &companyTemp.Drive, &companyTemp.TypeOfDrive, &companyTemp.FollowUp, &companyTemp.IsContacted, &companyTemp.Remarks, &companyTemp.ContactDetails, &companyTemp.HR1Details, &companyTemp.HR2Details, &companyTemp.Package, pq.Array(&assignedOfficerResult), &companyTemp.Status, &companyTemp.CreatedBy, &companyTemp.CreatedAt, &companyTemp.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	companyTemp.AssignedOfficer = assignedOfficerResult
	return &companyTemp, nil
}

func (r *Repository) ListCompanyTemps() ([]*entity.CompanyTemp, error) {
	query := `
		SELECT id, company_id, company_name, company_address, drive, type_of_drive, follow_up, is_contacted, remarks, contact_details, hr1_details, hr2_details, package, assigned_officer, status, created_by, created_at, updated_at 
		FROM companies_temp
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companyTemps []*entity.CompanyTemp
	for rows.Next() {
		var companyTemp entity.CompanyTemp
		var assignedOfficer []string
		err := rows.Scan(
			&companyTemp.ID, &companyTemp.CompanyID, &companyTemp.CompanyName, &companyTemp.CompanyAddress, &companyTemp.Drive, &companyTemp.TypeOfDrive, &companyTemp.FollowUp, &companyTemp.IsContacted, &companyTemp.Remarks, &companyTemp.ContactDetails, &companyTemp.HR1Details, &companyTemp.HR2Details, &companyTemp.Package, pq.Array(&assignedOfficer), &companyTemp.Status, &companyTemp.CreatedBy, &companyTemp.CreatedAt, &companyTemp.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		companyTemp.AssignedOfficer = assignedOfficer
		companyTemps = append(companyTemps, &companyTemp)
	}
	return companyTemps, nil
}

func (r *Repository) UpdateCompanyTempStatus(id string, status string) error {
	query := `UPDATE companies_temp SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.Exec(query, status, id)
	return err
}

func (r *Repository) ApproveCompanyTemp(id string) error {
	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get the company temp data
	var companyTemp entity.CompanyTemp
	var assignedOfficer []string
	err = tx.QueryRow(`
		SELECT company_id, company_name, company_address, drive, type_of_drive, follow_up, is_contacted, remarks, contact_details, hr1_details, hr2_details, package, assigned_officer 
		FROM companies_temp 
		WHERE id = $1`, id).Scan(
		&companyTemp.CompanyID, &companyTemp.CompanyName, &companyTemp.CompanyAddress, &companyTemp.Drive,
		&companyTemp.TypeOfDrive, &companyTemp.FollowUp, &companyTemp.IsContacted, &companyTemp.Remarks,
		&companyTemp.ContactDetails, &companyTemp.HR1Details, &companyTemp.HR2Details, &companyTemp.Package,
		pq.Array(&assignedOfficer))
	if err != nil {
		return err
	}
	companyTemp.AssignedOfficer = assignedOfficer

	// Update the company with the temp data
	_, err = tx.Exec(`
		UPDATE companies 
		SET company_name = $1,
			company_address = $2,
			drive = $3,
			type_of_drive = $4,
			follow_up = $5,
			is_contacted = $6,
			remarks = $7,
			contact_details = $8,
			hr1_details = $9,
			hr2_details = $10,
			package = $11,
			assigned_officer = $12,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $13`,
		companyTemp.CompanyName, companyTemp.CompanyAddress, companyTemp.Drive,
		companyTemp.TypeOfDrive, companyTemp.FollowUp, companyTemp.IsContacted,
		companyTemp.Remarks, companyTemp.ContactDetails, companyTemp.HR1Details,
		companyTemp.HR2Details, companyTemp.Package, pq.Array(companyTemp.AssignedOfficer),
		companyTemp.CompanyID)
	if err != nil {
		return err
	}

	// Delete the temp record
	_, err = tx.Exec("DELETE FROM companies_temp WHERE id = $1", id)
	if err != nil {
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

func (r *Repository) CreateEvent(date, eventType, title, description, createdBy string) (*entity.Event, error) {
	var event entity.Event

	err := r.db.QueryRow(`
		INSERT INTO events (id, date, type, title, description, created_by, created_at)
		VALUES (uuid_generate_v4(), $1, $2, $3, $4, $5, CURRENT_TIMESTAMP)
		RETURNING id, date, type, title, description, created_by, created_at`,
		date, eventType, title, description, createdBy,
	).Scan(
		&event.ID,
		&event.Date,
		&event.Type,
		&event.Title,
		&event.Description,
		&event.CreatedBy,
		&event.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *Repository) ListEvents() ([]*entity.Event, error) {
	query := `
		SELECT id, date, type, title, description, created_by, created_at 
		FROM events
		ORDER BY date DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*entity.Event
	for rows.Next() {
		var event entity.Event
		err := rows.Scan(
			&event.ID,
			&event.Date,
			&event.Type,
			&event.Title,
			&event.Description,
			&event.CreatedBy,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	return events, nil
}
