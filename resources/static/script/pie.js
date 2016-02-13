/**
 * Created by parth on 2/12/2016.
 */

function appendPieChartByMap(data, container) {
    var arrayData = mapTo2dArray(data);
    nv.addGraph(function () {
        var chart = nv.models.pieChart()
            .x(function (d) {
                return d[0]
            })
            .y(function (d) {
                return d[1]
            })
            .showLabels(true);

        container.datum(arrayData)
            .transition().duration(1200)
            .call(chart);

        return chart;
    });
}
