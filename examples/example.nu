// Nulang Example - JavaScript-like syntax

// Variables
let name = "Nulang";
const version = "0.1.0";
var greeting = "Hello";

// String operations
console.log(greeting + ", " + name + " v" + version + "!");

// Numbers and Math
let x = 10;
let y = 3;
console.log("Math operations:");
console.log("x + y =", x + y);
console.log("x - y =", x - y);
console.log("x * y =", x * y);
console.log("x / y =", x / y);
console.log("x % y =", x % y);
console.log("x ** y =", x ** y);

// Arrays
let numbers = [1, 2, 3, 4, 5];
console.log("\nArray:", numbers);
console.log("Length:", numbers.length);

// Array methods
let doubled = numbers.map(n => n * 2);
console.log("Doubled:", doubled);

let evens = numbers.filter(n => n % 2 === 0);
console.log("Evens:", evens);

let sum = numbers.reduce((acc, n) => acc + n, 0);
console.log("Sum:", sum);

// Objects
let person = {
    name: "John",
    age: 30,
    city: "New York"
};
console.log("\nPerson:", person);
console.log("Name:", person.name);

// Functions
function add(a, b) {
    return a + b;
}
console.log("\nadd(5, 3) =", add(5, 3));

// Arrow functions
const multiply = (a, b) => a * b;
console.log("multiply(4, 7) =", multiply(4, 7));

// Control flow
console.log("\nFizzBuzz from 1 to 15:");
for (let i = 1; i <= 15; i++) {
    if (i % 15 === 0) {
        console.log("FizzBuzz");
    } else if (i % 3 === 0) {
        console.log("Fizz");
    } else if (i % 5 === 0) {
        console.log("Buzz");
    } else {
        console.log(i);
    }
}

// Ternary operator
let age = 20;
let status = age >= 18 ? "adult" : "minor";
console.log("\nAge", age, "is:", status);

// String methods
let text = "  Hello, World!  ";
console.log("\nString methods:");
console.log("Original:", "'" + text + "'");
console.log("Trimmed:", "'" + text.trim() + "'");
console.log("Upper:", text.toUpperCase());
console.log("Lower:", text.toLowerCase());
console.log("Split:", "hello-world".split("-"));

// Math object
console.log("\nMath functions:");
console.log("Math.PI:", Math.PI);
console.log("Math.sqrt(16):", Math.sqrt(16));
console.log("Math.pow(2, 8):", Math.pow(2, 8));
console.log("Math.floor(3.7):", Math.floor(3.7));
console.log("Math.ceil(3.2):", Math.ceil(3.2));
console.log("Math.round(3.5):", Math.round(3.5));
console.log("Math.max(1, 5, 3):", Math.max(1, 5, 3));
console.log("Math.min(1, 5, 3):", Math.min(1, 5, 3));

// While loop
console.log("\nWhile loop countdown:");
let count = 5;
while (count > 0) {
    console.log(count);
    count--;
}
console.log("Blast off!");

// Typeof
console.log("\nTypeof:");
console.log("typeof 42:", typeof 42);
console.log("typeof 'hello':", typeof "hello");
console.log("typeof true:", typeof true);
console.log("typeof null:", typeof null);
console.log("typeof undefined:", typeof undefined);
console.log("typeof []:", typeof []);
console.log("typeof {}:", typeof {});
console.log("typeof function(){}:", typeof function(){});

// Nullish coalescing
let value = null;
let defaultValue = value ?? "default";
console.log("\nNullish coalescing:", defaultValue);

// Optional chaining simulation
let obj = { nested: { value: 42 } };
console.log("Nested value:", obj.nested.value);

// Try/catch
console.log("\nTry/catch:");
try {
    throw "This is an error!";
} catch (e) {
    console.log("Caught:", e);
}

console.log("\n✨ Nulang is working! ✨");
