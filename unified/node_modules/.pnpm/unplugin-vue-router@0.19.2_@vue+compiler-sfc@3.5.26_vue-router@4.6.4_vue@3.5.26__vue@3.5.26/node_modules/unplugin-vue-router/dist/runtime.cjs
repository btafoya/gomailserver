
//#region src/runtime.ts
/**
* Defines properties of the route for the current page component.
*
* @param route - route information to be added to this page
*/
const definePage = (route) => route;
/**
* Merges route records.
*
* @internal
*
* @param main - main route record
* @param routeRecords - route records to merge
* @returns merged route record
*/
function _mergeRouteRecord(main, ...routeRecords) {
	return routeRecords.reduce((acc, routeRecord) => {
		const meta = Object.assign({}, acc.meta, routeRecord.meta);
		const alias = [].concat(acc.alias || [], routeRecord.alias || []);
		Object.assign(acc, routeRecord);
		acc.meta = meta;
		acc.alias = alias;
		return acc;
	}, main);
}
/**
* TODO: native parsers ideas:
* - json -> just JSON.parse(value)
* - boolean -> 'true' | 'false' -> boolean
* - number -> Number(value) -> NaN if not a number
*/

//#endregion
exports._mergeRouteRecord = _mergeRouteRecord;
exports.definePage = definePage;