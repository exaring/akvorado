// SPDX-FileCopyrightText: 2023 Free Mobile
// SPDX-License-Identifier: AGPL-3.0-only

// Package routing fetches routing-related data (AS numbers, AS paths,
// communities). It is modular and accepts several kind of providers (including
// BMP).
package routing

import (
	"context"
	"net/netip"

	"akvorado/common/reporter"
	"akvorado/inlet/routing/provider"
)

// Component represents the metadata compomenent.
type Component struct {
	r        *reporter.Reporter
	provider provider.Provider
}

// Dependencies define the dependencies of the metadata component.
type Dependencies = provider.Dependencies

// New creates a new metadata component.
func New(r *reporter.Reporter, configuration Configuration, dependencies Dependencies) (*Component, error) {
	c := Component{
		r: r,
	}

	// Initialize the provider
	selectedProvider, err := configuration.Provider.Config.New(r, dependencies)
	if err != nil {
		return nil, err
	}
	c.provider = selectedProvider

	return &c, nil
}

// Start starts the routing component.
func (c *Component) Start() error {
	c.r.Info().Msg("starting routing component")
	if starterP, ok := c.provider.(starter); ok {
		if err := starterP.Start(); err != nil {
			return err
		}
	}
	return nil
}

// Stop stops the routing component
func (c *Component) Stop() error {
	c.r.Info().Msg("stopping routing component")
	if stopperP, ok := c.provider.(stopper); ok {
		if err := stopperP.Stop(); err != nil {
			return err
		}
	}
	return nil
}

type starter interface {
	Start() error
}
type stopper interface {
	Stop() error
}

// Lookup uses the selected provider to get an answer. It does not return an
// error, even when the context times out. Instead, it should just returns an
// empty answer.
func (c *Component) Lookup(ctx context.Context, ip netip.Addr, nh netip.Addr) provider.LookupResult {
	return c.provider.Lookup(ctx, ip, nh)
}