package tree

// RoutineName represent a function name in a UDF relevant statement, either
// DDL or DML statement. Similar to TableName, it is constructed for incoming
// SQL queries from an UnresolvedObjectName.
type RoutineName struct {
	objName
}

// Schema retrieves the unqualified schema name.
func (r *RoutineName) Schema() string {
	return string(r.SchemaName)
}

func (o *objName) Object() string {
	return string(o.ObjectName)
}
