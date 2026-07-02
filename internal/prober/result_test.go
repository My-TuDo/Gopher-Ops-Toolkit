package prober

import (
	"testing"
	"time"
)

type mockProber struct {
	result Result
	name   string
}

func (m mockProber) Probe(target string, timeout time.Duration) Result {
	return m.result
}

func (m mockProber) Name() string {
	return m.name
}

// 测试 MultiRoundProbe 函数在单轮探测时的行为
func TestMultiRoundProbe_SingleRound(t *testing.T) {
	// rounds = 1 时，应该直接调用 Probe 方法
	mock := mockProber{
		name: "Mock 测试",
		result: Result{
			Status:  "健康",
			Latency: 10,
			Detail:  "OK",
		},
	}

	res := MultiRoundProbe(mock, "test-target", 3*time.Second, 1)
	if res.Status != "健康" {
		t.Errorf("期望健康，得到 %s", res.Status)
	}
	if res.Latency != 10 {
		t.Errorf("期望延迟 10，得到 %d", res.Latency)
	}
}

// 测试 MultiRoundProbe 函数在多轮探测时的行为
func TestMultiRoundProbe_MultiRound(t *testing.T) {
	// rounds > 1 时，应该进行多轮探测并计算平均值
	mock := mockProber{
		name: "Mock 测试",
		result: Result{
			Status:  "健康",
			Latency: 10,
			Detail:  "OK",
		},
	}

	res := MultiRoundProbe(mock, "test-target", 3*time.Second, 5)
	if res.Status != "健康" {
		t.Errorf("期望健康，得到 %s", res.Status)
	}
}

// 测试 MultiRoundProbe 函数在多轮探测时的行为，模拟所有探测失败的情况
func TestMultiRoundProbe_AllFail(t *testing.T) {
	// rounds > 1 时，模拟所有探测失败的情况
	mock := mockProber{
		name: "Mock 测试",
		result: Result{
			Status: "不健康",
			Error:  "模拟失败",
		},
	}

	res := MultiRoundProbe(mock, "test-target", time.Second, 3)
	if res.Status != "不健康" {
		t.Errorf("期望不健康，得到 %s", res.Status)
	}
	if res.Error == "" {
		t.Error("期望有错误信息，但为空")
	}
}

// ———————————————————————————————————————————————————————————
// Benchmark: MultiRoundProbe 并发性能基准
// ———————————————————————————————————————————————————————————

// 模拟延时 5ms 的探针
type slowMockProber struct{}

func (slowMockProber) Probe(_ string, _ time.Duration) Result {
	time.Sleep(5 * time.Millisecond)
	return Result{Status: "健康", Latency: 5}
}
func (slowMockProber) Name() string { return "bench-prober" }

// 基准 1：串行单轮（baseline）
func BenchmarkMultiRoundProbe_Serial(b *testing.B) {
	p := slowMockProber{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MultiRoundProbe(p, "target", time.Second, 1)
	}
}

// 基准 2：并发 10 轮
func BenchmarkMultiRoundProbe_Concurrent10(b *testing.B) {
	p := slowMockProber{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MultiRoundProbe(p, "target", time.Second, 10)
	}
}

// 基准 3：并发 50 轮（熔断场景）
func BenchmarkMultiRoundProbe_Concurrent50(b *testing.B) {
	p := slowMockProber{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MultiRoundProbe(p, "target", time.Second, 50)
	}
}
