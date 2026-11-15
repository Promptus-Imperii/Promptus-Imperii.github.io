package main

import (
	"strings"
	"testing"
)

func TestErrorsHTMLEscapes(t *testing.T) {
	got := errorsHTML([]string{"<b>bad</b>"})
	if strings.Contains(got, "<b>") || !strings.Contains(got, "&lt;b&gt;bad&lt;/b&gt;") {
		t.Errorf("not escaped: %s", got)
	}
	if !strings.HasPrefix(got, `<hx-partial hx-target="#errors" hx-swap="outerHTML"><div id="errors"`) {
		t.Errorf("missing partial or id: %s", got)
	}
	if empty := errorsHTML(nil); empty != `<hx-partial hx-target="#errors" hx-swap="outerHTML"><div id="errors" class="flex flex-col text-red-400"></div></hx-partial>` {
		t.Errorf("empty = %s", empty)
	}
}
