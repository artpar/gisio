# Gisio

[![Build Status](https://travis-ci.org/artpar/gisio.svg?branch=master)](https://travis-ci.org/artpar/gisio)

Do you see a lot of CSVs around you but struggling your way to get the juice out of those ?

Just upload a CSV and get awesome automatic charts and graphs

![Gisio demo](https://github.com/artpar/gisio/raw/master/resources/static/images/gisio.png)

- Click on a column name to see the graphs/charts associated with that column

## Docker

I push latest builds to docker hub. You can always pull the latest image and run locally

```docker run -p 2299:2299 -v <path-to-csv-folder>:/opt/gisio/data gisio/gisio```


There is no "Upload CSV", put yours csv files in folder ```<path-to-csv-folder>```. New files will show up without restart.

## Build

```go get github.com/artpar/gisio```

## Run

```go run main.go```

- Put your CSV files inside the ```./data/``` folder
- Open http://localhost:2299/

## Work In Progress
