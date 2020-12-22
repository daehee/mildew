package mildew

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMildew_GetRoots(t *testing.T) {
	testSubs := []string{
		"www.public.navy.mil",
		"www.benning.army.mil",
		"www.192wg.ang.af.mil",
		"www.airforcemedicine.af.mil",
		"www.af.mil",
		"www.dau.mil",
	}
	want := []string{
		"navy.mil",
		"army.mil",
		"af.mil",
		"dau.mil",
	}

	mw := NewMildew()

	for _, v := range testSubs {
		mw.Subs.Insert(v)
	}

	got := mw.GetRoots()

	assert.Len(t, got, len(want))

	for _, v := range want {
		assert.Contains(t, got, v)
	}

}
