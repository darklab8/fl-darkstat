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
