import infra.kcls.models.github_workflow as gw
import infra.kcls.common.github as ghc
import infra.kcls.actions.checkout_freelancer as cf
import infra.kcls.actions.darkstat
import infra.kcls.models.github_action as ga

import .common as c

_workflow = gw.Workflow {
    _filename = "pages_staging.yml"
    name = "Deploy pages staging"
    on = {
        push.tags: ghc.TagsSemverStaging
        workflow_dispatch: {}
    }
    jobs = {
        job = _pages_staging_job
    }
}
_pages_staging_job: gw.Job = {
    name = "Github Pages Staging"
    permissions: ghc.PagesPermissions {}
    environment: dict(ghc.PagesEnvironment {})
    steps = [
        ghc.CheckoutRepo
        ga.Step {
            name: "Upload iframe redirect"
            run: """set -ex
mkdir build | true
cp docker/iframe-staging.html build/index.html
"""
        }
        ga.Step {
            name: "Upload artifacts for deployment to Github Pages"
            uses: "actions/upload-pages-artifact@v3"
            with: {
                name: "github-pages"
                path: "./build"
            }
        }
        ga.Step {
            name: "Deploy to pages"
            uses: "actions/deploy-pages@v4"
            id: "deployment"
        }
    ]
}
