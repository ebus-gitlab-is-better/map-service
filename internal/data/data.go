package data

import (
	"map-service/internal/conf"
	"map-service/pkg/valhalla"

	"github.com/Nerzal/gocloak/v13"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewValhallaClient, NewKeycloak, NewKeyCloakAPI)

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
	return client
}
