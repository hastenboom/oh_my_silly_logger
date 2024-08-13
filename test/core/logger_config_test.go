package core_test

import (
	"oh_my_logger/main/core"
	"testing"
)

func TestConstLevel(t *testing.T) {
	t.Logf("level: %d, type: %T", core.DEBUG, core.DEBUG)

}
