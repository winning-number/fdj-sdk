package draw

// CSV1 is the draw type v1.
// It is really near from the draw type v0.
// The only difference is the addition of the joker and the tirage order.
// So, CSV0 is embedded.
type CSV1 struct {
	CSV0

	Joker  string `csv:"numero_joker"`
	Tirage int32  `csv:"1er_ou_2eme_tirage"`
}

func (c CSV1) joker() Joker {
	return Joker{
		Plus:    c.JokerPlus,
		Base:    c.Joker,
		HasBase: true,
	}
}

func (c CSV1) metadata(drawT Type) (Metadata, error) {
	var meta Metadata
	var err error

	if meta, err = c.coreCSV.metadata(drawT); err != nil {
		return Metadata{}, err
	}
	meta.OldType = true
	meta.Version = V1
	meta.TirageOrder = c.Tirage
	meta.setID(c.WinOrder)

	return meta, nil
}
