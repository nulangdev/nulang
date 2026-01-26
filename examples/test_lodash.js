// Comprehensive Lodash Test Suite for Nulang
import _ from "lodash";

console.log("╔══════════════════════════════════════════════════════════════╗");
console.log("║           LODASH COMPREHENSIVE TEST SUITE                    ║");
console.log(
  "║                   Version:",
  _.VERSION,
  "                         ║",
);
console.log(
  "╚══════════════════════════════════════════════════════════════╝\n",
);

var passed = 0;
var failed = 0;

function test(name, fn) {
  try {
    var result = fn();
    if (result === true) {
      console.log("  ✓ " + name);
      passed++;
    } else {
      console.log("  ✗ " + name + " (returned: " + result + ")");
      failed++;
    }
  } catch (e) {
    console.log("  ✗ " + name + " (error: " + e.message + ")");
    failed++;
  }
}

// ═══════════════════════════════════════════════════════════════
console.log("📦 ARRAY METHODS");
console.log("─────────────────────────────────────────────────────────────");

test("_.head returns first element", function () {
  return _.head([1, 2, 3]) === 1;
});

test("_.last returns last element", function () {
  return _.last([1, 2, 3]) === 3;
});

test("_.first is alias for head", function () {
  return _.first([1, 2, 3]) === 1;
});

test("_.size returns array length", function () {
  return _.size([1, 2, 3, 4, 5]) === 5;
});

test("_.compact removes falsy values", function () {
  var result = _.compact([0, 1, false, 2, "", 3]);
  return (
    result.length === 3 && result[0] === 1 && result[1] === 2 && result[2] === 3
  );
});

test("_.flatten flattens array one level", function () {
  var result = _.flatten([
    [1, 2],
    [3, 4],
  ]);
  return result.length === 4 && result[0] === 1 && result[3] === 4;
});

test("_.uniq removes duplicates", function () {
  var result = _.uniq([1, 2, 1, 3, 2]);
  return result.length === 3;
});

test("_.indexOf finds index", function () {
  return _.indexOf([1, 2, 3, 4], 3) === 2;
});

test("_.take returns first n elements", function () {
  var result = _.take([1, 2, 3, 4, 5], 3);
  return result.length === 3 && result[2] === 3;
});

test("_.drop removes first n elements", function () {
  var result = _.drop([1, 2, 3, 4, 5], 2);
  return result.length === 3 && result[0] === 3;
});

test("_.initial returns all but last element", function () {
  var result = _.initial([1, 2, 3, 4]);
  return result.length === 3 && result[2] === 3;
});

test("_.tail returns all but first element", function () {
  var result = _.tail([1, 2, 3, 4]);
  return result.length === 3 && result[0] === 2;
});

test("_.nth returns element at index", function () {
  return _.nth([1, 2, 3, 4], 2) === 3;
});

test("_.chunk splits array into chunks", function () {
  var result = _.chunk([1, 2, 3, 4, 5], 2);
  return (
    result.length === 3 && result[0].length === 2 && result[2].length === 1
  );
});

test("_.fill fills array with value", function () {
  var result = _.fill([1, 2, 3], "a");
  return result[0] === "a" && result[1] === "a" && result[2] === "a";
});

// ═══════════════════════════════════════════════════════════════
console.log("\n📊 COLLECTION METHODS");
console.log("─────────────────────────────────────────────────────────────");

test("_.forEach iterates collection", function () {
  var sum = 0;
  _.forEach([1, 2, 3], function (n) {
    sum += n;
  });
  return sum === 6;
});

test("_.map transforms elements", function () {
  var result = _.map([1, 2, 3], function (n) {
    return n * 2;
  });
  return result[0] === 2 && result[1] === 4 && result[2] === 6;
});

test("_.filter selects matching elements", function () {
  var result = _.filter([1, 2, 3, 4, 5], function (n) {
    return n > 2;
  });
  return result.length === 3 && result[0] === 3;
});

