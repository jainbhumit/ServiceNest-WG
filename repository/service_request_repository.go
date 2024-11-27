package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"serviceNest/config"
	"serviceNest/errs"
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
	column := []string{"id", "householder_id", "householder_name", "householder_address", "householder_contact", "service_id", "requested_time", "scheduled_time", "status", "approve_status", "service_name", "description"}
	query := config.InsertQuery("service_requests", column)

	_, err := repo.db.Exec(query, request.ID, request.HouseholderID, request.HouseholderName, request.HouseholderAddress, request.HouseholderContact, request.ServiceID, request.RequestedTime.Format("2006-01-02 15:04:05"), request.ScheduledTime, request.Status, request.ApproveStatus, request.ServiceName, request.Description)
	return err
}

// GetServiceRequestByID retrieves a service request by its ID from MySQL
func (repo *ServiceRequestRepository) GetServiceRequestByID(requestID string) (*model.ServiceRequest, error) {
	firstTableColumn := []string{"id", "householder_id", "householder_name", "householder_address", "service_id", "requested_time", "scheduled_time", "status", "approve_status"}
	secondTableColumn := []string{"name"}
	query := config.SelectInnerJoinQuery("service_requests", "services", "service_requests.service_id = services.id", "service_requests.id", firstTableColumn, secondTableColumn)

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
			return nil, errors.New(errs.ServiceRequestNotFound)
		}
		return nil, err
	}

	// Parse the times
	request.RequestedTime, err = util.ParseTime(requestedTime)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", errs.ErrorParsingRequestTime, err)
	}

	request.ScheduledTime, err = util.ParseTime(scheduledTime)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", errs.ErrorParsingScheduleTime, err)
	}

	return &request, nil
}

