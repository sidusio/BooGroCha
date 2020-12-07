package timeedit

const chalmersCovid = "chalmers_covid"
const chalmers = "chalmers"

type version int

const (
	VersionChalmers version = iota
	VersionChalmersCovid
)

func (v version) String() string {
	return [...]string{chalmers, chalmersCovid}[v]
}
