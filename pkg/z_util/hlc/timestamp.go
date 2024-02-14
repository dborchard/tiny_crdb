package hlc

type Timestamp struct {
}

// ClockTimestamp is a Timestamp with the added capability of being able to
// update a peer's HLC clock. It possesses this capability because the clock
// timestamp itself is guaranteed to have come from an HLC clock somewhere in
// the system. As such, a clock timestamp is a promise that some node in the
// system has a clock with a reading equal to or above its value.
type ClockTimestamp Timestamp

func (t ClockTimestamp) ToTimestamp() Timestamp {
	return Timestamp{}
}

// NowAsClockTimestamp is like Now, but returns a ClockTimestamp instead
// of a raw Timestamp.
//
// This is the counterpart of Update, which is passed a ClockTimestamp
// received from another member of the distributed network. As such,
// callers that intend to use the returned timestamp to update a peer's
// HLC clock should use this method.
func (c *Clock) NowAsClockTimestamp() ClockTimestamp {
	return ClockTimestamp{}
}
