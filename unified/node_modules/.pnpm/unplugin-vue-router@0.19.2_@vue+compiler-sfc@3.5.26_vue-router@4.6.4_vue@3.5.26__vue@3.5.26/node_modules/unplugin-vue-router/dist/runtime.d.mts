import { RouteRecordRaw, TypesConfig } from "vue-router";

//#region src/runtime.d.ts

/**
 * Defines properties of the route for the current page component.
 *
 * @param route - route information to be added to this page
 */
declare const definePage: (route: DefinePage) => DefinePage;
/**
 * Merges route records.
 *
 * @internal
 *
 * @param main - main route record
 * @param routeRecords - route records to merge
 * @returns merged route record
 */
declare function _mergeRouteRecord(main: RouteRecordRaw, ...routeRecords: Partial<RouteRecordRaw>[]): RouteRecordRaw;
/**
 * Type to define a page. Can be augmented to add custom properties.
 */
interface DefinePage extends Partial<Omit<RouteRecordRaw, 'children' | 'components' | 'component' | 'name'>> {
  /**
   * A route name. If not provided, the name will be generated based on the file path.
   * Can be set to `false` to remove the name from types.
   */
  name?: string | false;
  /**
   * Custom parameters for the route. Requires `experimental.paramParsers` enabled.
   *
   * @experimental
   */
  params?: {
    path?: Record<string, ParamParserType>;
    /**
     * Parameters extracted from the query.
     */
    query?: Record<string, DefinePageQueryParamOptions | ParamParserType>;
  };
}
type ParamParserType_Native = 'int' | 'bool';
type ParamParserType = (TypesConfig extends Record<'ParamParsers', infer ParamParsers> ? ParamParsers : never) | ParamParserType_Native;
/**
 * Configures how to extract a route param from a specific query parameter.
 */
interface DefinePageQueryParamOptions<T = unknown> {
  /**
   * The type of the query parameter. Allowed values are native param parsers
   * and any parser in the {@link https://uvr.esm.is/TODO | params folder }. If
   * not provided, the value will kept as is.
   */
  parser?: ParamParserType;
  /**
   * Default value if the query parameter is missing or if the match fails
   * (e.g. a invalid number is passed to the int param parser). If not provided
   * and the param parser throws, the route will not match.
   */
  default?: (() => T) | T;
  /**
   * How to format the query parameter value.
   *
   * - 'value' - keep the first value only and pass that to parser
   * - 'array' - keep all values (even one or none) as an array and pass that to parser
   *
   * @default 'value'
   */
  format?: 'value' | 'array';
}
//#endregion
export { DefinePage, DefinePageQueryParamOptions, ParamParserType, ParamParserType_Native, _mergeRouteRecord, definePage };