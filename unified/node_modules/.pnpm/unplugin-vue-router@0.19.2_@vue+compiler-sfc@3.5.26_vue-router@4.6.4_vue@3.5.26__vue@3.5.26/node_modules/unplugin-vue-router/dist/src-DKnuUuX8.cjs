//#region rolldown:runtime
var __create = Object.create;
var __defProp = Object.defineProperty;
var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
var __getOwnPropNames = Object.getOwnPropertyNames;
var __getProtoOf = Object.getPrototypeOf;
var __hasOwnProp = Object.prototype.hasOwnProperty;
var __copyProps = (to, from, except, desc) => {
	if (from && typeof from === "object" || typeof from === "function") {
		for (var keys = __getOwnPropNames(from), i = 0, n = keys.length, key; i < n; i++) {
			key = keys[i];
			if (!__hasOwnProp.call(to, key) && key !== except) {
				__defProp(to, key, {
					get: ((k) => from[k]).bind(null, key),
					enumerable: !(desc = __getOwnPropDesc(from, key)) || desc.enumerable
				});
			}
		}
	}
	return to;
};
var __toESM = (mod, isNodeMode, target) => (target = mod != null ? __create(__getProtoOf(mod)) : {}, __copyProps(isNodeMode || !mod || !mod.__esModule ? __defProp(target, "default", {
	value: mod,
	enumerable: true
}) : target, mod));

//#endregion
const require_options = require('./options-CpfNkO-E.cjs');
let unplugin = require("unplugin");
let scule = require("scule");
let node_fs = require("node:fs");
let pathe = require("pathe");
pathe = __toESM(pathe);
let __vue_macros_common = require("@vue-macros/common");
let tinyglobby = require("tinyglobby");
let __vue_compiler_sfc = require("@vue/compiler-sfc");
let json5 = require("json5");
json5 = __toESM(json5);
let yaml = require("yaml");
let chokidar = require("chokidar");
let picomatch = require("picomatch");
picomatch = __toESM(picomatch);
let __babel_generator = require("@babel/generator");
let ast_walker_scope = require("ast-walker-scope");
let mlly = require("mlly");
let unplugin_utils = require("unplugin-utils");
let magic_string = require("magic-string");
magic_string = __toESM(magic_string);

//#region src/utils/encoding.ts
/**
* Encoding Rules (␣ = Space)
* - Path: ␣ " < > # ? { }
* - Query: ␣ " < > # & =
* - Hash: ␣ " < > `
*
* On top of that, the RFC3986 (https://tools.ietf.org/html/rfc3986#section-2.2)
* defines some extra characters to be encoded. Most browsers do not encode them
* in encodeURI https://github.com/whatwg/url/issues/369, so it may be safer to
* also encode `!'()*`. Leaving un-encoded only ASCII alphanumeric(`a-zA-Z0-9`)
* plus `-._~`. This extra safety should be applied to query by patching the
* string returned by encodeURIComponent encodeURI also encodes `[\]^`. `\`
* should be encoded to avoid ambiguity. Browsers (IE, FF, C) transform a `\`
* into a `/` if directly typed in. The _backtick_ (`````) should also be
* encoded everywhere because some browsers like FF encode it when directly
* written while others don't. Safari and IE don't encode ``"<>{}``` in hash.
*/
const HASH_RE = /#/g;
const IM_RE = /\?/g;
/**
* NOTE: It's not clear to me if we should encode the + symbol in queries, it
* seems to be less flexible than not doing so and I can't find out the legacy
* systems requiring this for regular requests like text/html. In the standard,
* the encoding of the plus character is only mentioned for
* application/x-www-form-urlencoded
* (https://url.spec.whatwg.org/#urlencoded-parsing) and most browsers seems lo
* leave the plus character as is in queries. To be more flexible, we allow the
* plus character on the query, but it can also be manually encoded by the user.
*
* Resources:
* - https://url.spec.whatwg.org/#urlencoded-parsing
* - https://stackoverflow.com/questions/1634271/url-encoding-the-space-character-or-20
*/
const ENC_BRACKET_OPEN_RE = /%5B/g;
const ENC_BRACKET_CLOSE_RE = /%5D/g;
const ENC_PIPE_RE = /%7C/g;
/**
* Encode characters that need to be encoded on the path, search and hash
* sections of the URL.
*
* @internal
* @param text - string to encode
* @returns encoded string
*/
function commonEncode(text) {
	return text == null ? "" : encodeURI("" + text).replace(ENC_PIPE_RE, "|").replace(ENC_BRACKET_OPEN_RE, "[").replace(ENC_BRACKET_CLOSE_RE, "]");
}
/**
* Encode characters that need to be encoded on the path section of the URL.
*
* @param text - string to encode
* @returns encoded string
*/
function encodePath(text) {
	return commonEncode(text).replace(HASH_RE, "%23").replace(IM_RE, "%3F");
}

