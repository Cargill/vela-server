// SPDX-License-Identifier: Apache-2.0

package org

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// DeleteOrg deletes an existing org from the database.
func (e *Engine) DeleteOrg(ctx context.Context, name string) error {
	e.logger.WithFields(logrus.Fields{
		"org": name,
	}).Tracef("deleting org row from orgs for %s", name)

	o := new(types.Organization)

	return e.client.
		WithContext(ctx).
		Table(constants.TableOrganization).
		Where("name = ?", name).
		Delete(o).
		Error
}
