package main

/* This example creates three bars: b1, b2, and b3.
 * b1 is created in the main thread but runs in a separate thread.
 * b2 is created in b1's thread, but runs in its own thread.
 * b3 is created and run in its own thread.
 * b1 and b3 appear immediately, but b2 comes later.
 * b2 is the fastest and should complete and be removed first.
 * b1 is second fastest and should complete and be removed second.
 * b3 is the slowest and should complete last.
 */

import (
	"time"

	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

const sleep = 100 * time.Millisecond

func main() {
	p := mpb.New()

	b1 := p.AddBar(int64(100), mpb.PrependDecorators(decor.StaticName("bar1", 0, 0)))
	// Run thread for b1.
	go func() {
		defer p.RemoveBar(b1)

		// Create and run thread for b2, which starts after a time.
		go func() {
			time.Sleep(10 * sleep)
			b2 := p.AddBar(int64(100), mpb.PrependDecorators(decor.StaticName("bar2", 0, 0)))
			defer p.RemoveBar(b2)
			for j := 0; !b2.Completed(); j++ {
				b2.IncrBy(10) // fastest
				time.Sleep(sleep)
			}
		}()

		for i := 0; !b1.Completed(); i++ {
			b1.IncrBy(2) // second fastest
			time.Sleep(sleep)
		}
	}()

	// Create and run thread for b3, which starts immediately.
	go func() {
		b3 := p.AddBar(int64(100), mpb.PrependDecorators(decor.StaticName("bar3", 0, 0)))
		defer p.RemoveBar(b3)
		for k := 0; !b3.Completed(); k++ {
			b3.IncrBy(1) // slowest
			time.Sleep(sleep)
		}
	}()

	p.Wait()
}
