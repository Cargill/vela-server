// SPDX-License-Identifier: Apache-2.0

package org

import (
	"context"
	"strings"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetOrg gets an org by name from the database.
func (e *Engine) GetOrg(ctx context.Context, name string) (*api.Organization, error) {
	e.logger.Tracef("getting org for %s", name)

	o := new(types.Organization)

	err := e.client.
		WithContext(ctx).
		Table(constants.TableOrganization).
		Where("name = ?", strings.ToLower(name)).
		Take(o).
		Error
	if err != nil {
		return nil, err
	}

	return o.ToAPI(), nil
}
