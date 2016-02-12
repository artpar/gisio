/**
 * Created by parth on 2/12/2016.
 */

function appendPieChartByMap(data, container) {

    var keys = Object.keys(data);
    var arrayData = [];
    for (var i = 0; i < keys.length; i++) {
        arrayData.push([keys[i], data[keys[i]]]);
    }
    return appendPieChartByArray(arrayData, container);
}
function appendPieChartByArray(data, container) {
    console.log("plot data ", data);
    var width = container.attr('width'),
        height = container.attr('height'),
        radius = Math.min(width, height) / 2;

    container = container.append("g")
        .attr("transform", "translate(" + width / 2 + "," + height / 2 + ")");


    var color = d3.scale.category10();
    var arc = d3.svg.arc()
        .outerRadius(radius - 10)
        .innerRadius(0);

    var labelArc = d3.svg.arc()
        .outerRadius(radius - 40)
        .innerRadius(radius - 40);

    var pie = d3.layout.pie()
        .sort(null)
        .value(function (d) {
            return d[1];
        });

    var g = container.selectAll(".arc")
        .data(pie(data))
        .enter().append("g")
        .attr("class", "arc");

    g.append("path")
        .attr("d", arc)
        .style("fill", function (d) {
            return color(d.data[0]);
        });

    g.append("text")
        .attr("transform", function (d) {
            return "translate(" + labelArc.centroid(d) + ")";
        })
        .attr("dy", ".35em")
        .text(function (d) {
            return d.data[0];
        });
}
