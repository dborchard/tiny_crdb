package colfetcher

import (
	"context"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/execinfra/execopnode"
	"github.com/dborchard/tiny_crdb/pkg/f_sql/temp/colexecop"
	"github.com/dborchard/tiny_crdb/pkg/y_col/coldata"
	"time"
)

// ScanOperator combines common interfaces between operators that perform KV
// scans, such as ColBatchScan and ColIndexJoin.
type ScanOperator interface {
	colexecop.KVReader
	colexecop.ClosableOperator
}

// ColBatchScan is the colexecop.Operator implementation of TableReader. It
// reads a table from the KV layer, presenting it as coldata.Batches via the
// colexecop.Operator interface.
type ColBatchScan struct {
	*colBatchScanBase
}

// colBatchScanBase is the common base for ColBatchScan and ColBatchDirectScan
// operators.
type colBatchScanBase struct {
}

var _ ScanOperator = &ColBatchScan{}

func (s *colBatchScanBase) GetBytesRead() int64 {
	//TODO implement me
	panic("implement me")
}

func (s *colBatchScanBase) GetKVPairsRead() int64 {
	//TODO implement me
	panic("implement me")
}

func (s *colBatchScanBase) GetRowsRead() int64 {
	//TODO implement me
	panic("implement me")
}

func (s *colBatchScanBase) GetBatchRequestsIssued() int64 {
	//TODO implement me
	panic("implement me")
}

func (s *colBatchScanBase) GetContentionTime() time.Duration {
	//TODO implement me
	panic("implement me")
}

func (s *colBatchScanBase) GetConsumedRU() uint64 {
	//TODO implement me
	panic("implement me")
}

func (s *colBatchScanBase) GetKVCPUTime() time.Duration {
	//TODO implement me
	panic("implement me")
}

func (s ColBatchScan) UsedStreamer() bool {
	//TODO implement me
	panic("implement me")
}

func (s *colBatchScanBase) Init(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (s *colBatchScanBase) Next() coldata.Batch {
	//TODO implement me
	panic("implement me")
}

func (s *colBatchScanBase) ChildCount(verbose bool) int {
	//TODO implement me
	panic("implement me")
}

func (s *colBatchScanBase) Child(nth int, verbose bool) execopnode.OpNode {
	//TODO implement me
	panic("implement me")
}

func (s *colBatchScanBase) Close(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
