package shortcodes

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	"regexp"
	"strings"

	"github.com/ponzu-cms/ponzu/system/addon"
	"github.com/ponzu-cms/ponzu/system/db"
	"github.com/ponzu-cms/ponzu/system/item"
)

// Addon Meta config
var meta = addon.Meta{
	PonzuAddonName:      "Shortcodes",
	PonzuAddonAuthor:    "Ollie Phillips",
	PonzuAddonAuthorURL: "https://slippytrumpet.io",
	PonzuAddonVersion:   "0.1.0",
}

type Shortcodes struct {
	addon.Addon
}

var _ = addon.Register(meta, func() interface{} { return new(Shortcodes) })

func init() {
	config, _ := getConfig()
	if config.Addon.Meta.PonzuAddonStatus == addon.StatusEnabled {
		item.Types["Shortcode"] = func() interface{} { return new(Shortcode) }
	}
}

// Replace is called from an item hook it accepts a byte slice and
// returns a mutated byte slice which includes shortcode replacements or error
func Replace(data []byte) ([]byte, error) {
	// get addon config
	config, err := getConfig()
	if err != nil {
		return data, errors.New("Shortcodes addon: could not obtain configuration")
	}

	// check enabled
	if config.Addon.Meta.PonzuAddonStatus != addon.StatusEnabled {
		return data, fmt.Errorf("Addon %s is not enabled", config.Addon.Meta.PonzuAddonName)
	}

	// parse out shortcodes with regex, format is [shortcode]
	re := regexp.MustCompile(`(?m)\[([a-z]*)\]`)
	matches := re.FindAllStringSubmatch(string(data), -1)

	// if none return data and nil error
	if len(matches) == 0 {
		return data, nil
	}

	// get active shortcodes
	activeSC, err := getActiveShortcodes()
	if err != nil {
		return data, errors.New("Shortcodes addon: could not obtain shortcodes")
	}

	// make replacements in thed data
	data = makeReplacements(data, matches, activeSC)

	return data, nil
}

// gets current config of addon
func getConfig() (*Shortcodes, error) {
	var ps = &Shortcodes{}

	key, err := addon.KeyFromMeta(meta)
	if err != nil {
		return ps, err
	}
	data, err := db.Addon(key)
	if err != nil {
		return ps, err
	}

	err = json.Unmarshal(data, ps)
	if err != nil {
		return ps, err
	}

	return ps, nil
}

// gets active shortcodes in system
func getActiveShortcodes() (map[string]Shortcode, error) {
	activeSC := make(map[string]Shortcode)
	items := db.ContentAll("Shortcode")
	for i := range items {
		item := Shortcode{}
		err := json.Unmarshal(items[i], &item)
		if err != nil {
			return activeSC, err
		}
		// if the shortcode active add it to data map
		// for later lookups of tag
		if item.Active == true {
			activeSC[item.Tag] = item
		}
	}

	return activeSC, nil
}

// performs the actual subtitution of shortcodes with their replacement data
func makeReplacements(data []byte, matches [][]string, sc map[string]Shortcode) []byte {
	content := string(data)
	for _, m := range matches {
		if val, ok := sc[m[1]]; ok {
			sub := fmt.Sprintf("[%s]", val.Tag)
			// sanitise every replacement
			rep := makeSafeToOutput(val.Replacement)
			content = strings.Replace(content, sub, rep, -1)
		}
	}

	return []byte(content)
}

// ensures shortcode replacements are safe to substitute into what
// is already marshalled and properly encoded/escaped data
func makeSafeToOutput(content string) string {
	buf := new(bytes.Buffer)
	esc := html.EscapeString(content)
	json.HTMLEscape(buf, []byte(esc))
	return string(buf.Bytes())
}

// AfterEnable is a hook which runs after addon is enabled.
// We use it to add a shortcode item for creation and storage
// of shortcodes
func (s *Shortcodes) AfterEnable(res http.ResponseWriter, req *http.Request) error {
	item.Types["Shortcode"] = func() interface{} { return new(Shortcode) }
	return nil
}

// AfterDisable is a hook which runs after addon is disabled.
// We use it to remove shortcode item from Ponzu server
func (s *Shortcodes) AfterDisable(res http.ResponseWriter, req *http.Request) error {
	delete(item.Types, "Shortcode")
	return nil
}
