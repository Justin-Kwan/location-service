package websocket

// import (
// 	"testing"
//
// 	"github.com/stretchr/testify/assert"
//
// 	"order-matching/internal/config"
// )

// func TestSetConfig(t *testing.T) {
// 	YAMLTextData1 := `
//   websocket_server:
//     read_deadline: 10
//     read_timeout: 10
//     write_timeout: 10
//     message_size_limit: 10
//     address: "127.0.0.1:8010"
//     path: "/path10"`
//
//   // setup
// 	YAMLByteData1 := []byte(YAMLTextData1)
// 	cfg := config.unmarshalYAML(YAMLByteData1)
//
// 	// function under test
// 	wsCfg := setConfig(&(*cfg).WsServer)
// 	assert.Equal(t, *wsCfg,
// 		WsServerConfig{
// 			ReadDeadline:     10,
// 			ReadTimeout:      10,
// 			WriteTimeout:     10,
// 			MessageSizeLimit: 10,
// 			Addr:             "127.0.0.1:8010",
// 			Path:             "/path10",
// 		}, "should be new converted web server config struct")
//
// }
