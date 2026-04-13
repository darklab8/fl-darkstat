
/**
 * Clipboard functionality for waypoint copying
 * @param {number} x
 * @param {number} y
 * @param {number} z
 * @param {HTMLElement} element
 */
function copyWaypointToClipboard(x, y, z, element) {
    const waypointCmd = `/wp ${x} ${y} ${z}`;
    console.log('Attempting to copy:', waypointCmd);

    // Try modern clipboard API first, fallback if it fails
    if (navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(waypointCmd).then(() => {
            console.log('Copied via modern API:', waypointCmd);
        }).catch(err => {
            console.error('Modern API failed, trying fallback:', err);
            if (fallbackCopyToClipboard(waypointCmd)) {
                console.log('Copied via fallback method:', waypointCmd);
            } else {
                console.error('Both methods failed');
            }
        });
    } else if (fallbackCopyToClipboard(waypointCmd)) {
        console.log('Copied via fallback method (no modern API):', waypointCmd);
    } else {
        console.error('No clipboard method available');
    }

    // Visual feedback
    if (element) {
        element.style.opacity = '1';
        element.style.transform = 'scale(1.2)';
        setTimeout(() => {
            element.style.opacity = '0.7';
            element.style.transform = 'scale(1)';
        }, 200);
    }
}

/**
 * Fallback copy to clipboard method for browsers without clipboard API
 * @param {string} text
 */
function fallbackCopyToClipboard(text) {
    try {
        const textArea = document.createElement('textarea');
        textArea.value = text;
        textArea.style.position = 'fixed';
        textArea.style.left = '-999999px';
        textArea.style.top = '-999999px';
        document.body.appendChild(textArea);
        textArea.focus();
        textArea.select();
        const successful = document.execCommand('copy');
        document.body.removeChild(textArea);
        return successful;
    } catch (err) {
        console.error('Fallback copy failed:', err);
        return false;
    }
}
