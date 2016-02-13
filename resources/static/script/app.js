/**
 * Created by parth on 2/12/2016.
 */

var height = 300, width = 300;

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

            for (var i = 0; i < data.ColumnInfo.length; i++) {
                var keys = [
                    data.ColumnInfo[i].TypeInfo + "&" + data.ColumnInfo[i].IsEnum.toString(),
                    "&" + data.ColumnInfo[i].IsEnum.toString()
                ];

                var f = undefined;
                for (var x = 0; x < keys.length; x++) {
                    var mapKey = keys[x];
                    f = functionMap[mapKey];
                    if (f != undefined) {
                        break
                    }
                }
                if (f == undefined) {
                    continue
                }
                f = f(data.ColumnInfo[i]);

                var container = addContainer(f.height(height), f.width(width), f.columnName());

                f(data.ColumnInfo[i].ValueCounts, container)
            }
        }
    });

    var functionMap = {
        '&true': function (colInfo) {
            var x;
            if (colInfo.DistinctValueCount > 15) {
                x = appendBarChart;
                x.height = function (h) {
                    return h;
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
                    return h;
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
        console.log("times", times)
        var col = $("<div class='col-md-" + (3 * times) + "'></div>");
        col.attr("id", "container-" + name);
        col.append("<span>" + name + "</span>");
        $("#chart").append(col);
        var container = d3.select("#container-" + name)
            .append("svg:svg")
            .style('height', height)
            .style('width', width)
            .attr('id', name);

        return container;
    }
});