//#endregion
//#region src/core/treeNodeValue.ts
let TreeNodeType = /* @__PURE__ */ function(TreeNodeType$1) {
	TreeNodeType$1[TreeNodeType$1["static"] = 0] = "static";
	TreeNodeType$1[TreeNodeType$1["group"] = 1] = "group";
	TreeNodeType$1[TreeNodeType$1["param"] = 2] = "param";
	return TreeNodeType$1;
}({});
const EDITS_OVERRIDE_NAME = "@@edits";
var _TreeNodeValueBase = class {
	/**
	* flag based on the type of the segment
	*/
	_type;
	parent;
	/**
	* segment as defined by the file structure e.g. keeps the `index` name, `(group-name)`
	*/
	rawSegment;
	/**
	* transformed version of the segment into a vue-router path. e.g. `'index'` becomes `''` and `[param]` becomes
	* `:param`, `prefix-[param]-end` becomes `prefix-:param-end`.
	*/
	pathSegment;
	/**
	* Array of sub segments. This is usually one single elements but can have more for paths like `prefix-[param]-end.vue`
	*/
	subSegments;
	/**
	* Overrides defined by each file. The map is necessary to handle named views.
	*/
	_overrides = /* @__PURE__ */ new Map();
	/**
	* View name (Vue Router feature) mapped to their corresponding file. By default, the view name is `default` unless
	* specified with a `@` e.g. `index@aux.vue` will have a view name of `aux`.
	*/
	components = /* @__PURE__ */ new Map();
	constructor(rawSegment, parent, pathSegment = rawSegment, subSegments = [pathSegment]) {
		this._type = 0;
		this.rawSegment = rawSegment;
		this.pathSegment = pathSegment;
		this.subSegments = subSegments;
		this.parent = parent;
	}
	/**
	* Path of the node. Can be absolute or not. If it has been overridden, it
	* will return the overridden path.
	*/
	get path() {
		return this.overrides.path ?? this.pathSegment;
	}
	/**
	* Full path of the node including parent nodes.
	*/
	get fullPath() {
		const pathSegment = this.path;
		if (pathSegment.startsWith("/")) return pathSegment;
		return require_options.joinPath(this.parent?.fullPath ?? "", pathSegment);
	}
	/**
	* Gets all the query params for the node. This does not include params from parent nodes.
	*/
	get queryParams() {
		const paramsQuery = this.overrides.params?.query;
		if (!paramsQuery) return [];
		const queryParams = [];
		for (var paramName in paramsQuery) {
			var config = paramsQuery[paramName];
			if (!config) continue;
			if (typeof config === "string") queryParams.push({
				paramName,
				parser: config,
				format: "value"
			});
			else queryParams.push({
				paramName,
				parser: config.parser || null,
				format: config.format || "value",
				defaultValue: config.default
			});
		}
		return queryParams;
	}
	/**
	* Gets all the params for the node including path and query params. This
	* does not include params from parent nodes.
	*/
	get params() {
		return [...this.isParam() ? this.pathParams : [], ...this.queryParams];
	}
	toString() {
		if (!this.pathSegment) return "<index>" + (this.rawSegment === "index" ? "" : " " + this.rawSegment);
		return this.pathSegment;
	}
	isParam() {
		return !!(this._type & TreeNodeType.param);
	}
	isStatic() {
		return this._type === TreeNodeType.static;
	}
	isGroup() {
		return this._type === TreeNodeType.group;
	}
	get overrides() {
		return [...this._overrides.entries()].sort(([nameA], [nameB]) => nameA === nameB ? 0 : nameA !== EDITS_OVERRIDE_NAME && (nameA < nameB || nameB === EDITS_OVERRIDE_NAME) ? -1 : 1).reduce((acc, [_path, routeBlock]) => {
			return require_options.mergeRouteRecordOverride(acc, routeBlock);
		}, {});
	}
	setOverride(filePath, routeBlock) {
		this._overrides.set(filePath, routeBlock || {});
	}
	/**
	* Remove all overrides for a given key.
	*
	* @param key - key to remove from the override, e.g. path, name, etc
	*/
	removeOverride(key) {
		for (const [_filePath, routeBlock] of this._overrides) delete routeBlock[key];
	}
	/**
	* Add an override to the current node by merging with the existing values.
	*
	* @param filePath - The file path to add to the override
	* @param routeBlock - The route block to add to the override
	*/
	mergeOverride(filePath, routeBlock) {
		const existing = this._overrides.get(filePath) || {};
		this._overrides.set(filePath, require_options.mergeRouteRecordOverride(existing, routeBlock));
	}
	/**
	* Add an override to the current node using the special file path `@@edits` that makes this added at build time.
	*
	* @param routeBlock -  The route block to add to the override
	*/
	addEditOverride(routeBlock) {
		return this.mergeOverride(EDITS_OVERRIDE_NAME, routeBlock);
	}
	/**
	* Set a specific value in the _edits_ override.
	*
	* @param key - key to set in the override, e.g. path, name, etc
	* @param value - value to set in the override
	*/
	setEditOverride(key, value) {
		if (!this._overrides.has(EDITS_OVERRIDE_NAME)) this._overrides.set(EDITS_OVERRIDE_NAME, {});
		const existing = this._overrides.get(EDITS_OVERRIDE_NAME);
		existing[key] = value;
	}
};
/**
* - Static
* - Static + Custom Param (subSegments)
* - Static + Param (subSegments)
* - Custom Param
* - Param
* - CatchAll
*/
/**
* Static path like `/users`, `/users/list`, etc
* @extends _TreeNodeValueBase
*/
var TreeNodeValueStatic = class extends _TreeNodeValueBase {
	_type = TreeNodeType.static;
	score = [300];
	constructor(rawSegment, parent, pathSegment = rawSegment) {
		super(rawSegment, parent, pathSegment);
	}
};
var TreeNodeValueGroup = class extends _TreeNodeValueBase {
	_type = TreeNodeType.group;
	groupName;
	score = [300];
	constructor(rawSegment, parent, pathSegment, groupName) {
		super(rawSegment, parent, pathSegment);
		this.groupName = groupName;
	}
};
/**
* Checks if a TreePathParam or TreeQueryParam is optional.
*
* @internal
*/
function isTreeParamOptional(param) {
	if ("optional" in param) return param.optional;
	return param.defaultValue !== void 0;
}
/**
* Checks if a TreePathParam or TreeQueryParam is repeatable (array).
*
* @internal
*/
function isTreeParamRepeatable(param) {
	if ("repeatable" in param) return param.repeatable;
	return param.format === "array";
}
/**
* Checks if a param is a TreePathParam.
*
* @internal
*/
function isTreePathParam(param) {
	return "modifier" in param;
}
/**
* To escape regex characters in the path segment.
* @internal
*/
const REGEX_CHARS_RE = /[.+*?^${}()[\]/\\]/g;
/**
* Escapes regex characters in a string to be used in a regex pattern.
* @param str - The string to escape.
*
* @internal
*/
const escapeRegex = (str) => str.replace(REGEX_CHARS_RE, "\\$&");
var TreeNodeValueParam = class extends _TreeNodeValueBase {
	_type = TreeNodeType.param;
	constructor(rawSegment, parent, pathParams, pathSegment, subSegments) {
		super(rawSegment, parent, pathSegment, subSegments);
		this.pathParams = pathParams;
	}
	get score() {
		return this.subSegments.map((segment) => {
			if (typeof segment === "string") return 300;
			else return 80 - (segment.isSplat ? 500 : (segment.optional ? 10 : 0) + (segment.repeatable ? 20 : 0));
		});
	}
	/**
	* Generates the regex pattern for the path segment.
	*/
	get re() {
		let regexp = "";
		for (var i = 0; i < this.subSegments.length; i++) {
			var segment = this.subSegments[i];
			if (!segment) continue;
			if (typeof segment === "string") regexp += escapeRegex(segment);
			else if (segment.isSplat) regexp += "(.*)";
			else {
				var re = segment.repeatable ? "(.+?)" : "([^/]+?)";
				if (segment.optional) {
					var prevSegment = this.subSegments[i - 1];
					if ((!prevSegment || typeof prevSegment === "string" && prevSegment.endsWith("/")) && this.subSegments.length > 1) {
						re = `(?:\\/${re})?`;
						regexp = regexp.slice(0, -2);
					} else re += "?";
				}
				regexp += re;
			}
		}
		return regexp;
	}
	toString() {
		const params = this.params.length > 0 ? ` 𝑥(` + this.params.map((p) => ("format" in p ? "?" : "") + `${p.paramName}${"modifier" in p ? p.modifier : ""}` + (p.parser ? "=" + p.parser : "")).join(", ") + ")" : "";
		return `${this.pathSegment}` + params;
	}
};
/**
* Resolves the options for the TreeNodeValue.
*
* @param options - options to resolve
* @returns resolved options
*/
function resolveTreeNodeValueOptions(options) {
	return {
		format: "file",
		dotNesting: true,
		...options
	};
}
/**
* Creates a new TreeNodeValue based on the segment. The result can be a static segment, group segment or a param segment.
*
* @param segment - path segment
* @param parent - parent node
* @param options - options
*/
function createTreeNodeValue(segment, parent, opts = {}) {
	if (!segment || segment === "index") return new TreeNodeValueStatic(segment, parent, "");
	const options = resolveTreeNodeValueOptions(opts);
	const openingPar = segment.indexOf("(");
	if (options.format === "file" && openingPar >= 0) {
		let groupName;
		const closingPar = segment.lastIndexOf(")");
		if (closingPar < 0 || closingPar < openingPar) {
			require_options.warn(`Segment "${segment}" is missing the closing ")". It will be treated as a static segment.`);
			return new TreeNodeValueStatic(segment, parent, segment);
		}
		groupName = segment.slice(openingPar + 1, closingPar);
		const before = segment.slice(0, openingPar);
		const after = segment.slice(closingPar + 1);
		if (!before && !after) return new TreeNodeValueGroup(segment, parent, "", groupName);
	}
	const [pathSegment, pathParams, subSegments] = options.format === "path" ? parseRawPathSegment(segment) : parseFileSegment(segment, options);
	if (pathParams.length) return new TreeNodeValueParam(segment, parent, pathParams, pathSegment, subSegments);
	return new TreeNodeValueStatic(segment, parent, pathSegment);
}
var ParseFileSegmentState = /* @__PURE__ */ function(ParseFileSegmentState$1) {
	ParseFileSegmentState$1[ParseFileSegmentState$1["static"] = 0] = "static";
	ParseFileSegmentState$1[ParseFileSegmentState$1["paramOptional"] = 1] = "paramOptional";
	ParseFileSegmentState$1[ParseFileSegmentState$1["param"] = 2] = "param";
	ParseFileSegmentState$1[ParseFileSegmentState$1["paramParser"] = 3] = "paramParser";
	ParseFileSegmentState$1[ParseFileSegmentState$1["modifier"] = 4] = "modifier";
	ParseFileSegmentState$1[ParseFileSegmentState$1["charCode"] = 5] = "charCode";
	return ParseFileSegmentState$1;
}(ParseFileSegmentState || {});
const IS_VARIABLE_CHAR_RE = /[0-9a-zA-Z_]/;
/**
* Parses a segment into the route path segment and the extracted params.
*
* @param segment - segment to parse without the extension
* @returns - the pathSegment and the params
*/
function parseFileSegment(segment, { dotNesting }) {
	let buffer = "";
	let paramParserBuffer = "";
	let state = ParseFileSegmentState.static;
	const params = [];
	let pathSegment = "";
	const subSegments = [];
	let currentTreeRouteParam = createEmptyRouteParam();
	let pos = 0;
	let c;
	function consumeBuffer() {
		if (state === ParseFileSegmentState.static) {
			const encodedBuffer = buffer.split("/").map((part) => encodePath(part)).join("/");
			pathSegment += encodedBuffer;
			subSegments.push(encodedBuffer);
		} else if (state === ParseFileSegmentState.modifier) {
			currentTreeRouteParam.paramName = buffer;
			currentTreeRouteParam.parser = paramParserBuffer || null;
			currentTreeRouteParam.modifier = currentTreeRouteParam.optional ? currentTreeRouteParam.repeatable ? "*" : "?" : currentTreeRouteParam.repeatable ? "+" : "";
			buffer = "";
			paramParserBuffer = "";
			pathSegment += `:${currentTreeRouteParam.paramName}${currentTreeRouteParam.isSplat ? "(.*)" : pos < segment.length - 1 && IS_VARIABLE_CHAR_RE.test(segment[pos + 1]) ? "()" : ""}${currentTreeRouteParam.modifier}`;
			params.push(currentTreeRouteParam);
			subSegments.push(currentTreeRouteParam);
			currentTreeRouteParam = createEmptyRouteParam();
		} else if (state === ParseFileSegmentState.charCode) {
			if (buffer.length !== 2) throw new SyntaxError(`Invalid character code in segment "${segment}". Hex code must be exactly 2 digits, got "${buffer}"`);
			const hexCode = parseInt(buffer, 16);
			if (!Number.isInteger(hexCode) || hexCode < 0 || hexCode > 255) throw new SyntaxError(`Invalid hex code "${buffer}" in segment "${segment}"`);
			pathSegment += String.fromCharCode(hexCode);
		}
		buffer = "";
	}
	for (pos = 0; pos < segment.length; pos++) {
		c = segment[pos];
		if (state === ParseFileSegmentState.static) if (c === "[") {
			if (buffer) consumeBuffer();
			state = ParseFileSegmentState.paramOptional;
		} else buffer += dotNesting && c === "." ? "/" : c;
		else if (state === ParseFileSegmentState.paramOptional) {
			if (c === "[") currentTreeRouteParam.optional = true;
			else if (c === ".") {
				currentTreeRouteParam.isSplat = true;
				pos += 2;
			} else buffer += c;
			state = ParseFileSegmentState.param;
		} else if (state === ParseFileSegmentState.param) if (c === "]") {
			if (currentTreeRouteParam.optional) pos++;
			state = ParseFileSegmentState.modifier;
		} else if (c === ".") {
			currentTreeRouteParam.isSplat = true;
			pos += 2;
		} else if (c === "=") {
			state = ParseFileSegmentState.paramParser;
			paramParserBuffer = "";
		} else if (c === "+" && buffer === "x" && !currentTreeRouteParam.isSplat && !currentTreeRouteParam.optional) {
			buffer = "";
			state = ParseFileSegmentState.charCode;
		} else buffer += c;
		else if (state === ParseFileSegmentState.modifier) {
			if (c === "+") currentTreeRouteParam.repeatable = true;
			else pos--;
			consumeBuffer();
			state = ParseFileSegmentState.static;
		} else if (state === ParseFileSegmentState.paramParser) if (c === "]") {
			if (currentTreeRouteParam.optional) pos++;
			state = ParseFileSegmentState.modifier;
		} else paramParserBuffer += c;
		else if (state === ParseFileSegmentState.charCode) if (c === "]") {
			consumeBuffer();
			state = ParseFileSegmentState.static;
		} else buffer += c;
	}
	if (state === ParseFileSegmentState.param || state === ParseFileSegmentState.paramOptional || state === ParseFileSegmentState.paramParser || state === ParseFileSegmentState.charCode) throw new SyntaxError(`Invalid segment: "${segment}"`);
	if (buffer) consumeBuffer();
	return [
		pathSegment,
		params,
		subSegments
	];
}
var ParseRawPathSegmentState = /* @__PURE__ */ function(ParseRawPathSegmentState$1) {
	ParseRawPathSegmentState$1[ParseRawPathSegmentState$1["static"] = 0] = "static";
	ParseRawPathSegmentState$1[ParseRawPathSegmentState$1["param"] = 1] = "param";
	ParseRawPathSegmentState$1[ParseRawPathSegmentState$1["regexp"] = 2] = "regexp";
	ParseRawPathSegmentState$1[ParseRawPathSegmentState$1["modifier"] = 3] = "modifier";
	return ParseRawPathSegmentState$1;
}(ParseRawPathSegmentState || {});
const IS_MODIFIER_RE = /[+*?]/;
/**
* Parses a raw path segment like the `:id` in a route `/users/:id`.
*
* @param segment - segment to parse without the extension
* @returns - the pathSegment and the params
*/
function parseRawPathSegment(segment) {
	let buffer = "";
	let state = ParseRawPathSegmentState.static;
	const params = [];
	const subSegments = [];
	let currentTreeRouteParam = createEmptyRouteParam();
	let pos = 0;
	let c;
	function consumeBuffer() {
		if (state === ParseRawPathSegmentState.static) subSegments.push(buffer);
		else if (state === ParseRawPathSegmentState.param || state === ParseRawPathSegmentState.regexp || state === ParseRawPathSegmentState.modifier) {
			if (!currentTreeRouteParam.paramName) {
				require_options.warn(`Invalid parameter in path "${segment}": parameter name cannot be empty. Using default name "pathMatch" for ':()'.`);
				currentTreeRouteParam.paramName = "pathMatch";
			}
			subSegments.push(currentTreeRouteParam);
			params.push(currentTreeRouteParam);
			currentTreeRouteParam = createEmptyRouteParam();
		}
		buffer = "";
	}
	for (pos = 0; pos < segment.length; pos++) {
		c = segment[pos];
		if (c === "\\") {
			pos++;
			buffer += segment[pos];
			continue;
		}
		if (state === ParseRawPathSegmentState.static) if (c === ":") {
			consumeBuffer();
			state = ParseRawPathSegmentState.param;
		} else buffer += c;
		else if (state === ParseRawPathSegmentState.param) if (c === "(") {
			currentTreeRouteParam.paramName = buffer;
			buffer = "";
			state = ParseRawPathSegmentState.regexp;
		} else if (IS_MODIFIER_RE.test(c)) {
			currentTreeRouteParam.modifier = c;
			currentTreeRouteParam.optional = c === "?" || c === "*";
			currentTreeRouteParam.repeatable = c === "+" || c === "*";
			consumeBuffer();
			state = ParseRawPathSegmentState.static;
		} else if (IS_VARIABLE_CHAR_RE.test(c)) {
			buffer += c;
			currentTreeRouteParam.paramName = buffer;
		} else {
			currentTreeRouteParam.paramName = buffer;
			consumeBuffer();
			pos--;
			state = ParseRawPathSegmentState.static;
		}
		else if (state === ParseRawPathSegmentState.regexp) if (c === ")") {
			if (buffer === ".*") currentTreeRouteParam.isSplat = true;
			state = ParseRawPathSegmentState.modifier;
		} else buffer += c;
		else if (state === ParseRawPathSegmentState.modifier) {
			if (IS_MODIFIER_RE.test(c)) {
				currentTreeRouteParam.modifier = c;
				currentTreeRouteParam.optional = c === "?" || c === "*";
				currentTreeRouteParam.repeatable = c === "+" || c === "*";
			} else pos--;
			consumeBuffer();
			state = ParseRawPathSegmentState.static;
		}
	}
	if (state === ParseRawPathSegmentState.regexp) throw new Error(`Invalid segment: "${segment}"`);
	if (buffer || state === ParseRawPathSegmentState.modifier) consumeBuffer();
	return [
		segment,
		params,
		subSegments
	];
}
/**
* Helper function to create an empty route param used by the parser.
*
* @returns an empty route param
*/
function createEmptyRouteParam() {
	return {
		paramName: "",
		parser: null,
		modifier: "",
		optional: false,
		repeatable: false,
		isSplat: false
	};
}

