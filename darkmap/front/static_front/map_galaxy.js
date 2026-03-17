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



function switch_systems(checked) {
    switch (checked) {
        case true:
            console.log(this.value, "turned on")

            let hidden_systems = document.querySelectorAll(".hidden_system");

            for (let row = 0; row < hidden_systems.length; row++) {
                hidden_systems[row].classList.remove('hidden_system');
                hidden_systems[row].classList.add("unhidden_system");
            }
            sessionStorage.setItem("checkbox_systems", "true");
            break;
        case false:
            console.log(this.value, "turned off")
            let unhidden_systems = document.querySelectorAll(".unhidden_system");

            for (let row = 0; row < unhidden_systems.length; row++) {
                unhidden_systems[row].classList.remove("unhidden_system");
                unhidden_systems[row].classList.add("hidden_system");
            }
            sessionStorage.setItem("checkbox_systems", "false");
            break;
    }
}

var checkbox_systems = document.querySelector("#checkbox_systems");
checkbox_systems.addEventListener('change', function () {
    switch_systems(this.checked);
});

let checkbox_systems_state = sessionStorage.getItem("checkbox_systems");
if (checkbox_systems_state !== null) {
    let checked = checkbox_systems_state === "true";
    switch_systems(checked);
    checkbox_systems.checked = checked;
}