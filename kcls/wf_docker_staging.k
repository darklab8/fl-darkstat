import infra.kcls.models.github_workflow as gw
import infra.kcls.common.github as ghc
import infra.kcls.common.docker
import infra.kcls.actions.install_go as go
import infra.kcls.actions.install_taskfile as tf
import infra.kcls.actions.install_autogit as autogit
import .common as c

_workflow_name = "Deploy docker staging"
_workflow = gw.Workflow {
    _filename = "docker-staging.yml"
    name = _workflow_name
    on = {
        push.tags: ghc.TagsSemverStaging
        workflow_dispatch: {}
    }
    env = ghc.DockerHostEnv
    jobs = {
        job = _docker_staging_job
    }
}
_docker_staging_job: gw.Job = {
    name = _workflow_name
    steps = [
        ghc.CheckoutRepoWithCommits
        go.InstallGo {}
        tf.InstallTaskfile {}
        autogit.InstallAutogit {}
        ghc.GetVersionFromTag
        ghc.DockerLogin
        ghc.InstallDarklabSshKey
        docker.Build(docker.BuildOpts {
            image_name = c._image_name
            tag = ghc.GetVersionOutput
            extra_tags = ["staging"]
            push = True
        })
        docker.DeploySwarm(docker.DeploySwarmOpts {
            image_name = c._image_name
            tag = "staging"
            service_name = "darkstat-staging"
        })
    ]
}