//#endregion
//#region src/core/tree.ts
var TreeNode = class TreeNode {
	/**
	* value of the node
	*/
	value;
	/**
	* children of the node
	*/
	children = /* @__PURE__ */ new Map();
	/**
	* Parent node.
	*/
	parent;
	/**
	* Plugin options taken into account by the tree.
	*/
	options;
	/**
	* Should this page import the page info
	*/
	hasDefinePage = false;
	/**
	* Creates a new tree node.
	*
	* @param options - TreeNodeOptions shared by all nodes
	* @param pathSegment - path segment of this node e.g. `users` or `:id`
	* @param parent
	*/
	constructor(options, pathSegment, parent) {
		this.options = options;
		this.parent = parent;
		this.value = createTreeNodeValue(pathSegment, parent?.value, options.treeNodeOptions || options.pathParser);
	}
	/**
	* Adds a path to the tree. `path` cannot start with a `/`.
	*
	* @param path - path segment to insert. **It shouldn't contain the file extension**
	* @param filePath - file path, must be a file (not a folder)
	*/
	insert(path$1, filePath) {
		const { tail, segment, viewName } = splitFilePath(path$1);
		if (!this.children.has(segment)) this.children.set(segment, new TreeNode(this.options, segment, this));
		const child = this.children.get(segment);
		if (!tail) child.value.components.set(viewName, filePath);
		else return child.insert(tail, filePath);
		return child;
	}
	/**
	* Adds a path that has already been parsed to the tree. `path` cannot start with a `/`. This method is similar to
	* `insert` but the path argument should be already parsed. e.g. `users/:id` for a file named `users/[id].vue`.
	*
	* @param path - path segment to insert, already parsed (e.g. users/:id)
	* @param filePath - file path, defaults to path for convenience and testing
	*/
	insertParsedPath(path$1, filePath = path$1) {
		const node = new TreeNode({
			...this.options,
			treeNodeOptions: {
				...this.options.pathParser,
				format: "path"
			}
		}, path$1, this);
		this.children.set(path$1, node);
		node.value.components.set("default", filePath);
		return node;
	}
	/**
	* Saves a custom route block for a specific file path. The file path is used as a key. Some special file paths will
	* have a lower or higher priority.
	*
	* @param filePath - file path where the custom block is located
	* @param routeBlock - custom block to set
	*/
	setCustomRouteBlock(filePath, routeBlock) {
		this.value.setOverride(filePath, routeBlock);
	}
	/**
	* Generator that yields all descendants without sorting.
	* Use with Array.from() for now, native .map() support in Node 22+.
	*/
	*getChildrenDeep() {
		for (const child of this.children.values()) {
			yield child;
			yield* child.getChildrenDeep();
		}
	}
	/**
	* Comparator function for sorting TreeNodes.
	*
	* @internal
	*/
	static compare(a, b) {
		return a.path.localeCompare(b.path, "en");
	}
	/**
	* Get the children of this node sorted by their path.
	*/
	getChildrenSorted() {
		return Array.from(this.children.values()).sort(TreeNode.compare);
	}
	/**
	* Calls {@link getChildrenDeep} and sorts the result by path in the end.
	*/
	getChildrenDeepSorted() {
		return Array.from(this.getChildrenDeep()).sort(TreeNode.compare);
	}
	/**
	* Delete and detach itself from the tree.
	*/
	delete() {
		if (!this.parent) throw new Error("Cannot delete the root node.");
		this.parent.children.delete(this.value.rawSegment);
		this.parent = void 0;
	}
	/**
	* Remove a route from the tree. The path shouldn't start with a `/` but it can be a nested one. e.g. `foo/bar`.
	* The `path` should be relative to the page folder.
	*
	* @param path - path segment of the file
	*/
	remove(path$1) {
		const { tail, segment, viewName } = splitFilePath(path$1);
		const child = this.children.get(segment);
		if (!child) throw new Error(`Cannot Delete "${path$1}". "${segment}" not found at "${this.path}".`);
		if (tail) {
			child.remove(tail);
			if (child.children.size === 0 && child.value.components.size === 0) this.children.delete(segment);
		} else {
			child.value.components.delete(viewName);
			if (child.children.size === 0 && child.value.components.size === 0) this.children.delete(segment);
		}
	}
	/**
	* Returns the route path of the node without parent paths. If the path was overridden, it returns the override.
	*/
	get path() {
		return this.value.overrides.path ?? (this.parent?.isRoot() ? "/" : "") + this.value.pathSegment;
	}
	/**
	* Returns the route path of the node including parent paths.
	*/
	get fullPath() {
		return this.value.fullPath;
	}
	/**
	* Object of components (filepaths) for this node.
	*/
	get components() {
		return Object.fromEntries(this.value.components.entries());
	}
	/**
	* Does this node render any component?
	*/
	get hasComponents() {
		return this.value.components.size > 0;
	}
	/**
	* Returns the route name of the node. If the name was overridden, it returns the override.
	*/
	get name() {
		const overrideName = this.value.overrides.name;
		return overrideName === void 0 ? this.options.getRouteName(this) : overrideName;
	}
	/**
	* Returns the meta property as an object.
	*/
	get metaAsObject() {
		return { ...this.value.overrides.meta };
	}
	/**
	* Returns the JSON string of the meta object of the node. If the meta was overridden, it returns the override. If
	* there is no override, it returns an empty string.
	*/
	get meta() {
		const overrideMeta = this.metaAsObject;
		return Object.keys(overrideMeta).length > 0 ? JSON.stringify(overrideMeta, null, 2) : "";
	}
	/**
	* Array of route params for this node. It includes **all** the params from the parents as well.
	*/
	get params() {
		const params = [...this.value.params];
		let node = this.parent;
		while (node) {
			params.unshift(...node.value.params);
			node = node.parent;
		}
		return params;
	}
	/**
	* Array of route params coming from the path. It includes all the params from the parents as well.
	*/
	get pathParams() {
		const params = this.value.isParam() ? [...this.value.pathParams] : [];
		let node = this.parent;
		while (node) {
			if (node.value.isParam()) params.unshift(...node.value.pathParams);
			node = node.parent;
		}
		return params;
	}
	/**
	* Array of query params extracted from definePage. Only returns query params from this specific node.
	*/
	get queryParams() {
		return this.value.queryParams;
	}
	/**
	* Generates a regexp based on this node and its parents. This regexp is used by the custom resolver
	*/
	get regexp() {
		let node = this;
		const nodeList = [];
		while (node && !node.isRoot()) {
			nodeList.unshift(node);
			node = node.parent;
		}
		let re = "";
		for (var i = 0; i < nodeList.length; i++) {
			node = nodeList[i];
			if (node.value.isParam()) {
				var nodeRe = node.value.re;
				if ((re || i < nodeList.length - 1) && node.value.subSegments.length === 1 && node.value.subSegments.at(0).optional) re += `(?:\\/${nodeRe.slice(0, -1)})?`;
				else re += (re ? "\\/" : "") + nodeRe;
			} else re += (re ? "\\/" : "") + escapeRegex(node.value.pathSegment);
		}
		return "/^" + (re.startsWith("(?:\\/") ? "" : "\\/") + re.replace(require_options.ESCAPED_TRAILING_SLASH_RE, "") + "$/i";
	}
	/**
	* Score of the path used for sorting routes.
	*/
	get score() {
		const scores = [];
		let node = this;
		while (node && !node.isRoot()) {
			scores.unshift(node.value.score);
			node = node.parent;
		}
		return scores;
	}
	/**
	* Is this node a splat (catch-all) param
	*/
	get isSplat() {
		return this.value.isParam() && this.value.pathParams.some((p) => p.isSplat);
	}
	/**
	* Returns an array of matcher parts that is consumed by
	* MatcherPatternPathDynamic to render the path.
	*/
	get matcherPatternPathDynamicParts() {
		const parts = [];
		let node = this;
		while (node && !node.isRoot()) {
			const subSegments = node.value.subSegments.map((segment) => typeof segment === "string" ? segment : segment.isSplat ? 0 : 1);
			if (subSegments.length > 1) parts.unshift(subSegments);
			else if (subSegments.length === 1) parts.unshift(subSegments[0]);
			node = node.parent;
		}
		return parts;
	}
	/**
	* Is this tree node matchable? A matchable node has at least one component
	* and a name.
	*/
	isMatchable() {
		return this.value.components.size > 0 && this.name !== false;
	}
	/**
	* Returns wether this tree node is the root node of the tree.
	*
	* @returns true if the node is the root node
	*/
	isRoot() {
		return !this.parent && this.value.fullPath === "/" && !this.value.components.size;
	}
	/**
	* Returns wether this tree node has a name. This allows to coerce the type
	* of TreeNode
	*/
	isNamed() {
		return !!this.name;
	}
	toString() {
		return `${this.isRoot() ? "·" : this.value}${this.value.components.size > 1 || this.value.components.size === 1 && !this.value.components.get("default") ? ` ⎈(${Array.from(this.value.components.keys()).join(", ")})` : ""}${this.hasDefinePage ? " ⚑ definePage()" : ""}`;
	}
};
/**
* Creates a new prefix tree. This is meant to only be the root node. It has access to extra methods that only make
* sense on the root node.
*/
var PrefixTree = class extends TreeNode {
	map = /* @__PURE__ */ new Map();
	constructor(options) {
		super(options, "");
	}
	insert(path$1, filePath) {
		const node = super.insert(path$1, filePath);
		this.map.set(filePath, node);
		return node;
	}
	/**
	* Returns the tree node of the given file path.
	*
	* @param filePath - file path of the tree node to get
	*/
	getChild(filePath) {
		return this.map.get(filePath);
	}
	/**
	* Removes the tree node of the given file path.
	*
	* @param filePath - file path of the tree node to remove
	*/
	removeChild(filePath) {
		if (this.map.has(filePath)) {
			this.map.get(filePath).delete();
			this.map.delete(filePath);
		}
	}
};
/**
* Splits a path into by finding the first '/' and returns the tail and segment. If it has an extension, it removes it.
* If it contains a named view, it returns the view name as well (otherwise it's default).
*
* @param filePath - filePath to split
*/
function splitFilePath(filePath) {
	const slashPos = filePath.indexOf("/");
	let head = slashPos < 0 ? filePath : filePath.slice(0, slashPos);
	const tail = slashPos < 0 ? "" : filePath.slice(slashPos + 1);
	let segment = head;
	let viewName = "default";
	const namedSeparatorPos = segment.indexOf("@");
	if (namedSeparatorPos > 0) {
		viewName = segment.slice(namedSeparatorPos + 1);
		segment = segment.slice(0, namedSeparatorPos);
	}
	return {
		segment,
		tail,
		viewName
	};
}

