package entry

type Interface interface {
	GetType() byte
	GetName() string
	GetPath() string
	GetHash() []byte
	GetData() []byte
}
