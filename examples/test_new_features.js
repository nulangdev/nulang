// Test file for new Nulang features
// =================================

// TEST 1: switch/case statement
console.log("=== Testing switch/case ===");
let fruit = "apple";
let result = "";

switch (fruit) {
  case "banana":
    result = "banana selected";
    break;
  case "apple":
    result = "apple selected";
    break;
  case "orange":
    result = "orange selected";
    break;
  default:
    result = "unknown fruit";
}
console.log("switch result:", result);

// TEST 2: switch with number
let num = 2;
let numResult = "";
switch (num) {
  case 1:
    numResult = "one";
    break;
  case 2:
    numResult = "two";
    break;
  case 3:
    numResult = "three";
    break;
  default:
    numResult = "other";
}
console.log("switch number result:", numResult);

// TEST 3: switch with default
let x = 99;
let defaultResult = "";
switch (x) {
  case 1:
    defaultResult = "one";
    break;
  default:
    defaultResult = "default triggered";
}
console.log("switch default result:", defaultResult);

// TEST 4: do...while loop
console.log("\n=== Testing do...while ===");
let count = 0;
do {
  count = count + 1;
  console.log("do-while iteration:", count);
} while (count < 3);
console.log("final count:", count);

// TEST 5: do...while runs at least once
let never = false;
let runOnce = 0;
do {
  runOnce = runOnce + 1;
} while (never);
console.log("do-while ran at least once:", runOnce);

// TEST 6: Symbol
console.log("\n=== Testing Symbol ===");
let sym1 = Symbol("description1");
console.log("Symbol created:", sym1);

let sym2 = Symbol("description2");
console.log("Symbols are unique:", sym1 !== sym2);

// Symbol.for and Symbol.keyFor
let globalSym = Symbol.for("global.key");
console.log("Symbol.for:", globalSym);

let sameGlobalSym = Symbol.for("global.key");
console.log("Same symbol from registry:", globalSym === sameGlobalSym);

let key = Symbol.keyFor(globalSym);
console.log("Symbol.keyFor:", key);

// TEST 7: BigInt
console.log("\n=== Testing BigInt ===");
let big1 = BigInt(123);
console.log("BigInt from number:", big1);

let big2 = BigInt("456");
console.log("BigInt from string:", big2);

let big3 = BigInt("789n");
console.log("BigInt from string with n:", big3);

console.log("\n=== All tests completed! ===");
