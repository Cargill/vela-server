// SPDX-License-Identifier: Apache-2.0

package org

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
)

type deleteMockDB struct {
	database.Interface
	getResult *api.OrgBuildLimit
	getErr    error
	deleteErr error
}

func (m *deleteMockDB) GetOrgBuildLimit(_ context.Context, _ string) (*api.OrgBuildLimit, error) {
	return m.getResult, m.getErr
}

func (m *deleteMockDB) DeleteOrgBuildLimit(_ context.Context, _ string) error {
	return m.deleteErr
}

func TestDeleteBuildLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		db, err := database.NewTest()
		if err != nil {
			t.Errorf("unable to create test database engine: %v", err)
		}

		defer db.Close()

		obl := new(api.OrgBuildLimit)
		obl.SetOrg("github")
		obl.SetBuildLimit(50)
		obl.SetCreatedAt(1)
		obl.SetUpdatedAt(1)
		obl.SetUpdatedBy("octocat")

		_, err = db.CreateOrgBuildLimit(context.TODO(), obl)
		if err != nil {
			t.Errorf("unable to create test org build limit: %v", err)
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodDelete, "/api/v1/repos/github/limit", nil)
		c.Set("logger", logrus.NewEntry(logrus.StandardLogger()))
		c.Set("org", "github")
		database.ToContext(c, db)

		DeleteBuildLimit(c)

		if w.Code != http.StatusOK {
			t.Errorf("DeleteBuildLimit returned %v, want %v", w.Code, http.StatusOK)
		}
	})

	t.Run("not-found", func(t *testing.T) {
		db, err := database.NewTest()
		if err != nil {
			t.Errorf("unable to create test database engine: %v", err)
		}

		defer db.Close()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodDelete, "/api/v1/repos/github/limit", nil)
		c.Set("logger", logrus.NewEntry(logrus.StandardLogger()))
		c.Set("org", "github")
		database.ToContext(c, db)

		DeleteBuildLimit(c)

		if w.Code != http.StatusNotFound {
			t.Errorf("DeleteBuildLimit returned %v, want %v", w.Code, http.StatusNotFound)
		}
	})

	t.Run("database-get-error", func(t *testing.T) {
		db, err := database.NewTest()
		if err != nil {
			t.Errorf("unable to create test database engine: %v", err)
		}

		db.Close()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodDelete, "/api/v1/repos/github/limit", nil)
		c.Set("logger", logrus.NewEntry(logrus.StandardLogger()))
		c.Set("org", "github")
		database.ToContext(c, db)

		DeleteBuildLimit(c)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("DeleteBuildLimit returned %v, want %v", w.Code, http.StatusInternalServerError)
		}
	})

	t.Run("database-delete-error", func(t *testing.T) {
		db, err := database.NewTest()
		if err != nil {
			t.Errorf("unable to create test database engine: %v", err)
		}

		defer db.Close()

		existing := new(api.OrgBuildLimit)
		existing.SetID(1)
		existing.SetOrg("github")
		existing.SetBuildLimit(50)
		existing.SetCreatedAt(1)
		existing.SetUpdatedAt(1)
		existing.SetUpdatedBy("octocat")

		mock := &deleteMockDB{
			Interface: db,
			getResult: existing,
			deleteErr: fmt.Errorf("delete failed"),
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodDelete, "/api/v1/repos/github/limit", nil)
		c.Set("logger", logrus.NewEntry(logrus.StandardLogger()))
		c.Set("org", "github")
		database.ToContext(c, mock)

		DeleteBuildLimit(c)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("DeleteBuildLimit returned %v, want %v", w.Code, http.StatusInternalServerError)
		}
	})
}
