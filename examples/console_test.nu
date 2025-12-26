// Teste completo do Console API

console.log("=== Teste do Console API ===");
console.log();

// log, info, debug
console.log("console.log:", "Hello", "World");
console.info("console.info:", "Information message");
console.debug("console.debug:", "Debug message");

// error, warn
console.error("console.error:", "This is an error");
console.warn("console.warn:", "This is a warning");

console.log();

// assert
console.log("=== console.assert ===");
console.assert(true, "This should NOT appear");
console.assert(false, "This assertion failed!");
console.assert(1 === 2, "Math is broken!");

console.log();

// count / countReset
console.log("=== console.count / countReset ===");
console.count("loop");
console.count("loop");
console.count("loop");
console.count(); // default
console.count(); // default
console.countReset("loop");
console.count("loop"); // should be 1

console.log();

// group / groupEnd
console.log("=== console.group ===");
console.log("Outside group");
console.group("Group 1");
console.log("Inside group 1");
console.group("Nested Group");
console.log("Inside nested group");
console.groupEnd();
console.log("Back to group 1");
console.groupEnd();
console.log("Outside group again");

console.log();

// table
console.log("=== console.table ===");
console.table(["apple", "banana", "cherry"]);
console.log();
console.table({ name: "John", age: 30, city: "NYC" });

console.log();

// time / timeEnd / timeLog
console.log("=== console.time / timeEnd / timeLog ===");
console.time("operation");

// Simulate some work
let sum = 0;
for (let i = 0; i < 10000; i++) {
  sum += i;
}

console.timeLog("operation", "sum is", sum);

// More work
for (let i = 0; i < 10000; i++) {
  sum += i;
}

console.timeEnd("operation");

console.log();

// trace
console.log("=== console.trace ===");
console.trace("Stack trace test");

console.log();

// dir / dirxml
console.log("=== console.dir / dirxml ===");
const obj = { a: 1, b: { c: 2, d: 3 } };
console.dir(obj);
console.dirxml(obj);

console.log();
console.log("=== Console API Test Complete ===");
