// Module that exports utilities
// usando exports CommonJS style

exports.add = function(a, b) {
    return a + b;
};

exports.multiply = function(a, b) {
    return a * b;
};

exports.PI = 3.14159;

exports.greet = function(name) {
    return "Hello, " + name + "!";
};

console.log("math_utils module loaded");