func (repo *ServiceRequestRepository) GetServiceRequestsByHouseholderID(householderID string, limit, offset int, status string) ([]model.ServiceRequest, error) {
	query := config.SelectJsonDataQuery(status != "")

	var rows *sql.Rows
	var err error
	if status != "" {
		rows, err = repo.db.Query(query, householderID, status, limit, offset)
	} else {
		rows, err = repo.db.Query(query, householderID, limit, offset)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []model.ServiceRequest
	for rows.Next() {
		var request model.ServiceRequest
		var requestedTime, scheduledTime []uint8
		var providerDetailsJSON sql.NullString
		var householderID sql.NullString
		var householderAddress sql.NullString
		var serviceName sql.NullString

		err := rows.Scan(
			&request.ID, &householderID, &request.HouseholderName, &householderAddress,
			&request.ServiceID, &requestedTime, &scheduledTime, &request.Status, &request.ApproveStatus, &serviceName,
			&providerDetailsJSON,
		)
		if err != nil {
			return nil, err
		}

		// Assign nullable values
		if householderID.Valid {
			request.HouseholderID = &householderID.String
		}
		if householderAddress.Valid {
			request.HouseholderAddress = &householderAddress.String
		}
		if serviceName.Valid {
			request.ServiceName = serviceName.String
		}

		// Parse the requested and scheduled times
		request.RequestedTime, err = util.ParseTime(requestedTime)
		if err != nil {
			return nil, fmt.Errorf("error parsing requested time: %v", err)
		}
		request.ScheduledTime, err = util.ParseTime(scheduledTime)
		if err != nil {
			return nil, fmt.Errorf("error parsing scheduled time: %v", err)
		}

		// Parse the JSON data for provider details
		if providerDetailsJSON.Valid {
			var providerDetails []model.ServiceProviderDetails
			err = json.Unmarshal([]byte(providerDetailsJSON.String), &providerDetails)
			if err != nil {
				return nil, fmt.Errorf("error parsing provider details JSON: %v", err)
			}
			request.ProviderDetails = providerDetails
		}

		requests = append(requests, request)
	}

	// Check for any errors encountered during iteration
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
func (repo *ServiceRequestRepository) GetAllServiceRequests(limit, offset int) ([]model.ServiceRequest, error) {
	firstTableColumn := []string{"id", "householder_id", "householder_name", "householder_address", "service_id", "requested_time", "scheduled_time", "status", "approve_status", "service_name"}
	secondTableColumn := []string{"service_provider_id", "name", "contact", "address", "price", "rating", "approve"}
	query := config.SelectLeftJoinQuery("service_requests", "service_provider_details", "service_requests.id = service_provider_details.service_request_id", "", firstTableColumn, secondTableColumn, limit, offset)

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
		var ServiceName sql.NullString

		err := rows.Scan(
			&request.ID, &request.HouseholderID, &request.HouseholderName, &request.HouseholderAddress,
			&request.ServiceID, &requestedTime, &scheduledTime, &request.Status, &request.ApproveStatus, &ServiceName,
			&providerID, &providerName, &providerContact, &providerAddress, &providerPrice,
			&providerRating, &providerApprove,
		)
		if err != nil {
			return nil, err
		}

		// Parse the requested_time
		request.RequestedTime, err = util.ParseTime(requestedTime)
		if err != nil {
			return nil, fmt.Errorf("%v: %v", errs.ErrorParsingRequestTime, err)
		}

		// Parse the scheduled_time
		request.ScheduledTime, err = util.ParseTime(scheduledTime)
		if err != nil {
			return nil, fmt.Errorf("%v: %v", errs.ErrorParsingScheduleTime, err)
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
			if providerApprove.Bool == true {
				provider.Approve = 1
			} else {
				provider.Approve = 0
			}
		}

		// Append provider details if a valid provider is found
		if providerID.Valid {
			request.ProviderDetails = append(request.ProviderDetails, provider)
		}
		if ServiceName.Valid {
			request.ServiceName = ServiceName.String
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
func (repo *ServiceRequestRepository) GetServiceRequestsByProviderID(providerID string, limit, offset int, sortOrder string) ([]model.ServiceRequest, error) {
	firstTableColumn := []string{"id", "householder_id", "householder_name", "householder_address", "service_id", "requested_time", "scheduled_time", "status", "approve_status", "householder_contact", "service_name"}
	secondTableColumn := []string{}
	query := config.SelectInnerJoinQueryPaginate("service_requests", "service_provider_details", "service_requests.id = service_provider_details.service_request_id AND service_provider_details.service_provider_id = ?", "service_provider_details.approve", firstTableColumn, secondTableColumn, limit, offset, "service_requests.scheduled_time", sortOrder)
	fmt.Println(query)
	rows, err := repo.db.Query(query, providerID, true)
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
			&request.HouseholderContact, &request.ServiceName,
		)
		if err != nil {
			return nil, err
		}
		// Parse the requested_time
		request.RequestedTime, err = util.ParseTime(requestedTime)
		if err != nil {
			return nil, fmt.Errorf("%v: %v", errs.ErrorParsingRequestTime, err)
		}

		// Parse the scheduled_time
		request.ScheduledTime, err = util.ParseTime(scheduledTime)
		if err != nil {
			return nil, fmt.Errorf("%v: %v", errs.ErrorParsingScheduleTime, err)
		}
		request.ProviderDetails = append(request.ProviderDetails, provider)
		requests = append(requests, request)
	}

	return requests, nil
}

func (repo *ServiceRequestRepository) GetServiceProviderByRequestID(requestID, providerID string) (*model.ServiceRequest, error) {
	firstTableColumn := []string{"id", "householder_id", "householder_name", "householder_address", "service_id", "requested_time", "scheduled_time", "status", "approve_status"}
	secondTableColumn := []string{"service_provider_id", "name", "contact", "address", "price", "rating", "approve"}
	query := config.SelectInnerJoinQuery("service_requests", "service_provider_details", "service_requests.id = service_provider_details.service_request_id", "service_provider_details.service_provider_id = ? AND service_requests.id", firstTableColumn, secondTableColumn)

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
			return nil, fmt.Errorf("%v: %v", errs.ErrorParsingRequestTime, err)
		}

		// Parse the scheduled_time
		request.ScheduledTime, err = util.ParseTime(scheduledTime)
		if err != nil {
			return nil, fmt.Errorf("%v: %v", errs.ErrorParsingScheduleTime, err)
		}

		// Append the provider details
		request.ProviderDetails = append(request.ProviderDetails, provider)

		return &request, nil
	}
	// If no rows found, return an error
	return nil, fmt.Errorf(errs.NoServiceProviderFoundForRequestId)
}

