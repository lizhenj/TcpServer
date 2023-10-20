package ztimer

import (
	"fmt"
	"testing"
	"time"
)

//定义一个超时函数
func myFunc(v ...interface{}) {
	fmt.Printf("No.(%d) function calld. delay %d second(s)\n", v[0].(int), v[1].(int))
}

func TestTimer(t *testing.T) {

	for i := 0; i < 5; i++ {
		go func(i int) {
			NewTimerAfter(NewDelayFunc(myFunc, []interface{}{1, 2 * i}), time.Duration(2*i)*time.Second).Run()
		}(i)
	}

	//阻塞
	time.Sleep(1 * time.Minute)
}
