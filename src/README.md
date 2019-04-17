<h1>Security Project</h1>

* Yo, just run 'go run main.go' to run the server

* Run 'go run cmd/ccnode/.go' to run the node

* Run 'make' to build both the server and the node

* Run 'make ccnode-linux' to build the default node. This can be used for testing in a separate terminal.

* Run 'make ccnode-android' to build an android emulator compatible version, but this isn't necessary for testing because if it compiles it'll work. Do all testing in the linux binary.

<h2>What's done now</h2>

* Everything is routed through a single endpoint

* C2 can create new nodes

* C2 can send jobs to specific nodes or all nodes

* C2 can steal files from a node (all nodes is untested)

* C2 can upload files to a node

<h2>Go Tips</h2>

Go is pythonic c++, most of the basic stuff is already implementated. The code for the Node is in node.go and that is consumed by both the server and cmd/ccnode/main.go.

* fmt.Println("string") is the same as cout with an endl in cpp

* All json is serialized, this means it's converted to binary before being sent. You'll need a marshaller and unmarsheller. Try to avoid in the beginning or message me for help. All json responses can be found in messages/messages.go

* Go is strongly-typed. Everything needs a type: int, int8, string, struct, bool, etc.

* Classes do not explicitly exist. A class in Go is a package. So package main or package messages at the top of the file designate what package it belongs to. All variables within the package belong to that class essentially.

* If you want a function or variable or const to be accessible from outside that package then the first letter should be capatilized, e.g. look in messages/messages.go for all exported msg types.

* Array designators come before the type, e.g. []string instead of string[]

* You can either create a new variable with 'var new = 128' or 'new := 128'

* Functions can have multiple returns. If you see '_, err := foo()' this just means that the first return value is completely ignored.
