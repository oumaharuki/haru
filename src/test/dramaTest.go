package test

import (
	"ctrl"
	"fmt"
	"testing"
)

func TestDrama(t *testing.T) {
	drama := ctrl.ExtractDramaHandle(`</li><li><a href="/vod/play/id/2158/sid/1/nid/3.html">第03集</a></li>`,
		`href="(?s:(.*?))">(?s:(.*?))<a href=`)
	for _, item := range drama {
		fmt.Println(item)
		t.Log(item)
	}

}
