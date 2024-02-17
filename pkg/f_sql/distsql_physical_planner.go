package sql

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/physicalplan"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/physicalplan/replicaoracle"
	kv "github.com/dborchard/tiny_crdb/pkg/g_kv"
	roachpb "github.com/dborchard/tiny_crdb/pkg/g_roachpb"
)

// DistributionType is an enum defining when a plan should be distributed.
type DistributionType int

const (
	// DistributionTypeNone does not distribute a plan across multiple instances.
	DistributionTypeNone = iota
	// DistributionTypeAlways distributes a plan across multiple instances whether
	// it is a system tenant or non-system tenant.
	DistributionTypeAlways
	// DistributionTypeSystemTenantOnly only distributes a plan if it is for a
	// system tenant. Plans on non-system tenants are not distributed.
	DistributionTypeSystemTenantOnly
)

// DistSQLPlanner is used to generate distributed plans from logical
// plans. A rough overview of the process:
//
//   - the plan is based on a planNode tree (in the future it will be based on an
//     intermediate representation tree). Only a subset of the possible trees is
//     supported (this can be checked via CheckSupport).
//
//   - we generate a PhysicalPlan for the planNode tree recursively. The
//     PhysicalPlan consists of a network of processors and streams, with a set
//     of unconnected "result routers". The PhysicalPlan also has information on
//     ordering and on the mapping planNode columns to columns in the result
//     streams (all result routers output streams with the same schema).
//
//     The PhysicalPlan for a scanNode leaf consists of TableReaders, one for each node
//     that has one or more ranges.
//
//   - for each an internal planNode we start with the plan of the child node(s)
//     and add processing stages (connected to the result routers of the children
//     node).
type DistSQLPlanner struct {
}

// NewPlanningCtx returns a new PlanningCtx. When distribute is false, a
// lightweight version PlanningCtx is returned that can be used when the caller
// knows plans will only be run on one node. On SQL tenants, the plan is only
// distributed if tenantDistributionEnabled is true. planner argument can be
// left nil.
func (dsp *DistSQLPlanner) NewPlanningCtx(
	ctx context.Context,
	evalCtx *extendedEvalContext,
	planner *planner,
	txn *kv.Txn,
	distributionType DistributionType,
) *PlanningCtx {
	return dsp.NewPlanningCtxWithOracle(
		ctx, evalCtx, planner, txn, distributionType, physicalplan.DefaultReplicaChooser, roachpb.Locality{},
	)
}

// PlanningCtx contains data used and updated throughout the planning process of
// a single query.
type PlanningCtx struct {
	ExtendedEvalCtx *extendedEvalContext
}

// NewPlanningCtxWithOracle is a variant of NewPlanningCtx that allows passing a
// replica choice oracle as well.
func (dsp *DistSQLPlanner) NewPlanningCtxWithOracle(
	ctx context.Context,
	evalCtx *extendedEvalContext,
	planner *planner,
	txn *kv.Txn,
	distributionType DistributionType,
	oracle replicaoracle.Oracle,
	localityFiler roachpb.Locality,
) *PlanningCtx {
	return &PlanningCtx{
		ExtendedEvalCtx: evalCtx,
	}
}

// PhysicalPlan is a partial physical plan which corresponds to a planNode
// (partial in that it can correspond to a planNode subtree and not necessarily
// to the entire planNode for a given query).
//
// It augments physicalplan.PhysicalPlan with information relating the physical
// plan to a planNode subtree.
//
// These plans are built recursively on a planNode tree.
type PhysicalPlan struct {
	physicalplan.PhysicalPlan

	// PlanToStreamColMap maps planNode columns (see planColumns()) to columns in
	// the result streams. These stream indices correspond to the streams
	// referenced in ResultTypes.
	//
	// Note that in some cases, not all columns in the result streams are
	// referenced in the map; for example, columns that are only required for
	// stream merges in downstream input synchronizers are not included here.
	// (This is due to some processors not being configurable to output only
	// certain columns and will be fixed.)
	//
	// Conversely, in some cases not all planNode columns have a corresponding
	// result stream column (these map to index -1); this is the case for scanNode
	// and indexJoinNode where not all columns in the table are actually used in
	// the plan, but are kept for possible use downstream (e.g., sorting).
	//
	// Before the query is run, the physical plan must be finalized, and during
	// the finalization a projection is added to the plan so that
	// DistSQLReceiver gets rows of the desired schema from the output
	// processor.
	PlanToStreamColMap []int
}
