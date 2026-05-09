var zoomInTreshold = 1.25;

var DRAG_THRESHOLD = 5;
var _pointerDownX = 0;
var _pointerDownY = 0;
var _didDrag = false;

// --- PANZOOM TOGGLE STATE ---
var _panzoom_instance = null;
var _panzoom_parent = null;
var _panzoom_listeners = null;
var _panzoom_map = null;
var _panzoom_on_change = null;

function panzoomDestroy() {
    if (_panzoom_instance) {
        var scale = _panzoom_instance.getScale();
        var pan = _panzoom_instance.getPan();
        _panzoom_instance.destroy();
        _panzoom_instance = null;
        if (_panzoom_parent && _panzoom_listeners) {
            _panzoom_parent.removeEventListener('wheel', _panzoom_listeners.wheel);
            _panzoom_parent.removeEventListener('pointerdown', _panzoom_listeners.pointerdown);
            _panzoom_parent.removeEventListener('pointermove', _panzoom_listeners.pointermove);
            _panzoom_parent.removeEventListener('pointerup', _panzoom_listeners.pointerup);
            document.body.removeEventListener('htmx:confirm', _panzoom_listeners.htmxConfirmer);
        }
        if (_panzoom_map && _panzoom_on_change) {
            _panzoom_map.removeEventListener('panzoomchange', _panzoom_on_change);
        }
        return { scale, pan };
    }
    return null;
}

/**
 * @param {boolean} allowed_moved 
 */
function PanzoomToggleContain(allowed_moved) {
    var snapshot = panzoomDestroy();
    InstallPanzoom(allowed_moved);
    if (snapshot) {
        _panzoom_instance.zoom(snapshot.scale, { animate: false, silent: true });
        _panzoom_instance.pan(snapshot.pan.x, snapshot.pan.y, { animate: false, silent: true });
    }
}
// --- PANZOOM TOGGLE STATE END ---

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
    _panzoom_instance = panzoom;

    // overrides to make possible moving on mouse click even when no longer selecting original "image"
    var parent = map.parentElement;
    _panzoom_parent = parent;

    var onWheel = function (e) { panzoom.zoomWithWheel(e); };
    var onPointerDown = function (e) {
        _pointerDownX = e.clientX;
        _pointerDownY = e.clientY;
        _didDrag = false;
        panzoom.handleDown(e);
    };
    var onPointerMove = function (e) {
        var dx = e.clientX - _pointerDownX;
        var dy = e.clientY - _pointerDownY;
        if (Math.sqrt(dx * dx + dy * dy) > DRAG_THRESHOLD) {
            _didDrag = true;
        }
        panzoom.handleMove(e);
    };
    var onPointerUp = function (e) { panzoom.handleUp(e); };

    // Block HTMX requests on obj- elements if the user was dragging
    const htmxConfirmer = function (e) {
        console.log("htmx confirmation running")
        var target = e.detail.elt;
        if (target.tagName.toLowerCase() === 'obj-' && _didDrag) {
            e.preventDefault(); // cancels hx-get
        }
    };

    _panzoom_listeners = {
        wheel: onWheel,
        pointerdown: onPointerDown,
        pointermove: onPointerMove,
        pointerup: onPointerUp,
        htmxConfirmer: htmxConfirmer,
    };

    parent.addEventListener('wheel', onWheel);

    parent.addEventListener('pointerdown', onPointerDown);
    parent.addEventListener('pointermove', onPointerMove);
    parent.addEventListener('pointerup', onPointerUp);

    document.body.addEventListener('htmx:confirm', htmxConfirmer);
    document.body.classList.add("zoomedOut");

    let timeout
    let timeout_ms = 500;
    const fullLabelsReset = () => {
        LabelsReset();
        InstallLabelOverlapper();
    }
    const onPanzoomChange = function (event) {
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
    };
    _panzoom_map = map;
    _panzoom_on_change = onPanzoomChange;
    map.addEventListener('panzoomchange', onPanzoomChange);
}
