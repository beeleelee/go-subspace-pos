package chiapos

type Table interface {
	YS() []uint32
	Position(pos uint32) []uint32
	Metadata(pos uint32) []byte
}
