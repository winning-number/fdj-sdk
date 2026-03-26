//nolint:dupl // The package source_* are similar but not identical, and the duplication is acceptable in this case.
package lotto

import "github.com/winning-number/fdj-sdk-lotto/draw"

// APIV3BasePath is the base endpoint for the FDJ history service (API V3).
const (
	APIV3BasePath = "https://www.sto.api.fdj.fr/anonymous/service-draw-info/v3/documentations"
)

// UUIDs for the different lotto history datasets from the FDJ API.
// Used in the URL to fetch the source.
const (
	GrandLotoV3UUID       = "1a2b3c4d-9876-4562-b3fc-2c963f66afg6"
	GrandLotoNoelV3UUID   = "1a2b3c4d-9876-4562-b3fc-2c963f66aff6"
	SuperLoto199605V3UUID = "1a2b3c4d-9876-4562-b3fc-2c963f66afh6"
	SuperLoto200810V3UUID = "1a2b3c4d-9876-4562-b3fc-2c963f66afi6"
	SuperLoto201703V3UUID = "1a2b3c4d-9876-4562-b3fc-2c963f66afj6"
	SuperLoto201907V3UUID = "1a2b3c4d-9876-4562-b3fc-2c963f66afk6"
	Loto197605V3UUID      = "1a2b3c4d-9876-4562-b3fc-2c963f66afl6"
	Loto200810V3UUID      = "1a2b3c4d-9876-4562-b3fc-2c963f66afm6"
	Loto201703V3UUID      = "1a2b3c4d-9876-4562-b3fc-2c963f66afn6"
	Loto201902V3UUID      = "1a2b3c4d-9876-4562-b3fc-2c963f66afo6"
	Loto201911V3UUID      = "1a2b3c4d-9876-4562-b3fc-2c963f66afp6"
)

var sourceInfoMapV3 = map[Source]SourceInfo{
	GrandLoto: {
		APIBase: APIV3BasePath,
		APIPath: GrandLotoV3UUID,
		Type:    draw.GrandLottoType,
		Version: draw.V3,
		Name:    GrandLoto,
	},
	GrandLotoNoel: {
		APIBase: APIV3BasePath,
		APIPath: GrandLotoNoelV3UUID,
		Type:    draw.XmasLottoType,
		Version: draw.V3,
		Name:    GrandLotoNoel,
	},
	SuperLoto199605: {
		APIBase: APIV3BasePath,
		APIPath: SuperLoto199605V3UUID,
		Type:    draw.SuperLottoType,
		Version: draw.V0,
		Name:    SuperLoto199605,
	},
	SuperLoto200810: {
		APIBase: APIV3BasePath,
		APIPath: SuperLoto200810V3UUID,
		Type:    draw.SuperLottoType,
		Version: draw.V2,
		Name:    SuperLoto200810,
	},
	SuperLoto201703: {
		APIBase: APIV3BasePath,
		APIPath: SuperLoto201703V3UUID,
		Type:    draw.SuperLottoType,
		Version: draw.V3,
		Name:    SuperLoto201703,
	},
	SuperLoto201907: {
		APIBase: APIV3BasePath,
		APIPath: SuperLoto201907V3UUID,
		Type:    draw.SuperLottoType,
		Version: draw.V3,
		Name:    SuperLoto201907,
	},
	Loto197605: {
		APIBase: APIV3BasePath,
		APIPath: Loto197605V3UUID,
		Type:    draw.LottoType,
		Version: draw.V1,
		Name:    Loto197605,
	},
	Loto200810: {
		APIBase: APIV3BasePath,
		APIPath: Loto200810V3UUID,
		Type:    draw.LottoType,
		Version: draw.V2,
		Name:    Loto200810,
	},
	Loto201703: {
		APIBase: APIV3BasePath,
		APIPath: Loto201703V3UUID,
		Type:    draw.LottoType,
		Version: draw.V3,
		Name:    Loto201703,
	},
	Loto201902: {
		APIBase: APIV3BasePath,
		APIPath: Loto201902V3UUID,
		Type:    draw.LottoType,
		Version: draw.V3,
		Name:    Loto201902,
	},
	Loto201911: {
		APIBase: APIV3BasePath,
		APIPath: Loto201911V3UUID,
		Type:    draw.LottoType,
		Version: draw.V4,
		Name:    Loto201911,
	},
}
