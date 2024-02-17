package flowinfra

import (
	"context"
	kv "github.com/dborchard/tiny_crdb/pkg/g_kv"
)

type Flow interface {
	// SetTxn is used to provide the transaction in which the flow will run.
	// It needs to be called after Setup() and before Start/Run.
	SetTxn(*kv.Txn)

	// Start starts the flow. Processors run asynchronously in their own
	// goroutines. Wait() needs to be called to wait for the flow to finish.
	// See Run() for a synchronous version.
	//
	// If errors are encountered during the setup part, they're returned.
	Start(context.Context) error
	// Run runs the flow to completion. The last processor is run in the current
	// goroutine; others may run in different goroutines depending on how the
	// flow was configured.
	//
	// f.Wait() is called internally, so the call blocks until all the flow's
	// goroutines are done.
	//
	// It is assumed that rowSyncFlowConsumer is set, so all errors encountered
	// when running this flow are sent to it.
	//
	// noWait is set true when the flow is bound to a pausable portal. With it set,
	// the function returns without waiting the all goroutines to finish. For a
	// pausable portal we will persist this flow and reuse it when re-executing
	// the portal. The flow will be cleaned when the portal is closed, rather than
	// when each portal execution finishes.
	//
	// The caller needs to call f.Cleanup().
	Run(ctx context.Context, noWait bool)

	// Wait waits for all the goroutines for this flow to exit. If the context gets
	// canceled before all goroutines exit, it calls f.cancel().
	Wait()

	// IsLocal returns whether this flow is being run as part of a local-only
	// query.
	IsLocal() bool

	// StatementSQL is the SQL statement for which this flow is executing. It is
	// populated on a best effort basis (only available for user-issued queries
	// that are also not like BulkIO/CDC related).
	StatementSQL() string

	// Cancel cancels the flow by canceling its context. Safe to be called from
	// any goroutine.
	Cancel()
}
