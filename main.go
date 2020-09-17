/**
* @Author : henry
* @Data: 2020-09-09 16:22
* @Note: redis demo演示
**/

package main

import (
	"context"
	_ "encoding/binary"
	"fmt"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"time"
)

var ctx context.Context
var rdx *redis.Client

func init() {
	ctx = context.Background()

	rdx = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	result, err := rdx.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	// 清空redis数据
	rdx.FlushAll(ctx)
}

func main() {

	//// string
	//RdsString(ctx, rdx)

	//// hash
	//RdsHash(ctx, rdx)

	//// list
	//RdsList(ctx, rdx)

	//// set
	//RdsSet(ctx, rdx)

	// ZSet
	RdsZSet(ctx, rdx)
}

// set / get / strlen  / append / del
func RdsString(ctx context.Context, rdx *redis.Client) {
	key := "name"
	value := "henry"

	result, err := rdx.Set(ctx, key, value, -1).Result()
	if err != nil {
		fmt.Printf("set %s %s failed,error:%s\n", key, value, err)
		return
	}
	fmt.Printf("set %s %s success,result: %s\n", key, value, result)

	result, err = rdx.Get(ctx, key).Result()
	if err != nil {
		fmt.Printf("get key:%s failed,error:%s\n", key, err)
		return
	}
	fmt.Printf("get %s success,result: %s\n", key, result)

	i, err := rdx.StrLen(ctx, key).Result()
	if err != nil {
		fmt.Printf("strlen %s failed,error:%s\n", key, err)
		return
	}
	fmt.Printf("strlen %s success,result: %d\n", key, i)

	i2, err := rdx.Append(ctx, key, "W").Result()
	if err != nil {
		fmt.Printf("append %s %s failed,error:%s\n", key, "W", err)
		return
	}
	fmt.Printf("append %s %s success,result:%d\n", key, "W", i2)

	i3, err := rdx.Del(ctx, key).Result()
	if err != nil {
		fmt.Printf("del %s failed,error:%s\n", key, err)
		return
	}
	fmt.Printf("del %s success,result:%d\n", key, i3)

}

// hset / hexists / hget / hmset / hmget / hlen / hdel
func RdsHash(ctx context.Context, rdx *redis.Client) {
	key := "user:1"
	field := "name"
	value := "henry"
	result, err := rdx.HSet(ctx, key, field, value).Result()
	if err != nil {
		fmt.Printf("hset %s %s %s failed,err:%s\n", key, field, value, err)
		return
	}
	fmt.Printf("hset %s %s %s success,result:%d\n", key, field, value, result)

	b, err := rdx.HExists(ctx, key, field).Result()
	if err != nil {
		fmt.Printf("hexists %s %s failed,err:%s\n", key, field, err)
		return
	}
	fmt.Printf("hexists %s %s success,result:%t\n", key, field, b)

	if b {
		s, err := rdx.HGet(ctx, key, field).Result()
		if err != nil {
			fmt.Printf("hget %s failed,err:%s\n", key, err)
			return
		}
		fmt.Printf("hset %s success,result:%s\n", key, s)
	}

	key2 := "user:2"
	vf := make(map[string]interface{})
	vf = map[string]interface{}{
		"name":   "henry",
		"age":    18,
		"weight": 78,
		"height": 1.78,
	}
	b2, err := rdx.HMSet(ctx, key2, vf).Result()
	if err != nil {
		fmt.Printf("hmset %s %s failed,err:%s\n", key, vf, err)
		return
	}
	fmt.Printf("hmset %s %v success,result:%t\n", key, vf, b2)

	i, err := rdx.HLen(ctx, key2).Result()
	if err != nil {
		fmt.Printf("hlen %s failed,err:%s\n", key2, err)
		return
	}
	fmt.Printf("hlen %s success,result:%d\n", key2, i)

	if i > 1 {
		m, err := rdx.HGetAll(ctx, key2).Result()
		if err != nil {
			fmt.Printf("hgetall %s failed,err:%s\n", key2, err)
			return
		}
		fmt.Printf("hgetall %s success,result:%s\n", key2, m)
	}

	i2, err := rdx.HDel(ctx, key, field).Result()
	if err != nil {
		fmt.Printf("hdel %s failed,err:%s\n", key, err)
		return
	}
	fmt.Printf("hdel %s success,result:%d\n", key, i2)
}

// hpush / llen / lpop / lrange
func RdsList(ctx context.Context, rdx *redis.Client) {
	key := "char"
	value := []string{"a", "b", "c", "d", "e"}

	i, err := rdx.LPush(ctx, key, value).Result()
	if err != nil {
		fmt.Printf("lpush %s %s failed,err:%s\n", key, value, err)
		return
	}
	fmt.Printf("hset %s %s success,result:%d\n", key, value, i)

	i1, err := rdx.LLen(ctx, key).Result()
	if err != nil {
		fmt.Printf("llen %s failed,err:%s\n", key, err)
		return
	}
	fmt.Printf("llen %s success,result:%d\n", key, i1)

	result, err := rdx.LPop(ctx, key).Result()
	if err != nil {
		fmt.Printf("lpop %s failed,err:%s\n", key, err)
		return
	}
	fmt.Printf("lpop %s success,result:%s\n", key, result)

	var start int64 = 0
	var stop int64 = -1
	strings, err := rdx.LRange(ctx, key, start, stop).Result()
	if err != nil {
		fmt.Printf("lrange %s %d %d failed,err:%s\n", key, err)
		return
	}
	fmt.Printf("lpop %s success,result:%s\n", key, strings)

}

