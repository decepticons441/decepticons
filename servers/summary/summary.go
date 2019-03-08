package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

//PreviewImage represents a preview image for a page
type PreviewImage struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secureURL,omitempty"`
	Type      string `json:"type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Alt       string `json:"alt,omitempty"`
}

//PreviewVideo represents a preview image for a page
type PreviewVideo struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secureURL,omitempty"`
	Type      string `json:"type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
}

//PageSummary represents summary properties for a web page
type PageSummary struct {
	Type        string          `json:"type,omitempty"`
	URL         string          `json:"url,omitempty"`
	Title       string          `json:"title,omitempty"`
	SiteName    string          `json:"siteName,omitempty"`
	Description string          `json:"description,omitempty"`
	Author      string          `json:"author,omitempty"`
	Keywords    []string        `json:"keywords,omitempty"`
	Icon        *PreviewImage   `json:"icon,omitempty"`
	Images      []*PreviewImage `json:"images,omitempty"`
	Videos      []*PreviewVideo `json:"videos,omitempty"`
}

//SummaryHandler handles requests for the page summary API.
//This API expects one query string parameter named `url`,
//which should contain a URL to a web page. It responds with
//a JSON-encoded PageSummary struct containing the page summary
//meta-data.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	/*TODO: add code and additional functions to do the following:
	Helpful Links:
	https://golang.org/pkg/net/http/#Request.FormValue
	https://golang.org/pkg/net/http/#Error
	https://golang.org/pkg/encoding/json/#NewEncoder
	*/

	// Add an HTTP header to the response with the name
	// `Access-Control-Allow-Origin` and a value of `*`. This will
	// allow cross-origin AJAX requests to your server.
	// w.Header().Add("Access-Control-Allow-Origin", "*")

	// Get the `url` query string parameter value from the request.
	// If not supplied, respond with an http.StatusBadRequest error.
	url := r.FormValue("url")
	if len(url) == 0 {
		http.Error(w, `Invalid Input Params!`, http.StatusBadRequest)
		return
	}

	// Call fetchHTML() to fetch the requested URL. See comments in that
	// function for more details.
	html, err := fetchHTML(url)
	if err != nil {
		http.Error(w, "Could not fetch html", http.StatusBadRequest)
	}

	// Call extractSummary() to extract the page summary meta-data,
	// as directed in the assignment. See comments in that function
	// for more details
	summary, err := extractSummary(url, html)
	w.Header().Set("Content-Type", "application/json")

	// if err is not EOF,  then do http.Error
	// check end of file after extractSummary
	if err != nil && err == io.EOF {
		http.Error(w, `Reached end of file`, 406)
	}

	// Close the response HTML stream so that you don't leak resources.
	defer html.Close()

	// Finally, respond with a JSON-encoded version of the PageSummary
	// struct. That way the client can easily parse the JSON back into
	// an object. Remember to tell the client that the response content
	// type is JSON.
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(summary); err != nil {
		fmt.Printf("error encoding struct into JSON %v\n", err)
	}
}

//fetchHTML fetches `pageURL` and returns the body stream or an error.
//Errors are returned if the response status code is an error (>=400),
//or if the content type indicates the URL is not an HTML page.
func fetchHTML(pageURL string) (io.ReadCloser, error) {
	/*TODO:
	To test your implementation of this function, run the TestFetchHTML
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestFetchHTML

	Helpful Links:
	https://golang.org/pkg/net/http/#Get
	*/

	// Do an HTTP GET for the page URL. If the response status
	// code is >= 400, return a nil stream and an error. If the response
	// content type does not indicate that the content is a web page, return
	// a nil stream and an error. Otherwise return the response body and
	// no (nil) error.
	response, err := http.Get(pageURL)

	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("Response Status Code is greater or equal to 400: %d", err)
	}

	if !strings.HasPrefix(response.Header.Get("Content-Type"), "text/html") {
		return nil, fmt.Errorf("Content Type is not a Web Page: %d", err)
	}

	return response.Body, nil
}

//extractSummary tokenizes the `htmlStream` and populates a PageSummary
//struct with the page's summary meta-data.
func extractSummary(pageURL string, htmlStream io.ReadCloser) (*PageSummary, error) {
	/*TODO: tokenize the `htmlStream` and extract the page summary meta-data
	according to the assignment description.

	To test your implementation of this function, run the TestExtractSummary
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestExtractSummary

	Helpful Links:
	https://drstearns.github.io/tutorials/tokenizing/
	http://ogp.me/
	https://developers.facebook.com/docs/reference/opengraph/
	https://golang.org/pkg/net/url/#URL.ResolveReference
	*/

	// Tokenize the `htmlStream` and extract the page summary meta-data
	// according to the assignment description.
	summary := &PageSummary{}
	tokenizer := html.NewTokenizer(htmlStream)

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			return summary, err
		}

		// Only looking at the header for summary
		if tokenType == html.EndTagToken {
			if tokenizer.Token().Data == "head" {
				break
			}
		} else if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			token := tokenizer.Token()
			data := token.Data

			if data == "link" {
				icon, err := makeIcon(pageURL, token)
				summary.Icon = icon
				if err != nil {
					return summary, fmt.Errorf("Could no parse icon's href url value: %v", err)
				}
			}

			if data == "meta" {
				err := checkMetaTag(token, summary, pageURL)
				if err != nil {
					return nil, err
				}
			}

			if data == "title" {
				tokenType2 := tokenizer.Next()
				token2 := tokenizer.Token()
				data2 := token2.Data

				if summary.Title == "" && tokenType2 == html.TextToken {
					summary.Title = data2
				}
			}
		}
	}
	return summary, nil
}

func checkMetaTag(token html.Token, summary *PageSummary, pageURL string) (error) {
	property := ""
	content := ""
	for _, attr := range token.Attr {
		if attr.Key == "property" || attr.Key == "name" {
			property = attr.Val
		}
		if attr.Key == "content" {
			content = attr.Val
		}
	}

	if property == "og:type" || property == "twitter:card" {
		// fmt.Printf(content)
		summary.Type = content
	} else if property == "og:url" {
		summary.URL = content
	} else if property == "og:title" || property == "twitter:title" {
		summary.Title = content
	} else if property == "og:site_name" {
		summary.SiteName = content
	} else if property == "og:description" || property == "twitter:description" {
		summary.Description = content
	} else if summary.Description == "" && property == "description" {
		summary.Description = content
	} else if property == "author" {
		summary.Author = content
	} else if property == "keywords" {
		makeKeywordSlice(content, summary)
	} else if strings.Contains(property, "og:image") || strings.Contains(property, "twitter:image") {
		if property == "og:image" || property == "twitter:image" {
			preview := &PreviewImage{}

			url, err := url.Parse(content)
			if err != nil {
				fmt.Printf("Error parsing content image url: %d", err)
				return err
			}
			base, err := url.Parse(pageURL)
			if err != nil {
				fmt.Printf("Error parsing image pageURL: %d", err)
				return err
			}

			finalURL := base.ResolveReference(url).String()
			preview.URL = finalURL
			summary.Images = append(summary.Images, preview)
		}
		makeImage(summary.Images[len(summary.Images)-1], content, property)
	} else if strings.Contains(property, "og:video") {
		if property == "og:video" {
			preview := &PreviewVideo{}

			url, err := url.Parse(content)
			if err != nil {
				fmt.Printf("Error parsing content video url: %d", err)
				return err
			}
			base, err := url.Parse(pageURL)
			if err != nil {
				fmt.Printf("Error parsing video pageURL: %d", err)
				return err
			}

			finalURL := base.ResolveReference(url).String()
			preview.URL = finalURL
			summary.Videos = append(summary.Videos, preview)
		}
		makeVideo(summary.Videos[len(summary.Videos)-1], content, property)
	}
	return nil
}

// Makes array of keywords and sets the summary's keyword field
func makeKeywordSlice(keywordStream string, summary *PageSummary) {
	keywordSplit := strings.Split(keywordStream, ",")
	keywordSlice := []string{}

	for i := range keywordSplit {
		keywordSlice = append(keywordSlice, strings.TrimSpace(keywordSplit[i]))
	}

	summary.Keywords = keywordSlice
}

// Sets some preview image fields for images that don't contain a url
func makeImage(preview *PreviewImage, content string, property string) {
	if property == "og:image:secure_url" {
		preview.SecureURL = content
	} else if property == "og:image:type" {
		preview.Type = content
	} else if property == "og:image:width" {
		width, err := strconv.Atoi(content)
		if err != nil {
			fmt.Errorf("Cannot convert Image's Width to an Integer: %d", err)
		}
		preview.Width = width
	} else if property == "og:image:height" {
		height, err := strconv.Atoi(content)
		if err != nil {
			fmt.Errorf("Cannot convert Image's Height to an Integer: %d", err)
		}
		preview.Height = height
	} else if property == "og:image:alt" || property == "twitter:image:alt" {
		preview.Alt = content
	}
}

// Sets some preview video fields for videos that don't contain a url
func makeVideo(preview *PreviewVideo, content string, property string) {
	if property == "og:video:secure_url" {
		preview.SecureURL = content
	} else if property == "og:video:type" {
		preview.Type = content
	} else if property == "og:video:width" {
		width, err := strconv.Atoi(content)
		if err != nil {
			fmt.Errorf("Cannot convert Video's Width to an Integer: %d", err)
		}
		preview.Width = width
	} else if property == "og:video:height" {
		height, err := strconv.Atoi(content)
		if err != nil {
			fmt.Errorf("Cannot convert Video's Height to an Integer: %d", err)
		}
		preview.Height = height
	}
}

// Sets icon fields for preview image
func makeIcon(pageURL string, token html.Token) (*PreviewImage, error) {
	preview := &PreviewImage{}
	for _, attr := range token.Attr {
		if attr.Key == "href" {
			url, err := url.Parse(attr.Val)
			if err != nil {
				return preview, fmt.Errorf("Could no parse icon's href url value: %v", err)
			}
			base, err := url.Parse(pageURL)
			if err != nil {
				return preview, fmt.Errorf("Could no parse icon's href base value: %v", err)
			}
			finalURL := base.ResolveReference(url).String()
			preview.URL = finalURL
		}
		if attr.Key == "type" {
			preview.Type = attr.Val
		}
		if attr.Key == "sizes" {
			sizes := []string{}
			if strings.Contains(attr.Val, "x") {
				sizes = strings.Split(attr.Val, "x")
			} else {
				sizes = strings.Split(attr.Val, "X")
			}

			if len(sizes) == 2 {
				height, err := strconv.Atoi(strings.TrimSpace(sizes[0]))
				width, err := strconv.Atoi(strings.TrimSpace(sizes[1]))
				if err != nil {
					return preview, fmt.Errorf("Could not convert height or width's string values into ints: %v", err)
				}
				preview.Height = height
				preview.Width = width
			}
		}
	}
	return preview, nil
}
