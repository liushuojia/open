package utils

import (
	"fmt"
	"math/big"
	"time"
)

/*
1. 生成方案
	前端js数字最大值 9007199254740991
		console.log(Number.MAX_SAFE_INTEGER)

		9007199254740991  最大值
		2025100801890001  方案值

	id 生成规则
		2025 		10 			08			01				01				0001
		年			月(01-12)	日(01-31)	00-23 (小时)		服务器(00-99)	自增
		max(8999)

	位	16-13 		12-11		10-9		8-7				6-5				4-1


2. 优劣
	如果出现时间回拨，可能存在相同的ID
	一小时一个服务器 最多生成10000记录
*/

const (
	yearStep  = 12
	monthStep = 10
	dayStep   = 8
	hourStep  = 6

	serverIDStep = 4
	autoIDStep   = 0

	//serverIDStep = 0
	//autoIDStep   = 2

	serverIDStepLength = 4
)

type ID struct {
	t        time.Time
	serverID uint     // 服务器 01-99
	idList   []uint64 // id 数组
}

func NewID(t time.Time, serverID uint) (*ID, error) {
	if serverID >= 100 {
		return nil, fmt.Errorf("serverID out of range")
	}
	id := &ID{
		t:        t,
		serverID: serverID,
		idList:   make([]uint64, 0),
	}
	base := uint64(t.Year())*id.calculate10ToN(yearStep) +
		uint64(t.Month())*id.calculate10ToN(monthStep) +
		uint64(t.Day())*id.calculate10ToN(dayStep) +
		uint64(t.Hour())*id.calculate10ToN(hourStep) +
		uint64(id.serverID)*id.calculate10ToN(serverIDStep)
	//uint64(0)*id.calculate10ToN(autoIDStep)

	for i := uint64(0); i < id.calculate10ToN(serverIDStepLength); i++ {
		n := base + i*id.calculate10ToN(autoIDStep)
		id.idList = append(id.idList, n)
	}
	return id, nil
}

func (u *ID) calculate10ToN(n uint) uint64 {
	// 初始化结果（10^0 = 1）
	result := big.NewInt(1)
	base := big.NewInt(10) // 底数10

	// 分n幕计算：每幕乘以10
	for i := 0; i < int(n); i++ {
		result.Mul(result, base) // 当前结果 = 结果 × 10
	}
	return result.Uint64()
}
func (u *ID) MaxID() uint64 {
	return u.idList[len(u.idList)-1]
}
func (u *ID) MinID() uint64 {
	return u.idList[0]
}
func (u *ID) IDList() []uint64 {
	return u.idList
}
