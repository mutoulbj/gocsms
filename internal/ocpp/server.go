package ocpp

import (
	"context"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"gocsms/internal/services"
)

type Server struct {
	addr    string
	handler *OCPPHandler
	log     *logrus.Logger
	server  *http.Server
	clients map[string]*websocket.Conn
	mu      sync.RWMutex
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity; adjust in production
	},
}

func NewOCPPServer(svc *services.ChargePointService, log *logrus.Logger) *Server {
	return &Server{
		handler: NewOCPPHandler(svc, log),
		log:     log,
		clients: make(map[string]*websocket.Conn),
	}
}

func (s *Server) Start(port string) {
	s.addr = ":" + port
	s.server = &http.Server{Addr: s.addr}
	http.HandleFunc("/ws", s.handleWebSocket)
	s.log.Infof("Starting OCPP server on %s", s.addr)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.log.Fatal("Failed to start OCPP server: ", err)
	}
}

func (s *Server) Stop() {
	s.log.Info("Stopping OCPP server")
	if err := s.server.Shutdown(context.Background()); err != nil {
		s.log.Error("Error shutting down OCPP server: ", err)
	}
	s.mu.Lock()
	for id, conn := range s.clients {
		conn.Close()
		delete(s.clients, id)
	}
	s.mu.Unlock()
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	chargePointID := r.URL.Query().Get("id")
	if chargePointID == "" {
		http.Error(w, "Missing charge point ID", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.log.Error("Failed to upgrade to WebSocket: ", err)
		return
	}

	s.mu.Lock()
	s.clients[chargePointID] = conn
	s.mu.Unlock()

	s.log.Infof("Charge point %s connected", chargePointID)
	defer func() {
		s.mu.Lock()
		delete(s.clients, chargePointID)
		s.mu.Unlock()
		conn.Close()
		s.log.Infof("Charge point %s disconnected", chargePointID)
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.log.Error("WebSocket read error: ", err)
			return
		}

		resp, err := s.handler.HandleMessage(r.Context(), chargePointID, msg)
		if err != nil {
			s.log.Error("Failed to handle OCPP message: ", err)
			continue
		}

		if err := conn.WriteMessage(websocket.TextMessage, resp); err != nil {
			s.log.Error("WebSocket write error: ", err)
			return
		}
	}
}