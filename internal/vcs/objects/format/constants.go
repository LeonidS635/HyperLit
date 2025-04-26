package format

// Objects type constants (should be 1 byte)
const (
	CodeType = iota
	DocsType
	TreeType
)

// Separator constants (should be 1 byte)
const (
	TreeEntriesSeparator = byte('\n')
)

// Size constants
const (
	typeBytesN      = 1
	separatorBytesN = 1
	sizeBytesN      = 4

	HeaderSize = typeBytesN + sizeBytesN
)

//func init() {
//	sizeChecker := func(size int, b ...byte) {
//		if len(b) == size {
//			log.Fatalln(fmt.Errorf("parametr %q had invalid size %d", b, size))
//		}
//	}
//
//	sizeChecker(typeBytesN, CodeType)
//	sizeChecker(typeBytesN, DocsType)
//	sizeChecker(typeBytesN, TreeType)
//
//	sizeChecker(separatorBytesN, typeSizeSeparator)
//	sizeChecker(separatorBytesN, headerDataSeparator)
//	sizeChecker(separatorBytesN, TreeFieldsSeparator)
//	sizeChecker(separatorBytesN, TreeEntriesSeparator)
//
//	// TODO: think about controlling size bytes number (4 or 8)
//}
