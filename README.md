<p align="center">
<img 
    src="logo.png" 
    width="240" height="78" border="0" alt="JSONed">
<br>
JSON Stream Editor
</p>

jsoned is a command line utility that provides a [fast](#performance) and simple way to retrieve or update values from JSON documents.
It uses [GJSON](https://github.com/tidwall/gjson) and [SJSON](https://github.com/tidwall/sjson) under the hood. 

It's [fast](#performance) because it avoids parsing irrelevant sections of json, skipping over values that do not apply, and aborts as soon as the target value has been found or updated.

Getting started
---------------

### Building

jsoned can be compiled and used on Linux, OSX, Windows, FreeBSD, and probably others since the codebase is 100% Go. 
We support both 32 bit and 64 bit systems. Go must be installed on the build machine.

To build simply:
```
$ make
```

Or [download a pre-built binary](https://github.com/tidwall/jsoned/releases) for Linux, OSX, Windows, or FreeBSD.


Usage menu:
```
$ jsoned -h

usage: jsoned [-v value] [-r] [-D] [-O] [-p] [-i infile] [-o outfile] keypath

examples: jsoned keypath                      read value from stdin
      or: jsoned -i infile keypath            read value from infile
      or: jsoned -v value keypath             edit value
      or: jsoned -v value -o outfile keypath  edit value and write to outfile

options:
      -v value             Edit JSON key path value
      -r                   Use raw values, otherwise types are auto-detected
      -O                   Performance boost for value updates.
	  -D                   Delete the value at the specified key path
	  -p                   Make json pretty, keypath is optional with this flag
      -i infile            Use input file instead of stdin
      -o outfile           Use output file instead of stdout
      keypath              JSON key path (like "name.last")

```


Examples
--------

### Getting a value 

jsoned uses a special [path syntax](https://github.com/tidwall/gjson#path-syntax) for finding values.

Get a string:
```sh
$ echo '{"name":{"first":"Tom","last":"Smith"}}' | jsoned name.last
Smith
```

Get a block of JSON:
```sh
$ echo '{"name":{"first":"Tom","last":"Smith"}}' | jsoned name
{"first":"Tom","last":"Smith"}
```

Try to get a non-existent key:
```sh
$ echo '{"name":{"first":"Tom","last":"Smith"}}' | jsoned name.middle
null
```

Get the raw string value:
```sh
$ echo '{"name":{"first":"Tom","last":"Smith"}}' | jsoned -r name.last
"Smith"
```

Get an array value by index:
```sh
$ echo '{"friends":["Tom","Jane","Carol"]}' | jsoned friends.1
Jane
```

### Setting a value

The [path syntax](https://github.com/tidwall/sjson#path-syntax) for setting values has a couple of tiny differences than for getting values.

The `-v value` option is auto-detected as a Number, Boolean, Null, or String. 
You can override the auto-detection and input raw JSON by including the `-r` option.
This is useful for raw JSON blocks such as object, arrays, or premarshalled strings.

Update a value:
```sh
$ echo '{"name":{"first":"Tom","last":"Smith"}}' | jsoned -v Andy name.first
{"name":{"first":"Andy","last":"Smith"}}
```

Set a new value:
```sh
$ echo '{"name":{"first":"Tom","last":"Smith"}}' | jsoned -v 46 age
{"age":46,"name":{"first":"Tom","last":"Smith"}}
```

Set a new nested value:
```sh
$ echo '{"name":{"first":"Tom","last":"Smith"}}' | jsoned -v relax task.today
{"task":{"today":"relax"},"name":{"first":"Tom","last":"Smith"}}
```

Replace an array value by index:
```sh
$ echo '{"friends":["Tom","Jane","Carol"]}' | jsoned -v Andy friends.1
{"friends":["Tom","Andy","Carol"]}
```

Append an array:
```sh
$ echo '{"friends":["Tom","Jane","Carol"]}' | jsoned -v Andy friends.-1
{"friends":["Tom","Andy","Carol","Andy"]}
```

Set an array value that's past the bounds:
```sh
$ echo '{"friends":["Tom","Jane","Carol"]}' | jsoned -v Andy friends.5
{"friends":["Tom","Andy","Carol",null,null,"Andy"]}
```

Set a raw block of JSON:
```sh
$ echo '{"name":"Carol"}' | jsoned -r -v '["Tom","Andy"]' friends
{"friends":["Tom","Andy"],"name":"Carol"}
```

Start new JSON document:
```sh
$ echo '' | jsoned -v 'Sam' name.first
{"name":{"first":"Sam"}}
```

### Deleting a value

Delete a value:
```sh
$ echo '{"age":46,"name":{"first":"Tom","last":"Smith"}}' | jsoned -D age
{"name":{"first":"Tom","last":"Smith"}}
```

Delete an array value by index:
```sh
$ echo '{"friends":["Andy","Carol"]}' | ./jsoned -D friends.0
{"friends":["Carol"]}
```

Delete last item in array:
```sh
$ echo '{"friends":["Andy","Carol"]}' | ./jsoned -D friends.-1
{"friends":["Andy"]}
```

### Optimistically update a value

The `-O` option can be used when the caller expects that a value at the
specified keypath already exists.

Using this option can speed up an operation by as much as 6x, but
slow down as much as 20% when the value does not exist.

For example:

```
echo '{"name":{"first":"Tom","last":"Smith"}}' | jsoned -v Tim -O name.first
```

The `-O` tells jsoned that the `name.first` likely exists so try a fasttrack operation first.

## Pretty printing

The `-p` flag will make the output json pretty.

```
$ echo '{"name":{"first":"Tom","last":"Smith"}}' | jsoned -p name
{
  "first": "Tom",
  "last": "Smith"
}
```

Also the keypath is optional when the `-p` flag is specified, allowing for the entire json document to be made pretty.

```
$ echo '{"name":{"first":"Tom","last":"Smith"}}' | jsoned -p
{
  "name": {
    "first": "Tom",
    "last": "Smith"
  }
}
```

## Performance

A quick comparison of jsoned to [jq](https://stedolan.github.io/jq/). The test [json file](https://github.com/zemirco/sf-city-lots-json) is 180MB file of 206,560 city parcels in San Francisco.

*Tested on a 2015 Macbook Pro running jq 1.5 and jsoned 0.2.1*

#### Get the lot number for the parcel at index 10000

jq:

```bash
$ time cat citylots.json | jq -cM .features[10000].properties.LOT_NUM
"091"

real    0m5.486s
user    0m4.870s
sys     0m0.686s
```

jsoned:

```bash
$ time cat citylots.json | jsoned -r features.10000.properties.LOT_NUM
"091"

real    0m0.344s
user    0m0.161s
sys	    0m0.312s
```

#### Update the lot number for the parcel at index 10000

jq:

```bash
$ time cat citylots.json | jq -cM '.features[10000].properties.LOT_NUM="12A"' > /dev/null

real    0m13.579s
user    0m16.484s
sys     0m1.310s
```

jsoned:

```bash
$ time cat citylots.json | jsoned -O -v 12A features.10000.properties.LOT_NUM > /dev/null

real    0m0.635s
user    0m0.343s
sys     0m0.563s
```


## Contact
Josh Baker [@tidwall](http://twitter.com/tidwall)

## License
jsoned source code is available under the MIT [License](/LICENSE).



