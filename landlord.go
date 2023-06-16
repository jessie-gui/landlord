package main

import (
	"context"
	"github/jessie-gui/landlord/core"
	"github/jessie-gui/landlord/global"

	"github.com/jessie-gui/x/xlog"
	"github.com/jessie-gui/x/xserver/xhttp"
	"github/jessie-gui/landlord/router"
)

/**
 *
 *
 * @author        Gavin Gui <guijiaxian@gmail.com>
 * @version       1.0.0
 * @copyright (c) 2022, Gavin Gui
 */
func main() {
	server := xhttp.NewServer(
		xhttp.Address(":8080"),
		xhttp.Handler(router.NewEcho()),
	)

	global.Table = core.NewTable()

	if err := server.Start(context.Background()); err != nil {
		xlog.Fatal("服务启动失败:", err)
	}
}
