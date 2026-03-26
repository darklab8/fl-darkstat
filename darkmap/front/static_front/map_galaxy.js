var map = document.querySelector('.panzoom');
var panzoom = Panzoom(map, {
    maxScale: 5,
    minScale: 1,
    // panOnlyWhenZoomed: false,
    // canvas: true,
    // contain: "outside",
    handleStartEvent: function (event) {
        event.preventDefault()
    },
    noBind: false,
});
map.parentElement.addEventListener('wheel', panzoom.zoomWithWheel)

var systems = document.querySelectorAll("system-");
for (let row = 0; row < systems.length; row++) {
    systems[row].addEventListener('mouseover', function () {
        let system_nickname = systems[row].attributes["nickname"].value
        let systems1 = document.querySelectorAll('connection-[data-system2-nickname="' + system_nickname + '"]')
        let systems2 = document.querySelectorAll('connection-[data-system1-nickname="' + system_nickname + '"]')
        for (let i = 0; i < systems1.length; i++) {
            systems1[i].firstElementChild.classList.add("conn_hover");
        }
        for (let i = 0; i < systems2.length; i++) {
            systems2[i].firstElementChild.classList.add("conn_hover");
        }
    });
    systems[row].addEventListener('mouseout', function () {
        let system_nickname = systems[row].attributes["nickname"].value
        let systems1 = document.querySelectorAll('connection-[data-system2-nickname="' + system_nickname + '"]')
        let systems2 = document.querySelectorAll('connection-[data-system1-nickname="' + system_nickname + '"]')
        for (let i = 0; i < systems1.length; i++) {
            systems1[i].firstElementChild.classList.remove("conn_hover");
        }
        for (let i = 0; i < systems2.length; i++) {
            systems2[i].firstElementChild.classList.remove("conn_hover");
        }
    });
}

function getOffset1(el) {
    const rect = el.getBoundingClientRect();
    return {
        x: rect.left + window.scrollX,
        y: rect.top + window.scrollY
    };
}

function getOffset2(el) {
    var _x = 0;
    var _y = 0;
    while (el && !isNaN(el.offsetLeft) && !isNaN(el.offsetTop)) {
        _x += el.offsetLeft - el.scrollLeft;
        _y += el.offsetTop - el.scrollTop;
        el = el.offsetParent;
    }
    return { top: _y, left: _x };
}

// function getCenter(el) {
//     const rect = getOffset2(el);
//     return {
//         x: rect.left,
//         y: rect.top,
//     };
// }

function getCenter(el) {
    const rect = el.getBoundingClientRect();
    return {
        x: (rect.right + rect.left) / 2.0,
        y: (rect.top + rect.bottom) / 2.0,
    };
}

function refresh_edges() {
    let edges = document.querySelectorAll("line-");

    for (let row = 0; row < edges.length; row++) {
        let sys1_nick = edges[row].attributes['data-system1-nickname'].value
        let sys2_nick = edges[row].attributes['data-system2-nickname'].value

        let system1 = document.querySelector("#system-" + sys1_nick)
        let system2 = document.querySelector("#system-" + sys2_nick)

        let p1 = getCenter(system1);
        let p2 = getCenter(system2);

        let dx = p2.x - p1.x;
        let dy = p2.y - p1.y;

        let range = Math.sqrt(dx * dx + dy * dy);

        edges[row].style.height = range + "px";
        edges[row].style.transform = "translateY(-" + range / 2.0 + "px)";

        edges[row].parentElement.style.transform = "rotate(" + (Math.atan2(dy, dx) - Math.PI / 2.0) + "rad)";
    }
}

refresh_edges();
window.addEventListener('resize', function (event) {
    refresh_edges();
});