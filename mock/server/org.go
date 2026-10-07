// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	api "github.com/go-vela/server/api/types"
)

const (
	// OrgResp represents a JSON return for a single organization.
	OrgResp = `{
  "id": 1,
  "org": "github",
  "build_limit": 30,
  "created_at": 1,
  "updated_at": 1,
  "updated_by": "octocat"
}`
)

// getOrg has a param :org returns mock JSON for a http GET.
func getOrg(c *gin.Context) {
	data := []byte(OrgResp)

	var body api.Organization

	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// updateOrg has a param :org returns mock JSON for a http PUT.
func updateOrg(c *gin.Context) {
	data := []byte(OrgResp)

	var body api.Organization

	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// deleteOrg has a param :org returns mock JSON for a http DELETE.
func deleteOrg(c *gin.Context) {
	o := c.Param("org")

	c.JSON(http.StatusOK, fmt.Sprintf("build limit override for organization %s deleted", o))
}
