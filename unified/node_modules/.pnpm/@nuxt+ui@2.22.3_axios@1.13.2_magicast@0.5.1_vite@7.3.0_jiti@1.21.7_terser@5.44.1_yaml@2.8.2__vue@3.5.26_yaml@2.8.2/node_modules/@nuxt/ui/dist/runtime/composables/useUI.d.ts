import type { Ref } from 'vue';
import type { DeepPartial, Strategy } from '../types/index.js';
export declare const useUI: <T>(key: any, $ui?: Ref<(DeepPartial<T> & {
    strategy?: Strategy;
}) | undefined>, $config?: Ref<T> | T, $wrapperClass?: Ref<string>, withAppConfig?: boolean) => {
    ui: import("vue").ComputedRef<T>;
    attrs: import("vue").ComputedRef<Pick<{
        [x: string]: unknown;
    }, string>>;
};
