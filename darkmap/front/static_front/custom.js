
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