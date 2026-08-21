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
	// OrgBuildLimitResp represents a JSON return for a single org build limit.
	OrgBuildLimitResp = `{
  "id": 1,
  "org": "github",
  "build_limit": 30,
  "created_at": 1,
  "updated_at": 1,
  "updated_by": "octocat"
}`
)

// getOrgBuildLimit has a param :org returns mock JSON for a http GET.
func getOrgBuildLimit(c *gin.Context) {
	data := []byte(OrgBuildLimitResp)

	var body api.OrgBuildLimit

	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// updateOrgBuildLimit has a param :org returns mock JSON for a http PUT.
func updateOrgBuildLimit(c *gin.Context) {
	data := []byte(OrgBuildLimitResp)

	var body api.OrgBuildLimit

	_ = json.Unmarshal(data, &body)

	c.JSON(http.StatusOK, body)
}

// deleteOrgBuildLimit has a param :org returns mock JSON for a http DELETE.
func deleteOrgBuildLimit(c *gin.Context) {
	o := c.Param("org")

	c.JSON(http.StatusOK, fmt.Sprintf("build limit override for org %s deleted", o))
}
