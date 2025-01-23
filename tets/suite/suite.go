package suite

import (
	"context"
	"gRPCAuthService/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
	"testing"

	ssov1 "github.com/21Timofei/Contracts/gen/go/sso"
)

const grpcHost = "localhost"

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.AuthClient
}

func NewSuite(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cnf := config.MustLoadByPath("../local_tests/config.yaml")

	ctx, cancelCtx := context.WithTimeout(context.Background(), cnf.GRPC.Timeout)

	t.Cleanup(
		func() {
			t.Helper()
			cancelCtx()
		})

	cc, err := grpc.DialContext(context.Background(), grpcAddress(cnf), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cnf,
		AuthClient: ssov1.NewAuthClient(cc),
	}
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}
