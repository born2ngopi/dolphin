# Dolpin 

> still development

Dolpin is auto generate test case and very simple to use it


it's generate with llama on local machine 

## Problem & Motivation

write test case is very boring and take a lot of time, so i create this tool to generate test case

## to run

you just run command

```shell
dolpin generate -r="." --mock-path="./mock" --mock-lib="gomock"
```

if you want to generate test case for spesific dir
    
```shell    
dolpin generate -r="." --dir="./user" --mock-path="./mock" --mock-lib="gomock"
```

or for spesific function

```shell
dolpin generate -r="." --mock-path="./mock" --mock-lib="gomock" --func="TestFunc" --file="./somefolder/test.go"
```


## example

### prompt result

with example code on [this](https://github.com/born2ngopi/example-dolpin)

and if we run command
```shell
dolpin generate --dir="." --mock-path="./mock" --mock-lib="gomock"
```

then we got prompt
``` text
can u write unit test on golang with heights coverage and multi scenario for this code

package check

import (
	"strconv"

	"github.com/born2ngopi/example-dolpin/types"
)

type Coba struct {
}

type Check interface {
	CheckFunction(msg types.Message, randNumber int) string
}

func NewCheck() Check {
	return &Coba{}
}

func (c *Coba) CheckFunction(msg types.Message, randNumber int) string {

	msg.Name = "CheckFunction"

	return msg.Name + strconv.Itoa(randNumber)
}


and i have some struct like this

type Message struct {
	Name string
	Status string
}
from "github.com/born2ngopi/example-dolpin/types"

and i use mock gomock and the dir is github.com/born2ngopi/example-dolpin/mock


i expect the unit test like this
func Test_[function_name](t *testing.T) {

	// add some preparation code here

	// add schenario here with []struct

	// looping schenario here and test the function
}
```

if you have some idea to improve this tool or better prompt, please create issue or pull request
