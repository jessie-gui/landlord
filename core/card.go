package core

// Card 纸牌
type Card struct {
	Suit  string // 花色
	Rank  int32  // 点数
	Name  string // 名称
	Index int    // 在玩家手中的索引
}

func NewCard(suit string, rank int32, name string) *Card {
	return &Card{Suit: suit, Rank: rank, Name: name}
}
