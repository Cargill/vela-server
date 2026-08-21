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
	getResult *api.Organization
	getErr    error
	deleteErr error
}

func (m *deleteMockDB) GetOrg(_ context.Context, _ string) (*api.Organization, error) {
	return m.getResult, m.getErr
}

func (m *deleteMockDB) DeleteOrg(_ context.Context, _ string) error {
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

		newOrg := new(api.Organization)
		newOrg.SetName("github")
		newOrg.SetBuildLimit(50)
		newOrg.SetCreatedAt(1)
		newOrg.SetUpdatedAt(1)
		newOrg.SetUpdatedBy("octocat")

		_, err = db.CreateOrg(context.TODO(), newOrg)
		if err != nil {
			t.Errorf("unable to create test organization: %v", err)
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodDelete, "/api/v1/orgs/github/limit", nil)
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
		c.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodDelete, "/api/v1/orgs/github/limit", nil)
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
		c.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodDelete, "/api/v1/orgs/github/limit", nil)
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

		existing := new(api.Organization)
		existing.SetID(1)
		existing.SetName("github")
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
		c.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodDelete, "/api/v1/orgs/github/limit", nil)
		c.Set("logger", logrus.NewEntry(logrus.StandardLogger()))
		c.Set("org", "github")
		database.ToContext(c, mock)

		DeleteBuildLimit(c)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("DeleteBuildLimit returned %v, want %v", w.Code, http.StatusInternalServerError)
		}
	})
}
