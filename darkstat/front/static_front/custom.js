function TableLen(table) {
    var count = 0;
    for (var i = 0, row; row = table.rows[i]; i++) {

        if (!row.classList.contains("hidden")) {
            count = count + 1;
        }
    }
    // console.log("count=" + count)
    return count;
}


function Hide(id) {
    var element = document.getElementById(id);
    // console.log("Hide.id=" + id)
    element.classList.add("hidden");
}

function Unhide(id) {
    var element = document.getElementById(id);
    // console.log("Unhide.id=" + id)
    element.classList.remove("hidden");
}

function LoadSelectedTractorID() {
    // console.log("triggered LoadSelectedTractorID")
    let selected_index = sessionStorage.getItem("tractor_id_selected_index");
    if (typeof (selected_index) != 'undefined' && selected_index != null) {
        tractor_id_elem = document.getElementById("tractor_id_selector");
        if (typeof (tractor_id_elem) != 'undefined' && tractor_id_elem != null) {
            tractor_id_elem.selectedIndex = selected_index;
        }
    }
}

function FilteringFunction() {
    // Declare variables
    // console.log("triggered FilteringFunction")
    var input, filter, table, tr, td, i, txtValue;
    input = document.getElementById("filterinput");
    filter = input.value.toUpperCase();
    table = document.querySelector("#table-top table");
    tr = table.getElementsByTagName("tr");

    // Select current ID tractor
    var tractor_id_elem, tractor_id_selected;
    tractor_id_selected = "";
    tractor_id_elem = document.getElementById("tractor_id_selector");
    if (typeof (tractor_id_elem) != 'undefined' && tractor_id_elem != null) {
        tractor_id_selected = tractor_id_elem.value;

        sessionStorage.setItem("tractor_id_selected_index", tractor_id_elem.selectedIndex);

    }

    // Loop through all table rows, and hide those who don't match the search query
    for (i = 1; i < tr.length; i++) {
        // row = document.getElementById("bottominfo_dsy_councilhf")
        row = tr[i];
        txtValue = row.textContent || row.innerText;

        // Refresh tech compat value
        var techcompat_visible = true;
        cell = row.getElementsByClassName("tech_compat")[0];
        if (typeof (cell) != 'undefined') {
            techcompats = JSON.parse(cell.attributes["techcompats"].textContent);
            compatibility = techcompats[tractor_id_selected] * 100;
            cell.innerHTML = compatibility + "%";


            techcompat_visible = compatibility > 10 || tractor_id_selected === ""
            // console.log("compatibility=", compatibility, "tractor_id_selected=", tractor_id_selected, "techcompat_visible=", techcompat_visible)
        }

        if (txtValue.toUpperCase().indexOf(filter) > -1 && techcompat_visible === true) {
            tr[i].style.display = "";
            // console.log("row-i", i, "is made visible");
        } else {
            tr[i].style.display = "none";
            // console.log("row-i", i, "is made invisible");
        }
    }
}