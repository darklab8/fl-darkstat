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