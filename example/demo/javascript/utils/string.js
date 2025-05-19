/* @@docs String Utils
This section has unusual comment placement.
*/
function capitalize(str) {
    return str.charAt(0).toUpperCase() + str.slice(1);
}

// @@docs Padding
// String padding functions
function padLeft(str, char, length) {
    while (str.length < length) {
        str = char + str;
    }
    return str;
}