import type { PropType } from 'vue';
import type { BadgeColor, BadgeSize, BadgeVariant, DeepPartial, Strategy } from '../../types/index.js';
declare const config: {
    base: string;
    rounded: string;
    font: string;
    size: {
        xs: string;
        sm: string;
        md: string;
        lg: string;
    };
    gap: {
        xs: string;
        sm: string;
        md: string;
        lg: string;
    };
    color: {
        white: {
            solid: string;
        };
        gray: {
            solid: string;
        };
        black: {
            solid: string;
        };
    };
    variant: {
        solid: string;
        outline: string;
        soft: string;
        subtle: string;
    };
    icon: {
        base: string;
        size: {
            xs: string;
            sm: string;
            md: string;
            lg: string;
        };
    };
    default: {
        size: string;
        variant: string;
        color: string;
    };
};
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
    size: {
        type: PropType<BadgeSize>;
        default: () => string;
        validator(value: string): boolean;
    };
    color: {
        type: PropType<BadgeColor>;
        default: () => string;
        validator(value: string): boolean;
    };
    variant: {
        type: PropType<BadgeVariant>;
        default: () => string;
        validator(value: string): boolean;
    };
    label: {
        type: (StringConstructor | NumberConstructor)[];
        default: any;
    };
    icon: {
        type: StringConstructor;
        default: any;
    };
    leadingIcon: {
        type: StringConstructor;
        default: any;
    };
    trailingIcon: {
        type: StringConstructor;
        default: any;
    };
    trailing: {
        type: BooleanConstructor;
        default: boolean;
    };
    leading: {
        type: BooleanConstructor;
        default: boolean;
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
    attrs: import("vue").ComputedRef<Pick<{
        [x: string]: unknown;
    }, string>>;
    isLeading: import("vue").ComputedRef<string | true>;
    isTrailing: import("vue").ComputedRef<string | true>;
    badgeClass: import("vue").ComputedRef<string>;
    leadingIconName: import("vue").ComputedRef<string>;
    trailingIconName: import("vue").ComputedRef<string>;
    leadingIconClass: import("vue").ComputedRef<string>;
    trailingIconClass: import("vue").ComputedRef<string>;
}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    size: {
        type: PropType<BadgeSize>;
        default: () => string;
        validator(value: string): boolean;
    };
    color: {
        type: PropType<BadgeColor>;
        default: () => string;
        validator(value: string): boolean;
    };
    variant: {
        type: PropType<BadgeVariant>;
        default: () => string;
        validator(value: string): boolean;
    };
    label: {
        type: (StringConstructor | NumberConstructor)[];
        default: any;
    };
    icon: {
        type: StringConstructor;
        default: any;
    };
    leadingIcon: {
        type: StringConstructor;
        default: any;
    };
    trailingIcon: {
        type: StringConstructor;
        default: any;
    };
    trailing: {
        type: BooleanConstructor;
        default: boolean;
    };
    leading: {
        type: BooleanConstructor;
        default: boolean;
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
    size: BadgeSize;
    class: any;
    ui: {
        base?: string;
        rounded?: string;
        font?: string;
        size?: DeepPartial<{
            xs: string;
            sm: string;
            md: string;
            lg: string;
        }, any>;
        gap?: DeepPartial<{
            xs: string;
            sm: string;
            md: string;
            lg: string;
        }, any>;
        color?: DeepPartial<{
            white: {
                solid: string;
            };
            gray: {
                solid: string;
            };
            black: {
                solid: string;
            };
        }, any>;
        variant?: DeepPartial<{
            solid: string;
            outline: string;
            soft: string;
            subtle: string;
        }, any>;
        icon?: DeepPartial<{
            base: string;
            size: {
                xs: string;
                sm: string;
                md: string;
                lg: string;
            };
        }, any>;
        default?: DeepPartial<{
            size: string;
            variant: string;
            color: string;
        }, any>;
    } & {
        [key: string]: any;
    } & {
        strategy?: Strategy;
    };
    color: BadgeColor;
    variant: BadgeVariant;
    leading: boolean;
    label: string | number;
    icon: string;
    leadingIcon: string;
    trailingIcon: string;
    trailing: boolean;
}, {}, {
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
