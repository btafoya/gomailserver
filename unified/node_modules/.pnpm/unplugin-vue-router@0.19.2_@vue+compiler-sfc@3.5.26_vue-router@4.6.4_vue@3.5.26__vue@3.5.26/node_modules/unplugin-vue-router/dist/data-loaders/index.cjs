const require_createDataLoader = require('./createDataLoader-CBwrbo71.cjs');
const require_utils = require('./utils-CL409u9D.cjs');
let vue_router = require("vue-router");
let vue = require("vue");

//#region src/data-loaders/navigation-guard.ts
/**
* Key to inject the global loading state for loaders used in `useIsDataLoading`.
* @internal
*/
const IS_DATA_LOADING_KEY = Symbol();
/**
* TODO: export functions that allow preloading outside of a navigation guard
*/
/**
* Setups the different Navigation Guards to collect the data loaders from the route records and then to execute them.
* @internal used by the `DataLoaderPlugin`
* @see {@link DataLoaderPlugin}
*
* @param router - the router instance
* @returns
*/
function setupLoaderGuard({ router, app, effect: scope, isSSR, errors: globalErrors = [], selectNavigationResult = (results) => results[0].value }) {
	if (router[require_utils.LOADER_ENTRIES_KEY] != null) {
		if (process.env.NODE_ENV !== "production") console.warn("[vue-router]: Data Loader was setup twice. Make sure to setup only once.");
		return () => {};
	}
	if (process.env.NODE_ENV === "development" && !isSSR) console.warn("[vue-router]: Data Loader is experimental and subject to breaking changes in the future.");
	router[require_utils.LOADER_ENTRIES_KEY] = /* @__PURE__ */ new WeakMap();
	router[require_utils.APP_KEY] = app;
	router[require_utils.DATA_LOADERS_EFFECT_SCOPE_KEY] = scope;
	router[require_utils.IS_SSR_KEY] = !!isSSR;
	const isDataLoading = scope.run(() => (0, vue.shallowRef)(false));
	app.provide(IS_DATA_LOADING_KEY, isDataLoading);
	const removeLoaderGuard = router.beforeEach((to) => {
		if (router[require_utils.PENDING_LOCATION_KEY]) router[require_utils.PENDING_LOCATION_KEY].meta[require_utils.ABORT_CONTROLLER_KEY]?.abort();
		router[require_utils.PENDING_LOCATION_KEY] = to;
		to.meta[require_utils.LOADER_SET_KEY] = /* @__PURE__ */ new Set();
		to.meta[require_utils.ABORT_CONTROLLER_KEY] = new AbortController();
		to.meta[require_utils.NAVIGATION_RESULTS_KEY] = [];
		for (const record of to.matched) record.meta[require_utils.LOADER_SET_KEY] ??= new Set(record.meta.loaders || []);
	});
	const removeDataLoaderGuard = router.beforeResolve((to, from) => {
		for (const record of to.matched) {
			for (const loader of record.meta[require_utils.LOADER_SET_KEY]) to.meta[require_utils.LOADER_SET_KEY].add(loader);
			for (const componentName in record.mods) {
				const viewModule = record.mods[componentName];
				for (const exportName in viewModule) {
					const exportValue = viewModule[exportName];
					if (require_utils.isDataLoader(exportValue)) to.meta[require_utils.LOADER_SET_KEY].add(exportValue);
				}
				const component = record.components?.[componentName];
				if (component && Array.isArray(component.__loaders)) {
					for (const loader of component.__loaders) if (require_utils.isDataLoader(loader)) to.meta[require_utils.LOADER_SET_KEY].add(loader);
				}
			}
		}
		const loaders = Array.from(to.meta[require_utils.LOADER_SET_KEY]);
		const { signal } = to.meta[require_utils.ABORT_CONTROLLER_KEY];
		require_utils.setCurrentContext([]);
		isDataLoading.value = true;
		return Promise.all(loaders.map((loader) => {
			const { server, lazy, errors } = loader._.options;
			if (!server && isSSR) return;
			const ret = scope.run(() => app.runWithContext(() => loader._.load(to, router, from)));
			return !isSSR && require_createDataLoader.toLazyValue(lazy, to, from) ? void 0 : ret.catch((reason) => {
				if (errors === true) {
					if (Array.isArray(globalErrors) ? globalErrors.some((Err) => reason instanceof Err) : globalErrors(reason)) return;
				} else if (errors && (Array.isArray(errors) ? errors.some((Err) => reason instanceof Err) : errors(reason))) return;
				throw reason;
			});
		})).then(() => {
			if (to.meta[require_utils.NAVIGATION_RESULTS_KEY].length) return selectNavigationResult(to.meta[require_utils.NAVIGATION_RESULTS_KEY]);
		}).catch((error) => error instanceof NavigationResult ? error.value : signal.aborted && error === signal.reason ? false : Promise.reject(error)).finally(() => {
			require_utils.setCurrentContext([]);
			isDataLoading.value = false;
		});
	});
	const removeAfterEach = router.afterEach((to, from, failure) => {
		if (failure) {
			to.meta[require_utils.ABORT_CONTROLLER_KEY]?.abort(failure);
			if ((0, vue_router.isNavigationFailure)(failure, 16)) for (const loader of to.meta[require_utils.LOADER_SET_KEY]) loader._.getEntry(router).resetPending();
		} else for (const loader of to.meta[require_utils.LOADER_SET_KEY]) {
			const { commit, lazy } = loader._.options;
			if (commit === "after-load") {
				const entry = loader._.getEntry(router);
				if (entry && (!require_createDataLoader.toLazyValue(lazy, to, from) || !entry.isLoading.value)) entry.commit(to);
			}
		}
		if (router[require_utils.PENDING_LOCATION_KEY] === to) router[require_utils.PENDING_LOCATION_KEY] = null;
	});
	const removeOnError = router.onError((error, to) => {
		to.meta[require_utils.ABORT_CONTROLLER_KEY]?.abort(error);
		if (router[require_utils.PENDING_LOCATION_KEY] === to) router[require_utils.PENDING_LOCATION_KEY] = null;
	});
	return () => {
		delete router[require_utils.LOADER_ENTRIES_KEY];
		delete router[require_utils.APP_KEY];
		removeLoaderGuard();
		removeDataLoaderGuard();
		removeAfterEach();
		removeOnError();
	};
}
/**
* Possible values to change the result of a navigation within a loader. Can be returned from a data loader and will
* appear in `selectNavigationResult`. If thrown, it will immediately cancel the navigation. It can only contain values
* that cancel the navigation.
*
* @example
* ```ts
* export const useUserData = defineLoader(async (to) => {
*   const user = await fetchUser(to.params.id)
*   if (!user) {
*     return new NavigationResult('/404')
*   }
*   return user
* })
* ```
*/
var NavigationResult = class {
	constructor(value) {
		this.value = value;
	}
};
/**
* Data Loader plugin to add data loading support to Vue Router.
*
* @example
* ```ts
* const router = createRouter({
*   routes,
*   history: createWebHistory(),
* })
*
* const app = createApp({})
* app.use(DataLoaderPlugin, { router })
* app.use(router)
* ```
*/
function DataLoaderPlugin(app, options) {
	const effect = (0, vue.effectScope)(true);
	const removeGuards = setupLoaderGuard(require_utils.assign({
		app,
		effect
	}, options));
	const { unmount } = app;
	app.unmount = () => {
		effect.stop();
		removeGuards();
		unmount.call(app);
	};
}
/**
* Return a ref that reflects the global loading state of all loaders within a navigation.
* This state doesn't update if `refresh()` is manually called.
*/
function useIsDataLoading() {
	return (0, vue.inject)(IS_DATA_LOADING_KEY);
}

