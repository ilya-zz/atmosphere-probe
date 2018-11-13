package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/ilya-zz/logrange/pkg/proto/atmosphere"
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

func main() {
	s, err := atmosphere.NewServer(&atmosphere.ServerConfig{
		ListenAddress: "0.0.0.0:10101",
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
