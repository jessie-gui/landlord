package core

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/jessie-gui/x/xlog"
	"github/jessie-gui/landlord/consts"
)

// Table 桌子对象。
type Table struct {
	Cards       []*Card         // 所有牌
	Players     map[int]*player // 玩家列表
	BottomCards []*Card         // 底牌
	TempCard    TempCard        // 最近一个出牌详情
	sync.RWMutex
}

// TempCard 最近一个出牌详情。
type TempCard struct {
	Position     int     // 上一个出牌玩家位置
	NextPosition int     // 下一个出牌玩家位置
	Cards        []*Card // 上一个出牌详情
	CardsIndex   []int   // 上一个出牌索引
}

// NewTable 新建桌子。
func NewTable() *Table {
	cards := make([]*Card, 0, 54)
	suits := []string{"♠", "♥", "♦", "♣"}
	ranks := []int32{3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	for _, suit := range suits {
		for _, rank := range ranks {
			var name string
			switch rank {
			case 3, 4, 5, 6, 7, 8, 9, 10:
				name = strconv.Itoa(int(rank))
			case 11:
				name = "J"
			case 12:
				name = "Q"
			case 13:
				name = "K"
			case 14:
				name = "A"
			case 15:
				name = "2"

			}
			card := NewCard(suit, rank, name)
			cards = append(cards, card)
		}
	}

	cards = append(cards, NewCard("小王", 16, "小王"))
	cards = append(cards, NewCard("大王", 17, "大王"))

	return &Table{Cards: cards, Players: make(map[int]*player)}
}

// AddPlayer 加入玩家。
func (t *Table) AddPlayer(p *player) *Table {
	if len(t.Players) > 3 {
		xlog.Error("房间已满！")

		return nil
	}

	if len(t.Players) <= 0 {
		t.Players = make(map[int]*player)
	}

	ps := make(map[int]*player)
	for _, v := range t.Players {
		ps[v.Position] = v
	}

	for i := 1; i <= 3; i++ {
		if _, ok := ps[i]; !ok {
			p.Position = i
			t.Players[p.Id] = p

			break
		}
	}

	xlog.Infof("玩家%s进入房间", p.Name)

	t.BroadCastMsg(p, consts.MSG_TYPE_OF_JOIN_TABLE, "玩家加入游戏")

	return t
}

// Deal 发牌。
func (t *Table) Deal() *Table {
	rand.Seed(time.Now().UnixNano())

	var ids []int
	for id, _ := range t.Players {
		ids = append(ids, id)
	}

	for i, card := range t.Cards {
		if i < 51 {
			switch i % 3 {
			case 0:
				t.Players[ids[0]].Card = append(t.Players[ids[0]].Card, card)
			case 1:
				t.Players[ids[1]].Card = append(t.Players[ids[1]].Card, card)
			case 2:
				t.Players[ids[2]].Card = append(t.Players[ids[2]].Card, card)
			}
		} else {
			t.BottomCards = append(t.BottomCards, card)
		}
	}

	xlog.Info("发牌完毕！")

	t.RandLord()

	for _, p := range t.Players {
		p.SendCard()

		if p.IsLord {
			p.SendMsgToPlayer(t, consts.MSG_TYPE_OF_PLAY_CARD, "玩家准备出牌")
		}
	}

	return t
}

// RandLord 随机叫地主。
func (t *Table) RandLord() *Table {
	rand.Seed(time.Now().UnixNano())

	var ids []int
	for id, _ := range t.Players {
		ids = append(ids, id)
	}

	k := rand.Intn(len(ids))

	t.Players[ids[k]].IsLord = true

	t.Players[ids[k]].Card = append(t.Players[ids[k]].Card, t.BottomCards...)

	t.TempCard.Position = t.Players[ids[k]].Position
	t.TempCard.NextPosition = t.Players[ids[k]].Position

	xlog.Infof("地主是%s！", t.Players[ids[k]].Name)

	return t
}

// BroadCastMsg 推送广播消息。
func (t *Table) BroadCastMsg(player *player, msgType int, hints string) {
	broadcastMsg := NewMsg()
	broadcastMsg.SubMsgType = msgType

	t.RLock()
	defer t.RUnlock()

	if player != nil {
		broadcastMsg.PlayerId = player.Id
		for i, p := range t.Players {
			if p != nil {
				broadcastMsg.PlayerIndexIdDic["id"+strconv.Itoa(p.Id)] = i
			}
		}
	}

	switch msgType {
	case consts.MSG_TYPE_OF_TIME_TICKER:
		broadcastMsg.Msg = hints
	case consts.MSG_TYPE_OF_READY:
		broadcastMsg.Msg = strconv.Itoa(player.Id) + "已准备"
	case consts.MSG_TYPE_OF_UN_READY:
		broadcastMsg.Msg = strconv.Itoa(player.Id) + "取消准备"
	case consts.MSG_TYPE_OF_JOIN_TABLE:
		broadcastMsg.Msg = strconv.Itoa(player.Id) + "加入游戏"
	case consts.MSG_TYPE_OF_LEAVE_TABLE:
		broadcastMsg.Msg = strconv.Itoa(player.Id) + "离开游戏"
	case consts.MSG_TYPE_OF_PLAY_CARD:
		broadcastMsg.Msg = strconv.Itoa(player.Id) + "出牌"
		for _, card := range t.TempCard.Cards {
			broadcastMsg.Cards = append(broadcastMsg.Cards, card)
		}

		broadcastMsg.CardsIndex = t.TempCard.CardsIndex
	case consts.MSG_TYPE_OF_PASS:
		broadcastMsg.Msg = strconv.Itoa(player.Id) + "过牌"
	case consts.MSG_TYPE_OF_CALL_SCORE:
		broadcastMsg.Msg = strconv.Itoa(player.Id) + "叫地主"
		broadcastMsg.Score = 0
	case consts.MSG_TYPE_OF_SCORE_CHANGE:
		broadcastMsg.Msg = "基础变动"
		broadcastMsg.Score = 0
	case consts.MSG_TYPE_OF_SEND_BOTTOM_CARDS:
		broadcastMsg.Msg = "发放底牌"
	case consts.MSG_TYPE_OF_GAME_OVER:
		broadcastMsg.Msg = "游戏结束，结算积分"
	default:
		broadcastMsg.Msg = hints
	}

	msgJson, err := json.Marshal(broadcastMsg)
	if err != nil {
		panic(err.Error())
	}

	for _, p := range t.Players {
		if p != nil {
			p.SendMsg(msgJson)
		}
	}
}

// PlayCard 出牌。
func (t *Table) PlayCard(position int, cards []*Card, cardsIndex []int) {
	if len(t.Players[position].Card) <= 0 {
		xlog.Error("游戏已结束！")
		return
	}

	lastCards := t.TempCard

	xlog.Info("当前出牌位置:", lastCards.NextPosition, position)

	// 没有轮到玩家出牌。
	if lastCards.NextPosition != position {
		xlog.Error("还没轮到你出牌！")
		return
	}

	// 检查所出牌是否都拥有。
	cs := t.Players[position].Card
	hitNum := 0
	for _, c := range cards {
		for _, card := range cs {
			if card.Rank == c.Rank && card.Suit == c.Suit {
				hitNum++
			}
		}
	}

	if len(cards) != hitNum {
		xlog.Error("跟牌错误，有未拥有的牌！")
		return
	}

	// 跟牌数量不匹配。
	if len(lastCards.Cards) != 0 && len(lastCards.Cards) != len(cards) && len(cards) != 2 && len(cards) != 4 && lastCards.Position != position {
		xlog.Error("跟牌数量不匹配:%d-%d", len(lastCards.Cards), len(cards))
		return
	}

	// 下一个出牌玩家位置。
	nextPosition := position + 1
	if nextPosition > 3 {
		nextPosition -= 3
	}

	// 构建出牌详情对象。
	tempCard := TempCard{Position: position, NextPosition: nextPosition, Cards: cards, CardsIndex: cardsIndex}

	switch len(cards) {
	case 1: // 单张
		if len(lastCards.Cards) > 0 && lastCards.Position != position && lastCards.Cards[0].Rank >= cards[0].Rank {
			xlog.Error("跟牌错误:", lastCards.Cards[0].Rank, cards[0].Rank)
			return
		}
	case 2: // 对子或王炸。
		if cards[0].Rank != cards[1].Rank && cards[0].Rank != 16 && cards[1].Rank != 17 && cards[0].Rank != 17 && cards[1].Rank != 16 {
			xlog.Error("跟牌错误:", cards[0].Rank, cards[1].Rank)
			return
		}

		if len(lastCards.Cards) > 0 && lastCards.Position != position {
			if lastCards.Cards[0].Rank == 16 && lastCards.Cards[1].Rank == 17 {
				xlog.Error("跟牌错误，上家王炸！")
				return
			}

			if lastCards.Cards[0].Rank >= cards[0].Rank {
				xlog.Error("跟牌错误:", lastCards.Cards[0].Rank, lastCards.Cards[1].Rank, cards[0].Rank, cards[1].Rank)
				return
			}
		}
	case 3: // 三张。
		if cards[0].Rank != cards[1].Rank || cards[1].Rank != cards[2].Rank {
			xlog.Error("跟牌错误:", cards[0].Rank, cards[1].Rank, cards[2].Rank)
			return
		}

		if len(lastCards.Cards) > 0 && lastCards.Cards[0].Rank >= cards[0].Rank && lastCards.Position != position {
			xlog.Error("跟牌错误:", lastCards.Cards[0].Rank, cards[0].Rank)
			return
		}
	case 4: // 炸弹。
		if cards[0].Rank != cards[1].Rank || cards[1].Rank != cards[2].Rank || cards[2].Rank != cards[3].Rank {
			xlog.Error("跟牌错误:", cards[0].Rank, cards[1].Rank, cards[2].Rank, cards[3].Rank)
			return
		}

		if IsBoom(lastCards.Cards) && lastCards.Cards[0].Rank >= cards[0].Rank && lastCards.Position != position {
			xlog.Error("跟牌错误:", lastCards.Cards[0].Rank, cards[0].Rank)
			return
		}
	default:
		if len(lastCards.Cards) != len(cards) && lastCards.Position != position {
			xlog.Error("跟牌数量错误:", len(lastCards.Cards), len(cards))
			return
		}

		if IsStraight(cards) || IsStraightPair(cards) {
			// 上一个出牌的不是自己。
			if lastCards.Position != position {
				if GetMaxCardRank(lastCards.Cards) >= GetMaxCardRank(cards) {
					xlog.Error("跟牌错误:", GetMaxCardRank(cards), GetMaxCardRank(cards))
					return
				}
			}
		} else {
			xlog.Error("跟牌错误!不是顺子也不是连对！")
			return
		}
	}

	t.TempCard = tempCard
	for _, i := range cards {
		for k, v := range t.Players[position].Card {
			if i.Suit == v.Suit && i.Rank == v.Rank {
				t.Players[position].Card = append(t.Players[position].Card[:k], t.Players[position].Card[k+1:]...)
			}
		}
	}

	t.BroadCastMsg(t.Players[position], consts.MSG_TYPE_OF_PLAY_CARD, "玩家出牌")

	// 下一个玩家准备出牌。
	if len(t.Players[position].Card) > 0 {
		t.Players[nextPosition].SendMsgToPlayer(t, consts.MSG_TYPE_OF_PLAY_CARD, "玩家准备出牌")
	} else {
		msg := "游戏结束,地主胜利"
		if t.Players[position].IsLord == true {
			msg = "游戏结束,农民胜利"
		}
		t.BroadCastMsg(nil, consts.MSG_TYPE_OF_GAME_OVER, msg)
	}
}
