package entry

type Interface interface {
	GetType() byte
	GetName() string
	GetHash() []byte
	GetData() []byte
}
