package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
	pb "github.com/tony-spark/metrico/gen/pb/api"
	"github.com/tony-spark/metrico/internal/crypto"
	"github.com/tony-spark/metrico/internal/dto"
	"github.com/tony-spark/metrico/internal/server/models"
	"github.com/tony-spark/metrico/internal/server/services"
	"google.golang.org/grpc"
)

type Controller struct {
	pb.UnimplementedMetricServiceServer

	srv           *grpc.Server
	listenAddress string
	ms            *services.MetricService
	dbm           models.DBManager
	h             dto.Hasher
	d             crypto.Decryptor
	trustedSubNet *net.IPNet
}

type Option func(c *Controller)

func WithListenAddress(addr string) Option {
	return func(c *Controller) {
		c.listenAddress = addr
	}
}

func WithHasher(h dto.Hasher) Option {
	return func(c *Controller) {
		c.h = h
	}
}

func WithDBManager(dbm models.DBManager) Option {
	return func(c *Controller) {
		c.dbm = dbm
	}
}

func WithDecryptor(d crypto.Decryptor) Option {
	return func(c *Controller) {
		c.d = d
	}
}

func WithTrustedSubNet(subnet *net.IPNet) Option {
	return func(c *Controller) {
		c.trustedSubNet = subnet
	}
}

func NewController(metricService *services.MetricService, options ...Option) *Controller {
	controller := &Controller{
		ms:  metricService,
		srv: grpc.NewServer(),
	}

	for _, opt := range options {
		opt(controller)
	}

	return controller
}

func (c *Controller) DBStatus(ctx context.Context, _ *pb.Empty) (*pb.Response, error) {
	var response pb.Response

	if c.dbm == nil {
		var errTxt = "database connection not configured"
		response.Status = pb.Status_ERROR
		response.Error = &errTxt
		return &response, nil
	}

	ok, err := c.dbm.Check(ctx)
	response.Status = pb.Status_OK
	if !ok || err != nil {
		var errTxt = "could not check DB or DB is not OK"
		response.Status = pb.Status_ERROR
		response.Error = &errTxt
	}

	return &response, nil
}

func (c *Controller) Run() error {
	listen, err := net.Listen("tcp", c.listenAddress)
	if err != nil {
		return fmt.Errorf("could not open listen socker: %w", err)
	}

	pb.RegisterMetricServiceServer(c.srv, c)

	err = c.srv.Serve(listen)
	if err != nil && err != grpc.ErrServerStopped {
		return fmt.Errorf("error running grpc server: %w", err)
	}

	return nil
}

func (c *Controller) Shutdown(ctx context.Context) error {
	stopped := make(chan struct{})
	go func() {
		log.Info().Msg("stopping grpc server gracefully...")
		c.srv.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		log.Warn().Msg("force stopping grpc server")
		c.srv.Stop()
	case <-stopped:
		log.Info().Msg("grpc server stopped gracefully")
	}

	return nil
}

func (c Controller) String() string {
	return fmt.Sprintf("GRPC controller at " + c.listenAddress)
}
