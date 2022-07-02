package week08

import (
	"context"
	"github.com/go-redis/redis/v8"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func IntRun(alsiz int) {
	sizes := []int{
		10000,
		25000,
		75000,
		100000,
		250000,
		350000,
		500000,
	}
	lsiz := 10
	if alsiz > 0 {
		lsiz = alsiz
	}

	graceTimeOut := 300000

	var str strings.Builder
	for i := 0; i < lsiz; i++ {
		str.WriteString("a")
	}

	timeOut := time.Duration(graceTimeOut) * time.Second
	ctx, cancle := context.WithTimeout(context.Background(), timeOut)
	client, err := InitRedisClient(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer cancle()
	defer client.Close()
	// flush db before insert
	client.FlushDB(ctx)

	var key string
	for _, size := range sizes {
		// log memory info before
		before, err := client.Info(ctx, "memory").Result()
		if err != nil {
			log.Fatal(err.Error())
		}
		for i := 0; i < size; i++ {
			key = strconv.Itoa(i)
			client.Set(ctx, key, str.String(), 0)
		}
		after, err := client.Info(ctx, "memory").Result()
		if err != nil {
			log.Fatal(err.Error())
		}

		name := "./docs/week08/q2/" + strconv.Itoa(lsiz) + "result_" + strconv.Itoa(size) + ".log"

		fer := ioutil.WriteFile(
			name,
			[]byte("before:  \n"+before+"\n\n"+"after: \n"+after),
			0644,
		)
		if fer != nil {
			log.Fatal(fer.Error())
		} else {
			tname, zerr := filepath.Abs(name)
			if zerr != nil {
				log.Println("写入文件成功", name)
			} else {
				log.Println("写入文件成功", tname)
			}
		}
		client.FlushDB(ctx)
		time.Sleep(5 * time.Second)
	}

}

func InitRedisClient(cxt context.Context) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "0932",
		DB:       0,
	})
	return client, client.Ping(cxt).Err()
}
