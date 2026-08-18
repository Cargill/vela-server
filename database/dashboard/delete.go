// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/Cargill/vela-server/api/types"
	"github.com/Cargill/vela-server/constants"
	"github.com/Cargill/vela-server/database/types"
)

// DeleteDashboard deletes an existing dashboard from the database.
func (e *Engine) DeleteDashboard(ctx context.Context, d *api.Dashboard) error {
	e.logger.WithFields(logrus.Fields{
		"dashboard": d.GetID(),
	}).Tracef("deleting dashboard %s", d.GetID())

	dashboard := types.DashboardFromAPI(d)

	// send query to the database
	return e.client.
		WithContext(ctx).
		Table(constants.TableDashboard).
		Delete(dashboard).
		Error
}
