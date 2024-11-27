package config

import (
	"fmt"
	"strings"
)

func SelectInnerJoinQuery(tableName, joinTable, joinCondition, condition string, firstTableColumns []string, secondTableColumns []string) string {
	for i, coloumn := range firstTableColumns {
		firstTableColumns[i] = tableName + "." + coloumn
	}
	firstColNames := strings.Join(firstTableColumns, ", ")
	var secondColNames string
	if len(secondTableColumns) > 0 {
		for i, coloumn := range secondTableColumns {
			secondTableColumns[i] = joinTable + "." + coloumn
		}
		secondColNames = strings.Join(secondTableColumns, ", ")
	}
	var query string
	if len(secondTableColumns) > 0 {
		query = fmt.Sprintf("SELECT %s, %s FROM %s", firstColNames, secondColNames, tableName)
	} else {
		query = fmt.Sprintf("SELECT %s FROM %s", firstColNames, tableName)
	}

	// Add INNER JOIN clause
	if joinTable != "" && joinCondition != "" {
		query += fmt.Sprintf(" INNER JOIN %s ON %s", joinTable, joinCondition)
	} else if joinCondition == "" {
		query += fmt.Sprintf(" INNER JOIN %s", joinTable)
	}

	if condition != "" {
		query += fmt.Sprintf(" WHERE %s = ?", condition)
	}

	return query
}
func SelectLeftJoinQuery(tableName, joinTable, joinCondition, condition string, firstTableColumns []string, secondTableColumns []string, limit, offset int) string {

	for i, column := range firstTableColumns {
		firstTableColumns[i] = tableName + "." + column
	}
	firstColNames := strings.Join(firstTableColumns, ", ")

	for i, column := range secondTableColumns {
		secondTableColumns[i] = joinTable + "." + column
	}
	secondColNames := strings.Join(secondTableColumns, ", ")

	query := fmt.Sprintf("SELECT %s, %s FROM %s", firstColNames, secondColNames, tableName)

	if joinTable != "" && joinCondition != "" {
		query += fmt.Sprintf(" LEFT JOIN %s ON %s", joinTable, joinCondition)
	}

	if condition != "" {
		query += fmt.Sprintf(" WHERE %s = ?", condition)
	}

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", offset)
	}

	return query
}

func SelectQuery(tableName, condition1, condition2 string, columns []string) string {
	colNames := strings.Join(columns, ", ")
	var query string
	if condition1 == "" && condition2 == "" {
		query = fmt.Sprintf("SELECT %s FROM %s", colNames, tableName)
	}
	if condition1 != "" && condition2 == "" {
		query = fmt.Sprintf("SELECT %s FROM %s WHERE %s = ?", colNames, tableName, condition1)
	}
	if condition1 != "" && condition2 != "" {
		query = fmt.Sprintf("SELECT %s FROM %s WHERE %s = ? AND %s = ?", colNames, tableName, condition1, condition2)
	}
	return query
}

func DeleteQuery(tableName, condition1, condition2 string) string {
	if condition2 == "" {
		query := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", tableName, condition1)
		return query
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = ? AND %s = ?", tableName, condition1, condition2)
	return query
}

func UpdateQuery(tableName, condition1, condition2 string, columns []string) string {
	setClause := make([]string, len(columns))
	for i, col := range columns {
		setClause[i] = fmt.Sprintf("%s = ?", col)
	}
	setClauseStr := strings.Join(setClause, ", ")
	if condition2 == "" {
		query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?", tableName, setClauseStr, condition1)
		return query
	}
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ? AND %s = ?", tableName, setClauseStr, condition1, condition2)
	return query
}
func InsertQuery(tableName string, columns []string) string {
	colNames := strings.Join(columns, ", ")
	placeholders := strings.Repeat("?, ", len(columns))
	placeholders = strings.TrimSuffix(placeholders, ", ")
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, colNames, placeholders)
	return query
}
func SelectCountQuery(tableName, condition string) string {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)

	// Add WHERE clause if a condition is provided
	if condition != "" {
		query += fmt.Sprintf(" WHERE %s = ?", condition)
	}

	return query
}
func SelectAverageQuery(tableName, column, condition string) string {
	query := fmt.Sprintf("SELECT AVG(%s) FROM %s", column, tableName)

	// Add WHERE clause if a condition is provided
	if condition != "" {
		query += fmt.Sprintf(" WHERE %s = ?", condition)
	}

	return query
}

func SelectJsonDataQuery(withStatus bool) string {
	query := `
    SELECT
        service_requests.id AS request_id,
        service_requests.householder_id,
        service_requests.householder_name,
        service_requests.householder_address,
        service_requests.service_id,
        service_requests.requested_time,
        service_requests.scheduled_time,
        service_requests.status,
        service_requests.approve_status,
        service_requests.service_name,
        CONCAT(
            '[',
            GROUP_CONCAT(
                JSON_OBJECT(
                    'service_provider_id', service_provider_details.service_provider_id,
                    'name', service_provider_details.name,
                    'contact', service_provider_details.contact,
                    'address', service_provider_details.address,
                    'price', service_provider_details.price,
                    'rating', service_provider_details.rating,
                    'approve', service_provider_details.approve
                )
                SEPARATOR ', '
            ),
            ']'
        ) AS provider_details
    FROM
        service_requests
    LEFT JOIN
        service_provider_details ON service_requests.id = service_provider_details.service_request_id
    WHERE
        service_requests.householder_id = ?
    `

	if withStatus {
		query += " AND service_requests.status = ? "
	}

	query += `
    GROUP BY
        service_requests.id
    LIMIT ? OFFSET ?;
    `
	return query
}

