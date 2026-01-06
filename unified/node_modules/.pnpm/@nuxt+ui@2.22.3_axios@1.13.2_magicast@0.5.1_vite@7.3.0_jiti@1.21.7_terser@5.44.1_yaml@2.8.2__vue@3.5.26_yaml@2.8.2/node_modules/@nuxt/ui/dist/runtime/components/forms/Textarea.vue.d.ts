import type { PropType } from 'vue';
import type { TextareaSize, TextareaColor, TextareaVariant, Strategy, DeepPartial } from '../../types/index.js';
declare const config: {
    form: string;
    default: {
        size: string;
        color: string;
        variant: string;
    };
    wrapper: string;
    base: string;
    rounded: string;
    placeholder: string;
    file: {
        base: string;
    };
    size: {
        '2xs': string;
        xs: string;
        sm: string;
        md: string;
        lg: string;
        xl: string;
    };
    gap: {
        '2xs': string;
        xs: string;
        sm: string;
        md: string;
        lg: string;
        xl: string;
    };
    padding: {
        '2xs': string;
        xs: string;
        sm: string;
        md: string;
        lg: string;
        xl: string;
    };
    leading: {
        padding: {
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
        };
    };
    trailing: {
        padding: {
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
        };
    };
    color: {
        white: {
            outline: string;
        };
        gray: {
            outline: string;
        };
    };
    variant: {
        outline: string;
        none: string;
    };
    icon: {
        base: string;
        color: string;
        loading: string;
        size: {
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
        };
        leading: {
            wrapper: string;
            pointer: string;
            padding: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
            };
        };
        trailing: {
            wrapper: string;
            pointer: string;
            padding: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
            };
        };
    };
};
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
    modelValue: {
        type: PropType<string | number | null>;
        default: string;
    };
    id: {
        type: StringConstructor;
        default: any;
    };
    name: {
        type: StringConstructor;
        default: any;
    };
    placeholder: {
        type: StringConstructor;
        default: any;
    };
    required: {
        type: BooleanConstructor;
        default: boolean;
    };
    disabled: {
        type: BooleanConstructor;
        default: boolean;
    };
    rows: {
        type: NumberConstructor;
        default: number;
    };
    maxrows: {
        type: NumberConstructor;
        default: number;
    };
    autoresize: {
        type: BooleanConstructor;
        default: boolean;
    };
    autofocus: {
        type: BooleanConstructor;
        default: boolean;
    };
    autofocusDelay: {
        type: NumberConstructor;
        default: number;
    };
    resize: {
        type: BooleanConstructor;
        default: boolean;
    };
    padded: {
        type: BooleanConstructor;
        default: boolean;
    };
    size: {
        type: PropType<TextareaSize>;
        default: any;
        validator(value: string): boolean;
    };
    color: {
        type: PropType<TextareaColor>;
        default: () => string;
        validator(value: string): boolean;
    };
    variant: {
        type: PropType<TextareaVariant>;
        default: () => string;
        validator(value: string): boolean;
    };
    textareaClass: {
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
    modelModifiers: {
        type: PropType<{
            trim?: boolean;
            lazy?: boolean;
            number?: boolean;
            nullify?: boolean;
        }>;
        default: () => {};
    };
}>, {
    ui: import("vue").ComputedRef<{
        form: string;
        default: {
            size: string;
            color: string;
            variant: string;
        };
        wrapper: string;
        base: string;
        rounded: string;
        placeholder: string;
        file: {
            base: string;
        };
        size: {
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
        };
        gap: {
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
        };
        padding: {
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
        };
        leading: {
            padding: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
            };
        };
        trailing: {
            padding: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
            };
        };
        color: {
            white: {
                outline: string;
            };
            gray: {
                outline: string;
            };
        };
        variant: {
            outline: string;
            none: string;
        };
        icon: {
            base: string;
            color: string;
            loading: string;
            size: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
            };
            leading: {
                wrapper: string;
                pointer: string;
                padding: {
                    '2xs': string;
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                    xl: string;
                };
            };
            trailing: {
                wrapper: string;
                pointer: string;
                padding: {
                    '2xs': string;
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                    xl: string;
                };
            };
        };
    }>;
    attrs: import("vue").ComputedRef<Pick<{
        [x: string]: unknown;
    }, string>>;
    name: import("vue").ComputedRef<string>;
    inputId: import("vue").ComputedRef<string>;
    textarea: import("vue").Ref<HTMLTextAreaElement, HTMLTextAreaElement>;
    textareaClass: import("vue").ComputedRef<string>;
    onInput: (event: Event) => void;
    onChange: (event: Event) => void;
    onBlur: (event: FocusEvent) => void;
}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, ("blur" | "change" | "update:modelValue")[], "blur" | "change" | "update:modelValue", import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    modelValue: {
        type: PropType<string | number | null>;
        default: string;
    };
    id: {
        type: StringConstructor;
        default: any;
    };
    name: {
        type: StringConstructor;
        default: any;
    };
    placeholder: {
        type: StringConstructor;
        default: any;
    };
    required: {
        type: BooleanConstructor;
        default: boolean;
    };
    disabled: {
        type: BooleanConstructor;
        default: boolean;
    };
    rows: {
        type: NumberConstructor;
        default: number;
    };
    maxrows: {
        type: NumberConstructor;
        default: number;
    };
    autoresize: {
        type: BooleanConstructor;
        default: boolean;
    };
    autofocus: {
        type: BooleanConstructor;
        default: boolean;
    };
    autofocusDelay: {
        type: NumberConstructor;
        default: number;
    };
    resize: {
        type: BooleanConstructor;
        default: boolean;
    };
    padded: {
        type: BooleanConstructor;
        default: boolean;
    };
    size: {
        type: PropType<TextareaSize>;
        default: any;
        validator(value: string): boolean;
    };
    color: {
        type: PropType<TextareaColor>;
        default: () => string;
        validator(value: string): boolean;
    };
    variant: {
        type: PropType<TextareaVariant>;
        default: () => string;
        validator(value: string): boolean;
    };
    textareaClass: {
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
    modelModifiers: {
        type: PropType<{
            trim?: boolean;
            lazy?: boolean;
            number?: boolean;
            nullify?: boolean;
        }>;
        default: () => {};
    };
}>> & Readonly<{
    onBlur?: (...args: any[]) => any;
    onChange?: (...args: any[]) => any;
    "onUpdate:modelValue"?: (...args: any[]) => any;
}>, {
    name: string;
    required: boolean;
    size: TextareaSize;
    class: any;
    placeholder: string;
    ui: {
        form?: string;
        default?: DeepPartial<{
            size: string;
            color: string;
            variant: string;
        }, any>;
        wrapper?: string;
        base?: string;
        rounded?: string;
        placeholder?: string;
        file?: DeepPartial<{
            base: string;
        }, any>;
        size?: DeepPartial<{
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
        }, any>;
        gap?: DeepPartial<{
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
        }, any>;
        padding?: DeepPartial<{
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
        }, any>;
        leading?: DeepPartial<{
            padding: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
            };
        }, any>;
        trailing?: DeepPartial<{
            padding: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
            };
        }, any>;
        color?: DeepPartial<{
            white: {
                outline: string;
            };
            gray: {
                outline: string;
            };
        }, any>;
        variant?: DeepPartial<{
            outline: string;
            none: string;
        }, any>;
        icon?: DeepPartial<{
            base: string;
            color: string;
            loading: string;
            size: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
            };
            leading: {
                wrapper: string;
                pointer: string;
                padding: {
                    '2xs': string;
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                    xl: string;
                };
            };
            trailing: {
                wrapper: string;
                pointer: string;
                padding: {
                    '2xs': string;
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                    xl: string;
                };
            };
        }, any>;
    } & {
        [key: string]: any;
    } & {
        strategy?: Strategy;
    };
    color: TextareaColor;
    variant: TextareaVariant;
    resize: boolean;
    disabled: boolean;
    id: string;
    modelValue: string | number;
    padded: boolean;
    rows: number;
    autofocus: boolean;
    autofocusDelay: number;
    modelModifiers: {
        trim?: boolean;
        lazy?: boolean;
        number?: boolean;
        nullify?: boolean;
    };
    maxrows: number;
    autoresize: boolean;
    textareaClass: string;
}, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
export default _default;