//#endregion
//#region src/codegen/generateParamParsers.ts
const NATIVE_PARAM_PARSERS = ["int", "bool"];
const NATIVE_PARAM_PARSERS_TYPES = {
	int: "number",
	bool: "boolean"
};
function warnMissingParamParsers(tree, paramParsers) {
	for (const node of tree.getChildrenDeepSorted()) for (const param of node.params) if (param.parser && !paramParsers.has(param.parser)) {
		if (!NATIVE_PARAM_PARSERS.includes(param.parser)) console.warn(`Parameter parser "${param.parser}" not found for route "${node.fullPath}".`);
	}
}
function generateParamParsersTypesDeclarations(paramParsers) {
	return Array.from(paramParsers.values()).map(({ typeName, relativePath }) => `type ${typeName} = ReturnType<NonNullable<typeof import('./${relativePath}').parser['get']>>`).join("\n");
}
function generateParamsTypes(params, parparsersMap) {
	return params.map((param) => {
		if (param.parser) {
			if (parparsersMap.has(param.parser)) return parparsersMap.get(param.parser).typeName;
			else if (param.parser in NATIVE_PARAM_PARSERS_TYPES) return NATIVE_PARAM_PARSERS_TYPES[param.parser];
		}
		return null;
	});
}
function generateParamParserOptions(param, importsMap, paramParsers) {
	if (!param.parser) return "";
	if (paramParsers.has(param.parser)) {
		const { name, absolutePath } = paramParsers.get(param.parser);
		const varName = `PARAM_PARSER__${name}`;
		importsMap.add(absolutePath, {
			name: "parser",
			as: varName
		});
		return varName;
	} else if (NATIVE_PARAM_PARSERS.includes(param.parser)) {
		const varName = `PARAM_PARSER_${param.parser.toUpperCase()}`;
		importsMap.add("vue-router/experimental", varName);
		return varName;
	}
	return "";
}
function generateParamParserCustomType(paramParsers) {
	const parserNames = Array.from(paramParsers.keys()).sort();
	if (parserNames.length === 0) return "never";
	if (parserNames.length === 1) return `'${parserNames[0]}'`;
	return parserNames.map((name) => `  | '${name}'`).join("\n");
}
function generatePathParamsOptions(params, importsMap, paramParsers) {
	const paramOptions = params.map((param) => {
		const optionList = [];
		const parser = generateParamParserOptions(param, importsMap, paramParsers);
		optionList.push(parser || `/* no parser */`);
		if (param.optional || param.repeatable) optionList.push(`/* repeatable: ` + (param.repeatable ? `*/ true` : `false */`));
		if (param.optional) optionList.push(`/* optional: ` + (param.optional ? `*/ true` : `false */`));
		return `
${param.paramName}: [${optionList.join(", ")}],
`.slice(1, -1);
	});
	return paramOptions.length === 0 ? "{}" : `{
      ${paramOptions.join("\n      ")}
    }`;
}

//#endregion
//#region src/codegen/generateRouteParams.ts
function generateRouteParams(node, isRaw) {
	const nodeParams = node.pathParams;
	return nodeParams.length > 0 ? `{ ${nodeParams.filter((param) => {
		if (!param.paramName) {
			console.warn(`Warning: A parameter without a name was found in the route "${node.fullPath}" in segment "${node.path}".\n‼️ This is a bug, please report it at https://github.com/posva/unplugin-vue-router`);
			return false;
		}
		return true;
	}).map((param) => `${param.paramName}${param.optional ? "?" : ""}: ` + (param.modifier === "+" ? `ParamValueOneOrMore<${isRaw}>` : param.modifier === "*" ? `ParamValueZeroOrMore<${isRaw}>` : param.modifier === "?" ? `ParamValueZeroOrOne<${isRaw}>` : `ParamValue<${isRaw}>`)).join(", ")} }` : "Record<never, never>";
}
function EXPERIMENTAL_generateRouteParams(node, types, isRaw) {
	const nodeParams = node.params;
	return nodeParams.length > 0 ? `{ ${nodeParams.map((param, i) => {
		const isOptional = isTreeParamOptional(param);
		const isRepeatable = isTreeParamRepeatable(param);
		const type = types[i];
		let extractedType;
		if (type?.startsWith("Param_")) extractedType = `${isRepeatable ? "Extract" : "Exclude"}<${type}, unknown[]>`;
		else extractedType = `${type ?? "string"}${isRepeatable ? "[]" : ""}`;
		extractedType += isTreePathParam(param) && isOptional && !isRepeatable ? " | null" : "";
		return `${param.paramName}${isRaw && isOptional ? "?" : ""}: ${extractedType}`;
	}).join(", ")} }` : "Record<never, never>";
}

//#endregion
//#region src/utils/index.ts
const ts = String.raw;
/**
* Pads a single-line string with spaces.
*
* @internal
*
* @param spaces The number of spaces to pad with.
* @param str The string to pad, none if omitted.
* @returns The padded string.
*/
function pad(spaces, str = "") {
	return " ".repeat(spaces) + str;
}
/**
* Formats an array of union items as a multiline union type.
*
* @internal
*
* @param items The items to format.
* @param spaces The number of spaces to indent each line.
* @returns The formatted multiline union type.
*/
function formatMultilineUnion(items, spaces) {
	return (items.length ? items : ["never"]).map((s) => `| ${s}`).join(`\n${pad(spaces)}`);
}
/**
* Converts a string value to a TS string literal type.
*
* @internal
*
* @param str the string to convert to a string type
* @returns The string wrapped in single quotes.
* @example
* stringToStringType('hello') // returns "'hello'"
*/
function stringToStringType(str) {
	return `'${str}'`;
}

//#endregion
//#region src/codegen/generateRouteMap.ts
function generateRouteNamedMap(node, options, paramParsersMap) {
	if (node.isRoot()) return `export interface RouteNamedMap {
${node.getChildrenSorted().map((n) => generateRouteNamedMap(n, options, paramParsersMap)).join("")}}`;
	return (node.value.components.size && node.isNamed() ? pad(2, `${stringToStringType(node.name)}: ${generateRouteRecordInfo(node, options, paramParsersMap)},\n`) : "") + (node.children.size > 0 ? node.getChildrenSorted().map((n) => generateRouteNamedMap(n, options, paramParsersMap)).join("\n") : "");
}
function generateRouteRecordInfo(node, options, paramParsersMap) {
	let paramParsers = [];
	if (options.experimental.paramParsers) paramParsers = generateParamsTypes(node.params, paramParsersMap);
	const typeParams = [
		stringToStringType(node.name),
		stringToStringType(node.fullPath),
		options.experimental.paramParsers ? EXPERIMENTAL_generateRouteParams(node, paramParsers, true) : generateRouteParams(node, true),
		options.experimental.paramParsers ? EXPERIMENTAL_generateRouteParams(node, paramParsers, false) : generateRouteParams(node, false)
	];
	const childRouteNames = node.children.size > 0 ? Array.from(node.getChildrenDeep()).reduce((acc, childRoute) => {
		if (childRoute.value.components.size && childRoute.isNamed()) acc.push(childRoute.name);
		return acc;
	}, []).sort() : [];
	typeParams.push(formatMultilineUnion(childRouteNames.map(stringToStringType), 4));
	return `RouteRecordInfo<
${typeParams.map((line) => pad(4, line)).join(",\n")}
  >`;
}

//#endregion
//#region src/codegen/generateRouteFileInfoMap.ts
function generateRouteFileInfoMap(node, { root }) {
	if (!node.isRoot()) throw new Error("The provided node is not a root node");
	const routesInfoList = node.getChildrenSorted().flatMap((child) => generateRouteFileInfoLines(child, root));
	const routesInfo = /* @__PURE__ */ new Map();
	for (const routeInfo of routesInfoList) {
		let info = routesInfo.get(routeInfo.key);
		if (!info) routesInfo.set(routeInfo.key, info = {
			routes: [],
			views: []
		});
		info.routes.push(...routeInfo.routeNames);
		info.views.push(...routeInfo.childrenNamedViews || []);
	}
	return `export interface _RouteFileInfoMap {
${Array.from(routesInfo.entries()).map(([file, { routes, views }]) => `
  '${file}': {
    routes:
      ${formatMultilineUnion(routes.sort().map(stringToStringType), 6)}
    views:
      ${formatMultilineUnion(views.sort().map(stringToStringType), 6)}
  }`).join("\n")}
}`;
}
/**
* Generate the route file info for a non-root node.
*/
function generateRouteFileInfoLines(node, rootDir) {
	const deepChildren = node.children.size > 0 ? node.getChildrenDeepSorted() : null;
	const deepChildrenNamedViews = deepChildren ? Array.from(new Set(deepChildren.flatMap((child) => Array.from(child.value.components.keys())))) : null;
	const routeNames = [node, ...deepChildren ?? []].reduce((acc, node$1) => {
		if (node$1.isNamed() && node$1.value.components.size > 0) acc.push(node$1.name);
		return acc;
	}, []);
	const currentRouteInfo = routeNames.length === 0 ? [] : Array.from(node.value.components.values()).map((file) => ({
		key: (0, pathe.relative)(rootDir, file).replaceAll("\\", "/"),
		routeNames,
		childrenNamedViews: deepChildrenNamedViews
	}));
	const childrenRouteInfo = node.getChildrenSorted().flatMap((child) => generateRouteFileInfoLines(child, rootDir));
	return currentRouteInfo.concat(childrenRouteInfo);
}

//#endregion
//#region src/core/moduleConstants.ts
const MODULE_ROUTES_PATH = `vue-router/auto-routes`;
const MODULE_RESOLVER_PATH = `vue-router/auto-resolver`;
let time = Date.now();
/**
* Last time the routes were loaded from MODULE_ROUTES_PATH
*/
const ROUTES_LAST_LOAD_TIME = {
	get value() {
		return time;
	},
	update(when = Date.now()) {
		time = when;
	}
};
const VIRTUAL_PREFIX = "\0";
const ROUTE_BLOCK_ID = asVirtualId("vue-router/auto/route-block");
function getVirtualId(id) {
	return id.startsWith(VIRTUAL_PREFIX) ? id.slice(1) : null;
}
const routeBlockQueryRE = /\?vue&type=route/;
function asVirtualId(id) {
	return VIRTUAL_PREFIX + id;
}
const DEFINE_PAGE_QUERY_RE = /\?.*\bdefinePage\&vue\b/;

//#endregion
//#region src/codegen/generateRouteRecords.ts
/**
* Generate the route records for the given node.
*
* @param node - the node to generate the route record for
* @param options - the options to use
* @param importsMap - the imports map to fill and use
* @param indent - the indent level
* @returns the code of the routes as a string
*/
function generateRouteRecords(node, options, importsMap, indent = 0) {
	if (node.isRoot()) return `[
${node.getChildrenSorted().map((child) => generateRouteRecords(child, options, importsMap, indent + 1)).join(",\n")}
]`;
	const definePageDataList = [];
	if (node.hasDefinePage) {
		for (const [name, filePath] of node.value.components) {
			const pageDataImport = `_definePage_${name}_${importsMap.size}`;
			definePageDataList.push(pageDataImport);
			const lang = (0, __vue_macros_common.getLang)(filePath);
			importsMap.addDefault(`${filePath}?definePage&` + (lang === "vue" ? "vue&lang.tsx" : `lang.${lang}`), pageDataImport);
		}
		if (definePageDataList.length > 0) indent++;
	}
	const startIndent = pad(indent * 2);
	const indentStr = pad((indent + 1) * 2);
	const overrides = node.value.overrides;
	const routeRecord = `${startIndent}{
${indentStr}path: '${node.path}',
${indentStr}${node.value.components.size ? node.isNamed() ? `name: ${stringToStringType(node.name)},` : `/* no name */` : `/* internal name: ${typeof node.name === "string" ? stringToStringType(node.name) : node.name} */`}
${indentStr}${node.value.components.size ? generateRouteRecordComponent$1(node, indentStr, options.importMode, importsMap) : "/* no component */"}
${overrides.props != null ? indentStr + `props: ${overrides.props},\n` : ""}${overrides.alias != null ? indentStr + `alias: ${JSON.stringify(overrides.alias)},\n` : ""}${indentStr}${node.children.size > 0 ? `children: [
${node.getChildrenSorted().map((child) => generateRouteRecords(child, options, importsMap, indent + 2)).join(",\n")}
${indentStr}],` : "/* no children */"}${formatMeta(node, indentStr)}
${startIndent}}`;
	if (definePageDataList.length > 0) {
		const mergeCallIndent = startIndent.slice(2);
		importsMap.add("unplugin-vue-router/runtime", "_mergeRouteRecord");
		return `${mergeCallIndent}_mergeRouteRecord(
${routeRecord},
${definePageDataList.map((s) => startIndent + s).join(",\n")}
${mergeCallIndent})`;
	}
	return routeRecord;
}
function generateRouteRecordComponent$1(node, indentStr, importMode, importsMap) {
	const files = Array.from(node.value.components);
	return files.length === 1 && files[0][0] === "default" ? `component: ${generatePageImport(files[0][1], importMode, importsMap)},` : `components: {
${files.map(([key, path$1]) => `${indentStr + "  "}'${key}': ${generatePageImport(path$1, importMode, importsMap)}`).join(",\n")}
${indentStr}},`;
}
/**
* Generate the import (dynamic or static) for the given filepath. If the filepath is a static import, add it to the importsMap.
*
* @param filepath - the filepath to the file
* @param importMode - the import mode to use
* @param importsMap - the import list to fill
* @returns
*/
function generatePageImport(filepath, importMode, importsMap) {
	if ((typeof importMode === "function" ? importMode(filepath) : importMode) === "async") return `() => import('${filepath}')`;
	const existingEntry = importsMap.getImportList(filepath).find((entry) => entry.name === "default");
	if (existingEntry) return existingEntry.as;
	const importName = `_page_${importsMap.size}`;
	importsMap.addDefault(filepath, importName);
	return importName;
}
function formatMeta(node, indent) {
	const meta = node.meta;
	const formatted = meta && meta.split("\n").map((line) => indent + line).join("\n") + ",";
	return formatted ? "\n" + indent + "meta: " + formatted.trimStart() : "";
}

//#endregion
//#region src/core/customBlock.ts
function getRouteBlock(path$1, content, options) {
	const blockStr = (0, __vue_compiler_sfc.parse)(content, { pad: "space" }).descriptor?.customBlocks.find((b) => b.type === "route");
	if (blockStr) return parseCustomBlock(blockStr, path$1, options);
}
function parseCustomBlock(block, filePath, options) {
	const lang = block.lang ?? options.routeBlockLang;
	if (lang === "json5") try {
		return json5.default.parse(block.content);
	} catch (err) {
		require_options.warn(`Invalid JSON5 format of <${block.type}> content in ${filePath}\n${err.message}`);
	}
	else if (lang === "json") try {
		return JSON.parse(block.content);
	} catch (err) {
		require_options.warn(`Invalid JSON format of <${block.type}> content in ${filePath}\n${err.message}`);
	}
	else if (lang === "yaml" || lang === "yml") try {
		return (0, yaml.parse)(block.content);
	} catch (err) {
		require_options.warn(`Invalid YAML format of <${block.type}> content in ${filePath}\n${err.message}`);
	}
	else require_options.warn(`Language "${lang}" for <${block.type}> is not supported. Supported languages are: json5, json, yaml, yml. Found in in ${filePath}.`);
}

//#endregion
//#region src/core/RoutesFolderWatcher.ts
var RoutesFolderWatcher = class {
	src;
	path;
	extensions;
	filePatterns;
	exclude;
	watcher;
	constructor(folderOptions) {
		this.src = folderOptions.src;
		this.path = folderOptions.path;
		this.exclude = folderOptions.exclude;
		this.extensions = folderOptions.extensions;
		this.filePatterns = folderOptions.pattern;
		const isMatch = (0, picomatch.default)(this.filePatterns, { ignore: this.exclude });
		this.watcher = (0, chokidar.watch)(".", {
			cwd: this.src,
			ignoreInitial: true,
			ignorePermissionErrors: true,
			awaitWriteFinish: !!process.env.CI,
			ignored: (filePath, stats) => {
				if (!stats || stats.isDirectory()) return false;
				return !isMatch(pathe.default.relative(this.src, filePath));
			}
		});
	}
	on(event, handler) {
		this.watcher.on(event, (filePath) => {
			filePath = (0, pathe.resolve)(this.src, filePath);
			handler({
				filePath,
				routePath: require_options.asRoutePath({
					src: this.src,
					path: this.path,
					extensions: this.extensions
				}, filePath)
			});
		});
		return this;
	}
	close() {
		return this.watcher.close();
	}
};
function resolveFolderOptions(globalOptions, folderOptions) {
	const extensions = overrideOption(globalOptions.extensions, folderOptions.extensions);
	const filePatterns = overrideOption(globalOptions.filePatterns, folderOptions.filePatterns);
	return {
		src: pathe.default.resolve(globalOptions.root, folderOptions.src),
		pattern: require_options.appendExtensionListToPattern(filePatterns, extensions),
		path: folderOptions.path || "",
		extensions,
		filePatterns,
		exclude: overrideOption(globalOptions.exclude, folderOptions.exclude).map((p) => p.startsWith("**") ? p : (0, pathe.resolve)(p))
	};
}
function overrideOption(existing, newValue) {
	const asArray = typeof existing === "string" ? [existing] : existing;
	if (typeof newValue === "function") return newValue(asArray);
	if (typeof newValue !== "undefined") return typeof newValue === "string" ? [newValue] : newValue;
	return asArray;
}

//#endregion
//#region src/codegen/generateDTS.ts
/**
* Removes empty lines and indent by two spaces to match the rest of the file.
*/
function normalizeLines(code) {
	return code.split("\n").filter((line) => line.length !== 0).map((line) => pad(2, line)).join("\n");
}
function generateDTS({ routesModule, routeNamedMap, routeFileInfoMap, paramsTypesDeclaration, customParamsType }) {
	return ts`