test("_.find returns first match", function () {
  var result = _.find([1, 2, 3, 4], function (n) {
    return n > 2;
  });
  return result === 3;
});

test("_.reduce accumulates values", function () {
  var result = _.reduce(
    [1, 2, 3, 4],
    function (sum, n) {
      return sum + n;
    },
    0,
  );
  return result === 10;
});

test("_.every checks all elements", function () {
  return (
    _.every([2, 4, 6], function (n) {
      return n % 2 === 0;
    }) === true
  );
});

test("_.some checks any element", function () {
  return (
    _.some([1, 2, 3], function (n) {
      return n > 2;
    }) === true
  );
});

test("_.includes checks membership", function () {
  return _.includes([1, 2, 3], 2) === true;
});

test("_.sortBy sorts collection", function () {
  var result = _.sortBy([3, 1, 2]);
  return result[0] === 1 && result[1] === 2 && result[2] === 3;
});

test("_.groupBy groups by criterion", function () {
  var result = _.groupBy([1.3, 2.1, 2.4], Math.floor);
  return result["1"].length === 1 && result["2"].length === 2;
});

test("_.countBy counts by criterion", function () {
  var result = _.countBy([1, 2, 3, 4, 5], function (n) {
    return n % 2 === 0 ? "even" : "odd";
  });
  return result["odd"] === 3 && result["even"] === 2;
});

// ═══════════════════════════════════════════════════════════════
console.log("\n🔧 OBJECT METHODS");
console.log("─────────────────────────────────────────────────────────────");

test("_.keys returns object keys", function () {
  var result = _.keys({ a: 1, b: 2, c: 3 });
  return result.length === 3;
});

test("_.values returns object values", function () {
  var result = _.values({ a: 1, b: 2, c: 3 });
  return result.length === 3;
});

test("_.entries returns key-value pairs", function () {
  var result = _.entries({ a: 1, b: 2 });
  return result.length === 2 && result[0].length === 2;
});

test("_.assign merges objects", function () {
  var result = _.assign({ a: 1 }, { b: 2 }, { c: 3 });
  return result.a === 1 && result.b === 2 && result.c === 3;
});

test("_.merge deep merges objects", function () {
  var result = _.merge({ a: { x: 1 } }, { a: { y: 2 } });
  return result.a.x === 1 && result.a.y === 2;
});

test("_.pick selects properties", function () {
  var result = _.pick({ a: 1, b: 2, c: 3 }, ["a", "c"]);
  return result.a === 1 && result.c === 3 && result.b === undefined;
});

test("_.omit excludes properties", function () {
  var result = _.omit({ a: 1, b: 2, c: 3 }, ["b"]);
  return result.a === 1 && result.c === 3 && result.b === undefined;
});

test("_.get retrieves nested value", function () {
  var obj = { a: { b: { c: 3 } } };
  return _.get(obj, "a.b.c") === 3;
});

test("_.get with default value", function () {
  var obj = { a: 1 };
  return _.get(obj, "b.c", "default") === "default";
});

test("_.has checks property existence", function () {
  return _.has({ a: { b: 2 } }, "a") === true;
});

test("_.invert swaps keys and values", function () {
  var result = _.invert({ a: "1", b: "2" });
  return result["1"] === "a" && result["2"] === "b";
});

test("_.mapKeys transforms keys", function () {
  var result = _.mapKeys({ a: 1, b: 2 }, function (v, k) {
    return k + k;
  });
  return result.aa === 1 && result.bb === 2;
});

test("_.mapValues transforms values", function () {
  var result = _.mapValues({ a: 1, b: 2 }, function (v) {
    return v * 2;
  });
  return result.a === 2 && result.b === 4;
});

// ═══════════════════════════════════════════════════════════════
console.log("\n🔍 TYPE CHECKS");
console.log("─────────────────────────────────────────────────────────────");

