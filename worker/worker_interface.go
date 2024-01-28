package message

import "sync"

type Worker interface {
	ListenMessageQueue(wg *sync.WaitGroup)
}
