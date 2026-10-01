// SPDX-License-Identifier: Apache-2.0

package types

import "fmt"

// Organization is the API representation of an SCM
// organization.
//
// swagger:model Organization
type Organization struct {
	ID         *int64  `json:"id"`
	Name       *string `json:"org,omitempty"`
	BuildLimit *int32  `json:"build_limit,omitempty"`
	CreatedAt  *int64  `json:"created_at,omitempty"`
	UpdatedAt  *int64  `json:"updated_at,omitempty"`
	UpdatedBy  *string `json:"updated_by,omitempty"`
}

// GetID returns the ID field.
//
// When the provided Organization type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (o *Organization) GetID() int64 {
	// return zero value if Organization type or ID field is nil
	if o == nil || o.ID == nil {
		return 0
	}

	return *o.ID
}

// GetName returns the Name field.
//
// When the provided Organization type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (o *Organization) GetName() string {
	// return zero value if Organization type or Name field is nil
	if o == nil || o.Name == nil {
		return ""
	}

	return *o.Name
}

// GetBuildLimit returns the BuildLimit field.
//
// When the provided Organization type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (o *Organization) GetBuildLimit() int32 {
	// return zero value if Organization type or BuildLimit field is nil
	if o == nil || o.BuildLimit == nil {
		return 0
	}

	return *o.BuildLimit
}

// GetCreatedAt returns the CreatedAt field.
//
// When the provided Organization type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (o *Organization) GetCreatedAt() int64 {
	// return zero value if Organization type or CreatedAt field is nil
	if o == nil || o.CreatedAt == nil {
		return 0
	}

	return *o.CreatedAt
}

// GetUpdatedAt returns the UpdatedAt field.
//
// When the provided Organization type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (o *Organization) GetUpdatedAt() int64 {
	// return zero value if Organization type or UpdatedAt field is nil
	if o == nil || o.UpdatedAt == nil {
		return 0
	}

	return *o.UpdatedAt
}

// GetUpdatedBy returns the UpdatedBy field.
//
// When the provided Organization type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (o *Organization) GetUpdatedBy() string {
	// return zero value if Organization type or UpdatedBy field is nil
	if o == nil || o.UpdatedBy == nil {
		return ""
	}

	return *o.UpdatedBy
}

// SetID sets the ID field.
//
// When the provided Organization type is nil, it
// will set nothing and immediately return.
func (o *Organization) SetID(v int64) {
	// return if Organization type is nil
	if o == nil {
		return
	}

	o.ID = &v
}

// SetName sets the Name field.
//
// When the provided Organization type is nil, it
// will set nothing and immediately return.
func (o *Organization) SetName(v string) {
	// return if Organization type is nil
	if o == nil {
		return
	}

	o.Name = &v
}

// SetBuildLimit sets the BuildLimit field.
//
// When the provided Organization type is nil, it
// will set nothing and immediately return.
func (o *Organization) SetBuildLimit(v int32) {
	// return if Organization type is nil
	if o == nil {
		return
	}

	o.BuildLimit = &v
}

// SetCreatedAt sets the CreatedAt field.
//
// When the provided Organization type is nil, it
// will set nothing and immediately return.
func (o *Organization) SetCreatedAt(v int64) {
	// return if Organization type is nil
	if o == nil {
		return
	}

	o.CreatedAt = &v
}

// SetUpdatedAt sets the UpdatedAt field.
//
// When the provided Organization type is nil, it
// will set nothing and immediately return.
func (o *Organization) SetUpdatedAt(v int64) {
	// return if Organization type is nil
	if o == nil {
		return
	}

	o.UpdatedAt = &v
}

// SetUpdatedBy sets the UpdatedBy field.
//
// When the provided Organization type is nil, it
// will set nothing and immediately return.
func (o *Organization) SetUpdatedBy(v string) {
	// return if Organization type is nil
	if o == nil {
		return
	}

	o.UpdatedBy = &v
}

// String implements the Stringer interface for the Organization type.
func (o *Organization) String() string {
	return fmt.Sprintf(`{
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
}
