/**
 * Created by parth on 2/18/2016.
 */

function arrayToObjects(arrayData) {
    var keys = argumentSplice(arguments, 1);
    var newData = [];
    for (var i = 0; i < arrayData.length; i++) {
        var row = {};
        for (var j = 0; j < keys.length; j++) {
            row[keys[j]] = +arrayData[i][j] || 3
        }
        newData.push(row)
    }
    return newData;
}

function argumentSplice(args, start) {
    var data = [];
    for (var i = start; i < args.length; i++) {
        data.push(args[i])
    }
    return data
}


function appendScatterChart(data, container, invert) {

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

    arrayData = arrayToObjects(arrayData, "x", "y", "size");
    console.log("scatter data ", arrayData);
    nv.addGraph(function () {
        var chart = nv.models.scatterChart()
            .showDistX(true)    //showDist, when true, will display those little distribution lines on the axis.
            .showDistY(true)
            .color(d3.scale.category10().range());

        //Configure how the tooltip looks.
        //chart.tooltip().contentGenerator(function (key) {
        //    return '<h3>' + key + '</h3>';
        //});

        //Axis settings
        chart.xAxis.tickFormat(d3.format('.02f'));
        chart.yAxis.tickFormat(d3.format('.02f'));

        //We want to show shapes other than circles.
        //chart.scatter.onlyCircles(false);

        //var myData = randomData(4, 40);
        container.datum([{
                "key": "Value",
                "values": arrayData
            }])
            .transition().duration(350)
            .call(chart);

        nv.utils.windowResize(chart.update);

        return chart;
    });


}
