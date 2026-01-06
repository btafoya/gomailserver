//#region src/data-loaders/createDataLoader.ts
const toLazyValue = (lazy, to, from) => typeof lazy === "function" ? lazy(to, from) : lazy;

//#endregion
export { toLazyValue as t };