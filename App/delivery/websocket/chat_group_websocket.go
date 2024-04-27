package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/syahyudi09/ChatPintar/App/config"
	"github.com/syahyudi09/ChatPintar/App/model"
	"github.com/syahyudi09/ChatPintar/App/usecase"
)

type GroupChatController struct {
	groupChatUsecase usecase.ChatGroupUsecase
	connections      map[string]map[string]*websocket.Conn
	mu               sync.Mutex
}

func NewGroupChatController(r *gin.Engine, groupChatUsecase usecase.ChatGroupUsecase) *GroupChatController {
	controller := &GroupChatController{
		groupChatUsecase: groupChatUsecase,
		connections:      make(map[string]map[string]*websocket.Conn),
	}

	r.GET("/ws/group", controller.HandleGroupWebSocket)
	return controller
}

func (gcc *GroupChatController) HandleGroupWebSocket(c *gin.Context) {
	ws, err := config.WebSocketConfig.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to upgrade to WebSocket: %v", err)})
		return
	}

	groupID := c.Query("group_id")
	userID := c.Query("user_id")

	if groupID == "" || userID == "" {
		ws.WriteMessage(websocket.TextMessage, []byte("Invalid group_id or user_id"))
		ws.Close()
		return
	}

	isMember, err := gcc.groupChatUsecase.IsUserMemberOfGroup(userID, groupID)
	if err != nil || !isMember {
		ws.WriteMessage(websocket.TextMessage, []byte("User is not a group member"))
		ws.Close()
		return
	}

	gcc.mu.Lock()
	if _, exists := gcc.connections[groupID]; !exists {
		gcc.connections[groupID] = make(map[string]*websocket.Conn)
	}
	gcc.connections[groupID][userID] = ws
	gcc.mu.Unlock()

	defer func() {
		gcc.mu.Lock()
		delete(gcc.connections[groupID], userID)
		gcc.mu.Unlock()
		ws.Close()
	}()

	for {
		_, messageContent, err := ws.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading WebSocket message: %v\n", err)
			break
		}

		var chatMessage model.MessageGroup
		if err := json.Unmarshal(messageContent, &chatMessage); err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte("Invalid message format"))
			continue
		}

		chatMessage.GroupId = groupID
		chatMessage.SenderID = userID
		chatMessage.CreatedAt = time.Now()

		// Simpan pesan ke database
		err = gcc.groupChatUsecase.CreateMessageGroup(chatMessage)
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Failed to save message: %v", err)))
			continue
		}

		// Kirim pesan ke semua anggota grup
		gcc.mu.Lock()
		for _, conn := range gcc.connections[groupID] {
			if conn != ws {
				conn.WriteMessage(websocket.TextMessage, messageContent)
			}
		}
		gcc.mu.Unlock()
	}
}
