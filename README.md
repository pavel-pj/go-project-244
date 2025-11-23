### Hexlet tests and linter status:
[![Actions Status](https://github.com/pavel-pj/go-project-244/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/pavel-pj/go-project-244/actions)
### Project tests and linter status:
[![tests & lint](https://github.com/pavel-pj/go-project-244/actions/workflows/main.yml/badge.svg)](https://github.com/pavel-pj/go-project-244/actions/workflows/main.yml)


# GOlang Console utilite "GenDiff" 

### Compares two files (json,yaml) and shows differencies in 3 types formats.


[![asciicast](https://asciinema.org/a/757726.svg)](https://asciinema.org/a/757726)

#### file01.json:
```
{
"common": {
  "setting1": "Value 1",
  "setting2": 200,
  "setting3": true,
  "setting6": {
    "key": "value",
    "doge": {
      "wow": ""
    }
  }
},
"group1": {
  "baz": "bas",
  "foo": "bar",
  "nest": {
    "key": "value"
  }
},
"group2": {
  "abc": 12345,
  "deep": {
    "id": 45
  }
}
}
```

### file02.json:
```
{
"common": {
  "follow": false,
  "setting1": "Value 1",
  "setting3": null,
  "setting4": "blah blah",
  "setting5": {
    "key5": "value5"
  },
  "setting6": {
    "key": "value",
    "ops": "vops",
    "doge": {
      "wow": "so much"
    }
  }
},
"group1": {
  "foo": "bar",
  "baz": "bars",
  "nest": "str"
},
"group3": {
  "deep": {
    "id": {
      "number": 45
    }
  },
  "fee": 100500
}
}
```

### result :
```
./bin/gendiff file1.json file2.json
{
  common: {
    + follow: false  
      setting1: Value 1
    - setting2: 200  
    - setting3: true  
    + setting3: null  
    + setting4: blah blah
    + setting5: {
          key5: value5
      }
      setting6: {
          doge: {
            - wow: 
            + wow: so much
          }
          key: value
        + ops: vops
      }
  }
  group1: {
    - baz: bas
    + baz: bars
      foo: bar
    - nest: {
          key: value
      }
    + nest: str
  }
- group2: {
      abc: 12345
      deep: {
          id: 45
      }
  }
+ group3: {
      deep: {
          id: {
              number: 45
          }
      }
      fee: 100500
  }
}
```