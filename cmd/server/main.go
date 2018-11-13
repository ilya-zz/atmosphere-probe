package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/logrange/logrange/pkg/proto/atmosphere"
	"github.com/sirupsen/logrus"
)

type listener struct {
	storage io.Writer
}

func newListener() atmosphere.ServerListener {
	_, err := os.OpenFile("store.db", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatal(err)
	}
	return &listener{
		//storage: bufio.NewWriterSize(f, 4096),
		storage: ioutil.Discard,
	}
}

func (l *listener) OnRead(r atmosphere.Reader, n int) error {
	if _, err := io.CopyN(l.storage, r, int64(n)); err != nil {
		panic(err)
	}
	r.ReadResponse(nil)
	return nil
}

var (
	port = flag.Int("p", 8080, "server port")
)

func main() {
	flag.Parse()

	addr := fmt.Sprintf("0.0.0.0:%d", *port)
	log.Println("Listen on ", addr)

	s, err := atmosphere.NewServer(&atmosphere.ServerConfig{
		ListenAddress: addr,
		SessTimeoutMs: 60 * 10000,
		ConnListener:  newListener(),
		Auth:          func(a, s string) bool { return true },
	})
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()
	for {
		time.Sleep(time.Second)
	}
}
