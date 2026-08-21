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

// swagger:operation GET /api/v1/repos/{org}/limit org GetBuildLimit
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
//     description: Successfully retrieved the org build limit
//     schema:
//       "$ref": "#/definitions/OrgBuildLimit"
//   '401':
//     description: Unauthorized
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
	defaultOrgBuildLimit := c.Value("defaultOrgBuildLimit").(int32)

	l.Debugf("reading build limit for org %s", o)

	limit, err := database.FromContext(c).GetOrgBuildLimit(ctx, o)
	if err != nil {
		// no override set for the org - return the effective default
		if errors.Is(err, gorm.ErrRecordNotFound) {
			limit = new(types.OrgBuildLimit)
			limit.SetOrg(o)
			limit.SetBuildLimit(defaultOrgBuildLimit)

			c.JSON(http.StatusOK, limit)

			return
		}

		retErr := fmt.Errorf("unable to read build limit for org %s: %w", o, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, limit)
}
