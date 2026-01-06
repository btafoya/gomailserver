import { t as toLazyValue } from "./createDataLoader-BK9Gdnky.mjs";
import { useRoute, useRouter } from "vue-router";
import { ABORT_CONTROLLER_KEY, APP_KEY, IS_SSR_KEY, IS_USE_DATA_LOADER_KEY, LOADER_ENTRIES_KEY, LOADER_SET_KEY, NAVIGATION_RESULTS_KEY, NavigationResult, PENDING_LOCATION_KEY, STAGED_NO_VALUE, getCurrentContext, setCurrentContext } from "unplugin-vue-router/data-loaders";
import { shallowRef } from "vue";
import "scule";

//#region src/core/utils.ts
function warn(msg, type = "warn") {
	console[type](`⚠️  [unplugin-vue-router]: ${msg}`);
}
/**
* Type safe alternative to Array.isArray
* https://github.com/microsoft/TypeScript/pull/48228
*/
const isArray = Array.isArray;

//#endregion
//#region src/data-loaders/defineLoader.ts
function defineBasicLoader(nameOrLoader, _loaderOrOptions, opts) {
	const loader = typeof nameOrLoader === "function" ? nameOrLoader : _loaderOrOptions;
	opts = typeof _loaderOrOptions === "object" ? _loaderOrOptions : opts;
	const options = {
		...DEFAULT_DEFINE_LOADER_OPTIONS,
		...opts,
		commit: opts?.commit || DEFAULT_DEFINE_LOADER_OPTIONS.commit
	};
	function load(to, router, from, parent) {
		const entries = router[LOADER_ENTRIES_KEY];
		const isSSR = router[IS_SSR_KEY];
		if (!entries.has(loader)) entries.set(loader, {
			data: shallowRef(),
			isLoading: shallowRef(false),
			error: shallowRef(null),
			to,
			options,
			children: /* @__PURE__ */ new Set(),
			resetPending() {
				this.pendingLoad = null;
				this.pendingTo = null;
			},
			pendingLoad: null,
			pendingTo: null,
			staged: STAGED_NO_VALUE,
			stagedError: null,
			stagedNavigationResult: null,
			commit
		});
		const entry = entries.get(loader);
		if (entry.pendingTo === to && entry.pendingLoad) return entry.pendingLoad;
		const { error, isLoading, data } = entry;
		const initialRootData = router[INITIAL_DATA_KEY];
		const key = options.key || "";
		let initialData = STAGED_NO_VALUE;
		if (initialRootData && key in initialRootData) {
			initialData = initialRootData[key];
			delete initialRootData[key];
		}
		if (initialData !== STAGED_NO_VALUE) {
			data.value = initialData;
			return entry.pendingLoad = Promise.resolve();
		}
		entry.pendingTo = to;
		isLoading.value = true;
		const currentContext = getCurrentContext();
		if (process.env.NODE_ENV !== "production") {
			if (parent !== currentContext[0]) console.warn(`❌👶 "${options.key}" has a different parent than the current context. This shouldn't be happening. Please report a bug with a reproduction to https://github.com/posva/unplugin-vue-router/`);
		}
		setCurrentContext([
			entry,
			router,
			to
		]);
		entry.staged = STAGED_NO_VALUE;
		entry.stagedError = error.value;
		entry.stagedNavigationResult = null;
		const currentLoad = Promise.resolve(loader(to, { signal: to.meta[ABORT_CONTROLLER_KEY]?.signal })).then((d) => {
			if (entry.pendingLoad === currentLoad) if (d instanceof NavigationResult) {
				to.meta[NAVIGATION_RESULTS_KEY].push(d);
				entry.stagedNavigationResult = d;
				if (process.env.NODE_ENV !== "production") warnNonExposedLoader({
					to,
					options,
					useDataLoader
				});
			} else {
				entry.staged = d;
				entry.stagedError = null;
			}
		}).catch((error$1) => {
			if (entry.pendingLoad === currentLoad) {
				if (process.env.NODE_ENV !== "production") {
					if (error$1 instanceof NavigationResult) warnNonExposedLoader({
						to,
						options,
						useDataLoader
					});
				}
				entry.stagedError = error$1;
				if (!toLazyValue(options.lazy, to, from) || isSSR) throw error$1;
			}
		}).finally(() => {
			setCurrentContext(currentContext);
			if (entry.pendingLoad === currentLoad) {
				isLoading.value = false;
				if (options.commit === "immediate" || !router[PENDING_LOCATION_KEY]) entry.commit(to);
			}
		});
		setCurrentContext(currentContext);
		entry.pendingLoad = currentLoad;
		return currentLoad;
	}
	function commit(to) {
		if (this.pendingTo === to) {
			if (process.env.NODE_ENV !== "production") {
				if (this.staged === STAGED_NO_VALUE && this.stagedError === null && this.stagedNavigationResult === null) console.warn(`Loader "${options.key}"'s "commit()" was called but there is no staged data.`);
			}
			if (this.staged !== STAGED_NO_VALUE) this.data.value = this.staged;
			this.error.value = this.stagedError;
			this.staged = STAGED_NO_VALUE;
			this.stagedError = this.error.value;
			this.pendingTo = null;
			this.to = to;
			for (const childEntry of this.children) childEntry.commit(to);
		}
	}
	const useDataLoader = () => {
		const currentContext = getCurrentContext();
		const [parentEntry, _router, _route] = currentContext;
		const router = _router || useRouter();
		const route = _route || useRoute();
		const entries = router[LOADER_ENTRIES_KEY];
		let entry = entries.get(loader);
		if (!entry || parentEntry && entry.pendingTo !== route || !entry.pendingLoad) router[APP_KEY].runWithContext(() => load(route, router, void 0, parentEntry));
		entry = entries.get(loader);
		if (parentEntry) {
			if (parentEntry === entry) console.warn(`👶❌ "${options.key}" has itself as parent. This shouldn't be happening. Please report a bug with a reproduction to https://github.com/posva/unplugin-vue-router/`);
			parentEntry.children.add(entry);
		}
		const { data, error, isLoading } = entry;
		const useDataLoaderResult = {
			data,
			error,
			isLoading,
			reload: (to = router.currentRoute.value) => router[APP_KEY].runWithContext(() => load(to, router)).then(() => entry.commit(to))
		};
		const promise = entry.pendingLoad.then(() => {
			return entry.staged === STAGED_NO_VALUE ? entry.stagedNavigationResult ? Promise.reject(entry.stagedNavigationResult) : data.value : entry.staged;
		}).catch((e) => parentEntry ? Promise.reject(e) : null);
		setCurrentContext(currentContext);
		return Object.assign(promise, useDataLoaderResult);
	};
	useDataLoader[IS_USE_DATA_LOADER_KEY] = true;
	useDataLoader._ = {
		load,
		options,
		getEntry(router) {
			return router[LOADER_ENTRIES_KEY].get(loader);
		}
	};
	return useDataLoader;
}
/**
* Dev only warning for loaders that return/throw NavigationResult but are not exposed
*
* @param to - target location
* @param options - options used to define the loader
* @param useDataLoader - the data loader composable
*/
function warnNonExposedLoader({ to, options, useDataLoader }) {
	const loaders = to.meta[LOADER_SET_KEY];
	if (loaders && !loaders.has(useDataLoader)) warn("A loader returned a NavigationResult but is not registered on the route. Did you forget to \"export\" it from the page component?" + (options.key ? ` (loader key: "${options.key}")` : ""));
}
const DEFAULT_DEFINE_LOADER_OPTIONS = {
	lazy: false,
	server: true,
	commit: "after-load"
};
/**
* Initial data generated on server and consumed on client.
* @internal
*/
const INITIAL_DATA_KEY = Symbol();

//#endregion
export { defineBasicLoader };