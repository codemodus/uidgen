package uidgen

import (
	"database/sql/driver"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid"
)

// ValueType ...
type ValueType uint8

// ValueType constants ...
const (
	BINARY16 ValueType = iota
	VARCHAR26
)

// UID ...
type UID struct {
	id  ulid.ULID
	bin bool
}

// Bytes ...
func (u UID) Bytes() []byte {
	return u.id[:]
}

func (u UID) String() string {
	return u.id.String()
}

// Scan ...
func (u *UID) Scan(src interface{}) error {
	return u.id.Scan(src)
}

// Value ...
func (u *UID) Value() (driver.Value, error) {
	if u.bin {
		return u.id.Value()
	}
	return u.id.String(), nil
}

// UIDGen ...
type UIDGen struct {
	ofs uint64
	erp *entropyReaderPool
	vt  ValueType
}

//func (id stringValuer) Value() (driver.Value, error) {
//		return ULID(id).String(), nil
//	}

// New ...
func New(offset uint64, vt ValueType) *UIDGen {
	g := UIDGen{
		ofs: offset,
		erp: newEntropyReaderPool(),
		vt:  vt,
	}

	return &g
}

// UID ...
func (g *UIDGen) UID() UID {
	ms := ulid.Timestamp(time.Now()) - g.ofs
	r := g.erp.get()

	lid, err := ulid.New(ms, r)
	g.erp.put(r)
	if err != nil {
		panic(fmt.Sprintf("uidgen: %s", err))
	}

	return UID{id: lid, bin: g.vt == BINARY16}
}

// Parse ...
func (g *UIDGen) Parse(s string) (UID, bool) {
	lid, err := ulid.Parse(s)
	if err != nil {
		return UID{}, false
	}
	return UID{id: lid, bin: g.vt == BINARY16}, true
}

type entropyReaderPool struct {
	p sync.Pool
}

func newEntropyReaderPool() *entropyReaderPool {
	return &entropyReaderPool{
		p: sync.Pool{
			New: func() interface{} {
				t := time.Now()
				rnr := rand.New(rand.NewSource(t.UnixNano()))
				return ulid.Monotonic(rnr, 0)
			},
		},
	}
}

func (p *entropyReaderPool) get() io.Reader {
	return p.p.Get().(io.Reader)
}

func (p *entropyReaderPool) put(r io.Reader) {
	p.p.Put(r)
}
