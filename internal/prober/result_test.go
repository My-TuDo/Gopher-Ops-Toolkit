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

	res := MultiRoundProbe(mock, "test-target", time.Second, 1)
	if res.Status != "健康" {
		t.Errorf("期望健康，得到 %s", res.Status)
	}
	if res.Latency != 10 {
		t.Errorf("期望延迟 10，得到 %d", res.Latency)
	}
}

// 测试 MultiRoundProbe 函数在多轮探测时的行为
func TestMultiRoundProbe_MultiRound(t *testing.T) {

}

// 测试 MultiRoundProbe 函数在多轮探测时的行为，模拟所有探测失败的情况
func TestMultiRoundProbe_AllFail(t *testing.T) {}
