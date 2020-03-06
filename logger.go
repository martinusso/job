package worker

type Logger interface {
	Error(err error)
}
