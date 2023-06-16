package core

import (
	"sort"
)

// IsBoom 判断是否为炸弹。
func IsBoom(c []*Card) bool {
	if len(c) == 4 && c[0].Rank == c[1].Rank && c[1].Rank == c[2].Rank && c[2].Rank == c[3].Rank {
		return true
	}

	return false
}

// GetMaxCardRank 获取最大牌的点数。
func GetMaxCardRank(c []*Card) int {
	if len(c) <= 0 {
		return 0
	}

	maxRank := 0
	for _, v := range c {
		if int(v.Rank) > maxRank {
			maxRank = int(v.Rank)
		}
	}

	return maxRank
}

// IsStraightPair 判断是否为连对。
func IsStraightPair(c []*Card) bool {
	// 连对至少有3对，且牌数必须为偶数。
	if len(c)%2 != 0 || len(c)/2 < 3 {
		return false
	}

	// 将牌按点数从小到大排序。
	sort.Slice(c, func(i, j int) bool {
		return c[i].Rank < c[j].Rank
	})

	// 如果相邻两对牌中有一对点数不同，或者两对牌的点数之差不为1，则不是连对。
	for i := 0; i < len(c); i += 2 {
		if c[i].Rank == 16 || c[i].Rank == 17 {
			return false
		}

		if c[i].Rank != c[i+1].Rank || c[i+1].Rank-c[i].Rank != 1 {
			return false
		}
	}

	return true
}

// IsStraight 判断是否为顺子。
func IsStraight(c []*Card) bool {
	// 顺子至少有5张牌，最多不超过A23456789TJQK
	if len(c) < 5 || len(c) > 12 {
		return false
	}

	// 将牌按点数从小到大排序。
	sort.Slice(c, func(i, j int) bool {
		return c[i].Rank < c[j].Rank
	})

	// 如果相邻两张牌点数之差不为1，则不是顺子。
	for i := 0; i < len(c)-1; i++ {
		if c[i].Rank == 16 || c[i].Rank == 17 {
			return false
		}

		if c[i+1].Rank-c[i].Rank != 1 {
			return false
		}
	}

	return true
}
