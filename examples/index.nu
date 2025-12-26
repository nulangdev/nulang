// Test 1: import with alias { a as b }
import { add, multiply as mult } from "./add.nu";

console.log("add(10, 20) =", add(10, 20));
console.log("mult(10, 20) =", mult(10, 20));

// Test 2: import * as namespace
import * as mathModule from "./add.nu";

console.log("mathModule.add(5, 3) =", mathModule.add(5, 3));
console.log("mathModule.multiply(5, 3) =", mathModule.multiply(5, 3));
