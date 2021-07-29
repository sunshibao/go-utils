// Author: Qingshan Luo <edoger@qq.com>
package uuid

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// UUID 返回一个32位的不重复字符串
func UUID() string {
	if id, err := uuid.NewRandom(); err == nil {
		return strings.ReplaceAll(id.String(), "-", "")
	}
	a := md5.Sum(strconv.AppendInt(strconv.AppendUint(nil, rand.Uint64(), 10), time.Now().UnixNano(), 10))
	return hex.EncodeToString(a[:])
}
