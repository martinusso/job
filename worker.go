package worker

import (
	"context"
	"sync"
	"time"
)

var (
	cancel      context.CancelFunc
	ctx         context.Context
	initialized bool
	initMutex   sync.Mutex
	jobs        []Job
)

func Close() {
	initMutex.Lock()
	defer initMutex.Unlock()
	if initialized {
		initialized = false
	}
}

func Init() error {
	initMutex.Lock()
	defer initMutex.Unlock()
	if !initialized {
		ctx, cancel = context.WithCancel(context.Background())
		gracefulStop(ctx, cancel)
		initialized = true
	}
	return nil
}

// Registar a Job
func Register(name string, fn jobFunc, interval time.Duration) {
	j := Job{
		Name:     name,
		Fn:       fn,
		Interval: interval,
	}
	jobs = append(jobs, j)
}

func run(ctx context.Context, monitor *sync.WaitGroup, job Job, logger Logger) {
	monitor.Add(1)

	go func() {
		defer func() {
			defer monitor.Done()
		}()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				done, err := job.Fn()
				if err != nil {
					logger.Error(err)
				} else if !done {
					continue
				}
				timeout := time.After(job.Interval)
				select {
				case <-ctx.Done():
					return
				case <-timeout:
				}
			}
		}
	}()

	return
}

// Work executes all registered jobs
func Work(logger Logger) error {
	if err := Init(); err != nil {
		return err
	}
	defer Close()

	var monitor sync.WaitGroup
	for _, j := range jobs {
		run(ctx, &monitor, j, logger)
	}
	monitor.Wait()
	return nil
}
