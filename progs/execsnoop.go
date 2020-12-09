/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package progs

import (
	"context"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v8"
)

func ExecSnoop() Program {
	return program{path: "/usr/share/bcc/tools/execsnoop",
		interpreter: execsnoop_interpret}
}

func execsnoop_interpret(line string, rdb *redis.Client) {
	fmt.Println("execsnoop_interpret: " + line)
	ctx := context.Background()
	var pcomm string
	var pid, ppid int64
	var ret int
	var arg1 string
	fmt.Sscanf(line, "%s %d %d %d %s", &pcomm, &pid, &ppid, &ret, &arg1)
	if pid == 0 {
		return
	}

	key := fmt.Sprintf("execsnoop:%d", pid)
	hset := rdb.HSet(ctx,
		key,
		[]string{"pcomm", pcomm,
			"pid", fmt.Sprintf("%d", pid),
			"ppid", fmt.Sprintf("%d", ppid),
			"ret", fmt.Sprintf("%d", ret),
			"arg1", arg1})

	if err := hset.Err(); err != nil {
		panic(err)
	}

	expire := rdb.Expire(ctx, key, time.Second*5)
	if err := expire.Err(); err != nil {
		panic(err)
	}
}
