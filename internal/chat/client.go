package chat

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/nhh-7/go-chat/internal/dao"
	"github.com/nhh-7/go-chat/internal/model"
	myKafka "github.com/nhh-7/go-chat/internal/service/kafka"
	"github.com/segmentio/kafka-go"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nhh-7/go-chat/internal/config"
	"github.com/nhh-7/go-chat/internal/dto/request"
	"github.com/nhh-7/go-chat/pkg/constants"
	"github.com/nhh-7/go-chat/pkg/enum/message/message_status_enum"
	"github.com/nhh-7/go-chat/utils/zlog"
)

type MessageBack struct {
	Message []byte
	Uuid    string
}

type Client struct {
	Conn     *websocket.Conn
	Uuid     string
	SendTo   chan []byte
	SendBack chan *MessageBack
}

// 将http连接升级为websocket连接
var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var ctx = context.Background()

var messageMode = config.GetConfig().KafkaConfig.MessageMode

func (c *Client) Read() {
	zlog.Info("ws read goroutine start")
	for {
		_, jsonMessage, err := c.Conn.ReadMessage()
		if err != nil {
			zlog.Error(err.Error())
			return
		} else {
			var message = request.ChatMessageRequest{}
			if err := json.Unmarshal(jsonMessage, &message); err != nil {
				zlog.Error(err.Error())
			}
			log.Println("接收到消息：", string(jsonMessage))
			if messageMode == "channel" {
				for len(ChatServer.Transmit) < constants.CHANNEL_SIZE && len(c.SendTo) > 0 {
					sendToMessage := <-c.SendTo
					ChatServer.SendMessageToTransmit(sendToMessage)
				}
				if len(ChatServer.Transmit) < constants.CHANNEL_SIZE {
					ChatServer.SendMessageToTransmit(jsonMessage)
				} else if len(c.SendTo) < constants.CHANNEL_SIZE {
					c.SendTo <- jsonMessage
				} else {
					// 否则考虑加宽channel size，或者使用kafka
					if err := c.Conn.WriteMessage(websocket.TextMessage, []byte("由于目前同一时间过多用户发送消息，消息发送失败，请稍后重试")); err != nil {
						zlog.Error(err.Error())
					}
				}
			} else {
				if err := myKafka.KafkaService.ChatWriter.WriteMessages(ctx, kafka.Message{
					Key:   []byte(strconv.Itoa(config.GetConfig().KafkaConfig.Partition)),
					Value: jsonMessage,
				}); err != nil {
					zlog.Error(err.Error())
				}
				zlog.Info("已发送消息：" + string(jsonMessage))
			}
		}
	}
}

func (c *Client) Write() {
	zlog.Info("ws write goroutine start")
	for messageBack := range c.SendBack {
		// 通过 WebSocket 发送消息
		err := c.Conn.WriteMessage(websocket.TextMessage, messageBack.Message)
		if err != nil {
			zlog.Error(err.Error())
			return // 直接断开websocket
		}
		log.Println("已发送消息：", messageBack.Message)
		// 说明顺利发送，修改状态为已发送
		if res := dao.GormDB.Model(&model.Message{}).Where("uuid = ?", messageBack.Uuid).Update("status", message_status_enum.Sent); res.Error != nil {
			zlog.Error(res.Error.Error())
		}
	}
}

func NewClientInit(c *gin.Context, clientId string) {
	kafkaConfig := config.GetConfig().KafkaConfig
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zlog.Error(err.Error())
	}
	client := &Client{
		Conn:     conn,
		Uuid:     clientId,
		SendTo:   make(chan []byte, constants.CHANNEL_SIZE),
		SendBack: make(chan *MessageBack, constants.CHANNEL_SIZE),
	}
	if kafkaConfig.MessageMode == "channel" {
		ChatServer.SendClientToLogin(client)
	} else {
		KafkaChatServer.SendClientToLogin(client)
	}
	go client.Read()
	go client.Write()
	zlog.Info("ws 连接成功")
}

func ClientLogout(clientId string) (string, int) {
	kafkaConfig := config.GetConfig().KafkaConfig
	client := ChatServer.Clients[clientId]
	if client != nil {
		if kafkaConfig.MessageMode == "channel" {
			ChatServer.SendClientToLogout(client)
		} else {
			KafkaChatServer.SendClientToLogout(client)
		}
		if err := client.Conn.Close(); err != nil {
			zlog.Error(err.Error())
			return constants.SYSTEM_ERROR, -1
		}
		close(client.SendTo)
		close(client.SendBack)
	}
	return "退出成功", 0
}
