/**
 * Created by parth on 2/12/2016.
 */

function appendBarChart(data, container) {
    var keys = Object.keys(data);
    var arrayData = [];
    for (var i = 0; i < keys.length; i++) {
        arrayData.push([keys[i], data[keys[i]]]);
    }
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
