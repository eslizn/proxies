package proxies

import (
	"context"
	"testing"
)

func TestAvailable(t *testing.T) {
	err := Available(context.Background(), "http://119.81.189.194:80")
	if err != nil {
		t.Error(err)
		return
	}
}
