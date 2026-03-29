function InstallEdgeHighlight() {
    var systems = document.querySelectorAll("system-");
    for (let row = 0; row < systems.length; row++) {
        systems[row].addEventListener('mouseover', function () {
            let system_nickname = systems[row].attributes["nickname"].value;
            let lines1 = document.querySelectorAll('line[data-system2-nickname="' + system_nickname + '"]');
            let lines2 = document.querySelectorAll('line[data-system1-nickname="' + system_nickname + '"]');
            for (let i = 0; i < lines1.length; i++) {
                lines1[i].classList.add("conn_hover");
            }
            for (let i = 0; i < lines2.length; i++) {
                lines2[i].classList.add("conn_hover");
            }
        });
        systems[row].addEventListener('mouseout', function () {
            let system_nickname = systems[row].attributes["nickname"].value;
            let lines1 = document.querySelectorAll('line[data-system2-nickname="' + system_nickname + '"]');
            let lines2 = document.querySelectorAll('line[data-system1-nickname="' + system_nickname + '"]');
            for (let i = 0; i < lines1.length; i++) {
                lines1[i].classList.remove("conn_hover");
            }
            for (let i = 0; i < lines2.length; i++) {
                lines2[i].classList.remove("conn_hover");
            }
        });
    }
}

function InstallGalaxy() {
    InstallLabelOverlapper();
    InstallMenu();
    InstallEdgeHighlight();
    InstallPanzoom();
}
