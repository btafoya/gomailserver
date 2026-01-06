import { defineNuxtModule, createResolver, addPlugin, addImports, getLayerDirectories, addImportsDir, addVitePlugin } from '@nuxt/kit';
import { fileURLToPath } from 'node:url';

function getStoreDeclaration(nodes) {
  return nodes?.find(
    (x) => x.init?.type === "CallExpression" && x.init.callee.type === "Identifier" && x.init.callee.name === "defineStore"
  );
}
function nameFromDeclaration(node) {
  return node?.id.type === "Identifier" && node.id.name;
}
function autoRegisterHMRPlugin(rootDir) {
  return {
    name: "pinia:auto-hmr-registration",
    transform(code, id) {
      if (id.startsWith("\0")) return;
      if (!id.startsWith(rootDir)) return;
      if (!code.includes("defineStore") || code.includes("acceptHMRUpdate")) {
        return;
      }
      const ast = this.parse(code);
      for (const n of ast.body) {
        if (n.type === "VariableDeclaration" || n.type === "ExportNamedDeclaration") {
          const storeDeclaration = getStoreDeclaration(
            n.type === "VariableDeclaration" ? n.declarations : n.declaration?.type === "VariableDeclaration" ? n.declaration?.declarations : void 0
          );
          const storeName = nameFromDeclaration(storeDeclaration);
          if (storeName) {
            return {
              code: [
                `import { acceptHMRUpdate } from 'pinia'`,
                code,
                "if (import.meta.hot) {",
                `  import.meta.hot.accept(acceptHMRUpdate(${storeName}, import.meta.hot))`,
                "}"
              ].join("\n")
            };
          }
        }
      }
    }
  };
}

const module = defineNuxtModule({
  meta: {
    name: "pinia",
    configKey: "pinia",
    compatibility: {
      nuxt: "^3.15.0 || ^4.0.0"
    }
  },
  defaults: {},
  setup(options, nuxt) {
    const { resolve } = createResolver(import.meta.url);
    const runtimeDir = fileURLToPath(new URL("./runtime", import.meta.url));
    nuxt.options.build.transpile.push(resolve(runtimeDir));
    nuxt.options.vite.optimizeDeps ??= {};
    nuxt.options.vite.optimizeDeps.exclude ??= [];
    if (!nuxt.options.vite.optimizeDeps.exclude.includes("pinia")) {
      nuxt.options.vite.optimizeDeps.exclude.push("pinia");
    }
    nuxt.hook("prepare:types", ({ references }) => {
      references.push({ types: "@pinia/nuxt" });
    });
    nuxt.hook("modules:done", () => {
      addPlugin(resolve(runtimeDir, "plugin.vue3"));
      addPlugin(resolve(runtimeDir, "payload-plugin"));
    });
    const composables = resolve(runtimeDir, "composables");
    addImports([
      { from: composables, name: "defineStore" },
      { from: composables, name: "acceptHMRUpdate" },
      { from: composables, name: "usePinia" },
      { from: composables, name: "storeToRefs" }
    ]);
    if (!options.storesDirs) {
      options.storesDirs = [resolve(nuxt.options.srcDir, "stores")];
    }
    if (options.storesDirs) {
      const layers = getLayerDirectories(nuxt);
      for (const storeDir of options.storesDirs) {
        for (const layer of layers) {
          addImportsDir(resolve(layer.app, storeDir));
        }
      }
    }
    if (nuxt.options.dev) {
      addVitePlugin(autoRegisterHMRPlugin(resolve(nuxt.options.rootDir)));
    }
  }
});

export { module as default };
