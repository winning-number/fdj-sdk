package loto

// UUIDs for the different lotto history datasets from the FDJ API.
// Used in the URL to fetch the source.
const (
	GrandLotoUUID       DatasetIdentifier = "1a2b3c4d-9876-4562-b3fc-2c963f66afg6"
	GrandLotoNoelUUID   DatasetIdentifier = "1a2b3c4d-9876-4562-b3fc-2c963f66aff6"
	SuperLoto199605UUID DatasetIdentifier = "1a2b3c4d-9876-4562-b3fc-2c963f66afh6"
	SuperLoto200810UUID DatasetIdentifier = "1a2b3c4d-9876-4562-b3fc-2c963f66afi6"
	SuperLoto201703UUID DatasetIdentifier = "1a2b3c4d-9876-4562-b3fc-2c963f66afj6"
	SuperLoto201907UUID DatasetIdentifier = "1a2b3c4d-9876-4562-b3fc-2c963f66afk6"
	Loto197605UUID      DatasetIdentifier = "1a2b3c4d-9876-4562-b3fc-2c963f66afl6"
	Loto200810UUID      DatasetIdentifier = "1a2b3c4d-9876-4562-b3fc-2c963f66afm6"
	Loto201703UUID      DatasetIdentifier = "1a2b3c4d-9876-4562-b3fc-2c963f66afn6"
	Loto201902UUID      DatasetIdentifier = "1a2b3c4d-9876-4562-b3fc-2c963f66afo6"
	Loto201911UUID      DatasetIdentifier = "1a2b3c4d-9876-4562-b3fc-2c963f66afp6"
)

// Datasets define a base of know history files which should be parsed.
var Datasets = map[DatasetIdentifier]DatasetInfo{
	GrandLotoUUID:       {Type: GrandLotto, Version: LotteryV3},
	GrandLotoNoelUUID:   {Type: XmasLotto, Version: LotteryV3},
	SuperLoto199605UUID: {Type: SuperLotto, Version: LotteryV0},
	SuperLoto200810UUID: {Type: SuperLotto, Version: LotteryV2},
	SuperLoto201703UUID: {Type: SuperLotto, Version: LotteryV3},
	SuperLoto201907UUID: {Type: SuperLotto, Version: LotteryV3},
	Loto197605UUID:      {Type: NewLotto, Version: LotteryV1},
	Loto200810UUID:      {Type: NewLotto, Version: LotteryV2},
	Loto201703UUID:      {Type: NewLotto, Version: LotteryV3},
	Loto201902UUID:      {Type: NewLotto, Version: LotteryV3},
	Loto201911UUID:      {Type: NewLotto, Version: LotteryV4},
}

// DatasetIdentifier is the url path for a specific history file.
type DatasetIdentifier string

// DatasetInfo define a type and version to target the rules to used for decode a csv history file.
type DatasetInfo struct {
	Type    LottoType
	Version LotteryVersion
}

// DatasetIdentifiers return all history files identifiers.
func DatasetIdentifiers() []DatasetIdentifier {
	return []DatasetIdentifier{
		GrandLotoUUID,
		GrandLotoNoelUUID,
		SuperLoto199605UUID,
		SuperLoto200810UUID,
		SuperLoto201703UUID,
		SuperLoto201907UUID,
		Loto197605UUID,
		Loto200810UUID,
		Loto201703UUID,
		Loto201902UUID,
		Loto201911UUID,
	}
}
