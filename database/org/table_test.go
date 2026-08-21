// SPDX-License-Identifier: Apache-2.0

package org

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestOrg_Engine_CreateOrganizationTable(t *testing.T) {
	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	_mock.ExpectExec(CreatePostgresTable).WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)

	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	tests := []struct {
		failure  bool
		name     string
		database *Engine
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.database.CreateOrganizationTable(context.TODO(), test.name)

			if test.failure {
				if err == nil {
					t.Errorf("CreateOrganizationTable for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CreateOrganizationTable for %s returned err: %v", test.name, err)
			}
		})
	}
}
