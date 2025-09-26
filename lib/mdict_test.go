package lib

import "testing"

func MddTst(t *testing.T) {
	mdict, err := New("/Users/eloxt/Downloads/牛津高阶英汉双解词典（第10版）V3/牛津高阶英汉双解词典（第10版）V3.mdd")
	if err != nil {
		t.Fatal(err)
	}

	err = mdict.BuildIndex()
	if err != nil {
		t.Fatal(err)
	}
}
