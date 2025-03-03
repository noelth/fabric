package searxng

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"

	"github.com/noelth/fabric/plugins"
)

// Structs for JSON response from SearXNG
type SearXNGSearchResult struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Content string `json:"content"`
}

type SearXNGResponse struct {
	Query   string                `json:"query"`
	Results []SearXNGSearchResult `json:"results"`
}

// SearXNG represents the tool structure (similar to YouTube tool)
type SearXNG struct {
	*plugins.PluginBase
	normalizeRegex *regexp.Regexp
}

// NewSearXNG initializes the tool (following YouTube's `NewYouTube`)
func NewSearXNG() (ret *SearXNG) {
	label := "SearXNG"

	ret = &SearXNG{
		PluginBase: &plugins.PluginBase{
			Name:             label,
			SetupDescription: "SearXNG Service - to query and retrieve search results",
			EnvNamePrefix:    plugins.BuildEnvVariablePrefix(label),
		},
	}

	ret.normalizeRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)
	return
}

// QuerySearXNG performs a search query using your local SearXNG instance
func (s *SearXNG) QuerySearXNG(query string) (ret string, err error) {
	searxngURL := fmt.Sprintf("http://localhost:3002/search?q=%s&format=json", url.QueryEscape(query))
	return s.request(searxngURL)
}

// request sends an HTTP GET request to SearXNG and parses the response
func (s *SearXNG) request(requestURL string) (ret string, err error) {
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	// Parse JSON response
	var parsedResponse SearXNGResponse
	if err = json.Unmarshal(body, &parsedResponse); err != nil {
		return "", fmt.Errorf("error parsing JSON response: %w", err)
	}

	// Format output
	if len(parsedResponse.Results) == 0 {
		return "No results found.", nil
	}

	output := "🔍 **SearXNG Search Results:**\n\n"
	for i, result := range parsedResponse.Results {
		output += fmt.Sprintf("%d. **%s**\n   🔗 %s\n   📄 %s\n\n", i+1, result.Title, result.URL, result.Content)
		if i == 4 { // Limit to 5 results
			break
		}
	}

	return output, nil
}

// NormalizeQuery ensures clean query input
func (s *SearXNG) NormalizeQuery(query string) string {
	return s.normalizeRegex.ReplaceAllString(query, "_")
}

// Grab performs a search and formats the result (aligned with YouTube's Grab function)
func (s *SearXNG) Grab(query string) (ret string, err error) {
	query = s.NormalizeQuery(query)
	return s.QuerySearXNG(query)
}
