package redis

import (
	"fmt"
	"github.com/go-redis/redis"
)

func main() {
	redisDB := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	var (
		getBitInt   *redis.IntCmd
		andBitInt   *redis.IntCmd
		posBitInt   *redis.IntCmd
		countBitInt *redis.IntCmd
		countBit    *redis.BitCount
	)

	redisDB.SetBit("bit_key", 1000, 1)
	getBitInt = redisDB.GetBit("bit_key", 1000)
	countBitInt = redisDB.BitCount("bit_key", countBit)
	// 对"bit_key1","bit_key2"做AND位运算，并保存到"dest_key"中
	andBitInt = redisDB.BitOpAnd("dest_key", "bit_key1", "bit_key2")
	// redisDB.BitPos("bit_key",false);
	posBitInt = redisDB.BitPos("bit_key", 1, 2)

	fmt.Println(getBitInt)
	fmt.Println(andBitInt)
	fmt.Println(countBitInt)
	fmt.Println(posBitInt)
}
