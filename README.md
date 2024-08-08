# Dolpin 

Dolpin is auto generate test case for golang and very simple to use it.


It's generate with llama on local machine (this verry dependend on your machine), and also can generate with google gemini and openai.
## Demo With Llama

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

with ollama
```shell
dolpin generate -r="." --mock-path="./mock" --mock-lib="gomock" --llm="ollama" --model="codegemma:7b"
```

if you have other host, you can use this command
```shell
dolpin generate -r="." --mock-path="./mock" --mock-lib="gomock" --llm="ollama" --model="codegemma:7b" --llm-host="https://yourhost.com"
```

with gemini
```shell
dolpin generate -r="." --mock-path="./mock" --mock-lib="gomock" --llm="gemini" --llm-token="abc" --model="gemini:7b"
```

with openai
```shell
dolpin generate -r="." --mock-path="./mocks" --mock-lib="gomock" --llm="openai" --model="gpt-4o-mini" --llm-token="abc"
```

if you want to generate test case for spesific dir
    
```shell    
dolpin generate -r="." --dir="./user" --mock-path="./mock" --mock-lib="gomock"
```

or for spesific function

```shell
dolpin generate -r="." --mock-path="./mock" --mock-lib="gomock" --func="TestFunc" --file="./somefolder/test.go"
```


