package utils

import (
	"fmt"
	"time"
)

func Countdown() {
	// 获取当前时间
	now := time.Now()
	fmt.Println("现在时间", now)
	// 计算今天20点的时间
	target := time.Date(now.Year(), now.Month(), now.Day(), 21, 0, 0, 0, now.Location())
	fmt.Println("计算今天6点的时间", target)
	// 如果当前时间已经是下午6点之后，则计算明天6点
	if now.After(target) {
		target = target.AddDate(0, 0, 1)
	}

	// 计算倒计时
	duration := target.Sub(now) //time.Sub方法用于计算两个时间点之间的时间差
	fmt.Println("时间差是：", duration)

	// 使用一个无限循环进行倒计时
	for {
		//fmt.Printf("下班倒计时: %d小时%d分钟%d秒\r", duration/time.Hour, duration/time.Minute%60, duration/time.Second%60)
		Logger.Errorf("下班倒计时: %d小时%d分钟%d秒", duration/time.Hour, duration/time.Minute%60, duration/time.Second%60)
		time.Sleep(1 * time.Second)
		duration = duration - time.Second*1
		// 当倒计时结束时退出循环
		if duration <= 0 {
			break
		}
	}
	fmt.Println("下班倒计时结束！")
}
