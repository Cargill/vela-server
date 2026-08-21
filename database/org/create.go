// SPDX-License-Identifier: Apache-2.0

package org

import (
	"context"
	"strings"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// CreateOrg creates a new org in the database.
func (e *Engine) CreateOrg(ctx context.Context, o *api.Organization) (*api.Organization, error) {
	e.logger.WithFields(logrus.Fields{
		"org": o.GetName(), //nolint:goconst // ignore making constant for log field key
	}).Tracef("creating org for %s", o.GetName())

	org := types.OrganizationFromAPI(o)

	// store org names in lowercase so lookups are case-insensitive
	org.Name.String = strings.ToLower(org.Name.String)

	err := org.Validate()
	if err != nil {
		return nil, err
	}

	result := e.client.
		WithContext(ctx).
		Table(constants.TableOrganization).
		Create(org)

	return org.ToAPI(), result.Error
}
