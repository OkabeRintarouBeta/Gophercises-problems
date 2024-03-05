package cyoa

type Story map[string]Chapter

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

func (s *Story) Get(title string) (Chapter, bool) {
	chapter, ok := (*s)[title]
	if ok {
		return chapter, ok
	} else {
		return Chapter{}, ok
	}
}
