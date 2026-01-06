import { l as RootCommand } from "./t-Cao2EVMz.js";
import { t as CompletionConfig } from "./shared-Bzh57-vW.js";
import { ArgsDef, CommandDef } from "citty";

//#region src/citty.d.ts
declare function tab<TArgs extends ArgsDef>(instance: CommandDef<TArgs>, completionConfig?: CompletionConfig): Promise<RootCommand>;
//#endregion
export { tab as default };