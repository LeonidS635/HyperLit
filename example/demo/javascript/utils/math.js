/* @@docs Vector Math
Operations _with_ 2D vectors.
This uses **object-oriented style**.
*/
class Vector {
    constructor(x, y) {
        this.x = x;
        this.y = y;
    }

    add(other) {
        return new Vector(this.x + other.x, this.y + other.y);
    }
}

/* @@docs Geometry
### Geometric calculations.

This section has the same indentation as the first one.
*/

function distance(a, b) {
    return Math.sqrt(Math.pow(a.x - b.x, 2) + Math.pow(a.y - b.y, 2));
}

// @@docs Constants
// Mathematical constants used throughout the app.
const PI = 3.141592;
const E = 2.718281;