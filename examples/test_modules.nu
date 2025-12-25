// Testing module imports with require

console.log("=== Testing Module System ===");
console.log("Current file:", __filename);
console.log("Current directory:", __dirname);

// Import a local module
let math = require("./modules/math_utils.nu");

console.log("\nUsing imported module:");
console.log("math.add(5, 3):", math.add(5, 3));
console.log("math.multiply(4, 7):", math.multiply(4, 7));
console.log("math.PI:", math.PI);
console.log("math.greet('Nulang'):", math.greet("Nulang"));

// Import built-in modules
let fs = require("fs");
let path = require("path");
let crypto = require("crypto");

console.log("\n=== Built-in modules via require ===");

// Use path module
console.log("path.basename(__filename):", path.basename(__filename));
console.log("path.dirname(__filename):", path.dirname(__filename));

// Use crypto module
let hash = crypto.createHash("md5");
hash.update("test");
console.log("MD5 of 'test':", hash.digest("hex"));

console.log("\n✅ Module system working!");
