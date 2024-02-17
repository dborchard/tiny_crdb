package username

// SQLUsername represents a username valid inside SQL.
//
// Note that SQL usernames are not just ASCII names: they can start
// with digits or contain only digits; they can contain certain
// punctuation, and they can contain non-ASCII unicode letters.
// For example, "123.-456" is a valid username.
// Therefore, care must be taken when assembling a string from a
// username for use in other contexts, e.g. to generate filenames:
// some escaping and/or quoting is likely necessary.
//
// Additionally, beware that usernames as manipulated client-side (in
// client drivers, in CLI commands) may not be the same as
// server-side; this is because usernames can be substituted during
// authentication. Additional care must be taken when deriving
// server-side strings in client code. It is always better to add an
// API server-side to assemble the string safely on the client's
// behalf.
//
// This datatype is more complex to a simple string so as to force
// usages to clarify when it is converted to/from strings.
// This complexity is necessary because in CockroachDB SQL, unlike in
// PostgreSQL, SQL usernames are case-folded and NFC-normalized when a
// user logs in, or when used as input to certain CLI commands or SQL
// statements. Then, "inside" CockroachDB, username strings are
// considered pre-normalized and can be used directly for comparisons,
// lookup etc.
//
//   - The constructor MakeSQLUsernameFromUserInput() creates
//     a username from "external input".
//
//   - The constructor MakeSQLUsernameFromPreNormalizedString()
//     creates a username when the caller can guarantee that
//     the input is already pre-normalized.
//
// For convenience, the SQLIdentifier() method also represents a
// username in the form suitable for input back by the SQL parser.
type SQLUsername struct {
	u string
}
