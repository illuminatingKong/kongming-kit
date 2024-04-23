package main

import (
	"context"
	pb "github.com/illuminatingKong/kongming-kit/grpc/task/node/grpc/task/proto"
	"github.com/illuminatingKong/kongming-kit/runner"
	"github.com/oklog/run"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	projectName             = "task-executer"
	addr                    = "127.0.0.1:8083"
	NoConfigurationInstance = "no find configuration instance"
	configDir               = "$HOME/workspace/illuminatingKong/kongming-kit/examples/example-grpc-task/task-executer"
)

var KeepAliveParams = keepalive.ServerParameters{
	MaxConnectionIdle: 30 * time.Second,
	Time:              30 * time.Second,
	Timeout:           3 * time.Second,
}

var KeepAlivePolicy = keepalive.EnforcementPolicy{
	MinTime:             10 * time.Second,
	PermitWithoutStream: true,
}

type GrpcBootstrap struct {
	Project runner.Options
	Srv     *GrpcSrv
}

type GrpcSrv struct {
	TcpServer        net.Listener
	GrpcServerEngine *grpc.Server
}

func (b GrpcBootstrap) Start() {
	o := b.Project
	var g run.Group
	//var once sync.Once
	type closeOnce struct {
		C     chan struct{}
		once  sync.Once
		Close func()
	}
	reloadReady := &closeOnce{
		C: make(chan struct{}),
	}
	reloadReady.Close = func() {
		reloadReady.once.Do(func() {
			close(reloadReady.C)
		})
	}

	{
		// termination handler.
		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM)
		cancel := make(chan struct{})
		g.Add(
			func() error {
				select {
				case <-term:
					o.Logger.Warn("received sigterm, exiting gracefully...")
					reloadReady.Close()
				case <-cancel:
					reloadReady.Close()
				}
				return nil
			},
			func(err error) {
				close(cancel)

			},
		)
	}

	{
		// http server handler.
		g.Add(
			func() error {
				defer o.Logger.Info("grpc server is returned")
				err := b.Srv.GrpcServerEngine.Serve(b.Srv.TcpServer)
				return err
			},
			func(err error) {
				defer func() {
					if r := recover(); r != nil {
						err = b.Srv.TcpServer.Close()
						if err != nil {
							o.Logger.Info("grpc serve force exit")
							os.Exit(127)
						}
					}
				}()
				b.Srv.GrpcServerEngine.GracefulStop()
			},
		)
	}

}

type NodeServer struct {
	pb.UnimplementedTaskServer
	Rpc *GrpcSrv
}

func main() {
	o := runner.NewContainer(projectName, addr).NewConfig(
		configDir, "yaml", projectName)
	var once sync.Once
	err := o.InitBase(context.Background(), &once)
	if err != nil {
		panic(err)
	}
	project := NewGrpcBind(o)
	project.Start()
}

func NewGrpcBind(o *runner.Options) *GrpcBootstrap {
	if o.Config == nil {
		panic(NoConfigurationInstance)
	}
	r := &GrpcBootstrap{
		Project: *o,
	}

	defer func() {
		if p := recover(); p != nil {
			o.Logger.Error(p)
		}
	}()

	var err error
	r.Srv.TcpServer, err = net.Listen("tcp", o.Addr)
	if err != nil {
		panic(err)
	}
	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(KeepAliveParams),
		grpc.KeepaliveEnforcementPolicy(KeepAlivePolicy),
	}
	r.Srv.GrpcServerEngine = grpc.NewServer(opts...)
	pb.RegisterTaskServer(r.Srv.GrpcServerEngine, NodeServer{Rpc: r.Srv})
	//engine, tcp := LoadGrpcService(o.Addr)

	return r
}

func LoadGrpcService(serverAddr string) (*grpc.Server, net.Listener) {
	s := &GrpcSrv{}
	var err error
	s.TcpServer, err = net.Listen("tcp", serverAddr)
	if err != nil {
		panic(err)
	}
	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(KeepAliveParams),
		grpc.KeepaliveEnforcementPolicy(KeepAlivePolicy),
	}
	s.GrpcServerEngine = grpc.NewServer(opts...)
	pb.RegisterTaskServer(s.GrpcServerEngine, NodeServer{Rpc: s})
	return s.GrpcServerEngine, s.TcpServer
}
