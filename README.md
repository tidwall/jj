<p align="center">
<img 
    src="logo.png" 
    width="240" height="78" border="0" alt="JSONed">
<br>
JSON Stream Editor
</p>

jsoned is a command line utility that provides a fast and simple way to retrieve or update values from JSON documents.
It uses [GJSON](https://github.com/tidwall/gjson) and [SJSON](https://github.com/tidwall/sjson) under the hood. 

It's fast because it avoids parsing irrelevant sections of json, skipping over values that do not apply, and aborts as soon as the target value has been found or updated.

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

usage: jsoned [-v value] [-s] [-i infile] [-o outfile] keypath

examples: jsoned keypath                      read value from stdin
      or: jsoned -i infile keypath            read value from infile
      or: jsoned -v value keypath             edit value
      or: jsoned -v value -o outfile keypath  edit value and write to outfile

options:
      -v value             Edit JSON key path value
      -i infile            Use input file instead of stdin
      -o outfile           Use output file instead of stdout
      -r                   Use raw values, otherwise types are auto-detected
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
{"age":46,"name":{"first":"Andy","last":"Smith"}}
```

Set a new nested value:
```sh
$ echo '{"name":{"first":"Tom","last":"Smith"}}' | jsoned -v relax task.today
{"task":{"today":"relax"},"name":{"first":"Andy","last":"Smith"}}
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

## Contact
Josh Baker [@tidwall](http://twitter.com/tidwall)

## License
jsoned source code is available under the MIT [License](/LICENSE).



