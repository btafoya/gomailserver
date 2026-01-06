import type { PropType } from 'vue';
import type { DeepPartial, Strategy } from '../../types/index.js';
declare const config: {
    base: string;
    background: string;
    divide: string;
    ring: string;
    rounded: string;
    shadow: string;
    body: {
        base: string;
        background: string;
        padding: string;
    };
    header: {
        base: string;
        background: string;
        padding: string;
    };
    footer: {
        base: string;
        background: string;
        padding: string;
    };
};
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
    as: {
        type: StringConstructor;
        default: string;
    };
    class: {
        type: PropType<any>;
        default: () => string;
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
        background: string;
        divide: string;
        ring: string;
        rounded: string;
        shadow: string;
        body: {
            base: string;
            background: string;
            padding: string;
        };
        header: {
            base: string;
            background: string;
            padding: string;
        };
        footer: {
            base: string;
            background: string;
            padding: string;
        };
    }>;
    attrs: import("vue").ComputedRef<Pick<{
        [x: string]: unknown;
    }, string>>;
    cardClass: import("vue").ComputedRef<string>;
}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    as: {
        type: StringConstructor;
        default: string;
    };
    class: {
        type: PropType<any>;
        default: () => string;
    };
    ui: {
        type: PropType<DeepPartial<typeof config> & {
            strategy?: Strategy;
        }>;
        default: () => {};
    };
}>> & Readonly<{}>, {
    class: any;
    ui: {
        base?: string;
        background?: string;
        divide?: string;
        ring?: string;
        rounded?: string;
        shadow?: string;
        body?: DeepPartial<{
            base: string;
            background: string;
            padding: string;
        }, any>;
        header?: DeepPartial<{
            base: string;
            background: string;
            padding: string;
        }, any>;
        footer?: DeepPartial<{
            base: string;
            background: string;
            padding: string;
        }, any>;
    } & {
        [key: string]: any;
    } & {
        strategy?: Strategy;
    };
    as: string;
}, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
export default _default;
