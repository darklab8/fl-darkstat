/**
 * Calculate count of visible elements in Table.
 * You should probably account elements with `tr[i].style.display = "none";` too may be here.
 * @param {HTMLTableElement} table
 */
function TableLen(table) {
    let count = 0;
    for (let i = 0, row; row = table.rows[i]; i++) {

        if (!row.classList.contains(HIDDEN_CLS)) {
            count = count + 1;
        }
    }
    // console.log("count=" + count)
    return count;
}

var HIDDEN_CLS = "hidden";

/**
 * hide, row or table or anything else
 * @param {string} id
 */
function Hide(id) {
    let element = document.getElementById(id);
    // console.log("Hide.id=" + id)
    if (!element.classList.contains(HIDDEN_CLS)) {
        element.classList.add(HIDDEN_CLS);
    }
}

/**
 * unhide, row or table or anything else
 * @param {string} id
 */
function Unhide(id) {
    let element = document.getElementById(id);
    // console.log("Unhide.id=" + id)
    if (element.classList.contains(HIDDEN_CLS)) {
        element.classList.remove(HIDDEN_CLS);
    }
}

/**
 * Function helping to persist selected ID
 * when user moves across different tabs
 */
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

/**
 * Implements functionality for filtering search bar
 * For table that has also filtering by selected ID tech compatibility, which is needed for Freelancer Discovery
 */
function FilteringFunction() {
    // Declare variables
    // console.log("triggered FilteringFunction")
    let input, filter, filter_infocard, table, tr, i, txtValue, txtValue_infocard;
    input = document.getElementById("filterinput");
    input_infocard = document.getElementById("filterinput_infocard");
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
        row = tr[0];
        cell = row.getElementsByClassName("tech_compat")[0];
        if (typeof (cell) != 'undefined') {
            cell.style.display = "none";
        }

    } else {
        row = tr[0];
        cell = row.getElementsByClassName("tech_compat")[0];
        if (typeof (cell) != 'undefined') {
            cell.style.display = "";
        }

    }

    // Loop through all table rows, and hide those who don't match the search query
    for (i = 1; i < tr.length; i++) {
        // row = document.getElementById("bottominfo_dsy_councilhf")
        row = tr[i];

        let txtValues = []
        let tds = row.getElementsByClassName("search-included")
        for (let elem of tds) {
            value = elem.textContent || elem.innerText;
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
        cell = row.getElementsByClassName("tech_compat")[0];
        if (typeof (cell) != 'undefined') {
            techcompats = JSON.parse(cell.attributes["techcompats"].textContent);
            compatibility = techcompats[tractor_id_selected] * 100;
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

    let infocards = document.getElementsByClassName("infocard")
    if (infocards.length > 0) {
        let infocard = infocards[0]
        highlight(input_infocard, infocard)
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

/**
 * Implements functionality for filtering search bar
 * @param {string} table_selector
 * @param {string} input_selector
 */
function FilteringForAnyTable(table_selector, input_selector) {
    // Declare variables
    // console.log("triggered FilteringFunction")
    let input, filter, table, tr, i, txtValue;
    input = document.getElementById(input_selector); // "filterinput"
    filter = input.value.toUpperCase();
    table = document.querySelector(table_selector); // "#table-top table"
    tr = table.getElementsByTagName("tr");


    // Loop through all table rows, and hide those who don't match the search query
    for (i = 1; i < tr.length; i++) {
        row = tr[i];
        txtValue = row.textContent || row.innerText;

        if (txtValue.toUpperCase().indexOf(filter) > -1) {
            tr[i].style.display = "";
            // console.log("row-i", i, "is made visible");
        } else {
            tr[i].style.display = "none";
            // console.log("row-i", i, "is made invisible");
        }
    }
}

/**
 * Highlights clicked table row
 * @param {HTMLTableRowElement} row
 */
function RowHighlighter(row) {
    let table = row.parentElement.parentElement;

    let selected_row_id = row.rowIndex;

    let rowsNotSelected = table.getElementsByTagName('tr');
    for (let row = 0; row < rowsNotSelected.length; row++) {
        rowsNotSelected[row].classList.remove('selected_row');
    }
    let rowSelected = table.getElementsByTagName('tr')[selected_row_id];
    rowSelected.classList.add("selected_row");
}
