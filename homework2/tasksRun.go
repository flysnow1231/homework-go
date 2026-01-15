package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Task func(ctx context.Context) error

type TaskResult struct {
	Index     int
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
	Err       error
}

// Scheduler 支持并发度控制
type Scheduler struct {
	Concurrency int // 并发度，<=0 表示不限制（等同于全部并发）
}

// Run 并发执行 tasks，返回每个任务的执行结果
func (s Scheduler) Run(ctx context.Context, tasks []Task) []TaskResult {
	n := len(tasks)
	results := make([]TaskResult, n)

	// 不限制并发：直接每个任务开一个 goroutine
	if s.Concurrency <= 0 || s.Concurrency >= n {
		var wg sync.WaitGroup
		wg.Add(n)
		for i, t := range tasks {
			i, t := i, t
			go func() {
				defer wg.Done()
				results[i] = runOne(ctx, i, t)
			}()
		}
		wg.Wait()
		return results
	}

	// 限制并发：worker pool
	type job struct {
		index int
		task  Task
	}
	jobs := make(chan job)

	var wg sync.WaitGroup
	worker := func() {
		defer wg.Done()
		for j := range jobs {
			results[j.index] = runOne(ctx, j.index, j.task)
		}
	}

	// 启动 worker
	wg.Add(s.Concurrency)
	for i := 0; i < s.Concurrency; i++ {
		go worker()
	}

	// 投递任务（支持 ctx 取消）
	for i, t := range tasks {
		select {
		case <-ctx.Done():
			// ctx 被取消：未投递的任务填充取消信息
			for k := i; k < n; k++ {
				results[k] = TaskResult{
					Index:     k,
					StartTime: time.Time{},
					EndTime:   time.Time{},
					Duration:  0,
					Err:       ctx.Err(),
				}
			}
			close(jobs)
			wg.Wait()
			return results
		case jobs <- job{index: i, task: t}:
		}
	}

	close(jobs)
	wg.Wait()
	return results
}

func runOne(ctx context.Context, idx int, t Task) TaskResult {
	start := time.Now()
	err := t(ctx)
	end := time.Now()
	return TaskResult{
		Index:     idx,
		StartTime: start,
		EndTime:   end,
		Duration:  end.Sub(start),
		Err:       err,
	}
}

// ======= demo =======
func runTask() {
	tasks := []Task{
		func(ctx context.Context) error {
			time.Sleep(300 * time.Millisecond)
			fmt.Printf("task 1 sleep 300 mili\n")
			return nil
		},
		func(ctx context.Context) error {
			time.Sleep(4000 * time.Millisecond)
			fmt.Printf("task 1 sleep 120 mili\n")
			return fmt.Errorf("task failed")
		},
		//func(ctx context.Context) error {
		//	select {
		//	case <-ctx.Done():
		//		return ctx.Err()
		//	case <-time.After(200 * time.Millisecond):
		//		return nil
		//	}
		//},
	}

	// 全局超时示例
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	s := Scheduler{Concurrency: 2}
	results := s.Run(ctx, tasks)

	// 打印结果
	var total time.Duration
	for _, r := range results {
		fmt.Printf("task=%d duration=%v err=%v\n", r.Index, r.Duration, r.Err)
		total += r.Duration
	}
	fmt.Printf("sum(duration)=%v\n", total)

}
