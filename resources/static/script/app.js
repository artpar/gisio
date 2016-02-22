/**
 * Created by parth on 2/12/2016.
 */

var height = 300, width = 300;
var chartCount = 0;

function generateGraphs(data, columnName) {
    console.log("plot charts", columnName);
    reset();
    k_combinations(data.ColumnInfo, 1).forEach(function (colInfo) {
            var col = colInfo[0];
            if (col.ColumnName != columnName || col.ColumnName == "") {
                return;
            }
            console.log(col.ColumnName);
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
    );

    k_combinations(data.ColumnInfo, 2).forEach(function (cols) {
            permutations(cols).forEach(function (colInfo) {
                //console.log("combinations of 2", colInfo);
                var colX = colInfo[0];
                var colY = colInfo[1];
                if (colX.ColumnName != columnName && colY.ColumnName != columnName || (colX.ColumnName == "" || colY.ColumnName == "") ) {
                    return;
                }
                if (colY.TypeInfo == "number" && !colY.IsUnique && !colY.IsEnum) {
                    $.ajax({
                        url: "operation",
                        data: {
                            'q': JSON.stringify({
                                "operation": "GroupBy",
                                "function": "sum",
                                "data": [
                                    {
                                        "ColumnName": colX.ColumnName
                                    },
                                    {
                                        "ColumnName": colY.ColumnName
                                    }
                                ]
                            })
                        },
                        success: function (d) {
                            console.log("2d chart for " + colX.ColumnName + " vs. " + colY.ColumnName, d);
                            var f = appendBarChart;
                            var h = height * 1.3;
                            var w = width * 3;
                            if (d.length < 7) {
                                h = height;
                                w = width;
                                f = appendPieChart
                            }
                            if (colX.TypeInfo == "number" && d.length > 40) {
                                f = appendScatterChart;
                            }
                            if (colX.TypeInfo == "date") {
                                for (var i = 0; i < d.length; i++) {
                                    d[i][0] = new Date(d[i][0])
                                }
                                f = appendAreaChart
                            }
                            var container = addContainer(h, w, colX.ColumnName + " vs. " + colY.ColumnName);
                            f(d, container, false)
                        }
                    })
                }
                //else if (colX.TypeInfo == "number" && colY.TypeInfo == "number") {
                //
                //}
                else if (colX.IsEnum && !colY.IsUnique) {
                    $.ajax({
                        url: "operation",
                        data: {
                            'q': JSON.stringify({
                                "operation": "GroupBy",
                                "function": "count",
                                "data": [
                                    {
                                        "ColumnName": colX.ColumnName
                                    },
                                    {
                                        "ColumnName": colY.ColumnName
                                    }
                                ]
                            })
                        },
                        success: function (d) {
                            var f = appendPieChart;
                            var h = height;
                            var w = width;
                            if (d.length > 6) {
                                h = height * 1.3;
                                w = width * 3;
                                f = appendBarChart
                            }
                            var container = addContainer(h, w, colX.ColumnName + " vs. " + colY.ColumnName);
                            console.log("Bar chart 2 for " + colX.ColumnName + " vs. " + colY.ColumnName, d);
                            f(d, container, false)
                        }
                    })
                }


            })
        }
    );
}
var functionMap = {
    '&true': function (colInfo) {
        var x;
        if (colInfo.DistinctValueCount > 7) {
            x = appendBarChart;
            x.height = function (h) {
                return h * 1.3;
            };
            x.width = function (w) {
                return w * 3;
            };
        } else {
            x = appendPieChart;
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

};

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

function reset() {
    $("#chart").html("");
}

function addContainer(height, width, name) {
    chartCount = chartCount + 1;
    var times = width / window.width;
    //console.log("times", width, height);
    var col = $("<div></div>");
    var cleanName = clean(name + chartCount);
    col.attr("id", "container-" + cleanName);
    col.css("width", width + "px");
    col.css("float", "left");
    col.css("height", height + "px");
    col.append("<p class='text-primary bg-info'>" + name + " Chart " + chartCount + "</p>");
    $("#chart").append(col);
    return d3.select("#container-" + cleanName)
        .append("svg:svg")
        .style('height', height)
        .style('width', width)
        .attr('id', cleanName);
}

$(document).ready(function () {
    $.ajax({
        //url: "http://localhost:2299/data/Catsup.csv/info",
        url: baseUrl + 'info',
        success: function (d) {
            var data = d;
            console.log("Data: ", data);
            d3.select("#columns")
                .append("div")
                .selectAll("div")
                .data(data.ColumnInfo)
                .enter()
                .append("div")
                .attr("class", "well")
                .attr("id", function (d) {
                    return "box-" + clean(d.ColumnName)
                })
                .on("click", function (d) {
                    console.log("clicked", d);
                    generateGraphs(data, d.ColumnName);
                })
                .text(function (d) {
                    return d.ColumnName + " - " + d.TypeInfo + (  d.IsEnum ? "(Enum)" : d.IsUnique ? "(Unique)" : ""  );
                });


        }
    });
});

function clean(name) {
    return name.replace(/[^a-zA-Z0-9]/g, '')
}