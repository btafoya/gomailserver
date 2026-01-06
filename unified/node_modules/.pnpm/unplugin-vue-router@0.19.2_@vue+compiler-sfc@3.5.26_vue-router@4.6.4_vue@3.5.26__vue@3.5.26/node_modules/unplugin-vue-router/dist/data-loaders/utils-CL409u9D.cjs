
//#region src/data-loaders/symbols.ts
/**
* Retrieves the internal version of loaders.
* @internal
*/
const LOADER_SET_KEY = Symbol("loaders");
/**
* Retrieves the internal version of loader entries.
* @internal
*/
const LOADER_ENTRIES_KEY = Symbol("loaderEntries");
/**
* Added to the loaders returned by `defineLoader()` to identify them.
* Allows to extract exported useData() from a component.
* @internal
*/
const IS_USE_DATA_LOADER_KEY = Symbol();
/**
* Symbol used to save the pending location on the router.
* @internal
*/
const PENDING_LOCATION_KEY = Symbol();
/**
* Symbol used to know there is no value staged for the loader and that commit should be skipped.
* @internal
*/
const STAGED_NO_VALUE = Symbol();
/**
* Gives access to the current app and it's `runWithContext` method.
* @internal
*/
const APP_KEY = Symbol();
/**
* Gives access to an AbortController that aborts when the navigation is canceled.
* @internal
*/
const ABORT_CONTROLLER_KEY = Symbol();
/**
* Gives access to the navigation results when the navigation is aborted by the user within a data loader.
* @internal
*/
const NAVIGATION_RESULTS_KEY = Symbol();
/**
* Symbol used to save the initial data on the router.
* @internal
*/
const IS_SSR_KEY = Symbol();
/**
* Symbol used to get the effect scope used for data loaders.
* @internal
*/
const DATA_LOADERS_EFFECT_SCOPE_KEY = Symbol();

//#endregion
//#region src/data-loaders/utils.ts
/**
* Check if a value is a `DataLoader`.
*
* @param loader - the object to check
*/
function isDataLoader(loader) {
	return loader && loader[IS_USE_DATA_LOADER_KEY];
}
/**
* @internal: data loaders authoring only. Use `getCurrentContext` instead.
*/
let currentContext;
function getCurrentContext() {
	return currentContext || [];
}
/**
* Sets the current context for data loaders. This allows for nested loaders to be aware of their parent context.
* INTERNAL ONLY.
*
* @param context - the context to set
* @internal
*/
function setCurrentContext(context) {
	currentContext = context ? context.length ? context : null : null;
}
/**
* Restore the current context after a promise is resolved.
* @param promise - promise to wrap
*/
function withLoaderContext(promise) {
	const context = currentContext;
	return promise.finally(() => currentContext = context);
}
const assign = Object.assign;
/**
* Track the reads of a route and its properties
* @internal
* @param route - route to track
*/
function trackRoute(route) {
	const [params, paramReads] = trackObjectReads(route.params);
	const [query, queryReads] = trackObjectReads(route.query);
	let hash = { v: null };
	return [
		{
			...route,
			get hash() {
				return hash.v = route.hash;
			},
			params,
			query
		},
		paramReads,
		queryReads,
		hash
	];
}
/**
*  Track the reads of an object (that doesn't change) and add the read properties to an object
* @internal
* @param obj - object to track
*/
function trackObjectReads(obj) {
	const reads = {};
	return [new Proxy(obj, { get(target, p, receiver) {
		const value = Reflect.get(target, p, receiver);
		reads[p] = value;
		return value;
	} }), reads];
}
/**
* Returns `true` if `inner` is a subset of `outer`. Used to check if a tr
*
* @internal
* @param outer - the bigger params
* @param inner - the smaller params
*/
function isSubsetOf(inner, outer) {
	for (const key in inner) {
		const innerValue = inner[key];
		const outerValue = outer[key];
		if (typeof innerValue === "string") {
			if (innerValue !== outerValue) return false;
		} else if (!innerValue || !outerValue) {
			if (innerValue !== outerValue) return false;
		} else if (!Array.isArray(outerValue) || outerValue.length !== innerValue.length || innerValue.some((value, i) => value !== outerValue[i])) return false;
	}
	return true;
}

//#endregion
Object.defineProperty(exports, 'ABORT_CONTROLLER_KEY', {
  enumerable: true,
  get: function () {
    return ABORT_CONTROLLER_KEY;
  }
});
Object.defineProperty(exports, 'APP_KEY', {
  enumerable: true,
  get: function () {
    return APP_KEY;
  }
});
Object.defineProperty(exports, 'DATA_LOADERS_EFFECT_SCOPE_KEY', {
  enumerable: true,
  get: function () {
    return DATA_LOADERS_EFFECT_SCOPE_KEY;
  }
});
Object.defineProperty(exports, 'IS_SSR_KEY', {
  enumerable: true,
  get: function () {
    return IS_SSR_KEY;
  }
});
Object.defineProperty(exports, 'IS_USE_DATA_LOADER_KEY', {
  enumerable: true,
  get: function () {
    return IS_USE_DATA_LOADER_KEY;
  }
});
Object.defineProperty(exports, 'LOADER_ENTRIES_KEY', {
  enumerable: true,
  get: function () {
    return LOADER_ENTRIES_KEY;
  }
});
Object.defineProperty(exports, 'LOADER_SET_KEY', {
  enumerable: true,
  get: function () {
    return LOADER_SET_KEY;
  }
});
Object.defineProperty(exports, 'NAVIGATION_RESULTS_KEY', {
  enumerable: true,
  get: function () {
    return NAVIGATION_RESULTS_KEY;
  }
});
Object.defineProperty(exports, 'PENDING_LOCATION_KEY', {
  enumerable: true,
  get: function () {
    return PENDING_LOCATION_KEY;
  }
});
Object.defineProperty(exports, 'STAGED_NO_VALUE', {
  enumerable: true,
  get: function () {
    return STAGED_NO_VALUE;
  }
});
Object.defineProperty(exports, 'assign', {
  enumerable: true,
  get: function () {
    return assign;
  }
});
Object.defineProperty(exports, 'currentContext', {
  enumerable: true,
  get: function () {
    return currentContext;
  }
});
Object.defineProperty(exports, 'getCurrentContext', {
  enumerable: true,
  get: function () {
    return getCurrentContext;
  }
});
Object.defineProperty(exports, 'isDataLoader', {
  enumerable: true,
  get: function () {
    return isDataLoader;
  }
});
Object.defineProperty(exports, 'isSubsetOf', {
  enumerable: true,
  get: function () {
    return isSubsetOf;
  }
});
Object.defineProperty(exports, 'setCurrentContext', {
  enumerable: true,
  get: function () {
    return setCurrentContext;
  }
});
Object.defineProperty(exports, 'trackRoute', {
  enumerable: true,
  get: function () {
    return trackRoute;
  }
});
Object.defineProperty(exports, 'withLoaderContext', {
  enumerable: true,
  get: function () {
    return withLoaderContext;
  }
});