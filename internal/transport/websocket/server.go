package websocket

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"location-service/internal/types"
)

type SocketHandler struct {
	config   WsServerConfig
	upgrader *websocket.Upgrader
	client   *websocket.Conn
	service  types.TrackingService
	drain    types.Drain
}

type WsServerConfig struct {
	ReadDeadline int
	ReadTimeout  int
	WriteTimeout int
	MsgSizeLimit int
	Addr         string
	Path         string
}

func NewSocketHandler(svc types.TrackingService, dn types.Drain, cfg types.WsServerConfig) *SocketHandler {
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return &SocketHandler{
		upgrader: upgrader,
		config:   setConfig(cfg),
		client:   nil,
		service:  svc,
		drain:    dn,
	}
}

func setConfig(cfg types.WsServerConfig) WsServerConfig {
	return WsServerConfig{
		ReadDeadline: cfg.ReadDeadline,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		MsgSizeLimit: cfg.MsgSizeLimit,
		Addr:         cfg.Addr,
		Path:         cfg.Path,
	}
}

func (sh *SocketHandler) Serve() {
	http.HandleFunc(sh.config.Path, sh.handleConnection)

	svr := &http.Server{
		Addr:         sh.config.Addr,
		ReadTimeout:  time.Duration(sh.config.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(sh.config.WriteTimeout) * time.Second,
	}

	log.Printf("Websocket server started...")
	log.Fatal(svr.ListenAndServe())
}

func (sh *SocketHandler) handleConnection(w http.ResponseWriter, r *http.Request) {
	// call controller?? -> auth before upgrading the connection

	conn, err := sh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Fprintf(w, "Error upgrading connection")
		log.Fatalf(err.Error())
	}

	if err := sh.handleRegisterCourier(w, r); err != nil {
		log.Fatalf(err.Error())
	}

	sh.client = conn
	sh.handleMessage()
}

func (sh *SocketHandler) handleMessage() {
	log.Printf("Client connected!")

	for {
		// conn.SetReadDeadline(time.Now().Add(sh.readDeadline * time.Second))

		msgType, msg, err := sh.client.ReadMessage()
		if err != nil {
			log.Printf("Client Disconnected!")
			return
		}

		log.Printf("%s sent: %s\n", sh.client.RemoteAddr(), string(msg))
		sh.handleTrackCourier(msgType, msg)

		if err = sh.client.WriteMessage(msgType, msg); err != nil {
			log.Printf("Client Disconnected!")
			return
		}
	}
}

func (sh *SocketHandler) handleRegisterCourier(w http.ResponseWriter, r *http.Request) error {
	// start background task to track courier on new connection
	return sh.service.TrackCourier("id")
}

func (sh *SocketHandler) handleTrackCourier(msgType int, msg []byte) {
	// auth and middleware schema validation
	log.Printf("MSG: " + string(msg))

	dto := types.TrackCourierDTO{}
	types.UnmarshalJSON(string(msg), dto)
	sh.drain.Send(dto)

	sh.client.WriteMessage(msgType, []byte("location updated!"))
}

// // map of key and callback function
// websocketHandlers = make(map[string]func(*websocket.Conn, []byte) error)
//
// // in wiring (functions similar to http handler, not service functions)
// // key in map is specified in websocket message (type:)
// func setupWebSocketHandlers() {
// 	websocketHandlers["broadcast"] = broadcastHandle
// 	websocketHandlers["remove"] = removeHandle
// 	websocketHandlers["update"] = updateHandle
// }
//
// // Check if the applied Type is a handler inside our handle MAP, if the map contains the handle key it will continue inside the if
// if fn := websocketHandlers[msg.Type]; f != nil {
// // If map contains the handle, it will return the function(handler) into fn.
// // Callback is the applied function so f is a function which we here run, we pass the msg.data field in as a Byte array.
// if err := f(ws, []byte(msg.Data)); err != nil {
