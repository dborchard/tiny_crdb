package sql

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/sessiondata"
	"time"
)

// NewInternalSessionData returns a session data for use in internal queries
// that are not run on behalf of a user session, such as those run during the
// steps of background jobs and schema changes. Each session variable is
// initialized using the correct default value.
func NewInternalSessionData(
	ctx context.Context, opName string,
) *sessiondata.SessionData {
	sd := &sessiondata.SessionData{}
	sd.Location = time.UTC
	return sd
}
