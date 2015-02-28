package offthegrid

type Config struct {
	// ListenSocket is the path to the UNIX socket to listen on.
	ListenSocket string `yaml:"socket"`

	// ConnectionURI is the MongoDB connection string to use to connect to the
	// server.
	//
	// Relevant documentation:
	//
	//    http://docs.mongodb.org/manual/reference/connection-string/
	ConnectionURI string `yaml:"database_uri"`

	// GridFSPrefix
	GridFSPrefix string `yaml:"gridfs_prefix"`

	// CORSHeader is the CORS header string to send along with the HTTP response.
	CORSHeader string `yaml:"cors_header"`

	// MaxAge is the number of seconds to set the MaxAge header to. Also controls Expiry.
	MaxAge int `yaml:"max_age"`
}
