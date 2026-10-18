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

    if (checkbox_name === "checkbox_systems") {
        LabelsReset();
        InstallLabelOverlapper();
    } else if (checkbox_name === "checkbox_map_labels") {
        LabelsReset();
        InstallLabelOverlapper();
    } else if (checkbox_name === "checkbox_wrecks_labels") {
        LabelsReset();
        InstallLabelOverlapper();
    } else if (checkbox_name === "checkbox_infocarded_labels") {
        LabelsReset();
        InstallLabelOverlapper();
    } else if (checkbox_name === "checkbox_objects") {
        LabelsReset();
        InstallLabelOverlapper();
    } else if (checkbox_name === "checkbox_obj_others") {
        LabelsReset();
        InstallLabelOverlapper();
    } else if (checkbox_name === "checkbox_obj_noname_labels") {
        LabelsReset();
        InstallLabelOverlapper();
    } else if (checkbox_name === "checkbox_pobs") {
        LabelsReset();
        InstallLabelOverlapper();
    }
}

const styleIdOn = 'tippy-toggle-on-style';
const styleIdOff = 'tippy-toggle-off-style';

function toggle_tippy_on(checkbox_name) {
    let style = document.getElementById(styleIdOn);
    if (!style) {
        style = document.createElement('style');
        style.id = styleIdOn;
        document.head.appendChild(style);
    }
    style.textContent = '.tippy-box { display: block !important; }';

    const existing2 = document.getElementById(styleIdOff);
    if (existing2) existing2.remove();

    sessionStorage.setItem(checkbox_name, "true");
}
function toggle_tippy_off(checkbox_name) {
    let style2 = document.getElementById(styleIdOff);
    if (!style2) {
        style2 = document.createElement('style');
        style2.id = styleIdOff;
        document.head.appendChild(style2);
    }
    style2.textContent = '.tippy-box.yes-height { display: block !important; }';

    const existing = document.getElementById(styleIdOn);
    if (existing) existing.remove();

    sessionStorage.setItem(checkbox_name, "false");
}

function toggle_tippy(checked, checkbox_name, hidden_class, unhidden_class) {
    switch (checked) {
        case true:
            console.log(this.value, "turned on")
            toggle_tippy_on(checkbox_name)
            break;
        case false:
            console.log(this.value, "turned off")
            toggle_tippy_off(checkbox_name)
            break;
    }
    toggle_option(checked, checkbox_name, hidden_class, unhidden_class);
}

/**
 * 
 * @param {string} button_id // id of a button without '#'
 * @param {string} default_state // default class state
 * @param {string} togglable_state // togglable class state
 */
function InstallButton(button_id, default_state, togglable_state, logic_func) {
    var checkbox_systems = document.querySelector("#" + button_id);
    console.log("installing button=", button_id);
    checkbox_systems.addEventListener('change', function () {
        logic_func(this.checked, button_id, default_state, togglable_state);
    });
    let checkbox_systems_state = sessionStorage.getItem(button_id);
    if (checkbox_systems_state !== null) {
        let checked = checkbox_systems_state === "true";
        if (checked === true) {
            logic_func(checked, button_id, default_state, togglable_state);
            checkbox_systems.checked = checked;
        }
    }
}


function toggle_grid_moving(checked, checkbox_name, hidden_class, unhidden_class) {
    switch (checked) {
        case true:
            console.log(this.value, "turned on pan")
            PanzoomToggleContain(true);
            sessionStorage.setItem(checkbox_name, "true");
            break;
        case false:
            console.log(this.value, "turned off pan")
            PanzoomToggleContain(false);
            sessionStorage.setItem(checkbox_name, "false");
            break;
    }
}

function InstallMenu() {
    InstallButton("checkbox_systems", "hidden_system", "unhidden_system", toggle_option);

    InstallButton("checkbox_map_labels", "unhidden_label", "hidden_label", toggle_option);

    InstallButton("checkbox_wrecks", "hidden_wreck", "unhidden_wreck", toggle_option);

    InstallButton("checkbox_wrecks_labels", "hidden_wreck_label", "unhidden_wreck_label", toggle_option);

    InstallButton("checkbox_objects", "hidden_obj", "unhidden_obj", toggle_option);

    InstallButton("checkbox_zones", "unhidden_zone", "hidden_zone", toggle_option);

    InstallButton("checkbox_all_zones", "hidden_all_zone", "unhidden_all_zone", toggle_option);

    InstallButton("checkbox_zone_cylinders", "hiddenCylinderZone", "unhiddenCylinderZone", toggle_option);

    InstallButton("checkbox_coords", "hidden_coords", "unhidden_coords", toggle_tippy);
    toggle_tippy_off("checkbox_coords")

    InstallButton("checkbox_pobs", "unhidden_pob", "hidden_pob", toggle_option);

    InstallButton("checkbox_obj_others", "hidden_obj_other_label", "unhidden_obj_other_label", toggle_option);

    InstallButton("checkbox_obj_noname_labels", "hidden_obj_noname_label", "unhidden_obj_noname_label", toggle_option);

    InstallButton("checkbox_infocarded_labels", "hidden_infocard_label", "unhidden_infocard_label", toggle_option);

    InstallButton("checkbox_move_out_of_grid", "none", "none", toggle_grid_moving);
}
