// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/Cargill/vela-server/api/types"
	"github.com/Cargill/vela-server/constants"
	"github.com/Cargill/vela-server/database/types"
)

// DeleteWorker deletes an existing worker from the database.
func (e *Engine) DeleteWorker(ctx context.Context, w *api.Worker) error {
	e.logger.WithFields(logrus.Fields{
		"worker": w.GetHostname(),
	}).Tracef("deleting worker %s", w.GetHostname())

	// cast the API type to database type
	worker := types.WorkerFromAPI(w)

	// send query to the database
	return e.client.
		WithContext(ctx).
		Table(constants.TableWorker).
		Delete(worker).
		Error
}
