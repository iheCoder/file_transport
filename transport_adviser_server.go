package file_tranport

func newServerTransportAdviser(bandwidth float64, blockTransportMS int) *transportAdviser {
	return &transportAdviser{
		bandwidth:        bandwidth,
		blockTransportMS: blockTransportMS,
		mgr:              newDataHandlerManager(),
	}
}

func (t *transportAdviser) GetServerDataHandler(info *FileInfo) (*dataHandler, error) {
	// update total size
	t.updateTotalSize(info.Size)

	// try get data handler from manager
	handler, ok := t.mgr.GetDataHandler(info.Key)
	if ok {
		return handler, nil
	}

	// create new server data handler
	mode := memoryMode
	if float64(info.Size)/t.bandwidth > float64(maxTransportSecondsPerBlock) {
		mode = fileMode
	}

	// create new server data handler
	dh, err := NewServerDataHandler(info.Path, info.Count, info.BlockSize, mode)
	if err != nil {
		return nil, err
	}

	// add data handler to manager
	t.mgr.RegisterDataHandler(info.Key, dh)

	return dh, nil
}
