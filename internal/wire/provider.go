package wire

import (
	"location-service/internal/storage/redis"
	"location-service/internal/storage/wrapper"

	"location-service/internal/stream"
	"location-service/internal/tracking"

	"location-service/internal/transport/grpc"
	"location-service/internal/transport/websocket"
	"location-service/internal/types"
)

type Provider struct {
	config         types.Config
	serviceDrain   types.Drain
	transportDrain types.Drain
}

func NewProvider(cfg types.Config) *Provider {
	return &Provider{
		config: cfg,
	}
}

///////////////////////////// providing db ///////////////////////////////////

func (p *Provider) ProvideKeyDB() *redis.KeyDB {
	return redis.NewKeyDB(
		redis.NewPool(&p.config.RedisKeyDB),
	)
}

func (p *Provider) ProvideGeoDB() *redis.GeoDB {
	return redis.NewGeoDB(
		redis.NewPool(&p.config.RedisGeoDB),
	)
}

///////////////////////////// providing store ////////////////////////////////

func (p *Provider) ProvideCourierStore() *wrapper.CourierStore {
	return wrapper.NewCourierStore(
		p.provideItemStore(&p.config.Stores.Courier),
	)
}

func (p *Provider) ProvideOrderStore() *wrapper.OrderStore {
	return wrapper.NewOrderStore(
		p.provideItemStore(&p.config.Stores.Order),
	)
}

func (p *Provider) provideItemStore(cfg *types.StoreConfig) *wrapper.ItemStore {
	return wrapper.NewItemStore(
		p.ProvideKeyDB(),
		p.ProvideGeoDB(),
		cfg,
	)
}

///////////////////////////// providing service //////////////////////////////

func (p *Provider) ProvideTrackingService() *tracking.Service {
	p.ProvideDrain()

	return tracking.NewService(
		p.ProvideCourierStore(),
		p.ProvideOrderStore(),
		p.serviceDrain,
	)
}

func (p *Provider) ProvideDrain() {
	p.serviceDrain = stream.NewDrain()
	p.transportDrain = stream.NewDrain()

	// set output channel of one drain as input channel of the other
	p.serviceDrain.SetInput(p.transportDrain.GetOutput())
	p.transportDrain.SetInput(p.serviceDrain.GetOutput())
}

///////////////////////////// providing transport ////////////////////////////

func (p *Provider) ProvideSocketHandler() *websocket.SocketHandler {
	// inject service to socket handler

	return websocket.NewSocketHandler(
		p.ProvideTrackingService(),
		p.transportDrain,
		p.config.WsServer,
	)
}

func (p *Provider) ProvideGrpcHandler() *grpc.GrpcHandler {
	return grpc.NewGrpcHandler(
		p.ProvideTrackingService(),
		&p.config.GrpcServer,
	)
}
