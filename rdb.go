package gocodecache

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type rdbSource struct {
	db              *sql.DB
	tableName       string
	keyColumnNames  [MaxKeyLength]string
	valueColumnName string
}

func RdbSource(db *sql.DB, tableName string,
	keyColumnNames [MaxKeyLength]string, valueColumnName string,
) Datasource {
	return &rdbSource{
		db:              db,
		tableName:       tableName,
		keyColumnNames:  keyColumnNames,
		valueColumnName: valueColumnName,
	}
}

func (d *rdbSource) ReadAll(ctx context.Context, keyLength int) (map[[MaxKeyLength]string]string, error) {
	if d.tableName == "" {
		return nil, errors.New("table name is empty")
	}

	keys := []string{}
	for i, k := range d.keyColumnNames {
		if keyLength <= i {
			if k != "" {
				return nil, fmt.Errorf("key column name is not empty (index = %d)", i)
			}
			break
		}
		if k == "" {
			return nil, fmt.Errorf("key column name is empty (index = %d)", i)
		}
		keys = append(keys, k)
	}

	if d.valueColumnName == "" {
		return nil, errors.New("value column name is empty")
	}

	query := fmt.Sprintf("SELECT %s, %s FROM %s", strings.Join(keys, ", "), d.valueColumnName, d.tableName)

	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to read from database: %w", err)
	}
	defer rows.Close()

	m := map[[MaxKeyLength]string]string{}
	dataPtr := make([]interface{}, keyLength+1)
	for rows.Next() {
		data := make([]string, keyLength+1)
		for i := range data {
			dataPtr[i] = &data[i]
		}

		if err := rows.Scan(dataPtr...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		keys := [MaxKeyLength]string{}
		for i := 0; i < keyLength; i++ {
			keys[i] = data[i]
		}

		m[keys] = data[keyLength]
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("failed to close row: %w", err)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred in row: %w", err)
	}

	return m, nil
}
