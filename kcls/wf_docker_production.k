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
        # i think it is messy approach
        # build = gw.JobDispatch {
        #     name = "Build docker image and save to latest"
        #     uses = "darklab8/fl-darkstat/.github/workflows/docker-staging.yml@master"
        # }
        job = _production_deploy_job
    }
}
_production_deploy_job: gw.Job = {
    name = _workflow_name
    # needs = ["build"]
    steps = [
        ghc.CheckoutRepo
        ghc.DockerLogin
        ghc.InstallDarklabSshKey
        ghc.GetVersionFromTag
        docker.Build(docker.BuildOpts {
            image_name = c._image_name
            tag = ghc.GetVersionOutput
            extra_tags = ["production"]
            push = True
        })
        docker.DeploySwarm(docker.DeploySwarmOpts {
            image_name = c._image_name
            tag = "production"
            service_name = "darkstat-production"
        })
    ]
}
