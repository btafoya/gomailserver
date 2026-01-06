
      import { defaultExtractor as createDefaultExtractor } from "tailwindcss/lib/lib/defaultExtractor.js";
      import { customSafelistExtractor, generateSafelist } from "/home/btafoya/projects/gomailserver/unified/node_modules/.pnpm/@nuxt+ui@2.22.3_axios@1.13.2_magicast@0.5.1_vite@7.3.0_jiti@1.21.7_terser@5.44.1_yaml@2.8.2__vue@3.5.26_yaml@2.8.2/node_modules/@nuxt/ui/dist/runtime/utils/colors";
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
            "/home/btafoya/projects/gomailserver/unified/node_modules/.pnpm/@nuxt+ui@2.22.3_axios@1.13.2_magicast@0.5.1_vite@7.3.0_jiti@1.21.7_terser@5.44.1_yaml@2.8.2__vue@3.5.26_yaml@2.8.2/node_modules/@nuxt/ui/dist/runtime/components/**/*.{vue,mjs,ts}",
            "/home/btafoya/projects/gomailserver/unified/node_modules/.pnpm/@nuxt+ui@2.22.3_axios@1.13.2_magicast@0.5.1_vite@7.3.0_jiti@1.21.7_terser@5.44.1_yaml@2.8.2__vue@3.5.26_yaml@2.8.2/node_modules/@nuxt/ui/dist/runtime/ui.config/**/*.{mjs,js,ts}"
          ],
          transform: {
            vue: (content) => {
              return content.replaceAll(/(?:\r\n|\r|\n)/g, ' ')
            }
          },
          extract: {
            vue: (content) => {
              return [
                ...defaultExtractor(content),
                ...customSafelistExtractor("U", content, ["red","orange","amber","yellow","lime","green","emerald","teal","cyan","sky","blue","indigo","violet","purple","fuchsia","pink","rose","primary"], ["primary"])
              ]
            }
          }
        },
        safelist: generateSafelist(["primary"], ["red","orange","amber","yellow","lime","green","emerald","teal","cyan","sky","blue","indigo","violet","purple","fuchsia","pink","rose","primary"]),
      }
    