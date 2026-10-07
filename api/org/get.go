// SPDX-License-Identifier: Apache-2.0

package org

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	oMiddleware "github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/util"
)

// swagger:operation GET /api/v1/orgs/{org}/limit org GetBuildLimit
//
// Get the concurrent build limit for an organization
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
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the organization
//     schema:
//       "$ref": "#/definitions/Organization"
//   '401':
//     description: Unauthorized (platform admin required)
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// GetBuildLimit represents the API handler to get the concurrent
// build limit for an organization. When no override has been set
// for the org, the effective default limit is returned.
func GetBuildLimit(c *gin.Context) {
	l := c.MustGet("logger").(*logrus.Entry)
	ctx := c.Request.Context()
	o := oMiddleware.Retrieve(c)
	defaultOrgBuildLimit := c.GetInt32("defaultOrgBuildLimit")

	l.Debugf("reading build limit for organization %s", o)

	limit, err := database.FromContext(c).GetOrg(ctx, o)
	if err != nil {
		// no override set for the organization - return the effective default
		if errors.Is(err, gorm.ErrRecordNotFound) {
			limit = new(types.Organization)
			limit.SetName(o)
			limit.SetBuildLimit(defaultOrgBuildLimit)

			c.JSON(http.StatusOK, limit)

			return
		}

		retErr := fmt.Errorf("unable to read build limit for organization %s: %w", o, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// a stored limit of zero means the effective default applies
	if limit.GetBuildLimit() <= 0 {
		limit.SetBuildLimit(defaultOrgBuildLimit)
	}

	c.JSON(http.StatusOK, limit)
}
