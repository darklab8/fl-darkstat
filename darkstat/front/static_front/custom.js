function TableLen(table) {
    var count = 0;
    for (var i = 0, row; row = table.rows[i]; i++) {

        if (!row.classList.contains("hidden")) {
            count = count + 1;
        }
    }
    console.log("count=" + count)
    return count;
}


function Hide(id) {
    var element = document.getElementById(id);
    console.log("Hide.id=" + id)
    element.classList.add("hidden");
}

function Unhide(id) {
    var element = document.getElementById(id);
    console.log("Unhide.id=" + id)
    element.classList.remove("hidden");
}

function FilteringFunction() {
    // Declare variables
    var input, filter, table, tr, td, i, txtValue;
    input = document.getElementById("filterinput");
    filter = input.value.toUpperCase();
    table = document.querySelector("#table-top table");
    tr = table.getElementsByTagName("tr");

    // Loop through all table rows, and hide those who don't match the search query
    for (i = 1; i < tr.length; i++) {
        row = tr[i];
        txtValue = row.textContent || row.innerText;
        if (txtValue.toUpperCase().indexOf(filter) > -1) {
            tr[i].style.display = "";
        } else {
            tr[i].style.display = "none";
        }
    }
}