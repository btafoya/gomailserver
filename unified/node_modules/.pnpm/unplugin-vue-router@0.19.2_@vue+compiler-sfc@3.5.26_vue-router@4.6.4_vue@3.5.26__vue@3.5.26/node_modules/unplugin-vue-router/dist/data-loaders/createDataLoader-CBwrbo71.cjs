
//#region src/data-loaders/createDataLoader.ts
const toLazyValue = (lazy, to, from) => typeof lazy === "function" ? lazy(to, from) : lazy;

//#endregion
Object.defineProperty(exports, 'toLazyValue', {
  enumerable: true,
  get: function () {
    return toLazyValue;
  }
});