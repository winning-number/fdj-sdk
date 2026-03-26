// Package draw contains the main types and functions to manipulate lotto's draws.
package draw

import (
	"sort"
)

// Order type.
const (
	OrderASC OrderType = iota
	OrderDESC
	OrderNone
)

// OrderType is a type to represent the order of a list of draws.
type OrderType uint8

// Draw represents a lotto draw.
type Draw struct {
	Metadata Metadata
	Roll     Roll
	Joker    Joker
	WinStats WinStats
	WinCode  WinCode
}

// OrderDraws order draws by date.
// If OrderType is OrderASC, order from less recent to more recent.
// If OrderType is OrderDESC, order from more recent to less recent.
// If OrderType is OrderNone, do nothing.
func OrderDraws(draws *[]Draw, order OrderType) {
	if order == OrderNone {
		return
	}
	sort.SliceStable(*draws, func(i, j int) bool {
		if order == OrderASC {
			return updateDrawOrder(draws, i, j)
		}

		// invert the index to match the updateDrawOrder comparison
		// with the DESC order.
		return updateDrawOrder(draws, j, i)
	})
}

func updateDrawOrder(draws *[]Draw, i, j int) bool {
	if (*draws)[i].Metadata.Date.After((*draws)[j].Metadata.Date) {
		return true
	}
	if (*draws)[i].Metadata.Date.Equal((*draws)[j].Metadata.Date) &&
		(*draws)[i].Metadata.TirageOrder > (*draws)[j].Metadata.TirageOrder {
		return true
	}

	return false
}

// Finder find a draw in a list of draws.
func Finder(list *[]Draw, draw Draw) bool {
	for _, d := range *list {
		if d.Metadata.ID == draw.Metadata.ID {
			return true
		}
	}

	return false
}
