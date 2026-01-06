import { a as legacyRootDirArgs, i as extendsArgs, n as dotEnvArgs, o as logLevelArgs, r as envNameArgs, t as cwdArgs } from "./_shared-BCYCnX0T.mjs";
import { n as logger } from "./logger-B4ge7MhP.mjs";
import { r as relativeToProcess, t as loadKit } from "./kit-4LNqcmNp.mjs";
import { t as clearBuildDir } from "./fs-BNfOTIDu.mjs";
import process from "node:process";
import { defineCommand } from "citty";
import { colors } from "consola/utils";
import { resolve } from "pathe";

//#region ../nuxi/src/commands/prepare.ts
var prepare_default = defineCommand({
	meta: {
		name: "prepare",
		description: "Prepare Nuxt for development/build"
	},
	args: {
		...dotEnvArgs,
		...cwdArgs,
		...logLevelArgs,
		...envNameArgs,
		...extendsArgs,
		...legacyRootDirArgs
	},
	async run(ctx) {
		process.env.NODE_ENV = process.env.NODE_ENV || "production";
		const cwd = resolve(ctx.args.cwd || ctx.args.rootDir);
		const { loadNuxt, buildNuxt, writeTypes } = await loadKit(cwd);
		const nuxt = await loadNuxt({
			cwd,
			dotenv: {
				cwd,
				fileName: ctx.args.dotenv
			},
			envName: ctx.args.envName,
			overrides: {
				_prepare: true,
				logLevel: ctx.args.logLevel,
				...ctx.args.extends && { extends: ctx.args.extends },
				...ctx.data?.overrides
			}
		});
		await clearBuildDir(nuxt.options.buildDir);
		await buildNuxt(nuxt);
		await writeTypes(nuxt);
		logger.success(`Types generated in ${colors.cyan(relativeToProcess(nuxt.options.buildDir))}.`);
	}
});

//#endregion
export { prepare_default as t };