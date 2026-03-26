package draw

import (
	"strconv"
	"time"
)

// Metadata is the metadata of a draw.
type Metadata struct {
	Date           time.Time
	ForclosureDate time.Time
	Version        Version
	DrawType       Type
	Day            Day
	Currency       Currency

	FDJID       string
	ID          string
	TirageOrder int32
	OldType     bool
}

// setID generate a unique ID for the metadata.
func (m *Metadata) setID(suffix string) {
	m.ID = m.FDJID + strconv.Itoa(int(m.TirageOrder)) + suffix
}
