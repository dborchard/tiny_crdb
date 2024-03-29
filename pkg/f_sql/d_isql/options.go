// Copyright 2022 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package isql

import "github.com/dborchard/tiny_crdb/pkg/f_sql/sessiondata"

type (

	// ExecutorOption configures an Executor.
	ExecutorOption interface{ applyEx(*ExecutorConfig) }

	// TxnOption is used to configure a Txn.
	TxnOption interface{ applyTxn(*TxnConfig) }

	// Option can configure both an Executor and a Txn.
	Option interface {
		TxnOption
		ExecutorOption
	}
)

// ExecutorConfig is the configuration used by the implementation of DB to
// set up the Executor.
type ExecutorConfig struct {
	sessionData *sessiondata.SessionData
}

func (ec *ExecutorConfig) GetSessionData() *sessiondata.SessionData {
	return ec.sessionData
}

// TxnConfig is the config to be set for txn.
type TxnConfig struct {
	ExecutorConfig
	steppingEnabled bool
}

func (tc *TxnConfig) Init(opts ...TxnOption) {
	for _, opt := range opts {
		opt.applyTxn(tc)
	}
}

// Init is used to initialize an ExecutorConfig.
func (ec *ExecutorConfig) Init(opts ...ExecutorOption) {
	for _, o := range opts {
		o.applyEx(ec)
	}
}

type sessionDataOption sessiondata.SessionData

// WithSessionData allows the user to configure the session data for the Txn or
// Executor.
func WithSessionData(sd *sessiondata.SessionData) Option {
	return (*sessionDataOption)(sd)
}

func (o *sessionDataOption) applyEx(cfg *ExecutorConfig) {
	cfg.sessionData = (*sessiondata.SessionData)(o)
}
func (o *sessionDataOption) applyTxn(cfg *TxnConfig) {
	cfg.sessionData = (*sessiondata.SessionData)(o)
}

// SteppingEnabled creates a TxnOption to determine whether the underlying
// transaction should have stepping enabled. If stepping is enabled, the
// transaction will implicitly use lower admission priority. However, the
// user will need to remember to Step the Txn to make writes visible. The
// Executor will automatically (for better or for worse) step the
// transaction when executing each statement.
func SteppingEnabled() TxnOption {
	return steppingEnabled(true)
}

type steppingEnabled bool

// GetSteppingEnabled return the steppingEnabled setting from the txn config.
func (tc *TxnConfig) GetSteppingEnabled() bool {
	return tc.steppingEnabled
}

func (s steppingEnabled) applyTxn(o *TxnConfig) { o.steppingEnabled = bool(s) }
