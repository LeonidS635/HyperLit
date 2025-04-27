package entry

type Entry struct {
	Path string
	Name string
	Type byte
	Size int
	Hash []byte
	Data []byte
}

func (e Entry) GetType() byte {
	return e.Type
}

func (e Entry) GetName() string {
	return e.Name
}

func (e Entry) GetPath() string {
	return e.Path
}

func (e Entry) GetHash() []byte {
	return e.Hash
}

func (e Entry) GetData() []byte {
	return e.Data
}
