package core

import (
	"encoding/json"
	"github/jessie-gui/landlord/consts"
)

// Msg 通用消息对象
type Msg struct {
	MsgType          int
	SubMsgType       int
	Msg              string
	Cards            []*Card
	CardsIndex       []int
	Score            int
	PlayerId         int
	SettleInfoDic    map[string]string
	PlayerIndexIdDic map[string]int
}

func NewMsg() *Msg {
	return &Msg{
		consts.MSG_TYPE_OF_TABLE_BROADCAST,
		-1,
		"",
		nil,
		[]int{},
		-1,
		-1,
		make(map[string]string),
		make(map[string]int),
	}
}

// LoginMsg 登录消息对象
type LoginMsg struct {
	MsgType int
	Msg     string
	ID      int
}

func NewLoginMsg(userID int, loginMsg string) ([]byte, error) {
	newMsg := LoginMsg{
		consts.MSG_TYPE_OF_LOGIN,
		loginMsg,
		-1,
	}
	newMsg.ID = userID
	return json.Marshal(newMsg)
}

type PlayMsg struct {
	MsgType int
	Msg     string
}

func NewPlayCardMsg() ([]byte, error) {
	msg := PlayMsg{
		consts.MSG_TYPE_OF_PLAY_CARD,
		"",
	}
	return json.Marshal(msg)
}
