// SPDX-License-Identifier: Apache-2.0

package org

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/testutils"
)

func TestOrg_Engine_UpdateOrg(t *testing.T) {
	// setup types
	_org := testOrganization()

	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "orgs" SET "name"=$1,"build_limit"=$2,"created_at"=$3,"updated_at"=$4,"updated_by"=$5 WHERE "id" = $6`).
		WithArgs("github", 30, 1, testutils.AnyArgument{}, "octocat", 1).
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
			got, err := test.database.UpdateOrg(context.TODO(), _org)
			got.SetUpdatedAt(_org.GetUpdatedAt())

			if test.failure {
				if err == nil {
					t.Errorf("UpdateOrg for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("UpdateOrg for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, _org) {
				t.Errorf("UpdateOrg for %s returned %s, want %s", test.name, got, _org)
			}
		})
	}
}
