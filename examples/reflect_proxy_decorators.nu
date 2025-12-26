const handler = {
  has: function (target, property) {
    if (property.startsWith("_")) {
      return false; // Esconde propriedades privadas
    }
    return property in target;
  },
};

const obj = { _private: 1, public: 2 };
const proxy = new Proxy(obj, handler);

console.log("public" in proxy); // true
console.log("_private" in proxy); // false
