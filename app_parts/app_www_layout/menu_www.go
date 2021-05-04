package app_www_layout

import (
	"fmt"
	"strings"

	server_http "github.com/pavlo67/tools/common/server/server_http_v2"
	"github.com/pavlo67/tools/common/server/server_http_v2/server_http_v2_jschmhr/wrapper_page"
	"github.com/pavlo67/tools/common/thread"

	"github.com/pavlo67/tools/common/kv"
)

// wrapper_page.CommonFragments ------------------------------------------------

var _ wrapper_page.CommonFragments = &SetMenu{}

type SetMenu struct {
	Process thread.KVGetString
}

func (sm *SetMenu) Set(fragments server_http.Fragments) (server_http.Fragments, error) {
	if fragments == nil {
		fragments = server_http.Fragments{}
	}

	menu, err := sm.Process.GetString()
	if err != nil {
		return fragments, err
	}

	fragments["left"] = menu
	return fragments, nil
}

// kv.ItemsProcess -------------------------------------------------------------

var _ kv.ItemsProcess = &MenuWWW{}

type MenuItemWWW struct {
	HRef  string
	Title string
}

type KVMenuItem struct {
	Key   []string
	Value MenuItemWWW
}

// !!! called under thread-mutex
type MenuWWW struct {
	items []KVMenuItem
	html  string
}

func (menuWWW *MenuWWW) ProcessOne(kvItem kv.Item) error {
	if menuWWW == nil {
		return fmt.Errorf("no menuWWW to .ProcessOne()")
	}
	menuItemWWW, ok := kvItem.Value.(MenuItemWWW)
	if !ok {
		return fmt.Errorf("wrong kvItem (%#v) to .ProcessOne()", kvItem)
	}

	menuItem := KVMenuItem{
		Key:   kvItem.Key,
		Value: menuItemWWW,
	}

POINTS:
	for i, p := range menuWWW.items {
		if len(p.Key) == len(menuItem.Key) {
			for i := range p.Key {
				if p.Key[i] != menuItem.Key[i] {
					continue POINTS
				}
			}
			menuWWW.items[i].Value = menuItem.Value
			return nil
		}
	}

	menuWWW.items = append(menuWWW.items, menuItem)

	var htmlPoints []string
	for _, menuItem := range menuWWW.items {

		htmlPoints = append(htmlPoints, fmt.Sprintf(`<li><a href="%s">%s</a></li>`, menuItem.Value.HRef, menuItem.Value.Title))
	}
	menuWWW.html = strings.Join(htmlPoints, "\n")

	return nil
}

func (menuWWW *MenuWWW) Finish() (string, error) {
	if menuWWW == nil {
		return "", fmt.Errorf("no menuWWW to .Finish()")
	}

	return menuWWW.html, nil
}
