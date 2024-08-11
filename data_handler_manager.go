package file_tranport

type dataHandlerManager struct {
	dataHandlers map[string]*dataHandler
}

func newDataHandlerManager() *dataHandlerManager {
	return &dataHandlerManager{
		dataHandlers: make(map[string]*dataHandler),
	}
}

func (m *dataHandlerManager) RegisterDataHandler(key string, handler *dataHandler) {
	m.dataHandlers[key] = handler
}

func (m *dataHandlerManager) GetDataHandler(key string) (*dataHandler, bool) {
	handler, ok := m.dataHandlers[key]
	return handler, ok
}
