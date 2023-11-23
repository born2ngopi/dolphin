# Dolpin 

> still development

Dolpin is auto generate test case and very simple to use it


it's generate with llama on local machine 

## Problem & Motivation

write test case is very boring and take a lot of time, so i create this tool to generate test case

## to run

you just run command

```shell
dolpin generate --dir="." --mock-path="./mock" --mock-lib="gomock"
```

or for spesific function

```shell
dolpin generate --dir="." --mock-path="./mock" --mock-lib="gomock" --func="TestFunc" --file="./somefolder/test.go"
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
i have a function like this

func CheckFunction(msg types.Message, check Coba) string {

	msg.Name = "CheckFunction"

	return msg.Name
}


and i have a struct like this


type Message struct {

  Name string

  Status string

}

type Coba struct {

  Name string

}




and i use mock gomock and the dir is ./mock


can you write unit test with heights coverage and looping test case , So there can be looping a positive case and a negative case for this function . And only return to me the function unit test without package name and import package?
```

