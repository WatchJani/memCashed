package constants

const (
	HeaderSize      = 10
	SetOperation    = 'S'
	GetOperation    = 'G'
	DeleteOperation = 'D'
)

var (
	ObjectInserted = []byte("object inserted")
	ObjectDeleted  = []byte("Deleted")

	ErrObjectNotFound = []byte("object not found")
	ErrTimeExpire     = []byte("time expire")
)
