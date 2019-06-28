# Introduction

First of all, simple_dependency_finder(sdf) is built to figure out **what test suits to run**.

Big Project has a lot of unit tests. It takes long time to run all of test suits.

Some possible solution for faster test execution is..

1. Run chunk of test on multiple CPU
2. Run chunk of test on different machine
3. Make less IO by refactor test files
4. Run test suits as little as possible

sdf helps to do **4. Run test suits as little as possible**

# Install

TODO

# How it works

## Graph
![dependency_graph](https://user-images.githubusercontent.com/11783802/60318268-58387e00-99ad-11e9-8103-ac2becd76393.png)


Above graph show dependency graph.

Each Node represents **main module unit of any Programing Language**. For example, class for Ruby, function for JS, package for Go, etc.

Arrow Between Nodes represents **Which module knows which module**.
For example, `OrderItem` knows both `Order` and `Item` directly, and knows `User` Indirectly through `Order`

Below table shows **what to test** or **what is effected** based on change on this dependency graph.

Basically, we need to test only effected modules.

|diff|what to test|
|----|------------|
|User|User, Order, OrderItem|
|Order|Order, OrderItem|
|Item|Item, OrderItem|
|OrderItem|OrderItem|
|CircularA|CircularA, CircularB|
|CircularC|CircularC, CircularD, CircularE|

## Function

sdf is cli tool, which is like function below.

Outputed effected modules are **what we need to test**.

It's important to understand that feeding these input is your responsibility.

How to parse code is very different based on context like Programing Language.
So sdf does not care about it by design.

![sdf_func](https://user-images.githubusercontent.com/11783802/60318569-92eee600-99ae-11e9-9090-bb9cb54a72fa.png)

cli parameters...

|option name|desc|required|default|type|
|-----------|----|--------|-------|----|
|mpath|path to json file which has modules info|true|NA|string|
|rpath|path to json file which has relation info|true|NA|string|
|diffs|file diffs comma separated|true|NA|string|
|depth|how deep you want to search|false|0(unlimited)|int|

each json scheme will be found in `testdata` directory

example usage

``` shell
sdf \
  -mpath testdata/modules.json \
  -rpath testdata/relations.json \
  -diffs order.rb
```

# Example Usage

## Case: Ruby project

```
sdf \
  -mpath testdata/modules.json \
  -rpath testdata/relations.json \
  -diffs $(git ls-files -omd --exclude-standard | sed "s,^,$PWD/,") \
  | jq -r '.[] | .name' > patterns.txt

grep -rl -f patterns.txt spec/ | xargs rspec
rm patterns.txt
```
