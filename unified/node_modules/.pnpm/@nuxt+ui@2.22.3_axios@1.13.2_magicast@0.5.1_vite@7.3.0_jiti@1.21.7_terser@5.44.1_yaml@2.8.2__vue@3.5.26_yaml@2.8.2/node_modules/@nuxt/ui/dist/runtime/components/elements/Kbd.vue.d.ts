import type { PropType } from 'vue';
import type { DeepPartial, KbdSize, Strategy } from '../../types/index.js';
declare const config: {
    base: string;
    padding: string;
    size: {
        xs: string;
        sm: string;
        md: string;
    };
    rounded: string;
    font: string;
    background: string;
    ring: string;
    default: {
        size: string;
    };
};
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
    value: {
        type: StringConstructor;
        default: any;
    };
    size: {
        type: PropType<KbdSize>;
        default: () => string;
        validator(value: string): boolean;
    };
    class: {
        type: PropType<any>;
        default: () => "";
    };
    ui: {
        type: PropType<DeepPartial<typeof config> & {
            strategy?: Strategy;
        }>;
        default: () => {};
    };
}>, {
    ui: import("vue").ComputedRef<{
        base: string;
        padding: string;
        size: {
            xs: string;
            sm: string;
            md: string;
        };
        rounded: string;
        font: string;
        background: string;
        ring: string;
        default: {
            size: string;
        };
    }>;
    attrs: import("vue").ComputedRef<Pick<{
        [x: string]: unknown;
    }, string>>;
    kbdClass: import("vue").ComputedRef<string>;
}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    value: {
        type: StringConstructor;
        default: any;
    };
    size: {
        type: PropType<KbdSize>;
        default: () => string;
        validator(value: string): boolean;
    };
    class: {
        type: PropType<any>;
        default: () => "";
    };
    ui: {
        type: PropType<DeepPartial<typeof config> & {
            strategy?: Strategy;
        }>;
        default: () => {};
    };
}>> & Readonly<{}>, {
    size: KbdSize;
    class: any;
    ui: {
        base?: string;
        padding?: string;
        size?: DeepPartial<{
            xs: string;
            sm: string;
            md: string;
        }, any>;
        rounded?: string;
        font?: string;
        background?: string;
        ring?: string;
        default?: DeepPartial<{
            size: string;
        }, any>;
    } & {
        [key: string]: any;
    } & {
        strategy?: Strategy;
    };
    value: string;
}, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
export default _default;
