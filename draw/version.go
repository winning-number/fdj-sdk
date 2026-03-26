package draw

// Version is the version of the draw to parse in adequation with the draw type.
// The version is used to know how to parse the draw.
// V0 is the oldest version and V4 the most recent.
const (
	V0 Version = "v0"
	V1 Version = "v1"
	V2 Version = "v2"
	V3 Version = "v3"
	V4 Version = "v4"
)

// Version is the version of the draw to parse in adequation with the draw type.
type Version string
