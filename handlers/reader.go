// Copyright 2014 Canonical Ltd.
// Licensed under the LGPLv3 with static-linking exception.
// See LICENCE file for details.

package handlers

import (
	"../ratelimit"
	"io"
)

type reader struct {
	r      io.Reader
	bucket *ratelimit.Bucket
}

// Reader returns a reader that is rate limited by
// the given token bucket. Each token in the bucket
// represents one byte.
func Reader(r io.Reader, bucket *ratelimit.Bucket) CustomReader {
	return &reader{
		r:      r,
		bucket: bucket,
	}
}

func (r *reader) Read(buf []byte) (int, error) {
	n, err := r.r.Read(buf)
	if n <= 0 {
		return n, err
	}
	r.bucket.Wait(int64(n))
	return n, err
}

func (r *reader) ReadInjectedWithChannel(buf []byte, signal func(int64)) (int, error) {
	n, err := r.r.Read(buf)
	if n <= 0 {
		return n, err
	}
	latestTick := r.bucket.Wait(int64(n))
	if latestTick != 0 {
		signal(latestTick)
	}
	return n, err
}

type writer struct {
	w      io.Writer
	bucket *ratelimit.Bucket
}

// Writer returns a reader that is rate limited by
// the given token bucket. Each token in the bucket
// represents one byte.
func Writer(w io.Writer, bucket *ratelimit.Bucket) io.Writer {
	return &writer{
		w:      w,
		bucket: bucket,
	}
}

func (w *writer) Write(buf []byte) (int, error) {
	w.bucket.Wait(int64(len(buf)))
	return w.w.Write(buf)
}
