package jsutils

// BytesPool  is a pool of byte slice that can ben used
type BytesPool struct {
	pool chan []byte
}

func NewBytesPool(max int) *BytesPool {
	return &BytesPool{
		pool: make(chan []byte, max),
	}
}

// Get returns a byte slice with at least sz capacity.
func (p *BytesPool) GetFifo(sz int) []byte {
	var c []byte
	select {
	case c = <-p.pool:
	default:
		return make([]byte, sz)
	}

	if cap(c) < sz {
		return make([]byte, sz)
	}

	return c[:sz]
}

// Put returns a slice back to the pool
func (p *BytesPool) PutFifo(c []byte) {
	select {
	case p.pool <- c:
	default:
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////

type EntryFifo struct {
	Key    int32
	Values []byte
}

func NewEntryFifo(key int32, values []byte) EntryFifo {
	bbb := make([]byte,len(values))
	copy(bbb,values)
	return EntryFifo{
		Key:    key,
		Values: bbb,
	}
}

type LimitedEntryFifo struct {
	pool chan EntryFifo
}

func NewLimitedEntryFifo(capacity int) *LimitedEntryFifo {
	return &LimitedEntryFifo{
		pool: make(chan EntryFifo, capacity),
	}
}

func (p *LimitedEntryFifo) GetEntryFifo() (EntryFifo, bool) {
	var e EntryFifo

	select {
	case e = <-p.pool:
		return e, true
	default:
	}
	return e, false
}

func (p *LimitedEntryFifo) PutEntryFifo(e EntryFifo) {
	select {
	case p.pool <- e:
	default:
	}
}

func (p *LimitedEntryFifo) EntryFifoLen() int {
	return len(p.pool)
}
