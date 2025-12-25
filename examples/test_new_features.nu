// Test ES6 Features: Classes, Timers, HTTP, Streams
// ==================================================

console.log("=== Testing Classes ===");

// Basic class
class Animal {
  constructor(name) {
    this.name = name;
  }

  speak() {
    console.log(this.name + " makes a sound");
  }
}

let animal = new Animal("Generic Animal");
animal.speak();
console.log("Animal name:", animal.name);

// Class inheritance (without super() call - set properties directly)
class Dog extends Animal {
  constructor(name, breed) {
    this.name = name;
    this.breed = breed;
  }

  speak() {
    console.log(this.name + " barks!");
  }

  getInfo() {
    return this.name + " is a " + this.breed;
  }
}

let dog = new Dog("Rex", "German Shepherd");
dog.speak();
console.log("Dog info:", dog.getInfo());

// Static methods
class Calculator {
  static add(a, b) {
    return a + b;
  }

  static multiply(a, b) {
    return a * b;
  }
}

console.log("Calculator.add(5, 3):", Calculator.add(5, 3));
console.log("Calculator.multiply(4, 7):", Calculator.multiply(4, 7));

console.log("\n=== Testing Timers ===");

// setTimeout example
console.log("Setting timeout for 100ms...");
let timeoutId = setTimeout(function () {
  console.log("Timeout executed!");
}, 100);

// setInterval example (will run 3 times then clear)
let count = 0;
let intervalId = setInterval(function () {
  count++;
  console.log("Interval count:", count);
  if (count >= 3) {
    clearInterval(intervalId);
    console.log("Interval cleared");
  }
}, 50);

// Sleep function
console.log("Sleeping for 200ms...");
sleep(200);
console.log("Woke up!");

console.log("\n=== Testing HTTP Module ===");

// HTTP module is available
console.log("HTTP module available:", typeof http);
console.log("http.get:", typeof http.get);
console.log("http.post:", typeof http.post);

// Fetch function is available globally
console.log("fetch available:", typeof fetch);

console.log("\n=== Testing Streams ===");

// Create a readable stream
let readable = stream.Readable();
console.log("Created readable stream:", typeof readable);
console.log("readable.push:", typeof readable.push);
console.log("readable.pipe:", typeof readable.pipe);

// Create a writable stream
let writable = stream.Writable();
console.log("Created writable stream:", typeof writable);
console.log("writable.write:", typeof writable.write);

// Create a transform stream
let transform = stream.Transform();
console.log("Created transform stream:", typeof transform);

// Test writing to writable stream
writable.write("Hello ");
writable.write("World!");
console.log("Writable contents:", writable.toString());

console.log("\n=== Testing URL Module ===");

let parsedUrl = url.parse(
  "https://example.com:8080/path/to/resource?query=value#hash"
);
console.log("Parsed URL:");
console.log("  protocol:", parsedUrl.protocol);
console.log("  host:", parsedUrl.host);
console.log("  pathname:", parsedUrl.pathname);
console.log("  search:", parsedUrl.search);
console.log("  hash:", parsedUrl.hash);

console.log("\n=== Testing QueryString Module ===");

let obj = { name: "John", age: "30", city: "NYC" };
let qs = querystring.stringify(obj);
console.log("Stringified:", qs);

let parsed = querystring.parse("foo=bar&baz=qux");
console.log("Parsed foo:", parsed.foo);
console.log("Parsed baz:", parsed.baz);

// Wait for async timers to complete
sleep(500);

console.log("\n=== All Tests Complete ===");
