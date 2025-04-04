package handle

import (
	"blog-go/internal/application/service"
	"blog-go/internal/domain/entity"
	"blog-go/internal/interfaces/response"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type ChatHandle struct {
	chatService service.ChatService
	logger      *zap.Logger
}

func NewChatHandle(chatService service.ChatService, logger *zap.Logger) *ChatHandle {
	return &ChatHandle{
		chatService: chatService,
		logger:      logger,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ClientRequest struct {
	Option string          `json:"option"`
	Data   json.RawMessage `json:"data"`
}

type Client struct {
	conn     *websocket.Conn
	username string
	send     chan []byte
}

type Hub struct {
	clients    map[*Client]struct{}
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewClient(conn *websocket.Conn, username string) *Client {
	return &Client{
		conn:     conn,
		username: username,
		send:     make(chan []byte),
	}
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]struct{}),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = struct{}{}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					h.unregister <- client
				}
			}
		}
	}
}

var hub *Hub

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

func (c *ChatHandle) Chat(ctx *gin.Context) {
	if hub == nil {
		hub = newHub()
		go hub.run()
	}

	username, ok := ctx.Get("username")

	if !ok {
		ctx.JSON(response.Unauthorized(nil))
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	client := NewClient(conn, username.(string))
	hub.register <- client
	c.logger.Info("[CHAT] 用户 " + client.username + " 进入聊天室")

	go client.readMsg(c)
	go client.sendMsg(c)
	go client.sendHistory(c)
}

func (client *Client) readMsg(chatHandle *ChatHandle) {
	defer func() {
		hub.unregister <- client
		client.conn.Close()
		chatHandle.logger.Info("[CHAT] 用户 " + client.username + " 退出聊天室")
	}()
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			return
		}

		req := ClientRequest{}
		err = json.Unmarshal(msg, &req)
		if err != nil {
			continue
		}
		if req.Option == "chat" {
			var data entity.Chat
			err := json.Unmarshal(req.Data, &data)
			if err != nil {
				continue
			}
			data.Username = client.username
			data.CreatedAt = time.Now()
			data.Deleted = false
			data.GenChatId()
			chatHandle.logger.Info("[CHAT] 用户 " + client.username + " 发送了消息" + " -> Content: " + data.Content + " Reply: " + data.ReplyChatId)
			go chatHandle.chatService.AddChat(&data)
			dataRaw, _ := json.Marshal(data)
			dataByte, _ := json.Marshal(ClientRequest{
				Option: "chat",
				Data:   dataRaw,
			})
			chatHandle.logger.Info("[CHAT] 广播消息 -> ALL")
			hub.broadcast <- dataByte
		}
	}
}

func (client *Client) sendMsg(chatHandle *ChatHandle) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
		chatHandle.logger.Info("[CHAT] 用户 " + client.username + " 退出聊天室")
	}()
	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := client.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) sendHistory(handle *ChatHandle) {
	chats, err := handle.chatService.GetAll()
	if err != nil {
		return
	}
	chatsRaw, _ := json.Marshal(chats)
	data := ClientRequest{
		Option: "chats",
		Data:   chatsRaw,
	}
	dataByte, _ := json.Marshal(data)
	handle.logger.Info("[CHAT] 发送历史消息 -> " + c.username)
	c.send <- dataByte
}
