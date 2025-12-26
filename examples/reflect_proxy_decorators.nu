// ==================================================
// Test file for Reflect, Proxy and Decorators
// ==================================================

console.log("=== Testing Reflect API ===\n");

// Test Reflect.get and Reflect.set
const obj = { name: "John", age: 30 };
console.log("Original object:", obj);
console.log("Reflect.get(obj, 'name'):", Reflect.get(obj, "name"));

Reflect.set(obj, "age", 31);
console.log("After Reflect.set(obj, 'age', 31):", obj);

// Test Reflect.has
console.log("Reflect.has(obj, 'name'):", Reflect.has(obj, "name"));
console.log("Reflect.has(obj, 'email'):", Reflect.has(obj, "email"));

// Test Reflect.ownKeys
console.log("Reflect.ownKeys(obj):", Reflect.ownKeys(obj));

// Test Reflect.deleteProperty
Reflect.deleteProperty(obj, "age");
console.log("After Reflect.deleteProperty(obj, 'age'):", obj);

// Test Reflect.getPrototypeOf and setPrototypeOf
const proto = {
  greet: function () {
    return "Hello!";
  },
};
Reflect.setPrototypeOf(obj, proto);
console.log(
  "After setting prototype, Reflect.getPrototypeOf(obj):",
  Reflect.getPrototypeOf(obj)
);

// Test Reflect.apply
function greet(greeting, punctuation) {
  return `${greeting}, ${this.name}${punctuation}`;
}
const result = Reflect.apply(greet, obj, ["Hello", "!"]);
console.log("Reflect.apply result:", result);

// Test Reflect.construct
class Person {
  constructor(name, age) {
    this.name = name;
    this.age = age;
  }

  sayHello() {
    return `Hi, I'm ${this.name}`;
  }
}

const person = Reflect.construct(Person, ["Alice", 25]);
console.log("Reflect.construct result:", person);
console.log("person.sayHello():", person.sayHello());

// Test Reflect.getOwnPropertyDescriptor
console.log(
  "Reflect.getOwnPropertyDescriptor(obj, 'name'):",
  Reflect.getOwnPropertyDescriptor(obj, "name")
);

console.log("\n=== Testing Proxy ===\n");

// Basic Proxy with get trap
const target = { message: "Hello World", count: 0 };
const handler = {
  get: function (target, property, receiver) {
    console.log(`[Proxy] Getting property: ${property}`);
    return Reflect.get(target, property);
  },
  set: function (target, property, value, receiver) {
    console.log(`[Proxy] Setting property: ${property} = ${value}`);
    return Reflect.set(target, property, value);
  },
};

const proxy = new Proxy(target, handler);
console.log("proxy.message:", proxy.message);
proxy.count = 42;
console.log("proxy.count:", proxy.count);

// Proxy with validation
const validationHandler = {
  set: function (target, property, value, receiver) {
    if (property === "age" && typeof value === "number" && value < 0) {
      console.log("Error: Age cannot be negative!");
      return false;
    }
    return Reflect.set(target, property, value);
  },
};

const validatedPerson = new Proxy({ name: "Bob", age: 25 }, validationHandler);
console.log("\nValidated Proxy:");
console.log("Setting valid age...");
validatedPerson.age = 30;
console.log("Age is now:", validatedPerson.age);

console.log("Trying to set negative age...");
validatedPerson.age = -5;
console.log("Age is still:", validatedPerson.age);

// Proxy.revocable
console.log("\n=== Testing Proxy.revocable ===");
const revocableResult = Proxy.revocable(
  { data: "secret" },
  {
    get: function (target, prop) {
      console.log(`[Revocable] Accessing ${prop}`);
      return target[prop];
    },
  }
);

const revocableProxy = revocableResult.proxy;
const revoke = revocableResult.revoke;

console.log("Before revoke - data:", revocableProxy.data);
revoke();
console.log("Proxy revoked!");

console.log("\n=== Testing Decorators ===\n");

// Define decorator functions
function log(target) {
  console.log(`[Decorator] Class ${target.Name} was decorated!`);
  return target;
}

function withTimestamp(target) {
  console.log("[Decorator] Adding timestamp functionality...");
  // Note: In real implementation, this would add methods/properties
  return target;
}

// Decorator factory
function Metadata(key, value) {
  return function (target) {
    console.log(`[Decorator Factory] Setting metadata: ${key} = ${value}`);
    return target;
  };
}

// Using decorators on classes
@log
@withTimestamp
@Metadata("version", "1.0")
class MyComponent {
  constructor(name) {
    this.name = name;
  }

  render() {
    return `<div>${this.name}</div>`;
  }
}

const component = new MyComponent("TestComponent");
console.log("Component name:", component.name);
console.log("Component render:", component.render());

console.log("\n=== All tests completed! ===");
