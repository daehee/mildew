package mildew

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMildew_ScrapeDirs(t *testing.T) {
	want := []string{
		"www.public.navy.mil",
		"www.benning.army.mil",
		"www.192wg.ang.af.mil",
		"www.airforcemedicine.af.mil",
		"www.af.mil",
		"www.dau.mil",
	}
	notWant := []string{
		"www.facebook.com",
		"www.linkedin.com",
		"nationalguard.com",
		"www.armymwr.com",
		"diversity.defense.gov",
		"www.afneurope.net",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mw := NewMildew()
	err := mw.ScrapeDirs(ctx)
	assert.NoError(t, err)
	res := mw.Subs.Slice()
	fmt.Printf("unique subdomains: %d\n", mw.Subs.Len())
	for _, v := range want {
		assert.Contains(t, res, v)
	}
	for _, v := range notWant {
		assert.NotContains(t, res, v)
	}
}
