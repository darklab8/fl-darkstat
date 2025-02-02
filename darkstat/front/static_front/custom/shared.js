/**
 * Implements functionality for filtering search bar
 * @param {HTMLElement} item
 * @param {string} excepted_filter
 */
function IsHavingLocksFromOtherFilters(item, excepted_filter) { // eslint-disable-line no-unused-vars
    const filterings = ["darkstat_filtering1", "darkstat_filtering2"];

    for (let i = 0; i < filterings.length; i++) {
        if (filterings[i] in item && excepted_filter !== filterings[i]) {
            return true
        }
    }

    return false
}