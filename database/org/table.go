// SPDX-License-Identifier: Apache-2.0

package org

import (
	"context"

	"github.com/go-vela/server/constants"
)

const (
	// CreatePostgresTable represents a query to create the Postgres orgs table.
	CreatePostgresTable = `
CREATE TABLE
IF NOT EXISTS
orgs (
	id          SERIAL PRIMARY KEY,
	name        VARCHAR(250),
	build_limit INTEGER,
	created_at  BIGINT,
	updated_at  BIGINT,
	updated_by  VARCHAR(250),
	UNIQUE(name)
);
`

	// CreateSqliteTable represents a query to create the Sqlite orgs table.
	CreateSqliteTable = `
CREATE TABLE
IF NOT EXISTS
orgs (
	id          INTEGER PRIMARY KEY AUTOINCREMENT,
	name        TEXT,
	build_limit INTEGER,
	created_at  INTEGER,
	updated_at  INTEGER,
	updated_by  TEXT,
	UNIQUE(name)
);
`
)

// CreateOrganizationTable creates the `orgs` table in the database.
func (e *Engine) CreateOrganizationTable(ctx context.Context, driver string) error {
	e.logger.Tracef("creating orgs table for %s driver", driver)

	// handle the driver provided to create the table
	switch driver {
	case constants.DriverPostgres:
		// create the orgs table for Postgres
		return e.client.
			WithContext(ctx).
			Exec(CreatePostgresTable).Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// create the orgs table for Sqlite
		return e.client.
			WithContext(ctx).
			Exec(CreateSqliteTable).Error
	}
}
