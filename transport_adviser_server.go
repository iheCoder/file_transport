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
	if float64(info.Size)/t.bandwidth > float64(maxTransportSecondsPerBlock) {
		return NewServerDataHandler(info.Path, info.Count, info.BlockSize, fileMode)
	} else {
		return NewServerDataHandler(info.Path, info.Count, info.BlockSize, memoryMode)
	}
}
