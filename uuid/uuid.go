package uuid

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"os"
	"sync/atomic"
	"time"
)

var Continuity = "20060102150405"

func init() {
	rand.Seed(time.Now().UnixNano())
}

var Generate = func() string {
	return generate(time.Now())
}

func ID() string {
	return uuid.New().String()
}

// 生成24位订单号
// 前面17位代表时间精确到毫秒，中间3位代表进程id，最后4位代表序号
func generate(t time.Time) string {
	var num int64
	s := t.Format(Continuity)
	m := t.UnixNano()/1e6 - t.UnixNano()/1e9*1e3
	ms := sup(m, 3)
	p := os.Getpid() % 1000
	ps := sup(int64(p), 3)
	i := atomic.AddInt64(&num, 1)
	r := i % 10000
	rs := sup(r, 4)
	n := fmt.Sprintf("%s%s%s%s", s, ms, ps, rs)
	return n
}

// 对长度不足n的数字前面补0
func sup(i int64, n int) string {
	m := fmt.Sprintf("%d", i)
	for len(m) < n {
		m = fmt.Sprintf("0%s", m)
	}
	return m
}
