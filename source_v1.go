//nolint:dupl // The package source_* are similar but not identical, and the duplication is acceptable in this case.
package lotto

import "github.com/winning-number/fdj-sdk-lotto/draw"

// APIBasePath is the base endpoint for the FDJ history service (API V1).
const (
	APIBasePath = "https://media.fdj.fr/static/csv/loto"
)

// ZIP file names for the different lotto history datasets from the FDJ API.
// Used in the URL to fetch the source.
const (
	GrandLotoZipName       = "grandloto_201912.zip"
	GrandLotoNoelZipName   = "lotonoel_201703.zip"
	SuperLoto199605ZipName = "superloto_199605.zip"
	SuperLoto200810ZipName = "superloto_200810.zip"
	SuperLoto201703ZipName = "superloto_201703.zip"
	SuperLoto201907ZipName = "superloto_201907.zip"
	Loto197605ZipName      = "loto_197605.zip"
	Loto200810ZipName      = "loto_200810.zip"
	Loto201703ZipName      = "loto_201703.zip"
	Loto201902ZipName      = "loto_201902.zip"
	Loto201911ZipName      = "loto_201911.zip"
)

var sourceInfoMap = map[Source]SourceInfo{
	GrandLoto: {
		APIBase: APIBasePath,
		APIPath: GrandLotoZipName,
		Type:    draw.GrandLottoType,
		Version: draw.V3,
		Name:    GrandLoto,
	},
	GrandLotoNoel: {
		APIBase: APIBasePath,
		APIPath: GrandLotoNoelZipName,
		Type:    draw.XmasLottoType,
		Version: draw.V3,
		Name:    GrandLotoNoel,
	},
	SuperLoto199605: {
		APIBase: APIBasePath,
		APIPath: SuperLoto199605ZipName,
		Type:    draw.SuperLottoType,
		Version: draw.V0,
		Name:    SuperLoto199605,
	},
	SuperLoto200810: {
		APIBase: APIBasePath,
		APIPath: SuperLoto200810ZipName,
		Type:    draw.SuperLottoType,
		Version: draw.V2,
		Name:    SuperLoto200810,
	},
	SuperLoto201703: {
		APIBase: APIBasePath,
		APIPath: SuperLoto201703ZipName,
		Type:    draw.SuperLottoType,
		Version: draw.V3,
		Name:    SuperLoto201703,
	},
	SuperLoto201907: {
		APIBase: APIBasePath,
		APIPath: SuperLoto201907ZipName,
		Type:    draw.SuperLottoType,
		Version: draw.V3,
		Name:    SuperLoto201907,
	},
	Loto197605: {
		APIBase: APIBasePath,
		APIPath: Loto197605ZipName,
		Type:    draw.LottoType,
		Version: draw.V1,
		Name:    Loto197605,
	},
	Loto200810: {
		APIBase: APIBasePath,
		APIPath: Loto200810ZipName,
		Type:    draw.LottoType,
		Version: draw.V2,
		Name:    Loto200810,
	},
	Loto201703: {
		APIBase: APIBasePath,
		APIPath: Loto201703ZipName,
		Type:    draw.LottoType,
		Version: draw.V3,
		Name:    Loto201703,
	},
	Loto201902: {
		APIBase: APIBasePath,
		APIPath: Loto201902ZipName,
		Type:    draw.LottoType,
		Version: draw.V3,
		Name:    Loto201902,
	},
	Loto201911: {
		APIBase: APIBasePath,
		APIPath: Loto201911ZipName,
		Type:    draw.LottoType,
		Version: draw.V4,
		Name:    Loto201911,
	},
}