test("_.isArray detects arrays", function () {
  return _.isArray([1, 2, 3]) === true && _.isArray("abc") === false;
});

test("_.isObject detects objects", function () {
  return _.isObject({}) === true && _.isObject([]) === true;
});

test("_.isString detects strings", function () {
  return _.isString("abc") === true && _.isString(123) === false;
});

test("_.isNumber detects numbers", function () {
  return _.isNumber(42) === true && _.isNumber("42") === false;
});

test("_.isBoolean detects booleans", function () {
  return _.isBoolean(true) === true && _.isBoolean(1) === false;
});

test("_.isFunction detects functions", function () {
  return _.isFunction(function () {}) === true;
});

test("_.isNil detects null/undefined", function () {
  return (
    _.isNil(null) === true &&
    _.isNil(undefined) === true &&
    _.isNil(0) === false
  );
});

test("_.isNull detects null", function () {
  return _.isNull(null) === true && _.isNull(undefined) === false;
});

test("_.isUndefined detects undefined", function () {
  return _.isUndefined(undefined) === true && _.isUndefined(null) === false;
});

test("_.isEmpty checks emptiness", function () {
  return (
    _.isEmpty([]) === true && _.isEmpty({}) === true && _.isEmpty([1]) === false
  );
});

test("_.isEqual performs deep equality", function () {
  return (
    _.isEqual({ a: 1 }, { a: 1 }) === true &&
    _.isEqual({ a: 1 }, { a: 2 }) === false
  );
});

test("_.isPlainObject detects plain objects", function () {
  return _.isPlainObject({}) === true && _.isPlainObject([]) === false;
});

// ═══════════════════════════════════════════════════════════════
console.log("\n📝 STRING METHODS");
console.log("─────────────────────────────────────────────────────────────");

test("_.camelCase converts to camelCase", function () {
  return _.camelCase("hello world") === "helloWorld";
});

test("_.capitalize capitalizes string", function () {
  return _.capitalize("hello") === "Hello";
});

test("_.kebabCase converts to kebab-case", function () {
  return _.kebabCase("hello world") === "hello-world";
});

test("_.snakeCase converts to snake_case", function () {
  return _.snakeCase("hello world") === "hello_world";
});

test("_.upperCase converts to UPPER CASE", function () {
  return _.upperCase("hello-world") === "HELLO WORLD";
});

test("_.lowerCase converts to lower case", function () {
  return _.lowerCase("HELLO WORLD") === "hello world";
});

test("_.trim removes whitespace", function () {
  return _.trim("  hello  ") === "hello";
});

test("_.pad pads string", function () {
  return _.pad("abc", 7) === "  abc  ";
});

test("_.repeat repeats string", function () {
  return _.repeat("ab", 3) === "ababab";
});

test("_.split splits string", function () {
  var result = _.split("a-b-c", "-");
  return result.length === 3 && result[0] === "a";
});

test("_.startsWith checks prefix", function () {
  return _.startsWith("hello world", "hello") === true;
});

test("_.endsWith checks suffix", function () {
  return _.endsWith("hello world", "world") === true;
});

// ═══════════════════════════════════════════════════════════════
console.log("\n🔢 MATH METHODS");
console.log("─────────────────────────────────────────────────────────────");

test("_.max finds maximum", function () {
  return _.max([1, 5, 3, 9, 2]) === 9;
});

test("_.min finds minimum", function () {
  return _.min([1, 5, 3, 9, 2]) === 1;
});

test("_.sum calculates sum", function () {
  return _.sum([1, 2, 3, 4]) === 10;
});

test("_.mean calculates average", function () {
  return _.mean([1, 2, 3, 4, 5]) === 3;
});

test("_.ceil rounds up", function () {
  return _.ceil(4.2) === 5;
});

test("_.floor rounds down", function () {
  return _.floor(4.8) === 4;
});

