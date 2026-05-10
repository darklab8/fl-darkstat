
var zoomInTreshold = 1.25;



var DRAG_THRESHOLD = 5;
var _pointerDownX = 0;
var _pointerDownY = 0;
var _didDrag = false;

/**
 * @param {boolean} allowed_moved 
 */
function InstallPanzoom(allowed_moved) {
    var map = document.querySelector('.panzoom');

    let contain = null
    if (!allowed_moved) {
        contain = "outside"
    }

    var options = {
        maxScale: 7.5,
        minScale: 1,
        handleStartEvent: function (event) {
            event.preventDefault();
        },
        noBind: false,
        contain: contain,
    };

    var panzoom = Panzoom(map, options);

    // overrides to make possible moving on mouse click even when no longer selecting original "image"
    var parent = map.parentElement;
    parent.addEventListener('wheel', panzoom.zoomWithWheel);

    // ANTI INFOCARD CLICKER if PANNING/MOVING START
    parent.addEventListener('pointerdown', function (e) {
        _pointerDownX = e.clientX;
        _pointerDownY = e.clientY;
        _didDrag = false;
        panzoom.handleDown(e);
    });
    parent.addEventListener('pointermove', function (e) {
        var dx = e.clientX - _pointerDownX;
        var dy = e.clientY - _pointerDownY;
        if (Math.sqrt(dx * dx + dy * dy) > DRAG_THRESHOLD) {
            _didDrag = true;
        }
        panzoom.handleMove(e);
    });
    parent.addEventListener('pointerup', panzoom.handleUp);
    // Block HTMX requests on obj- elements if the user was dragging
    document.body.addEventListener('htmx:confirm', function (e) {
        var target = e.detail.elt;
        if (target.tagName.toLowerCase() === 'obj-' && _didDrag) {
            e.preventDefault(); // cancels hx-get
        }
    });
    // ANTI INFOCARD CLICKER if PANNING/MOVING END

    document.body.classList.add("zoomedOut");

    let timeout
    let timeout_ms = 500;
    const fullLabelsReset = () => {
        LabelsReset();
        InstallLabelOverlapper();
    }
    map.addEventListener('panzoomchange', function (event) {
        console.log("event.detail.scale=", event.detail.scale);
        document.documentElement.style.setProperty('--panzoom-scale', `${event.detail.scale}`);

        clearTimeout(timeout);
        timeout = setTimeout(fullLabelsReset, timeout_ms);

        if (event.detail.scale > zoomInTreshold) {
            document.body.classList.add("zoomedIn");
            document.body.classList.remove("zoomedOut");
        } else {
            document.body.classList.remove("zoomedIn");
            document.body.classList.add("zoomedOut");
        }

        if (event.detail.scale > 1.2) {
            document.body.classList.add("zoomed08");
        } else {
            document.body.classList.remove("zoomed08");
        }
        if (event.detail.scale > 1.5) {
            document.body.classList.add("zoomed06");
        } else {
            document.body.classList.remove("zoomed06");
        }
        if (event.detail.scale > 2.0) {
            document.body.classList.add("zoomed04");
        } else {
            document.body.classList.remove("zoomed04");
        }
    });
}
