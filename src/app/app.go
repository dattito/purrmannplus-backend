package app

import (
	"github.com/datti-to/purrmannplus-backend/app/commands"
	"github.com/datti-to/purrmannplus-backend/config"
)

func Init() {
	if config.ENABLE_SUBSTITUTIONS_SCHEDULER {
		commands.EnableSubstitutionUpdater()
	}
}
