import { a as legacyRootDirArgs, i as extendsArgs, n as dotEnvArgs, o as logLevelArgs, r as envNameArgs, t as cwdArgs } from "./_shared-BCYCnX0T.mjs";
import { n as logger } from "./logger-B4ge7MhP.mjs";
import { t as overrideEnv } from "./env-BRiVSJMz.mjs";
import { t as loadKit } from "./kit-4LNqcmNp.mjs";
import "./versions-DD7jbeRR.mjs";
import { n as showVersions } from "./banner-BrAj2w2K.mjs";
import { t as clearBuildDir } from "./fs-BNfOTIDu.mjs";
import process from "node:process";
import { defineCommand } from "citty";
import { colors } from "consola/utils";
import { intro, outro } from "@clack/prompts";
import { relative, resolve } from "pathe";

//#region ../nuxi/src/commands/build.ts
var build_default = defineCommand({
	meta: {
		name: "build",
		description: "Build Nuxt for production deployment"
	},
	args: {
		...cwdArgs,
		...logLevelArgs,
		prerender: {
			type: "boolean",
			description: "Build Nuxt and prerender static routes"
		},
		preset: {
			type: "string",
			description: "Nitro server preset"
		},
		...dotEnvArgs,
		...envNameArgs,
		...extendsArgs,
		...legacyRootDirArgs
	},
	async run(ctx) {
		overrideEnv("production");
		const cwd = resolve(ctx.args.cwd || ctx.args.rootDir);
		intro(colors.cyan("Building Nuxt for production..."));
		const kit = await loadKit(cwd);
		await showVersions(cwd, kit);
		const nuxt = await kit.loadNuxt({
			cwd,
			dotenv: {
				cwd,
				fileName: ctx.args.dotenv
			},
			envName: ctx.args.envName,
			overrides: {
				logLevel: ctx.args.logLevel,
				_generate: ctx.args.prerender,
				nitro: {
					static: ctx.args.prerender,
					preset: ctx.args.preset || process.env.NITRO_PRESET || process.env.SERVER_PRESET
				},
				...ctx.args.extends && { extends: ctx.args.extends },
				...ctx.data?.overrides
			}
		});
		let nitro;
		try {
			nitro = kit.useNitro?.();
			if (nitro) logger.info(`Nitro preset: ${colors.cyan(nitro.options.preset)}`);
		} catch {}
		await clearBuildDir(nuxt.options.buildDir);
		await kit.writeTypes(nuxt);
		nuxt.hook("build:error", (err) => {
			logger.error(`Nuxt build error: ${err}`);
			process.exit(1);
		});
		await kit.buildNuxt(nuxt);
		if (ctx.args.prerender) {
			if (!nuxt.options.ssr) {
				logger.warn(`HTML content not prerendered because ${colors.cyan("ssr: false")} was set.`);
				logger.info(`You can read more in ${colors.cyan("https://nuxt.com/docs/getting-started/deployment#static-hosting")}.`);
			}
			const dir = nitro?.options.output.publicDir;
			const publicDir = dir ? relative(process.cwd(), dir) : ".output/public";
			outro(`✨ You can now deploy ${colors.cyan(publicDir)} to any static hosting!`);
		} else outro("✨ Build complete!");
	}
});

//#endregion
export { build_default as default };