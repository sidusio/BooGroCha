package timeedit

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type form struct {
	Values url.Values
	Action string
}

func getForm(resp *http.Response, selector string) (form, error) {
	// Extract the form element
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return form{}, err
	}
	formElement := doc.Find(fmt.Sprintf("body %s", selector))

	// Get the action url for submitting the form
	formAction, found := formElement.Attr("action")
	if !found {
		return form{}, err
	}

	// Build url for relative paths
	if formAction[0] == '/' {
		formAction = fmt.Sprintf("%s://%s%s", resp.Request.URL.Scheme, resp.Request.URL.Host, formAction)
	}

	// Extract all in the form values
	formData := url.Values{}
	inputs := formElement.Find("input")
	for i := range inputs.Nodes {
		input := inputs.Eq(i)
		name, foundName := input.Attr("name")
		value, foundValue := input.Attr("value")
		if foundName && foundValue && value != "" {
			formData.Add(name, value)
		}
	}

	return form{
		Values: formData,
		Action: formAction,
	}, nil
}

func (f *form) Post(client *http.Client) (*http.Response, error) {
	return client.PostForm(f.Action, f.Values)
}