/* eslint-disable */
/* prettier-ignore */
// @ts-nocheck
// noinspection ES6UnusedImports
// Generated by unplugin-vue-router. !! DO NOT MODIFY THIS FILE !!
// It's recommended to commit this file.
// Make sure to add this file to your tsconfig.json file as an "includes" or "files" entry.

${paramsTypesDeclaration ? `
// Custom route params parsers
${paramsTypesDeclaration}

`.trimStart() : ""}declare module 'vue-router/auto-resolver' {
  export type ParamParserCustom = ${customParamsType}
}

declare module '${routesModule}' {
  import type {
    RouteRecordInfo,
    ParamValue,
    ParamValueOneOrMore,
    ParamValueZeroOrMore,
    ParamValueZeroOrOne,
  } from 'vue-router'

  /**
   * Route name map generated by unplugin-vue-router
   */
${normalizeLines(routeNamedMap)}

  /**
   * Route file to route info map by unplugin-vue-router.
   * Used by the \`sfc-typed-router\` Volar plugin to automatically type \`useRoute()\`.
   *
   * Each key is a file path relative to the project root with 2 properties:
   * - routes: union of route names of the possible routes when in this page (passed to useRoute<...>())
   * - views: names of nested views (can be passed to <RouterView name="...">)
   *
   * @internal
   */
${normalizeLines(routeFileInfoMap)}

