package client

import (
	"errors"
	"fmt"
	"github.com/illuminatingKong/kongming-kit/base/logx"
	pb "github.com/illuminatingKong/kongming-kit/grpc/task/node/grpc/task/proto"
	"github.com/illuminatingKong/kongming-kit/grpc/task/pool"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
	"time"
)

var (
	taskMap sync.Map
)

func generateTaskUniqueKey(ip string, port int, id int64) string {
	return fmt.Sprintf("%s:%d:%d", ip, port, id)
}

func Send(ip string, port int, taskReq *pb.TaskRequest, log logx.Logger) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("panic#rpc/client.go:Exec#", err)
		}
	}()
	addr := fmt.Sprintf("%s:%d", ip, port)
	c, err := pool.Pool.Get(addr)
	if err != nil {
		return "", err
	}
	if taskReq.Timeout <= 0 || taskReq.Timeout > 86400 {
		taskReq.Timeout = 86400
	}
	timeout := time.Duration(taskReq.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	taskUniqueKey := generateTaskUniqueKey(ip, port, taskReq.Id)
	taskMap.Store(taskUniqueKey, cancel)
	defer taskMap.Delete(taskUniqueKey)

	resp, err := c.Run(ctx, taskReq)
	if err != nil {
		return parseGRPCError(err)
	}

	if resp.Error == "" {
		return resp.Output, nil
	}

	return resp.Output, errors.New(resp.Error)
}

func parseGRPCError(err error) (string, error) {
	switch status.Code(err) {
	case codes.Unavailable:
		return "", errors.New("服务不可用")
	case codes.DeadlineExceeded:
		return "", errors.New("执行超时, 强制结束")
	case codes.Canceled:
		return "", errors.New("手动停止")
	}
	return "", err
}
