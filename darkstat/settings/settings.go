package settings

import (
	"fmt"
	"os"
	"strings"

	_ "embed"

	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

var FreelancerFolder utils_types.FilePath

var ToolName = "fldarkstat"

func init() {
	FreelancerFolder = utils_types.FilePath(os.Getenv(strings.ToUpper(ToolName) + "_FREELANCER_FOLDER"))
	if FreelancerFolder == "" {
		workdir, _ := os.Getwd()
		FreelancerFolder = utils_types.FilePath(workdir)
	}

	fmt.Println("FreelancerFolder=", FreelancerFolder)
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
