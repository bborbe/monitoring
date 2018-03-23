package lock

import (
	"os"
	"testing"
	"time"

	. "github.com/bborbe/assert"
)

func TestLocking(t *testing.T) {
	var err error
	lockName := os.TempDir() + "/bla.lock"
	l1 := NewLock(lockName)
	result := true
	err = l1.Lock()
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		l2 := NewLock(lockName)
		erro := l2.Lock()
		if erro != nil {
			t.Fatal(erro)
		}
		result = false
		erro = l2.Unlock()
		if erro != nil {
			t.Fatal(erro)
		}
	}()
	err = AssertThat(result, Is(true))
	if err != nil {
		t.Fatal(err)
	}
	err = l1.Unlock()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(100 * time.Millisecond)
	err = AssertThat(result, Is(false))
	if err != nil {
		t.Fatal(err)
	}
}
