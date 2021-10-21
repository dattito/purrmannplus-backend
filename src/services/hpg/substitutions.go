package hpg

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dattito/purrmannplus-backend/utils"
)

var weekdays = [5]string{"Mo", "Di", "Mi", "Do", "Fr"}

// Returns true if the given string begins with a weekday
func beginsWithAWeekday(s string) bool {
	b := strings.Fields(s)[0]
	for _, weekday := range weekdays {
		if b == weekday {
			return true
		}
	}
	return false
}

func GetSubstituationOfStudent(authid, authpw string) (map[string][]string, error) {

	// Request the HTML page.
	res, err := http.PostForm(fmt.Sprintf("https://vertretungsplan.hpg-speyer.de/pmwiki/pmwiki.php?n=Main.%s", strings.ToLower(authid)),
		url.Values{
			"authid": {authid},
			"authpw": {authpw},
		},
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer res.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	if !strings.Contains(doc.Text(), "abmelden") {
		return nil, WrongCredentialsError
	}

	// Find the review items
	s := doc.Find("table") // if s.Length()=4, there are new substituations

	if s.Length() < 4 {
		return map[string][]string{}, nil
	}

	sp := s.Eq(1)

	spMap := map[string][]string{}

	lastW := ""
	var spErr error
	sp.Find("tr").Each(func(i int, s *goquery.Selection) {
		if spErr != nil {
			return
		}

		outp := ""
		s.Find("td").Each(func(j int, k *goquery.Selection) {
			t := strings.TrimSpace(k.Text())

			if beginsWithAWeekday(t) {
				lastW = t
			} else {
				outp += t + " "
			}

		})
		outpt, err := utils.ConvertStringToLatin1(strings.TrimSuffix(outp, " "))
		if err != nil {
			spErr = err
			return
		}
		if outpt != "" {
			spMap[lastW] = append(spMap[lastW], outpt)
		}
	})

	if spErr != nil {
		return map[string][]string{}, spErr
	}

	return spMap, nil
}
