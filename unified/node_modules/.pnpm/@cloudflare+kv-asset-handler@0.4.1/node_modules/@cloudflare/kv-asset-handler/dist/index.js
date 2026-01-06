import * as mime from 'mime';

var __defProp = Object.defineProperty;
var __name = (target, value) => __defProp(target, "name", { value, configurable: true });

// src/types.ts
var KVError = class _KVError extends Error {
  static {
    __name(this, "KVError");
  }
  constructor(message, status = 500) {
    super(message);
    Object.setPrototypeOf(this, new.target.prototype);
    this.name = _KVError.name;
    this.status = status;
  }
  status;
};
var MethodNotAllowedError = class extends KVError {
  static {
    __name(this, "MethodNotAllowedError");
  }
  constructor(message = `Not a valid request method`, status = 405) {
    super(message, status);
  }
};
var NotFoundError = class extends KVError {
  static {
    __name(this, "NotFoundError");
  }
  constructor(message = `Not Found`, status = 404) {
    super(message, status);
  }
};
var InternalError = class extends KVError {
  static {
    __name(this, "InternalError");
  }
  constructor(message = `Internal Error in KV Asset Handler`, status = 500) {
    super(message, status);
  }
};

// src/index.ts
var defaultCacheControl = {
  browserTTL: null,
  edgeTTL: 2 * 60 * 60 * 24,
  // 2 days
  bypassCache: false
  // do not bypass Cloudflare's cache
};
var parseStringAsObject = /* @__PURE__ */ __name((maybeString) => typeof maybeString === "string" ? JSON.parse(maybeString) : maybeString, "parseStringAsObject");
function getAssetFromKVDefaultOptions() {
  return {
    ASSET_NAMESPACE: typeof __STATIC_CONTENT !== "undefined" ? __STATIC_CONTENT : void 0,
    ASSET_MANIFEST: typeof __STATIC_CONTENT_MANIFEST !== "undefined" ? parseStringAsObject(__STATIC_CONTENT_MANIFEST) : {},
    cacheControl: defaultCacheControl,
    defaultMimeType: "text/plain",
    defaultDocument: "index.html",
    pathIsEncoded: false,
    defaultETag: "strong"
  };
}
__name(getAssetFromKVDefaultOptions, "getAssetFromKVDefaultOptions");
function assignOptions(options) {
  return Object.assign({}, getAssetFromKVDefaultOptions(), options);
}
__name(assignOptions, "assignOptions");
var mapRequestToAsset = /* @__PURE__ */ __name((request, options) => {
  options = assignOptions(options);
  const parsedUrl = new URL(request.url);
  let pathname = parsedUrl.pathname;
  if (pathname.endsWith("/")) {
    pathname = pathname.concat(options.defaultDocument);
  } else if (!mime.getType(pathname)) {
    pathname = pathname.concat("/" + options.defaultDocument);
  }
  parsedUrl.pathname = pathname;
  return new Request(parsedUrl.toString(), request);
}, "mapRequestToAsset");
function serveSinglePageApp(request, options) {
  options = assignOptions(options);
  request = mapRequestToAsset(request, options);
  const parsedUrl = new URL(request.url);
  if (parsedUrl.pathname.endsWith(".html")) {
    return new Request(
      `${parsedUrl.origin}/${options.defaultDocument}`,
      request
    );
  } else {
    return request;
  }
}
__name(serveSinglePageApp, "serveSinglePageApp");
var getAssetFromKV = /* @__PURE__ */ __name(async (event, options) => {
  options = assignOptions(options);
  const request = event.request;
  const ASSET_NAMESPACE = options.ASSET_NAMESPACE;
  const ASSET_MANIFEST = parseStringAsObject(
    options.ASSET_MANIFEST
  );
  if (typeof ASSET_NAMESPACE === "undefined") {
    throw new InternalError(`there is no KV namespace bound to the script`);
  }
  const rawPathKey = new URL(request.url).pathname.replace(/^\/+/, "");
  let pathIsEncoded = options.pathIsEncoded;
  let requestKey;
  if (options.mapRequestToAsset) {
    requestKey = options.mapRequestToAsset(request);
  } else if (ASSET_MANIFEST[rawPathKey]) {
    requestKey = request;
  } else if (ASSET_MANIFEST[decodeURIComponent(rawPathKey)]) {
    pathIsEncoded = true;
    requestKey = request;
  } else {
    const mappedRequest = mapRequestToAsset(request);
    const mappedRawPathKey = new URL(mappedRequest.url).pathname.replace(
      /^\/+/,
      ""
    );
    if (ASSET_MANIFEST[decodeURIComponent(mappedRawPathKey)]) {
      pathIsEncoded = true;
      requestKey = mappedRequest;
    } else {
      requestKey = mapRequestToAsset(request, options);
    }
  }
  const SUPPORTED_METHODS = ["GET", "HEAD"];
  if (!SUPPORTED_METHODS.includes(requestKey.method)) {
    throw new MethodNotAllowedError(
      `${requestKey.method} is not a valid request method`
    );
  }
  const parsedUrl = new URL(requestKey.url);
  const pathname = pathIsEncoded ? decodeURIComponent(parsedUrl.pathname) : parsedUrl.pathname;
  let pathKey = pathname.replace(/^\/+/, "");
  const cache = caches.default;
  let mimeType = mime.getType(pathKey) || options.defaultMimeType;
  if (mimeType.startsWith("text") || mimeType === "application/javascript") {
    mimeType += "; charset=utf-8";
  }
  let shouldEdgeCache = false;
  if (typeof ASSET_MANIFEST !== "undefined") {
    if (ASSET_MANIFEST[pathKey]) {
      pathKey = ASSET_MANIFEST[pathKey];
      shouldEdgeCache = true;
    }
  }
  const cacheKey = new Request(`${parsedUrl.origin}/${pathKey}`, request);
  const evalCacheOpts = (() => {
    switch (typeof options.cacheControl) {
      case "function":
        return options.cacheControl(request);
      case "object":
        return options.cacheControl;
      default:
        return defaultCacheControl;
    }
  })();
  const formatETag = /* @__PURE__ */ __name((entityId = pathKey, validatorType = options.defaultETag) => {
    if (!entityId) {
      return "";
    }
    switch (validatorType) {
      case "weak":
        if (!entityId.startsWith("W/")) {
          if (entityId.startsWith(`"`) && entityId.endsWith(`"`)) {
            return `W/${entityId}`;
          }
          return `W/"${entityId}"`;
        }
        return entityId;
      case "strong":
        if (entityId.startsWith(`W/"`)) {
          entityId = entityId.replace("W/", "");
        }
        if (!entityId.endsWith(`"`)) {
          entityId = `"${entityId}"`;
        }
        return entityId;
      default:
        return "";
    }
  }, "formatETag");
  options.cacheControl = Object.assign({}, defaultCacheControl, evalCacheOpts);
  if (options.cacheControl.bypassCache || options.cacheControl.edgeTTL === null || request.method == "HEAD") {
    shouldEdgeCache = false;
  }
  const shouldSetBrowserCache = typeof options.cacheControl.browserTTL === "number";
  let response = null;
  if (shouldEdgeCache) {
    response = await cache.match(cacheKey);
  }
  if (response) {
    if (response.status > 300 && response.status < 400) {
      if (response.body && "cancel" in Object.getPrototypeOf(response.body)) {
        response.body.cancel();
      }
      response = new Response(null, response);
    } else {
      const opts = {
        headers: new Headers(response.headers),
        status: 0,
        statusText: ""
      };
      opts.headers.set("cf-cache-status", "HIT");
      if (response.status) {
        opts.status = response.status;
        opts.statusText = response.statusText;
      } else if (opts.headers.has("Content-Range")) {
        opts.status = 206;
        opts.statusText = "Partial Content";
      } else {
        opts.status = 200;
        opts.statusText = "OK";
      }
      response = new Response(response.body, opts);
    }
  } else {
    const body = await ASSET_NAMESPACE.get(pathKey, "arrayBuffer");
    if (body === null) {
      throw new NotFoundError(
        `could not find ${pathKey} in your content namespace`
      );
    }
    response = new Response(body);
    if (shouldEdgeCache) {
      response.headers.set("Accept-Ranges", "bytes");
      response.headers.set("Content-Length", String(body.byteLength));
      if (!response.headers.has("etag")) {
        response.headers.set("etag", formatETag(pathKey));
      }
      response.headers.set(
        "Cache-Control",
        `max-age=${options.cacheControl.edgeTTL}`
      );
      event.waitUntil(cache.put(cacheKey, response.clone()));
      response.headers.set("CF-Cache-Status", "MISS");
    }
  }
  response.headers.set("Content-Type", mimeType);
  if (response.status === 304) {
    const etag = formatETag(response.headers.get("etag"));
    const ifNoneMatch = cacheKey.headers.get("if-none-match");
    const proxyCacheStatus = response.headers.get("CF-Cache-Status");
    if (etag) {
      if (ifNoneMatch && ifNoneMatch === etag && proxyCacheStatus === "MISS") {
        response.headers.set("CF-Cache-Status", "EXPIRED");
      } else {
        response.headers.set("CF-Cache-Status", "REVALIDATED");
      }
      response.headers.set("etag", formatETag(etag, "weak"));
    }
  }
  if (shouldSetBrowserCache) {
    response.headers.set(
      "Cache-Control",
      `max-age=${options.cacheControl.browserTTL}`
    );
  } else {
    response.headers.delete("Cache-Control");
  }
  return response;
}, "getAssetFromKV");

export { InternalError, MethodNotAllowedError, NotFoundError, getAssetFromKV, mapRequestToAsset, serveSinglePageApp };
//# sourceMappingURL=index.js.map
//# sourceMappingURL=index.js.map