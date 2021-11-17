package moodle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"

	"github.com/dattito/purrmannplus-backend/app/models"
	"github.com/dattito/purrmannplus-backend/config"
	"github.com/dattito/purrmannplus-backend/utils/logging"
)

type reponse struct {
	Token        string `json:"token"`
	ErrorMessage string `json:"error"`
	ErrorCode    string `json:"errorcode"`
}

func GetToken(username, password string) (string, error) {
	if config.MOODLE_URL == "" {
		return "", fmt.Errorf("moodle URL not set")
	}

	if username == "" || password == "" {
		return "", nil
	}

	resp, err := http.PostForm(fmt.Sprintf("%s/login/token.php", config.MOODLE_URL),
		url.Values{
			"username": {username},
			"password": {password},
			"service":  {"moodle_mobile_app"},
		})
	if err != nil {
		logging.Errorf("Error while checking moodle credentials: %s", err)
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logging.Errorf("Error while checking moodle credentials: %s", err)
		return "", err
	}

	var r reponse
	if err := json.Unmarshal(body, &r); err != nil {
		logging.Errorf("Error while checking moodle credentials: %s", err)
		return "", err
	}

	if r.ErrorCode == "invalidlogin" {
		return "", nil
	} else if r.ErrorCode != "" {
		logging.Errorf("Error while checking moodle credentials: %s", r.ErrorMessage)
		return "", fmt.Errorf("error while checking moodle credentials: %s", r.ErrorMessage)
	}

	if r.Token == "" {
		logging.Errorf("Error while checking moodle credentials: token is empty")
		return "", fmt.Errorf("error while checking moodle credentials: token is empty")
	}

	return r.Token, nil
}

// Checks if the credentials are correct, should be the same as substitutions.CheckCredentials()
func CheckCredentials(username, password string) (bool, error) {
	token, err := GetToken(username, password)
	if err != nil {
		return false, err
	}

	return token != "", nil
}

func GetRawAssignments(token string) (models.MoodleCourse, error) {
	if config.MOODLE_URL == "" {
		return models.MoodleCourse{}, fmt.Errorf("moodle URL not set")
	}

	if token == "" {
		return models.MoodleCourse{}, fmt.Errorf("token is empty")
	}

	resp, err := http.Get(fmt.Sprintf("%s/webservice/rest/server.php?wstoken=%s&wsfunction=mod_assign_get_assignments&moodlewsrestformat=json", config.MOODLE_URL, token))
	if err != nil {
		logging.Errorf("Error while getting moodle assignments: %s", err)
		return models.MoodleCourse{}, err
	}

	defer resp.Body.Close()

	var r models.MoodleCourse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return models.MoodleCourse{}, err
	}

	return r, nil
}

func GetRawAssignmentsByCredentials(username, password string) (models.MoodleCourse, error) {
	token, err := GetToken(username, password)
	if err != nil {
		return models.MoodleCourse{}, err
	}

	return GetRawAssignments(token)
}

func GetAssignmentIDs(assingments models.MoodleCourse) []int {
	var ids []int
	for _, course := range assingments.Courses {
		for _, assignment := range course.Assignments {
			ids = append(ids, assignment.ID)
		}
	}
	sort.Ints(ids)
	return ids
}

func GetAssignmentIDsByCredentials(username, password string) ([]int, error) {
	token, err := GetToken(username, password)
	if err != nil {
		return []int{}, err
	}

	assingments, err := GetRawAssignments(token)
	if err != nil {
		return []int{}, err
	}

	return GetAssignmentIDs(assingments), nil
}

func GetAssignmentIdToCourseNameMap(assingments models.MoodleCourse) map[int]string {
	var idsToCourses map[int]string = make(map[int]string)
	for _, course := range assingments.Courses {
		for _, assignment := range course.Assignments {
			idsToCourses[assignment.ID] = course.FullName
		}
	}
	return idsToCourses
}
