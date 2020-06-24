package wire

// import (
// 	"location-service/internal/storage/redis"
// 	"location-service/internal/storage/wrappers"
// 	"location-service/internal/tracking"
// 	"location-service/internal/transport"
// 	"location-service/internal/types"
// )
//
// type Provider struct {
// 	config types.Config
// 	// logger?
// }
//
// func NewProvider(cfg types.Config) *Provider {
// 	return &Provider{
// 		config: cfg,
// 	}
// }
//
// ///////////////////////////// providing db ///////////////////////////////////
//
// func (p *Provider) ProvideKeyDB() *redis.KeyDB {
// 	return redis.NewKeyDB(
// 		p.ProvideRedisPool(p.config.RedisKeyDB),
// 	)
// }
//
// func (p *Provider) ProvideGeoDB() *redis.GeoDB {
// 	return redis.NewGeoDB(
// 		p.ProvideRedisPool(p.config.RedisGeoDB),
// 	)
// }
//
// func (p *Provider) ProvideRedisPool() *redis.Pool {
// 	return redis.NewPool(p.config.RedisCfg)
// }
//
// ///////////////////////////// providing store ////////////////////////////////
//
// func (p *Provider) ProvideCourierStore() *wrappers.ItemStore {
//   return p.provideItemStore(p.config.Stores.Courier)
// }
//
// func (p *Provider) ProvideOrderStore() *wrappers.ItemStore {
//   return p.provideItemStore(p.config.Stores.Order)
// }
//
// func (p *Provider) provideItemStore(cfg types.StoreConfig) *wrappers.ItemStore {
//   return wrappers.NewItemStore(
//     p.ProvideKeyDB(),
//     p.ProvideGeoDB(),
//     cfg,
//   )
// }
//
// ///////////////////////////// providing service //////////////////////////////
//
// func (p *Provider) ProvideTrackingService() *tracking.Service {
//   return tracking.NewService(
//     ProvideCourierStore(),
//     ProvideOrderStore(),
//     ProvideDrain(),
//   )
// }
//
// func (p *Provider) ProvideDrain() {
//
// }
//
// ///////////////////////////// providing transport ////////////////////////////
