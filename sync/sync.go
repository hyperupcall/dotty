package sync

import (
	"github.com/eankeen/globe/config"
	"github.com/eankeen/globe/internal/util"
	"github.com/eankeen/globe/validate"
)

// Sync project with all bootstrap files
func Sync(validatedArgs validate.ValidatedArgs) {
	project := config.GetConfig(validatedArgs.StoreDir)
	util.PrintInfo("Project located at %s\n", project.ProjectLocation)

	for _, file := range project.BootstrapFiles.Files {
		util.PrintInfo("Processing file %s\n", file.RelPath)

		if file.Op == "add" {
			copyFile(project, file)
			continue
		} else if file.Op == "remove" {
			removeFile(project, file)
			continue
		}

		util.PrintError("File '%s's operation could not be read. Exiting.\n", file.RelPath)
	}
}
