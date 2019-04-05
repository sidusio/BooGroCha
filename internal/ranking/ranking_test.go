package ranking

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSort(t *testing.T) {
	ranking := Rankings{
		"a": 2,
		"d": 2,
		"b": 1,
		"c": 1,
	}

	sorted := ranking.Sort([]string{"a", "d", "b", "c", "e"})
	assert.Equal(t, len(sorted), 5, "Length should be conserved when sorting")

	// Make sure list was sorted with lowest ranking first and reverse alphabetic order
	assert.Equal(t, sorted[0], "e", "Failed to sort")
	assert.Equal(t, sorted[1], "c", "Failed to sort")
	assert.Equal(t, sorted[2], "b", "Failed to sort")
	assert.Equal(t, sorted[3], "d", "Failed to sort")
	assert.Equal(t, sorted[4], "a", "Failed to sort")

}

func TestUpdate(t *testing.T) {
	ranking := Rankings{
		"q": 5,
		"r": 0,
		"a": 1,
	}

	ranking.Update("a", []string{"q", "w", "e"})

	assert.Equal(t, ranking["a"], uint64(1), "Selected rooms ranking should not be effected")
	assert.Equal(t, ranking["r"], uint64(0), "Elements outside of pool should not be effected")
	assert.True(t, ranking["w"] > (ranking["q"]-5), "Lesser elements should not be effected as mush as greater elements")
	assert.True(t, ranking["e"] > 0, "Not selected elements should be punished")
}