  /**
   * Get a union of possible route names in a certain route component file.
   * Used by the \`sfc-typed-router\` Volar plugin to automatically type \`useRoute()\`.
   *
   * @internal
   */
  export type _RouteNamesForFilePath<FilePath extends string> =
    _RouteFileInfoMap extends Record<FilePath, infer Info>
      ? Info['routes']
      : keyof RouteNamedMap
}
`.trimStart();
}

//#endregion
//#region src/core/definePage.ts
const MACRO_DEFINE_PAGE = "definePage";
const MACRO_DEFINE_PAGE_QUERY = /[?&]definePage\b/;
/**
* Generate the ast from a code string and an id. Works with SFC and non-SFC files.
*/
function getCodeAst(code, id) {
	let offset = 0;
	let ast;
	const lang = (0, __vue_macros_common.getLang)(id.split(MACRO_DEFINE_PAGE_QUERY)[0]);
	if (lang === "vue") {
		const sfc = (0, __vue_macros_common.parseSFC)(code, id);
		if (sfc.scriptSetup) {
			ast = sfc.getSetupAst();
			offset = sfc.scriptSetup.loc.start.offset;
		} else if (sfc.script) {
			ast = sfc.getScriptAst();
			offset = sfc.script.loc.start.offset;
		}
	} else if (/[jt]sx?$/.test(lang)) ast = (0, __vue_macros_common.babelParse)(code, lang);
	const definePageNodes = (ast?.body || []).map((node) => {
		const definePageCallNode = node.type === "ExpressionStatement" ? node.expression : node;
		return (0, __vue_macros_common.isCallOf)(definePageCallNode, MACRO_DEFINE_PAGE) ? definePageCallNode : null;
	}).filter((node) => !!node);
	return {
		ast,
		offset,
		definePageNodes
	};
}
function definePageTransform({ code, id }) {
	const isExtractingDefinePage = MACRO_DEFINE_PAGE_QUERY.test(id);
	if (!code.includes(MACRO_DEFINE_PAGE)) return isExtractingDefinePage ? "export default {}" : void 0;
	let ast;
	let offset;
	let definePageNodes;
	try {
		const result = getCodeAst(code, id);
		ast = result.ast;
		offset = result.offset;
		definePageNodes = result.definePageNodes;
	} catch (error) {
		require_options.warn(`[${id}]: Failed to process definePage: ${error instanceof Error ? error.message : "Unknown error"}`);
		return isExtractingDefinePage ? "export default {}" : void 0;
	}
	if (!ast) return;
	if (!definePageNodes.length) return isExtractingDefinePage ? "export default {}" : null;
	else if (definePageNodes.length > 1) throw new SyntaxError(`duplicate definePage() call`);
	const definePageNode = definePageNodes[0];
	if (isExtractingDefinePage) {
		const s = new __vue_macros_common.MagicString(code);
		const routeRecord = definePageNode.arguments[0];
		if (!routeRecord) throw new SyntaxError(`[${id}]: definePage() expects an object expression as its only argument`);
		const scriptBindings = ast.body ? getIdentifiers(ast.body) : [];
		try {
			(0, __vue_macros_common.checkInvalidScopeReference)(routeRecord, MACRO_DEFINE_PAGE, scriptBindings);
		} catch (error) {
			require_options.warn(`[${id}]: ${error instanceof Error ? error.message : "Invalid scope reference in definePage"}`);
			return "export default {}";
		}
		s.remove(offset + routeRecord.end, code.length);
		s.remove(0, offset + routeRecord.start);
		s.prepend(`export default `);
		const staticImports = (0, mlly.findStaticImports)(code);
		const usedIds = /* @__PURE__ */ new Set();
		const localIds = /* @__PURE__ */ new Set();
		(0, ast_walker_scope.walkAST)(routeRecord, {
			enter(node) {
				if (this.parent?.type === "ObjectProperty" && this.parent.key === node && !this.parent.computed && node.type === "Identifier") this.skip();
				else if (this.parent?.type === "MemberExpression" && this.parent.property === node && !this.parent.computed && node.type === "Identifier") this.skip();
				else if (node.type === "TSTypeAnnotation") this.skip();
				else if (node.type === "Identifier" && !localIds.has(node.name)) usedIds.add(node.name);
				else if ("scopeIds" in node && node.scopeIds instanceof Set) for (const id$1 of node.scopeIds) localIds.add(id$1);
			},
			leave(node) {
				if ("scopeIds" in node && node.scopeIds instanceof Set) for (const id$1 of node.scopeIds) localIds.delete(id$1);
			}
		});
		for (const imp of staticImports) {
			const importCode = generateFilteredImportStatement((0, mlly.parseStaticImport)(imp), usedIds);
			if (importCode) s.prepend(importCode + "\n");
		}
		return (0, __vue_macros_common.generateTransform)(s, id);
	} else {
		const s = new __vue_macros_common.MagicString(code);
		s.remove(offset + definePageNode.start, offset + definePageNode.end);
		return (0, __vue_macros_common.generateTransform)(s, id);
	}
}
/**
* Extracts name, path, and params from definePage(). Those do not require
* extracting the whole definePage object as a different import
*/
function extractDefinePageInfo(sfcCode, id) {
	if (!sfcCode.includes(MACRO_DEFINE_PAGE)) return;
	let ast;
	let definePageNodes;
	try {
		const result = getCodeAst(sfcCode, id);
		ast = result.ast;
		definePageNodes = result.definePageNodes;
	} catch (error) {
		require_options.warn(`[${id}]: Failed to extract definePage info: ${error instanceof Error ? error.message : "Unknown error"}`);
		return;
	}
	if (!ast) return;
	if (!definePageNodes.length) return;
	else if (definePageNodes.length > 1) throw new SyntaxError(`duplicate definePage() call`);
	const routeRecord = definePageNodes[0].arguments[0];
	if (!routeRecord) throw new SyntaxError(`[${id}]: definePage() expects an object expression as its only argument`);
	if (routeRecord.type !== "ObjectExpression") throw new SyntaxError(`[${id}]: definePage() expects an object expression as its only argument`);
	const routeInfo = {};
	for (const prop of routeRecord.properties) if (prop.type === "ObjectProperty" && prop.key.type === "Identifier") {
		if (prop.key.name === "name") if (prop.value.type !== "StringLiteral" && (prop.value.type !== "BooleanLiteral" || prop.value.value !== false)) require_options.warn(`route name must be a string literal or false. Found in "${id}".`);
		else routeInfo.name = prop.value.value;
		else if (prop.key.name === "path") if (prop.value.type !== "StringLiteral") require_options.warn(`route path must be a string literal. Found in "${id}".`);
		else routeInfo.path = prop.value.value;
		else if (prop.key.name === "params") {
			if (prop.value.type === "ObjectExpression") routeInfo.params = extractParamsInfo(prop.value, id);
		}
	}
	return routeInfo;
}
function extractParamsInfo(paramsObj, id) {
	const params = {};
	for (const prop of paramsObj.properties) if (prop.type === "ObjectProperty" && prop.key.type === "Identifier") {
		if (prop.key.name === "query" && prop.value.type === "ObjectExpression") params.query = extractQueryParams(prop.value, id);
		else if (prop.key.name === "path" && prop.value.type === "ObjectExpression") params.path = extractPathParams(prop.value, id);
	}
	return params;
}
function extractQueryParams(queryObj, _id) {
	const queryParams = {};
	for (const prop of queryObj.properties) if (prop.type === "ObjectProperty" && prop.key.type === "Identifier") {
		const paramName = prop.key.name;
		if (prop.value.type === "StringLiteral") queryParams[paramName] = { parser: prop.value.value };
		else if (prop.value.type === "ObjectExpression") {
			const paramInfo = {};
			for (const paramProp of prop.value.properties) if (paramProp.type === "ObjectProperty" && paramProp.key.type === "Identifier") {
				if (paramProp.key.name === "parser" && paramProp.value.type === "StringLiteral") paramInfo.parser = paramProp.value.value;
				else if (paramProp.key.name === "format" && paramProp.value.type === "StringLiteral") paramInfo.format = paramProp.value.value;
				else if (paramProp.key.name === "default") if (typeof paramProp.value.extra?.raw === "string") paramInfo.default = paramProp.value.extra.raw;
				else if (paramProp.value.type === "NumericLiteral") paramInfo.default = String(paramProp.value.value);
				else if (paramProp.value.type === "StringLiteral") paramInfo.default = JSON.stringify(paramProp.value.value);
				else if (paramProp.value.type === "BooleanLiteral") paramInfo.default = String(paramProp.value.value);
				else if (paramProp.value.type === "NullLiteral") paramInfo.default = "null";
				else if (paramProp.value.type === "UnaryExpression" && (paramProp.value.operator === "-" || paramProp.value.operator === "+" || paramProp.value.operator === "!" || paramProp.value.operator === "~") && paramProp.value.argument.type === "NumericLiteral") paramInfo.default = `${paramProp.value.operator}${paramProp.value.argument.value}`;
				else if (paramProp.value.type === "ArrowFunctionExpression") paramInfo.default = (0, __babel_generator.generate)(paramProp.value).code;
				else require_options.warn(`Unrecognized default value in definePage() for query param "${paramName}". Typeof value: "${paramProp.value.type}". This is a bug or a missing type of value, open an issue on https://github.com/posva/unplugin-vue-router and provide the definePage() code.`);
			}
			queryParams[paramName] = paramInfo;
		}
	}
	return queryParams;
}
function extractPathParams(pathObj, _id) {
	const pathParams = {};
	for (const prop of pathObj.properties) if (prop.type === "ObjectProperty" && prop.key.type === "Identifier" && prop.value.type === "StringLiteral") pathParams[prop.key.name] = prop.value.value;
	return pathParams;
}
const getIdentifiers = (stmts) => {
	let ids = [];
	(0, ast_walker_scope.walkAST)({
		type: "Program",
		body: stmts,
		directives: [],
		sourceType: "module"
	}, {
		enter(node) {
			if (node.type === "BlockStatement") this.skip();
		},
		leave(node) {
			if (node.type !== "Program") return;
			ids = Object.keys(this.scope);
		}
	});
	return ids;
};
/**
* Generate a filtere import statement based on a set of identifiers that should be kept.
*
* @param parsedImports - parsed imports with mlly
* @param usedIds - set of used identifiers
* @returns `null` if no import statement should be generated, otherwise the import statement as a string without a newline
*/
function generateFilteredImportStatement(parsedImports, usedIds) {
	if (!parsedImports || usedIds.size < 1) return null;
	const { namedImports, defaultImport, namespacedImport } = parsedImports;
	if (namespacedImport && usedIds.has(namespacedImport)) return `import * as ${namespacedImport} from '${parsedImports.specifier}'`;
	let importListCode = "";
	if (defaultImport && usedIds.has(defaultImport)) importListCode += defaultImport;
	let namedImportListCode = "";
	for (const importName in namedImports) if (usedIds.has(importName)) {
		namedImportListCode += namedImportListCode ? `, ` : "";
		namedImportListCode += importName === namedImports[importName] ? importName : `${importName} as ${namedImports[importName]}`;
	}
	importListCode += importListCode && namedImportListCode ? ", " : "";
	importListCode += namedImportListCode ? `{${namedImportListCode}}` : "";
	if (!importListCode) return null;
	return `import ${importListCode} from '${parsedImports.specifier}'`;
}

//#endregion
//#region src/core/extendRoutes.ts
/**
* A route node that can be modified by the user. The tree can be iterated to be traversed.
* @example
* ```js
* [...node] // creates an array of all the children
* for (const child of node) {
*   // do something with the child node
* }
* ```
*
* @experimental
*/
var EditableTreeNode = class EditableTreeNode {
	node;
	constructor(node) {
		this.node = node;
	}
	/**
	* Remove and detach the current route node from the tree. Subsequently, its children will be removed as well.
	*/
	delete() {
		return this.node.delete();
	}
	/**
	* Inserts a new route as a child of this route. This route cannot use `definePage()`. If it was meant to be included,
	* add it to the `routesFolder` option.
	*
	* @param path - path segment to insert. Note this is relative to the current route. **It shouldn't start with `/`**. If it does, it will be added to the root of the tree.
	* @param filePath - file path
	* @returns the new editable route node
	*/
	insert(path$1, filePath) {
		let addBackLeadingSlash = false;
		if (path$1.startsWith("/")) {
			path$1 = path$1.slice(1);
			addBackLeadingSlash = !this.node.isRoot();
		}
		const node = this.node.insertParsedPath(path$1, filePath);
		const editable = new EditableTreeNode(node);
		if (addBackLeadingSlash) editable.path = "/" + node.path;
		return editable;
	}
	/**
	* Get an editable version of the parent node if it exists.
	*/
	get parent() {
		return this.node.parent && new EditableTreeNode(this.node.parent);
	}
	/**
	* Return a Map of the files associated to the current route. The key of the map represents the name of the view (Vue
	* Router feature) while the value is the **resolved** file path.
	* By default, the name of the view is `default`.
	*/
	get components() {
		return this.node.value.components;
	}
	/**
	* Alias for `route.components.get('default')`.
	*/
	get component() {
		return this.node.value.components.get("default");
	}
	/**
	* Name of the route. Note that **all routes are named** but when the final `routes` array is generated, routes
	* without a `component` will not include their `name` property to avoid accidentally navigating to them and display
	* nothing.
	* @see {@link isPassThrough}
	*/
	get name() {
		return this.node.name;
	}
	/**
	* Override the name of the route.
	*/
	set name(name) {
		this.node.value.addEditOverride({ name });
	}
	/**
	* Whether the route is a pass-through route. A pass-through route is a route that does not have a component and is
	* used to group other routes under the same prefix `path` and/or `meta` properties.
	*/
	get isPassThrough() {
		return this.node.value.components.size === 0;
	}
	/**
	* Meta property of the route as an object. Note this property is readonly and will be serialized as JSON. It won't contain the meta properties defined with `definePage()` as it could contain expressions **but it does contain the meta properties defined with `<route>` blocks**.
	*/
	get meta() {
		return this.node.metaAsObject;
	}
	/**
	* Override the meta property of the route. This will discard any other meta property defined with `<route>` blocks or
	* through other means. If you want to keep the existing meta properties, use `addToMeta`.
	* @see {@link addToMeta}
	*/
	set meta(meta) {
		this.node.value.removeOverride("meta");
		this.node.value.setEditOverride("meta", meta);
	}
	/**
	* Add meta properties to the route keeping the existing ones. The passed object will be deeply merged with the
	* existing meta object if any. Note that the meta property is later on serialized as JSON so you can't pass functions
	* or any other non-serializable value.
	*/
	addToMeta(meta) {
		this.node.value.addEditOverride({ meta });
	}
	/**
	* Path of the route without parent paths.
	*/
	get path() {
		return this.node.path;
	}
	/**
	* Override the path of the route. You must ensure `params` match with the existing path.
	*/
	set path(path$1) {
		if ((!this.node.parent || this.node.parent.isRoot()) && !path$1.startsWith("/")) path$1 = "/" + path$1;
		this.node.value.addEditOverride({ path: path$1 });
	}
	/**
	* Alias of the route.
	*/
	get alias() {
		return this.node.value.overrides.alias;
	}
	/**
	* Add an alias to the route.
	*
	* @param alias - Alias to add to the route
	*/
	addAlias(alias) {
		this.node.value.addEditOverride({ alias });
	}
	/**
	* Array of the route params and all of its parent's params. Note that
	* changing the params will not update the path, you need to update both.
	*/
	get params() {
		return this.node.pathParams;
	}
	/**
	* Path of the route including parent paths.
	*/
	get fullPath() {
		return this.node.fullPath;
	}
	/**
	* Computes an array of EditableTreeNode from the current node. Differently from iterating over the tree, this method
	* **only returns direct children**.
	*/
	get children() {
		return [...this.node.children.values()].map((node) => new EditableTreeNode(node));
	}
	/**
	* DFS traversal of the tree.
	* @example
	* ```ts
	* for (const node of tree) {
	*   // ...
	* }
	* ```
	*/
	*traverseDFS() {
		if (!this.node.isRoot()) yield this;
		for (const [_name, child] of this.node.children) yield* new EditableTreeNode(child).traverseDFS();
	}
	*[Symbol.iterator]() {
		yield* this.traverseBFS();
	}
	/**
	* BFS traversal of the tree as a generator.
	*
	* @example
	* ```ts
	* for (const node of tree) {
	*   // ...
	* }
	* ```
	*/
	*traverseBFS() {
		for (const [_name, child] of this.node.children) yield new EditableTreeNode(child);
		for (const [_name, child] of this.node.children) yield* new EditableTreeNode(child).traverseBFS();
	}
};

