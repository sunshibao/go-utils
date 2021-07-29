// Copyright 2021 Sunshibao <664588619@qq.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package cache

import (
	"fmt"
)

// ----------------------------------------
//  Redis 键
// ----------------------------------------
type RedisKey string

func (k RedisKey) String() string {
	return string(k)
}

func (k RedisKey) Format(args ...interface{}) string {
	if len(args) == 0 {
		return k.String()
	} else {
		return fmt.Sprintf(k.String(), args...)
	}
}
