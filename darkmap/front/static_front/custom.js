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
    noBind: true,
});
map.parentElement.addEventListener('wheel', panzoom.zoomWithWheel)

map.addEventListener('panzoomchange', function (event) {
    if (event.detail.scale == 1.0 && event.detail.x != 0 && event.detail.y != 0) {
        console.log("event", event)
        panzoom.reset({ scale: event.detail.scale })
    }
});
