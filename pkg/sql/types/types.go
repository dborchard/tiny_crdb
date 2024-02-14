package types

import (
	githubcomlibpqoid "github.com/lib/pq/oid"
)

// T is an instance of a SQL scalar, array, or tuple type. It describes the
// domain of possible values which a column can return, or to which an
// expression can evaluate. The type system does not differentiate between
// nullable and non-nullable types. It is up to the caller to store that
// information separately if it is needed. Here are some example types:
//
//	INT4                     - any 32-bit integer
//	DECIMAL(10, 3)           - any base-10 value with at most 10 digits, with
//	                           up to 3 to right of decimal point
//	FLOAT[]                  - array of 64-bit IEEE 754 floating-point values
//	TUPLE[TIME, VARCHAR(20)] - any pair of values where first value is a time
//	                           of day and the second value is a string having
//	                           up to 20 characters
//
// Fundamentally, a type consists of the following attributes, each of which has
// a corresponding accessor method. Some of these attributes are only defined
// for a subset of types. See the method comments for more details.
//
//	Family        - equivalence group of the type (enumeration)
//	Oid           - Postgres Object ID that describes the type (enumeration)
//	Precision     - maximum accuracy of the type (numeric)
//	Width         - maximum size or scale of the type (numeric)
//	Locale        - location which governs sorting, formatting, etc. (string)
//	ArrayContents - array element type (T)
//	TupleContents - slice of types of each tuple field ([]*T)
//	TupleLabels   - slice of labels of each tuple field ([]string)
//
// Some types are not currently allowed as the type of a column (e.g. nested
// arrays). Other usages of the types package may have similar restrictions.
// Each such caller is responsible for enforcing their own restrictions; it's
// not the concern of the types package.
//
// Implementation-wise, types.T wraps a protobuf-generated InternalType struct.
// The generated protobuf code defines the struct fields, marshals/unmarshals
// them, formats a string representation, etc. Meanwhile, the wrapper types.T
// struct overrides the Marshal/Unmarshal methods in order to map to/from older
// persisted InternalType representations. For example, older versions of
// InternalType (previously called ColumnType) used a VisibleType field to
// represent INT2, whereas newer versions use Width/Oid. Unmarshal upgrades from
// this old format to the new, and Marshal downgrades, thus preserving backwards
// compatibility.
//
// Simple (unary) scalars types
// ----------------------------
//
// | SQL type          | Family         | Oid           | Precision | Width |
// |-------------------|----------------|---------------|-----------|-------|
// | NULL (unknown)    | UNKNOWN        | T_unknown     | 0         | 0     |
// | BOOL              | BOOL           | T_bool        | 0         | 0     |
// | DATE              | DATE           | T_date        | 0         | 0     |
// | TIMESTAMP         | TIMESTAMP      | T_timestamp   | 0         | 0     |
// | INTERVAL          | INTERVAL       | T_interval    | 0         | 0     |
// | TIMESTAMPTZ       | TIMESTAMPTZ    | T_timestamptz | 0         | 0     |
// | OID               | OID            | T_oid         | 0         | 0     |
// | UUID              | UUID           | T_uuid        | 0         | 0     |
// | INET              | INET           | T_inet        | 0         | 0     |
// | TIME              | TIME           | T_time        | 0         | 0     |
// | TIMETZ            | TIMETZ         | T_timetz      | 0         | 0     |
// | JSON              | JSONB          | T_json        | 0         | 0     |
// | JSONB             | JSONB          | T_jsonb       | 0         | 0     |
// |                   |                |               |           |       |
// | BYTES             | BYTES          | T_bytea       | 0         | 0     |
// |                   |                |               |           |       |
// | STRING            | STRING         | T_text        | 0         | 0     |
// | STRING(N)         | STRING         | T_text        | 0         | N     |
// | VARCHAR           | STRING         | T_varchar     | 0         | 0     |
// | VARCHAR(N)        | STRING         | T_varchar     | 0         | N     |
// | CHAR              | STRING         | T_bpchar      | 0         | 1     |
// | CHAR(N)           | STRING         | T_bpchar      | 0         | N     |
// | "char"            | STRING         | T_char        | 0         | 0     |
// | NAME              | STRING         | T_name        | 0         | 0     |
// |                   |                |               |           |       |
// | STRING COLLATE en | COLLATEDSTRING | T_text        | 0         | 0     |
// | STRING(N) COL...  | COLLATEDSTRING | T_text        | 0         | N     |
// | VARCHAR COL...    | COLLATEDSTRING | T_varchar     | 0         | N     |
// | VARCHAR(N) COL... | COLLATEDSTRING | T_varchar     | 0         | N     |
// | CHAR COL...       | COLLATEDSTRING | T_bpchar      | 0         | 1     |
// | CHAR(N) COL...    | COLLATEDSTRING | T_bpchar      | 0         | N     |
// | "char" COL...     | COLLATEDSTRING | T_char        | 0         | 0     |
// |                   |                |               |           |       |
// | DECIMAL           | DECIMAL        | T_decimal     | 0         | 0     |
// | DECIMAL(N)        | DECIMAL        | T_decimal     | N         | 0     |
// | DECIMAL(N,M)      | DECIMAL        | T_decimal     | N         | M     |
// |                   |                |               |           |       |
// | FLOAT8            | FLOAT          | T_float8      | 0         | 0     |
// | FLOAT4            | FLOAT          | T_float4      | 0         | 0     |
// |                   |                |               |           |       |
// | BIT               | BIT            | T_bit         | 0         | 1     |
// | BIT(N)            | BIT            | T_bit         | 0         | N     |
// | VARBIT            | BIT            | T_varbit      | 0         | 0     |
// | VARBIT(N)         | BIT            | T_varbit      | 0         | N     |
// |                   |                |               |           |       |
// | INT,INTEGER       | INT            | T_int8        | 0         | 64    |
// | INT2,SMALLINT     | INT            | T_int2        | 0         | 16    |
// | INT4              | INT            | T_int4        | 0         | 32    |
// | INT8,INT64,BIGINT | INT            | T_int8        | 0         | 64    |
//
// Tuple types
// -----------
//
// These cannot (yet) be used in tables but are used in DistSQL flow
// processors for queries that have tuple-typed intermediate results.
//
// | Field           | Description                                             |
// |-----------------|---------------------------------------------------------|
// | Family          | TupleFamily                                             |
// | Oid             | T_record                                                |
// | TupleContents   | Contains tuple field types (can be recursively defined) |
// | TupleLabels     | Contains labels for each tuple field                    |
//
// Array types
// -----------
//
// | Field           | Description                                             |
// |-----------------|---------------------------------------------------------|
// | Family          | ArrayFamily                                             |
// | Oid             | T__XXX (double underscores), where XXX is the Oid name  |
// |                 | of a scalar type                                        |
// | ArrayContents   | Type of array elements (scalar, array, or tuple)        |
//
// There are two special ARRAY types:
//
// | SQL type          | Family         | Oid           | ArrayContents |
// |-------------------|----------------|---------------|---------------|
// | INT2VECTOR        | ARRAY          | T_int2vector  | Int           |
// | OIDVECTOR         | ARRAY          | T_oidvector   | Oid           |
//
// When these types are themselves made into arrays, the Oids become T__int2vector and
// T__oidvector, respectively.
//
// User defined types
// ------------------
//
// * Enums
// | Field         | Description                                |
// |---------------|--------------------------------------------|
// | Family        | EnumFamily                                 |
// | Oid           | A unique OID generated upon enum creation  |
//
// See types.proto for the corresponding proto definition. Its automatic
// type declaration is suppressed in the proto so that it is possible to
// add additional fields to T without serializing them.
type T struct {
	// InternalType should never be directly referenced outside this package. The
	// only reason it is exported is because gogoproto panics when printing the
	// string representation of an unexported field. This is a problem when this
	// struct is embedded in a larger struct (like a ColumnDescriptor).
	InternalType InternalType
}

