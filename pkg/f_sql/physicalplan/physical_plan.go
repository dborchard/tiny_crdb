package physicalplan

// PhysicalPlan represents a network of processors and streams along with
// information about the results output by this network. The results come from
// unconnected output routers of a subset of processors; all these routers
// output the same kind of data (same schema).
type PhysicalPlan struct {
}
