package coldataext

import (
	"github.com/dborchard/tiny_crdb/pkg/i_sql/sem/eval"
	"github.com/dborchard/tiny_crdb/pkg/i_sql/sem/tree"
	"github.com/dborchard/tiny_crdb/pkg/i_sql/types"
	"github.com/dborchard/tiny_crdb/pkg/y_col/coldata"
)

// datumVec is a vector of tree.Datums of the same type.
type datumVec struct {
	t       *types.T
	data    []tree.Datum
	scratch []byte

	evalCtx *eval.Context
	da      tree.DatumAlloc
}

var _ coldata.DatumVec = &datumVec{}

// newDatumVec returns a datumVec struct with capacity of n.
func newDatumVec(t *types.T, n int, evalCtx *eval.Context) coldata.DatumVec {
	return &datumVec{
		t:       t,
		data:    make([]tree.Datum, n),
		evalCtx: evalCtx,
	}
}

func (d *datumVec) Get(i int) coldata.Datum {
	//TODO implement me
	panic("implement me")
}

func (d *datumVec) Set(i int, v coldata.Datum) {
	//TODO implement me
	panic("implement me")
}

func (d *datumVec) Window(start, end int) coldata.DatumVec {
	//TODO implement me
	panic("implement me")
}

func (d *datumVec) CopySlice(src coldata.DatumVec, destIdx, srcStartIdx, srcEndIdx int) {
	//TODO implement me
	panic("implement me")
}

func (d *datumVec) AppendSlice(src coldata.DatumVec, destIdx, srcStartIdx, srcEndIdx int) {
	//TODO implement me
	panic("implement me")
}

func (d *datumVec) AppendVal(v coldata.Datum) {
	//TODO implement me
	panic("implement me")
}

func (d *datumVec) SetLength(l int) {
	//TODO implement me
	panic("implement me")
}

func (d *datumVec) Len() int {
	//TODO implement me
	panic("implement me")
}

func (d *datumVec) Cap() int {
	//TODO implement me
	panic("implement me")
}

func (d *datumVec) MarshalAt(appendTo []byte, i int) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (d *datumVec) UnmarshalTo(i int, b []byte) error {
	//TODO implement me
	panic("implement me")
}

func (d *datumVec) Size(startIdx int) int64 {
	//TODO implement me
	panic("implement me")
}

func (d *datumVec) SetEvalCtx(evalCtx interface{}) {
	//TODO implement me
	panic("implement me")
}
