# Dolpin 

Dolpin is auto generate test case and very simple to use it


It's generate with llama on local machine. The process will be depend with your hardware

## Demo

Hardware for testing:
* Processor AMD Ryzen 7 7840HS
* GPU Nvidia RTX 4050 6GB laptop
* RAM 16GB

Source code for test [here](https://github.com/born2ngopi/example-dolpin) and result in [here](https://github.com/born2ngopi/example-dolpin/tree/generate-test)

[![WATCH](https://img.youtube.com/vi/oRNMYKI5nR8/hqdefault.jpg)](https://youtu.be/oRNMYKI5nR8)


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
dolpin generate --dir="." --mock-path="./mocks" --mock-lib="gomock" --model="codegemma:7b"
```

then we got prompt
``` text
WRITE UNIT TEST GOLANG

can u write unit test on golang with heights coverage and multi scenario for this code

package service

import (
        "strconv"

        "github.com/born2ngopi/example-dolpin/types"
)

type service struct {
}

type Service interface {
        Sum(a, b int) int
        SumFromStr(data types.SumField) (int, error)
}

func New() Service {
        return &service{}
}

func (s *service) Sum(a, b int) int {
        return a + b
}

func (s *service) SumFromStr(data types.SumField) (int, error) {
        aInt, err := strconv.Atoi(data.A)
        if err != nil {
                return 0, err
        }

        bInt, err := strconv.Atoi(data.B)
        if err != nil {
                return 0, err
        }

        return aInt + bInt, nil
}


and i have some struct like this


type SumField struct {

        A string

        B string

}
from "github.com/born2ngopi/example-dolpin/types"

and i use mock gomock and the dir is github.com/born2ngopi/example-dolpin/mocks

i expect the unit test like this
func Test_[function_name](t *testing.T) {

        // add some preparation code here include mock, var, and etc

        // add schenario here with []struct
        /*
                example:
                type arg struct {
                        // this field must be parameter function
                }

                tests := []struct{
                        name string
                        arg arg // arg is parameter function,
                        wantError error
                        wantResponse [response function]
                        prepare func([parameter function]) // prepare for expected mock function
                }{
                        {
                                // fill hire with success scenario and posibility negative/error scenario
                        }
                }
        /*

        // looping schenario here and test the function
        /*
                example:
                for _, tt := range tests {
                        t.Run(tt.name, func(t *testing.T){
                                // some test logic here
                        })
                }
        /*
}
```

if you have some idea to improve this tool or better prompt, please create issue or pull request
