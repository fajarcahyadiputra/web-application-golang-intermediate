package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebScoketConnection struct {
	*websocket.Conn
}
type WSPayload struct {
	Action      string              `json:"action"`
	Message     string              `json:"message"`
	UserName    string              `json:"username"`
	MessageType string              `json:"message_type"`
	UserID      int                 `json:"user_id"`
	Conn        WebScoketConnection `json:"-"`
}

type WSJsonResponse struct {
	Action  string `json:"action"`
	Message string `json:"message"`
	UserID  int    `json:"user_id"`
}

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var clients = make(map[WebScoketConnection]string)

var WSChan = make(chan WSPayload)

func (app *application) WSEndpoint(w http.ResponseWriter, r *http.Request) {
	wupgrade := w
	if u, ok := w.(interface{ Unwrap() http.ResponseWriter }); ok {
		wupgrade = u.Unwrap()
	}

	ws, err := upgradeConnection.Upgrade(wupgrade, r, nil)

	if err != nil {
		app.errorLog.Println(err)
		return
	}

	app.infoLog.Println(fmt.Printf("Client connected from %s", r.RemoteAddr))
	var response WSJsonResponse
	response.Message = "Connected To Server"

	err = ws.WriteJSON(response)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	conn := WebScoketConnection{Conn: ws}
	clients[conn] = ""

	go app.ListenForWS(&conn)
}

func (app *application) ListenForWS(conn *WebScoketConnection) {
	defer func() {
		if r := recover(); r != nil {
			app.errorLog.Println("EROR:", fmt.Sprintf("%v", r))
		}
	}()

	var payload WSPayload

	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			//do nothing
		} else {
			payload.Conn = *conn
			WSChan <- payload
		}
	}
}

func (app *application) ListenToWSChannel() {
	var response WSJsonResponse

	for {
		e := <-WSChan
		switch e.Action {
		case "deleteUser":
			response.Action = "logout"
			response.Message = "Your Account Has Been deleted"
			response.UserID = e.UserID
			app.broadcastToAll(response)
		default:
		}
	}
}

func (app *application) broadcastToAll(response WSJsonResponse) {
	for client := range clients {
		//broadcast to every connection client
		err := client.WriteJSON(response)
		if err != nil {
			app.errorLog.Printf("Websocket err on %s: %s", response.Action, err)
			_ = client.Close()
			delete(clients, client)
		}
	}
}
