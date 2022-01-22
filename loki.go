package loki_client_go

import (
	"encoding/json"
	"github.com/ShugetsuSoft/loki-client-go/lib"
	"github.com/dghubble/sling"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const LogPushPath = "loki/api/v1/push"

type LokiClient struct {
	base   *sling.Sling
	labels lib.LabelInfos
	logs   lib.LogInfos
	lock   sync.Mutex
}

func NewLokiClient(uri string) *LokiClient {
	base := sling.New().Base(uri)
	return &LokiClient{
		base:   base,
		labels: lib.LabelInfos{},
		logs:   lib.LogInfos{},
	}
}

func (cli *LokiClient) Push() error {
	if len(cli.labels) < 1 {
		return nil
	}
	cli.lock.Lock()
	defer cli.lock.Unlock()
	streams := make([]lib.PushStream, len(cli.labels))
	i := 0
	for key := range cli.labels {
		values := make([]lib.Value, len(cli.logs[key]))
		for j := range cli.logs[key] {
			values[j][0] = strconv.FormatInt(cli.logs[key][j].Time.UnixNano(), 10)
			values[j][1] = cli.logs[key][j].Log
		}
		streams[i].Stream = cli.labels[key]
		streams[i].Values = values
		i++
	}
	cli.labels = make(lib.LabelInfos)
	cli.logs = make(lib.LogInfos)

	body := lib.PushRequest{
		Streams: streams,
	}
	req, err := cli.base.New().Post(LogPushPath).BodyJSON(body).Request()
	if err != nil {
		return err
	}
	client := &http.Client{}
	_, err = client.Do(req)
	return err
}

func (cli *LokiClient) RunPush() chan error {
	errs := make(chan error)
	go func() {
		err := cli.Push()
		if err != nil {
			errs <- err
		}
		time.Sleep(time.Second * 10)
	}()
	return errs
}

func (cli *LokiClient) WriteLog(lable lib.Label, log string, times time.Time) error {
	lableBytes, err := json.Marshal(lable)
	if err != nil {
		return err
	}
	lableStr := lib.StringOut(lableBytes)
	logEntry := lib.Log{
		Time: times,
		Log:  log,
	}
	cli.lock.Lock()
	defer cli.lock.Unlock()
	if _, ok := cli.labels[lableStr]; ok {
		cli.logs[lableStr] = append(cli.logs[lableStr], logEntry)
	} else {
		cli.labels[lableStr] = lable
		cli.logs[lableStr] = []lib.Log{logEntry}
	}
	return nil
}
