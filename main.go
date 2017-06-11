package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

// a lot of this borrows directly from:
// 	https://github.com/prometheus/prometheus/blob/master/documentation/examples/remote_storage/remote_storage_adapter/main.go

type config struct {
	//tcp://host1:9000?username=user&password=qwerty&database=clicks&read_timeout=10&write_timeout=20&alt_hosts=host2:9000,host3:9000
	ChDSN           string
	ChDB            string
	ChTable         string
	ChBatch         int
	ChanSize        int
	HTTPTimeout     time.Duration
	HTTPAddr        string
	HTTPWritePath   string
	HTTPMetricsPath string
}

var (
	versionFlag bool
)

func main() {
	excode := 0

	conf := parseFlags()

	if versionFlag {
		fmt.Println("Git Commit:", GitCommit)
		fmt.Println("Version:", Version)
		if VersionPrerelease != "" {
			fmt.Println("Version PreRelease:", VersionPrerelease)
		}
		os.Exit(excode)
	}

	fmt.Println("Starting up..")

	srv, err := NewP2CServer(conf)
	if err != nil {
		fmt.Printf("Error: could not create server: %s\n", err.Error())
		excode = 1
		os.Exit(excode)
	}
	err = srv.Start()
	if err != nil {
		fmt.Printf("Error: http server returned error: %s\n", err.Error())
		excode = 1
	}

	fmt.Println("Shutting down..")
	srv.Shutdown()
	fmt.Println("Exiting..")
	os.Exit(excode)
}

func parseFlags() *config {
	cfg := new(config)

	// print version?
	flag.BoolVar(&versionFlag, "version", false, "Version")

	// clickhouse dsn
	ddsn := "tcp://127.0.0.1:9000?username=&password=&database=metrics&" +
		"read_timeout=10&write_timeout=10&alt_hosts="
	flag.StringVar(&cfg.ChDSN, "ch.dsn", ddsn,
		"The clickhouse server DSN to write to eg. "+
			"tcp://host1:9000?username=user&password=qwerty&database=clicks&"+
			"read_timeout=10&write_timeout=20&alt_hosts=host2:9000,host3:9000 "+
			"(see https://github.com/kshvakov/clickhouse).",
	)

	// clickhouse db
	flag.StringVar(&cfg.ChDB, "ch.db", "metrics",
		"The clickhouse database to write to.",
	)

	// clickhouse table
	flag.StringVar(&cfg.ChTable, "ch.table", "samples",
		"The clickhouse table to write to.",
	)

	// clickhouse insertion batch size
	flag.IntVar(&cfg.ChBatch, "ch.batch", 8192,
		"Clickhouse write batch size (n metrics).",
	)

	// channel buffer size between http server => clickhouse writer(s)
	flag.IntVar(&cfg.ChanSize, "ch.buffer", 8192,
		"Maximum internal channel buffer size (n requests).",
	)

	// http listen address
	flag.StringVar(&cfg.HTTPAddr, "web.address", ":9201",
		"Address to listen on for web endpoints.",
	)

	// http prometheus remote write endpoint
	flag.StringVar(&cfg.HTTPWritePath, "web.write", "/write",
		"Address to listen on for remote write requests.",
	)

	// http prometheus metrics endpoint
	flag.StringVar(&cfg.HTTPMetricsPath, "web.metrics", "/metrics",
		"Address to listen on for metric requests.",
	)

	// http shutdown and request timeout
	flag.DurationVar(&cfg.HTTPTimeout, "web.timeout", 30*time.Second,
		"The timeout to use for HTTP requests and server shutdown. Defaults to 30s.",
	)

	flag.Parse()

	return cfg
}
