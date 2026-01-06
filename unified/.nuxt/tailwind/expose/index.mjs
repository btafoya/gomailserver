import theme from "#build/tailwind/expose/theme.mjs"
import corePlugins from "#build/tailwind/expose/corePlugins.mjs"
import plugins from "#build/tailwind/expose/plugins.mjs"
import content from "#build/tailwind/expose/content.mjs"
import safelist from "#build/tailwind/expose/safelist.mjs"
import darkMode from "#build/tailwind/expose/darkMode.mjs"
import variants from "#build/tailwind/expose/variants.mjs"
import presets from "#build/tailwind/expose/presets.mjs"
import prefix from "#build/tailwind/expose/prefix.mjs"
import important from "#build/tailwind/expose/important.mjs"
import separator from "#build/tailwind/expose/separator.mjs"
import blocklist from "#build/tailwind/expose/blocklist.mjs"
const config = { theme, corePlugins, plugins, content, safelist, darkMode, variants, presets, prefix, important, separator, blocklist }
export { config as default, theme, corePlugins, plugins, content, safelist, darkMode, variants, presets, prefix, important, separator, blocklist }