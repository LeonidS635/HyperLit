package entry

type Entry struct {
	Type byte
	Size int
	Name string
	Hash []byte
	Data []byte
}

func (e Entry) GetType() byte {
	return e.Type
}

func (e Entry) GetName() string {
	return e.Name
}

func (e Entry) GetHash() []byte {
	return e.Hash
}

func (e Entry) GetData() []byte {
	return e.Data
}

func (e Entry) GetContent() []byte {
	return e.Data
}
