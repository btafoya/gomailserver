<template>
  <component
    :is="as"
    v-if="!to"
    :type="type"
    :disabled="disabled"
    v-bind="$attrs"
    :class="active ? activeClass : inactiveClass"
  >
    <slot v-bind="{ isActive: active }" />
  </component>
  <NuxtLink
    v-else
    v-slot="{ route, href, target, rel, navigate, isActive, isExactActive, isExternal }"
    v-bind="$props"
    custom
  >
    <a
      v-bind="$attrs"
      :href="!disabled ? href : void 0"
      :aria-disabled="disabled ? 'true' : void 0"
      :role="disabled ? 'link' : void 0"
      :rel="rel"
      :target="target"
      :class="active !== void 0 ? active ? activeClass : inactiveClass : resolveLinkClass(route, $route, { isActive, isExactActive })"
      :tabindex="disabled ? -1 : void 0"
      @click="(e) => !isExternal && !disabled && navigate(e)"
    >
      <slot v-bind="{ isActive: active !== void 0 ? active : exact ? isExactActive : isActive }" />
    </a>
  </NuxtLink>
</template>

<script>
import { isEqual, diff } from "ohash/utils";
import { defineComponent } from "vue";
import { nuxtLinkProps } from "../../utils";
export default defineComponent({
  inheritAttrs: false,
  props: {
    ...nuxtLinkProps,
    as: {
      type: String,
      default: "button"
    },
    type: {
      type: String,
      default: "button"
    },
    disabled: {
      type: Boolean,
      default: null
    },
    active: {
      type: Boolean,
      default: void 0
    },
    exact: {
      type: Boolean,
      default: false
    },
    exactQuery: {
      type: [Boolean, String],
      default: false
    },
    exactHash: {
      type: Boolean,
      default: false
    },
    inactiveClass: {
      type: String,
      default: void 0
    }
  },
  setup(props) {
    function isPartiallyEqual(item1, item2) {
      const diffedKeys = diff(item1, item2).reduce((filtered, q) => {
        if (q.type === "added") {
          filtered.add(q.key);
        }
        return filtered;
      }, /* @__PURE__ */ new Set());
      const item1Filtered = Object.fromEntries(Object.entries(item1).filter(([key]) => !diffedKeys.has(key)));
      const item2Filtered = Object.fromEntries(Object.entries(item2).filter(([key]) => !diffedKeys.has(key)));
      return isEqual(item1Filtered, item2Filtered);
    }
    function resolveLinkClass(route, $route, { isActive, isExactActive }) {
      if (props.exactQuery === "partial") {
        if (!isPartiallyEqual(route.query, $route.query)) return props.inactiveClass;
      } else if (props.exactQuery === true) {
        if (!isEqual(route.query, $route.query)) return props.inactiveClass;
      }
      if (props.exactHash && route.hash !== $route.hash) {
        return props.inactiveClass;
      }
      if (props.exact && isExactActive) {
        return props.activeClass;
      }
      if (!props.exact && isActive) {
        return props.activeClass;
      }
      return props.inactiveClass;
    }
    return {
      resolveLinkClass
    };
  }
});
</script>
