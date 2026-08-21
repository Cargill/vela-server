// SPDX-License-Identifier: Apache-2.0

package org

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/go-vela/server/constants"
)

type (
	// config represents the settings required to create the engine that implements the OrgInterface interface.
	config struct {
		// specifies to skip creating tables and indexes for the Org engine
		SkipCreation bool
	}

	// Engine represents the org functionality that implements the OrgInterface interface.
	Engine struct {
		// engine configuration settings used in org functions
		config *config

		ctx context.Context

		// gorm.io/gorm database client used in org functions
		client *gorm.DB

		// sirupsen/logrus logger used in org functions
		logger *logrus.Entry
	}
)

// New creates and returns a Vela service for integrating with orgs in the database.
func New(opts ...EngineOpt) (*Engine, error) {
	e := new(Engine)

	e.client = new(gorm.DB)
	e.config = new(config)
	e.logger = new(logrus.Entry)

	for _, opt := range opts {
		err := opt(e)
		if err != nil {
			return nil, err
		}
	}

	if e.config.SkipCreation {
		e.logger.Warningf("skipping creation of %s table", constants.TableOrganization)

		return e, nil
	}

	err := e.CreateOrganizationTable(e.ctx, e.client.Config.Dialector.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s table: %w", constants.TableOrganization, err)
	}

	return e, nil
}
