
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
    InstallLabelOverlapper();
}

function InstallMenu() {
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
}

var zoomInTreshold = 1.25;

function InstallPanzoom() {
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
        noBind: false,
    });
    map.parentElement.addEventListener('wheel', panzoom.zoomWithWheel);

    document.body.classList.add("zoomedOut");

    map.addEventListener('panzoomchange', function (event) {

        console.log("event.detail.scale=", event.detail.scale);
        if (event.detail.scale > zoomInTreshold) {
            document.body.classList.add("zoomedIn");
            document.body.classList.remove("zoomedOut");
        } else {
            document.body.classList.remove("zoomedIn");
            document.body.classList.add("zoomedOut");
        }
    });
}

/* anti-overlap code start */
function objectTerritorialConflictResolver(objects) {
    var currentDiffSum = "nope";
    var prevDiffSum = "nope";
    var prevPrevDiffSum;
    var iterationCount = 1;
    while (prevPrevDiffSum != 0 && iterationCount < 8) {
        prevPrevDiffSum = prevDiffSum;
        prevDiffSum = currentDiffSum;
        currentDiffSum = 0;
        for (i = 0; i < objects.length; i++) {
            var objectArray = objects;
            var currentObject = objectArray[i];
            currentDiffSum += moveIfOverlapsAndReturnDiff(currentObject, objectArray);
        }
        iterationCount++;
    }
    console.log("Labels settled after " + iterationCount + " iterations");
}

function moveIfOverlapsAndReturnDiff(currentObject, objectArray) {
    var diffSum = 0;
    var reducedObjectArray = objectArray; //objectArray.splice(i, 1);
    for (o = 0; o < reducedObjectArray.length; o++) {
        if (overlaps(currentObject, reducedObjectArray[o]) && currentObject != reducedObjectArray[o]) {
            if ((currentObject.getBoundingClientRect().top) <= (reducedObjectArray[o].getBoundingClientRect().top)) {
                var currentTransform;
                if (reducedObjectArray[o].style.marginTop.match(/([\d\.]+)/g) && reducedObjectArray[o].style.marginTop.match(/([\d\.]+)/g) != null) {
                    currentTransform = parseFloat(reducedObjectArray[o].style.marginTop.match(/([\d\.]+)/g));
                } else {
                    currentTransform = 0;
                }
                reducedObjectArray[o].style.marginTop = Math.abs(currentTransform + currentObject.getBoundingClientRect().bottom - reducedObjectArray[o].getBoundingClientRect().top) + "px";
                diffSum += Math.abs(currentObject.getBoundingClientRect().bottom - reducedObjectArray[o].getBoundingClientRect().top);
                /*moveIfOverlaps(reducedObjectArray[o], reducedObjectArray);*/
            } else {
                var currentTransform;
                if (currentObject.style.marginTop.match(/([\d\.]+)/g) && currentObject.style.marginTop.match(/([\d\.]+)/g) != null) {
                    currentTransform = parseFloat(currentObject.style.marginTop.match(/([\d\.]+)/g));
                } else {
                    currentTransform = 0;
                }
                currentObject.style.marginTop = Math.abs(currentTransform + reducedObjectArray[o].getBoundingClientRect().bottom - currentObject.getBoundingClientRect().top) + "px";
                diffSum += Math.abs(reducedObjectArray[o].getBoundingClientRect().bottom - currentObject.getBoundingClientRect().top);
                /*moveIfOverlaps(currentObject, reducedObjectArray);*/
            }
        }
    }
    return diffSum;
}

function overlaps(objectA, objectB) {
    var a = objectA.getBoundingClientRect();
    var b = objectB.getBoundingClientRect();

    var al = a.left;
    var ar = a.left + a.width;
    var bl = b.left;
    var br = b.left + b.width;

    var at = a.top;
    var ab = a.top + a.height;
    var bt = b.top;
    var bb = b.top + b.height;

    if (bl > ar || br < al) { return false; } /*overlap not possible*/
    if (bt > ab || bb < at) { return false; } /*overlap not possible*/

    if (bl > al && bl < ar) { return true; }
    if (br > al && br < ar) { return true; }

    if (bt > at && bt < ab) { return true; }
    if (bb > at && bb < ab) { return true; }

    return false;
}

function InstallLabelOverlapper() {
    let labels = document.querySelectorAll("system-label")
    objectTerritorialConflictResolver(labels);
}
/* anti-overlap code end */
