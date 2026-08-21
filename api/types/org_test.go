// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTypes_Organization_Getters(t *testing.T) {
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
			org:  new(Organization),
			want: new(Organization),
		},
	}

	// run tests
	for _, test := range tests {
		if test.org.GetID() != test.want.GetID() {
			t.Errorf("GetID is %v, want %v", test.org.GetID(), test.want.GetID())
		}

		if test.org.GetName() != test.want.GetName() {
			t.Errorf("GetName is %v, want %v", test.org.GetName(), test.want.GetName())
		}

		if test.org.GetBuildLimit() != test.want.GetBuildLimit() {
			t.Errorf("GetBuildLimit is %v, want %v", test.org.GetBuildLimit(), test.want.GetBuildLimit())
		}

		if test.org.GetCreatedAt() != test.want.GetCreatedAt() {
			t.Errorf("GetCreatedAt is %v, want %v", test.org.GetCreatedAt(), test.want.GetCreatedAt())
		}

		if test.org.GetUpdatedAt() != test.want.GetUpdatedAt() {
			t.Errorf("GetUpdatedAt is %v, want %v", test.org.GetUpdatedAt(), test.want.GetUpdatedAt())
		}

		if test.org.GetUpdatedBy() != test.want.GetUpdatedBy() {
			t.Errorf("GetUpdatedBy is %v, want %v", test.org.GetUpdatedBy(), test.want.GetUpdatedBy())
		}
	}
}

func TestTypes_Organization_Setters(t *testing.T) {
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
			want: new(Organization),
		},
	}

	// run tests
	for _, test := range tests {
		test.org.SetID(test.want.GetID())
		test.org.SetName(test.want.GetName())
		test.org.SetBuildLimit(test.want.GetBuildLimit())
		test.org.SetCreatedAt(test.want.GetCreatedAt())
		test.org.SetUpdatedAt(test.want.GetUpdatedAt())
		test.org.SetUpdatedBy(test.want.GetUpdatedBy())

		if test.org.GetID() != test.want.GetID() {
			t.Errorf("SetID is %v, want %v", test.org.GetID(), test.want.GetID())
		}

		if test.org.GetName() != test.want.GetName() {
			t.Errorf("SetName is %v, want %v", test.org.GetName(), test.want.GetName())
		}

		if test.org.GetBuildLimit() != test.want.GetBuildLimit() {
			t.Errorf("SetBuildLimit is %v, want %v", test.org.GetBuildLimit(), test.want.GetBuildLimit())
		}

		if test.org.GetCreatedAt() != test.want.GetCreatedAt() {
			t.Errorf("SetCreatedAt is %v, want %v", test.org.GetCreatedAt(), test.want.GetCreatedAt())
		}

		if test.org.GetUpdatedAt() != test.want.GetUpdatedAt() {
			t.Errorf("SetUpdatedAt is %v, want %v", test.org.GetUpdatedAt(), test.want.GetUpdatedAt())
		}

		if test.org.GetUpdatedBy() != test.want.GetUpdatedBy() {
			t.Errorf("SetUpdatedBy is %v, want %v", test.org.GetUpdatedBy(), test.want.GetUpdatedBy())
		}
	}
}

func TestTypes_Organization_String(t *testing.T) {
	// setup types
	o := testOrganization()

	want := fmt.Sprintf(`{
  ID: %d,
  Name: %s,
  BuildLimit: %d,
  CreatedAt: %d,
  UpdatedAt: %d,
  UpdatedBy: %s,
}`,
		o.GetID(),
		o.GetName(),
		o.GetBuildLimit(),
		o.GetCreatedAt(),
		o.GetUpdatedAt(),
		o.GetUpdatedBy(),
	)

	// run test
	got := o.String()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("String is %v, want %v", got, want)
	}
}

// testOrganization is a test helper function to create an Organization
// type with all fields set to a fake value.
func testOrganization() *Organization {
	o := new(Organization)

	o.SetID(1)
	o.SetName("github")
	o.SetBuildLimit(30)
	o.SetCreatedAt(1)
	o.SetUpdatedAt(1)
	o.SetUpdatedBy("octocat")

	return o
}
