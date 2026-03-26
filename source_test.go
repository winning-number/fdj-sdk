package lotto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/winning-number/fdj-sdk-lotto/draw"
)

func TestSourceInfo_URL(t *testing.T) {
	t.Run("Should return the url for the source 'Grand lotto'", func(t *testing.T) {
		expected := "https://media.fdj.fr/static/csv/loto/grandloto_201912.zip"

		got := GetSourceInfo(GrandLoto, APIVersion1).URL()

		assert.Equal(t, expected, got)
	})
}

func TestGetSourceInfo(t *testing.T) {
	t.Run("Should return the source info 'Grand lotto'", func(t *testing.T) {
		expected := SourceInfo{
			APIBase: APIBasePath,
			APIPath: GrandLotoZipName,
			Type:    draw.GrandLottoType,
			Version: draw.V3,
			Name:    GrandLoto,
		}

		got := GetSourceInfo(GrandLoto, APIVersion1)

		assert.Equal(t, expected, got)
	})
	t.Run("Should return the source info 'Grand lotto noel'", func(t *testing.T) {
		expected := SourceInfo{
			APIBase: APIBasePath,
			APIPath: GrandLotoNoelZipName,
			Type:    draw.XmasLottoType,
			Version: draw.V3,
			Name:    GrandLotoNoel,
		}

		got := GetSourceInfo(GrandLotoNoel, APIVersion1)

		assert.Equal(t, expected, got)
	})
	t.Run("Should return the source info 'Super lotto 1996-05'", func(t *testing.T) {
		expected := SourceInfo{
			APIBase: APIBasePath,
			APIPath: SuperLoto199605ZipName,
			Type:    draw.SuperLottoType,
			Version: draw.V0,
			Name:    SuperLoto199605,
		}

		got := GetSourceInfo(SuperLoto199605, APIVersion1)

		assert.Equal(t, expected, got)
	})
	t.Run("Should return the source info 'Super lotto 2008-10'", func(t *testing.T) {
		expected := SourceInfo{
			APIBase: APIBasePath,
			APIPath: SuperLoto200810ZipName,
			Type:    draw.SuperLottoType,
			Version: draw.V2,
			Name:    SuperLoto200810,
		}

		got := GetSourceInfo(SuperLoto200810, APIVersion1)

		assert.Equal(t, expected, got)
	})
	t.Run("Should return the source info 'Super lotto 2017-03'", func(t *testing.T) {
		expected := SourceInfo{
			APIBase: APIBasePath,
			APIPath: SuperLoto201703ZipName,
			Type:    draw.SuperLottoType,
			Version: draw.V3,
			Name:    SuperLoto201703,
		}

		got := GetSourceInfo(SuperLoto201703, APIVersion1)

		assert.Equal(t, expected, got)
	})
	t.Run("Should return the source info 'Super lotto 2019-07'", func(t *testing.T) {
		expected := SourceInfo{
			APIBase: APIBasePath,
			APIPath: SuperLoto201907ZipName,
			Type:    draw.SuperLottoType,
			Version: draw.V3,
			Name:    SuperLoto201907,
		}

		got := GetSourceInfo(SuperLoto201907, APIVersion1)

		assert.Equal(t, expected, got)
	})
	t.Run("Should return the source info 'Lotto 1976-05'", func(t *testing.T) {
		expected := SourceInfo{
			APIBase: APIBasePath,
			APIPath: Loto197605ZipName,
			Type:    draw.LottoType,
			Version: draw.V1,
			Name:    Loto197605,
		}

		got := GetSourceInfo(Loto197605, APIVersion1)

		assert.Equal(t, expected, got)
	})

	t.Run("Should return the source info 'Lotto 2008-10'", func(t *testing.T) {
		expected := SourceInfo{
			APIBase: APIBasePath,
			APIPath: Loto200810ZipName,
			Type:    draw.LottoType,
			Version: draw.V2,
			Name:    Loto200810,
		}

		got := GetSourceInfo(Loto200810, APIVersion1)

		assert.Equal(t, expected, got)
	})
	t.Run("Should return the source info 'Lotto 2017-03'", func(t *testing.T) {
		expected := SourceInfo{
			APIBase: APIBasePath,
			APIPath: Loto201703ZipName,
			Type:    draw.LottoType,
			Version: draw.V3,
			Name:    Loto201703,
		}

		got := GetSourceInfo(Loto201703, APIVersion1)

		assert.Equal(t, expected, got)
	})
	t.Run("Should return the source info 'Lotto 2019-02'", func(t *testing.T) {
		expected := SourceInfo{
			APIBase: APIBasePath,
			APIPath: Loto201902ZipName,
			Type:    draw.LottoType,
			Version: draw.V3,
			Name:    Loto201902,
		}

		got := GetSourceInfo(Loto201902, APIVersion1)

		assert.Equal(t, expected, got)
	})
	t.Run("Should return the source info 'Lotto 20019-11'", func(t *testing.T) {
		expected := SourceInfo{
			APIBase: APIBasePath,
			APIPath: Loto201911ZipName,
			Type:    draw.LottoType,
			Version: draw.V4,
			Name:    Loto201911,
		}

		got := GetSourceInfo(Loto201911, APIVersion1)

		assert.Equal(t, expected, got)
	})
	t.Run("Should return the default source info 'unknow source'", func(t *testing.T) {
		expected := SourceInfo{
			APIBase: APIBasePath,
			APIPath: Loto201911ZipName,
			Type:    draw.LottoType,
			Version: draw.V4,
			Name:    Loto201911,
		}

		got := GetSourceInfo(9999, APIVersion1)

		assert.Equal(t, expected, got)
	})
	t.Run("Should return the source info 'Grand lotto' V3", func(t *testing.T) {
		expected := SourceInfo{
			APIBase: APIV3BasePath,
			APIPath: GrandLotoV3UUID,
			Type:    draw.GrandLottoType,
			Version: draw.V3,
			Name:    GrandLoto,
		}

		got := GetSourceInfo(GrandLoto, APIVersion3)

		assert.Equal(t, expected, got)
	})
	t.Run("Should return the source info 'Grand lotto noel' V3", func(t *testing.T) {
		expected := SourceInfo{
			APIBase: APIV3BasePath,
			APIPath: GrandLotoNoelV3UUID,
			Type:    draw.XmasLottoType,
			Version: draw.V3,
			Name:    GrandLotoNoel,
		}

		got := GetSourceInfo(GrandLotoNoel, APIVersion3)

		assert.Equal(t, expected, got)
	})
	t.Run("Should return the source info 'Super lotto 1996-05' V3", func(t *testing.T) {
		expected := SourceInfo{
			APIBase: APIV3BasePath,
			APIPath: SuperLoto199605V3UUID,
			Type:    draw.SuperLottoType,
			Version: draw.V0,
			Name:    SuperLoto199605,
		}

		got := GetSourceInfo(SuperLoto199605, APIVersion3)

		assert.Equal(t, expected, got)
	})
	t.Run("Should return the source info 'Super lotto 2008-10' V3", func(t *testing.T) {
		expected := SourceInfo{
			APIBase: APIV3BasePath,
			APIPath: SuperLoto200810V3UUID,
			Type:    draw.SuperLottoType,
			Version: draw.V2,
			Name:    SuperLoto200810,
		}

		got := GetSourceInfo(SuperLoto200810, APIVersion3)

		assert.Equal(t, expected, got)
	})
	t.Run("Should return the source info 'Super lotto 2017-03' V3", func(t *testing.T) {
		expected := SourceInfo{
			APIBase: APIV3BasePath,
			APIPath: SuperLoto201703V3UUID,
			Type:    draw.SuperLottoType,
			Version: draw.V3,
			Name:    SuperLoto201703,
		}

		got := GetSourceInfo(SuperLoto201703, APIVersion3)

		assert.Equal(t, expected, got)
	})
	t.Run("Should return the source info 'Super lotto 2019-07' V3", func(t *testing.T) {
		expected := SourceInfo{
			APIBase: APIV3BasePath,
			APIPath: SuperLoto201907V3UUID,
			Type:    draw.SuperLottoType,
			Version: draw.V3,
			Name:    SuperLoto201907,
		}

		got := GetSourceInfo(SuperLoto201907, APIVersion3)

		assert.Equal(t, expected, got)
	})
	t.Run("Should return the source info 'Lotto 1976-05' V3", func(t *testing.T) {
		expected := SourceInfo{
			APIBase: APIV3BasePath,
			APIPath: Loto197605V3UUID,
			Type:    draw.LottoType,
			Version: draw.V1,
			Name:    Loto197605,
		}

		got := GetSourceInfo(Loto197605, APIVersion3)

		assert.Equal(t, expected, got)
	})

	t.Run("Should return the source info 'Lotto 2008-10' V3", func(t *testing.T) {
		expected := SourceInfo{
			APIBase: APIV3BasePath,
			APIPath: Loto200810V3UUID,
			Type:    draw.LottoType,
			Version: draw.V2,
			Name:    Loto200810,
		}

		got := GetSourceInfo(Loto200810, APIVersion3)

		assert.Equal(t, expected, got)
	})
	t.Run("Should return the source info 'Lotto 2017-03' V3", func(t *testing.T) {
		expected := SourceInfo{
			APIBase: APIV3BasePath,
			APIPath: Loto201703V3UUID,
			Type:    draw.LottoType,
			Version: draw.V3,
			Name:    Loto201703,
		}

		got := GetSourceInfo(Loto201703, APIVersion3)

		assert.Equal(t, expected, got)
	})
	t.Run("Should return the source info 'Lotto 2019-02' V3", func(t *testing.T) {
		expected := SourceInfo{
			APIBase: APIV3BasePath,
			APIPath: Loto201902V3UUID,
			Type:    draw.LottoType,
			Version: draw.V3,
			Name:    Loto201902,
		}

		got := GetSourceInfo(Loto201902, APIVersion3)

		assert.Equal(t, expected, got)
	})
	t.Run("Should return the source info 'Lotto 2019-11' V3", func(t *testing.T) {
		expected := SourceInfo{
			APIBase: APIV3BasePath,
			APIPath: Loto201911V3UUID,
			Type:    draw.LottoType,
			Version: draw.V4,
			Name:    Loto201911,
		}

		got := GetSourceInfo(Loto201911, APIVersion3)

		assert.Equal(t, expected, got)
	})
	t.Run("Should return the default source info 'unknow source' V3", func(t *testing.T) {
		expected := SourceInfo{
			APIBase: APIV3BasePath,
			APIPath: Loto201911V3UUID,
			Type:    draw.LottoType,
			Version: draw.V4,
			Name:    Loto201911,
		}

		got := GetSourceInfo(9999, APIVersion3)

		assert.Equal(t, expected, got)
	})
}

func TestSourceInfoAll(t *testing.T) {
	t.Run("Should return all the sources v1", func(t *testing.T) {
		infos := SourceInfoAll(APIVersion1)

		assert.Len(t, infos, 11)
	})
	t.Run("Should return all the sources v3", func(t *testing.T) {
		infos := SourceInfoAll(APIVersion3)

		assert.Len(t, infos, 11)
	})
}

func TestSourceAll(t *testing.T) {
	t.Run("Should return all the sources", func(t *testing.T) {
		expected := []Source{
			GrandLoto, GrandLotoNoel, SuperLoto199605, SuperLoto200810, SuperLoto201703, SuperLoto201907,
			Loto197605, Loto200810, Loto201703, Loto201902, Loto201911,
		}

		got := SourceAll()

		assert.Equal(t, expected, got)
	})
}
