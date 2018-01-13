// Package ctxutils defines context utility functions.
package ctxutils

import (
	"context"
)

// WithFail returns a copy of parent with a new Done channel. The returned
// context's Done channel is closed when the returned fail function is
// called or when the parent context's Done channel is closed, whichever
// happens first. Fail semantic is the following: when the context has yet
// to be cancelled - so fail is called for the first time - it cancel the
// associated context and return the error; subsequent calls to fail result
// in no further action, returning the first error.
//
// Calling fail on this context releases resources associated with it, so
// code should call it as soon as the operations running in this Context
// complete. Fail is thread safe
func WithFail(parent context.Context) (ctx context.Context, fail func(error) error) {
	//fail context
	ctx, cancel := context.WithCancel(parent)

	//failErr is used to store the error for which the context is cancelled
	var failErr error
	//channel for permission to set failErr error - one time only use
	once := make(chan struct{}, 1)
	once <- struct{}{}
	//channel for signaling failErr being set
	done := make(chan struct{})

	//the returned fail function
	fail = func(err error) error {
		select {
		case <-once: //ask permission
			//set error state, possibly giving priority to parent error
			if failErr = parent.Err(); failErr == nil {
				failErr = err
			}
			//cancel the context - NB: the context may be already cancelled by the parent
			cancel()

			//notify others of the status change
			close(done)
		case <-done: //wait for the failErr to be setted
			//do nothing
		}
		return failErr
	}

	return
}
