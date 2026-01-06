import type { PropType } from 'vue';
import type { ToggleSize, ToggleColor, Strategy, DeepPartial } from '../../types/index.js';
declare const config: {
    base: string;
    rounded: string;
    ring: string;
    active: string;
    inactive: string;
    size: {
        '2xs': string;
        xs: string;
        sm: string;
        md: string;
        lg: string;
        xl: string;
        '2xl': string;
    };
    container: {
        base: string;
        active: {
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
            '2xl': string;
        };
        inactive: string;
        size: {
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
            '2xl': string;
        };
    };
    icon: {
        base: string;
        active: string;
        inactive: string;
        size: {
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
            '2xl': string;
        };
        on: string;
        off: string;
        loading: string;
    };
    default: {
        onIcon: any;
        offIcon: any;
        loadingIcon: string;
        color: string;
        size: string;
    };
};
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
    id: {
        type: StringConstructor;
        default: any;
    };
    name: {
        type: StringConstructor;
        default: any;
    };
    modelValue: {
        type: PropType<boolean | null>;
        default: boolean;
    };
    disabled: {
        type: BooleanConstructor;
        default: boolean;
    };
    loading: {
        type: BooleanConstructor;
        default: boolean;
    };
    onIcon: {
        type: StringConstructor;
        default: () => any;
    };
    offIcon: {
        type: StringConstructor;
        default: () => any;
    };
    loadingIcon: {
        type: StringConstructor;
        default: () => string;
    };
    color: {
        type: PropType<ToggleColor>;
        default: () => string;
        validator(value: string): any;
    };
    size: {
        type: PropType<ToggleSize>;
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
        rounded: string;
        ring: string;
        active: string;
        inactive: string;
        size: {
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
            '2xl': string;
        };
        container: {
            base: string;
            active: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
                '2xl': string;
            };
            inactive: string;
            size: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
                '2xl': string;
            };
        };
        icon: {
            base: string;
            active: string;
            inactive: string;
            size: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
                '2xl': string;
            };
            on: string;
            off: string;
            loading: string;
        };
        default: {
            onIcon: any;
            offIcon: any;
            loadingIcon: string;
            color: string;
            size: string;
        };
    }>;
    attrs: import("vue").ComputedRef<Pick<{
        [x: string]: unknown;
    }, string>>;
    name: import("vue").ComputedRef<string>;
    inputId: import("vue").ComputedRef<string>;
    active: import("vue").WritableComputedRef<boolean, boolean>;
    switchClass: import("vue").ComputedRef<string>;
    containerClass: import("vue").ComputedRef<string>;
    onIconClass: import("vue").ComputedRef<string>;
    offIconClass: import("vue").ComputedRef<string>;
    loadingIconClass: import("vue").ComputedRef<string>;
}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, ("change" | "update:modelValue")[], "change" | "update:modelValue", import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    id: {
        type: StringConstructor;
        default: any;
    };
    name: {
        type: StringConstructor;
        default: any;
    };
    modelValue: {
        type: PropType<boolean | null>;
        default: boolean;
    };
    disabled: {
        type: BooleanConstructor;
        default: boolean;
    };
    loading: {
        type: BooleanConstructor;
        default: boolean;
    };
    onIcon: {
        type: StringConstructor;
        default: () => any;
    };
    offIcon: {
        type: StringConstructor;
        default: () => any;
    };
    loadingIcon: {
        type: StringConstructor;
        default: () => string;
    };
    color: {
        type: PropType<ToggleColor>;
        default: () => string;
        validator(value: string): any;
    };
    size: {
        type: PropType<ToggleSize>;
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
}>> & Readonly<{
    onChange?: (...args: any[]) => any;
    "onUpdate:modelValue"?: (...args: any[]) => any;
}>, {
    name: string;
    size: ToggleSize;
    class: any;
    ui: {
        base?: string;
        rounded?: string;
        ring?: string;
        active?: string;
        inactive?: string;
        size?: DeepPartial<{
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
            '2xl': string;
        }, any>;
        container?: DeepPartial<{
            base: string;
            active: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
                '2xl': string;
            };
            inactive: string;
            size: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
                '2xl': string;
            };
        }, any>;
        icon?: DeepPartial<{
            base: string;
            active: string;
            inactive: string;
            size: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
                '2xl': string;
            };
            on: string;
            off: string;
            loading: string;
        }, any>;
        default?: DeepPartial<{
            onIcon: any;
            offIcon: any;
            loadingIcon: string;
            color: string;
            size: string;
        }, any>;
    } & {
        [key: string]: any;
    } & {
        strategy?: Strategy;
    };
    color: "primary" | "red" | "orange" | "amber" | "yellow" | "lime" | "green" | "emerald" | "teal" | "cyan" | "sky" | "blue" | "indigo" | "violet" | "purple" | "fuchsia" | "pink" | "rose";
    disabled: boolean;
    id: string;
    modelValue: boolean;
    loading: boolean;
    loadingIcon: string;
    onIcon: string;
    offIcon: string;
}, {}, {
    HSwitch: import("vue").DefineComponent<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
        modelValue: {
            type: BooleanConstructor;
            default: undefined;
        };
        defaultChecked: {
            type: BooleanConstructor;
            optional: boolean;
        };
        form: {
            type: StringConstructor;
            optional: boolean;
        };
        name: {
            type: StringConstructor;
            optional: boolean;
        };
        value: {
            type: StringConstructor;
            optional: boolean;
        };
        id: {
            type: StringConstructor;
            default: null;
        };
        disabled: {
            type: BooleanConstructor;
            default: boolean;
        };
        tabIndex: {
            type: NumberConstructor;
            default: number;
        };
    }, () => import("vue").VNode<import("vue").RendererNode, import("vue").RendererElement, {
        [key: string]: any;
    }>, unknown, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {
        'update:modelValue': (_value: boolean) => true;
    }, string, import("vue").VNodeProps & import("vue").AllowedComponentProps & import("vue").ComponentCustomProps, Readonly<import("vue").ExtractPropTypes<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
        modelValue: {
            type: BooleanConstructor;
            default: undefined;
        };
        defaultChecked: {
            type: BooleanConstructor;
            optional: boolean;
        };
        form: {
            type: StringConstructor;
            optional: boolean;
        };
        name: {
            type: StringConstructor;
            optional: boolean;
        };
        value: {
            type: StringConstructor;
            optional: boolean;
        };
        id: {
            type: StringConstructor;
            default: null;
        };
        disabled: {
            type: BooleanConstructor;
            default: boolean;
        };
        tabIndex: {
            type: NumberConstructor;
            default: number;
        };
    }>> & {
        "onUpdate:modelValue"?: ((_value: boolean) => any) | undefined;
    }, {
        id: string;
        tabIndex: number;
        as: string | Record<string, any>;
        disabled: boolean;
        modelValue: boolean;
        defaultChecked: boolean;
    }, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
    UIcon: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
        name: {
            type: StringConstructor;
            required: true;
        };
        mode: {
            type: PropType<"svg" | "css">;
            required: false;
            default: any;
        };
        size: {
            type: (StringConstructor | NumberConstructor)[];
            required: false;
            default: any;
        };
        customize: {
            type: FunctionConstructor;
            required: false;
            default: any;
        };
    }>, {}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
        name: {
            type: StringConstructor;
            required: true;
        };
        mode: {
            type: PropType<"svg" | "css">;
            required: false;
            default: any;
        };
        size: {
            type: (StringConstructor | NumberConstructor)[];
            required: false;
            default: any;
        };
        customize: {
            type: FunctionConstructor;
            required: false;
            default: any;
        };
    }>> & Readonly<{}>, {
        mode: "svg" | "css";
        size: string | number;
        customize: Function;
    }, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
export default _default;
