import { C as createTreeNodeValue, S as TreeNodeValueStatic, _ as getPascalCaseRouteName, a as ResolvedOptions, b as TreeNodeValueGroup, g as getFileBasedRouteName, h as EditableTreeNode, l as ServerContext, r as Options, t as DEFAULT_OPTIONS, v as TreeNode, x as TreeNodeValueParam, y as TreeNodeValue } from "./options-DG3niQXy.mjs";
import "./types-CF4U0ZkR.mjs";
import * as unplugin0 from "unplugin";
import { StringFilter } from "unplugin";
import { Plugin } from "vite";

//#region src/core/context.d.ts
declare function createRoutesContext(options: ResolvedOptions): {
  scanPages: (startWatchers?: boolean) => Promise<void>;
  writeConfigFiles: () => void;
  setServerContext: (_server: ServerContext) => void;
  stopWatcher: () => void;
  generateRoutes: () => string;
  generateResolver: () => string;
  definePageTransform(code: string, id: string): unplugin0.Thenable<unplugin0.TransformResult>;
};
//#endregion
//#region src/data-loaders/auto-exports.d.ts
/**
 * {@link AutoExportLoaders} options.
 */
interface AutoExportLoadersOptions {
  /**
   * Filter page components to apply the auto-export. Passed to `transform.filter.id`.
   */
  transformFilter: StringFilter;
  /**
   * Globs to match the paths of the loaders.
   */
  loadersPathsGlobs: string | string[];
  /**
   * Root of the project. All paths are resolved relatively to this one.
   * @default `process.cwd()`
   */
  root?: string;
}
/**
 * Vite Plugin to automatically export loaders from page components.
 *
 * @param options Options
 * @experimental - This API is experimental and can be changed in the future. It's used internally by `experimental.autoExportsDataLoaders`

 */
declare function AutoExportLoaders({
  transformFilter,
  loadersPathsGlobs,
  root
}: AutoExportLoadersOptions): Plugin;
//#endregion
//#region src/index.d.ts
declare const _default: unplugin0.UnpluginInstance<Options | undefined, boolean>;
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
declare const VueRouterAutoImports: Record<string, Array<string | [importName: string, alias: string]>>;
//#endregion
export { AutoExportLoaders, type AutoExportLoadersOptions, DEFAULT_OPTIONS, EditableTreeNode, Options, TreeNode, TreeNodeValue, TreeNodeValueGroup, TreeNodeValueParam, TreeNodeValueStatic, VueRouterAutoImports, createRoutesContext, createTreeNodeValue, _default as default, getFileBasedRouteName, getPascalCaseRouteName };