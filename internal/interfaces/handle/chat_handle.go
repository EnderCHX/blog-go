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
	"sync"
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

type Msg struct {
	Option string          `json:"option"`
	Data   json.RawMessage `json:"data"`
}

type Client struct {
	conn     *websocket.Conn
	username string
	role     string
	send     chan []byte
	ctrl     chan Msg
	logger   *zap.Logger
}

type Hub struct {
	clients    map[*Client]struct{}
	userconn   map[string]*Client
	userlist   map[string]string
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	logger     *zap.Logger
}

func newClient(conn *websocket.Conn, username string, role string, logger *zap.Logger) *Client {
	return &Client{
		conn:     conn,
		username: username,
		role:     role,
		send:     make(chan []byte),
		ctrl:     make(chan Msg),
		logger:   logger,
	}
}

func newHub(logger *zap.Logger) *Hub {
	return &Hub{
		clients:    make(map[*Client]struct{}),
		userconn:   make(map[string]*Client),
		userlist:   make(map[string]string),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		logger:     logger,
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = struct{}{}
			h.userlist[client.username] = client.role
			h.userconn[client.username] = client
			h.logger.Info("[CHAT] 用户 " + client.username + " 进入聊天室")
			go h.sendUserJoin(client.username)
			go hub.sendUserNum()
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.userlist, client.username)
				delete(h.userconn, client.username)
				close(client.send)
				close(client.ctrl)
				client.conn.Close()
				h.logger.Info("[CHAT] 用户 " + client.username + " 退出聊天室")
				go h.sendUserLeave(client.username)
				go hub.sendUserNum()
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				client.send <- message
				h.logger.Info("[CHAT] 发送消息 -> " + client.username)
			}
		}
	}
}

var hub *Hub
var hubOnce = sync.Once{}

func (c *ChatHandle) Chat(ctx *gin.Context) {
	hubOnce.Do(func() {
		hub = newHub(c.logger)
		go hub.run()
	})

	username, ok := ctx.Get("username")
	role, _ := ctx.Get("role")

	if !ok {
		ctx.JSON(response.Unauthorized(nil))
		return
	}

	if client, ok := hub.userconn[username.(string)]; ok {
		hub.unregister <- client
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	client := newClient(conn, username.(string), role.(string), c.logger)
	hub.register <- client

	go client.readMsg(c)
	go client.sendMsg()
	go client.sendHistory(c)
	go client.heartBeat()
}

func (client *Client) heartBeat() {
	defer func() {
		if err := recover(); err != nil {
			client.logger.Error("[CHAT] 心跳失败 -> " + client.username + " -> " + err.(error).Error())
		}
		hub.unregister <- client
	}()
	ticker1 := time.NewTicker(time.Second * 120)
	defer ticker1.Stop()
	for {
		select {
		case <-ticker1.C:
			client.send <- func() (data []byte) {
				data, _ = json.Marshal(Msg{Option: "ping"})
				return
			}()
			client.logger.Info("[CHAT] 发送心跳 -> " + client.username)
			ticker2 := time.NewTicker(time.Second * 90)
			defer ticker2.Stop()
			select {
			case <-ticker2.C:
				client.logger.Info("[CHAT] 心跳失败 -> " + client.username)
				return
			case ctrl := <-client.ctrl:
				if ctrl.Option == "pong" {
					client.logger.Info("[CHAT] 用户 " + client.username + " 响应心跳")
					continue
				}
			}
		}
	}
}
func (client *Client) readMsg(chatHandle *ChatHandle) {
	defer func() {
		hub.unregister <- client
	}()
	for {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			return
		}

		req := Msg{}
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
			dataByte, _ := json.Marshal(Msg{
				Option: "chat",
				Data:   dataRaw,
			})
			chatHandle.logger.Info("[CHAT] 广播消息 -> ALL")
			hub.broadcast <- dataByte
		} else if req.Option == "pong" {
			client.ctrl <- req
		}
	}
}

func (client *Client) sendMsg() {
	defer func() {
		hub.unregister <- client
	}()
	for {
		select {
		case message, ok := <-client.send:
			if !ok {
				return
			}
			if err := client.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
			client.logger.Info("[CHAT] 发送消息 -> " + client.username)
		}
	}
}

func (c *Client) sendHistory(handle *ChatHandle) {
	defer func() {
		if err := recover(); err != nil {
			handle.logger.Error("[CHAT] 发送历史消息失败 -> " + c.username + " -> " + err.(error).Error())
		}
	}()
	chats, err := handle.chatService.GetAll()
	if err != nil {
		return
	}
	chatsRaw, _ := json.Marshal(chats)
	data := Msg{
		Option: "chats",
		Data:   chatsRaw,
	}
	dataByte, _ := json.Marshal(data)
	handle.logger.Info("[CHAT] 发送历史消息 -> " + c.username)

	c.send <- dataByte
}

func (h *Hub) sendUserNum() {
	defer func() {
		if err := recover(); err != nil {
			h.logger.Error("[CHAT] 发送用户数量消息失败 -> " + "ALL" + " -> " + err.(error).Error())
		}
	}()
	data := struct {
		UserNum int               `json:"user_num"`
		Users   map[string]string `json:"users"`
	}{
		UserNum: len(hub.clients),
		Users:   h.userlist,
	}
	msg := Msg{
		Option: "user_num",
		Data: func() json.RawMessage {
			dataByte, _ := json.Marshal(data)
			return dataByte
		}(),
	}
	dataByte, _ := json.Marshal(msg)
	h.broadcast <- dataByte
}

func (h *Hub) sendUserJoin(username string) {
	defer func() {
		if err := recover(); err != nil {
			h.logger.Error("[CHAT] 发送用户加入消息失败 -> " + username + " -> " + err.(error).Error())
		} else {
			h.logger.Info("[CHAT] 发送用户加入消息 -> " + username)
		}
	}()
	msg := Msg{
		Option: "user_join",
		Data: func() json.RawMessage {
			dataByte, _ := json.Marshal(struct {
				Username string `json:"username"`
			}{
				Username: username,
			})
			return dataByte
		}(),
	}
	dataByte, _ := json.Marshal(msg)
	h.broadcast <- dataByte
}

func (h *Hub) sendUserLeave(username string) {
	defer func() {
		if err := recover(); err != nil {
			h.logger.Error("[CHAT] 发送用户离开消息失败 -> " + username + " -> " + err.(error).Error())
		}
	}()
	msg := Msg{
		Option: "user_leave",
		Data: func() json.RawMessage {
			dataByte, _ := json.Marshal(struct {
				Username string `json:"username"`
			}{
				Username: username,
			})
			return dataByte
		}(),
	}
	dataByte, _ := json.Marshal(msg)
	h.broadcast <- dataByte
}
