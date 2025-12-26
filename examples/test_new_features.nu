// ============================================================
// Nulang Test - New Features
// Testing: String Interpolation, Blob, File, Implements
// ============================================================

console.log("=== Testing Nulang New Features ===\n");

// ============================================================
// 1. String Interpolation (Template Literals)
// ============================================================
console.log("1. String Interpolation (Template Literals)");
console.log("-------------------------------------------");

let name = "Nulang";
let version = 1.0;
let greeting = `Hello, ${name}!`;
console.log(greeting);

let complex = `Version: ${version}, Name: ${name.toUpperCase()}`;
console.log(complex);

let math = `2 + 2 = ${2 + 2}`;
console.log(math);

let nested = `Items: ${[1, 2, 3].join(", ")}`;
console.log(nested);

let multiline = `This is
a multiline
string with ${name}`;
console.log(multiline);

// Expression inside template
let obj = { a: 1, b: 2 };
let objStr = `Object sum: ${obj.a + obj.b}`;
console.log(objStr);

console.log("✅ String interpolation working!\n");

// ============================================================
// 2. Blob Class
// ============================================================
console.log("2. Blob Class");
console.log("-------------------------------------------");

// Create a simple Blob
let blob1 = new Blob(["Hello, World!"]);
console.log("Blob created:", blob1);
console.log("Blob size:", blob1.size);
console.log("Blob type:", blob1.type);

// Create Blob with content type
let blob2 = new Blob(["<html><body>Test</body></html>"], { type: "text/html" });
console.log("HTML Blob size:", blob2.size);
console.log("HTML Blob type:", blob2.type);

// Create Blob from multiple parts
let blob3 = new Blob(["Part 1", " - ", "Part 2"]);
console.log("Multi-part Blob size:", blob3.size);

// Blob.slice()
let sliced = blob1.slice(0, 5);
console.log("Sliced Blob size:", sliced.size);

console.log("✅ Blob class working!\n");

// ============================================================
// 3. File Class
// ============================================================
console.log("3. File Class");
console.log("-------------------------------------------");

// Create a File
let file1 = new File(["File content here"], "test.txt", { type: "text/plain" });
console.log("File created:", file1);
console.log("File name:", file1.name);
console.log("File size:", file1.size);
console.log("File type:", file1.type);
console.log("File lastModified:", file1.lastModified);
console.log("File webkitRelativePath:", file1.webkitRelativePath);

// Create File with custom lastModified
let customDate = Date.now() - 86400000; // 1 day ago
let file2 = new File(["Old content"], "old.txt", {
  type: "text/plain",
  lastModified: customDate,
});
console.log("File2 lastModified:", file2.lastModified);

console.log("✅ File class working!\n");

// ============================================================
// 4. Implements in Classes
// ============================================================
console.log("4. Implements in Classes");
console.log("-------------------------------------------");

// Define an interface (TypeScript-like, no runtime validation)
interface Printable {
  print(): string;
}

interface Loggable {
  log(): void;
}

// Class implementing interfaces
class Document implements Printable, Loggable {
  constructor(content) {
    this.content = content;
  }

  print() {
    return `Document: ${this.content}`;
  }

  log() {
    console.log(this.print());
  }
}

let doc = new Document("Hello from Document!");
doc.log();
console.log("Print result:", doc.print());

// Class with extends and implements
interface Serializable {
  serialize(): string;
}

class Animal {
  constructor(name) {
    this.name = name;
  }

  speak() {
    return `${this.name} makes a sound`;
  }
}

class Dog extends Animal implements Serializable {
  constructor(name, breed) {
    super(name);
    this.breed = breed;
  }

  speak() {
    return `${this.name} barks`;
  }

  serialize() {
    return JSON.stringify({ name: this.name, breed: this.breed });
  }
}

let dog = new Dog("Rex", "German Shepherd");
console.log("Dog speaks:", dog.speak());
console.log("Dog serialized:", dog.serialize());

console.log("✅ Implements in classes working!\n");

// ============================================================
// 5. Combined Features Test
// ============================================================
console.log("5. Combined Features Test");
console.log("-------------------------------------------");

// Use template literals with Blob and File
let jsonData = { name: "Test", value: 42 };
let jsonContent = JSON.stringify(jsonData);
let jsonBlob = new Blob([jsonContent], { type: "application/json" });
console.log(`Created JSON Blob with ${jsonBlob.size} bytes`);

let fileName = "data.json";
let jsonFile = new File([jsonContent], fileName, { type: "application/json" });
console.log(`Created ${jsonFile.name} with ${jsonFile.size} bytes`);

console.log("\n=== All Tests Passed! ===");
