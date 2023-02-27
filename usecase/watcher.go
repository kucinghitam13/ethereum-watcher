package usecase

type Watcher interface {
	StartWatching()
	StopWatching()
}
