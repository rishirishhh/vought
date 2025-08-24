package controllers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	jsonDTO "github.com/rishirishhh/vought/src/cmd/api/dto/json"

	"github.com/rishirishhh/vought/src/cmd/api/config"
	"github.com/rishirishhh/vought/src/cmd/api/dto/protobuf"
	"github.com/rishirishhh/vought/src/pkg/clients"
	contracts "github.com/rishirishhh/vought/src/pkg/contracts/v1"
)

type WSHandler struct {
	Config                config.Config
	AmqpVideoStatusUpdate clients.AmqpClient
}

func (wsh WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug("Ws Wshandler new connection", r.Host)

	upgrader := websocket.Upgrader{}

	upgrader.CheckOrigin = func(r *http.Request) bool {

		decodedData, err := decodeAuthorization(r)
		if err != nil {
			log.Error("Could not decode data", err)
			return false
		}

		givenUser, givenPass := extractCredentials(decodedData)

		if givenUser == wsh.Config.UserAuth &&
			givenPass == wsh.Config.PwdAuth {
			return true
		}
		return false
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Cannot Upgrade : ", err)
		return
	}

	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte("Connection is a success"))
	if err != nil {
		log.Error("Cannot send message : ", err)
		return
	}

	randomQueueName := wsh.AmqpVideoStatusUpdate.GetRandomQueueName()
	HandleMessage(context.Background(), &wsh, randomQueueName, conn)
}

func decodeAuthorization(r *http.Request) (decodedData []byte, err error) {
	authCookie, err := r.Cookie("Authorization")
	if err != nil {
		log.Error("No Such cookie", err)
	}

	strCookie := authCookie.String()
	auth := strCookie[len("Authorization="):]
	decodedData, err = base64.StdEncoding.DecodeString(auth[len("Basic%20"):])
	if err != nil {
		return nil, err
	}

	return decodedData, nil
}

func extractCredentials(data []byte) (username string, password string) {
	creds := bytes.SplitN(data, []byte(":"), 2)
	givenUser := string(creds[0])
	givenPass := string(creds[1])
	return givenUser, givenPass
}

var HandleMessage = func(ctx context.Context, wsh *WSHandler, randomQueueName string, conn *websocket.Conn) {
	ctx, clear := context.WithCancel(ctx)

	conn.SetCloseHandler(func(code int, text string) error {
		log.Debugf("Connection closed with code %v : %v", code, text)
		clear()
		return nil
	})

	// Read message from client
	go wsh.handleClientMessage(ctx, clear, randomQueueName, conn)

	// Transfer message from queue to Client
	go wsh.handleUpdateMessage(ctx, randomQueueName, conn)

	// Ping client to ensure connection is still needed
	wsh.pingClient(ctx, clear, conn, time.Duration(5)*time.Second)

	conn.Close()
}

func (wsh *WSHandler) handleClientMessage(ctx context.Context, clear context.CancelFunc, randomQueueName string, conn *websocket.Conn) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, msg, err := conn.ReadMessage()
			if err != nil {
				if _, ok := err.(*websocket.CloseError); ok {
					log.Debugf("Close message received.")
					clear()
				} else {
					log.Error("Could not read message : ", err)
				}
			}

			err = wsh.AmqpVideoStatusUpdate.QueueBind(randomQueueName, string(msg))
			if err != nil {
				log.Error("Could no bind queue : ", err)
			}
		}
	}
}

func (wsh *WSHandler) handleUpdateMessage(ctx context.Context, randomQueueName string, conn *websocket.Conn) {
	session := wsh.AmqpVideoStatusUpdate.WithRedial()

	for {
		var client clients.AmqpClient
		select {
		case <-ctx.Done():
			return

		case client = <-session:
			msgs, err := client.Consume(randomQueueName)
			if err != nil {
				log.Error("Failed to consume RabbitMQ client: ", err)
				continue
			}

			for d := range msgs {
				videoProto := &contracts.Video{}
				if err := proto.Unmarshal([]byte(d.Body), videoProto); err != nil {
					log.Error("Fail to unmarshal video event : ", err)
					continue
				}
				video := protobuf.VideoProtobufToVideo(videoProto)
				video.Title = d.RoutingKey
				msg, err := json.Marshal(jsonDTO.VideoToStatusJson(video))
				if err != nil {
					log.Error("Failed to marshall response to front :", err)
				}

				err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
				if err != nil {
					log.Error("Cannot send message : ", err)
					return
				}

				if err := d.Acknowledger.Ack(d.DeliveryTag, false); err != nil {
					log.Error("Failed to Ack message ", video.ID, " - ", err)
					continue
				}
			}
		}

		client.Close()
	}
}

func (wsh *WSHandler) pingClient(ctx context.Context, clear context.CancelFunc, conn *websocket.Conn, timeout time.Duration) {
	lastCheck := time.Now()
	conn.SetPongHandler(func(appData string) error {
		lastCheck = time.Now()
		return nil
	})
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if time.Now().After(lastCheck.Add(timeout)) {
				clear()
				return
			} else {
				err := conn.WriteMessage(websocket.PingMessage, []byte("pingClient"))
				if err != nil {
					log.Error("Could not ping the client : ", err)
				}
			}
		}
		time.Sleep(timeout * 9 / 10)
	}
}
