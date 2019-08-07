package shortcodes

import (
	"fmt"

	"github.com/ponzu-cms/ponzu/management/editor"
	"github.com/ponzu-cms/ponzu/system/item"
)

type Shortcode struct {
	item.Item

	Tag         string `json:"tag"`
	Description string `json:"description"`
	Replacement string `json:"replacement"`
	Active      bool   `json:"active"`
}

// MarshalEditor writes a buffer of html to edit a Shortcode within the CMS
// and implements editor.Editable
func (s *Shortcode) MarshalEditor() ([]byte, error) {
	view, err := editor.Form(s,
		editor.Field{
			View: editor.Input("Tag", s, map[string]string{
				"label":       "Tag",
				"type":        "text",
				"placeholder": "Enter the name of shortcode tag",
			}),
		},

		editor.Field{
			View: editor.Input("Description", s, map[string]string{
				"label":       "Description",
				"type":        "text",
				"placeholder": "Provide a brief description of shortcode behaviour",
			}),
		},

		editor.Field{
			View: editor.Textarea("Replacement", s, map[string]string{
				"label":       "Replacement",
				"placeholder": "Enter the replacement text",
			}),
		},

		editor.Field{
			View: editor.Checkbox("Active", s, map[string]string{
				"label": "Active",
			}, map[string]string{
				"true": "Yes",
			}),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to render Shortcode editor view: %s", err.Error())
	}

	return view, nil
}

// String is implemented to provider user friendly naming
func (s *Shortcode) String() string {
	return fmt.Sprintf("[%s] - %s", s.Tag, s.Description)
}
