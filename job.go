package worker

import "time"

type jobFunc func() (bool, error)

type Job struct {
	Name     string
	Fn       jobFunc
	Interval time.Duration
}
