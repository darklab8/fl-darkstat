var map = document.querySelector('.panzoom');
var panzoom = Panzoom(map, {
    maxScale: 5,
    minScale: 1,
    startScale: 1,
    // panOnlyWhenZoomed: false,
    // canvas: true,
    // contain: "outside",
    handleStartEvent: function (event) {
        event.preventDefault()
    },
    noBind: false,
});
map.parentElement.addEventListener('wheel', panzoom.zoomWithWheel)
