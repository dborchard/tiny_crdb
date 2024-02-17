package physicalplan

import "github.com/dborchard/tiny_crdb/pkg/f_sql/physicalplan/replicaoracle"

// DefaultReplicaChooser is a nil replicaoracle.Oracle which can be passed in
// place of a replica oracle to some APIs to indicate they can use their default
// replica oracle.
var DefaultReplicaChooser replicaoracle.Oracle
