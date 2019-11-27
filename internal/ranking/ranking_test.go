package ranking

import (
	"github.com/stretchr/testify/assert"
	"sidus.io/boogrocha/internal/booking"
	"testing"
)

func TestSort(t *testing.T) {
	r1 := booking.AvailableRoom{
		Provider: "A",
		Name:       "A",
	}
	r2 := booking.AvailableRoom{
		Provider: "A",
		Name:       "B",
	}
	r3 := booking.AvailableRoom{
		Provider: "B",
		Name:       "C",
	}
	r4 := booking.AvailableRoom{
		Provider: "B",
		Name:       "D",
	}
	r5 := booking.AvailableRoom{
		Provider: "B",
		Name:       "E",
	}
	
	ranking := Rankings{
		r1: 2,
		r4: 2,
		r2: 1,
		r3: 1,
	}

	sorted := ranking.Sort([]booking.AvailableRoom{r1, r4, r2, r3, r5})
	assert.Equal(t, len(sorted), 5, "Length should be conserved when sorting")

	// Make sure list was sorted with lowest ranking first and reverse alphabetic order
	assert.Equal(t, sorted[0], r5, "Failed to sort")
	assert.Equal(t, sorted[1], r3, "Failed to sort")
	assert.Equal(t, sorted[2], r2, "Failed to sort")
	assert.Equal(t, sorted[3], r4, "Failed to sort")
	assert.Equal(t, sorted[4], r1, "Failed to sort")

}

func TestUpdate(t *testing.T) {
	r1 := booking.AvailableRoom{
		Provider: "A",
		Name:       "Q",
	}
	r2 := booking.AvailableRoom{
		Provider: "A",
		Name:       "R",
	}
	r3 := booking.AvailableRoom{
		Provider: "B",
		Name:       "A",
	}
	r4 := booking.AvailableRoom{
		Provider: "B",
		Name:       "W",
	}
	r5 := booking.AvailableRoom{
		Provider: "B",
		Name:       "E",
	}
	
	ranking := Rankings{
		r1: 5,
		r2: 0,
		r3: 1,
	}

	ranking.Update(r3, []booking.AvailableRoom{r1, r4, r5})

	assert.Equal(t, ranking[r3], uint64(1), "Selected rooms ranking should not be effected")
	assert.Equal(t, ranking[r2], uint64(0), "Elements outside of pool should not be effected")
	assert.True(t, ranking[r4] > (ranking[r1]-5), "Lesser elements should not be effected as mush as greater elements")
	assert.True(t, ranking[r5] > 0, "Not selected elements should be punished")
}
