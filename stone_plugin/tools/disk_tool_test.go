package tools

import (
	"testing"
)

func TestDf(t *testing.T) {
	dfs, err := ParseDf()
	if err != nil {
		t.Error(err)
	}
	if len(dfs) == 0 {
		t.Errorf("no found df info")
	}

	for _, df := range dfs {
		t.Logf("%+v", df)
	}
}
