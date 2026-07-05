function InstallSystem() {
    highlightFromQuery();
    InstallPanzoom(false);
    InstallGridLabelsPanzoom();
    InstallMenu();
    InstallLabelOverlapper();
}

document.addEventListener("DOMContentLoaded", (event) => {
    InstallSystem();
});