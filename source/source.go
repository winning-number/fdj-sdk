// Package source define data(s) relative to the game history files gotten from the FDJ's API.
// It allow to create a game decoder to get the concrete game data type.
package source

import (
	"time"

	"github.com/gofast-pkg/zip"
)

// Source contains the metadata and the data of the source.
type Source struct {
	Metadata

	Data zip.Reader
}

// Metadata contains the metadata of the source.
type Metadata struct {
	Identifier  string
	Size        int64
	RequestedAt time.Time
	FileName    string
}
