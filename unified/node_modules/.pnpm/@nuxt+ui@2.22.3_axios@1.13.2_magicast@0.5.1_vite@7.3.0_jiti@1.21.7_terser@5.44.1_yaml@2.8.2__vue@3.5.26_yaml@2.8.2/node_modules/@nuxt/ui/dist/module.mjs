import { createRequire } from 'node:module';
import { useNuxt, addTemplate, createResolver, installModule, defineNuxtModule, addPlugin, addComponentsDir, addImportsDir } from '@nuxt/kit';
import { join } from 'pathe';
import { defu } from 'defu';
import { setGlobalColors } from '../dist/runtime/utils/colors.js';

const name = "@nuxt/ui";
const version = "2.22.3";

function createTemplates(nuxt = useNuxt()) {
  const template = addTemplate({
    filename: "ui.colors.mjs",
    getContents: () => `export default ${JSON.stringify(nuxt.options.appConfig.ui.colors)};`,
    write: true
  });
  const typesTemplate = addTemplate({
    filename: "ui.colors.d.ts",
    getContents: () => `declare module '#ui-colors' { const defaultExport: ${JSON.stringify(nuxt.options.appConfig.ui.colors)}; export default defaultExport; }`,
    write: true
  });
  nuxt.options.alias["#ui-colors"] = template.dst;
  nuxt.hook("prepare:types", (opts) => {
    opts.references.push({ path: typesTemplate.dst });
  });
}

function installTailwind(moduleOptions, nuxt = useNuxt(), resolve = createResolver(import.meta.url).resolve) {
  const runtimeDir = resolve("./runtime");
  nuxt.hook("tailwindcss:config", function(tailwindConfig) {
    tailwindConfig.theme = tailwindConfig.theme || {};
    tailwindConfig.theme.extend = tailwindConfig.theme.extend || {};
    tailwindConfig.theme.extend.colors = tailwindConfig.theme.extend.colors || {};
    const colors = setGlobalColors(tailwindConfig.theme);
    nuxt.options.appConfig.ui = {
      primary: "green",
      gray: "cool",
      colors,
      strategy: "merge"
    };
  });
  const configTemplate = addTemplate({
    filename: "nuxtui-tailwind.config.mjs",
    write: true,
    getContents: ({ nuxt: nuxt2 }) => `
      import { defaultExtractor as createDefaultExtractor } from "tailwindcss/lib/lib/defaultExtractor.js";
      import { customSafelistExtractor, generateSafelist } from ${JSON.stringify(resolve(runtimeDir, "utils", "colors"))};
      import formsPlugin from "@tailwindcss/forms";
      import aspectRatio from "@tailwindcss/aspect-ratio";
      import typography from "@tailwindcss/typography";
      import containerQueries from "@tailwindcss/container-queries";
      import headlessUi from "@headlessui/tailwindcss";

      const defaultExtractor = createDefaultExtractor({ tailwindConfig: { separator: ':' } });

      export default {
        plugins: [
          formsPlugin({ strategy: 'class' }),
          aspectRatio,
          typography,
          containerQueries,
          headlessUi
        ],
        content: {
          files: [
            ${JSON.stringify(resolve(runtimeDir, "components/**/*.{vue,mjs,ts}"))},
            ${JSON.stringify(resolve(runtimeDir, "ui.config/**/*.{mjs,js,ts}"))}
          ],
          transform: {
            vue: (content) => {
              return content.replaceAll(/(?:\\r\\n|\\r|\\n)/g, ' ')
            }
          },
          extract: {
            vue: (content) => {
              return [
                ...defaultExtractor(content),
                ...customSafelistExtractor(${JSON.stringify(moduleOptions.prefix)}, content, ${JSON.stringify(nuxt2.options.appConfig.ui.colors)}, ${JSON.stringify(moduleOptions.safelistColors)})
              ]
            }
          }
        },
        safelist: generateSafelist(${JSON.stringify(moduleOptions.safelistColors || [])}, ${JSON.stringify(nuxt2.options.appConfig.ui.colors)}),
      }
    `
  });
  const { configPath: userTwConfigPath = [], ...twModuleConfig } = nuxt.options.tailwindcss ?? {};
  const twConfigPaths = [
    configTemplate.dst,
    join(nuxt.options.rootDir, "tailwind.config")
  ];
  if (typeof userTwConfigPath === "string") {
    twConfigPaths.push(userTwConfigPath);
  } else {
    twConfigPaths.push(...userTwConfigPath);
  }
  return installModule("@nuxtjs/tailwindcss", defu({
    exposeConfig: true,
    config: {
      darkMode: "class"
    },
    configPath: twConfigPaths
  }, twModuleConfig));
}

const _require = createRequire(import.meta.url);
const defaultColors = _require("tailwindcss/colors.js");
delete defaultColors.lightBlue;
delete defaultColors.warmGray;
delete defaultColors.trueGray;
delete defaultColors.coolGray;
delete defaultColors.blueGray;
const module = defineNuxtModule({
  meta: {
    name,
    version,
    configKey: "ui",
    compatibility: {
      nuxt: ">=3.10.0"
    }
  },
  defaults: {
    prefix: "U",
    colorMode: true,
    safelistColors: ["primary"],
    disableGlobalStyles: false
  },
  async setup(options, nuxt) {
    const { resolve } = createResolver(import.meta.url);
    const runtimeDir = resolve("./runtime");
    nuxt.options.build.transpile.push(runtimeDir);
    nuxt.options.build.transpile.push("@popperjs/core", "@headlessui/vue");
    nuxt.options.alias["#ui"] = runtimeDir;
    if (!options.disableGlobalStyles) {
      nuxt.options.css.push(resolve(runtimeDir, "ui.css"));
    }
    createTemplates(nuxt);
    await installModule("@nuxt/icon");
    if (options.colorMode) {
      await installModule("@nuxtjs/color-mode", { classSuffix: "" });
    }
    await installTailwind(options, nuxt, resolve);
    addPlugin({
      src: resolve(runtimeDir, "plugins", "colors")
    });
    addPlugin({
      src: resolve(runtimeDir, "plugins", "modals")
    });
    addPlugin({
      src: resolve(runtimeDir, "plugins", "slideovers")
    });
    addComponentsDir({
      path: resolve(runtimeDir, "components", "elements"),
      prefix: options.prefix,
      global: options.global,
      watch: false
    });
    addComponentsDir({
      path: resolve(runtimeDir, "components", "forms"),
      prefix: options.prefix,
      global: options.global,
      watch: false
    });
    addComponentsDir({
      path: resolve(runtimeDir, "components", "data"),
      prefix: options.prefix,
      global: options.global,
      watch: false
    });
    addComponentsDir({
      path: resolve(runtimeDir, "components", "layout"),
      prefix: options.prefix,
      global: options.global,
      watch: false
    });
    addComponentsDir({
      path: resolve(runtimeDir, "components", "navigation"),
      prefix: options.prefix,
      global: options.global,
      watch: false
    });
    addComponentsDir({
      path: resolve(runtimeDir, "components", "overlays"),
      prefix: options.prefix,
      global: options.global,
      watch: false
    });
    addImportsDir(resolve(runtimeDir, "composables"));
  }
});

export { module as default };
