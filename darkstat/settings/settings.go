package settings

import (
	"os"
	"strings"

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
}
