function InstallSystem() {
    InstallGridLabelsPanzoom();
    highlightFromQuery();
    InstallPanzoom(false);
    InstallMenu();
    InstallLabelOverlapper();
    SetCursorPosition();
}

document.addEventListener("DOMContentLoaded", (event) => {
    InstallSystem();
});