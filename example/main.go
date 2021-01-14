package main

import (
	"fmt"
	"log"
	"time"

	"github.com/18211167516/gocache"
	_ "github.com/18211167516/gocache/store"
)

func main() {
	cache, err := gocache.New("Memory")
	if err != nil {
		log.Fatalf("初始化缓存管理器 失败 %s", err)
	}

	fmt.Println("获取store type:", cache.GetStoreName())

	fmt.Println(time.Now())

	if err := cache.Set("aaaa", "123123", 1); err != nil {
		log.Fatalf("设置键 %s 的value 失败 %s", "aaaa", err)
	}

	fmt.Println("判断键值：", fmt.Sprintln(cache.Has("aaaa")))
	fmt.Println("获取键值：", fmt.Sprintln(cache.Get("aaaa")))
	fmt.Println("获取ttl：", fmt.Sprintln(cache.GetTTl("aaaa")))
	/* fmt.Println("设置永久键值：", cache.Forever("bbb", "bbbb"))
	fmt.Println("获取ttl：", fmt.Sprintln(cache.GetTTl("bbb")))
	fmt.Println("获取键值：", fmt.Sprintln(cache.Get("bbb"))) */

	_ = cache.Set("cccc", "cccc", 5)
	_ = cache.Set("dddd", "dddd", 10)
	fmt.Println("获取键值：", fmt.Sprintln(cache.Get("cccc")))
	time.Sleep(time.Duration(5) * time.Second)

	fmt.Println("延迟5秒获取键值：", fmt.Sprintln(cache.Get("cccc")))

	fmt.Println("获取键值：", fmt.Sprintln(cache.Get("dddd")))
	fmt.Println("获取ttl：", fmt.Sprintln(cache.GetTTl("dddd")))
	//fmt.Println("删除全部：", fmt.Sprintln(cache.Clear()))
	fmt.Println("缓存个数：", cache.Size())
	//select {}

	ticker := time.NewTicker(5 * time.Second)

exit:
	for {
		select {
		case <-ticker.C:
			size := cache.Size()
			fmt.Println("缓存个数：", size)
			if size == 0 {
				ticker.Stop()
				break exit
			}
		}
	}

}
