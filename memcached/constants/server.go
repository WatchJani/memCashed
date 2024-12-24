package constants

import (
	"errors"
	"runtime"
)

const (
	SetOperation    = 'S'
	GetOperation    = 'G'
	DeleteOperation = 'D'

	HeaderSize = 10
	MiB        = 1024 * 1024
	TCP        = "tcp"

	KiB                       = 1024 // 1 MiB in bytes
	MinimumNumberOfConnection = 5    // Minimum number of connections to the server
	IntDefaultValue           = 0    // Default value for integers
	DefaultPort               = 5000 // Default server port
	BufferSizeTCP             = 4
)

var (
	DefaultNumberOfWorkers = runtime.NumCPU() // Default number of worker threads

	ObjectInserted = []byte("object inserted")
	ObjectDeleted  = []byte("deleted")

	// ErrOperationIsNotSupported is the error returned when an unsupported operation is attempted.
	ErrOperationIsNotSupported = errors.New("operation is not supported")

	// ErrNotEnoughSpace is the error returned when there is not enough space to allocate memory.
	ErrNotEnoughSpace = errors.New("there is not enough space")

	// NoReq is a buffer used for reading unprocessed data when memory allocation fails.
	NoReq             []byte = make([]byte, MiB+HeaderSize) // Buffer to read unprocessed data
	ErrObjectNotFound        = []byte("object not found")
	ErrTimeExpire            = []byte("time expire")

	InfoServerClose     = "server closed"
	InfoConnectionClose = "connection is close"
)
