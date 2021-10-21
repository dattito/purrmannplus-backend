package hpg

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/dattito/purrmannplus-backend/config"
)

// Checks if the credentials are correct
func CheckCredentials(authId, authPw string) (bool, error) {

	data := url.Values{
		"authid": {authId},
		"authpw": {authPw},
	}

	resp, err := http.PostForm(fmt.Sprintf("%s/pmwiki/pmwiki.php?n=Main.%s", config.SUBSTITUTION_URL, authId), data)

	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return false, err
	}

	sb := string(body)

	return strings.Contains(sb, "abmelden"), nil
}
