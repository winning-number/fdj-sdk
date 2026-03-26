package lotto

import (
	"fmt"

	"github.com/winning-number/fdj-sdk-lotto/draw"
)

// APIVersion is the version used to call fdj's API.
type APIVersion int16

// APIVersion supported.
const (
	APIVersion1 = iota
	APIVersion3
)

// Source represents a lotto history source.
type Source int16

// Source constants for the different lotto history dataset from the FDJ API.
const (
	GrandLoto Source = iota
	GrandLotoNoel
	SuperLoto199605
	SuperLoto200810
	SuperLoto201703
	SuperLoto201907
	Loto197605
	Loto200810
	Loto201703
	Loto201902
	// most recent lotto history.
	Loto201911
)

// SourceInfo contains metadata about a source.
// APIBase is the base URL of the source.
// APIPath is the path to the ZIP file containing the history from the FDJ api.
// Type is the draw type.
// Version is the draw version used for the conversion.
// Name identifies the source.
type SourceInfo struct {
	APIBase string
	APIPath string
	Type    draw.Type
	Version draw.Version
	Name    Source
}

// URL returns the full URL to download the source.
func (s SourceInfo) URL() string {
	return fmt.Sprintf("%s/%s", s.APIBase, s.APIPath)
}

// SourceInfoAll returns all the source info.
func SourceInfoAll(version APIVersion) []SourceInfo {
	sources := SourceAll()
	infos := make([]SourceInfo, len(sources))

	for i, source := range sources {
		infos[i] = GetSourceInfo(source, version)
	}

	return infos
}

// GetSourceInfo returns metadata for the given source.
// If the source is not found, it returns the most recent source (Loto201911, November 2019).
func GetSourceInfo(source Source, version APIVersion) SourceInfo {
	if version == APIVersion3 {
		if info, ok := sourceInfoMapV3[source]; ok {
			return info
		}

		return sourceInfoMapV3[Loto201911]
	}

	if info, ok := sourceInfoMap[source]; ok {
		return info
	}

	return sourceInfoMap[Loto201911]
}

// SourceAll returns all the available sources.
func SourceAll() []Source {
	return []Source{
		GrandLoto,
		GrandLotoNoel,
		SuperLoto199605,
		SuperLoto200810,
		SuperLoto201703,
		SuperLoto201907,
		Loto197605,
		Loto200810,
		Loto201703,
		Loto201902,
		Loto201911,
	}
}
