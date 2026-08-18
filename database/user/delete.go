// SPDX-License-Identifier: Apache-2.0

package user

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/Cargill/vela-server/api/types"
	"github.com/Cargill/vela-server/constants"
	"github.com/Cargill/vela-server/database/types"
)

// DeleteUser deletes an existing user from the database.
func (e *Engine) DeleteUser(ctx context.Context, u *api.User) error {
	e.logger.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Tracef("deleting user %s", u.GetName())

	// cast the API type to database type
	user := types.UserFromAPI(u)

	// send query to the database
	return e.client.
		WithContext(ctx).
		Table(constants.TableUser).
		Delete(user).
		Error
}
