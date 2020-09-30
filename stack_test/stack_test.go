package stack_test

import (
	"sync"
	"testing"

	"github.com/liangyaopei/goroutineinfo"
)

const (
	cnt = 100
)

func TestGetInfoSingle(t *testing.T) {
	stacks := goroutineinfo.GetInfo(false)
	for _, stack := range stacks {
		t.Logf("id:[%d],state:[%s]", stack.ID(), stack.State())
	}
}

func TestGetInfoAll(t *testing.T) {
	var (
		wg     = sync.WaitGroup{}
		stacks []goroutineinfo.Stack
	)
	wg.Add(cnt)
	for i := 0; i < cnt; i++ {
		go func() {
			wg.Done()
		}()
	}
	stacks = goroutineinfo.GetInfo(true)
	for _, stack := range stacks {
		t.Logf("wait,id:[%d],state:[%s]", stack.ID(), stack.State())
	}
	wg.Wait()
	stacks = goroutineinfo.GetInfo(true)
	for _, stack := range stacks {
		t.Logf("receive,id:[%d],state:[%s]", stack.ID(), stack.State())
	}
}

func BenchmarkGetInfoSingle(b *testing.B) {
	_ = goroutineinfo.GetInfo(false)
}

func BenchmarkGetInfoAll(b *testing.B) {
	var wg = sync.WaitGroup{}
	wg.Add(cnt)
	for i := 0; i < cnt; i++ {
		go func() {
			wg.Done()
		}()
	}
	_ = goroutineinfo.GetInfo(true)
	wg.Wait()
	_ = goroutineinfo.GetInfo(true)
}
