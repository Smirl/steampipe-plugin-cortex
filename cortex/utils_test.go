package cortex

import (
	"testing"
)

func TestUnescapeSurrogatePairs(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "rocket emoji surrogate pair is resolved to real UTF-8",
			in:   `{"description":"it's shipped \uD83D\uDE80"}`,
			want: "{\"description\":\"it's shipped 🚀\"}",
		},
		{
			name: "multiple surrogate pairs in the same string",
			in:   `{"description":"\uD83D\uDC4D/\uD83D\uDC4E feedback"}`,
			want: "{\"description\":\"👍/👎 feedback\"}",
		},
		{
			name: "non-surrogate escapes are left untouched",
			in:   `{"name":"Skyscanner\/hotel-itinerary-scoring-ml"}`,
			want: `{"name":"Skyscanner\/hotel-itinerary-scoring-ml"}`,
		},
		{
			name: "plain ASCII is left untouched",
			in:   `{"name":"trackolding"}`,
			want: `{"name":"trackolding"}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := string(unescapeSurrogatePairs([]byte(c.in)))
			if got != c.want {
				t.Errorf("unescapeSurrogatePairs(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

func TestYamlUnmarshalJSONHandlesEmoji(t *testing.T) {
	type entity struct {
		Description string `yaml:"description"`
	}

	// This is exactly the shape of payload that previously made yaml.Unmarshal
	// fail with "found invalid Unicode character escape code" for the whole
	// response, aborting the Cortex entity listing.
	data := []byte(`{"description":"A Claude Code skill that takes you from \"here's an idea\" to \"it's shipped\" \uD83D\uDE80"}`)

	var e entity
	if err := yamlUnmarshalJSON(data, &e); err != nil {
		t.Fatalf("yamlUnmarshalJSON returned an error: %v", err)
	}

	want := "A Claude Code skill that takes you from \"here's an idea\" to \"it's shipped\" 🚀"
	if e.Description != want {
		t.Errorf("Description = %q, want %q", e.Description, want)
	}
}
