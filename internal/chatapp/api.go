package chatapp

import (
	"errors"
	"finalproject/config"
	"log"
	"net/http"
	"strings"

	"fmt"

	"github.com/gorilla/websocket"

	"github.com/labstack/echo/v4"
)

type handler struct {
	cfg     *config.Config
	service IService
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func RegisterAPI(r echo.Group, cfg *config.Config, service IService) {
	handler := handler{cfg: cfg, service: IService(service)}

	r.GET("/", handler.homePage)
	r.GET("/ws", handler.webSocket)
}

func (h handler) homePage(c echo.Context) error {
	return c.Render(200, "index.html", nil)
}

func (h handler) webSocket(c echo.Context) error {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := websocket.Upgrade(c.Response().Writer, c.Request(), c.Response().Header(), 1024, 1024)

	if err != nil {
		return errors.New("error connection")
	}
	defer ws.Close()

	username := c.QueryParam("username")
	currentConn := WebSocketConnection{Conn: ws, Username: username}
	connections = append(connections, &currentConn)

	go h.handleIO(&currentConn, connections)
	return nil
}

func (h handler) handleIO(currentConn *WebSocketConnection, connections []*WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("ERROR", fmt.Sprintf("%v", r))
		}
	}()

	h.broadcastMessage(currentConn, MESSAGE_NEW_USER, "")

	for {
		payload := SocketPayload{}
		err := currentConn.ReadJSON(&payload)
		if err != nil {
			if strings.Contains(err.Error(), "websocket: close") {
				h.broadcastMessage(currentConn, MESSAGE_LEAVE, "")
				return
			}

			log.Println("ERROR", err.Error())
			continue
		}

		h.broadcastMessage(currentConn, MESSAGE_CHAT, payload.Message)
	}
}

func (h handler) broadcastMessage(currentConn *WebSocketConnection, kind, message string) {
	for _, eachConn := range connections {
		if eachConn == currentConn {
			continue
		}

		eachConn.WriteJSON(SocketResponse{
			From:    currentConn.Username,
			Type:    kind,
			Message: message,
		})
	}
}
