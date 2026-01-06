import { t as toLazyValue } from "./createDataLoader-BK9Gdnky.mjs";
import { _ as PENDING_LOCATION_KEY, a as isSubsetOf, c as withLoaderContext, d as DATA_LOADERS_EFFECT_SCOPE_KEY, f as IS_SSR_KEY, g as NAVIGATION_RESULTS_KEY, h as LOADER_SET_KEY, i as isDataLoader, l as ABORT_CONTROLLER_KEY, m as LOADER_ENTRIES_KEY, n as currentContext, o as setCurrentContext, p as IS_USE_DATA_LOADER_KEY, r as getCurrentContext, s as trackRoute, t as assign, u as APP_KEY, v as STAGED_NO_VALUE } from "./utils-C5t8q_dT.mjs";
import { isNavigationFailure } from "vue-router";
import { effectScope, inject, shallowRef } from "vue";

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
	if (router[LOADER_ENTRIES_KEY] != null) {
		if (process.env.NODE_ENV !== "production") console.warn("[vue-router]: Data Loader was setup twice. Make sure to setup only once.");
		return () => {};
	}
	if (process.env.NODE_ENV === "development" && !isSSR) console.warn("[vue-router]: Data Loader is experimental and subject to breaking changes in the future.");
	router[LOADER_ENTRIES_KEY] = /* @__PURE__ */ new WeakMap();
	router[APP_KEY] = app;
	router[DATA_LOADERS_EFFECT_SCOPE_KEY] = scope;
	router[IS_SSR_KEY] = !!isSSR;
	const isDataLoading = scope.run(() => shallowRef(false));
	app.provide(IS_DATA_LOADING_KEY, isDataLoading);
	const removeLoaderGuard = router.beforeEach((to) => {
		if (router[PENDING_LOCATION_KEY]) router[PENDING_LOCATION_KEY].meta[ABORT_CONTROLLER_KEY]?.abort();
		router[PENDING_LOCATION_KEY] = to;
		to.meta[LOADER_SET_KEY] = /* @__PURE__ */ new Set();
		to.meta[ABORT_CONTROLLER_KEY] = new AbortController();
		to.meta[NAVIGATION_RESULTS_KEY] = [];
		for (const record of to.matched) record.meta[LOADER_SET_KEY] ??= new Set(record.meta.loaders || []);
	});
	const removeDataLoaderGuard = router.beforeResolve((to, from) => {
		for (const record of to.matched) {
			for (const loader of record.meta[LOADER_SET_KEY]) to.meta[LOADER_SET_KEY].add(loader);
			for (const componentName in record.mods) {
				const viewModule = record.mods[componentName];
				for (const exportName in viewModule) {
					const exportValue = viewModule[exportName];
					if (isDataLoader(exportValue)) to.meta[LOADER_SET_KEY].add(exportValue);
				}
				const component = record.components?.[componentName];
				if (component && Array.isArray(component.__loaders)) {
					for (const loader of component.__loaders) if (isDataLoader(loader)) to.meta[LOADER_SET_KEY].add(loader);
				}
			}
		}
		const loaders = Array.from(to.meta[LOADER_SET_KEY]);
		const { signal } = to.meta[ABORT_CONTROLLER_KEY];
		setCurrentContext([]);
		isDataLoading.value = true;
		return Promise.all(loaders.map((loader) => {
			const { server, lazy, errors } = loader._.options;
			if (!server && isSSR) return;
			const ret = scope.run(() => app.runWithContext(() => loader._.load(to, router, from)));
			return !isSSR && toLazyValue(lazy, to, from) ? void 0 : ret.catch((reason) => {
				if (errors === true) {
					if (Array.isArray(globalErrors) ? globalErrors.some((Err) => reason instanceof Err) : globalErrors(reason)) return;
				} else if (errors && (Array.isArray(errors) ? errors.some((Err) => reason instanceof Err) : errors(reason))) return;
				throw reason;
			});
		})).then(() => {
			if (to.meta[NAVIGATION_RESULTS_KEY].length) return selectNavigationResult(to.meta[NAVIGATION_RESULTS_KEY]);
		}).catch((error) => error instanceof NavigationResult ? error.value : signal.aborted && error === signal.reason ? false : Promise.reject(error)).finally(() => {
			setCurrentContext([]);
			isDataLoading.value = false;
		});
	});
	const removeAfterEach = router.afterEach((to, from, failure) => {
		if (failure) {
			to.meta[ABORT_CONTROLLER_KEY]?.abort(failure);
			if (isNavigationFailure(failure, 16)) for (const loader of to.meta[LOADER_SET_KEY]) loader._.getEntry(router).resetPending();
		} else for (const loader of to.meta[LOADER_SET_KEY]) {
			const { commit, lazy } = loader._.options;
			if (commit === "after-load") {
				const entry = loader._.getEntry(router);
				if (entry && (!toLazyValue(lazy, to, from) || !entry.isLoading.value)) entry.commit(to);
			}
		}
		if (router[PENDING_LOCATION_KEY] === to) router[PENDING_LOCATION_KEY] = null;
	});
	const removeOnError = router.onError((error, to) => {
		to.meta[ABORT_CONTROLLER_KEY]?.abort(error);
		if (router[PENDING_LOCATION_KEY] === to) router[PENDING_LOCATION_KEY] = null;
	});
	return () => {
		delete router[LOADER_ENTRIES_KEY];
		delete router[APP_KEY];
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
	const effect = effectScope(true);
	const removeGuards = setupLoaderGuard(assign({
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
	return inject(IS_DATA_LOADING_KEY);
}

//#endregion
export { ABORT_CONTROLLER_KEY, APP_KEY, DATA_LOADERS_EFFECT_SCOPE_KEY, DataLoaderPlugin, IS_SSR_KEY, IS_USE_DATA_LOADER_KEY, LOADER_ENTRIES_KEY, LOADER_SET_KEY, NAVIGATION_RESULTS_KEY, NavigationResult, PENDING_LOCATION_KEY, STAGED_NO_VALUE, assign, currentContext, getCurrentContext, isSubsetOf, setCurrentContext, toLazyValue, trackRoute, useIsDataLoading, withLoaderContext };