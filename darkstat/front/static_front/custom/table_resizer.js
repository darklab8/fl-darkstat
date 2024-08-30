function makeTopBottomTablesResizable() {
    const element_top = document.querySelector("#table-top")
    const element_bottom = document.querySelector("#table-bottom")
    const currentResizer = document.querySelector('.resizer-top-bottom')

    let top_height_perc = 0;
    let botttom_height_perc = 0;

    let original_top_height = 0;
    let original_botttom_height = 0;

    let top_rect_top = 0;
    let bottom_rect_bottom = 0;

    currentResizer.addEventListener('mousedown', function (e) {
        e.preventDefault()
        top_height_perc = element_top.style.height.replace('%', '');
        botttom_height_perc = element_bottom.style.height.replace('%', '');

        original_top_height = parseFloat(getComputedStyle(element_top, null).getPropertyValue('height').replace('px', ''));
        original_botttom_height = parseFloat(getComputedStyle(element_bottom, null).getPropertyValue('height').replace('px', ''));

        top_rect_top = element_top.getBoundingClientRect().top;
        bottom_rect_bottom = element_bottom.getBoundingClientRect().bottom;

        window.addEventListener('mousemove', resize)
        window.addEventListener('mouseup', stopResize)
    })

    function resize(e) {

        var new_top_height = (e.pageY - top_rect_top) / original_top_height * top_height_perc
        var new_bottom_height = (bottom_rect_bottom - e.pageY) / original_botttom_height * botttom_height_perc

        element_top.style.height = new_top_height + "%";
        element_bottom.style.height = new_bottom_height + "%";
    }

    function stopResize() {
        window.removeEventListener('mousemove', resize)
    }
}

function makeLeftRightTablesResizable() {
    const element_left = document.querySelector("#table-wrapper")
    const element_right = document.querySelector("#infocard_view")
    const currentResizer = document.querySelector('.resizer-left-right')

    let left_width_perc = 0;
    let right_width_perc = 0;

    let original_left_width = 0;
    let original_right_width = 0;

    let left_rect_left = 0;
    let right_rect_right = 0;

    function resize_start(e) {
        e.preventDefault()
        left_width_perc = element_left.style.width.replace('%', '');
        right_width_perc = element_right.style.width.replace('%', '');

        original_left_width = parseFloat(getComputedStyle(element_left, null).getPropertyValue('width').replace('px', ''));
        original_right_width = parseFloat(getComputedStyle(element_right, null).getPropertyValue('width').replace('px', ''));

        left_rect_left = element_left.getBoundingClientRect().left;
        right_rect_right = element_right.getBoundingClientRect().right;

        window.addEventListener('mousemove', resize)
        window.addEventListener('mouseup', stopResize)
        window.addEventListener('touchmove', resize)
        window.addEventListener('touchup', stopResize)
    }

    currentResizer.addEventListener('mousedown', resize_start)
    currentResizer.addEventListener('touchdown', resize_start)

    function resize(e) {

        var new_left_width = (e.pageX - left_rect_left) / original_left_width * left_width_perc
        var new_right_width = (right_rect_right - e.pageX) / original_right_width * right_width_perc

        element_left.style.width = new_left_width + "%";
        element_right.style.width = new_right_width + "%";
    }

    function stopResize() {
        window.removeEventListener('mousemove', resize)
    }
}