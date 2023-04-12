package grpc

import (
	"fmt"
	"net"

	"github.com/dhikaroofi/stock.git/internal/core"
	"github.com/dhikaroofi/stock.git/pkg/logger"
	"github.com/dhikaroofi/stock.git/pkg/protobuf/pb"
	"google.golang.org/grpc"
)

type Task interface {
	ListenAndServe(exitSignal chan bool)
}

type server struct {
	port          string
	grpc          *grpc.Server
	coreContainer *core.Container
}

// New is used for initiate data streamer from file or it can be kafka
func New(port string, coreContainer *core.Container) Task {
	return &server{
		port:          port,
		coreContainer: coreContainer,
		grpc:          grpc.NewServer(),
	}
}

func (s server) ListenAndServe(extitSignalch chan bool) {
	go func() {
		lis, err := net.Listen("tcp", s.port)
		if err != nil {
			logger.Fatal(fmt.Sprintf("failed to listen: %s", err.Error()))
		}

		registerServiceServer(s.grpc, s.coreContainer)

		logger.SysInfo("server is running")
		if err := s.grpc.Serve(lis); err != nil {
			logger.Fatal(fmt.Sprintf("failed to serve: %s", err.Error()))
		}

	}()

	go func() {
		<-extitSignalch
		s.grpc.GracefulStop()
		logger.SysInfo("server is shutting down")
		extitSignalch <- true
	}()
}

func registerServiceServer(grpcServer *grpc.Server, coreContainer *core.Container) {
	pb.RegisterOhlcServiceServer(grpcServer, NewOhlcHandler(coreContainer.UseCase.OHLC))
}
