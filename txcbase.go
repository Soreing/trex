package trex

import "fmt"

// Basic functionality for a Transaction Context
type TxContextBase struct {
	trcprnt string
	tid     string
	pid     string
	rid     string
	flg     string
}

func (b *TxContextBase) IsParent() bool {
	return b.pid == ""
}

func (b *TxContextBase) GetTraceparent() string {
	return b.trcprnt
}

func (b *TxContextBase) GetTransactionId() string {
	return b.tid
}

func (b *TxContextBase) GetParentId() string {
	return b.pid
}

func (b *TxContextBase) GetResourceId() string {
	return b.rid
}

func (b *TxContextBase) GetTraceFlags() string {
	return b.flg

}

// Basic constructor for TxContextBase
func NewTxContextBase(
	ver string,
	tid string,
	pid string,
	rid string,
	flg string,
) TxContextBase {
	return TxContextBase{
		trcprnt: fmt.Sprintf("%s-%s-%s-%s", ver, tid, rid, flg),
		tid:     tid,
		pid:     pid,
		rid:     rid,
		flg:     flg,
	}
}
