package util

import (
	"fmt"
	"time"
)

// RetryTimes 重试，限制次数
func RetryTimes(name string, tryTimes int, sleep time.Duration, callback func() error) (err error) {
	for i := 1; i <= tryTimes; i++ {
		err = callback()
		if err == nil {
			return nil
		}
		fmt.Printf("[%v]失败，第%v次重试， 错误信息:%s \n", name, i, err)
		time.Sleep(sleep)
	}
	err = fmt.Errorf("[%v]失败，共重试%d次, 最近一次错误:%s \n", name, tryTimes, err)
	fmt.Println(err)
	return err
}
