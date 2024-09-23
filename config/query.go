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
func SelectLeftJoinQuery(tableName, joinTable, joinCondition, condition string, firstTableColumns []string, secondTableColumns []string) string {
	for i, coloumn := range firstTableColumns {
		firstTableColumns[i] = tableName + "." + coloumn
	}
	firstColNames := strings.Join(firstTableColumns, ", ")
	for i, coloumn := range secondTableColumns {
		secondTableColumns[i] = joinTable + "." + coloumn
	}
	secondColNames := strings.Join(secondTableColumns, ", ")
	query := fmt.Sprintf("SELECT %s, %s FROM %s", firstColNames, secondColNames, tableName)

	// Add INNER JOIN clause
	if joinTable != "" && joinCondition != "" {
		query += fmt.Sprintf(" LEFT JOIN %s ON %s", joinTable, joinCondition)
	}

	if condition != "" {
		query += fmt.Sprintf(" WHERE %s = ?", condition)
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
