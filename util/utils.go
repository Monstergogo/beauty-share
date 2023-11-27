package util

import (
	"context"
	"errors"
	"gopkg.in/yaml.v3"
	"net"
	"os"
	"sync"
	"time"
)

const (
	workerBits  uint8 = 10
	numberBits  uint8 = 12
	workerMax   int64 = -1 ^ (-1 << workerBits)
	numberMax   int64 = -1 ^ (-1 << numberBits)
	timeShift   uint8 = workerBits + numberBits
	workerShift uint8 = numberBits
	startTime   int64 = 1698400286000
)

type Worker struct {
	mu        sync.Mutex
	timestamp int64
	workerId  int64
	number    int64
}

// NewWorker 雪花算法初始化，用于生成唯一id
func NewWorker(workerId int64) (Worker, error) {
	if workerId < 0 || workerId > workerMax {
		return Worker{}, errors.New("worker ID excess of quantity")
	}
	// 生成一个新节点
	return Worker{
		timestamp: 0,
		workerId:  workerId,
		number:    0,
	}, nil
}

func (w Worker) GetId() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := time.Now().UnixNano() / 1e6
	if w.timestamp == now {
		w.number++
		if w.number > numberMax {
			for now <= w.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		w.number = 0
		w.timestamp = now
	}
	ID := (now-startTime)<<timeShift | (w.workerId << workerShift) | (w.number)
	return ID
}

func GetTraceIdFromCtx(ctx context.Context) string {
	traceId := ctx.Value(CtxTraceID)
	if traceId == nil {
		return ""
	}
	if traceIdStr, ok := traceId.(string); ok {
		return traceIdStr
	}
	return ""
}

// GetOutboundIP 获取本机的出口IP
func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

// GetConfFromYaml
// @desc 从yaml文件中读取配置
func GetConfFromYaml(path string, config interface{}) error {
	dataBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(dataBytes, config)
	return err
}
