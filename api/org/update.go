// SPDX-License-Identifier: Apache-2.0

package org

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	orgMiddleware "github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
)

// swagger:operation PUT /api/v1/orgs/{org}/limit org UpdateBuildLimit
//
// Update the concurrent build limit for an organization
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the organization
//   required: true
//   type: string
// - in: body
//   name: body
//   description: The organization build limit to apply
//   required: true
//   schema:
//     "$ref": "#/definitions/Organization"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the organization
//     schema:
//       "$ref": "#/definitions/Organization"
//   '201':
//     description: Successfully created the organization
//     schema:
//       "$ref": "#/definitions/Organization"
//   '400':
//     description: Invalid request payload
//     schema:
//       "$ref": "#/definitions/Error"
//   '401':
//     description: Unauthorized (platform admin required)
//     schema:
//       "$ref": "#/definitions/Error"
//   '403':
//     description: Organization build limits are disabled
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateBuildLimit represents the API handler to set the concurrent
// build limit for an organization. A value that is not positive is
// stored as zero, which means the platform default limit applies at
// enforcement time, and the record is created when no override yet
// exists for the organization. The endpoint is only available when
// org build limits are enabled on the server.
func UpdateBuildLimit(c *gin.Context) {
	l := c.MustGet("logger").(*logrus.Entry)
	ctx := c.Request.Context()
	o := orgMiddleware.Retrieve(c)
	u := user.Retrieve(c)

	l.Debugf("updating build limit for organization %s", o)

	// organization build limits must be enabled by platform admins (a default of zero means disabled)
	if c.GetInt32("defaultOrgBuildLimit") <= 0 {
		retErr := fmt.Errorf("organization build limits have been disabled by Vela admins")

		util.HandleError(c, http.StatusForbidden, retErr)

		return
	}

	// capture body from API request
	input := new(types.Organization)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for organization: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// a limit of zero means the platform default applies at enforcement time
	limit := max(input.GetBuildLimit(), 0)

	// look up any existing organization record
	existing, err := database.FromContext(c).GetOrg(ctx, o)

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		// create a new organization record
		newOrg := new(types.Organization)
		newOrg.SetName(o)
		newOrg.SetBuildLimit(limit)
		newOrg.SetCreatedAt(time.Now().UTC().Unix())
		newOrg.SetUpdatedAt(time.Now().UTC().Unix())
		newOrg.SetUpdatedBy(u.GetName())

		newOrg, err = database.FromContext(c).CreateOrg(ctx, newOrg)
		if err != nil {
			retErr := fmt.Errorf("unable to create build limit for organization %s: %w", o, err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		l.Infof("build limit for organization %s created with value %d", o, limit)

		c.JSON(http.StatusCreated, newOrg)

	case err != nil:
		retErr := fmt.Errorf("unable to read build limit for organization %s: %w", o, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

	default:
		// update the existing organization record
		existing.SetBuildLimit(limit)
		existing.SetUpdatedAt(time.Now().UTC().Unix())
		existing.SetUpdatedBy(u.GetName())

		existing, err = database.FromContext(c).UpdateOrg(ctx, existing)
		if err != nil {
			retErr := fmt.Errorf("unable to update build limit for organization %s: %w", o, err)

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return
		}

		l.Infof("build limit for organization %s updated to %d", o, limit)

		c.JSON(http.StatusOK, existing)
	}
}
