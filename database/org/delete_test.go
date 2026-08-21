// SPDX-License-Identifier: Apache-2.0

package org

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestOrg_Engine_DeleteOrg(t *testing.T) {
	// setup types
	_org := testOrganization()

	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`DELETE FROM "orgs" WHERE name = $1`).
		WithArgs("github").
		WillReturnResult(sqlmock.NewResult(1, 1))

	_sqlite := testSqlite(t)

	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	_, err := _sqlite.CreateOrg(context.TODO(), _org)
	if err != nil {
		t.Errorf("unable to create test org for sqlite: %v", err)
	}

	// setup tests
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

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.database.DeleteOrg(context.TODO(), "github")

			if test.failure {
				if err == nil {
					t.Errorf("DeleteOrg for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("DeleteOrg for %s returned err: %v", test.name, err)
			}
		})
	}
}
