/* anti-overlap code start */
function objectTerritorialConflictResolver(objects) {
    let iterationCount = 0;
    const MAX_ITERATIONS = 16;

    while (iterationCount < MAX_ITERATIONS) {
        const rects = objects.map(el => el.getBoundingClientRect());
        const margins = objects.map(el => parseFloat(el.style.marginTop) || 0);

        let totalDiff = 0;

        for (let i = 0; i < objects.length; i++) {
            // j = i+1 avoids duplicate pairs
            for (let j = i + 1; j < objects.length; j++) {
                if (!rectsOverlap(rects[i], rects[j])) continue;

                const a = rects[i], b = rects[j];

                if (a.top <= b.top) {
                    const diff = a.bottom - b.top;
                    if (diff > 0) {
                        margins[j] += diff;
                        totalDiff += diff;
                    }
                } else {
                    const diff = b.bottom - a.top;
                    if (diff > 0) {
                        margins[i] += diff;
                        totalDiff += diff;
                    }
                }
            }
        }

        // Batch write
        objects.forEach((el, i) => {
            el.style.marginTop = margins[i] + "px";
        });

        iterationCount++;
        if (totalDiff === 0) break;
    }

    console.log(`Labels settled after ${iterationCount} iterations`);
}

function rectsOverlap(a, b) {
    if (b.left > a.right || b.right < a.left) return false;
    if (b.top > a.bottom || b.bottom < a.top) return false;
    return true;
}

function InstallLabelOverlapper() {
    let labels = [...document.querySelectorAll("system-label")]
        .filter(el => getComputedStyle(el).display !== "none");

    objectTerritorialConflictResolver(labels);
}
function LabelsReset() {
    let labels = document.querySelectorAll("system-label")
    for (let row = 0; row < labels.length; row++) {
        labels[row].style.marginTop = "0px";
    }
}

/* anti-overlap code end */