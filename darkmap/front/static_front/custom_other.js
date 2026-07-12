function highlightFromQuery() {
    var params = new URLSearchParams(window.location.search);
    var q = params.get('q');
    if (!q) return;

    var targetEl = document.getElementById(q);
    if (!targetEl) return;

    targetEl.classList.add('is-target');

    targetEl.addEventListener('animationend', function onEntranceDone() {
        targetEl.classList.add('pulse-normal');
        targetEl.removeEventListener('animationend', onEntranceDone);
    });
}


function CloseInfocard() {
    document.querySelector("#remodal-bg").style.display = "none";
}

function updateQuery(el) {
    const url = new URL(window.location.href);
    url.searchParams.set('q', el.dataset.nickname);
    window.history.pushState({}, '', url);

    document.querySelectorAll('.is-target').forEach(obj => obj.classList.remove('is-target'));
}

function SetCursorPosition() {
    const panzoomEl = document.querySelector('.panzoom');
    const outputEl = document.querySelector('cursor-pos-view-');

    if (!panzoomEl || !outputEl) {
        console.warn('cursor-pos-view: required elements not found');
        return;
    }

    function getFlLength() {
        const val = parseFloat(outputEl.getAttribute('data-fl-length'));
        return isNaN(val) ? 0 : val;
    }

    function formatWithUnderscores(num) {
        const rounded = Math.round(num);
        const isNegative = rounded < 0;
        const digits = Math.abs(rounded).toString();

        let result = '';
        for (let i = 0; i < digits.length; i++) {
            const posFromEnd = digits.length - i;
            if (i > 0 && posFromEnd % 3 === 0) {
                result += '_';
            }
            result += digits[i];
        }

        return (isNegative ? '-' : '') + result;
    }

    function handleMouseMove(e) {
        const rect = panzoomEl.getBoundingClientRect();
        const flLength = getFlLength();

        if (rect.width === 0 || rect.height === 0 || flLength === 0) {
            return;
        }

        // Cursor position as a fraction (0 -> left/top edge, 1 -> right/bottom edge)
        const fracX = (e.clientX - rect.left) / rect.width;
        const fracY = (e.clientY - rect.top) / rect.height;

        // Center of .panzoom = (0,0) in game coords, full width/height = flLength
        const x = (fracX - 0.5) * flLength;
        const z = (fracY - 0.5) * flLength;

        outputEl.textContent = `X: ${formatWithUnderscores(x)}, Z: ${formatWithUnderscores(z)}`;
    }
    document.addEventListener('mousemove', handleMouseMove);
}
