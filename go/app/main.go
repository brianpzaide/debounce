package main

import (
	"context"
	"debounce"
	"fmt"
	"time"
)

func greet(name string) string {
	return fmt.Sprintf("Hi %s, Have a nice day!!", name)
}

func main() {
	ctx := context.Background()

	db := debounce.NewDebouncer(ctx, greet, 2*time.Second)

	db.Input <- "Alice"
	time.Sleep(500 * time.Millisecond)
	db.Input <- "Bob"
	time.Sleep(500 * time.Millisecond)
	db.Input <- "Charlie"

	select {
	case res := <-db.Promise():
		fmt.Println("Got:", res)
	case <-time.After(3 * time.Second):
		fmt.Println("Timeout!")
	}
}
