package python

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseModule(t *testing.T) {
	testCases := []struct {
		line   string
		want   []string
		wantOk bool
	}{
		{
			line:   "import math",
			want:   []string{"math"},
			wantOk: true,
		},
		{
			line:   "  import os, sys",
			want:   []string{"os", "sys"},
			wantOk: true,
		},
		{
			line:   "from django.shortcuts import (render, redirect)",
			want:   []string{"django.shortcuts", "render", "redirect"},
			wantOk: true,
		},
		{
			line:   "from numpy import *",
			want:   []string{"numpy", "*"},
			wantOk: true,
		},
		{
			line:   "import pandas as pd",
			want:   []string{"pandas"},
			wantOk: true,
		},
		{
			line:   "something import",
			want:   nil,
			wantOk: false,
		},
	}

	for _, c := range testCases {
		m, ok := parseModule(c.line)
		assert.Equal(t, ok, c.wantOk)
		assert.EqualValues(t, m, c.want)
	}
}
