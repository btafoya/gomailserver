<template>
  <component :is="as" :class="skeletonClass" v-bind="attrs" />
</template>

<script>
import { computed, toRef, defineComponent } from "vue";
import { twJoin } from "tailwind-merge";
import { useUI } from "../../composables/useUI";
import { mergeConfig, twMerge } from "../../utils";
import appConfig from "#build/app.config";
import { skeleton } from "#ui/ui.config";
const config = mergeConfig(appConfig.ui.strategy, appConfig.ui.skeleton, skeleton);
export default defineComponent({
  inheritAttrs: false,
  props: {
    as: {
      type: String,
      default: "div"
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
    const { ui, attrs } = useUI("skeleton", toRef(props, "ui"), config);
    const skeletonClass = computed(() => {
      return twMerge(twJoin(
        ui.value.base,
        ui.value.background,
        ui.value.rounded
      ), props.class);
    });
    return {
      // eslint-disable-next-line vue/no-dupe-keys
      ui,
      attrs,
      skeletonClass
    };
  }
});
</script>
