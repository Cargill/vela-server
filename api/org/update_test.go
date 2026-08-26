// SPDX-License-Identifier: Apache-2.0

package org

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/api/types/settings"
	"github.com/go-vela/server/database"
)

type updateMockDB struct {
	database.Interface
	getResult *api.OrgBuildLimit
	getErr    error
	createErr error
	updateErr error
}

func (m *updateMockDB) GetOrgBuildLimit(_ context.Context, _ string) (*api.OrgBuildLimit, error) {
	return m.getResult, m.getErr
}

func (m *updateMockDB) CreateOrgBuildLimit(_ context.Context, _ *api.OrgBuildLimit) (*api.OrgBuildLimit, error) {
	return nil, m.createErr
}

func (m *updateMockDB) UpdateOrgBuildLimit(_ context.Context, _ *api.OrgBuildLimit) (*api.OrgBuildLimit, error) {
	return nil, m.updateErr
}

func TestUpdateBuildLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("disabled", func(t *testing.T) {
		ps := new(settings.Platform)
		ps.SetEnableOrgBuildLimit(false)

		u := new(api.User)
		u.SetName("octocat")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodPut, "/api/v1/repos/github/limit", nil)
		c.Set("logger", logrus.NewEntry(logrus.StandardLogger()))
		c.Set("org", "github")
		c.Set("user", u)
		c.Set("settings", ps)
		c.Set("defaultOrgBuildLimit", int32(30))
		c.Set("maxOrgBuildLimit", int32(100))

		UpdateBuildLimit(c)

		if w.Code != http.StatusForbidden {
			t.Errorf("UpdateBuildLimit returned %v, want %v", w.Code, http.StatusForbidden)
		}
	})

	t.Run("invalid-body", func(t *testing.T) {
		ps := new(settings.Platform)
		ps.SetEnableOrgBuildLimit(true)

		u := new(api.User)
		u.SetName("octocat")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodPut, "/api/v1/repos/github/limit", strings.NewReader("not json"))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("logger", logrus.NewEntry(logrus.StandardLogger()))
		c.Set("org", "github")
		c.Set("user", u)
		c.Set("settings", ps)
		c.Set("defaultOrgBuildLimit", int32(30))
		c.Set("maxOrgBuildLimit", int32(100))

		UpdateBuildLimit(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("UpdateBuildLimit returned %v, want %v", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("create", func(t *testing.T) {
		db, err := database.NewTest()
		if err != nil {
			t.Errorf("unable to create test database engine: %v", err)
		}

		defer db.Close()

		ps := new(settings.Platform)
		ps.SetEnableOrgBuildLimit(true)

		u := new(api.User)
		u.SetName("octocat")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodPut, "/api/v1/repos/github/limit", strings.NewReader(`{"build_limit": 50}`))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("logger", logrus.NewEntry(logrus.StandardLogger()))
		c.Set("org", "github")
		c.Set("user", u)
		c.Set("settings", ps)
		c.Set("defaultOrgBuildLimit", int32(30))
		c.Set("maxOrgBuildLimit", int32(100))
		database.ToContext(c, db)

		UpdateBuildLimit(c)

		if w.Code != http.StatusCreated {
			t.Errorf("UpdateBuildLimit returned %v, want %v", w.Code, http.StatusCreated)
		}

		var got api.OrgBuildLimit

		if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
			t.Errorf("unable to unmarshal response: %v", err)
		}

		if got.GetBuildLimit() != 50 {
			t.Errorf("UpdateBuildLimit build_limit is %v, want %v", got.GetBuildLimit(), int32(50))
		}

		if got.GetOrg() != "github" {
			t.Errorf("UpdateBuildLimit org is %v, want %v", got.GetOrg(), "github")
		}

		if got.GetUpdatedBy() != "octocat" {
			t.Errorf("UpdateBuildLimit updated_by is %v, want %v", got.GetUpdatedBy(), "octocat")
		}
	})

	t.Run("update", func(t *testing.T) {
		db, err := database.NewTest()
		if err != nil {
			t.Errorf("unable to create test database engine: %v", err)
		}

		defer db.Close()

		existing := new(api.OrgBuildLimit)
		existing.SetOrg("github")
		existing.SetBuildLimit(30)
		existing.SetCreatedAt(1)
		existing.SetUpdatedAt(1)
		existing.SetUpdatedBy("admin")

		_, err = db.CreateOrgBuildLimit(context.TODO(), existing)
		if err != nil {
			t.Errorf("unable to create test org build limit: %v", err)
		}

		ps := new(settings.Platform)
		ps.SetEnableOrgBuildLimit(true)

		u := new(api.User)
		u.SetName("octocat")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodPut, "/api/v1/repos/github/limit", strings.NewReader(`{"build_limit": 75}`))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("logger", logrus.NewEntry(logrus.StandardLogger()))
		c.Set("org", "github")
		c.Set("user", u)
		c.Set("settings", ps)
		c.Set("defaultOrgBuildLimit", int32(30))
		c.Set("maxOrgBuildLimit", int32(100))
		database.ToContext(c, db)

		UpdateBuildLimit(c)

		if w.Code != http.StatusOK {
			t.Errorf("UpdateBuildLimit returned %v, want %v", w.Code, http.StatusOK)
		}

		var got api.OrgBuildLimit

		if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
			t.Errorf("unable to unmarshal response: %v", err)
		}

		if got.GetBuildLimit() != 75 {
			t.Errorf("UpdateBuildLimit build_limit is %v, want %v", got.GetBuildLimit(), int32(75))
		}

		if got.GetUpdatedBy() != "octocat" {
			t.Errorf("UpdateBuildLimit updated_by is %v, want %v", got.GetUpdatedBy(), "octocat")
		}
	})

	t.Run("clamp-to-max", func(t *testing.T) {
		db, err := database.NewTest()
		if err != nil {
			t.Errorf("unable to create test database engine: %v", err)
		}

		defer db.Close()

		ps := new(settings.Platform)
		ps.SetEnableOrgBuildLimit(true)

		u := new(api.User)
		u.SetName("octocat")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodPut, "/api/v1/repos/github/limit", strings.NewReader(`{"build_limit": 200}`))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("logger", logrus.NewEntry(logrus.StandardLogger()))
		c.Set("org", "github")
		c.Set("user", u)
		c.Set("settings", ps)
		c.Set("defaultOrgBuildLimit", int32(30))
		c.Set("maxOrgBuildLimit", int32(100))
		database.ToContext(c, db)

		UpdateBuildLimit(c)

		if w.Code != http.StatusCreated {
			t.Errorf("UpdateBuildLimit returned %v, want %v", w.Code, http.StatusCreated)
		}

		var got api.OrgBuildLimit

		if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
			t.Errorf("unable to unmarshal response: %v", err)
		}

		if got.GetBuildLimit() != 100 {
			t.Errorf("UpdateBuildLimit build_limit is %v, want %v", got.GetBuildLimit(), int32(100))
		}
	})

	t.Run("clamp-to-default", func(t *testing.T) {
		db, err := database.NewTest()
		if err != nil {
			t.Errorf("unable to create test database engine: %v", err)
		}

		defer db.Close()

		ps := new(settings.Platform)
		ps.SetEnableOrgBuildLimit(true)

		u := new(api.User)
		u.SetName("octocat")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodPut, "/api/v1/repos/github/limit", strings.NewReader(`{"build_limit": 0}`))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("logger", logrus.NewEntry(logrus.StandardLogger()))
		c.Set("org", "github")
		c.Set("user", u)
		c.Set("settings", ps)
		c.Set("defaultOrgBuildLimit", int32(30))
		c.Set("maxOrgBuildLimit", int32(100))
		database.ToContext(c, db)

		UpdateBuildLimit(c)

		if w.Code != http.StatusCreated {
			t.Errorf("UpdateBuildLimit returned %v, want %v", w.Code, http.StatusCreated)
		}

		var got api.OrgBuildLimit

		if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
			t.Errorf("unable to unmarshal response: %v", err)
		}

		if got.GetBuildLimit() != 30 {
			t.Errorf("UpdateBuildLimit build_limit is %v, want %v", got.GetBuildLimit(), int32(30))
		}
	})

	t.Run("database-get-error", func(t *testing.T) {
		db, err := database.NewTest()
		if err != nil {
			t.Errorf("unable to create test database engine: %v", err)
		}

		db.Close()

		ps := new(settings.Platform)
		ps.SetEnableOrgBuildLimit(true)

		u := new(api.User)
		u.SetName("octocat")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodPut, "/api/v1/repos/github/limit", strings.NewReader(`{"build_limit": 50}`))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("logger", logrus.NewEntry(logrus.StandardLogger()))
		c.Set("org", "github")
		c.Set("user", u)
		c.Set("settings", ps)
		c.Set("defaultOrgBuildLimit", int32(30))
		c.Set("maxOrgBuildLimit", int32(100))
		database.ToContext(c, db)

		UpdateBuildLimit(c)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("UpdateBuildLimit returned %v, want %v", w.Code, http.StatusInternalServerError)
		}
	})

	t.Run("database-create-error", func(t *testing.T) {
		db, err := database.NewTest()
		if err != nil {
			t.Errorf("unable to create test database engine: %v", err)
		}

		defer db.Close()

		mock := &updateMockDB{
			Interface: db,
			getErr:    gorm.ErrRecordNotFound,
			createErr: fmt.Errorf("create failed"),
		}

		ps := new(settings.Platform)
		ps.SetEnableOrgBuildLimit(true)

		u := new(api.User)
		u.SetName("octocat")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodPut, "/api/v1/repos/github/limit", strings.NewReader(`{"build_limit": 50}`))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("logger", logrus.NewEntry(logrus.StandardLogger()))
		c.Set("org", "github")
		c.Set("user", u)
		c.Set("settings", ps)
		c.Set("defaultOrgBuildLimit", int32(30))
		c.Set("maxOrgBuildLimit", int32(100))
		database.ToContext(c, mock)

		UpdateBuildLimit(c)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("UpdateBuildLimit returned %v, want %v", w.Code, http.StatusInternalServerError)
		}
	})

	t.Run("database-update-error", func(t *testing.T) {
		db, err := database.NewTest()
		if err != nil {
			t.Errorf("unable to create test database engine: %v", err)
		}

		defer db.Close()

		existing := new(api.OrgBuildLimit)
		existing.SetID(1)
		existing.SetOrg("github")
		existing.SetBuildLimit(30)
		existing.SetCreatedAt(1)
		existing.SetUpdatedAt(1)
		existing.SetUpdatedBy("admin")

		mock := &updateMockDB{
			Interface: db,
			getResult: existing,
			updateErr: fmt.Errorf("update failed"),
		}

		ps := new(settings.Platform)
		ps.SetEnableOrgBuildLimit(true)

		u := new(api.User)
		u.SetName("octocat")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequestWithContext(t.Context(), http.MethodPut, "/api/v1/repos/github/limit", strings.NewReader(`{"build_limit": 75}`))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("logger", logrus.NewEntry(logrus.StandardLogger()))
		c.Set("org", "github")
		c.Set("user", u)
		c.Set("settings", ps)
		c.Set("defaultOrgBuildLimit", int32(30))
		c.Set("maxOrgBuildLimit", int32(100))
		database.ToContext(c, mock)

		UpdateBuildLimit(c)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("UpdateBuildLimit returned %v, want %v", w.Code, http.StatusInternalServerError)
		}
	})
}
