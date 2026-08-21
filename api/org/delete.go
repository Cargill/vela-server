// SPDX-License-Identifier: Apache-2.0

package org

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/go-vela/server/database"
	oMiddleware "github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/util"
)

// swagger:operation DELETE /api/v1/repos/{org}/limit org DeleteBuildLimit
//
// Delete the concurrent build limit override for an organization
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
//     description: Successfully deleted the org build limit
//     schema:
//       type: string
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: No override exists for the org
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// DeleteBuildLimit represents the API handler to remove a concurrent
// build limit override for an organization. After deletion, the org
// falls back to the platform-configured default.
func DeleteBuildLimit(c *gin.Context) {
	l := c.MustGet("logger").(*logrus.Entry)
	ctx := c.Request.Context()
	o := oMiddleware.Retrieve(c)

	l.Debugf("deleting build limit for org %s", o)

	_, err := database.FromContext(c).GetOrgBuildLimit(ctx, o)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			retErr := fmt.Errorf("no build limit override exists for org %s", o)

			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		retErr := fmt.Errorf("unable to read build limit for org %s: %w", o, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	err = database.FromContext(c).DeleteOrgBuildLimit(ctx, o)
	if err != nil {
		retErr := fmt.Errorf("unable to delete build limit for org %s: %w", o, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	l.Infof("build limit override for org %s deleted", o)

	c.JSON(http.StatusOK, fmt.Sprintf("build limit override for org %s deleted", o))
}
