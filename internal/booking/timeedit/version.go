package timeedit

const chalmersTest = "chalmers_test"
const chalmers = "chalmers"

type version int

const (
	Chalmers version = iota
	ChalmersTest
)

func (v version) String() string {
	return [...]string{chalmers, chalmersTest}[v]
}
