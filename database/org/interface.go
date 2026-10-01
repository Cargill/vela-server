// SPDX-License-Identifier: Apache-2.0

package org

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

// OrgInterface represents the Vela interface for org
// functions with the supported Database backends.
type OrgInterface interface {
	// Org Data Definition Language Functions

	// CreateOrganizationTable defines a function that creates the orgs table.
	CreateOrganizationTable(context.Context, string) error

	// Org Data Manipulation Language Functions

	// CreateOrg defines a function that creates an org.
	CreateOrg(context.Context, *api.Organization) (*api.Organization, error)
	// DeleteOrg defines a function that deletes an org by name.
	DeleteOrg(context.Context, string) error
	// GetOrg defines a function that gets an org by name.
	GetOrg(context.Context, string) (*api.Organization, error)
	// UpdateOrg defines a function that updates an org.
	UpdateOrg(context.Context, *api.Organization) (*api.Organization, error)
}
