/**
 * Created by parth on 2/12/2016.
 */

function appendBarChart(data, container) {
    var keys = Object.keys(data);
    var arrayData = [];
    for (var i = 0; i < keys.length; i++) {
        arrayData.push([keys[i], data[keys[i]]]);
    }
    return appendBarChartArray(arrayData, container);
}

function appendBarChartArray(data, container) {
    var leftAxisWidth = 30;
    var width = container.attr('width'),
        height = container.attr('height'),
        radius = Math.min(width, height) / 2;

    var margin = {top: 20 + (height / 100), right: width / 100, bottom: 20 + (height / 100), left: leftAxisWidth + 10 + (width / 100)};
    width = width - margin.left - margin.right;
    height = height - margin.top - margin.bottom;

    var barWidth = width / data.length;
    var barGap = (barWidth * 20) / 100;
    barWidth = barWidth - barGap;
    var max = d3.max(data, function (d) {
        return d[1];
    });

    var y0 = d3.scale.linear().range([height, 0]).domain([0, max]);
    var yAxisLeft = d3.svg.axis().scale(y0)
        .orient("left").ticks(5);

    var valueline = d3.svg.line()
        .x(function(d) { return d[0]; })
        .y(function(d) { return y0(d[1]); });

    container.append("path")        // Add the valueline path.
        .attr("transform", "translate(" + ((margin.left/2) + 14) + "," + margin.top + ")")
        .attr("d", valueline(data));

    container.append("g")
        .attr("transform", "translate(" + ((margin.left/2) + 14) + "," + margin.top + ")")
        .attr("class", "y axis")
        .style("fill", "steelblue")
        .call(yAxisLeft);


    var g = container.append("g")
        .attr("transform", "translate(" + margin.left + "," + margin.top + ")")
        .attr("width", width)
        .attr("height", height);

    var scale = d3.scale.linear().range([0, height]);

    scale.domain([0, max]);
    var color = d3.scale.linear()
        .domain([0, max / 2, max])
        .range(['yellow', 'green', 'red']);


    g.selectAll("rect")
        .data(data)
        .enter().append("svg:rect")
        .attr("fill", function (d, i) {
            return color(d[1])
        })
        .attr("width", function (d, i) {
            return barWidth;
        })
        .attr("height", function (d, i) {
            return scale(d[1])
        })
        .attr("transform", function (d, i) {
            return "translate(" + ((barWidth + barGap) * (i)) + "," + (height - scale(d[1])) + ")";
        });
    container.append("g")
        .attr("transform", "translate(" + margin.left + "," + 10 + ")")
        .attr("width", width)
        .attr("height", margin.bottom)
        .selectAll("text")
        .data(data)
        .enter().append("text")
        .attr("transform", function (d, i) {
            return "translate(" + ((barWidth + barGap) * (i) + 10) + "," + (margin.top + height + 2) + ")rotate(45)"
        })
        .style("font-size", "8px")
        .text(function (d) {
            return d[0];
        })


}