type Family int32

const (
	UnknownFamily Family = iota
	IntFamily
	FloatFamily
	StringFamily
	BoolFamily
)

var (
	emptyLocale = ""
)

type InternalType struct {
	Family             Family
	Oid                githubcomlibpqoid.Oid
	Locale             *string
	Width              int32
	Precision          int32
	TimePrecisionIsSet bool
}

// Convenience list of pre-constructed types. Caller code can use any of these
// types, or use the MakeXXX methods to construct a custom type that is not
// listed here (e.g. if a custom width is needed).
var (
	// Unknown is the type of an expression that statically evaluates to NULL.
	// This type should never be returned for an expression that does not *always*
	// evaluate to NULL.
	Unknown = &T{InternalType: InternalType{
		Family: UnknownFamily, Oid: githubcomlibpqoid.T_unknown, Locale: &emptyLocale}}

	// Bool is the type of a boolean true/false value.
	Bool = &T{InternalType: InternalType{
		Family: BoolFamily, Oid: githubcomlibpqoid.T_bool, Locale: &emptyLocale}}

	// Int is the type of a 64-bit signed integer. This is the canonical type
	// for IntFamily.
	Int = &T{InternalType: InternalType{
		Family: IntFamily, Width: 64, Oid: githubcomlibpqoid.T_int8, Locale: &emptyLocale}}

	// Int4 is the type of a 32-bit signed integer.
	Int4 = &T{InternalType: InternalType{
		Family: IntFamily, Width: 32, Oid: githubcomlibpqoid.T_int4, Locale: &emptyLocale}}

	// Int2 is the type of a 16-bit signed integer.
	Int2 = &T{InternalType: InternalType{
		Family: IntFamily, Width: 16, Oid: githubcomlibpqoid.T_int2, Locale: &emptyLocale}}

	// Float is the type of a 64-bit base-2 floating-point number (IEEE 754).
	// This is the canonical type for FloatFamily.
	Float = &T{InternalType: InternalType{
		Family: FloatFamily, Width: 64, Oid: githubcomlibpqoid.T_float8, Locale: &emptyLocale}}

	// Float4 is the type of a 32-bit base-2 floating-point number (IEEE 754).
	Float4 = &T{InternalType: InternalType{
		Family: FloatFamily, Width: 32, Oid: githubcomlibpqoid.T_float4, Locale: &emptyLocale}}

	// String is the type of a Unicode string, with no specified limit on the
	// count of characters. This is the canonical type for StringFamily. It is
	// reported as STRING in SHOW CREATE but "text" in introspection for
	// compatibility with PostgreSQL.
	String = &T{InternalType: InternalType{
		Family: StringFamily, Oid: githubcomlibpqoid.T_text, Locale: &emptyLocale}}

	// VarChar is equivalent to String, but has a differing OID (T_varchar),
	// which makes it show up differently when displayed. It is reported as
	// VARCHAR in SHOW CREATE and "character varying" in introspection for
	// compatibility with PostgreSQL.
	VarChar = &T{InternalType: InternalType{
		Family: StringFamily, Oid: githubcomlibpqoid.T_varchar, Locale: &emptyLocale}}
)
