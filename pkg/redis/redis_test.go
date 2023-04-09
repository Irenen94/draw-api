package redis

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"testing"
)

func TestConnect(t *testing.T) {
	//维护服务执行状态的全局变量
	conn, err := Connect()
	if err != nil {
		return
	}
	defer conn.Close()

	r, err := redis.String(conn.Do("Get", "user_from_token_beb35604cf2aca25a403b7fbf6aa6c1e"))
	if err != nil {
		fmt.Println("conn.Do err=", err)
		return
	}

	fmt.Println("取出的数据 r=", r)

	type User struct {
		Name string `json:"user"`
		Age  int    `json:"age"`
	}
	user := User{
		Name: "test-name",
		Age:  18,
	}

	var ub []byte
	ub, err = json.Marshal(user)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	conn.Send("SET", "user", ub)
	v, err := conn.Receive()
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	fmt.Printf("%#v\n", v)

}

//func TestConnectPool(t *testing.T) {
//
//	RedisInit()
//	defer RedisClose()
//
//	var err error
//	_, err = RedisConn.Do("SETEX", "name22", 60*60*24, "abc") //redis set命令，超时时间是24小时
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	do, err := RedisConn.Do("get", "name22")
//	if err != nil || do == nil {
//		return
//	}
//
//	res, err := redis.String(do, err) //redis get命令
//
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println("res:", res)
//
//}
