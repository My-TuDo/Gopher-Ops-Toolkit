package prober

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

// ——— 测试用探针 ———

// fastProbe 快速返回（5ms），用于高并发测试
type fastProbe struct{}

func (fastProbe) Probe(_ string, _ time.Duration) Result {
	time.Sleep(5 * time.Millisecond)
	return Result{Status: "健康", Latency: 5}
}
func (fastProbe) Name() string { return "fast" }

// ———————————————————————————————————————————————
// 好版本：MultiRoundProbe 内部通过信号量限流
// ———————————————————————————————————————————————
func TestGoroutineControl_Good(t *testing.T) {
	baseline := runtime.NumGoroutine()
	t.Logf("启动前 goroutine: %d", baseline)

	p := fastProbe{}

	// 单次 MultiRoundProbe 内部发起 500 轮并发探测，
	// 受 maxConcurrentProbes=5 信号量控制，不会炸 goroutine
	var peakDelta int
	var mu sync.Mutex

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = MultiRoundProbe(p, "test-target", 30*time.Second, 500)
	}()

	// 探测期间持续采样 goroutine 数量
	for i := 0; i < 8; i++ {
		time.Sleep(60 * time.Millisecond)
		current := runtime.NumGoroutine()
		delta := current - baseline
		mu.Lock()
		if delta > peakDelta {
			peakDelta = delta
		}
		mu.Unlock()
		t.Logf("  采样 %d: goroutine=%d (增量=%d)", i+1, current, delta)
	}

	wg.Wait()
	time.Sleep(300 * time.Millisecond)
	endDelta := runtime.NumGoroutine() - baseline
	t.Logf("结束后 goroutine=%d (增量=%d)", runtime.NumGoroutine(), endDelta)

	// 信号量 maxConcurrentProbes=5，加上系统开销，峰值增量应 ≤ 20
	if peakDelta > 20 {
		t.Errorf("❌ 信号量失效？峰值增量 %d > 20", peakDelta)
	} else {
		fmt.Printf("✅ 好版本: 500 轮并发，峰值增量仅 %d goroutine\n", peakDelta)
	}

	if endDelta > 10 {
		t.Errorf("❌ 可能 goroutine 泄漏: 结束后增量 %d", endDelta)
	} else {
		fmt.Printf("✅ 无 goroutine 泄漏: 结束后增量 %d\n", endDelta)
	}
}

// ———————————————————————————————————————————————
// 坏版本：不用信号量，直接起 goroutine
// ———————————————————————————————————————————————
func TestGoroutineControl_Bad(t *testing.T) {
	baseline := runtime.NumGoroutine()
	t.Logf("启动前 goroutine: %d", baseline)

	rounds := 500
	done := make(chan struct{}, rounds)

	for i := 0; i < rounds; i++ {
		go func() {
			time.Sleep(200 * time.Millisecond)
			done <- struct{}{}
		}()
	}

	time.Sleep(100 * time.Millisecond)
	peak := runtime.NumGoroutine()
	peakDelta := peak - baseline
	t.Logf("坏版本峰值 goroutine=%d (增量=%d)", peak, peakDelta)

	for i := 0; i < rounds; i++ {
		<-done
	}
	time.Sleep(200 * time.Millisecond)

	fmt.Printf("  坏版本: %d 轮无控制 → 峰值增量 %d goroutine\n", rounds, peakDelta)
	fmt.Printf("  好版本: 500 轮 MultiRoundProbe → 峰值增量 ≤20 goroutine\n")
	fmt.Printf("  📊 结论: 信号量将 goroutine 从 %d 控制在 20 以内, 降低 %.1f%%\n",
		peakDelta, (1-20.0/float64(peakDelta))*100)
}
