package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/goGo-service/back/apiproto"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type Client struct {
	centrifugo apiproto.CentrifugoApiClient
}

func NewCentrifugo(centrifugo apiproto.CentrifugoApiClient) *Client {
	return &Client{centrifugo: centrifugo}
}

func (c *Client) GetPresence(channel string) ([]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := c.centrifugo.Presence(ctx, &apiproto.PresenceRequest{Channel: channel})
	if err != nil {
		logrus.Errorf("Transport level error: %v", err)
		return nil, err
	}

	if resp.GetError() != nil {
		logrus.Errorf("centrifugo error: %s", resp.GetError().Message)

		return nil, fmt.Errorf("centrifugo error: %s", resp.GetError().Message)
	}

	var userIDs []int
	for _, info := range resp.GetResult().Presence {
		userId, err := strconv.Atoi(info.GetUser())
		if err != nil {
			logrus.Errorf("error cast string: %s; user: %s", err, info.GetUser())
			continue
		}
		userIDs = append(userIDs, userId)
	}

	return userIDs, nil
}

func (c *Client) PublishMessage(channel string, data interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var payload []byte
	var err error

	// Преобразуем входящие данные
	switch v := data.(type) {
	case []byte:
		payload = v
	case string:
		payload = []byte(v)
	default:
		payload, err = json.Marshal(v)
		if err != nil {
			logrus.Errorf("failed to serialize data: %v", err)
			return fmt.Errorf("failed to serialize data: %v", err)
		}
	}
	resp, err := c.centrifugo.Publish(ctx, &apiproto.PublishRequest{
		Channel: channel,
		Data:    payload,
	})
	if err != nil {
		logrus.Errorf("Transport level error: %v", err)
		return err
	}

	if resp.GetError() != nil {
		logrus.Errorf("centrifugo error: %s", resp.GetError().Message)

		return fmt.Errorf("centrifugo error: %s", resp.GetError().Message)
	}

	return nil

}
