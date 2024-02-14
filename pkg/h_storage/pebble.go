package storage

import (
	"context"
	roachpb "github.com/dborchard/tiny_crdb/pkg/f_roachpb"
)

var _ Engine = &Pebble{}

// Pebble is a wrapper around a Pebble database instance.
type Pebble struct {
}

func NewPebble(ctx context.Context, cfg engineConfig) (p *Pebble, err error) {
	return &Pebble{}, nil
}

func (p *Pebble) Close() {
	//TODO implement me
	panic("implement me")
}

func (p *Pebble) Closed() bool {
	//TODO implement me
	panic("implement me")
}

func (p *Pebble) MVCCIterate(ctx context.Context, start, end roachpb.Key, iterKind MVCCIterKind, keyTypes IterKeyType, readCategory ReadCategory, f func(MVCCKeyValue, MVCCRangeKeyStack) error) error {
	//TODO implement me
	panic("implement me")
}

func (p *Pebble) NewMVCCIterator(ctx context.Context, iterKind MVCCIterKind, opts IterOptions) (MVCCIterator, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Pebble) PutMVCC(key MVCCKey, value MVCCValue) error {
	//TODO implement me
	panic("implement me")
}

func (p *Pebble) BufferedSize() int {
	//TODO implement me
	panic("implement me")
}

func (p *Pebble) Compact() error {
	//TODO implement me
	panic("implement me")
}

func (p *Pebble) Flush() error {
	//TODO implement me
	panic("implement me")
}
