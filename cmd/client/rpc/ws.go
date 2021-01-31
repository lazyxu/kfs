package rpc

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type Handle func(ctx context.Context, messageType int, data []byte, send func(messageType int, data []byte) error) error

func wsHandler(handle Handle) func(echo.Context) error {
	return func(c echo.Context) error {
		ctx, cancel := context.WithCancel(context.Background())
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer func() {
			cancel()
			ws.Close()
		}()

		for {
			// Read
			typ, data, err := ws.ReadMessage()
			if err != nil {
				return err
			}

			go func() {
				err := handle(ctx, typ, data, func(messageType int, data []byte) error {
					select {
					case <-ctx.Done():
					default:
						return ws.WriteMessage(messageType, data)
					}
					return nil
				})
				if err != nil {
					c.Logger().Error(err)
				}
			}()
		}
	}
}

func startWS(handle Handle) {
	e := echo.New()
	//e.Use(middleware.Logger())
	//e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.GET("/ws", wsHandler(handle))
	e.Logger.Fatal(e.Start(":1323"))
}
