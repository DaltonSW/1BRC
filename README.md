# 1BRC

## Info

Running the 1 Billion Row Challenge in a few different languages, to faciliate learning optimization in a few different languages.

The script `createMeasurements.py` will create the measurement file:
```
usage: createMeasurements.py [-h] [-o OUTPUT] [-r RECORDS]

Create measurement file

optional arguments:
  -h, --help            show this help message and exit
  -o OUTPUT, --output OUTPUT
                        Measurement file name (default is "measurements.txt")
  -r RECORDS, --records RECORDS
                        Number of records to create (default is 1_000_000_000)
Shoutouts to github.com/ifnesi for the script
```

You'll also need the following Python modules: `numpy`, `polars`, and `tqdm`

## Rules 

This snippet is taken from the original repo.

The text file contains temperature values for a range of weather stations.
Each row is one measurement in the format `<string: station name>;<double: measurement>`, with the measurement value having exactly one fractional digit.
The following shows ten rows as an example:

```
Hamburg;12.0
Bulawayo;8.9
Palembang;38.8
St. John's;15.2
Cracow;12.6
Bridgetown;26.9
Istanbul;6.2
Roseau;34.4
Conakry;31.2
Istanbul;23.0
```

The task is to write a program which reads the file, calculates the min, mean, and max temperature value per weather station, and emits the results sorted alphabetically by station name, and the result values per station in the format `<min>/<mean>/<max>`, rounded to one fractional digit

## Test Environment 

Unless otherwise stated, all of these tests are performed in the environment below: 

OS: WSL 2 Ubuntu on Windows 11
CPU: Intel i9-12900K (24 Cores) @ 3.19 GHz
RAM: 32GB @ 3200MHz (Though WSL only reports 16GB, so maybe that)
