/**
 * Created by parth on 2/12/2016.
 */

var pieHeight = 200, pieWidth = 200;


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


                if (!data.ColumnInfo[i].IsEnum) {
                    continue
                }

                var f = appendPieChartByMap;
                var width = pieWidth, height = pieHeight;

                if (data.ColumnInfo[i].DistinctValueCount > 15) {
                    f = appendBarChart;
                    width = width * 3;
//                        height = height * 1;
                }

                var columnName = data.ColumnInfo[i].ColumnName;
                var container = d3.select("#pies")
                    .append("svg:svg")
                    .attr('height', height)
                    .attr('width', width)
                    .attr('id', columnName);

                (function (name) {
                    container.append("text")
                        .attr("dy", ".35em")
                        .attr("transform", "translate(30,10)")
                        .text(function (d) {
                            return name;
                        });
                })(columnName);

                f(data.ColumnInfo[i].ValueCounts, container)

            }
        }
    })
})