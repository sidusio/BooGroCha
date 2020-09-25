package timeedit

const chalmersTest = "chalmers_test"
const chalmers = "chalmers"

type version int

const (
	VersionChalmers version = iota
	VersionChalmersTest
)

func (v version) String() string {
	return [...]string{chalmers, chalmersTest}[v]
}
