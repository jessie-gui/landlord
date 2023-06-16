package controller

import (
	"github.com/jessie-gui/x/xlog"
	"github.com/jessie-gui/x/xserver/xwebsocket"
	"github.com/labstack/echo/v4"
	"github/jessie-gui/landlord/consts"
	"github/jessie-gui/landlord/core"
	"github/jessie-gui/landlord/global"
	"net/http"
	"strconv"
	"sync"
)

var (
	userId int
	wg     sync.WaitGroup
)

// Index 游戏主页面。
func Index(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, "index.html", map[string]int{
		"ok": 1,
	})
}

// Connect 连接websocket。
func Connect(ctx echo.Context) error {
	wsClient, err := xwebsocket.NewWebSocketClient(ctx.Response(), ctx.Request())
	if err != nil {
		xlog.Error("NewWebSocketClient:", err)
		return nil
	}

	defer wsClient.Close()

	userId++
	xlog.Info("玩家：" + strconv.Itoa(userId) + "登陆游戏")

	player := core.NewPlayer(core.Id(userId), core.Name(strconv.Itoa(userId)), core.WsClient(wsClient))

	player.SendMsgToPlayer(global.Table, consts.MSG_TYPE_OF_LOGIN, "登录游戏")

	if len(global.Table.Players) < 3 {
		global.Table.AddPlayer(player)

		wg.Add(1)

		//启动一个goroutine监听该客户端发来的消息。
		go player.HandlerMsg(&wg, global.Table)

		wg.Wait()
	}

	return nil
}
