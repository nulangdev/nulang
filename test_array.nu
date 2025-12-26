// Test Array Constructor
console.log("--- Array Constructor ---");
const arr1 = new Array(3);
console.log("new Array(3):", arr1, "length:", arr1.length);
if (arr1.length !== 3) {
  console.error("FAIL: new Array(3) length mismatch");
}
if (arr1[0] !== undefined) {
  console.error("FAIL: new Array(3) elements undefined");
}

const arr2 = new Array(1, 2, 3);
console.log("new Array(1, 2, 3):", arr2);
if (arr2.length !== 3) {
  console.error("FAIL: new Array(1, 2, 3) length mismatch");
}
if (arr2[2] !== 3) {
  console.error("FAIL: new Array(1, 2, 3) content mismatch");
}

// Test Array.from
console.log("\n--- Array.from ---");
const fromStr = Array.from("foo");
console.log("Array.from('foo'):", fromStr);
if (fromStr[0] !== "f" || fromStr[1] !== "o") {
  console.error("FAIL: Array.from string");
}

const fromArr = Array.from([1, 2, 3], (x) => x * 2);
console.log("Array.from([1, 2, 3], x => x * 2):", fromArr);
if (fromArr[0] !== 2 || fromArr[2] !== 6) {
  console.error("FAIL: Array.from map");
}

const helper = { factor: 3 };
const fromMapThis = Array.from(
  [1, 2],
  function (x) {
    return x * this.factor;
  },
  helper
);
console.log("Array.from(..., thisArg):", fromMapThis);
if (fromMapThis[0] !== 3 || fromMapThis[1] !== 6) {
  console.error("FAIL: Array.from thisArg");
}

// Test map with thisArg
console.log("\n--- map with thisArg ---");
const mapper = { val: 10 };
const mapped = [1, 2].map(function (x) {
  return x + this.val;
}, mapper);
console.log("[1, 2].map(..., thisArg):", mapped);
if (mapped[0] !== 11) {
  console.error("FAIL: map thisArg");
}

// Test includes with fromIndex
console.log("\n--- includes with fromIndex ---");
const incArr = [1, 2, 3];
console.log("includes(2):", incArr.includes(2)); // true
console.log("includes(2, 2):", incArr.includes(2, 2)); // false, starts at index 2 (val 3)
console.log("includes(3, -1):", incArr.includes(3, -1)); // true
if (!incArr.includes(2)) {
  console.error("FAIL: includes basic");
}
if (incArr.includes(2, 2)) {
  console.error("FAIL: includes fromIndex");
}
if (!incArr.includes(3, -1)) {
  console.error("FAIL: includes negative fromIndex");
}

// Test splice
console.log("\n--- splice ---");
const sArr = [1, 2, 3, 4, 5];
const removed = sArr.splice(1, 2);
console.log("splice(1, 2) removed:", removed);
console.log("sArr after splice:", sArr);

if (removed.length !== 2 || removed[0] !== 2) {
  console.error("FAIL: splice removed");
}
if (sArr.length !== 3 || sArr[1] !== 4) {
  console.error("FAIL: splice remaining");
}

// Test other methods briefly
console.log("\n--- other methods ---");
const fArr = [1, 2, 3, 4];
const filtered = fArr.filter((x) => x > 2);
console.log("filter > 2:", filtered);

const found = fArr.find((x) => x > 3);
console.log("find > 3:", found);

const idx = fArr.findIndex((x) => x > 3);
console.log("findIndex > 3:", idx);

const sum = fArr.reduce((acc, x) => acc + x, 0);
console.log("reduce sum:", sum);

console.log("Tests Completed");
