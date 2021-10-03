package app

import (
	"github.com/dattito/purrmannplus-backend/app/commands"
	"github.com/dattito/purrmannplus-backend/config"
)

func Init() {
	if config.ENABLE_SUBSTITUTIONS_SCHEDULER {
		commands.EnableSubstitutionUpdater()
	}
}
