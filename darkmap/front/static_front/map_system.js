function InstallSystem() {
    highlightFromQuery();
    InstallPanzoom(false);
    InstallMenu();
    InstallLabelOverlapper();
}

document.addEventListener("DOMContentLoaded", (event) => {
    InstallSystem();
});