function getCenter(el) {
    const rect = el.getBoundingClientRect();
    return {
        x: (rect.right + rect.left) / 2.0,
        y: (rect.top + rect.bottom) / 2.0,
    };
}

function RefreshEdgePositions() {
    // let edges = document.querySelectorAll("line-");

    // for (let row = 0; row < edges.length; row++) {
    //     let sys1_nick = edges[row].attributes['data-system1-nickname'].value
    //     let sys2_nick = edges[row].attributes['data-system2-nickname'].value

    //     let system1 = document.querySelector("#system-" + sys1_nick)
    //     let system2 = document.querySelector("#system-" + sys2_nick)

    //     let p1 = getCenter(system1);
    //     let p2 = getCenter(system2);

    //     let dx = p2.x - p1.x;
    //     let dy = p2.y - p1.y;

    //     let range = Math.sqrt(dx * dx + dy * dy);

    //     edges[row].style.height = range + "px";
    //     edges[row].style.transform = "translateY(-" + range / 2.0 + "px)";

    //     edges[row].parentElement.style.transform = "rotate(" + (Math.atan2(dy, dx) - Math.PI / 2.0) + "rad)";
    // }
}

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
    InstallMenu();
    InstallEdgeHighlight();
    InstallLabelOverlapper();
    window.addEventListener('resize', function (event) {
        RefreshEdgePositions();
    });
    InstallPanzoom();
}
