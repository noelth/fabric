package searxng

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/noelth/fabric/plugins"
)

// SearXNGSearchResult represents a single search result from SearXNG
type SearXNGSearchResult struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Content string `json:"content"`
}

// SearXNGResponse represents the full API response
type SearXNGResponse struct {
	Query   string                `json:"query"`
	Results []SearXNGSearchResult `json:"results"`
}

// Client implements the Fabric Tool interface
type Client struct {
	*plugins.PluginBase
}

// Ensure Client implements the Tool interface
var _ plugins.Tool = (*Client)(nil)

// NewClient initializes the tool
func NewClient() (ret *Client) {
	label := "SearXNG"

	ret = &Client{
		PluginBase: &plugins.PluginBase{
			Name:             label,
			SetupDescription: "SearXNG Service - to query and retrieve search results",
			EnvNamePrefix:    plugins.BuildEnvVariablePrefix(label),
		},
	}

	return
}

// Run executes the tool (Fabric calls this function when invoking the tool)
func (sc *Client) Run(input string) (ret string, err error) {
	searxngURL := fmt.Sprintf("http://localhost:8080/search?q=%s&format=json", url.QueryEscape(input))
	return sc.request(searxngURL)
}

func (sc *Client) request(requestURL string) (ret string, err error) {
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

// Register the tool with Fabric
func init() {
	plugins.RegisterTool(NewClient())
}
