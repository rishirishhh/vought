package transformer

import (
	"context"
	"fmt"
	"net"
	"os/exec"

	log "github.com/sirupsen/logrus"

	"github.com/rishirishhh/vought/src/pkg/clients"
	"github.com/rishirishhh/vought/src/pkg/transformer/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const MAX_CHUNK_SIZE int = 32000

type ITransformerServer interface {
	StartRPCServer(ctx context.Context, srv transformer.TransformerServiceServer, port uint32) error
	TransformVideo(ctx context.Context, args *transformer.TransformVideoRequest, stream transformer.TransformerService_TransformVideoServer)
	Stop()
}

type TransformerServer struct {
	CreateTransformationCmd func(ctx context.Context) *exec.Cmd
	DiscoveryClient         clients.ServiceDiscovery
	S3Client                clients.IS3Client
}

func (t TransformerServer) createRPCClient(clientName string) (transformer.TransformerServiceClient, error) {
	// Retrieve service address and port
	tfServices, err := t.DiscoveryClient.GetTransformationService(clientName)
	if err != nil {
		log.Errorf("Cannot get address for service name %v : %v", clientName, err)
		return nil, err
	}

	conn, err := grpc.Dial(tfServices, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("Cannot open TCP connection with grpc %v transformer server : %v", clientName, err)
		return nil, err
	}

	return transformer.NewTransformerServiceClient(conn), nil

}

func (t TransformerServer) StartRPCServer(ctx context.Context, srv transformer.TransformerServiceServer, port uint32) error {
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		log.Error("failed to listen : ", err)
		return err
	}

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	// check for context
	go func() {
		<-ctx.Done()
		log.Info("Gracefully shutdown grpcServer\n")
		grpcServer.Stop()
	}()

	transformer.RegisterTransformerServiceServer(grpcServer, srv)
	if err := grpcServer.Serve(listen); err != nil {
		log.Error("Cannot create gRPC server : ", err)
		return err
	}

	return nil

}
