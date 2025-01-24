import infra.kcls.models.github_workflow as gw
import infra.kcls.common.github as ghc
import infra.kcls.actions.checkout_freelancer as cf
import infra.kcls.actions.darkstat

import .common as c

_workflow = gw.Workflow {
    _filename = "test-build.yml"
    name = "Test build"
    on = {
        push.branches = ["master"]
        schedule: [{cron = "0 9 * * *"}]
        pull_request: {}
    }
    jobs = {
        job = _test_job
    }
}
_test_job: gw.Job = {
    name = "Test build"
    steps = [
        ghc.CheckoutRepo
        cf.NewCheckoutFreelancer(cf.CheckoutArgs {
            freelancer_mod = cf._mod_discovery
            freelancer_folder = c._freelancer_folder
        })
    ] + c.DarkstatBuildAssetsSteps(darkstat.BuildArgs {
        site_host = ""
        site_root = "/fl-darkstat/"
        freelancer_folder = c._freelancer_folder
        heading = '<span style="font-weight:1000;">DEV ENV</span> <a href="https://github.com/darklab8/fl-darkstat">fl-darkstat</a> for <a href="https://github.com/darklab8/fl-data-discovery">Freelancer Discovery</a>'
        relay_host = "https://darkrelay-staging.dd84ai.com"
    })
}
