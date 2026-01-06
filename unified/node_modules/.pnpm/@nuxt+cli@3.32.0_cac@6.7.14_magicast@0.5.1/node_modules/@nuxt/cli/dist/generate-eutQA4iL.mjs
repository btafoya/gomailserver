import { a as legacyRootDirArgs, i as extendsArgs, n as dotEnvArgs, o as logLevelArgs, r as envNameArgs, t as cwdArgs } from "./_shared-BCYCnX0T.mjs";
import "./logger-B4ge7MhP.mjs";
import "./env-BRiVSJMz.mjs";
import "./kit-4LNqcmNp.mjs";
import "./versions-DD7jbeRR.mjs";
import "./banner-BrAj2w2K.mjs";
import "./fs-BNfOTIDu.mjs";
import build_default from "./build-DF96Rh_E.mjs";
import { defineCommand } from "citty";

//#region ../nuxi/src/commands/generate.ts
var generate_default = defineCommand({
	meta: {
		name: "generate",
		description: "Build Nuxt and prerender all routes"
	},
	args: {
		...cwdArgs,
		...logLevelArgs,
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
		ctx.args.prerender = true;
		await build_default.run(ctx);
	}
});

//#endregion
export { generate_default as default };