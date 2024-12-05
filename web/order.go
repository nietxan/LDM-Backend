package web

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{}
	users    = make(map[uint]*websocket.Conn)
	mutex    = &sync.Mutex{}                
)

func OrderUpdateSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "WebSocket upgrade failed"})
		return
	}
	defer conn.Close()

	userID, err := getUserIdFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	mutex.Lock()
	users[userID] = conn
	mutex.Unlock()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

	mutex.Lock()
	delete(users, userID)
	mutex.Unlock()
}

func getUserIdFromRequest(c *gin.Context) (uint, error) {
	userIDParam := c.Query("user_id")
	if userIDParam == "" {
		return 0, fmt.Errorf("user_id not found in request")
	}

	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid user_id format")
	}

	return uint(userID), nil
}

func notifyUser(userID uint, message string) {
	mutex.Lock()
	defer mutex.Unlock()

	conn, exists := users[userID]
	if !exists {
		fmt.Printf("User %d not connected\n", userID)
		return
	}

	err := conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		fmt.Printf("Failed to send message to user %d: %v\n", userID, err)
		conn.Close()
		delete(users, userID)
	}
}

func getUserIdByOrderId(orderID uint) (uint, error) {
	var order Order
	result := DB.First(&order, orderID)
	if result.Error != nil {
		return 0, result.Error
	}
	return order.UserID, nil
}

