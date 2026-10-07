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

// UpdateOrg updates an existing org in the database.
func (e *Engine) UpdateOrg(ctx context.Context, o *api.Organization) (*api.Organization, error) {
	e.logger.WithFields(logrus.Fields{
		"org": o.GetName(), //nolint:goconst // ignore making constant for log field key
	}).Tracef("updating org for %s", o.GetName())

	org := types.OrganizationFromAPI(o)

	// store org names in lowercase so lookups are case-insensitive
	org.Name.String = strings.ToLower(org.Name.String)

	err := org.Validate()
	if err != nil {
		return nil, err
	}

	err = e.client.
		WithContext(ctx).
		Table(constants.TableOrganization).
		Save(org).Error
	if err != nil {
		return nil, err
	}

	return org.ToAPI(), nil
}
