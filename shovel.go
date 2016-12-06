package gosshtun

import (
	"io"
	"os"
)

// Shovel shovels data from an io.ReadCloser to an io.WriteCloser
// in an independent go routine started by Shovel::Start().
// You can request that the shovel stop by closing ReqStop,
// and wait until Done is closed to know that it is finished.
type shovel struct {
	Done    chan bool
	ReqStop chan bool
	Ready   chan bool

	// logging functionality, off by default
	DoLog     bool
	LogReads  io.Writer
	LogWrites io.Writer
}

// make a new Shovel
func newShovel(doLog bool) *shovel {
	return &shovel{
		Done:      make(chan bool),
		ReqStop:   make(chan bool),
		Ready:     make(chan bool),
		DoLog:     doLog,
		LogReads:  os.Stdout,
		LogWrites: os.Stdout,
	}
}

type readerNilCloser struct{ io.Reader }

func (rc *readerNilCloser) Close() error { return nil }

type writerNilCloser struct{ io.Writer }

func (wc *writerNilCloser) Close() error { return nil }

// Start starts the shovel doing an io.Copy from r to w. The
// goroutine that is running the copy will close the Ready
// channel just before starting the io.Copy. The
// label parameter allows reporting on when a specific shovel
// was shut down.
func (s *shovel) Start(w io.WriteCloser, r io.ReadCloser, label string) {

	if s.DoLog {
		// TeeReader returns a Reader that writes to w what it reads from r.
		// All reads from r performed through it are matched with
		// corresponding writes to w. There is no internal buffering -
		// the write must complete before the read completes.
		// Any error encountered while writing is reported as a read error.
		r = &readerNilCloser{io.TeeReader(r, s.LogReads)}
		w = &writerNilCloser{io.MultiWriter(w, s.LogWrites)}
	}

	go func() {
		var err error
		var n int64
		defer func() {
			close(s.Done)
			p("shovel %s copied %d bytes before shutting down", label, n)
		}()
		close(s.Ready)
		n, err = io.Copy(w, r)
		if err != nil {
			// don't freak out, the network connection got closed most likely.
			// e.g. read tcp 127.0.0.1:33631: use of closed network connection
			//panic(fmt.Sprintf("in Shovel '%s', io.Copy failed: %v\n", label, err))
			return
		}
	}()
	go func() {
		<-s.ReqStop
		r.Close() // causes io.Copy to finish
		w.Close()
	}()
}

// stop the shovel goroutine. returns only once the goroutine is done.
func (s *shovel) Stop() {
	// avoid double closing ReqStop here
	select {
	case <-s.ReqStop:
	default:
		close(s.ReqStop)
	}
	<-s.Done
}

// a shovelPair manages the forwarding of a bidirectional
// channel, such as that in forwarding an ssh connection.
type shovelPair struct {
	AB      *shovel
	BA      *shovel
	Done    chan bool
	ReqStop chan bool
	Ready   chan bool

	DoLog bool
}

// make a new shovelPair
func newShovelPair(doLog bool) *shovelPair {
	return &shovelPair{
		AB:      newShovel(doLog),
		BA:      newShovel(doLog),
		Done:    make(chan bool),
		ReqStop: make(chan bool),
		Ready:   make(chan bool),
	}
}

// Start the pair of shovels. abLabel will label the a<-b shovel. baLabel will
// label the b<-a shovel.
func (s *shovelPair) Start(a io.ReadWriteCloser, b io.ReadWriteCloser, abLabel string, baLabel string) {
	s.AB.Start(a, b, abLabel)
	<-s.AB.Ready
	s.BA.Start(b, a, baLabel)
	<-s.BA.Ready
	close(s.Ready)

	// if one stops, shut down the other
	go func() {
		select {
		case <-s.AB.Done:
			s.BA.Stop()
		case <-s.BA.Done:
			s.AB.Stop()
		}
	}()
}

func (s *shovelPair) Stop() {
	s.AB.Stop()
	s.BA.Stop()
}
