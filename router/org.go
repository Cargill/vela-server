// SPDX-License-Identifier: Apache-2.0

package router

import (
	"github.com/gin-gonic/gin"

	orgapi "github.com/go-vela/server/api/org"
	"github.com/go-vela/server/router/middleware"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/perm"
)

// OrgHandlers is a function that extends the provided base router group
// with the API handlers for org functionality.
//
// GET    /api/v1/orgs/:org/limit
// PUT    /api/v1/orgs/:org/limit
// DELETE /api/v1/orgs/:org/limit .
func OrgHandlers(base *gin.RouterGroup) {
	// Orgs endpoints
	_orgs := base.Group("/orgs")
	{
		// Org endpoints
		_org := _orgs.Group("/:org", org.Establish())
		{
			_org.GET("/limit", perm.MustPlatformAdmin(), orgapi.GetBuildLimit)
			_org.PUT("/limit", perm.MustPlatformAdmin(), middleware.Payload(), orgapi.UpdateBuildLimit)
			_org.DELETE("/limit", perm.MustPlatformAdmin(), orgapi.DeleteBuildLimit)
		}
	}
}
