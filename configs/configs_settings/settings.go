package configs_settings

import (
	"os"

	"github.com/darklab8/go-utils/utils/enverant"
	"github.com/darklab8/go-utils/utils/utils_settings"
	"github.com/darklab8/go-utils/utils/utils_types"
)

type ConfEnvVars struct {
	utils_settings.UtilsEnvs
	FreelancerFolder         utils_types.FilePath
	FreelancerFolderFailback utils_types.FilePath
	FullBasesAPIURL          *string
	Enver                    *enverant.Enverant
}

var Env ConfEnvVars

func init() {
	Env = GetEnvs()
}

func GetEnvs() ConfEnvVars {
	envs := enverant.NewEnverant(enverant.WithPrefix("CONFIGS_"), enverant.WithDescription("CONFIGS set of envs for freelancer configs parsing library"))
	Env = ConfEnvVars{
		UtilsEnvs:                utils_settings.GetEnvs(),
		FreelancerFolder:         getGameLocation(envs),
		FreelancerFolderFailback: utils_types.FilePath(envs.GetStrOr("FREELANCER_FOLDER_FAILBACK", "", enverant.WithDesc("if some configs aren't defined in first freelancer folder, grab from this one. Useful for FLSR usage in CI"))),
		FullBasesAPIURL:          envs.GetPtrStr("DISCO_BASES_FULL_URL", enverant.WithDesc("base url that has all pobs but no pob goods. useful to enchance data")),
		Enver:                    envs,
	}

	return Env
}

func getGameLocation(envs *enverant.Enverant) utils_types.FilePath {
	var folder utils_types.FilePath = utils_types.FilePath(
		envs.GetStr("FREELANCER_FOLDER", enverant.OrStr(""), enverant.WithDesc("path to Freelancer folder root for data parsing. By default grabs current workdir")),
	)

	if folder == "" {
		workdir, _ := os.Getwd()
		folder = utils_types.FilePath(workdir)
	}

	return folder
}