func (repo *ServiceRequestRepository) GetApproveServiceRequestsByHouseholderID(householderID string, limit, offset int, sortOrder string) ([]model.ServiceRequest, error) {
	query := config.SelectJsonDataQueryWithApprove(sortOrder)

	rows, err := repo.db.Query(query, householderID, true, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []model.ServiceRequest
	for rows.Next() {
		var request model.ServiceRequest
		var requestedTime, scheduledTime []uint8
		var providerDetailsJSON sql.NullString
		var householderID sql.NullString
		var householderAddress sql.NullString
		var serviceName sql.NullString

		err := rows.Scan(
			&request.ID, &householderID, &request.HouseholderName, &householderAddress,
			&request.ServiceID, &requestedTime, &scheduledTime, &request.Status, &request.ApproveStatus, &serviceName,
			&providerDetailsJSON,
		)
		if err != nil {
			return nil, err
		}

		// Assign nullable values
		if householderID.Valid {
			request.HouseholderID = &householderID.String
		}
		if householderAddress.Valid {
			request.HouseholderAddress = &householderAddress.String
		}
		if serviceName.Valid {
			request.ServiceName = serviceName.String
		}

		// Parse the requested and scheduled times
		request.RequestedTime, err = util.ParseTime(requestedTime)
		if err != nil {
			return nil, fmt.Errorf("error parsing requested time: %v", err)
		}
		request.ScheduledTime, err = util.ParseTime(scheduledTime)
		if err != nil {
			return nil, fmt.Errorf("error parsing scheduled time: %v", err)
		}

		// Parse the JSON data for provider details
		if providerDetailsJSON.Valid {
			var providerDetails []model.ServiceProviderDetails
			err = json.Unmarshal([]byte(providerDetailsJSON.String), &providerDetails)
			if err != nil {
				return nil, fmt.Errorf("error parsing provider details JSON: %v", err)
			}
			request.ProviderDetails = providerDetails
		}

		requests = append(requests, request)
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}

func (repo *ServiceRequestRepository) GetAllPendingRequestsByProvider(providerId string, serviceID string, limit, offset int) ([]model.ServiceRequest, error) {

	// Update the query to select only service_requests that do not have a matching service_provider_details entry for the given providerId
	query := config.ViewPendingRequestByProvider(serviceID)

	// Execute the query with the providerId, limit, and offset parameters
	var rows *sql.Rows
	var err error

	// Execute query with parameters based on presence of serviceID
	if serviceID != "" {
		rows, err = repo.db.Query(query, providerId, serviceID, limit, offset)
	} else {
		rows, err = repo.db.Query(query, providerId, limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []model.ServiceRequest
	for rows.Next() {
		var requestedTime, scheduledTime []uint8
		var request model.ServiceRequest

		// Scan data from rows into request fields
		err := rows.Scan(
			&request.ID, &request.HouseholderID, &request.HouseholderName, &request.HouseholderAddress,
			&request.ServiceID, &requestedTime, &scheduledTime, &request.Description,
			&request.Status, &request.ApproveStatus, &request.ServiceName,
		)
		if err != nil {
			return nil, err
		}

		// Parse the requested and scheduled time fields
		request.RequestedTime, err = util.ParseTime(requestedTime)
		if err != nil {
			return nil, fmt.Errorf("%v: %v", errs.ErrorParsingRequestTime, err)
		}

		request.ScheduledTime, err = util.ParseTime(scheduledTime)
		if err != nil {
			return nil, fmt.Errorf("%v: %v", errs.ErrorParsingScheduleTime, err)
		}

		// Append the parsed request to the list of requests
		requests = append(requests, request)
	}

	// Check for any errors encountered during row iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}
