// SPDX-License-Identifier: Apache-2.0

package org

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
)

func TestGetBuildLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("found", func(t *testing.T) {
		db, err := database.NewTest()
		if err != nil {
			t.Errorf("unable to create test database engine: %v", err)
		}

		defer db.Close()

		want := new(api.Organization)
		want.SetName("github")
		want.SetBuildLimit(50)
		want.SetCreatedAt(1)
		want.SetUpdatedAt(1)
		want.SetUpdatedBy("octocat")

		want, err = db.CreateOrg(context.TODO(), want)
		if err != nil {
			t.Errorf("unable to create test organization: %v", err)
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/repos/github/limit", nil)
		c.Set("logger", logrus.NewEntry(logrus.StandardLogger()))
		c.Set("org", "github")
		c.Set("defaultOrgBuildLimit", int32(30))
		database.ToContext(c, db)

		GetBuildLimit(c)

		if w.Code != http.StatusOK {
			t.Errorf("GetBuildLimit returned %v, want %v", w.Code, http.StatusOK)
		}

		var got api.Organization

		if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
			t.Errorf("unable to unmarshal response: %v", err)
		}

		if got.GetBuildLimit() != want.GetBuildLimit() {
			t.Errorf("GetBuildLimit build_limit is %v, want %v", got.GetBuildLimit(), want.GetBuildLimit())
		}

		if got.GetName() != want.GetName() {
			t.Errorf("GetBuildLimit name is %v, want %v", got.GetName(), want.GetName())
		}
	})

	t.Run("not-found-returns-default", func(t *testing.T) {
		db, err := database.NewTest()
		if err != nil {
			t.Errorf("unable to create test database engine: %v", err)
		}

		defer db.Close()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/repos/github/limit", nil)
		c.Set("logger", logrus.NewEntry(logrus.StandardLogger()))
		c.Set("org", "github")
		c.Set("defaultOrgBuildLimit", int32(30))
		database.ToContext(c, db)

		GetBuildLimit(c)

		if w.Code != http.StatusOK {
			t.Errorf("GetBuildLimit returned %v, want %v", w.Code, http.StatusOK)
		}

		var got api.Organization

		if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
			t.Errorf("unable to unmarshal response: %v", err)
		}

		if got.GetBuildLimit() != 30 {
			t.Errorf("GetBuildLimit build_limit is %v, want %v", got.GetBuildLimit(), int32(30))
		}

		if got.GetName() != "github" {
			t.Errorf("GetBuildLimit name is %v, want %v", got.GetName(), "github")
		}
	})

	t.Run("database-error", func(t *testing.T) {
		db, err := database.NewTest()
		if err != nil {
			t.Errorf("unable to create test database engine: %v", err)
		}

		db.Close()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodGet, "/api/v1/repos/github/limit", nil)
		c.Set("logger", logrus.NewEntry(logrus.StandardLogger()))
		c.Set("org", "github")
		c.Set("defaultOrgBuildLimit", int32(30))
		database.ToContext(c, db)

		GetBuildLimit(c)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("GetBuildLimit returned %v, want %v", w.Code, http.StatusInternalServerError)
		}
	})
}
