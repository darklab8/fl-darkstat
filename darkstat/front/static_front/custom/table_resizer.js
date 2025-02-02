function makeTopBottomTablesResizable() { // eslint-disable-line no-unused-vars
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

        top_height_perc = 65
        if (element_top.hasOwnProperty("height_perc")) {  // eslint-disable-line no-prototype-builtins
            top_height_perc = element_top.height_perc
        }

        botttom_height_perc = 35
        if (element_bottom.hasOwnProperty("height_perc")) {  // eslint-disable-line no-prototype-builtins
            botttom_height_perc = element_bottom.height_perc
        }

        original_top_height = parseFloat(getComputedStyle(element_top, null).getPropertyValue('height').replace('px', ''));
        original_botttom_height = parseFloat(getComputedStyle(element_bottom, null).getPropertyValue('height').replace('px', ''));

        top_rect_top = element_top.getBoundingClientRect().top;
        bottom_rect_bottom = element_bottom.getBoundingClientRect().bottom;

        window.addEventListener('mousemove', resize1)
        window.addEventListener('mouseup', stopResize1)
    })

    function resize1(e) {
        var new_top_height = (e.pageY - top_rect_top) / original_top_height * top_height_perc
        var new_bottom_height = (bottom_rect_bottom - e.pageY) / original_botttom_height * botttom_height_perc

        element_top.style.height = "calc(" + new_top_height + "% - 7px)";
        element_bottom.style.height = "calc(" + new_bottom_height + "% - 7px)";

        element_top.height_perc = new_top_height
        element_bottom.height_perc = new_bottom_height
    }

    function stopResize1() {
        window.removeEventListener('mousemove', resize1)
    }
}

function makeLeftRightTablesResizable() { // eslint-disable-line no-unused-vars
    const currentResizer = document.querySelector('.resizer-left-right')
    const element_left = document.querySelector("#table-wrapper")
    const element_right = document.querySelector("#infocard_view")

    let left_width_perc = 0;
    let right_width_perc = 0;

    let original_left_width = 0;
    let original_right_width = 0;

    let left_rect_left = 0;
    let right_rect_right = 0;

    currentResizer.addEventListener('mousedown', function (e) {
        e.preventDefault()
        left_width_perc = element_left.style.width.replace('%', '');
        right_width_perc = element_right.style.width.replace('%', '');

        original_left_width = parseFloat(getComputedStyle(element_left, null).getPropertyValue('width').replace('px', ''));
        original_right_width = parseFloat(getComputedStyle(element_right, null).getPropertyValue('width').replace('px', ''));

        left_rect_left = element_left.getBoundingClientRect().left;
        right_rect_right = element_right.getBoundingClientRect().right;

        window.addEventListener('mousemove', resize)
        window.addEventListener('mouseup', stopResize)
    })

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