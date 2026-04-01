package loto

import (
	"testing"

	"github.com/gofast-pkg/zip/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winning-number/fdj-sdk/v2/loto/testdata"
	"github.com/winning-number/fdj-sdk/v2/source"
)

const (
	assertMockTimeExpectation = 1
)

func TestNewDecoder(t *testing.T) {
	t.Run("should return en error with an unexpected identifier", func(t *testing.T) {
		dec, err := NewDecoder(&source.Source{
			Metadata: source.Metadata{
				Identifier: "invalid-identifier",
			},
		})
		require.ErrorIs(t, err, ErrInvalidDataset)
		assert.Nil(t, dec)
	})
	t.Run("should return an error with a nil source", func(t *testing.T) {
		dec, err := NewDecoder(nil)

		require.ErrorIs(t, err, ErrNilSource)
		assert.Nil(t, dec)
	})
	t.Run("should return an error with a nil source data", func(t *testing.T) {
		dec, err := NewDecoder(&source.Source{
			Metadata: source.Metadata{
				Identifier: string(GrandLotoUUID),
			},
			Data: nil,
		})
		require.ErrorIs(t, err, ErrInvalidSourceReader)
		assert.Nil(t, dec)
	})
	t.Run("should create a decoder", func(t *testing.T) {
		zipReader := mocks.NewMockReader(t)

		dec, err := NewDecoder(&source.Source{
			Metadata: source.Metadata{
				Identifier: string(GrandLotoUUID),
			},
			Data: zipReader,
		})
		require.NoError(t, err)
		assert.NotNil(t, dec)
	})
}

func TestNewDecoderWithDataset(t *testing.T) {
	t.Run("should return an error with a nil source", func(t *testing.T) {
		dec, err := NewDecoderWithDataset(nil, DatasetInfo{})

		require.ErrorIs(t, err, ErrNilSource)
		assert.Nil(t, dec)
	})
	t.Run("should return an error with a nil source data", func(t *testing.T) {
		dec, err := NewDecoderWithDataset(&source.Source{
			Data: nil,
		}, DatasetInfo{})
		require.ErrorIs(t, err, ErrInvalidSourceReader)
		assert.Nil(t, dec)
	})
	t.Run("should create a decoder", func(t *testing.T) {
		zipReader := mocks.NewMockReader(t)

		dec, err := NewDecoderWithDataset(&source.Source{
			Data: zipReader,
		}, DatasetInfo{})
		require.NoError(t, err)
		assert.NotNil(t, dec)
	})
}

