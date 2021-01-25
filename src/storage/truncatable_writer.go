package storage

// A TruncatableWriter is a buffer that supports flexible operations.
//
// The behaviour of all functions is that of os.File (os.file fulfills this interface)
type TruncatableWriter interface {
	Truncate(n int64) error
	Write(b []byte) (n int, err error)
	Read(p []byte) (n int, err error)
	Seek(offset int64, whence int) (ret int64, err error)
	Sync() (err error)
}
