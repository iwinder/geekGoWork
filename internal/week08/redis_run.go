package week08

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func IntRun(alsiz int, flag bool) {
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
		client.FlushDB(ctx)
		// log memory info before
		before, err := client.Info(ctx, "memory").Result()
		if err != nil {
			log.Fatal(err.Error())
		}
		alen := len("used_memory_dataset:")
		beforeBeginPos := strings.Index(before, "used_memory_dataset:") + alen
		beforeEndPos := strings.Index(before, "\r\nused_memory_dataset_perc:")
		bsize := before[beforeBeginPos:beforeEndPos]
		for i := 0; i < size; i++ {
			key = strconv.Itoa(i)
			if flag {
				key = uuid.New().String()
			}

			client.Set(ctx, key, str.String(), 0)
		}
		after, err := client.Info(ctx, "memory").Result()
		if err != nil {
			log.Fatal(err.Error())
		}
		afterBeginPos := strings.Index(after, "used_memory_dataset:") + alen
		afterEndPos := strings.Index(after, "\r\nused_memory_dataset_perc:")
		asize := after[afterBeginPos:afterEndPos]
		var typeFlag string = "type2"
		if flag {
			typeFlag = "type1"
		}
		name := "./docs/week08/q2/" + typeFlag + "/" + strconv.Itoa(lsiz) + "/" + strconv.Itoa(lsiz) + "result_" + strconv.Itoa(size) + ".log"

		fer := ioutil.WriteFile(
			name,
			[]byte("before:  \n"+before+"\n\n"+"after: \n"+after),
			0644,
		)
		afsize, aferr := strconv.ParseFloat(asize, 64)
		if aferr != nil {
			log.Fatal(aferr.Error())
		}
		bfsize, bferr := strconv.ParseFloat(bsize, 64)
		if bferr != nil {
			log.Fatal(bferr.Error())
		}
		avg := (afsize - bfsize - float64(alsiz*size)) / float64(size)
		if fer != nil {
			log.Fatal(fer.Error())
		} else {
			tname, zerr := filepath.Abs(name)
			if zerr != nil {
				log.Println("写入文件成功", name)
			} else {
				log.Println("写入文件成功", tname)
			}
			log.Println("before:", bsize, "after:", afsize, "avg:", avg)
			log.Println("=====================>>>>>>>>>>>>>>>>>")
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
