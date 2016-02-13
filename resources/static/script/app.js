/**
 * Created by parth on 2/12/2016.
 */

var height = 400, width = 400;

$(document).ready(function () {
    $.ajax({
        url: baseUrl + 'info',
        success: function (d) {
            var data = JSON.parse(d);
            console.log("Data: ", data);
            d3.select("#columns")
                .append("ul")
                .selectAll("li")
                .data(data.ColumnInfo)
                .enter()
                .append("li")
                .text(function (d) {
                    return d.ColumnName + " - " + d.TypeInfo + (  d.IsEnum ? "(Enum)" : ""  );
                });

            k_combinations(data.ColumnInfo, 1).forEach(function (colInfo) {
                    var col = colInfo[0];
                    console.log("try - ", col);
                    var keys = [
                        col.TypeInfo + "&" + col.IsEnum.toString(),
                        "&" + col.IsEnum.toString()
                    ];
                    var f = getFunction(keys);
                    if (f == undefined) {
                        return
                    }
                    f = f(col);
                    var container = addContainer(f.height(height), f.width(width), f.columnName());
                    f(col.ValueCounts, container)
                }
            )
        }
    });

    function getFunction(keys) {
        var f = undefined;
        for (var x = 0; x < keys.length; x++) {
            var mapKey = keys[x];
            f = functionMap[mapKey];
            if (f != undefined) {
                break
            }
        }
        return f;
    }

    var functionMap = {
        '&true': function (colInfo) {
            var x;
            if (colInfo.DistinctValueCount > 15) {
                x = appendBarChart;
                x.height = function (h) {
                    return h * 1.3;
                };
                x.width = function (w) {
                    return w * 3;
                };
            } else {
                x = appendPieChartByMap;
                x.height = function (h) {
                    return h;
                };
                x.width = function (w) {
                    return w;
                };
            }
            x.columnName = function () {
                return colInfo.ColumnName
            };
            return x;
        },
        'number&true': function (colInfo) {
            var x;
            if (colInfo.DistinctValueCount < 5) {
                x = appendPieChartByMap;
                x.height = function (h) {
                    return h;
                };
                x.width = function (w) {
                    return w;
                };
            } else {
                x = appendBarChart;
                x.height = function (h) {
                    return h * 1.3;
                };
                x.width = function (w) {
                    return w * 3;
                };
            }
            x.columnName = function () {
                return colInfo.ColumnName
            };
            return x;
        }

    };

    function addContainer(height, width, name) {
        var times = width / window.width;
        console.log("times", width, height);
        var col = $("<div></div>");
        col.attr("id", "container-" + name);
        col.css("width", width + "px");
        col.css("float", "left");
        col.css("height", height + "px");
        col.append("<span>" + name + "</span>");
        $("#chart").append(col);
        return d3.select("#container-" + name)
            .append("svg:svg")
            .style('height', height)
            .style('width', width)
            .attr('id', name);
    }
});