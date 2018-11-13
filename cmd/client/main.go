package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/ilya-zz/logrange/pkg/proto/atmosphere"
	"github.com/sirupsen/logrus"
)

var (
	url   = flag.String("url", "", "log server url")
	bs    = flag.Int("bs", 1024, "block size")
	count = flag.Int("count", 10000, "block count")
	hbeat = 10000
)

func main() {
	flag.Parse()

	if *url == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	cl, err := atmosphere.NewClient(*url, &atmosphere.ClientConfig{
		HeartBeatMs: hbeat,
	})
	if err != nil {
		logrus.Fatal(err)
	}

	log.Printf("send %d packets (%s each) to %s\n", *count, h(*bs), *url)

	msg := []byte(strings.Repeat("X", *bs))

	var sent, recv int

	t0 := time.Now()
	for i := 0; i < *count; i++ {
		n, err := cl.Write(atmosphere.Message(msg), atmosphere.Message([]byte{}))
		if err != nil {
			logrus.Fatal(err)
		}
		sent += len(msg)
		recv += n
	}
	fmt.Printf("Sent/recv %s / %s bytes in %f secs\n\n", h(sent), h(recv), time.Since(t0).Seconds())
}

func h(v int) string {
	return humanize.Bytes(uint64(v))
}
