# DigoLang

Toy programming language based of the Thorsten Ball's incredible book, [Writing an Interpeter in Go](https://interpreterbook.com).

```
let a = 1;
let b = 1;

let sum = fn(a, b) {a+b};

let result = sum(5, 5);

let max = fn(x, y) {
    if(x > y) {
        return x;
    }

    return y;
};

fn(x) {x==5}(9)

let whichIsMax = max(10, 20);

return whichIsMax;

```