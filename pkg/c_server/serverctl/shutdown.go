package serverctl

// MakeShutdownRequest constructs a serverctl.ShutdownRequest.
func MakeShutdownRequest(err error) ShutdownRequest {
	return ShutdownRequest{
		Err: err,
	}
}
