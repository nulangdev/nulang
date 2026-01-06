// Proxy Test Suite based on doc/proxy.md

// 1. Basic get/set
console.log("--- Basic Get/Set ---");
const target1 = {
  message: "Hello, World!",
};

const handler1 = {
  get: function (target, property, receiver) {
    if (property == "message") {
      return "Intercepted: " + target[property];
    }
    return target[property];
  },
  set: function (target, property, value, receiver) {
    console.log("Setting " + property + " to " + value);
    target[property] = value;
    return true;
  },
};

const proxy1 = new Proxy(target1, handler1);
console.log(proxy1.message); // Intercepted: Hello, World!
proxy1.newProp = 123; // Setting newProp to 123
console.log(proxy1.newProp); // 123

// 2. Has trap
console.log("\n--- Has Trap ---");
const handler2 = {
  has: function (target, property) {
    if (property.startsWith("_")) {
      return false; // Esconde propriedades privadas
    }
    return property in target;
  },
};

const obj2 = { _private: 1, public: 2 };
const proxy2 = new Proxy(obj2, handler2);

console.log("public" in proxy2); // true
console.log("_private" in proxy2); // false

// 3. DeleteProperty trap
console.log("\n--- DeleteProperty Trap ---");
const handler3 = {
  deleteProperty: function (target, property) {
    if (property.startsWith("_")) {
      console.log("Não é possível deletar propriedades privadas");
      return false;
    }
    delete target[property];
    return true;
  },
};

const obj3 = { _id: 1, name: "test" };
const proxy3 = new Proxy(obj3, handler3);

delete proxy3.name; // OK
console.log(obj3.name); // undefined
delete proxy3._id; // "Não é possível deletar propriedades privadas"
console.log(obj3._id); // 1

// 4. Apply trap
console.log("\n--- Apply Trap ---");
const handler4 = {
  apply: function (target, thisArg, argumentsList) {
    console.log(`Chamando com args: ${argumentsList}`);
    return target.apply(thisArg, argumentsList);
  },
};

function sum(a, b) {
  return a + b;
}

const proxy4 = new Proxy(sum, handler4);
console.log(proxy4(1, 2));
// Output:
// Chamando com args: 1,2
// 3

// 5. Construct trap
console.log("\n--- Construct Trap ---");
const handler5 = {
  construct: function (target, argumentsList, newTarget) {
    console.log(`Construindo com args: ${argumentsList}`);
    return new target(...argumentsList);
  },
};

class Person {
  constructor(name) {
    this.name = name;
  }
}

const ProxyPerson = new Proxy(Person, handler5);
const p = new ProxyPerson("João");
// Output: Construindo com args: João
console.log(p.name); // "João"

// 6. getPrototypeOf trap
console.log("\n--- GetPrototypeOf Trap ---");
const handler6 = {
  getPrototypeOf: function (target) {
    return { custom: true };
  },
};

const obj6 = {};
const proxy6 = new Proxy(obj6, handler6);
const proto = Reflect.getPrototypeOf(proxy6); // or Object.getPrototypeOf if implemented
console.log(proto.custom); // true

// 7. setPrototypeOf trap
console.log("\n--- SetPrototypeOf Trap ---");
const handler7 = {
  setPrototypeOf: function (target, prototype) {
    console.log("Tentativa de alterar prototype bloqueada");
    return false;
  },
};

const obj7 = {};
const proxy7 = new Proxy(obj7, handler7);
Reflect.setPrototypeOf(proxy7, {}); // "Tentativa de alterar prototype bloqueada"

// 8. isExtensible trap
console.log("\n--- IsExtensible Trap ---");
const handler8 = {
  isExtensible: function (target) {
    return false; // Lie about it
  },
};

const obj8 = {};
const proxy8 = new Proxy(obj8, handler8);
console.log(Reflect.isExtensible(proxy8)); // true (Wait, we returned false in trap, should be false?)
// If the target IS extensible, the proxy MUST return true.
// My implementation doesn't enforce invariants yet, so it should return what the trap returns.
// Wait, my implementation returns what the trap returns directly?
// Evaluator: return applyProxyTrap(...) -> returns result of trap.
// If trap returns false, it returns false.

// 9. preventExtensions trap
console.log("\n--- PreventExtensions Trap ---");
const handler9 = {
  preventExtensions: function (target) {
    console.log("Preventing extensions");
    return true;
  },
};
const proxy9 = new Proxy({}, handler9);
Reflect.preventExtensions(proxy9);

// 10. getOwnPropertyDescriptor trap
console.log("\n--- GetOwnPropertyDescriptor Trap ---");
const handler10 = {
  getOwnPropertyDescriptor: function (target, property) {
    return {
      value: "virtual",
      writable: true,
      enumerable: true,
      configurable: true,
    };
  },
};
const proxy10 = new Proxy({}, handler10);
const desc = Reflect.getOwnPropertyDescriptor(proxy10, "foo");
console.log(desc.value); // virtual

// 11. defineProperty trap
console.log("\n--- DefineProperty Trap ---");
const handler11 = {
  defineProperty: function (target, property, descriptor) {
    console.log(`Definindo ${property}`);
    return true;
  },
};
const proxy11 = new Proxy({}, handler11);
Reflect.defineProperty(proxy11, "bar", { value: 1 });

// 12. ownKeys trap
console.log("\n--- OwnKeys Trap ---");
const handler12 = {
  ownKeys: function (target) {
    return ["visible"];
  },
};

const obj12 = { _hidden: 1, visible: 2 };
const proxy12 = new Proxy(obj12, handler12);
const keys = Reflect.ownKeys(proxy12);
console.log(keys); // ["visible"]
console.log(Object.keys(proxy12)); // ["visible"]

// 13. Revocable
console.log("\n--- Revocable ---");
const target13 = { name: "Test" };
const handler13 = {
  get: (t, p) => t[p],
};

const revocable = Proxy.revocable(target13, handler13);
const proxy13 = revocable.proxy;
const revoke = revocable.revoke;

console.log(proxy13.name); // "Test"

revoke();

try {
  console.log(proxy13.name);
} catch (e) {
  console.log(e.message); // "Cannot perform 'get' on a proxy that has been revoked"
}