func SelectJsonDataQueryWithApprove(sortOrder string) string {
	baseQuery := `
		SELECT
			service_requests.id AS request_id,
			service_requests.householder_id,
			service_requests.householder_name,
			service_requests.householder_address,
			service_requests.service_id,
			service_requests.requested_time,
			service_requests.scheduled_time,
			service_requests.status,
			service_requests.approve_status,
			service_requests.service_name,
			CONCAT(
				'[',
				GROUP_CONCAT(
					JSON_OBJECT(
						'service_provider_id', service_provider_details.service_provider_id,
						'name', service_provider_details.name,
						'contact', service_provider_details.contact,
						'address', service_provider_details.address,
						'price', service_provider_details.price,
						'rating', service_provider_details.rating,
						'approve', service_provider_details.approve
					)
					SEPARATOR ', '
				),
				']'
			) AS provider_details
		FROM
			service_requests
		LEFT JOIN
			service_provider_details ON service_requests.id = service_provider_details.service_request_id
		WHERE
			service_requests.householder_id = ?
			AND service_requests.approve_status = ?
		GROUP BY
			service_requests.id
	`
	// Append ORDER BY clause if sortOrder is valid
	if sortOrder == "ASC" || sortOrder == "DESC" {
		baseQuery += fmt.Sprintf(" ORDER BY service_requests.scheduled_time %s", sortOrder)
	}

	// Add LIMIT and OFFSET for pagination
	baseQuery += " LIMIT ? OFFSET ?;"

	return baseQuery
}

func ViewPendingRequestByProvider(serviceID string) string {
	baseQuery := `
        SELECT sr.id, sr.householder_id, sr.householder_name, sr.householder_address, sr.service_id, 
               sr.requested_time, sr.scheduled_time, sr.description, sr.status, sr.approve_status, sr.service_name
        FROM service_requests sr
        LEFT JOIN service_provider_details spd 
        ON sr.id = spd.service_request_id AND spd.service_provider_id = ?
        WHERE spd.service_request_id IS NULL AND (sr.status="pending" OR sr.status="accepted")`

	// Add a filter for service_id if provided
	if serviceID != "" {
		baseQuery += ` AND sr.service_id = ?`
	}

	// Append limit and offset
	baseQuery += ` LIMIT ? OFFSET ?;`
	return baseQuery
}

func SelectInnerJoinQueryPaginate(tableName, joinTable, joinCondition, condition string, firstTableColumns []string, secondTableColumns []string, limit, offset int, sortColumn, sortOrder string) string {
	// Prefix columns with their table names
	for i, column := range firstTableColumns {
		firstTableColumns[i] = tableName + "." + column
	}
	firstColNames := strings.Join(firstTableColumns, ", ")

	var secondColNames string
	if len(secondTableColumns) > 0 {
		for i, column := range secondTableColumns {
			secondTableColumns[i] = joinTable + "." + column
		}
		secondColNames = strings.Join(secondTableColumns, ", ")
	}

	// Base query with table columns
	var query string
	if len(secondTableColumns) > 0 {
		query = fmt.Sprintf("SELECT %s, %s FROM %s", firstColNames, secondColNames, tableName)
	} else {
		query = fmt.Sprintf("SELECT %s FROM %s", firstColNames, tableName)
	}

	// Add INNER JOIN clause
	if joinTable != "" && joinCondition != "" {
		query += fmt.Sprintf(" INNER JOIN %s ON %s", joinTable, joinCondition)
	}

	// Add WHERE condition if provided
	if condition != "" {
		query += fmt.Sprintf(" WHERE %s = ?", condition)
	}

	// Add ORDER BY clause if sortColumn and sortOrder are provided
	if sortColumn != "" && (sortOrder == "ASC" || sortOrder == "DESC") {
		query += fmt.Sprintf(" ORDER BY %s %s", sortColumn, sortOrder)
	}

	// Add LIMIT and OFFSET for pagination
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", offset)
	}

	return query
}

func SelectQueryWithLimit(tableName, condition1, condition2 string, columns []string, limit, offset int) string {
	colNames := strings.Join(columns, ", ")
	query := fmt.Sprintf("SELECT %s FROM %s", colNames, tableName)

	// Add conditions only if they are provided
	if condition1 != "" && condition2 != "" {
		query += fmt.Sprintf(" WHERE %s = ? AND %s = ?", condition1, condition2)
	} else if condition1 != "" {
		query += fmt.Sprintf(" WHERE %s = ?", condition1)
	}

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", offset)
	}

	return query
}

func CountReviewAddedQuery() string {
	return `SELECT COUNT(*) FROM reviews WHERE provider_id = ? AND service_id = ? AND householder_id = ?`
}
