package hpg

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var weekdays = [5]string{"Mo", "Di", "Mi", "Do", "Fr"}

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

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	log.Println(sb)

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("wrong credentials")
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	// Checked for wrong crendentials
	if doc.Find("form").Length() > 0 {
		return nil, fmt.Errorf("wrong credentials")
	}

	// Find the review items
	s := doc.Find("table") // if s.Length()=3, there are new substituations

	if s.Length() < 3 {
		return make(map[string][]string), nil
	}

	sp := s.Eq(1)

	spMap := map[string][]string{}

	lastW := ""
	sp.Find("tr").Each(func(i int, s *goquery.Selection) {
		outp := ""
		s.Find("td").Each(func(j int, k *goquery.Selection) {
			t := strings.TrimSpace(k.Text())

			if beginsWithAWeekday(t) {
				lastW = t
			} else {
				outp += t + " "
			}

		})
		outpt := strings.TrimSuffix(outp, " ")
		if outpt != "" {
			spMap[lastW] = append(spMap[lastW], outpt)
		}
	})

	return spMap, nil
}
