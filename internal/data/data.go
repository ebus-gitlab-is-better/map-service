package data

import (
	"map-service/internal/conf"

	"github.com/Nerzal/gocloak/v13"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/mojixcoder/gosrm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewOSRMClient, NewKeycloak, NewKeyCloakAPI)

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

func NewOSRMClient(c *conf.Data) gosrm.OSRMClient {
	client, err := gosrm.New(c.Osrm)
	if err != nil {
		panic(err)
	}
	return client
}

func NewKeycloak(c *conf.Data) *gocloak.GoCloak {
	client := gocloak.NewClient(c.Keycloak.Hostname)
	return client
}
