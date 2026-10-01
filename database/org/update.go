// SPDX-License-Identifier: Apache-2.0

package org

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// UpdateOrg updates an existing org in the database.
func (e *Engine) UpdateOrg(ctx context.Context, o *api.Organization) (*api.Organization, error) {
	e.logger.WithFields(logrus.Fields{
		"org": o.GetName(),
	}).Tracef("updating org for %s", o.GetName())

	org := types.OrganizationFromAPI(o)

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
