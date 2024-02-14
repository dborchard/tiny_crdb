package tree

// FmtCtx is suitable for passing to Format() methods.
// It also exposes the underlying bytes.Buffer interface for
// convenience.
//
// FmtCtx cannot be copied by value.
type FmtCtx struct {
}

// NodeFormatter is implemented by nodes that can be pretty-printed.
type NodeFormatter interface {
	// Format performs pretty-printing towards a bytes buffer. The flags member
	// of ctx influences the results. Most callers should use FormatNode instead.
	Format(ctx *FmtCtx)
}
