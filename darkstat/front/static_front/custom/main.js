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

function cloneAttributes(target, source) {
    [...source.attributes].forEach(attr => { target.setAttribute(attr.nodeName, attr.nodeValue) })
}

function LoadTechCompat() {
    console.log("loading tech compat attempt")
    let targets1 = document.querySelectorAll("[data-target-1]")
    for (let i = 0; i < targets1.length; i++) {
        let target = targets1[i]
        let cache_key = target.attributes["data-target-1"].value
        let source = document.querySelector("[data-source-1='" + cache_key + "']")
        target.innerHTML = source.innerHTML
        cloneAttributes(target, source)
        target.removeAttribute("data-target-1")
    }

    let targets2 = document.querySelectorAll("[data-target-2]")
    for (let i = 0; i < targets2.length; i++) {
        let target = targets2[i]
        let cache_key = target.attributes["data-target-2"].value
        let source = document.querySelector("[data-source-2='" + cache_key + "']")
        target.innerHTML = source.innerHTML
        cloneAttributes(target, source)
        target.removeAttribute("data-target-2")
    }
}
