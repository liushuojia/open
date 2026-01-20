package worker

import (
	"fmt"
	"testing"
	"time"
)

func Test_POOL(t *testing.T) {

	Run(func() error {
		time.Sleep(time.Second)
		fmt.Println("1 second")
		return nil
	}, func() error {
		time.Sleep(3 * time.Second)
		fmt.Println("3 second")
		return nil
	}, func() error {
		time.Sleep(3 * time.Second)
		fmt.Println("5 second")
		return nil
	})
	fmt.Println("no wait")
	time.Sleep(10 * time.Second)

}
