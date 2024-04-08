package settings

import (
	"os"
	"strings"

	_ "embed"

	"github.com/darklab8/go-utils/goutils/utils"
	"github.com/darklab8/go-utils/goutils/utils/utils_filepath"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

var ProjectFolder utils_types.FilePath

var FreelancerFolder utils_types.FilePath

var ToolName = "fldarkstat"

func init() {
	ProjectFolder = utils_filepath.Dir(utils_filepath.Dir(utils.GetCurrentFolder()))

	FreelancerFolder = utils_types.FilePath(os.Getenv(strings.ToUpper(ToolName) + "_FREELANCER_FOLDER"))
	if FreelancerFolder == "" {
		workdir, _ := os.Getwd()
		FreelancerFolder = utils_types.FilePath(workdir)
	}
}

//go:embed version.txt
var version string

func GetVersion() string {
	// cleaning up version from... debugging logs used during dev env
	lines := strings.Split(version, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "v") {
			return line
		}
	}
	return version
}
