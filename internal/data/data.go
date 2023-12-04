package data

import (
	"context"
	"crypto/tls"
	"map-service/internal/conf"
	"map-service/pkg/valhalla"
	"time"

	accidentS "map-service/api/accident/v1"

	"github.com/Nerzal/gocloak/v13"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewValhallaClient,
	NewKeycloak,
	NewKeyCloakAPI,
	NewAccidentService,
)

// Data .
type Data struct {
	// TODO wrapped database client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{}, cleanup, nil
}

func NewValhallaClient(c *conf.Data) *valhalla.Client {
	client := valhalla.New(c.Osrm)
	return client
}

func NewKeycloak(c *conf.Data) *gocloak.GoCloak {
	client := gocloak.NewClient(c.Keycloak.Hostname)
	restyClient := client.RestyClient()
	restyClient.SetDebug(true)
	restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	return client
}

func NewAccidentService(c *conf.Data) accidentS.AccidentClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(c.AccidentService),
		grpc.WithMiddleware(
			tracing.Client(),
			recovery.Recovery()),
		grpc.WithTimeout(2*time.Second),
	)
	if err != nil {
		panic(err)
	}
	return accidentS.NewAccidentClient(conn)
}
