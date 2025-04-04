/*
How to implement?
Insert json with data of routes into every row. Right into cell needing recalculation ;)
[{"time": smth, "profit": "smth"}], format like this

In Input Value Change
    Grab input value.
    Per each ship category:
        Find max proffit per distance
        Which have required minimum distance.
        Insert updated to the Cell

On Render:
    Grab Input Value
    Hide Rows Which have distances in all Three cells less than minimum.

P.S. how to make it playing nice with other filters? mm... check some flag if it was already filtered by smth else.
U can https://stackoverflow.com/questions/4258466/can-i-add-arbitrary-properties-to-dom-objects
*/


function FilteringForDistances() { // eslint-disable-line no-unused-vars
    // Declare variables
    let input, table, tr, max_profit;

    input = document.getElementById("input_route_min_dist");
    let min_distance_threshold = input.value;
    if (min_distance_threshold === '') {
        min_distance_threshold = 0
    }

    table = document.querySelector("#table-top table");
    tr = table.getElementsByTagName("tr");

    // Loop through all table rows, and hide those who don't match the search query
    for (let i = 1; i < tr.length; i++) {
        let row = tr[i];

        for (let r = 0; r < route_types.length; r++) { // eslint-disable-line no-undef
            let cell = row.getElementsByClassName(route_types[r])[0]; // eslint-disable-line no-undef

            let routesinfo = JSON.parse(cell.attributes["routesinfo"].textContent);

            if (routesinfo === null) {
                continue
            }
            // list of { ProffitPetTime TotalSeconds } number values
            // renamed to { p s } for client side not overloading reasons. otherwise html was taking 155mb
            max_profit = 0
            for (let j = 0; j < routesinfo.length; j++) {
                if (routesinfo[j].S > min_distance_threshold) {
                    if (routesinfo[j].P > max_profit) {
                        max_profit = routesinfo[j].P
                    }
                }
            }

            cell.innerHTML = (100 * max_profit).toFixed(2);
        }
    }
}
