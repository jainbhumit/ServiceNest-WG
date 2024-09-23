package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"serviceNest/config"
	"serviceNest/interfaces"
	"serviceNest/model"
	"serviceNest/util"
)

type ServiceRequestRepository struct {
	db *sql.DB
}

// NewServiceRequestRepository initializes a new ServiceRequestRepository with MySQL
func NewServiceRequestRepository(db *sql.DB) interfaces.ServiceRequestRepository {
	return &ServiceRequestRepository{db: db}
}

// SaveServiceRequest saves a service request to the MySQL database
func (repo *ServiceRequestRepository) SaveServiceRequest(request model.ServiceRequest) error {
	column := []string{"id", "householder_id", "householder_name", "householder_address", "service_id", "requested_time", "scheduled_time", "status", "approve_status"}
	query := config.InsertQuery("service_requests", column)
	//query := `
	//	INSERT INTO service_requests
	//	(id, householder_id, householder_name, householder_address, service_id, requested_time, scheduled_time, status, approve_status)
	//	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	//`

	_, err := repo.db.Exec(query, request.ID, request.HouseholderID, request.HouseholderName, request.HouseholderAddress, request.ServiceID, request.RequestedTime, request.ScheduledTime, request.Status, request.ApproveStatus)
	return err
}

// GetServiceRequestByID retrieves a service request by its ID from MySQL
func (repo *ServiceRequestRepository) GetServiceRequestByID(requestID string) (*model.ServiceRequest, error) {
	firstTableColumn := []string{"id", "householder_id", "householder_name", "householder_address", "service_id", "requested_time", "scheduled_time", "status", "approve_status"}
	secondTableColumn := []string{"name"}
	query := config.SelectInnerJoinQuery("service_requests", "services", "service_requests.service_id = services.id", "service_requests.id", firstTableColumn, secondTableColumn)

	//query := `
	//	SELECT sr.id, sr.householder_id, sr.householder_name, sr.householder_address, sr.service_id,
	//	       sr.requested_time, sr.scheduled_time, sr.status, sr.approve_status ,s.name as service_name,
	//	FROM service_requests sr
	//	INNER JOIN services s ON sr.service_id = s.id
	//	WHERE sr.id = ?
	//`

	var request model.ServiceRequest
	var requestedTime []uint8
	var scheduledTime []uint8

	// Execute the query
	err := repo.db.QueryRow(query, requestID).Scan(
		&request.ID, &request.HouseholderID, &request.HouseholderName, &request.HouseholderAddress,
		&request.ServiceID, &requestedTime, &scheduledTime, &request.Status, &request.ApproveStatus, &request.ServiceName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("service request not found")
		}
		return nil, err
	}

	// Parse the times
	request.RequestedTime, err = util.ParseTime(requestedTime)
	if err != nil {
		return nil, fmt.Errorf("error parsing requested_time: %v", err)
	}

	request.ScheduledTime, err = util.ParseTime(scheduledTime)
	if err != nil {
		return nil, fmt.Errorf("error parsing scheduled_time: %v", err)
	}

	return &request, nil
}

