package consts

const (
	MSG_TYPE_OF_READY       = iota //准备
	MSG_TYPE_OF_UN_READY           //取消准备
	MSG_TYPE_OF_JOIN_TABLE         //加入桌子
	MSG_TYPE_OF_LEAVE_TABLE        //离开桌子

	MSG_TYPE_OF_HINT      //提示
	MSG_TYPE_OF_PLAY_CARD //出牌
	MSG_TYPE_OF_PASS      //过牌

	MSG_TYPE_OF_AUTO                //托管
	MSG_TYPE_OF_SEND_CARD           //发牌
	MSG_TYPE_OF_CALL_SCORE          //抢地主叫分
	MSG_TYPE_OF_CONFIRM             //客户端出牌等操作确认信息
	MSG_TYPE_OF_CALL_SCORE_TIME_OUT //叫地主超时
	MSG_TYPE_OF_PLAY_ERROR          //出牌错误
	MSG_TYPE_OF_PLAY_CARD_SUCCESS   //出牌成功
	MSG_TYPE_OF_TABLE_BROADCAST     //桌子广播消息
	MSG_TYPE_OF_SCORE_CHANGE        //牌局分数变化
	MSG_TYPE_OF_SETTLE_SCORE        //结算玩家分数
	MSG_TYPE_OF_GAME_OVER           //游戏结束
	MSG_TYPE_OF_LOGIN               //登陆消息
	MSG_TYPE_OF_SEND_BOTTOM_CARDS   //发底牌
	MSG_TYPE_OF_TIME_TICKER         //倒计时数
	MSG_TYPE_OF_POKER_RECORDER      //推送记牌器消息
)
