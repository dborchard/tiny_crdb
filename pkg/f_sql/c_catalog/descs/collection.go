package descs

// Collection is a collection of descriptors held by a single session that
// serves SQL requests, or a background job using descriptors. The
// collection is cleared using ReleaseAll() which is called at the
// end of each transaction on the session, or on hitting conditions such
// as errors, or retries that result in transaction timestamp changes.
//
// TODO(ajwerner): Remove the txn argument from the Collection by more tightly
// binding a collection to a *kv.Txn.
type Collection struct {
}
