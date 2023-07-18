//go:generate go run github.com/dmarkham/enumer -type=Status

package module

type Status uint16

const (
	starting Status = iota
	running
	failed
)
