// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/Cargill/vela-server/api/types"
	"github.com/Cargill/vela-server/constants"
	"github.com/Cargill/vela-server/database/types"
)

// DeleteHook deletes an existing hook from the database.
func (e *Engine) DeleteHook(ctx context.Context, h *api.Hook) error {
	e.logger.WithFields(logrus.Fields{
		"hook": h.GetNumber(),
	}).Tracef("deleting hook %d", h.GetNumber())

	hook := types.HookFromAPI(h)

	// send query to the database
	return e.client.
		WithContext(ctx).
		Table(constants.TableHook).
		Delete(hook).
		Error
}