// sadd / scard / sismember / srem / sdiff / sinter / sunion
func RdsSet(ctx context.Context, rdx *redis.Client) {
	keyOne := "set1"
	valueOne := []string{"a", "b", "c", "d"}

	keyTwo := "set2"
	valueTwo := []string{"c", "d", "e", "f"}

	i, err := rdx.SAdd(ctx, keyOne, valueOne).Result()
	if err != nil {
		fmt.Printf("sadd %s %s failed,err:%s\n", keyOne, valueOne, err)
		return
	}
	fmt.Printf("sadd %s %s success,result:%d\n", keyOne, valueOne, i)

	i, err = rdx.SAdd(ctx, keyTwo, valueTwo).Result()
	if err != nil {
		fmt.Printf("sadd %s %s failed,err:%s\n", keyTwo, valueTwo, err)
		return
	}
	fmt.Printf("sadd %s %s success,result:%d\n", keyTwo, valueTwo, i)

	Set1Num, err := rdx.SCard(ctx, keyOne).Result()
	if err != nil {
		fmt.Printf("scard %s failed,err:%s\n", keyOne, err)
		return
	}
	fmt.Printf("scard %s success,result:%d\n", keyOne, Set1Num)

	if i > 0 {
		charDel := "a"
		isExists, err := rdx.SIsMember(ctx, keyOne, "a").Result()
		if err != nil {
			fmt.Printf("sismember %s %s failed,err:%s\n", keyOne, charDel, err)
			return
		}
		fmt.Printf("sismember %s %s success,result:%t\n", keyOne, charDel, isExists)

		if isExists {
			i, err := rdx.SRem(ctx, keyOne, charDel).Result()
			if err != nil {
				fmt.Printf("srem %s %s failed,err:%s\n", keyOne, charDel, err)
				return
			}
			fmt.Printf("srem %s %s success,result:%d\n", keyOne, charDel, i)
		}
	}

	result, err := rdx.SDiff(ctx, keyOne, keyTwo).Result()
	if err != nil {
		fmt.Printf("sdiff %s %s failed,err:%s\n", keyOne, keyTwo, err)
		return
	}
	fmt.Printf("sdiff %s %s success,result:%s\n", keyOne, keyTwo, result)

	result, err = rdx.SInter(ctx, keyOne, keyTwo).Result()
	if err != nil {
		fmt.Printf("sinter %s %s failed,err:%s\n", keyOne, keyTwo, err)
		return
	}
	fmt.Printf("sinter %s %s success,result:%s\n", keyOne, keyTwo, result)

	result, err = rdx.SUnion(ctx, keyOne, keyTwo).Result()
	if err != nil {
		fmt.Printf("sunion %s %s failed,err:%s\n", keyOne, keyTwo, err)
		return
	}
	fmt.Printf("sunion %s %s success,result:%s\n", keyOne, keyTwo, result)
}

func RdsZSet(ctx context.Context, rdx *redis.Client) {
	rand.Seed(time.Now().UnixNano())
	zSet1 := "zSet1"
	zSet2 := "zSet2"

	z1s := make([]*redis.Z, 0)
	z2s := make([]*redis.Z, 0)

	for i := 0; i < 5; i++ {
		z1s = append(z1s, &redis.Z{Score: rand.Float64(), Member: fmt.Sprintf("one%s", rand.Float64())})
		z2s = append(z2s, &redis.Z{Score: rand.Float64(), Member: fmt.Sprintf("one%s", rand.Float64())})
	}

	i, err := rdx.ZAdd(ctx, zSet1, z1s...).Result()
	if err != nil {
		fmt.Printf("sadd %s failed,err:%s\n", zSet1, err)
		return
	}
	fmt.Printf("sadd %s success,result:%d\n", zSet1, i)

	i, err = rdx.ZAdd(ctx, zSet2, z2s...).Result()
	if err != nil {
		fmt.Printf("sadd %s failed,err:%s\n", zSet2, err)
		return
	}
	fmt.Printf("sadd %s success,result:%d\n", zSet2, i)

	i, err = rdx.ZCard(ctx, zSet1).Result()
	if err != nil {
		fmt.Printf("zcard %s failed,err:%s\n", zSet1, err)
		return
	}
	fmt.Printf("zcard %s success,result:%d\n", zSet1, i)

	var start int64 = 0
	var stop int64 = -1
	str, err := rdx.ZRange(ctx, zSet1, start, stop).Result()
	if err != nil {
		fmt.Printf("zrange %s failed,err:%s\n", zSet1, err)
		return
	}
	fmt.Printf("zrange %s success,result:", zSet1)
	for i := 0; i < len(str); i++ {
		fmt.Printf("%v", str[i])
	}
}
