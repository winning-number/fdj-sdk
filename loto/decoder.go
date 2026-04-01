package loto

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/gofast-pkg/csv"
	xcsv "github.com/winning-number/fdj-sdk/v2/loto/csv"
	"github.com/winning-number/fdj-sdk/v2/source"
)

const (
	contextKeyDataset = "contextKeyDataset"
	contextKeyData    = "contextKeyData"
	contextKeyWarning = "contextKeyWarning"
)

// Public error for decoder operations.
var (
	ErrNilSource                = errors.New("source is nil")
	ErrInvalidSourceReader      = errors.New("invalid source reader")
	ErrInvalidDataset           = errors.New("invalid dataset")
	ErrFailedToReadFile         = errors.New("failed to read file from source")
	ErrFailedToParseCSV         = errors.New("failed to parse csv file")
	ErrFailedToCreateCSVReader  = errors.New("failed to create csv reader")
	ErrFailedToCreateCSVDecoder = errors.New("failed to create csv decoder")
	ErrFailedToDecodeCSVData    = errors.New("failed to decode csv data")

	ErrDataMissingInCSVContext = errors.New("data missing in the csv context")
	ErrDataTypeInCSVContext    = errors.New("data with an invalid type in csv context")
	ErrInvalidDatasetVersion   = errors.New("invalid dataset version to create new instance")

	ErrMissingWarningInCSVContext = errors.New("warning missing in the csv context")
	ErrInvalidWarningInCSVContext = errors.New("invalid warning in csv context")
)

// Decoder is a factory to decode a csv file.
// Create a new Decoder with the NewDecoder function or NewDecoderWithDataset if you want target custom csv file.
type Decoder interface {
	// Decode the CSV file configured during the call to NewDecoder / NewDecoderWithDataset.
	Decode() ([]*LotteryDraw, error)
	// UnusedFields return slice of header which are not used by LotteryDraw.Metadata.Identifier.
	UnusedFields() map[string][]string
}

type decoder struct {
	src          *source.Source
	dataset      DatasetInfo
	unusedFields map[string][]string
}

// NewDecoder decode the source parameter.
// It automatically get the associated dataset info to know the type and version used to decode data.
// If you decode a custom csv file, prefer calling the NewDecoderWithDataset function to provide the
// appropriate DatasetInfo.
func NewDecoder(src *source.Source) (Decoder, error) {
	if src == nil {
		return nil, ErrNilSource
	}

	data, ok := Datasets[DatasetIdentifier(src.Identifier)]
	if !ok {
		return nil, ErrInvalidDataset
	}

	return NewDecoderWithDataset(src, data)
}

// NewDecoderWithDataset take the source parameter and try to decode it with the DatasetInfo.
func NewDecoderWithDataset(src *source.Source, dataset DatasetInfo) (Decoder, error) {
	if src == nil {
		return nil, ErrNilSource
	}

	if src.Data == nil {
		return nil, ErrInvalidSourceReader
	}

	return &decoder{
		src:          src,
		dataset:      dataset,
		unusedFields: make(map[string][]string, 0),
	}, nil
}

func (d *decoder) Decode() ([]*LotteryDraw, error) {
	var lotteryDraws []*LotteryDraw

	numFile := d.src.Data.NumFile()
	for i := 0; i < numFile; i++ {
		buf, err := d.src.Data.Read(i)
		if err != nil {
			return nil, errors.Join(err, ErrFailedToReadFile)
		}

		draws, err := d.parseCSV(bytes.NewBuffer(buf), i)
		if err != nil {
			return nil, errors.Join(err, ErrFailedToParseCSV)
		}

		lotteryDraws = append(lotteryDraws, draws...)
	}

	return lotteryDraws, nil
}

func (d *decoder) UnusedFields() map[string][]string {
	return d.unusedFields
}

func (d *decoder) parseCSV(input io.Reader, index int) ([]*LotteryDraw, error) {
	var err error
	var csvReader csv.CSV
	var csvDecoder csv.Decoder
	var lotteryDraws []*LotteryDraw

	if csvReader, err = csv.New(input, ';'); err != nil {
		return nil, errors.Join(err, ErrFailedToCreateCSVReader)
	}
	if csvDecoder, err = csv.NewDecoder(csv.ConfigDecoder{
		NewInstanceFunc:  newInstanceFunc,
		SaveInstanceFunc: saveInstanceFunc,
		WarningInstanceFunc: func(decoder csv.Decoder, warn csv.Warning) error {
			return warningInstanceFunc(decoder, warn, index)
		},
	}); err != nil {
		return nil, errors.Join(err, ErrFailedToCreateCSVDecoder)
	}

	csvDecoder.ContextSet(contextKeyDataset, d.dataset)
	csvDecoder.ContextSet(contextKeyData, &lotteryDraws)
	csvDecoder.ContextSet(contextKeyWarning, d.unusedFields)

	if err = csvReader.DecodeWithDecoder(csvDecoder); err != nil {
		return nil, errors.Join(err, ErrFailedToDecodeCSVData)
	}

	return lotteryDraws, nil
}

func newInstanceFunc(decoder csv.Decoder) (any, error) {
	info, _, err := parseDecoderContext(decoder)
	if err != nil {
		return nil, err
	}

	switch info.Version {
	case LotteryV0:
		return &xcsv.Version0{}, nil
	case LotteryV1:
		return &xcsv.Version1{}, nil
	case LotteryV2:
		return &xcsv.Version2{}, nil
	case LotteryV3:
		return &xcsv.Version3{}, nil
	case LotteryV4:
		return &xcsv.Version4{}, nil
	default:
		return nil, ErrInvalidDatasetVersion
	}
}

func saveInstanceFunc(decoder csv.Decoder, obj any) error {
	info, draws, err := parseDecoderContext(decoder)
	if err != nil {
		return err
	}

	draw, err := CSVConverter(info.Type, info.Version, obj)
	if err != nil {
		return err
	}

	*draws = append((*draws), draw)

	return nil
}

func warningInstanceFunc(decoder csv.Decoder, warn csv.Warning, rawIndex int) error {
	val, ok := decoder.ContextGet(contextKeyWarning)
	if !ok {
		return ErrMissingWarningInCSVContext
	}
	unusedFields, ok := val.(map[string][]string)
	if !ok {
		return ErrInvalidWarningInCSVContext
	}

	for k, v := range warn {
		key := fmt.Sprintf("%d-%s", rawIndex, k)
		unusedFields[key] = append(unusedFields[key], v)
	}

	return nil
}

func parseDecoderContext(decoder csv.Decoder) (DatasetInfo, *[]*LotteryDraw, error) {
	var dataset DatasetInfo
	var draws *[]*LotteryDraw

	val, ok := decoder.ContextGet(contextKeyDataset)
	if !ok {
		return DatasetInfo{}, nil, fmt.Errorf("%s: %w", contextKeyDataset, ErrDataMissingInCSVContext)
	}
	dataset, ok = val.(DatasetInfo)
	if !ok {
		return DatasetInfo{}, nil, fmt.Errorf("DatasetInfo: %w", ErrDataTypeInCSVContext)
	}

	val, ok = decoder.ContextGet(contextKeyData)
	if !ok {
		return DatasetInfo{}, nil, fmt.Errorf("%s: %w", contextKeyData, ErrDataMissingInCSVContext)
	}
	draws, ok = val.(*[]*LotteryDraw)
	if !ok {
		return DatasetInfo{}, nil, fmt.Errorf("*[]*LotteryDraw:%w", ErrDataTypeInCSVContext)
	}

	return dataset, draws, nil
}
