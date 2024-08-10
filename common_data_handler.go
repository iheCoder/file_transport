package file_tranport

type DataParser func(data []byte) (any, error)

var globalCommonDataHandler = &CommonDataHandler{
	dataParsers: make(map[TransportDataType]DataParser),
}

type CommonDataHandler struct {
	dataParsers map[TransportDataType]DataParser
}

func RegisterDataParser(dataType TransportDataType, parser DataParser) {
	globalCommonDataHandler.dataParsers[dataType] = parser
}
