package offthegrid

import (
	"gopkg.in/mgo.v2"
	"net"
	"net/http"
	"time"
)

const VERSION = 1

type OTGServer struct {
	Config   Config
	dbsess   *mgo.Session
	listener net.Listener
}

func (h *OTGServer) DB() *mgo.Session {
	return h.dbsess.Clone()
}

func (h *OTGServer) Run() (cstarted chan bool, cerr chan error) {
	cstarted = make(chan bool)
	cerr = make(chan error)

	go func(cerr chan error) {
		err := http.Serve(h.listener, h)
		if err != nil {
			cerr <- err
		}
	}(cerr)

	go func(cstarted chan bool, cerr chan error) {
		// wait 100ms before checking if we started successfully
		select {
		case _ = <-time.After(100 * time.Millisecond):
		}

		select {
		case err := <-cerr:
			cerr <- err
		default:
			cstarted <- true
		}
	}(cstarted, cerr)

	return
}

func NewServer(cfg Config) (*OTGServer, error) {
	var err error

	// Open a UNIX socket at the configured path
	listener, err := net.Listen("unix", cfg.ListenSocket)
	if err != nil {
		return nil, err
	}

	// Connect to MongoDB
	session, err := mgo.Dial(cfg.ConnectionURI)
	if err != nil {
		return nil, err
	}

	// And create the server object.
	server := &OTGServer{
		Config:   cfg,
		dbsess:   session,
		listener: listener,
	}

	return server, nil
}
