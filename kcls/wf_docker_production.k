import infra.kcls.models.github_workflow as gw
import infra.kcls.common.github as ghc
import infra.kcls.common.docker
import .common as c

_workflow_name = "Deploy docker Production"
_workflow = gw.Workflow {
    _filename = "docker-production.yml"
    name = _workflow_name
    on = {
        workflow_dispatch: {}
        push.tags: ["*"]
    }
    env = ghc.DockerHostEnv
    jobs = {
        build = gw.JobDispatch {
            name = "Build docker image and save to latest"
            uses = "darklab8/fl-darkstat/.github/workflows/deploy-docker-staging.yml@master"
        }
        job = {
            name = _workflow_name
            needs = ["build"]
            steps = [
                ghc.CheckoutRepo
                ghc.GetVersionFromTag
                ghc.DockerLogin
                ghc.InstallDarklabSshKey
                docker.Retag(docker.RetagOpts {
                    image_name = c._image_name
                    old_tag = "staging"
                    new_tag = ghc.GetVersionOutput
                    pull = True
                    push = True
                })
                docker.Retag(docker.RetagOpts {
                    image_name = c._image_name
                    old_tag = ghc.GetVersionOutput
                    new_tag = "production"
                    push = True
                })
                docker.DeploySwarm(docker.DeploySwarmOpts {
                    image_name = c._image_name
                    service_name = "darkstat-production"
                    tag = "production"
                })
            ]
        }
    }
}