test("_.round rounds to nearest", function () {
  return _.round(4.5) === 5;
});

test("_.clamp clamps value", function () {
  return _.clamp(10, 0, 5) === 5 && _.clamp(-5, 0, 10) === 0;
});

test("_.inRange checks range", function () {
  return _.inRange(3, 2, 5) === true && _.inRange(5, 2, 5) === false;
});

// ═══════════════════════════════════════════════════════════════
console.log("\n🛠️ UTILITY METHODS");
console.log("─────────────────────────────────────────────────────────────");

test("_.identity returns input", function () {
  return _.identity(5) === 5;
});

test("_.constant returns constant function", function () {
  var fn = _.constant(42);
  return fn() === 42;
});

test("_.noop does nothing", function () {
  return _.noop() === undefined;
});

test("_.times calls n times", function () {
  var count = 0;
  _.times(5, function () {
    count++;
  });
  return count === 5;
});

test("_.range creates range array", function () {
  var result = _.range(5);
  return result.length === 5 && result[0] === 0 && result[4] === 4;
});

test("_.range with start and end", function () {
  var result = _.range(1, 5);
  return result.length === 4 && result[0] === 1 && result[3] === 4;
});

test("_.uniqueId generates unique ids", function () {
  var id1 = _.uniqueId("prefix_");
  var id2 = _.uniqueId("prefix_");
  return id1 !== id2;
});

test("_.clone shallow clones", function () {
  var obj = { a: 1, b: { c: 2 } };
  var clone = _.clone(obj);
  return clone.a === 1 && clone !== obj;
});

test("_.cloneDeep deep clones", function () {
  var obj = { a: { b: { c: 3 } } };
  var clone = _.cloneDeep(obj);
  clone.a.b.c = 999;
  return obj.a.b.c === 3;
});

test("_.defaultsDeep provides defaults", function () {
  var result = _.defaultsDeep({ a: { b: 1 } }, { a: { b: 2, c: 3 } });
  return result.a.b === 1 && result.a.c === 3;
});

// ═══════════════════════════════════════════════════════════════
console.log("\n⚡ FUNCTION METHODS");
console.log("─────────────────────────────────────────────────────────────");

test("_.once creates one-time function", function () {
  var count = 0;
  var fn = _.once(function () {
    count++;
    return count;
  });
  fn();
  fn();
  fn();
  return count === 1;
});

test("_.memoize caches results", function () {
  var count = 0;
  var fn = _.memoize(function (n) {
    count++;
    return n * 2;
  });
  fn(5);
  fn(5);
  fn(5);
  return count === 1 && fn(5) === 10;
});

test("_.negate inverts predicate", function () {
  var isEven = function (n) {
    return n % 2 === 0;
  };
  var isOdd = _.negate(isEven);
  return isOdd(3) === true && isOdd(4) === false;
});

test("_.partial applies partial args", function () {
  var add = function (a, b) {
    return a + b;
  };
  var add5 = _.partial(add, 5);
  return add5(3) === 8;
});

// ═══════════════════════════════════════════════════════════════
// SUMMARY
// ═══════════════════════════════════════════════════════════════
console.log(
  "\n╔══════════════════════════════════════════════════════════════╗",
);
console.log("║                      TEST SUMMARY                            ║");
console.log("╠══════════════════════════════════════════════════════════════╣");
console.log(
  "║  Passed: " +
    passed +
    "                                                   ║",
);
console.log(
  "║  Failed: " +
    failed +
    "                                                    ║",
);
console.log(
  "║  Total:  " +
    (passed + failed) +
    "                                                   ║",
);
console.log("╚══════════════════════════════════════════════════════════════╝");

if (failed === 0) {
  console.log("\n🎉 ALL TESTS PASSED! Lodash is fully compatible with Nulang!");
} else {
  console.log("\n⚠️  Some tests failed. See above for details.");
}
