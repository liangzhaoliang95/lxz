/**
 * @author  zhaoliang.liang
 * @date  2025/8/4 11:20
 */

package database_drivers

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
)

type MySQLDriver struct {
	*DatabaseConn
}

func (_this *MySQLDriver) formatTableName(database, table string) string {
	return fmt.Sprintf("`%s`.`%s`", database, table)
}

func (_this *MySQLDriver) GetDBConn() (*DatabaseConn, error) {
	if _this.dbConn == nil {
		err := _this.InitConnect()
		if err != nil {
			return nil, err
		}
	}
	return _this.DatabaseConn, nil
}

// GetDbList 获取当前MySQL实例的数据库列表
func (_this *MySQLDriver) GetDbList() ([]string, error) {
	if _this.dbConn == nil {
		err := _this.InitConnect()
		if err != nil {
			return nil, err
		}
	}
	var databases []string

	sqlDB, err := _this.dbConn.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from gorm.DB: %w", err)
	}
	rows, err := sqlDB.Query("SHOW DATABASES")
	if err != nil {
		return nil, fmt.Errorf("failed to query databases: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			return nil, fmt.Errorf("failed to scan database name: %w", err)
		}
		databases = append(databases, dbName)
	}
	return databases, nil
}

// GetTableList 获取当前数据库的表列表
func (_this *MySQLDriver) GetTableList(dbName string) ([]string, error) {
	if _this.dbConn == nil {
		err := _this.InitConnect()
		if err != nil {
			return nil, err
		}
	}
	var tableList []string
	sqlDB, err := _this.dbConn.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from gorm.DB: %w", err)
	}
	rows, err := sqlDB.Query(fmt.Sprintf("SHOW TABLES FROM `%s`", dbName))
	if err != nil {
		return nil, fmt.Errorf("failed to query tables: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}
		tableList = append(tableList, tableName)
	}
	return tableList, nil
}

func (_this *MySQLDriver) GetRecords(
	database, table, where, sort string,
	offset, limit int,
) (paginatedResults [][]string, totalRecords int, err error) {
	if _this.dbConn == nil {
		err := _this.InitConnect()
		if err != nil {
			return nil, 0, err
		}
	}
	if table == "" {
		return nil, 0, errors.New("table name is required")
	}

	if database == "" {
		return nil, 0, errors.New("database name is required")
	}

	if limit == 0 {
		limit = DefaultRowLimit
	}

	sqlDB, err := _this.dbConn.DB()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get sql.DB from gorm.DB: %w", err)
	}

	query := "SELECT * FROM "
	query += _this.formatTableName(database, table)

	if where != "" {
		query += fmt.Sprintf(" %s", where)
	}

	if sort != "" {
		query += fmt.Sprintf(" ORDER BY %s", sort)
	}

	query += " LIMIT ?, ?"

	slog.Debug("Executing query", "query", query, "offset", offset, "limit", limit)

	paginatedRows, err := sqlDB.Query(query, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		_ = paginatedRows.Close()
	}()

	columns, err := paginatedRows.Columns()
	if err != nil {
		return nil, 0, err
	}

	paginatedResults = append(paginatedResults, columns)

	for paginatedRows.Next() {
		nullStringSlice := make([]sql.NullString, len(columns))

		rowValues := make([]interface{}, len(columns))
		for i := range nullStringSlice {
			rowValues[i] = &nullStringSlice[i]
		}

		err = paginatedRows.Scan(rowValues...)
		if err != nil {
			return nil, 0, err
		}

		var row []string
		for _, col := range nullStringSlice {
			if col.Valid {
				if col.String == "" {
					row = append(row, "EMPTY&")
				} else {
					row = append(row, col.String)
				}
			} else {
				row = append(row, "NULL&")
			}
		}

		paginatedResults = append(paginatedResults, row)
	}
	if err := paginatedRows.Err(); err != nil {
		return nil, 0, err
	}
	// close to release the connection
	if err := paginatedRows.Close(); err != nil {
		return nil, 0, err
	}

	countQuery := "SELECT COUNT(*) FROM "
	countQuery += fmt.Sprintf("`%s`.", database)
	countQuery += fmt.Sprintf("`%s`", table)
	row := sqlDB.QueryRow(countQuery)
	if err := row.Scan(&totalRecords); err != nil {
		return nil, 0, err
	}

	return paginatedResults, totalRecords, nil
}

func (_this *MySQLDriver) ExecuteQuery(query string) ([][]string, int, error) {
	if _this.dbConn == nil {
		err := _this.InitConnect()
		if err != nil {
			return nil, 0, err
		}
	}

	sqlDB, err := _this.dbConn.DB()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get sql.DB from gorm.DB: %w", err)
	}

	rows, err := sqlDB.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		_ = rows.Close()
	}()

	columns, err := rows.Columns()
	if err != nil {
		return nil, 0, err
	}

	records := make([][]string, 0)
	for rows.Next() {
		rowValues := make([]interface{}, len(columns))
		for i := range columns {
			rowValues[i] = new(sql.RawBytes)
		}

		err = rows.Scan(rowValues...)
		if err != nil {
			return nil, 0, err
		}

		var row []string
		for _, col := range rowValues {
			row = append(row, string(*col.(*sql.RawBytes)))
		}

		records = append(records, row)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	// Prepend the columns to the records.
	results := append([][]string{columns}, records...)

	return results, len(records), nil
}
