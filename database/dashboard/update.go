// SPDX-License-Identifier: Apache-2.0

package dashboard

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/Cargill/vela-server/api/types"
	"github.com/Cargill/vela-server/constants"
	"github.com/Cargill/vela-server/database/types"
)

// UpdateDashboard updates an existing dashboard in the database.
func (e *Engine) UpdateDashboard(ctx context.Context, d *api.Dashboard) (*api.Dashboard, error) {
	e.logger.WithFields(logrus.Fields{
		"dashboard": d.GetID(),
	}).Tracef("creating dashboard %s", d.GetID())

	dashboard := types.DashboardFromAPI(d)

	err := dashboard.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	err = e.client.
		WithContext(ctx).
		Table(constants.TableDashboard).
		Save(dashboard).Error
	if err != nil {
		return nil, err
	}

	return dashboard.ToAPI(), nil
}
