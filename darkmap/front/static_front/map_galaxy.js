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
    noBind: true,
});
map.parentElement.addEventListener('wheel', panzoom.zoomWithWheel)

map.addEventListener('panzoomchange', function (event) {
    if (event.detail.scale == 1.0 && event.detail.x != 0 && event.detail.y != 0) {
        console.log("event", event)
        panzoom.reset({ scale: event.detail.scale })
    }
});



function toggle_option(checked, checkbox_name, hidden_class, unhidden_class) {
    switch (checked) {
        case true:
            console.log(this.value, "turned on")

            let hidden_systems = document.querySelectorAll("." + hidden_class);

            for (let row = 0; row < hidden_systems.length; row++) {
                hidden_systems[row].classList.remove(hidden_class);
                hidden_systems[row].classList.add(unhidden_class);
            }
            sessionStorage.setItem(checkbox_name, "true");
            break;
        case false:
            console.log(this.value, "turned off")
            let unhidden_systems = document.querySelectorAll("." + unhidden_class);

            for (let row = 0; row < unhidden_systems.length; row++) {
                unhidden_systems[row].classList.remove(unhidden_class);
                unhidden_systems[row].classList.add(hidden_class);
            }
            sessionStorage.setItem(checkbox_name, "false");
            break;
    }
}


var checkbox_systems = document.querySelector("#checkbox_systems");
checkbox_systems.addEventListener('change', function () {
    toggle_option(this.checked, "checkbox_systems", "hidden_system", "unhidden_system");
});

let checkbox_systems_state = sessionStorage.getItem("checkbox_systems");
if (checkbox_systems_state !== null) {
    let checked = checkbox_systems_state === "true";
    toggle_option(checked, "checkbox_systems", "hidden_system", "unhidden_system");
    checkbox_systems.checked = checked;
}


var checkbox_labels = document.querySelector("#checkbox_map_labels");
checkbox_labels.addEventListener('change', function () {
    toggle_option(this.checked, "checkbox_map_labels", "unhidden_label", "hidden_label");
});

let checkbox_label_state = sessionStorage.getItem("checkbox_map_labels");
if (checkbox_label_state !== null) {
    let checked = checkbox_label_state === "true";
    toggle_option(checked, "checkbox_map_labels", "unhidden_label", "hidden_label");
    checkbox_labels.checked = checked;
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