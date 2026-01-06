let muggle_string = require("muggle-string");
let pathe = require("pathe");

//#region src/volar/utils/augment-vls-ctx.ts
/**
* Augments the VLS context (volar) with additianal type information.
*
* @param content - content retrieved from the volar pluign
* @param codes - codes to add to the VLS context
*/
function augmentVlsCtx(content, codes) {
	let from = -1;
	for (let i = 0; i < content.length; i++) {
		const code = content[i];
		if (typeof code !== "string") continue;
		if (from === -1 && code.startsWith(`const __VLS_ctx`)) from = i;
		else if (from !== -1) {
			if (code === `}`) {
				content.splice(i, 0, ...codes.map((code$1) => `...${code$1},\n`));
				break;
			} else if (code === `;\n`) {
				content.splice(from + 1, i - from, `{\n`, `...`, ...content.slice(from + 1, i), `,\n`, ...codes.map((code$1) => `...${code$1},\n`), `}`, `;\n`);
				break;
			}
		}
	}
}

//#endregion
//#region src/volar/entries/sfc-typed-router.ts
const plugin = ({ compilerOptions, modules: { typescript: ts } }) => {
	const RE = { DOLLAR_ROUTE: { VLS_CTX: /\b__VLS_ctx.\$route\b/g } };
	return {
		version: 2.1,
		resolveEmbeddedCode(fileName, sfc, embeddedCode) {
			if (!embeddedCode.id.startsWith("script_")) return;
			const useRouteNameTypeParam = `<${`import('vue-router/auto-routes')._RouteNamesForFilePath<'${compilerOptions.rootDir ? (0, pathe.relative)(compilerOptions.rootDir, fileName) : fileName}'>`}>`;
			if (sfc.scriptSetup) visit(sfc.scriptSetup.ast);
			function visit(node) {
				if (ts.isCallExpression(node) && ts.isIdentifier(node.expression) && node.expression.text === "useRoute" && !node.typeArguments && !node.arguments.length) if (!sfc.scriptSetup.lang.startsWith("js")) (0, muggle_string.replaceSourceRange)(embeddedCode.content, sfc.scriptSetup.name, node.expression.end, node.expression.end, useRouteNameTypeParam);
				else {
					const start = node.getStart(sfc.scriptSetup.ast);
					(0, muggle_string.replaceSourceRange)(embeddedCode.content, sfc.scriptSetup.name, start, start, `(`);
					(0, muggle_string.replaceSourceRange)(embeddedCode.content, sfc.scriptSetup.name, node.end, node.end, ` as ReturnType<typeof useRoute${useRouteNameTypeParam}>)`);
				}
				else ts.forEachChild(node, visit);
			}
			const contentStr = (0, muggle_string.toString)(embeddedCode.content);
			const vlsCtxAugmentations = [];
			if (contentStr.match(RE.DOLLAR_ROUTE.VLS_CTX)) vlsCtxAugmentations.push(`{} as { $route: ReturnType<typeof import('vue-router').useRoute${useRouteNameTypeParam}> }`);
			if (vlsCtxAugmentations.length) augmentVlsCtx(embeddedCode.content, vlsCtxAugmentations);
		}
	};
};
var sfc_typed_router_default = plugin;

//#endregion
module.exports = sfc_typed_router_default;