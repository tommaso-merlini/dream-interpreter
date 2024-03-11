package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/kr/pretty"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"

	"github.com/tommaso-merlini/dream-interpreter/claude"
	"github.com/tommaso-merlini/dream-interpreter/view/chat"
)

var emptyMessage = errors.New("empty Message")

type Message struct {
	role    string
	content string
}

var clients = make(map[*websocket.Conn][]Message)

func ChatShow(c echo.Context) error {
	return render(c, chat.Chat())
}

func DeleteThinkingMessage(c echo.Context) error {
	return nil
}

func ChatWS(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		clients[ws] = []Message{}
		for {
			in := ""
			err := websocket.Message.Receive(ws, &in)
			if err != nil {
				c.Logger().Error(err)
				return
			}
			msg, err := getMessage(in)
			if err != nil {
				c.Logger().Error(err)
				continue
			}
			err = resendMessageToUser(ws, msg)
			if err != nil {
				c.Logger().Error(err)
				continue
			}
			err = respondToUser(ws, msg)
			if err != nil {
				c.Logger().Error(err)
				continue
			}
			pretty.Println(clients[ws])
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func respondToUser(user *websocket.Conn, msg string) error {
	bufferTextInput := &bytes.Buffer{}
	resp, err := claude.GetMessage(getHistoryMap(user), msg)
	if err != nil {
		return err
	}
	m := chat.Message(resp, false)
	buffer := &bytes.Buffer{}
	m.Render(context.Background(), buffer)
	err = websocket.Message.Send(
		user,
		buffer.String()+bufferTextInput.String(),
	)
	if err != nil {
		return err
	}
	clients[user] = append(clients[user], Message{role: "assistant", content: resp})
	return nil
}

func resendMessageToUser(user *websocket.Conn, msg string) error {
	freshTextInput := chat.Input("")
	bufferTextInput := &bytes.Buffer{}
	freshTextInput.Render(context.Background(), bufferTextInput)

	m := chat.Message(msg, true)
	buffer := &bytes.Buffer{}
	m.Render(context.Background(), buffer)

	t := chat.ThinkingMessage()
	bufferThinking := &bytes.Buffer{}
	t.Render(context.Background(), bufferThinking)

	err := websocket.Message.Send(
		user,
		buffer.String()+bufferTextInput.String()+bufferThinking.String(),
	)
	if err != nil {
		return err
	}
	clients[user] = append(clients[user], Message{role: "user", content: msg})
	return nil
}

func getMessage(msg string) (string, error) {
	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(msg), &jsonMap)
	Message, ok := jsonMap["chat_message"].(string)
	if !ok {
		return "", errors.New("invalid Message")
	}
	if Message == "" {
		return "", emptyMessage
	}
	return Message, nil
}

func broadcastMessage(msg string) {
	for client := range clients {
		err := websocket.Message.Send(client, msg)
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}

func getHistoryMap(ws *websocket.Conn) []map[string]string {
	history := []map[string]string{}
	client := clients[ws]
	for _, message := range client {
		history = append(
			history,
			map[string]string{"role": message.role, "content": message.content},
		)
	}
	return history
}
