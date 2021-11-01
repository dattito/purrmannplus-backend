package moodle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/dattito/purrmannplus-backend/config"
	"github.com/dattito/purrmannplus-backend/utils/logging"
)

type reponse struct {
	Token        string `json:"token"`
	ErrorMessage string `json:"error"`
	ErrorCode    string `json:"errorcode"`
}

// Checks if the credentials are correct, should be the same as substitutions.CheckCredentials()
func CheckCredentials(username, password string) (bool, error) {
	if config.MOODLE_URL == "" {
		return false, fmt.Errorf("moodle URL not set")
	}

	if username == "" || password == "" {
		return false, nil
	}

	resp, err := http.PostForm(fmt.Sprintf("%s/login/token.php", config.MOODLE_URL),
		url.Values{
			"username": {username},
			"password": {password},
			"service":  {"moodle_mobile_app"},
		})
	if err != nil {
		logging.Errorf("Error while checking moodle credentials: %s", err)
		return false, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logging.Errorf("Error while checking moodle credentials: %s", err)
		return false, err
	}

	var r reponse
	if err := json.Unmarshal(body, &r); err != nil {
		logging.Errorf("Error while checking moodle credentials: %s", err)
		return false, err
	}

	if r.ErrorCode == "invalidlogin" {
		return false, nil
	} else if r.ErrorCode != "" {
		logging.Errorf("Error while checking moodle credentials: %s", r.ErrorMessage)
		return false, fmt.Errorf("error while checking moodle credentials: %s", r.ErrorMessage)
	}

	if r.Token == "" {
		logging.Errorf("Error while checking moodle credentials: token is empty")
		return false, fmt.Errorf("error while checking moodle credentials: token is empty")
	}

	return true, nil
}
