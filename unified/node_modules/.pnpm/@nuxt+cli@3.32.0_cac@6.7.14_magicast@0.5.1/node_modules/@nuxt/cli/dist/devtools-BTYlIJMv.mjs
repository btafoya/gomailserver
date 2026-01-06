import { a as legacyRootDirArgs, t as cwdArgs } from "./_shared-BCYCnX0T.mjs";
import { n as logger } from "./logger-B4ge7MhP.mjs";
import process from "node:process";
import { defineCommand } from "citty";
import { colors } from "consola/utils";
import { resolve } from "pathe";
import { x } from "tinyexec";

//#region ../nuxi/src/commands/devtools.ts
var devtools_default = defineCommand({
	meta: {
		name: "devtools",
		description: "Enable or disable devtools in a Nuxt project"
	},
	args: {
		...cwdArgs,
		command: {
			type: "positional",
			description: "Command to run",
			valueHint: "enable|disable"
		},
		...legacyRootDirArgs
	},
	async run(ctx) {
		const cwd = resolve(ctx.args.cwd || ctx.args.rootDir);
		if (!["enable", "disable"].includes(ctx.args.command)) {
			logger.error(`Unknown command ${colors.cyan(ctx.args.command)}.`);
			process.exit(1);
		}
		await x("npx", [
			"@nuxt/devtools-wizard@latest",
			ctx.args.command,
			cwd
		], {
			throwOnError: true,
			nodeOptions: {
				stdio: "inherit",
				cwd
			}
		});
	}
});

//#endregion
export { devtools_default as default };