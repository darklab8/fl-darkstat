function InstallSystem() {
    highlightFromQuery();
    InstallPanzoom(false);
    InstallGridLabelsPanzoom();
    InstallMenu();
    InstallLabelOverlapper();
    SetCursorPosition();
}

document.addEventListener("DOMContentLoaded", (event) => {
    InstallSystem();
});