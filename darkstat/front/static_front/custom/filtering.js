/**
 * Implements functionality for filtering search bar
 * For table that has also filtering by selected ID tech compatibility, which is needed for Freelancer Discovery
 */
function FilteringFunction() { // eslint-disable-line no-unused-vars
    // Declare variables
    // console.log("triggered FilteringFunction")
    let input, filter, filter_infocard, table, tr, txtValue, txtValue_infocard;
    input = document.getElementById("filterinput");
    if (typeof (input) === 'undefined' || input === null) {
        return;
    }
    let input_infocard = document.getElementById("filterinput_infocard");
    filter_infocard = input_infocard.value.toUpperCase();
    filter = input.value.toUpperCase();
    table = document.querySelector("#table-top table");
    tr = table.getElementsByTagName("tr");

    // Select current ID tractor
    let tractor_id_elem, tractor_id_selected;
    tractor_id_selected = "";
    tractor_id_elem = document.getElementById("tractor_id_selector");
    if (typeof (tractor_id_elem) != 'undefined' && tractor_id_elem != null) {
        tractor_id_selected = tractor_id_elem.value;

        sessionStorage.setItem("tractor_id_selected_index", tractor_id_elem.selectedIndex);

    }

    // making invisible info about ID Compatibility if no ID is selected
    if (tractor_id_selected === "") {
        let row = tr[0];
        let cell = row.getElementsByClassName("tech_compat")[0];
        if (typeof (cell) != 'undefined') {
            cell.style.display = "none";
        }

    } else {
        let row = tr[0];
        let cell = row.getElementsByClassName("tech_compat")[0];
        if (typeof (cell) != 'undefined') {
            cell.style.display = "";
        }

    }

    // Loop through all table rows, and hide those who don't match the search query
    for (let i = 1; i < tr.length; i++) {
        // row = document.getElementById("bottominfo_dsy_councilhf")
        let row = tr[i];

        let txtValues = []
        let tds = row.getElementsByClassName("seo")
        for (let elem of tds) {
            let value = elem.textContent || elem.innerText;
            txtValues.push(value)
        }
        txtValue = txtValues.join('');

        let infocards = row.getElementsByClassName("search-infocard");
        txtValue_infocard = '';
        if (infocards.length > 0) {
            txtValue_infocard = infocards[0].textContent || infocards[0].innerText
        }

        // Refresh tech compat value
        let techcompat_visible = true;
        let compatibility;
        let cell = row.getElementsByClassName("tech_compat")[0];
        if (typeof (cell) != 'undefined') {
            let techcompats = JSON.parse(cell.attributes["techcompats"].textContent.replaceAll("'", '"'));

            if (tractor_id_selected in techcompats) {
                compatibility = techcompats[tractor_id_selected] * 100;
            } else {
                compatibility = 0;
            }
            cell.innerHTML = compatibility + "%";


            techcompat_visible = compatibility > 10 || tractor_id_selected === ""

            // making invisible info about ID Compatibility if no ID is selected
            if (tractor_id_selected === "") {
                cell.style.display = "none";
            } else {
                cell.style.display = "";
            }

            // console.log("compatibility=", compatibility, "tractor_id_selected=", tractor_id_selected, "techcompat_visible=", techcompat_visible)
        }

        if ((txtValue.toUpperCase().indexOf(filter) > -1 && txtValue_infocard.toUpperCase().indexOf(filter_infocard) > -1) && techcompat_visible === true) {
            tr[i].style.display = "";
            // console.log("row-i", i, "is made visible");
        } else {
            tr[i].style.display = "none";
            // console.log("row-i", i, "is made invisible");
        }
    }
}

/**
 * Implements functionality for filtering search bar
 * @param {string} table_selector
 * @param {string} input_selector
 */
function FilteringForAnyTable(table_selector, input_selector) { // eslint-disable-line no-unused-vars
    // Declare variables
    // console.log("triggered FilteringFunction")
    let input, filter, table, tr, txtValue;
    input = document.getElementById(input_selector); // "filterinput"
    filter = input.value.toUpperCase();
    table = document.querySelector(table_selector); // "#table-top table"
    tr = table.getElementsByTagName("tr");


    // Loop through all table rows, and hide those who don't match the search query
    for (let i = 1; i < tr.length; i++) {
        let row = tr[i];
        txtValue = row.textContent || row.innerText;

        if (IsHavingLocksFromOtherFilters(row, 'darkstat_filtering1')) { // eslint-disable-line no-undef
            continue
        }

        if (txtValue.toUpperCase().indexOf(filter) > -1) {
            tr[i].style.display = "";
            if ('darkstat_filtering1' in row) {
                delete row.darkstat_filtering1
            }
            // console.log("row-i", i, "is made visible");
        } else {
            tr[i].style.display = "none";
            row.darkstat_filtering1 = true
            // console.log("row-i", i, "is made invisible");
        }
    }
}

/**
 * Useful to highlight searched text in an infocard
 * @param {HTMLElement} inputText
 * @param {HTMLElement} text
 */
function highlight(input_infocard, infocard) {
    let innerHTML = infocard.innerHTML;

    innerHTML = innerHTML.replaceAll("<highlight>", "");
    innerHTML = innerHTML.replaceAll("</highlight>", "");
    document.getElementsByClassName("infocard")[0].innerHTML = innerHTML

    // let index = infocard.innerHTML.toUpperCase().indexOf(input_infocard.value.toUpperCase());

    if (input_infocard.value.length < 2) {
        return
    }

    var searchMask = input_infocard.value;
    var regEx = new RegExp(searchMask, "ig");
    var replaceMask = "<highlight>" + input_infocard.value + "</highlight>";
    innerHTML = innerHTML.replace(regEx, replaceMask);
    document.getElementsByClassName("infocard")[0].innerHTML = innerHTML
}

function highlightInfocardHook() {  // eslint-disable-line no-unused-vars
    let input_infocard = document.getElementById("filterinput_infocard");
    let infocards = document.getElementsByClassName("infocard")
    if (infocards.length > 0) {
        let infocard = infocards[0]
        highlight(input_infocard, infocard)
    }
}