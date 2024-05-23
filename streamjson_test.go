package streamjson

import (
	"fmt"
	"log"
	"net"
	"testing"
	"time"
)

func TestStreamJson(t *testing.T) {

	cr, cw := net.Pipe()

	js := `{"aa":{"a":1,"b":{"c":2,"d":5}},"bb":[{"a":123},{"a":456}],"cc":[["e","f"],["g","h"]], "ee":999}`

	go func() {
		// 模拟流式发送json数据
		now := time.Now()
		for i := 0; i < len(js); i++ {
			cw.Write(append([]byte(nil), js[i]))
			time.Sleep(time.Millisecond * 10)
		}
		fmt.Printf("send end: %vms\n", time.Since(now).Milliseconds())
		cw.Close()
	}()

	sj := NewStreamJson()

	start := time.Now()
	sj.AddMonitor("name", func(a any, err error) {
		fmt.Printf("timestamp: %v name: %v  error: %v\n", time.Since(start).Milliseconds(), a, err)
	})
	sj.AddMonitor("ee", func(a any, err error) {
		fmt.Printf("timestamp: %v ee: %v  error: %v\n", time.Since(start).Milliseconds(), a, err)
	})
	sj.AddMonitor("aa.b.c", func(a any, err error) {
		fmt.Printf("timestamp: %v aa.b.c: %v  error: %v\n", time.Since(start).Milliseconds(), a, err)
	})
	sj.AddMonitor("cc.*.*", func(a any, err error) {
		fmt.Printf("timestamp: %v cc.*.*: %v  error: %v\n", time.Since(start).Milliseconds(), a, err)
	})
	sj.AddMonitor("bb.*.a", func(a any, err error) {
		fmt.Printf("timestamp: %v bb.*.a: %v  error: %v\n", time.Since(start).Milliseconds(), a, err)
	})

	err := sj.ProcessStream(cr)
	if err != nil {
		log.Println("process stream error: ", err)
	}
}
