import type { PropType } from 'vue';
import type { DeepPartial, Strategy } from '../../types/index.js';
import type colors from '#ui-colors';
declare const config: {
    wrapper: string;
    container: string;
    base: string;
    form: string;
    rounded: string;
    color: string;
    background: string;
    border: string;
    ring: string;
    inner: string;
    label: string;
    required: string;
    help: string;
    default: {
        color: string;
    };
};
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
    id: {
        type: StringConstructor;
        default: () => any;
    };
    value: {
        type: (StringConstructor | NumberConstructor | BooleanConstructor | ObjectConstructor)[];
        default: any;
    };
    modelValue: {
        type: PropType<boolean | Array<any> | null>;
        default: any;
    };
    name: {
        type: StringConstructor;
        default: any;
    };
    disabled: {
        type: BooleanConstructor;
        default: boolean;
    };
    indeterminate: {
        type: BooleanConstructor;
        default: any;
    };
    help: {
        type: StringConstructor;
        default: any;
    };
    label: {
        type: StringConstructor;
        default: any;
    };
    required: {
        type: BooleanConstructor;
        default: boolean;
    };
    color: {
        type: PropType<(typeof colors)[number]>;
        default: () => string;
        validator(value: string): any;
    };
    inputClass: {
        type: StringConstructor;
        default: string;
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
        wrapper: string;
        container: string;
        base: string;
        form: string;
        rounded: string;
        color: string;
        background: string;
        border: string;
        ring: string;
        inner: string;
        label: string;
        required: string;
        help: string;
        default: {
            color: string;
        };
    }>;
    attrs: import("vue").ComputedRef<Pick<{
        [x: string]: unknown;
    }, string>>;
    toggle: import("vue").WritableComputedRef<boolean | any[], boolean | any[]>;
    inputId: string;
    name: import("vue").ComputedRef<string>;
    inputClass: import("vue").ComputedRef<string>;
    onChange: (event: Event) => void;
}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, ("change" | "update:modelValue")[], "change" | "update:modelValue", import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    id: {
        type: StringConstructor;
        default: () => any;
    };
    value: {
        type: (StringConstructor | NumberConstructor | BooleanConstructor | ObjectConstructor)[];
        default: any;
    };
    modelValue: {
        type: PropType<boolean | Array<any> | null>;
        default: any;
    };
    name: {
        type: StringConstructor;
        default: any;
    };
    disabled: {
        type: BooleanConstructor;
        default: boolean;
    };
    indeterminate: {
        type: BooleanConstructor;
        default: any;
    };
    help: {
        type: StringConstructor;
        default: any;
    };
    label: {
        type: StringConstructor;
        default: any;
    };
    required: {
        type: BooleanConstructor;
        default: boolean;
    };
    color: {
        type: PropType<(typeof colors)[number]>;
        default: () => string;
        validator(value: string): any;
    };
    inputClass: {
        type: StringConstructor;
        default: string;
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
}>> & Readonly<{
    onChange?: (...args: any[]) => any;
    "onUpdate:modelValue"?: (...args: any[]) => any;
}>, {
    name: string;
    required: boolean;
    class: any;
    ui: {
        wrapper?: string;
        container?: string;
        base?: string;
        form?: string;
        rounded?: string;
        color?: string;
        background?: string;
        border?: string;
        ring?: string;
        inner?: string;
        label?: string;
        required?: string;
        help?: string;
        default?: DeepPartial<{
            color: string;
        }, any>;
    } & {
        [key: string]: any;
    } & {
        strategy?: Strategy;
    };
    color: "primary" | "red" | "orange" | "amber" | "yellow" | "lime" | "green" | "emerald" | "teal" | "cyan" | "sky" | "blue" | "indigo" | "violet" | "purple" | "fuchsia" | "pink" | "rose";
    disabled: boolean;
    label: string;
    id: string;
    modelValue: boolean | any[];
    value: string | number | boolean | Record<string, any>;
    help: string;
    indeterminate: boolean;
    inputClass: string;
}, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
export default _default;
