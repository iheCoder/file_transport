package file_tranport

type bitmap struct {
	data []byte
	size int
}

func NewBitmap(size int) *bitmap {
	return &bitmap{
		data: make([]byte, (size+7)/8),
		size: size,
	}
}

func (b *bitmap) Set(index int) {
	b.data[index/8] |= 1 << uint(index%8)
}

func (b *bitmap) Unset(index int) {
	b.data[index/8] &= ^(1 << uint(index%8))
}

func (b *bitmap) IsSet(index int) bool {
	return b.data[index/8]&(1<<uint(index%8)) != 0
}

func (b *bitmap) IsAllSet() bool {
	for _, v := range b.data {
		if v != 0xff {
			return false
		}
	}
	return true
}
