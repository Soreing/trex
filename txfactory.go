package trex

// Describes the implementation of a TxFactory
// TxFactories generate transaction contexts for tracking
type TxFactory interface {
	Generate(
		ver string,
		tid string,
		pid string,
		rid string,
		flg string,
	) (interface{}, error)
}
