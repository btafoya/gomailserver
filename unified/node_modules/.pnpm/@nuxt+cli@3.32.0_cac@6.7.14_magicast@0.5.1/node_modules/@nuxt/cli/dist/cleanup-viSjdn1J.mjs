import { a as legacyRootDirArgs, t as cwdArgs } from "./_shared-BCYCnX0T.mjs";
import { n as logger } from "./logger-B4ge7MhP.mjs";
import { t as loadKit } from "./kit-4LNqcmNp.mjs";
import "./fs-BNfOTIDu.mjs";
import { t as cleanupNuxtDirs } from "./nuxt-4ILhz74C.mjs";
import { defineCommand } from "citty";
import { resolve } from "pathe";

//#region ../nuxi/src/commands/cleanup.ts
var cleanup_default = defineCommand({
	meta: {
		name: "cleanup",
		description: "Clean up generated Nuxt files and caches"
	},
	args: {
		...cwdArgs,
		...legacyRootDirArgs
	},
	async run(ctx) {
		const cwd = resolve(ctx.args.cwd || ctx.args.rootDir);
		const { loadNuxtConfig } = await loadKit(cwd);
		const nuxtOptions = await loadNuxtConfig({
			cwd,
			overrides: { dev: true }
		});
		await cleanupNuxtDirs(nuxtOptions.rootDir, nuxtOptions.buildDir);
		logger.success("Cleanup complete!");
	}
});

//#endregion
export { cleanup_default as default };