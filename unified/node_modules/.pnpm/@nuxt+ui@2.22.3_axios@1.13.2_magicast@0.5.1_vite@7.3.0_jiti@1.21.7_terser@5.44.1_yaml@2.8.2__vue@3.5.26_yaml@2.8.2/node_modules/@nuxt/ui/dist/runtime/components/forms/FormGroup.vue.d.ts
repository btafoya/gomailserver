import type { Ref, PropType } from 'vue';
import type { FormGroupSize, Strategy, DeepPartial } from '../../types/index.js';
declare const config: {
    wrapper: string;
    inner: string;
    label: {
        wrapper: string;
        base: string;
        required: string;
    };
    size: {
        '2xs': string;
        xs: string;
        sm: string;
        md: string;
        lg: string;
        xl: string;
    };
    container: string;
    description: string;
    hint: string;
    help: string;
    error: string;
    default: {
        size: string;
    };
};
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
    name: {
        type: StringConstructor;
        default: any;
    };
    size: {
        type: PropType<FormGroupSize>;
        default: any;
        validator(value: string): boolean;
    };
    label: {
        type: StringConstructor;
        default: any;
    };
    description: {
        type: StringConstructor;
        default: any;
    };
    required: {
        type: BooleanConstructor;
        default: boolean;
    };
    help: {
        type: StringConstructor;
        default: any;
    };
    error: {
        type: (StringConstructor | BooleanConstructor)[];
        default: any;
    };
    hint: {
        type: StringConstructor;
        default: any;
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
    eagerValidation: {
        type: BooleanConstructor;
        default: boolean;
    };
}>, {
    ui: import("vue").ComputedRef<{
        wrapper: string;
        inner: string;
        label: {
            wrapper: string;
            base: string;
            required: string;
        };
        size: {
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
        };
        container: string;
        description: string;
        hint: string;
        help: string;
        error: string;
        default: {
            size: string;
        };
    }>;
    attrs: import("vue").ComputedRef<Pick<{
        [x: string]: unknown;
    }, string>>;
    inputId: Ref<string, string>;
    size: import("vue").ComputedRef<any>;
    error: import("vue").ComputedRef<string | boolean>;
}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    name: {
        type: StringConstructor;
        default: any;
    };
    size: {
        type: PropType<FormGroupSize>;
        default: any;
        validator(value: string): boolean;
    };
    label: {
        type: StringConstructor;
        default: any;
    };
    description: {
        type: StringConstructor;
        default: any;
    };
    required: {
        type: BooleanConstructor;
        default: boolean;
    };
    help: {
        type: StringConstructor;
        default: any;
    };
    error: {
        type: (StringConstructor | BooleanConstructor)[];
        default: any;
    };
    hint: {
        type: StringConstructor;
        default: any;
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
    eagerValidation: {
        type: BooleanConstructor;
        default: boolean;
    };
}>> & Readonly<{}>, {
    name: string;
    required: boolean;
    size: FormGroupSize;
    class: any;
    ui: {
        wrapper?: string;
        inner?: string;
        label?: DeepPartial<{
            wrapper: string;
            base: string;
            required: string;
        }, any>;
        size?: DeepPartial<{
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
        }, any>;
        container?: string;
        description?: string;
        hint?: string;
        help?: string;
        error?: string;
        default?: DeepPartial<{
            size: string;
        }, any>;
    } & {
        [key: string]: any;
    } & {
        strategy?: Strategy;
    };
    label: string;
    error: string | boolean;
    description: string;
    help: string;
    hint: string;
    eagerValidation: boolean;
}, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
export default _default;
