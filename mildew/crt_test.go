package mildew

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMildew_ScrapeCrts(t *testing.T) {
	testSubs := []string{
		// "www.public.navy.mil",
		// "www.benning.army.mil",
		// "www.192wg.ang.af.mil",
		// "www.airforcemedicine.af.mil",
		// "www.af.mil",
		"www.dau.mil",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mw := NewMildew()
	for _, v := range testSubs {
		mw.Subs.Insert(v)
	}
	err := mw.ScrapeCrts(ctx)
	assert.NoError(t, err)
	fmt.Printf("%v\n", mw.Subs.Slice())
}
