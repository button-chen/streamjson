# streamjson
Parse json streams in real time

## Example

```go
func TestStreamJson(t *testing.T) {

	cr, cw := net.Pipe()

	js := `{"aa":{"a":1,"b":{"c":2,"d":5}},"bb":[{"a":123},{"a":456}],"cc":[["e","f"],["g","h"]], "ee":999}`

	go func() {
		// Simulates streaming json data
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


output：
- timestamp: 353     aa.b.c: 2     error: <nil>
- timestamp: 713     bb.*.a: 123   error: <nil>
- timestamp: 869     bb.*.a: 456   error: <nil>
- timestamp: 1071    cc.*.*: e     error: <nil>
- timestamp: 1134    cc.*.*: f     error: <nil>
- timestamp: 1226    cc.*.*: g     error: <nil>
- timestamp: 1289    cc.*.*: h     error: <nil>
- timestamp: 1477    ee: 999       error: <nil>
- send end:  1492ms
- timestamp: 1492    name: <nil>   error: not find pattern: name

```

Or use MonitorAll：

```go
sj.MonitorAll(func(key string, a any) {
	fmt.Printf("timestamp: %v key: %v  val: %v\n", time.Since(start).Milliseconds(), key, a)
})

output:
- timestamp: 204   key: aa.a    val: 1
- timestamp: 388   key: aa.b.c  val: 2
- timestamp: 508   key: aa.b.d  val: 5
- timestamp: 786   key: bb.*.a  val: 123
- timestamp: 960   key: bb.*.a  val: 456
- timestamp: 1190  key: cc.*.*  val: e
- timestamp: 1258  key: cc.*.*  val: f
- timestamp: 1357  key: cc.*.*  val: g
- timestamp: 1427  key: cc.*.*  val: h
- timestamp: 1617  key: ee      val: 999
```



