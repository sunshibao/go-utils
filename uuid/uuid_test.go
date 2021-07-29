package uuid

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestUUID(t *testing.T) {
	uuid2 := UUID()

	random, _ := uuid.NewRandom()
	fmt.Println(uuid2)
	fmt.Println(random)

	a := md5.Sum(strconv.AppendInt(strconv.AppendUint(nil, rand.Uint64(), 10), time.Now().UnixNano(), 10))
	fmt.Println(hex.EncodeToString(a[:]))
}
