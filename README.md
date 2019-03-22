# uidgen

    go get -u github.com/codemodus/uidgen

## Usage

```go
type UID
    func (u UID) Bytes() []byte
    func (u *UID) Scan(src interface{}) error
    func (u UID) String() string
    func (u *UID) Value() (driver.Value, error)
type UIDGen
    func New(offset uint64, vt ValueType) *UIDGen
    func (g *UIDGen) Parse(s string) (UID, bool)
    func (g *UIDGen) UID() UID
type ValueType
```

```go
type ValueType uint8

const (
    BINARY16 ValueType = iota
    VARCHAR26
)
```
