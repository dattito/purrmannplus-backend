package commands

import (
	"errors"
	"fmt"

	"github.com/dattito/purrmannplus-backend/app/models"
	"github.com/dattito/purrmannplus-backend/config"
	"github.com/dattito/purrmannplus-backend/database"
	db_errors "github.com/dattito/purrmannplus-backend/database/errors"
	"github.com/dattito/purrmannplus-backend/services/moodle"
	"github.com/dattito/purrmannplus-backend/services/scheduler"
	"github.com/dattito/purrmannplus-backend/services/signal_message_sender"
	"github.com/dattito/purrmannplus-backend/utils"
	"github.com/dattito/purrmannplus-backend/utils/logging"
)

func moodleAssignmentsDifferenceAmount(mayNewAssignments []int, old_assignments []int) []int {
	var difference []int
	for _, new_assignment := range mayNewAssignments {
		found := false
		for _, old_assignment := range old_assignments {
			if new_assignment == old_assignment {
				found = true
				break
			}
		}
		if !found {
			difference = append(difference, new_assignment)
		}
	}
	return difference
}

func moodleAssignmentsToTextMessage(newAssignments []int, assignmentIdToCourseNameMap map[int]string) string {
	if len(newAssignments) == 0 {
		return "Du hast keine neuen Moodle-Aufgaben"
	}

	var courseNamesThatHaveBeenNamed []string

	var text string = "Du hast neue Moodle-Aufgaben in: \n"
	for _, assignmentId := range newAssignments {
		courseName := assignmentIdToCourseNameMap[assignmentId]
		if !utils.Contains(courseNamesThatHaveBeenNamed, courseName) {
			courseNamesThatHaveBeenNamed = append(courseNamesThatHaveBeenNamed, courseName)
			text += fmt.Sprintf("%s\n", courseName)
		}
	}
	return text
}

// Returns error produced by user; error not produced by user
func AddAccountToMoodleAssignmentUpdater(accountId string) (error, error) {
	ai, err := database.DB.GetAccountInfo(accountId)
	if err != nil {
		if errors.Is(err, &db_errors.ErrRecordNotFound) {
			return errors.New("phone number has to be added first"), nil
		}
		return nil, err
	}

	if ai.PhoneNumber == "" {
		return errors.New("phone number has to be added first"), nil
	}

	if _, err := database.DB.GetMoodleAssignments(accountId); err != nil || !errors.Is(err, &db_errors.ErrRecordNotFound) {
		return errors.New("moodle assignments already exist"), nil
	}

	_, err = database.DB.SetMoodleAssignments(accountId, []int{}, true)
	if err != nil {
		return nil, err
	}

	return nil, UpdateMoodleAssignmentsByAccountId(accountId)
}

func RemoveAccountFromMoodleAssignmentUpdater(accountId string) error {
	return database.DB.RemoveAccountFromMoodleAssignmentUpdater(accountId)
}

func UpdateMoodleAssignments(m models.MoodleAssignmentUpdateInfos) error {
	logging.Debugf("Updating moodle assignments of account %s (id: %s)", m.AuthId, m.AccountId)

	rawAssignments, err := moodle.GetRawAssignmentsByCredentials(m.AuthId, m.AuthPw)
	if err != nil {
		return err
	}
	mayNewAssignments := moodle.GetAssignmentIDs(rawAssignments)

	old_assignments := m.AssignmentIds

	newAssignments := moodleAssignmentsDifferenceAmount(mayNewAssignments, old_assignments)

	// If there are no new assignments, we don't need to do anything
	if len(newAssignments) == 0 {
		return nil
	}

	_, err = database.DB.SetMoodleAssignments(m.AccountId, newAssignments, false)
	if err != nil {
		return err
	}

	logging.Debugf("Successfully updated moodle assignments of %s", m.AuthId)

	if m.NotSetYet || len(m.AssignmentIds) == 0 {
		return nil
	}

	// Send a message to the user if there are new assignments
	return signal_message_sender.SignalMessageSender.Send(moodleAssignmentsToTextMessage(newAssignments, moodle.GetAssignmentIdToCourseNameMap(rawAssignments)), m.PhoneNumber)
}

func UpdateMoodleAssignmentsByAccountId(accountId string) error {
	mdb, err := database.DB.GetAccountCredentialsAndPhoneNumberAndSMoodleAssignments(accountId)
	if err != nil {
		return err
	}
	m := models.AccountCredentialsAndPhoneNumberAndMoodleUserAssignmentsDBModelToMoodleAssignmentUpdateInfos(&mdb)

	return UpdateMoodleAssignments(m)
}

func UpdateAllMoodleAssignments() error {
	mdbs, err := database.DB.GetAllAccountCredentialsAndPhoneNumberAndSMoodleAssignments()
	if err != nil {
		return err
	}

	errCount := 0

	for _, mdb := range mdbs {
		m := models.AccountCredentialsAndPhoneNumberAndMoodleUserAssignmentsDBModelToMoodleAssignmentUpdateInfos(&mdb)
		err = UpdateMoodleAssignments(m)
		if err != nil {
			logging.Errorf("Error while updating moodle assignments of %s: %s", m.AuthId, err.Error())
			errCount++
			if errCount > config.MAX_ERROS_TO_STOP_UPDATING_MOODLE_ASSIGNMENTS {
				return errors.New("got too many errors while updating moodle assignments, stopping")
			}
		}
	}

	return nil
}

func EnableMoodleAssignmentUpdater() {
	scheduler.AddJob(config.MOODLE_UPDATECRON, func() {
		if err := UpdateAllMoodleAssignments(); err != nil {
			logging.Errorf("Error while updating moodle assignments: %s", err.Error())
		}
	})
}