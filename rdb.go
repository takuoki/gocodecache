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
	keyColumnNames  []string
	valueColumnName string
}

func RdbSource(db *sql.DB, tableName string,
	keyColumnNames []string, valueColumnName string,
) Datasource {
	return &rdbSource{
		db:              db,
		tableName:       tableName,
		keyColumnNames:  keyColumnNames,
		valueColumnName: valueColumnName,
	}
}

func (d *rdbSource) ReadAll(ctx context.Context, keyLength int) (Codes, error) {
	return d.ReadFirstKeys(ctx, keyLength, nil)
}

func (d *rdbSource) ReadFirstKeys(ctx context.Context,
	keyLength int, firstKeys map[string]struct{}) (Codes, error) {
	return d.read(ctx, keyLength, firstKeys)
}

func (d *rdbSource) read(ctx context.Context,
	keyLength int, firstKeys map[string]struct{}) (Codes, error) {

	if d.tableName == "" {
		return nil, errors.New("table name is empty")
	}
	if len(d.keyColumnNames) != keyLength {
		return nil, errors.New("length of keyColumnNames doesn't match keyLength")
	}

	firstKey := ""
	keys := []string{}
	for i, k := range d.keyColumnNames {
		if k == "" {
			return nil, fmt.Errorf("key column name is empty (index = %d)", i)
		}
		if i == 0 {
			firstKey = k
		}
		keys = append(keys, k)
	}

	if d.valueColumnName == "" {
		return nil, errors.New("value column name is empty")
	}

	query := fmt.Sprintf("SELECT %s, %s FROM %s", strings.Join(keys, ", "), d.valueColumnName, d.tableName)
	firstKeyList := []interface{}{}
	if firstKeys != nil {
		ps := []string{}
		i := 0
		for k := range firstKeys {
			i++
			firstKeyList = append(firstKeyList, k)
			ps = append(ps, fmt.Sprintf("$%d", i))
		}
		query += fmt.Sprintf(" WHERE %s IN (%s)", firstKey, strings.Join(ps, ", "))
	}

	rows, err := d.db.QueryContext(ctx, query, firstKeyList...)
	if err != nil {
		return nil, fmt.Errorf("failed to read from database: %w", err)
	}
	defer rows.Close()

	m := Codes{}
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
