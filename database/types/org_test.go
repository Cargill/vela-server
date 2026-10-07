// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
)

func TestTypes_Organization_Nullify(t *testing.T) {
	// setup types
	var o *Organization

	// setup tests
	tests := []struct {
		org  *Organization
		want *Organization
	}{
		{
			org:  testOrganization(),
			want: testOrganization(),
		},
		{
			org:  o,
			want: nil,
		},
		{
			org:  new(Organization),
			want: new(Organization),
		},
	}

	// run tests
	for _, test := range tests {
		got := test.org.Nullify()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Nullify is %v, want %v", got, test.want)
		}
	}
}

func TestTypes_Organization_ToAPI(t *testing.T) {
	// setup types
	want := new(api.Organization)
	want.SetID(1)
	want.SetName("github")
	want.SetBuildLimit(30)
	want.SetCreatedAt(1)
	want.SetUpdatedAt(1)
	want.SetUpdatedBy("octocat")

	// run test
	got := testOrganization().ToAPI()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToAPI is %v, want %v", got, want)
	}
}

func TestTypes_Organization_Validate(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		org     *Organization
	}{
		{
			failure: false,
			org:     testOrganization(),
		},
		{ // no Org set
			failure: true,
			org: &Organization{
				ID:         sql.NullInt64{Int64: 1, Valid: true},
				BuildLimit: sql.NullInt32{Int32: 30, Valid: true},
			},
		},
		{ // no BuildLimit set (zero means the platform default applies)
			failure: false,
			org: &Organization{
				ID:   sql.NullInt64{Int64: 1, Valid: true},
				Name: sql.NullString{String: "github", Valid: true},
			},
		},
		{ // negative BuildLimit set
			failure: true,
			org: &Organization{
				ID:         sql.NullInt64{Int64: 1, Valid: true},
				Name:       sql.NullString{String: "github", Valid: true},
				BuildLimit: sql.NullInt32{Int32: -1, Valid: true},
			},
		},
	}

	// run tests
	for _, test := range tests {
		err := test.org.Validate()

		if test.failure {
			if err == nil {
				t.Errorf("Validate should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Validate returned err: %v", err)
		}
	}
}

func TestTypes_Organization_OrganizationFromAPI(t *testing.T) {
	// setup types
	o := new(api.Organization)
	o.SetID(1)
	o.SetName("github")
	o.SetBuildLimit(30)
	o.SetCreatedAt(1)
	o.SetUpdatedAt(1)
	o.SetUpdatedBy("octocat")

	want := testOrganization()

	// run test
	got := OrganizationFromAPI(o)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("OrganizationFromAPI is %v, want %v", got, want)
	}
}

// testOrganization is a test helper function to create an Organization
// type with all fields set to a fake value.
func testOrganization() *Organization {
	return &Organization{
		ID:         sql.NullInt64{Int64: 1, Valid: true},
		Name:       sql.NullString{String: "github", Valid: true},
		BuildLimit: sql.NullInt32{Int32: 30, Valid: true},
		CreatedAt:  sql.NullInt64{Int64: 1, Valid: true},
		UpdatedAt:  sql.NullInt64{Int64: 1, Valid: true},
		UpdatedBy:  sql.NullString{String: "octocat", Valid: true},
	}
}
