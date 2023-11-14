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

