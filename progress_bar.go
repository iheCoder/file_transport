package file_tranport

import "sync"

type progressBar struct {
	// total is the total number of tasks.
	total int
	// completedIndex is the index of the completed boundary.
	completedIndex int
	// bm is the bitmap to record the completion status of each task.
	bm *bitmap
	// rwlock is the lock to protect the completedIndex and bm.
	rwlock sync.RWMutex
}

func NewProgressBar(total int) *progressBar {
	return &progressBar{
		total:  total,
		bm:     NewBitmap(total),
		rwlock: sync.RWMutex{},
		// The initial value of completedIndex is -1, which means no task has been completed.
		completedIndex: -1,
	}
}

func (p *progressBar) Set(index int) {
	p.rwlock.Lock()
	defer p.rwlock.Unlock()

	p.bm.Set(index)

	// If the index is the next one to be completed, update the completedIndex.
	if index == p.completedIndex+1 {
		p.completedIndex = index
	}
}

func (p *progressBar) IsAllSet() bool {
	p.rwlock.RLock()
	defer p.rwlock.RUnlock()

	return p.bm.IsAllSet()
}

func (p *progressBar) FindNextUnset() int {
	p.rwlock.RLock()
	defer p.rwlock.RUnlock()

	for i := p.completedIndex + 1; i < p.total; i++ {
		if !p.bm.IsSet(i) {
			return i
		}
	}

	return -1
}
