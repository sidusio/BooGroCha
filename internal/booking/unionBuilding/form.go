package unionBuilding

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

func loginForm(client *http.Client) (url.Values, error) {
	resp, err := client.Get(loginURL)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	inputs := doc.Find("body > form > input")

	values := url.Values{}

	/*
	  Traverses through each input field and extracts the name and values
	  which is used as parameters in the login call.
	*/
	inputs.Each(func(i int, s *goquery.Selection) {
		name, exists := s.Attr("name")
		if exists {
			value, exists := s.Attr("value")
			if exists {
				values.Add(name, value)
			}
		}
	})

	for _, key := range []string{
		viewStateKey,
		viewStateGeneratorKey,
		eventValidationKey,
	} {
		if _, ok := values[key]; !ok {
			return nil, errors.New(fmt.Sprintf("value not found in form: %s", key))
		}
	}
	values.Add(loginButtonKey, "Enter")

	return values, nil
}
