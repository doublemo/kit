package networks

import (
	"fmt"
	"testing"
)

func TestHttp(t *testing.T) {
	c := NewWebConfig(":8080")
	fmt.Println(c.MaxHeaderBytes)
	w := NewWeb()
	w.Serve(c)
}
