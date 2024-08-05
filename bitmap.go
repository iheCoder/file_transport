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
	var i int
	for i = 0; i < len(b.data)-1; i++ {
		if b.data[i] != 0xff {
			return false
		}
	}

	for j := 0; j < b.size%8; j++ {
		if b.data[i]&(1<<uint(j)) == 0 {
			return false
		}
	}
	return true
}