//#endregion
//#region src/codegen/generateRouteResolver.ts
/**
* Compare two score arrays for sorting routes by priority.
* Higher scores should come first (more specific routes).
*/
function compareRouteScore(a, b) {
	const maxLength = Math.max(a.length, b.length);
	for (let i = 0; i < maxLength; i++) {
		const aSegment = a[i] || [];
		const bSegment = b[i] || [];
		const aMinScore = aSegment.length > 0 ? Math.min(...aSegment) : 0;
		const bMinScore = bSegment.length > 0 ? Math.min(...bSegment) : 0;
		if (aMinScore !== bMinScore) return bMinScore - aMinScore;
		const aAvgScore = aSegment.length > 0 ? aSegment.reduce((sum, s) => sum + s, 0) / aSegment.length : 0;
		const bAvgScore = bSegment.length > 0 ? bSegment.reduce((sum, s) => sum + s, 0) / bSegment.length : 0;
		if (aAvgScore !== bAvgScore) return bAvgScore - aAvgScore;
		if (aSegment.length !== bSegment.length) return aSegment.length - bSegment.length;
	}
	return a.length - b.length;
}
const ROUTE_RECORD_VAR_PREFIX = "__route_";
function generateRouteResolver(tree, options, importsMap, paramParsersMap) {
	const state = {
		id: 0,
		matchableRecords: []
	};
	const records = tree.getChildrenSorted().map((node) => generateRouteRecord({
		node,
		parentVar: null,
		parentNode: null,
		state,
		options,
		importsMap,
		paramParsersMap
	}));
	importsMap.add("vue-router/experimental", "createFixedResolver");
	importsMap.add("vue-router/experimental", "MatcherPatternPathStatic");
	importsMap.add("vue-router/experimental", "MatcherPatternPathDynamic");
	importsMap.add("vue-router/experimental", "normalizeRouteRecord");
	return ts`
${records.join("\n\n")}

export const resolver = createFixedResolver([
${state.matchableRecords.sort((a, b) => compareRouteScore(a.score, b.score)).map(({ varName, path: path$1 }) => `  ${varName},  ${" ".repeat(String(state.id).length - varName.length + 8)}// ${path$1}`).join("\n")}
])
`;
}
/**
* Generates the route record in the format expected by the static resolver.
*/
function generateRouteRecord({ node, parentVar, parentNode, state, options, importsMap, paramParsersMap }) {
	const isMatchable = node.isMatchable();
	const shouldSkipNode = !isMatchable && !node.meta && !node.hasComponents;
	let varName = null;
	let recordDeclaration = "";
	const definePageDataList = [];
	if (node.hasDefinePage) for (const [name, filePath] of node.value.components) {
		const pageDataImport = `_definePage_${name}_${importsMap.size}`;
		definePageDataList.push(pageDataImport);
		const lang = (0, __vue_macros_common.getLang)(filePath);
		importsMap.addDefault(`${filePath}?definePage&` + (lang === "vue" ? "vue&lang.tsx" : `lang.${lang}`), pageDataImport);
	}
	if (!shouldSkipNode) {
		varName = `${ROUTE_RECORD_VAR_PREFIX}${state.id++}`;
		let recordName;
		const recordComponents = generateRouteRecordComponent(node, "  ", options.importMode, importsMap);
		if (isMatchable) {
			state.matchableRecords.push({
				path: node.fullPath,
				varName,
				score: node.score
			});
			recordName = `name: '${node.name}',`;
		} else recordName = node.name ? `/* (internal) name: '${node.name}' */` : `/* (removed) name: false */`;
		const queryProperty = generateRouteRecordQuery({
			node,
			importsMap,
			paramParsersMap
		});
		const routeRecordObject = `{
  ${recordName}
  ${generateRouteRecordPath({
			node,
			importsMap,
			paramParsersMap,
			parentVar,
			parentNode
		})}${queryProperty ? `\n  ${queryProperty}` : ""}${formatMeta(node, "  ")}
  ${recordComponents}${parentVar ? `\n  parent: ${parentVar},` : ""}
}`;
		recordDeclaration = definePageDataList.length > 0 ? `
const ${varName} = normalizeRouteRecord(
  ${generateRouteRecordMerge(routeRecordObject, definePageDataList, importsMap)}
)
` : `
const ${varName} = normalizeRouteRecord(${routeRecordObject})
`.trim().split("\n").filter((l) => l.trimStart().length > 0).join("\n");
	}
	const children = node.getChildrenSorted().map((child) => generateRouteRecord({
		node: child,
		parentVar: shouldSkipNode ? parentVar : varName,
		parentNode: shouldSkipNode ? parentNode : node,
		state,
		options,
		importsMap,
		paramParsersMap
	}));
	return recordDeclaration + (children.length ? (recordDeclaration ? "\n" : "") + children.join("\n") : "");
}
function generateRouteRecordComponent(node, indentStr, importMode, importsMap) {
	if (!node.hasComponents) return "";
	return `components: {
${Array.from(node.value.components).map(([key, path$1]) => `${indentStr + "  "}'${key}': ${generatePageImport(path$1, importMode, importsMap)}`).join(",\n")}
${indentStr}},`;
}
/**
* Generates the `path` property of a route record for the static resolver.
*/
function generateRouteRecordPath({ node, importsMap, paramParsersMap, parentVar, parentNode }) {
	if (!node.isMatchable() && node.name) return "";
	if (parentVar && parentNode && node.regexp === parentNode.regexp) return `path: ${parentVar}.path,`;
	const params = node.pathParams;
	if (params.length > 0) return `path: new MatcherPatternPathDynamic(
    ${node.regexp},
    ${generatePathParamsOptions(params, importsMap, paramParsersMap)},
    ${JSON.stringify(node.matcherPatternPathDynamicParts)},
    ${node.isSplat ? "null," : "/* trailingSlash */"}
  ),`;
	else return `path: new MatcherPatternPathStatic('${node.fullPath}'),`;
}
/**
* Generates the `query` property of a route record for the static resolver.
*/
function generateRouteRecordQuery({ node, importsMap, paramParsersMap }) {
	const queryParams = node.queryParams;
	if (queryParams.length === 0) return "";
	importsMap.add("vue-router/experimental", "MatcherPatternQueryParam");
	return `query: [
${queryParams.map((param) => {
		const parserOptions = generateParamParserOptions(param, importsMap, paramParsersMap);
		const args = [
			`'${param.paramName}'`,
			`'${param.paramName}'`,
			`'${param.format}'`
		];
		if (parserOptions || param.defaultValue !== void 0) args.push(parserOptions || "{}");
		if (param.defaultValue !== void 0) args.push(param.defaultValue);
		return `    new MatcherPatternQueryParam(${args.join(", ")})`;
	}).join(",\n")}
  ],`;
}
/**
* Generates a merge call for route records with definePage data in the experimental resolver format.
*/
function generateRouteRecordMerge(routeRecordObject, definePageDataList, importsMap) {
	if (definePageDataList.length === 0) return routeRecordObject;
	importsMap.add("vue-router/experimental", "_mergeRouteRecord");
	return `_mergeRouteRecord(
${routeRecordObject.split("\n").map((line) => {
		return line && `    ${line}`;
	}).join("\n")},
${definePageDataList.map((name) => `    ${name}`).join(",\n")}
  )`;
}

//#endregion
//#region src/core/context.ts
function createRoutesContext(options) {
	const { dts: preferDTS, root, routesFolder } = options;
	const dts = preferDTS === false ? false : preferDTS === true ? (0, pathe.resolve)(root, "typed-router.d.ts") : (0, pathe.resolve)(root, preferDTS);
	const routeTree = new PrefixTree(options);
	const editableRoutes = new EditableTreeNode(routeTree);
	const logger = new Proxy(console, { get(target, prop) {
		const res = Reflect.get(target, prop);
		if (typeof res === "function") return options.logs ? res : () => {};
		return res;
	} });
	const watchers = [];
	const paramParsersMap = /* @__PURE__ */ new Map();
	async function scanPages(startWatchers = true) {
		if (options.extensions.length < 1) throw new Error("\"extensions\" cannot be empty. Please specify at least one extension.");
		if (watchers.length > 0) return;
		const PARAM_PARSER_GLOB = "*.{ts,js}";
		const isParamParserMatch = (0, picomatch.default)(PARAM_PARSER_GLOB);
		await Promise.all([...routesFolder.map((folder) => resolveFolderOptions(options, folder)).map((folder) => {
			if (startWatchers) watchers.push(setupWatcher(new RoutesFolderWatcher(folder)));
			const ignorePattern = folder.exclude.map((f) => f.startsWith("**") ? f : (0, pathe.relative)(folder.src, f));
			return (0, tinyglobby.glob)(folder.pattern, {
				cwd: folder.src,
				ignore: ignorePattern,
				expandDirectories: false
			}).then((files) => Promise.all(files.map((file) => (0, pathe.resolve)(folder.src, file)).map((file) => addPage({
				routePath: require_options.asRoutePath(folder, file),
				filePath: file
			}))));
		}), ...options.experimental.paramParsers?.dir.map((folder) => {
			if (startWatchers) watchers.push(setupParamParserWatcher((0, chokidar.watch)(".", {
				cwd: folder,
				ignoreInitial: true,
				ignorePermissionErrors: true,
				ignored: (filePath, stats) => {
					if (!stats || stats.isDirectory()) return false;
					return !isParamParserMatch((0, pathe.relative)(folder, filePath));
				}
			}), folder));
			return (0, tinyglobby.glob)(PARAM_PARSER_GLOB, {
				cwd: folder,
				onlyFiles: true,
				expandDirectories: false
			}).then((paramParserFiles) => {
				for (const file of paramParserFiles) {
					const fileName = (0, pathe.parse)(file).name;
					const name = (0, scule.camelCase)(fileName);
					const absolutePath = (0, pathe.resolve)(folder, file);
					paramParsersMap.set(fileName, {
						name,
						typeName: `Param_${name}`,
						absolutePath,
						relativePath: (0, pathe.relative)(options.root, absolutePath)
					});
				}
				logger.log("Parsed param parsers", [...paramParsersMap].map((p) => p[0]));
			});
		}) || []]);
		for (const route of editableRoutes) await options.extendRoute?.(route);
		await _writeConfigFiles();
	}
	async function writeRouteInfoToNode(node, filePath) {
		const content = await node_fs.promises.readFile(filePath, "utf8");
		node.hasDefinePage ||= content.includes("definePage");
		const definedPageInfo = extractDefinePageInfo(content, filePath);
		const routeBlock = getRouteBlock(filePath, content, options);
		node.setCustomRouteBlock(filePath, {
			...routeBlock,
			...definedPageInfo
		});
		server?.invalidatePage(filePath);
	}
	async function addPage({ filePath, routePath }, triggerExtendRoute = false) {
		logger.log(`added "${routePath}" for "${filePath}"`);
		const node = routeTree.insert(routePath, filePath);
		await writeRouteInfoToNode(node, filePath);
		if (triggerExtendRoute) await options.extendRoute?.(new EditableTreeNode(node));
		server?.updateRoutes();
	}
	async function updatePage({ filePath, routePath }) {
		logger.log(`updated "${routePath}" for "${filePath}"`);
		const node = routeTree.getChild(filePath);
		if (!node) {
			logger.warn(`Cannot update "${filePath}": Not found.`);
			return;
		}
		await writeRouteInfoToNode(node, filePath);
		await options.extendRoute?.(new EditableTreeNode(node));
		server?.updateRoutes();
	}
	function removePage({ filePath, routePath }) {
		logger.log(`remove "${routePath}" for "${filePath}"`);
		routeTree.removeChild(filePath);
		server?.updateRoutes();
	}
	function setupParamParserWatcher(watcher, cwd) {
		logger.log(`🤖 Scanning param parsers in ${cwd}`);
		return watcher.on("add", (file) => {
			const fileName = (0, pathe.parse)(file).name;
			const name = (0, scule.camelCase)(fileName);
			const absolutePath = (0, pathe.resolve)(cwd, file);
			paramParsersMap.set(fileName, {
				name,
				typeName: `Param_${name}`,
				absolutePath,
				relativePath: "./" + (0, pathe.relative)(options.root, absolutePath)
			});
			writeConfigFiles();
		}).on("unlink", (file) => {
			paramParsersMap.delete((0, pathe.parse)(file).name);
			writeConfigFiles();
		});
	}
	function setupWatcher(watcher) {
		logger.log(`🤖 Scanning files in ${watcher.src}`);
		return watcher.on("change", async (ctx) => {
			await updatePage(ctx);
			writeConfigFiles();
		}).on("add", async (ctx) => {
			await addPage(ctx, true);
			writeConfigFiles();
		}).on("unlink", (ctx) => {
			removePage(ctx);
			writeConfigFiles();
		});
	}
	function generateResolver() {
		const importsMap = new require_options.ImportsMap();
		const resolverCode = generateRouteResolver(routeTree, options, importsMap, paramParsersMap);
		let imports = importsMap.toString();
		if (imports) imports += "\n";
		const hmr = ts`
