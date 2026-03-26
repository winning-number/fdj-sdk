package lotto

import (
	"github.com/gofast-pkg/csv"
	"github.com/winning-number/fdj-sdk-lotto/draw"
)

const keyDrawType = "drawType"

func newInstanceFunc(version draw.Version) draw.Converter {
	switch version {
	case draw.V0:
		return &draw.CSV0{}
	case draw.V1:
		return &draw.CSV1{}
	case draw.V2:
		return &draw.CSV2{}
	case draw.V3:
		return &draw.CSV3{}
	case draw.V4:
		return &draw.CSV4{}
	default:
		return &draw.CSV4{}
	}
}

func saveInstanceFunc(instance any, decoder csv.Decoder) (draw.Draw, error) {
	var drawConv draw.Converter
	var err error
	var d draw.Draw
	var ok bool
	var drawT string

	if drawT, ok = decoder.ContextGet(keyDrawType); !ok {
		return draw.Draw{}, ErrDrawTypeKeyNotFound
	}
	if drawConv, ok = instance.(draw.Converter); !ok {
		return draw.Draw{}, ErrInvalidDrawType
	}
	if d, err = draw.Convert(drawConv, (draw.Type)(drawT)); err != nil {
		return draw.Draw{}, err
	}

	return d, nil
}
