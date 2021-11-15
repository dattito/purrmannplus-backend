package commands

import (
	"fmt"

	"github.com/dattito/purrmannplus-backend/app/models"
	"github.com/dattito/purrmannplus-backend/database"
	"github.com/dattito/purrmannplus-backend/services/moodle"
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
