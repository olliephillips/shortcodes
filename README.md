# Shortcodes

## A Ponzu CMS addon

Shortcodes allows content items to contain shortcode tags i.e ```[myshortcode]``` which can be automatically substituted with alternative at the point the content is served over the Ponzu Server API.

Search and replacement is possible across the entire response as it is made on the byte slice output itself. 

Replacements are sanitized and JSON escaped to prevent injection and the JSON object becoming malformed.

## Usage

## Adding

Use the ```ponzu``` cli tool in your project to add the addon.

```bash
ponzu add github.com/olliephillips/shortcodes
ponzu build
ponzu run
```

As with all Ponzu addons, you must have the addon included as an import
for it to be built into Ponzu server. In the short term you can do this to include in the your Ponzu server application:

```go
import _ "github.com/olliephillips/shortcodes"
```

Once installed, and with your Ponzu server application running, visit the Addons link and enable the addon.

## Implementing on a content item

In the content item, for example ```review```, override the hookable interface ```BeforeAPIResponse``` hook with something like the following:

```go
func (r *Review) BeforeAPIResponse(res http.ResponseWriter, req *http.Request,      data []byte) ([]byte, error) {
	var replaced []byte
	replaced, err := shortcodes.Replace(data)
	if err != nil {
		return replaced, err
	}
	return replaced, nil
}
```

Note this hook differs from most other Ponzu hooks in that in addition to ```http.ResponseWriter``` and ```*http.Request``` it receives a ```[]byte``` slice containing the reponse that will be output by Ponzu API.

Similarly this hook expects a ```[]byte``` slice return in addition to the ```error``` type, whether modified or not.

## Creating shortcodes

On installation, a ```Shortcode``` content item is added, with the following implementation:

```go
type Shortcode struct {
	item.Item

	Tag         string `json:"tag"`
	Description string `json:"description"`
	Replacement string `json:"replacement"`
	Active      bool   `json:"active"`
}
```

This type implements ```MarshalEditor``` so Ponzu provides the standard CRUD interface through which shortcodes can be added and maintained.

In the tag, you should not include the wrapping square brackets, however in content items the tag should be wrapped in sqaure brackets.