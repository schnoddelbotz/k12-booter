package ws

import (
	"context"
	"time"

	"nhooyr.io/websocket"
)

func ConnectToTeacherWS(url string) (*websocket.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	c, _, err := websocket.Dial(ctx, url+"/subscribe", nil)
	if err != nil {
		return nil, err
	}
	return c, nil
}
