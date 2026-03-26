package draw

import "github.com/pkg/errors"

// error list.
var (
	ErrCSVDay     = errors.New("day for dayConverter is invalid")
	ErrCSVDate    = errors.New("date for dateConverter is invalid")
	ErrCSVType    = errors.New("type for typeConverter is invalid")
	ErrCSVPrice   = errors.New("price for priceConverter is invalid")
	ErrEmptyMoney = errors.New("money is empty to parse")
)
