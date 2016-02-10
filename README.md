# qcli
Lightweight command line interface flag handling

This is a project that I wanted to make because I was a bit tired of CLI helper libraries trying to do too much.
All I really needed was just an easier way to setup flags, and then get their values :P Thus, this package was born.
Also I like the idea of organizing my flags in JSON format.

## Usage
First, you'll need to setup your flag organization in a `flags.json` file. If you don't like the filename, too
bad (or you could just modify the source if it REALLY bothers you). An example structure:
```JSON
{
"flags":[{
  "name":"force",
  "type":"bool",
  "default":false,
  "usage":"Forces something..."
  }]
}
```

This file says make a flag named `force` of type `bool` and a default value of `false` with a very unhelpful usage description
This file also does the same thing:
```JSON
{
"flags":[{
  "name":"force",
  "default":false,
  "usage":"Forces something..."
  }]
}
```

So does this one:
```JSON
{
"flags":[{
  "name":"force",
  "type":"bool",
  "usage":"Forces something..."
  }]
}
```

It's okay if you omit either the type or the default value. The library will try to do its best to figure it out on init().
Now to retrieve the value of this flag, we can make a simple golang application:

```go
package main

import (
  "fmt"
  
  "github.com/corvuscrypto/qcli"
)

func main(){
  fmt.Println(qcli.Flags.Get("force"))
}
```

If you run that program with the flag --force or -force, you will see the output of `true`. Otherwise, it will return `false`.
And that's what I call easy and q(uick)cli flag handling.
