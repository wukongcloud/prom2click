package main

import (
	"fmt"
	"github.com/prometheus/prometheus/prompb"
	"testing"
	"time"
)

func TestReader(t *testing.T) {
	conf := &config{
		ChDSN:           "tcp://127.0.0.1:9001?debug=true&username=user1&password=123456&database=default&read_timeout=10&write_timeout=10&alt_hosts=",
		ChDB:            "default",
		ChTable:         "samples",
		ChBatch:         1,
		ChanSize:        8192,
		CHQuantile:      0.75,
		CHMaxSamples:    8192,
		CHMinPeriod:     10,
		HTTPTimeout:     30000000000,
		HTTPAddr:        ":9201",
		HTTPWritePath:   "/write",
		HTTPMetricsPath: "/metrics",
	}

	reader, err := NewP2CReader(conf)
	if err != nil {
		fmt.Printf("Error creating clickhouse writer: %s\n", err.Error())
		t.Error(err)
	}

	req := &prompb.ReadRequest{
		Queries: []*prompb.Query{
			{
				StartTimestampMs: time.Now().Add(-5 * time.Hour).UnixMilli(),
				EndTimestampMs:   time.Now().UnixMilli(),
			},
		},
	}
	reader.Read(req)
	time.Sleep(5 * time.Minute)
}