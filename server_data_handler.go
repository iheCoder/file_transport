package file_tranport

func newServerDataHandler(path string) (*dataHandler, error) {
	fh, err := NewFileWriterHelper(path)
	if err != nil {
		return nil, err
	}

	return &dataHandler{
		fh: fh,
	}, nil
}

func (d *dataHandler) initServerDataHandler(bd *blockData) {
	d.bds = make([]blockData, bd.Count)
	d.pb = NewProgressBar(bd.Count)
	d.blockSize = bd.BlockSize
	d.fh.InitWriterFileHelper(bd.BlockSize)
}

func (d *dataHandler) WriteBlock(bd *blockData) error {
	if !d.initialized {
		d.initServerDataHandler(bd)
		d.initialized = true
	}

	return d.writeBlock(bd)
}

func (d *dataHandler) writeBlock(bd *blockData) error {
	switch d.mode {
	case fileMode:
		return d.fh.WriteBlock(bd.Index, bd.Raw)
	default:
		d.bds[bd.Index] = *bd
		d.pb.Set(bd.Index)
	}

	return nil
}

func (d *dataHandler) combinedOne() []byte {
	total := make([]byte, 0)
	for _, bd := range d.bds {
		total = append(total, bd.Raw...)
	}
	return total
}

func (d *dataHandler) Persist() error {
	// if file mode, do nothing
	if d.mode == fileMode {
		return nil
	}

	return d.fh.WriteAll(d.combinedOne())
}

func (d *dataHandler) IsAllBlockCompleted() bool {
	return d.pb.IsAllSet()
}
