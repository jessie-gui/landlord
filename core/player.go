package core

import (
	"encoding/json"
	"fmt"
	"github.com/camry/g/glog"
	"github.com/gorilla/websocket"
	"github.com/jessie-gui/x/xlog"
	"github.com/jessie-gui/x/xserver/xwebsocket"
	"github.com/tidwall/gjson"
	"github/jessie-gui/landlord/consts"
	"log"
	"runtime/debug"
	"sort"
	"strconv"
	"sync"
)

// PlayerOption /**
type PlayerOption func(p *player)

type player struct {
	Id       int                         // 玩家编号
	Name     string                      // 玩家名字
	Position int                         // 玩家位置
	Card     []*Card                     // 玩家手中的牌
	IsLord   bool                        // 是否是地主
	IsReady  bool                        // 是否准备好
	wsClient *xwebsocket.WebSocketClient // ws连接对象
	Score    int                         // 叫地主分
}

func Id(id int) PlayerOption {
	return func(p *player) {
		p.Id = id
	}
}

func Name(name string) PlayerOption {
	return func(p *player) {
		p.Name = name
	}
}

func WsClient(wsClient *xwebsocket.WebSocketClient) PlayerOption {
	return func(p *player) {
		p.wsClient = wsClient
	}
}

func NewPlayer(o ...PlayerOption) *player {
	p := &player{}
	for _, opt := range o {
		opt(p)
	}

	return p
}

func (p *player) SendMsg(msg []byte) {
	if err := p.wsClient.WriteMessage(msg); err != nil {
		xlog.Errorf("玩家%s发送消息失败：%v", p.Name, err)
	}
}

func (p *player) Ready(t *Table) {
	p.IsReady = true

	readyNum := 0
	for _, v := range t.Players {
		if v.IsReady {
			readyNum++
		}
	}

	if readyNum == 3 {
		t.Deal()
	}

	t.BroadCastMsg(p, consts.MSG_TYPE_OF_READY, "玩家准备")
}

func (p *player) UnReady(t *Table) {
	p.IsReady = false

	t.BroadCastMsg(p, consts.MSG_TYPE_OF_UN_READY, "玩家准备")
}

func (p *player) Pass(t *Table) {
	// 更新下一个出牌玩家位置
	NextPosition := t.TempCard.NextPosition + 1
	if NextPosition > 3 {
		NextPosition -= 3
	}

	xlog.Info("当前过牌位置:", t.TempCard.NextPosition)
	xlog.Info("下一个出牌位置:", NextPosition)

	t.TempCard.NextPosition = NextPosition

	t.BroadCastMsg(p, consts.MSG_TYPE_OF_PASS, "玩家过牌")

	for _, v := range t.Players {
		if v.Position == NextPosition {
			v.SendMsgToPlayer(t, consts.MSG_TYPE_OF_PLAY_CARD, "玩家准备出牌")
		}
	}
}

func (p *player) CallScore(score int) {
	p.Score = score
}

// SendCard 推送发牌消息
func (p *player) SendCard() {
	newMsg := NewMsg()
	newMsg.MsgType = consts.MSG_TYPE_OF_SEND_CARD

	cards := p.Card
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Rank < cards[j].Rank
	})

	for k, card := range cards {
		card.Index = k
		newMsg.CardsIndex = append(newMsg.CardsIndex, k)
	}

	newMsg.Cards = cards
	newMsg.PlayerId = p.Id

	msgJson, err := json.Marshal(newMsg)
	if err != nil {
		panic(err.Error())
	}

	p.SendMsg(msgJson)
}

func (p *player) ResolveMsg(msgB []byte, t *Table) error {
	msgType, err := strconv.Atoi(gjson.Get(string(msgB), "MsgType").String())
	if err != nil {
		p.SendMsg(msgB)
		return err
	}

	switch msgType {
	case consts.MSG_TYPE_OF_AUTO:

	case consts.MSG_TYPE_OF_UN_READY:
		go p.UnReady(t)
	case consts.MSG_TYPE_OF_READY:
		go p.Ready(t)
	case consts.MSG_TYPE_OF_PLAY_CARD:
		cardIndex := gjson.Get(string(msgB), "Data.CardIndex").Array()
		var cards []*Card
		var cardsIndex []int
		for _, cardIndex := range cardIndex {
			for _, v := range p.Card {
				if v.Index == int(cardIndex.Int()) {
					cards = append(cards, v)
				}
			}

			cardsIndex = append(cardsIndex, int(cardIndex.Int()))
		}

		go t.PlayCard(p.Position, cards, cardsIndex)
	case consts.MSG_TYPE_OF_PASS:
		go p.Pass(t)
	case consts.MSG_TYPE_OF_LEAVE_TABLE:

	case consts.MSG_TYPE_OF_JOIN_TABLE:

	case consts.MSG_TYPE_OF_HINT:

	case consts.MSG_TYPE_OF_CALL_SCORE:
		score, _ := strconv.Atoi(gjson.Get(string(msgB), "Data.Score").String())
		go p.CallScore(score)

	default:
		p.wsClient.WriteMessageWithType(msgType, msgB)
	}

	return nil
}

func (p *player) HandlerMsg(wg *sync.WaitGroup, room *Table) {
	defer wg.Done()
	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("panic recover! p: %v", e)
			debug.PrintStack()
		}
	}()

	for {
		msgType, msg, err := p.wsClient.ReadMessageWithType()
		log.Println("消息类型:", msgType)
		if err == nil {
			switch msgType {
			case websocket.TextMessage:
				//同桌用户交流，包含对话流程和出牌流程
				p.ResolveMsg(msg, room)
			case -1:
				glog.Info("玩家：" + strconv.Itoa(p.Id) + "断开链接")
				break
				//离开桌子流程，后续包含断线保持，自动出牌
			default:

			}
		} else {
			log.Println("获取客户端消息错误:", err)
			break
		}
	}
}

func (p *player) SendMsgToPlayer(t *Table, msgType int, hints string) {

	var newMsg []byte
	var err error
	switch msgType {
	case consts.MSG_TYPE_OF_CALL_SCORE:

	case consts.MSG_TYPE_OF_CALL_SCORE_TIME_OUT:

	case consts.MSG_TYPE_OF_PLAY_CARD:
		newMsg, err = NewPlayCardMsg()
	case consts.MSG_TYPE_OF_PLAY_ERROR:

	case consts.MSG_TYPE_OF_PLAY_CARD_SUCCESS:

	case consts.MSG_TYPE_OF_LOGIN:
		newMsg, err = NewLoginMsg(p.Id, "登陆成功")
	default:
		return
	}

	if err == nil {
		p.SendMsg(newMsg)
	} else {
		log.Fatal(err.Error())
	}
}
