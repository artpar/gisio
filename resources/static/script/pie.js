/**
 * Created by parth on 2/12/2016.
 */

function appendPieChartByMap(data, container, invert) {
    var arrayData;
    if (!(data[0] instanceof Array )) {
        arrayData = mapTo2dArray(data);
    } else {
        arrayData = data;
    }
    if (invert) {
        console.log("invert data - ", arrayData);
        for (var i = 0; i < arrayData.length; i++) {
            var temp = arrayData[i][0];
            arrayData[i][0] = arrayData[i][1];
            arrayData[i][1] = temp;
        }
    }
    console.log("pie chart data ", JSON.stringify(data), JSON.stringify(arrayData));
    nv.addGraph(function () {
        var chart = nv.models.pieChart()
            .x(function (d) {
                console.log("for zero", d);
                if (d[0] == "0") {
                    return "zero"
                }
                if (d[0] == "1") {
                    return "one"
                }
                return d[0]
            })
            .y(function (d) {
                return d[1]
            })
            .color(d3.scale.category20())
            .showLabels(true);

        container.datum(arrayData)
            .transition().duration(1200)
            .call(chart);

        return chart;
    }, function () {
        console.log("pie chart complete")
    });
}
