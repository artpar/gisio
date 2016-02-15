/**
 * Created by parth on 2/12/2016.
 */

var height = 300, width = 300;

$(document).ready(function () {
    $.ajax({
        //url: "http://localhost:2299/data/Catsup.csv/info",
        url: baseUrl + 'info',
        success: function (d) {
            var data = JSON.parse(d);
            console.log("Data: ", data);
            d3.select("#columns")
                .append("div")
                .selectAll("pre")
                .data(data.ColumnInfo)
                .enter()
                .append("pre")
                .text(function (d) {
                    return d.ColumnName + " - " + d.TypeInfo + (  d.IsEnum ? "(Enum)" : d.IsUnique ? "(Unique)" : ""  );
                });

            k_combinations(data.ColumnInfo, 1).forEach(function (colInfo) {
                    var col = colInfo[0];
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
                    // console.log(f.columnName(), " for ", f);
                    f(col.ValueCounts, container)
                }
            );

            k_combinations(data.ColumnInfo, 2).forEach(function (cols) {
                    permutations(cols).forEach(function (colInfo) {
                        console.log("combinations of 2", colInfo);
                        var colX = colInfo[0];
                        var colY = colInfo[1];
                        if (colX.TypeInfo == "number" && colX.IsUnique && colY.TypeInfo == "number" && !colY.IsUnique && !colY.IsEnum) {
                            var f = appendAreaChart;
                            var container = addContainer(height * 1.3, width * 3, colX.ColumnName + " vs. " + colY.ColumnName);
                            $.ajax({
                                url: "operation",
                                data: {
                                    "operation": "MapColumnValue",
                                    "data": [
                                        {
                                            "ColumnName": colX.ColumnName
                                        },
                                        {
                                            "ColumnName": colY.ColumnName
                                        }
                                    ]
                                },
                                success: function (d) {
                                    f(d, container)
                                }
                            })
                        }
                        //f = f(colX);
                        //var container = addContainer(f.height(height), f.width(width), f.columnName());
                        //console.log(f.columnName(), " for ", f);
                        //f(colX.ValueCounts, container)
                    })
                }
            );


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
            if (colInfo.DistinctValueCount > 10) {
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
            if (colInfo.DistinctValueCount < 7) {
                x = appendLineChart;
                x.height = function (h) {
                    return h;
                };
                x.width = function (w) {
                    return w;
                };
            } else {
                x = appendLineChart;
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
        //'date&true': function (colInfo) {
        //    var x;
        //    if (colInfo.DistinctValueCount < 7) {
        //        x = appendLineChart;
        //        x.height = function (h) {
        //            return h;
        //        };
        //        x.width = function (w) {
        //            return w;
        //        };
        //    } else {
        //        x = appendLineChart;
        //        x.height = function (h) {
        //            return h * 1.3;
        //        };
        //        x.width = function (w) {
        //            return w * 3;
        //        };
        //    }
        //    x.columnName = function () {
        //        return colInfo.ColumnName
        //    };
        //    return x;
        //}

    };

    function addContainer(height, width, name) {
        var times = width / window.width;
        //console.log("times", width, height);
        var col = $("<div></div>");
        var cleanName = clean(name);
        col.attr("id", "container-" + cleanName);
        col.css("width", width + "px");
        col.css("float", "left");
        col.css("height", height + "px");
        col.append("<span>" + name + "</span>");
        $("#chart").append(col);
        return d3.select("#container-" + cleanName)
            .append("svg:svg")
            .style('height', height)
            .style('width', width)
            .attr('id', cleanName);
    }
});

function clean(name) {
    return name.replace(/[^a-zA-Z0-9]/g, '')
}