func (repo *ServiceRequestRepository) GetServiceRequestsByHouseholderID(householderID string) ([]model.ServiceRequest, error) {
	firstTableColumn := []string{"id", "householder_id", "householder_name", "householder_address", "service_id", "requested_time", "scheduled_time", "status", "approve_status"}
	secondTableColumn := []string{"service_provider_id", "name", "contact", "address", "price", "rating", "approve"}
	query := config.SelectLeftJoinQuery("service_requests", "service_provider_details", "service_requests.id = service_provider_details.service_request_id", "service_requests.householder_id", firstTableColumn, secondTableColumn)

	//query := `
	//	SELECT sr.id, sr.householder_id, sr.householder_name, sr.householder_address, sr.service_id,
	//	       sr.requested_time, sr.scheduled_time, sr.status, sr.approve_status, spd.service_provider_id, spd.name,
	//	       spd.contact, spd.address, spd.price, spd.rating, spd.approve
	//	FROM service_requests AS sr
	//	LEFT JOIN service_provider_details AS spd ON sr.id = spd.service_request_id
	//	WHERE householder_id = ?
	//`

	rows, err := repo.db.Query(query, householderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []model.ServiceRequest
	for rows.Next() {
		var request model.ServiceRequest
		var requestedTime []uint8
		var scheduledTime []uint8

		// Using sql.NullString and sql.NullFloat64 for nullable columns
		var providerID sql.NullString
		var providerName sql.NullString
		var providerContact sql.NullString
		var providerAddress sql.NullString
		var providerPrice sql.NullString
		var providerRating sql.NullFloat64
		var providerApprove sql.NullBool

		// Scan the row data
		err := rows.Scan(
			&request.ID, &request.HouseholderID, &request.HouseholderName, &request.HouseholderAddress,
			&request.ServiceID, &requestedTime, &scheduledTime, &request.Status, &request.ApproveStatus,
			&providerID, &providerName, &providerContact, &providerAddress, &providerPrice, &providerRating,
			&providerApprove,
		)
		if err != nil {
			return nil, err
		}

		// Parse the requested and scheduled times
		request.RequestedTime, err = util.ParseTime(requestedTime)
		if err != nil {
			return nil, fmt.Errorf("error parsing requested_time: %v", err)
		}

		request.ScheduledTime, err = util.ParseTime(scheduledTime)
		if err != nil {
			return nil, fmt.Errorf("error parsing scheduled_time: %v", err)
		}

		// Check if provider details are valid (non-NULL)
		if providerID.Valid {
			provider := model.ServiceProviderDetails{
				ServiceProviderID: providerID.String,
				Name:              providerName.String,
				Contact:           providerContact.String,
				Address:           providerAddress.String,
				Price:             providerPrice.String,
				Rating:            providerRating.Float64,
				Approve:           providerApprove.Bool,
			}
			request.ProviderDetails = append(request.ProviderDetails, provider)
		}

		// Append the request to the slice
		requests = append(requests, request)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}

// UpdateServiceRequest updates an existing service request in MySQL
func (repo *ServiceRequestRepository) UpdateServiceRequest(updatedRequest *model.ServiceRequest) error {
	query := `
		UPDATE service_requests 
		SET householder_id = ?, householder_name = ?, householder_address = ?, service_id = ?, requested_time = ?, scheduled_time = ?, status = ?, approve_status = ? 
		WHERE id = ?
	`

	_, err := repo.db.Exec(query, updatedRequest.HouseholderID, updatedRequest.HouseholderName, updatedRequest.HouseholderAddress, updatedRequest.ServiceID, updatedRequest.RequestedTime, updatedRequest.ScheduledTime, updatedRequest.Status, updatedRequest.ApproveStatus, updatedRequest.ID)
	return err
}

// GetAllServiceRequests retrieves all service requests from MySQL
func (repo *ServiceRequestRepository) GetAllServiceRequests() ([]model.ServiceRequest, error) {
	firstTableColumn := []string{"id", "householder_id", "householder_name", "householder_address", "service_id", "requested_time", "scheduled_time", "status", "approve_status"}
	secondTableColumn := []string{"service_provider_id", "name", "contact", "address", "price", "rating", "approve"}
	query := config.SelectLeftJoinQuery("service_requests", "service_provider_details", "service_requests.id = service_provider_details.service_request_id", "", firstTableColumn, secondTableColumn)

	//query := `
	//	SELECT sr.id, sr.householder_id, sr.householder_name, sr.householder_address, sr.service_id, sr.requested_time, sr.scheduled_time, sr.status, sr.approve_status,
	//	       spd.service_provider_id, spd.name, spd.contact, spd.address, spd.price, spd.rating, spd.approve
	//	FROM service_requests AS sr
	//	LEFT JOIN service_provider_details AS spd ON sr.id = spd.service_request_id
	//`

	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []model.ServiceRequest
	for rows.Next() {
		var requestedTime, scheduledTime []uint8
		var request model.ServiceRequest
		var provider model.ServiceProviderDetails

		// Use sql.NullString and other nullable types for fields that may contain NULLs
		var providerID, providerName, providerContact, providerAddress, providerPrice sql.NullString
		var providerRating sql.NullFloat64
		var providerApprove sql.NullBool

		err := rows.Scan(
			&request.ID, &request.HouseholderID, &request.HouseholderName, &request.HouseholderAddress,
			&request.ServiceID, &requestedTime, &scheduledTime, &request.Status, &request.ApproveStatus,
			&providerID, &providerName, &providerContact, &providerAddress, &providerPrice,
			&providerRating, &providerApprove,
		)
		if err != nil {
			return nil, err
		}

		// Parse the requested_time
		request.RequestedTime, err = util.ParseTime(requestedTime)
		if err != nil {
			return nil, fmt.Errorf("error parsing requested_time: %v", err)
		}

		// Parse the scheduled_time
		request.ScheduledTime, err = util.ParseTime(scheduledTime)
		if err != nil {
			return nil, fmt.Errorf("error parsing scheduled_time: %v", err)
		}

		// Assign values to provider struct only if they are not NULL
		if providerID.Valid {
			provider.ServiceProviderID = providerID.String
		}
		if providerName.Valid {
			provider.Name = providerName.String
		}
		if providerContact.Valid {
			provider.Contact = providerContact.String
		}
		if providerAddress.Valid {
			provider.Address = providerAddress.String
		}
		if providerPrice.Valid {
			provider.Price = providerPrice.String
		}
		if providerRating.Valid {
			provider.Rating = providerRating.Float64
		}
		if providerApprove.Valid {
			provider.Approve = providerApprove.Bool
		}

		// Append provider details if a valid provider is found
		if providerID.Valid {
			request.ProviderDetails = append(request.ProviderDetails, provider)
		}

		requests = append(requests, request)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}

// GetServiceRequestsByProviderID retrieves service requests by the provider ID from MySQL
func (repo *ServiceRequestRepository) GetServiceRequestsByProviderID(providerID string) ([]model.ServiceRequest, error) {
	firstTableColumn := []string{"id", "householder_id", "householder_name", "householder_address", "service_id", "requested_time", "scheduled_time", "status", "approve_status"}
	secondTableColumn := []string{"service_provider_id", "name", "contact", "address", "price", "rating", "approve"}
	query := config.SelectInnerJoinQuery("service_requests", "service_provider_details", "service_requests.id = service_provider_details.service_request_id", "service_provider_details.service_provider_id", firstTableColumn, secondTableColumn)

	//	query := `
	//	SELECT sr.id, sr.householder_id, sr.householder_name, sr.householder_address, sr.service_id, sr.requested_time, sr.scheduled_time, sr.status, sr.approve_status
	//,spd.service_provider_id,spd.name,spd.contact,spd.address
	//,spd.price,spd.rating,spd.approve FROM service_requests as sr inner join service_provider_details as spd  on sr.id=spd.service_request_id
	//		WHERE spd.service_provider_id=?;
	//	`

	rows, err := repo.db.Query(query, providerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []model.ServiceRequest
	var requestedTime []uint8
	var scheduledTime []uint8
	for rows.Next() {
		var request model.ServiceRequest
		var provider model.ServiceProviderDetails
		err := rows.Scan(
			&request.ID, &request.HouseholderID, &request.HouseholderName, &request.HouseholderAddress,
			&request.ServiceID, &requestedTime, &scheduledTime, &request.Status, &request.ApproveStatus,
			&provider.ServiceProviderID, &provider.Name, &provider.Contact, &provider.Address, &provider.Price,
			&provider.Rating, &provider.Approve,
		)
		if err != nil {
			return nil, err
		}
		// Parse the requested_time
		request.RequestedTime, err = util.ParseTime(requestedTime)
		if err != nil {
			return nil, fmt.Errorf("error parsing requested_time: %v", err)
		}

		// Parse the scheduled_time
		request.ScheduledTime, err = util.ParseTime(scheduledTime)
		if err != nil {
			return nil, fmt.Errorf("error parsing scheduled_time: %v", err)
		}
		request.ProviderDetails = append(request.ProviderDetails, provider)
		requests = append(requests, request)
	}

	return requests, nil
}

func (repo *ServiceRequestRepository) GetServiceProviderByRequestID(requestID, providerID string) (*model.ServiceRequest, error) {
	firstTableColumn := []string{"id", "householder_id", "householder_name", "householder_address", "service_id", "requested_time", "scheduled_time", "status", "approve_status"}
	secondTableColumn := []string{"service_provider_id", "name", "contact", "address", "price", "rating", "approve"}
	query := config.SelectInnerJoinQuery("service_requests", "service_provider_details", "service_requests.id = service_provider_details.service_request_id", "service_provider_details.service_provider_id = ? AND service_requests", firstTableColumn, secondTableColumn)

	//query := `SELECT sr.id, sr.householder_id, sr.householder_name, sr.householder_address, sr.service_id, sr.requested_time, sr.scheduled_time, sr.status, sr.approve_status,
	//spd.service_provider_id, spd.name, spd.contact, spd.address, spd.price, spd.rating, spd.approve
	//FROM service_requests AS sr
	//INNER JOIN service_provider_details AS spd ON sr.id = spd.service_request_id
	//WHERE spd.service_provider_id = ? AND sr.id = ?`

	rows, err := repo.db.Query(query, providerID, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Make sure to call Next() before scanning
	if rows.Next() {
		var request model.ServiceRequest
		var requestedTime []uint8
		var scheduledTime []uint8
		var provider model.ServiceProviderDetails

		// Scan the row data
		err = rows.Scan(
			&request.ID, &request.HouseholderID, &request.HouseholderName, &request.HouseholderAddress,
			&request.ServiceID, &requestedTime, &scheduledTime, &request.Status, &request.ApproveStatus,
			&provider.ServiceProviderID, &provider.Name, &provider.Contact, &provider.Address, &provider.Price,
			&provider.Rating, &provider.Approve,
		)
		if err != nil {
			return nil, err
		}

		// Parse the requested_time
		request.RequestedTime, err = util.ParseTime(requestedTime)
		if err != nil {
			return nil, fmt.Errorf("error parsing requested_time: %v", err)
		}

		// Parse the scheduled_time
		request.ScheduledTime, err = util.ParseTime(scheduledTime)
		if err != nil {
			return nil, fmt.Errorf("error parsing scheduled_time: %v", err)
		}

		// Append the provider details
		request.ProviderDetails = append(request.ProviderDetails, provider)

		return &request, nil
	}
	// If no rows found, return an error
	return nil, fmt.Errorf("no service request found for request ID: %s and provider ID: %s", requestID, providerID)
}