//#endregion
exports.ABORT_CONTROLLER_KEY = require_utils.ABORT_CONTROLLER_KEY;
exports.APP_KEY = require_utils.APP_KEY;
exports.DATA_LOADERS_EFFECT_SCOPE_KEY = require_utils.DATA_LOADERS_EFFECT_SCOPE_KEY;
exports.DataLoaderPlugin = DataLoaderPlugin;
exports.IS_SSR_KEY = require_utils.IS_SSR_KEY;
exports.IS_USE_DATA_LOADER_KEY = require_utils.IS_USE_DATA_LOADER_KEY;
exports.LOADER_ENTRIES_KEY = require_utils.LOADER_ENTRIES_KEY;
exports.LOADER_SET_KEY = require_utils.LOADER_SET_KEY;
exports.NAVIGATION_RESULTS_KEY = require_utils.NAVIGATION_RESULTS_KEY;
exports.NavigationResult = NavigationResult;
exports.PENDING_LOCATION_KEY = require_utils.PENDING_LOCATION_KEY;
exports.STAGED_NO_VALUE = require_utils.STAGED_NO_VALUE;
exports.assign = require_utils.assign;
Object.defineProperty(exports, 'currentContext', {
  enumerable: true,
  get: function () {
    return require_utils.currentContext;
  }
});
exports.getCurrentContext = require_utils.getCurrentContext;
exports.isSubsetOf = require_utils.isSubsetOf;
exports.setCurrentContext = require_utils.setCurrentContext;
exports.toLazyValue = require_createDataLoader.toLazyValue;
exports.trackRoute = require_utils.trackRoute;
exports.useIsDataLoading = useIsDataLoading;
exports.withLoaderContext = require_utils.withLoaderContext;