/**
 * Created by parth on 2/12/2016.
 */

function appendBarChart(data, container) {
    var arrayData = mapTo2dArray(data);
    //console.log("bar chart", arrayData);
    nv.addGraph(function () {
        var chart = nv.models.discreteBarChart()
            .x(function (d) {
                return d[0]
            })
            .y(function (d) {
                return d[1]
            })
            .staggerLabels(true)
            //.staggerLabels(historicalBarChart[0].values.length > 8)
            .showValues(true)
            .duration(250);

        container.datum([{key: "", values: arrayData}])
            .call(chart);

        nv.utils.windowResize(chart.update);

        return chart;
    });
}

function appendAreaChart(data, container) {
    var arrayData = mapTo2dArray(data);
    var total = 0;
    for (var i = 0; i < arrayData.length; i++) {
        var newVal = parseInt(arrayData[i][0]);
        if (isNaN(newVal)) {
            newVal = total / (i + 1);
        }
        total = total + newVal;
        arrayData[i][0] = newVal
    }
    arrayData.sort(function (a, b) {
        return a[0] - b[0];
    });
    console.log("chart", arrayData);

    var colors = d3.scale.category10();
    var chart;
    nv.addGraph(function () {
        chart = nv.models.lineChart()
            .useInteractiveGuideline(true)
            .color(colors)
            .x(function (d) {
                return d[0]
            })
            .y(function (d) {
                return d[1]
            })
            .controlLabels({stacked: "Stacked"})
            .duration(300);
        //chart.xAxis.tickFormat(d3.format('%d'));
        //chart.yAxis.tickFormat(d3.format(',.4f'));
        chart.legend.vers('furious');
        container.datum([{
                "key": "Val",
                "values": arrayData
            }])
            .transition()
            .duration(1000)
            .call(chart)
            .each('start', function () {
                setTimeout(function () {
                    container.selectAll("*").each(function () {
                        if (this.__transition__)
                            this.__transition__.duration = 1;
                    })
                }, 0)
            });
        nv.utils.windowResize(chart.update);
        return chart;
    });

}


function appendLineChart(data, container) {
    var arrayData = mapTo2dArray(data);
    var total = 0;
    for (var i = 0; i < arrayData.length; i++) {
        var newVal = parseInt(arrayData[i][0]);
        if (isNaN(newVal)) {
            newVal = total / (i + 1);
        }
        total = total + newVal;
        arrayData[i][0] = newVal
    }
    arrayData.sort(function (a, b) {
        return a[0] - b[0];
    });
    console.log("line chart", arrayData);

    var colors = d3.scale.category10();
    var chart;
    nv.addGraph(function () {
        chart = nv.models.lineChart()
            //.useInteractiveGuideline(true)
            //.color(colors)
            .x(function (d) {
                return d[0]
            })
            .y(function (d) {
                return d[1]
            })
            //.controlLabels({stacked: "Stacked"})
            .duration(300);
        //chart.xAxis.tickFormat(d3.format('%d'));
        //chart.yAxis.tickFormat(d3.format(',.4f'));
        chart.legend.vers('furious');
        container.datum([{
                "key": "Val",
                "values": arrayData
            }])
            .transition()
            .duration(1000)
            .call(chart)
            .each('start', function () {
                setTimeout(function () {
                    container.selectAll("*").each(function () {
                        if (this.__transition__)
                            this.__transition__.duration = 1;
                    })
                }, 0)
            });
        nv.utils.windowResize(chart.update);
        return chart;
    });

}


