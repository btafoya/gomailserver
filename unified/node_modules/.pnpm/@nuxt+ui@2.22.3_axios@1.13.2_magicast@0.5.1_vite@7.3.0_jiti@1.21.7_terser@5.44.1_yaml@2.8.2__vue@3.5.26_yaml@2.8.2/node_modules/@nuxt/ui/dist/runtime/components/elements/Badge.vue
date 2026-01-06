<template>
  <span :class="badgeClass" v-bind="attrs">
    <slot name="leading">
      <UIcon v-if="isLeading && leadingIconName" :name="leadingIconName" :class="leadingIconClass" aria-hidden="true" />
    </slot>

    <slot>
      <span v-if="label !== void 0 && label !== null">
        {{ label }}
      </span>
    </slot>

    <slot name="trailing">
      <UIcon v-if="isTrailing && trailingIconName" :name="trailingIconName" :class="trailingIconClass" aria-hidden="true" />
    </slot>
  </span>
</template>

<script>
import { computed, toRef, defineComponent } from "vue";
import { twJoin } from "tailwind-merge";
import UIcon from "../elements/Icon.vue";
import { useUI } from "../../composables/useUI";
import { mergeConfig, twMerge } from "../../utils";
import { useInjectButtonGroup } from "../../composables/useButtonGroup";
import appConfig from "#build/app.config";
import { badge } from "#ui/ui.config";
const config = mergeConfig(appConfig.ui.strategy, appConfig.ui.badge, badge);
export default defineComponent({
  components: {
    UIcon
  },
  inheritAttrs: false,
  props: {
    size: {
      type: String,
      default: () => config.default.size,
      validator(value) {
        return Object.keys(config.size).includes(value);
      }
    },
    color: {
      type: String,
      default: () => config.default.color,
      validator(value) {
        return [...appConfig.ui.colors, ...Object.keys(config.color)].includes(value);
      }
    },
    variant: {
      type: String,
      default: () => config.default.variant,
      validator(value) {
        return [
          ...Object.keys(config.variant),
          ...Object.values(config.color).flatMap((value2) => Object.keys(value2))
        ].includes(value);
      }
    },
    label: {
      type: [String, Number],
      default: null
    },
    icon: {
      type: String,
      default: null
    },
    leadingIcon: {
      type: String,
      default: null
    },
    trailingIcon: {
      type: String,
      default: null
    },
    trailing: {
      type: Boolean,
      default: false
    },
    leading: {
      type: Boolean,
      default: false
    },
    class: {
      type: [String, Object, Array],
      default: () => ""
    },
    ui: {
      type: Object,
      default: () => ({})
    }
  },
  setup(props) {
    const { ui, attrs } = useUI("badge", toRef(props, "ui"), config);
    const { size, rounded } = useInjectButtonGroup({ ui, props });
    const isLeading = computed(() => {
      return props.icon && props.leading || props.icon && !props.trailing || !props.trailing || props.leadingIcon;
    });
    const isTrailing = computed(() => {
      return props.icon && props.trailing || props.trailing || props.trailingIcon;
    });
    const badgeClass = computed(() => {
      const variant = ui.value.color?.[props.color]?.[props.variant] || ui.value.variant[props.variant];
      return twMerge(twJoin(
        ui.value.base,
        ui.value.font,
        rounded.value,
        ui.value.size[size.value],
        ui.value.gap[size.value],
        variant?.replaceAll("{color}", props.color)
      ), props.class);
    });
    const leadingIconName = computed(() => {
      return props.leadingIcon || props.icon;
    });
    const trailingIconName = computed(() => {
      return props.trailingIcon || props.icon;
    });
    const leadingIconClass = computed(() => {
      return twJoin(
        ui.value.icon.base,
        ui.value.icon.size[size.value]
      );
    });
    const trailingIconClass = computed(() => {
      return twJoin(
        ui.value.icon.base,
        ui.value.icon.size[size.value]
      );
    });
    return {
      attrs,
      isLeading,
      isTrailing,
      badgeClass,
      leadingIconName,
      trailingIconName,
      leadingIconClass,
      trailingIconClass
    };
  }
});
</script>
