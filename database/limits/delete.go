// SPDX-License-Identifier: Apache-2.0

package limits

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// DeleteOrgBuildLimit deletes an existing org build limit from the database.
func (e *Engine) DeleteOrgBuildLimit(ctx context.Context, org string) error {
	e.logger.WithFields(logrus.Fields{
		"org": org,
	}).Tracef("deleting org build limit for %s", org)

	o := new(types.OrgBuildLimit)

	return e.client.
		WithContext(ctx).
		Table(constants.TableOrgBuildLimit).
		Where("org = ?", org).
		Delete(o).
		Error
}
