// Test TypeScript-like syntax in Nulang

interface User {
  id: string;
  name: string;
  age?: number;
  sayHello(msg: string): void;
}

type ID = string | number;

declare global {
  let appVersion: string;
  function initialize(): void;
}

declare module "my-module" {
  export function doWork(data: any): boolean;
}

// Variable with type annotation
let username: string = "antigravity";
const pi: number = 3.14159;

// Function with type annotations
function greet(name: string, age: number): string {
  return "Hello, " + name + "! You are " + age + " years old.";
}

// Arrow function with type annotations
const add = (a: number, b: number): number => a + b;

class Person {
  name: string;

  constructor(name: string) {
    this.name = name;
  }

  say(msg: string): void {
    console.log(this.name + " says: " + msg);
  }
}

console.log("🚀 Testing TS-like syntax...");
console.log(greet(username, 25));
console.log("3 + 5 =", add(3, 5));

let p = new Person("Antigravity");
p.say("Nulang is growing!");

console.log("✅ All TS-like syntax parsed and ignored correctly at runtime.");
