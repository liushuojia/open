package utils

import (
	"fmt"
	"testing"
	"time"
)

func Test_ID(t *testing.T) {

	tt := time.Now()

	id1, _ := NewID(time.Now(), 99)
	fmt.Println("cost", time.Now().UnixMilli()-tt.UnixMilli(), "Millisecond")
	fmt.Println("", id1.MinID(), " - ", id1.MaxID())
	fmt.Println("", id1.IDList())
	fmt.Println("")
	fmt.Println("")

	tt = time.Now()
	id2, _ := NewID(time.Now(), 22)
	fmt.Println("cost", time.Now().UnixMilli()-tt.UnixMilli(), "Millisecond")
	fmt.Println("", id2.MinID(), " - ", id2.MaxID())
	fmt.Println("", id2.IDList())
	fmt.Println("")

}