export function handleHotUpdate(_router, _hotUpdateCallback) {
  if (import.meta.hot) {
    import.meta.hot.data.router = _router
    import.meta.hot.data.router_hotUpdateCallback = _hotUpdateCallback
  }
}

if (import.meta.hot) {
  import.meta.hot.accept((mod) => {
    const router = import.meta.hot.data.router
    if (!router) {
      import.meta.hot.invalidate('[unplugin-vue-router:HMR] Cannot replace the resolver because there is no active router. Reloading.')
      return
    }
    router._hmrReplaceResolver(mod.resolver)
    // call the hotUpdateCallback for custom updates
    import.meta.hot.data.router_hotUpdateCallback?.(mod.resolver)
    const route = router.currentRoute.value
    router.replace({
      path: route.path,
      query: route.query,
      hash: route.hash,
      force: true
    })
  })
}`;
		return `${imports}${resolverCode}\n${hmr}`;
	}
	function generateRoutes() {
		const importsMap = new require_options.ImportsMap();
		const routeList = `export const routes = ${generateRouteRecords(routeTree, options, importsMap)}\n`;
		const hmr = ts`
export function handleHotUpdate(_router, _hotUpdateCallback) {
  if (import.meta.hot) {
    import.meta.hot.data.router = _router
    import.meta.hot.data.router_hotUpdateCallback = _hotUpdateCallback
  }
}

if (import.meta.hot) {
  import.meta.hot.accept((mod) => {
    const router = import.meta.hot.data.router
    if (!router) {
      import.meta.hot.invalidate('[unplugin-vue-router:HMR] Cannot replace the routes because there is no active router. Reloading.')
      return
    }
    router.clearRoutes()
    for (const route of mod.routes) {
      router.addRoute(route)
    }
    // call the hotUpdateCallback for custom updates
    import.meta.hot.data.router_hotUpdateCallback?.(mod.routes)
    const route = router.currentRoute.value
    router.replace({
      ...route,
      // NOTE: we should be able to just do ...route but the router
      // currently skips resolving and can give errors with renamed routes
      // so we explicitly set remove matched and name
      name: undefined,
      matched: undefined,
      force: true
    })
  })
}
`;
		let imports = importsMap.toString();
		if (imports) imports += "\n";
		return `${imports}${routeList}${hmr}\n`;
	}
	function generateDTS$1() {
		if (options.experimental.paramParsers?.dir.length) warnMissingParamParsers(routeTree, paramParsersMap);
		return generateDTS({
			routesModule: MODULE_ROUTES_PATH,
			routeNamedMap: generateRouteNamedMap(routeTree, options, paramParsersMap),
			routeFileInfoMap: generateRouteFileInfoMap(routeTree, { root }),
			paramsTypesDeclaration: generateParamParsersTypesDeclarations(paramParsersMap),
			customParamsType: generateParamParserCustomType(paramParsersMap)
		});
	}
	let lastDTS;
	async function _writeConfigFiles() {
		logger.time("writeConfigFiles");
		if (options.beforeWriteFiles) {
			await options.beforeWriteFiles(editableRoutes);
			logger.timeLog("writeConfigFiles", "beforeWriteFiles()");
		}
		require_options.logTree(routeTree, logger.log);
		if (dts) {
			const content = generateDTS$1();
			if (lastDTS !== content) {
				await node_fs.promises.mkdir((0, pathe.dirname)(dts), { recursive: true });
				await node_fs.promises.writeFile(dts, content, "utf-8");
				logger.timeLog("writeConfigFiles", "wrote dts file");
				lastDTS = content;
				server?.updateRoutes();
			}
		}
		logger.timeEnd("writeConfigFiles");
	}
	const writeConfigFiles = require_options.throttle(_writeConfigFiles, 500, 100);
	function stopWatcher() {
		if (watchers.length) {
			logger.log("🛑 stopping watcher");
			watchers.forEach((watcher) => watcher.close());
		}
	}
	let server;
	function setServerContext(_server) {
		server = _server;
	}
	return {
		scanPages,
		writeConfigFiles,
		setServerContext,
		stopWatcher,
		generateRoutes,
		generateResolver,
		definePageTransform(code, id) {
			return definePageTransform({
				code,
				id
			});
		}
	};
}

//#endregion
//#region src/core/vite/index.ts
function createViteContext(server) {
	function invalidate(path$1) {
		const foundModule = server.moduleGraph.getModuleById(path$1);
		if (foundModule) return server.reloadModule(foundModule);
		return !!foundModule;
	}
	function invalidatePage(filepath) {
		const pageModules = server.moduleGraph.getModulesByFile(filepath);
		if (pageModules) return Promise.all([...pageModules].map((mod) => server.reloadModule(mod))).then(() => {});
		return false;
	}
	function reload() {
		server.ws.send({
			type: "full-reload",
			path: "*"
		});
	}
	/**
	* Triggers HMR for the vue-router/auto-routes module.
	*/
	async function updateRoutes() {
		const autoRoutesMod = server.moduleGraph.getModuleById(asVirtualId(MODULE_ROUTES_PATH));
		const autoResolvedMod = server.moduleGraph.getModuleById(asVirtualId(MODULE_RESOLVER_PATH));
		await Promise.all([autoRoutesMod && server.reloadModule(autoRoutesMod), autoResolvedMod && server.reloadModule(autoResolvedMod)]);
	}
	return {
		invalidate,
		invalidatePage,
		updateRoutes,
		reload
	};
}

//#endregion
//#region src/data-loaders/auto-exports.ts
function extractLoadersToExport(code, filterPaths, root) {
	return (0, mlly.findStaticImports)(code).flatMap((i) => {
		const parsed = (0, mlly.parseStaticImport)(i);
		if (!filterPaths((0, pathe.resolve)(root, parsed.specifier.startsWith("/") ? parsed.specifier.slice(1) : parsed.specifier))) return [];
		return [parsed.defaultImport, ...Object.values(parsed.namedImports || {})].filter((v) => !!v && !v.startsWith("_"));
	});
}
const PLUGIN_NAME = "unplugin-vue-router:data-loaders-auto-export";
/**
* Vite Plugin to automatically export loaders from page components.
*
* @param options Options
* @experimental - This API is experimental and can be changed in the future. It's used internally by `experimental.autoExportsDataLoaders`

*/
function AutoExportLoaders({ transformFilter, loadersPathsGlobs, root = process.cwd() }) {
	const filterPaths = (0, unplugin_utils.createFilter)(loadersPathsGlobs);
	return {
		name: PLUGIN_NAME,
		transform: {
			order: "post",
			filter: { id: transformFilter },
			handler(code) {
				const loadersToExports = extractLoadersToExport(code, filterPaths, root);
				if (loadersToExports.length <= 0) return;
				const s = new magic_string.default(code);
				s.append(`\nexport const __loaders = [\n${loadersToExports.join(",\n")}\n];\n`);
				return {
					code: s.toString(),
					map: s.generateMap()
				};
			}
		}
	};
}
function createAutoExportPlugin(options) {
	return {
		name: PLUGIN_NAME,
		vite: AutoExportLoaders(options)
	};
}

//#endregion
//#region src/index.ts
var src_default = (0, unplugin.createUnplugin)((opt = {}, _meta) => {
	const options = require_options.resolveOptions(opt);
	const ctx = createRoutesContext(options);
	function getVirtualId$1(id) {
		if (options._inspect) return id;
		return getVirtualId(id);
	}
	function asVirtualId$1(id) {
		if (options._inspect) return id;
		return asVirtualId(id);
	}
	const pageFilePattern = require_options.appendExtensionListToPattern(options.filePatterns, require_options.mergeAllExtensions(options));
	const IDS_TO_INCLUDE = options.routesFolder.flatMap((routeOption) => pageFilePattern.map((pattern) => (0, pathe.join)(routeOption.src, pattern)));
	const plugins = [{
		name: "unplugin-vue-router",
		enforce: "pre",
		resolveId: {
			filter: { id: { include: [
				/* @__PURE__ */ new RegExp(`^${MODULE_ROUTES_PATH}$`),
				/* @__PURE__ */ new RegExp(`^${MODULE_RESOLVER_PATH}$`),
				routeBlockQueryRE
			] } },
			handler(id) {
				if (id === MODULE_ROUTES_PATH || id === MODULE_RESOLVER_PATH) return asVirtualId$1(id);
				return ROUTE_BLOCK_ID;
			}
		},
		async buildStart() {
			await ctx.scanPages(options.watch);
		},
		buildEnd() {
			ctx.stopWatcher();
		},
		transform: {
			filter: { id: {
				include: [...IDS_TO_INCLUDE, DEFINE_PAGE_QUERY_RE],
				exclude: options.exclude
			} },
			handler(code, id) {
				return ctx.definePageTransform(code, id);
			}
		},
		load: {
			filter: { id: { include: [
				/* @__PURE__ */ new RegExp(`^${ROUTE_BLOCK_ID}$`),
				/* @__PURE__ */ new RegExp(`^${VIRTUAL_PREFIX}${MODULE_ROUTES_PATH}$`),
				/* @__PURE__ */ new RegExp(`^${VIRTUAL_PREFIX}${MODULE_RESOLVER_PATH}$`)
			] } },
			handler(id) {
				if (id === ROUTE_BLOCK_ID) return {
					code: `export default {}`,
					map: null
				};
				const resolvedId = getVirtualId$1(id);
				if (resolvedId === MODULE_ROUTES_PATH) {
					ROUTES_LAST_LOAD_TIME.update();
					return ctx.generateRoutes();
				}
				if (resolvedId === MODULE_RESOLVER_PATH) {
					ROUTES_LAST_LOAD_TIME.update();
					return ctx.generateResolver();
				}
			}
		},
		vite: { configureServer(server) {
			ctx.setServerContext(createViteContext(server));
		} }
	}];
	if (options.experimental.autoExportsDataLoaders) plugins.push(createAutoExportPlugin({
		transformFilter: {
			include: IDS_TO_INCLUDE,
			exclude: options.exclude
		},
		loadersPathsGlobs: options.experimental.autoExportsDataLoaders,
		root: options.root
	}));
	return plugins;
});
/**
* Adds useful auto imports to the AutoImport config:
* @example
* ```js
* import { VueRouterAutoImports } from 'unplugin-vue-router'
*
* AutoImport({
*   imports: [VueRouterAutoImports],
* }),
* ```
*/
const VueRouterAutoImports = {
	"vue-router": [
		"useRoute",
		"useRouter",
		"onBeforeRouteUpdate",
		"onBeforeRouteLeave"
	],
	"unplugin-vue-router/runtime": []
};

//#endregion
Object.defineProperty(exports, 'AutoExportLoaders', {
  enumerable: true,
  get: function () {
    return AutoExportLoaders;
  }
});
Object.defineProperty(exports, 'EditableTreeNode', {
  enumerable: true,
  get: function () {
    return EditableTreeNode;
  }
});
Object.defineProperty(exports, 'VueRouterAutoImports', {
  enumerable: true,
  get: function () {
    return VueRouterAutoImports;
  }
});
Object.defineProperty(exports, '__toESM', {
  enumerable: true,
  get: function () {
    return __toESM;
  }
});
Object.defineProperty(exports, 'createRoutesContext', {
  enumerable: true,
  get: function () {
    return createRoutesContext;
  }
});
Object.defineProperty(exports, 'createTreeNodeValue', {
  enumerable: true,
  get: function () {
    return createTreeNodeValue;
  }
});
Object.defineProperty(exports, 'src_default', {
  enumerable: true,
  get: function () {
    return src_default;
  }
});