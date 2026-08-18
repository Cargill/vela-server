// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/Cargill/vela-server/api/types"
	"github.com/Cargill/vela-server/constants"
	"github.com/Cargill/vela-server/database/types"
)

// CreateWorker creates a new worker in the database.
func (e *Engine) CreateWorker(ctx context.Context, w *api.Worker) (*api.Worker, error) {
	e.logger.WithFields(logrus.Fields{
		"worker": w.GetHostname(),
	}).Tracef("creating worker %s", w.GetHostname())

	// cast the API type to database type
	worker := types.WorkerFromAPI(w)

	// validate the necessary fields are populated
	err := worker.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	result := e.client.
		WithContext(ctx).
		Table(constants.TableWorker).
		Create(worker)

	return worker.ToAPI(w.GetRunningBuilds()), result.Error
}
