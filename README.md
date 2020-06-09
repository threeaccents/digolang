# DigoLang

Toy programming language based of the Thorsten Ball's incredible book, [Writing an Interpeter in Go](https://interpreterbook.com).

```
let greeter = fn(name) {
    return "hello" + " " + name
}


let msg = greeter("world")

println(msg)

let a = "hello"

let b;

if (isNull(b)) {
    println("b IS NULL!")
} else {
    println("b NOT NULL", "it's", b)
}

if (isNull(a)) {
    println("a IS NULL!")
} else {
    println("a NOT NULL", "it's", a)
}

let add = fn(x) {
    return fn(y) {
        x + y;
    }
}

let adder = add(2)

let four = adder(2)

println(four)

let howLong = "how long is this string?"

println(howLong, len(howLong))

let myarr = [1,2,3]

println(myarr[0])

println(myarr.len)
println(myarr.first)
println(myarr.last)
println(myarr.rest)
```