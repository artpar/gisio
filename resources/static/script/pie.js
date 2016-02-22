/**
 * Created by parth on 2/12/2016.
 */

function appendPieChart(data, container, invert) {
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
        chart.pie.dispatch.on("elementClick", function (d) {
            console.log("you clicked ", d)
        });
        return chart;
    });
}
