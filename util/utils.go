package util

import (
	"context"
	"errors"
	"net"
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
	CtxTraceId        = "traceId"
)

type Worker struct {
	mu        sync.Mutex
	timestamp int64
	workerId  int64
	number    int64
}

// NewWorker 雪花算法初始化，用于生成唯一id
func NewWorker(workerId int64) (*Worker, error) {
	if workerId < 0 || workerId > workerMax {
		return nil, errors.New("worker ID excess of quantity")
	}
	// 生成一个新节点
	return &Worker{
		timestamp: 0,
		workerId:  workerId,
		number:    0,
	}, nil
}

func (w *Worker) GetId() int64 {
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
	traceId := ctx.Value(CtxTraceId)
	if traceId == nil {
		return ""
	}
	if traceIdStr, ok := traceId.(string); ok {
		return traceIdStr
	}
	return ""
}

// GetCurrIp 获取ip地址
func GetCurrIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", errors.New("can not find the client ip address")
}
