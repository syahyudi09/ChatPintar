package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/syahyudi09/ChatPintar/App/config"
	"github.com/syahyudi09/ChatPintar/App/model"
	"github.com/syahyudi09/ChatPintar/App/usecase"
)

type WebSocketController struct {
	privateChatUsecase usecase.PrivateChatUsecase
	usecase usecase.AuthUsecase
	connections        map[string]*websocket.Conn
	mu                 sync.Mutex
}

func NewWebSocketController(r *gin.Engine, privateChatUsecase usecase.PrivateChatUsecase, usecase usecase.AuthUsecase) *WebSocketController {
	wsController := &WebSocketController{
		privateChatUsecase: privateChatUsecase,
		connections:        make(map[string]*websocket.Conn),
		usecase: usecase,
	}

	r.GET("/ws", wsController.HandleWebSocket)

	return wsController
}

func (wsc *WebSocketController) HandleWebSocket(c *gin.Context) {
	// Meng-upgrade HTTP ke WebSocket
	ws, err := config.WebSocketConfig.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to upgrade to WebSocket: %v", err)})
		return
	}

	senderID := c.Query("user_id")
	receiverID := c.Query("receiver_id") 
	fmt.Println("WebSocket connection initiated with user_id:", senderID)
	if senderID == "" || receiverID == ""{
		ws.WriteMessage(websocket.TextMessage, []byte("Error: Missing user_id"))
		ws.Close()
		return
	}

	senderExists, err := wsc.usecase.PhoneNumberExits(senderID)
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error checking sender phone number"))
		ws.Close()
		return
	}

	receiverExists, err := wsc.usecase.PhoneNumberExits(receiverID)
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Error checking receiver phone number"))
		ws.Close()
		return
	}

	if !senderExists || !receiverExists {
		ws.WriteMessage(websocket.TextMessage, []byte("Error: Invalid phone number for sender or receiver"))
		ws.Close()
		return
	}

	wsc.mu.Lock()
	wsc.connections[senderID] = ws
	wsc.mu.Unlock()

	defer func() {
		wsc.mu.Lock()
		delete(wsc.connections, senderID)
		wsc.mu.Unlock()
		ws.Close()
	}()

	ws.WriteMessage(websocket.TextMessage, []byte("Connected to WebSocket"))

	for {
		_, messageContent, err := ws.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading WebSocket message: %v\n", err)
			break
		}

		// Menguraikan pesan JSON
		var msg model.InputMessageModel
		if err := json.Unmarshal(messageContent, &msg); err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte("Error: Invalid message format"))
			continue
		}

		msg.SenderID = senderID
		msg.ReceiverID = receiverID
		msg.Status = model.Pending
		fmt.Printf("Sending message to receiver_id: %s %s\n", msg.SenderID , senderID)

		// Simpan pesan ke database
		if err := wsc.privateChatUsecase.CreateMessage(msg); err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte("Error creating message: " + err.Error()))
			continue
		}

		// Mengirim pesan ke penerima
		wsc.mu.Lock()
        if receiverConn, exists := wsc.connections[msg.ReceiverID]; exists {
            err = receiverConn.WriteMessage(websocket.TextMessage, messageContent) // Mengirim pesan ke penerima
            if err == nil {
                msg.Status = model.Send// Jika pengiriman berhasil
            } else {
                msg.Status = model.Failed// Jika pengiriman gagal
            }
        } else {
            msg.Status = model.Failed // Jika penerima tidak terhubung
            ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error: Receiver with ID %s not connected", msg.ReceiverID)))
        }
        wsc.mu.Unlock()

        err = wsc.privateChatUsecase.UpdateMessageStatusBySender(msg.SenderID, string(msg.Status))
			if err != nil{
				ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Error updating message status: %v", err)))
				continue
			}

	}
}
