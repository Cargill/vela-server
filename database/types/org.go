// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"errors"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/util"
)

var (
	// ErrEmptyOrganizationNameField defines the error type when an
	// Organization type has an empty Name field provided.
	ErrEmptyOrganizationNameField = errors.New("empty organization name provided")

	// ErrInvalidOrganizationBuildLimit defines the error type when an
	// Organization type has a BuildLimit field that is not positive.
	ErrInvalidOrganizationBuildLimit = errors.New("organization build limit must be greater than zero")
)

// Organization is the database representation of an SCM organization
// lazily created for override configurations such as build limits.
type Organization struct {
	ID         sql.NullInt64  `sql:"id"`
	Name       sql.NullString `sql:"name"`
	BuildLimit sql.NullInt32  `sql:"build_limit"`
	CreatedAt  sql.NullInt64  `sql:"created_at"`
	UpdatedAt  sql.NullInt64  `sql:"updated_at"`
	UpdatedBy  sql.NullString `sql:"updated_by"`
}

// Nullify ensures the valid flag for
// the sql.Null types are properly set.
//
// When a field within the Organization type is the zero
// value for the field, the valid flag is set to
// false causing it to be NULL in the database.
func (o *Organization) Nullify() *Organization {
	if o == nil {
		return nil
	}

	// check if the ID field should be false
	if o.ID.Int64 == 0 {
		o.ID.Valid = false
	}

	// check if the Name field should be false
	if len(o.Name.String) == 0 {
		o.Name.Valid = false
	}

	// check if the CreatedAt field should be false
	if o.CreatedAt.Int64 < 0 {
		o.CreatedAt.Valid = false
	}

	// check if the UpdatedAt field should be false
	if o.UpdatedAt.Int64 < 0 {
		o.UpdatedAt.Valid = false
	}

	// check if the UpdatedBy field should be false
	if len(o.UpdatedBy.String) == 0 {
		o.UpdatedBy.Valid = false
	}

	return o
}

// ToAPI converts the Organization type
// to an API Organization type.
func (o *Organization) ToAPI() *api.Organization {
	oAPI := new(api.Organization)

	oAPI.SetID(o.ID.Int64)
	oAPI.SetName(o.Name.String)
	oAPI.SetBuildLimit(o.BuildLimit.Int32)
	oAPI.SetCreatedAt(o.CreatedAt.Int64)
	oAPI.SetUpdatedAt(o.UpdatedAt.Int64)
	oAPI.SetUpdatedBy(o.UpdatedBy.String)

	return oAPI
}

// Validate verifies the necessary fields for
// the Organization type are populated correctly.
func (o *Organization) Validate() error {
	// verify the Name field is populated
	if len(o.Name.String) == 0 {
		return ErrEmptyOrganizationNameField
	}

	// verify the BuildLimit field is positive
	if o.BuildLimit.Int32 <= 0 {
		return ErrInvalidOrganizationBuildLimit
	}

	// ensure that the Name field is sanitized
	// to avoid unsafe HTML content
	o.Name = sql.NullString{String: util.Sanitize(o.Name.String), Valid: o.Name.Valid}

	return nil
}

// OrganizationFromAPI converts the API Organization type
// to a database Organization type.
func OrganizationFromAPI(o *api.Organization) *Organization {
	org := &Organization{
		ID:         sql.NullInt64{Int64: o.GetID(), Valid: true},
		Name:       sql.NullString{String: o.GetName(), Valid: true},
		BuildLimit: sql.NullInt32{Int32: o.GetBuildLimit(), Valid: true},
		CreatedAt:  sql.NullInt64{Int64: o.GetCreatedAt(), Valid: true},
		UpdatedAt:  sql.NullInt64{Int64: o.GetUpdatedAt(), Valid: true},
		UpdatedBy:  sql.NullString{String: o.GetUpdatedBy(), Valid: true},
	}

	return org.Nullify()
}