func TestDecoder_Decode(t *testing.T) {
	t.Run("when decoder fails to read data", func(t *testing.T) {
		zipReader := mocks.NewMockReader(t)

		zipReader.EXPECT().
			NumFile().
			Return(1).
			Times(assertMockTimeExpectation)

		zipReader.EXPECT().
			Read(0).
			Return(nil, assert.AnError).
			Times(assertMockTimeExpectation)

		dec := &decoder{
			src: &source.Source{Data: zipReader},
		}

		draws, err := dec.Decode()
		require.ErrorIs(t, err, ErrFailedToReadFile)
		assert.Nil(t, draws)
	})
	t.Run("when csv reader fails to be create", func(t *testing.T) {
		zipReader := mocks.NewMockReader(t)

		zipReader.EXPECT().
			NumFile().
			Return(1).
			Times(assertMockTimeExpectation)

		zipReader.EXPECT().
			Read(0).
			Return(nil, nil).
			Times(assertMockTimeExpectation)

		dec := &decoder{
			src: &source.Source{Data: zipReader},
		}

		draws, err := dec.Decode()
		require.ErrorIs(t, err, ErrFailedToCreateCSVReader)
		assert.Nil(t, draws)
	})
	t.Run("when csv reader fails to decode with decoder", func(t *testing.T) {
		zipReader := mocks.NewMockReader(t)

		zipReader.EXPECT().
			NumFile().
			Return(1).
			Times(assertMockTimeExpectation)

		zipReader.EXPECT().
			Read(0).
			Return(testdata.InvalidCSV, nil).
			Times(assertMockTimeExpectation)

		dec := &decoder{
			src: &source.Source{
				Metadata: source.Metadata{
					Identifier: string(Loto197605UUID),
				},
				Data: zipReader,
			},
			dataset: DatasetInfo{
				Type:    NewLotto,
				Version: LotteryV0,
			},
		}

		draws, err := dec.Decode()
		require.ErrorIs(t, err, ErrFailedToDecodeCSVData)
		assert.Nil(t, draws)
	})
	t.Run("when csv reader try to decode with an unknown version", func(t *testing.T) {
		zipReader := mocks.NewMockReader(t)

		zipReader.EXPECT().
			NumFile().
			Return(1).
			Times(assertMockTimeExpectation)

		zipReader.EXPECT().
			Read(0).
			Return(testdata.ValidCSVV0, nil).
			Times(assertMockTimeExpectation)

		dec := &decoder{
			src: &source.Source{
				Metadata: source.Metadata{
					Identifier: string(Loto197605UUID),
				},
				Data: zipReader,
			},
			dataset: DatasetInfo{
				Type:    NewLotto,
				Version: testdata.NotValidValue,
			},
		}

		draws, err := dec.Decode()
		require.ErrorIs(t, err, ErrInvalidDatasetVersion)
		assert.Nil(t, draws)
	})
	t.Run("when csv reader return a warning", func(t *testing.T) {
		zipReader := mocks.NewMockReader(t)

		zipReader.EXPECT().
			NumFile().
			Return(1).
			Times(assertMockTimeExpectation)

		zipReader.EXPECT().
			Read(0).
			Return(testdata.ValidCSVV0WithWarning, nil).
			Times(assertMockTimeExpectation)

		dec := &decoder{
			src: &source.Source{
				Metadata: source.Metadata{
					Identifier: string(Loto197605UUID),
				},
				Data: zipReader,
			},
			dataset: DatasetInfo{
				Type:    NewLotto,
				Version: LotteryV0,
			},
			unusedFields: make(map[string][]string),
		}

		expectedWarn := map[string][]string{
			"0-truc": {"42", "21"},
		}

		draws, err := dec.Decode()
		require.NoError(t, err)
		assert.Len(t, draws, 2)
		assert.Equal(t, expectedWarn, dec.unusedFields)
	})
	t.Run("when csv reader decode data", func(t *testing.T) {
		testCases := map[string]struct {
			version LotteryVersion
			data    []byte
			NumDraw int
		}{
			"with version 0": {
				version: LotteryV0,
				data:    testdata.ValidCSVV0,
				NumDraw: 2,
			},
			"with version 1 (no data)": {
				version: LotteryV1,
				data:    testdata.ValidCSVHeader,
				NumDraw: 0,
			},
			"with version 2 (no data)": {
				version: LotteryV2,
				data:    testdata.ValidCSVHeader,
				NumDraw: 0,
			},
			"with version 3 (no data)": {
				version: LotteryV3,
				data:    testdata.ValidCSVHeader,
				NumDraw: 0,
			},
			"with version 4 (no data)": {
				version: LotteryV4,
				data:    testdata.ValidCSVHeader,
				NumDraw: 0,
			},
		}

		for name, tt := range testCases {
			t.Run(name, func(t *testing.T) {
				zipReader := mocks.NewMockReader(t)

				zipReader.EXPECT().
					NumFile().
					Return(1).
					Times(assertMockTimeExpectation)

				zipReader.EXPECT().
					Read(0).
					Return(tt.data, nil).
					Times(assertMockTimeExpectation)

				dec := &decoder{
					src: &source.Source{
						Metadata: source.Metadata{
							Identifier: string(Loto197605UUID),
						},
						Data: zipReader,
					},
					dataset: DatasetInfo{
						Type:    NewLotto,
						Version: tt.version,
					},
				}

				draws, err := dec.Decode()
				require.NoError(t, err)
				assert.Len(t, draws, tt.NumDraw)
				assert.Empty(t, dec.unusedFields)
			})
		}
	})
}
