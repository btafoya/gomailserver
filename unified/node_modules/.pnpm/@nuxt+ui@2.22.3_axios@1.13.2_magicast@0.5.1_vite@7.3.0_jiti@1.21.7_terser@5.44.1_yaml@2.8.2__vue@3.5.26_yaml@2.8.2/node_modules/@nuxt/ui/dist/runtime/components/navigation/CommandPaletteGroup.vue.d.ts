import type { PropType } from 'vue';
import type { Command, Group } from '../../types/index.js';
import type { commandPalette } from '#ui/ui.config';
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
    group: {
        type: PropType<Group>;
        required: true;
    };
    query: {
        type: StringConstructor;
        default: string;
    };
    groupAttribute: {
        type: StringConstructor;
        required: true;
    };
    commandAttribute: {
        type: StringConstructor;
        required: true;
    };
    selectedIcon: {
        type: StringConstructor;
        required: true;
    };
    ui: {
        type: PropType<typeof commandPalette>;
        required: true;
    };
}>, {
    label: import("vue").ComputedRef<any>;
    highlight: (text: string, { indices, value }: Command["matches"][number]) => string;
}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    group: {
        type: PropType<Group>;
        required: true;
    };
    query: {
        type: StringConstructor;
        default: string;
    };
    groupAttribute: {
        type: StringConstructor;
        required: true;
    };
    commandAttribute: {
        type: StringConstructor;
        required: true;
    };
    selectedIcon: {
        type: StringConstructor;
        required: true;
    };
    ui: {
        type: PropType<typeof commandPalette>;
        required: true;
    };
}>> & Readonly<{}>, {
    query: string;
}, {}, {
    HComboboxOption: import("vue").DefineComponent<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
        value: {
            type: PropType<string | number | boolean | object | null>;
        };
        disabled: {
            type: BooleanConstructor;
            default: boolean;
        };
        order: {
            type: NumberConstructor[];
            default: null;
        };
    }, () => import("vue").VNode<import("vue").RendererNode, import("vue").RendererElement, {
        [key: string]: any;
    }> | import("vue").VNode<import("vue").RendererNode, import("vue").RendererElement, {
        [key: string]: any;
    }>[] | null, unknown, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, Record<string, any>, string, import("vue").VNodeProps & import("vue").AllowedComponentProps & import("vue").ComponentCustomProps, Readonly<import("vue").ExtractPropTypes<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
        value: {
            type: PropType<string | number | boolean | object | null>;
        };
        disabled: {
            type: BooleanConstructor;
            default: boolean;
        };
        order: {
            type: NumberConstructor[];
            default: null;
        };
    }>>, {
        as: string | Record<string, any>;
        disabled: boolean;
        order: number;
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
    UAvatar: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
        as: {
            type: (StringConstructor | ObjectConstructor)[];
            default: string;
        };
        src: {
            type: (StringConstructor | BooleanConstructor)[];
            default: any;
        };
        alt: {
            type: StringConstructor;
            default: any;
        };
        text: {
            type: StringConstructor;
            default: any;
        };
        icon: {
            type: StringConstructor;
            default: () => any;
        };
        size: {
            type: PropType<import("../../types/avatar.js").AvatarSize>;
            default: () => string;
            validator(value: string): boolean;
        };
        chipColor: {
            type: PropType<import("../../types/avatar.js").AvatarChipColor>;
            default: () => any;
            validator(value: string): boolean;
        };
        chipPosition: {
            type: PropType<import("../../types/avatar.js").AvatarChipPosition>;
            default: () => string;
            validator(value: string): boolean;
        };
        chipText: {
            type: (StringConstructor | NumberConstructor)[];
            default: any;
        };
        imgClass: {
            type: StringConstructor;
            default: string;
        };
        class: {
            type: PropType<any>;
            default: () => "";
        };
        ui: {
            type: PropType<import("../../types/utils.js").DeepPartial<{
                wrapper: string;
                background: string;
                rounded: string;
                text: string;
                placeholder: string;
                size: {
                    '3xs': string;
                    '2xs': string;
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                    xl: string;
                    '2xl': string;
                    '3xl': string;
                };
                chip: {
                    base: string;
                    background: string;
                    position: {
                        'top-right': string;
                        'bottom-right': string;
                        'top-left': string;
                        'bottom-left': string;
                    };
                    size: {
                        '3xs': string;
                        '2xs': string;
                        xs: string;
                        sm: string;
                        md: string;
                        lg: string;
                        xl: string;
                        '2xl': string;
                        '3xl': string;
                    };
                };
                icon: {
                    base: string;
                    size: {
                        '3xs': string;
                        '2xs': string;
                        xs: string;
                        sm: string;
                        md: string;
                        lg: string;
                        xl: string;
                        '2xl': string;
                        '3xl': string;
                    };
                };
                default: {
                    size: string;
                    icon: any;
                    chipColor: any;
                    chipPosition: string;
                };
            }> & {
                strategy?: import("../../types/utils.js").Strategy;
            }>;
            default: () => {};
        };
    }>, {
        ui: import("vue").ComputedRef<{
            wrapper: string;
            background: string;
            rounded: string;
            text: string;
            placeholder: string;
            size: {
                '3xs': string;
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
                '2xl': string;
                '3xl': string;
            };
            chip: {
                base: string;
                background: string;
                position: {
                    'top-right': string;
                    'bottom-right': string;
                    'top-left': string;
                    'bottom-left': string;
                };
                size: {
                    '3xs': string;
                    '2xs': string;
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                    xl: string;
                    '2xl': string;
                    '3xl': string;
                };
            };
            icon: {
                base: string;
                size: {
                    '3xs': string;
                    '2xs': string;
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                    xl: string;
                    '2xl': string;
                    '3xl': string;
                };
            };
            default: {
                size: string;
                icon: any;
                chipColor: any;
                chipPosition: string;
            };
        }>;
        attrs: import("vue").ComputedRef<Pick<{
            [x: string]: unknown;
        }, string>>;
        wrapperClass: import("vue").ComputedRef<string>;
        imgClass: import("vue").ComputedRef<string>;
        iconClass: import("vue").ComputedRef<string>;
        chipClass: import("vue").ComputedRef<string>;
        url: import("vue").ComputedRef<string>;
        placeholder: import("vue").ComputedRef<string>;
        error: import("vue").Ref<boolean, boolean>;
        onError: () => void;
    }, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
        as: {
            type: (StringConstructor | ObjectConstructor)[];
            default: string;
        };
        src: {
            type: (StringConstructor | BooleanConstructor)[];
            default: any;
        };
        alt: {
            type: StringConstructor;
            default: any;
        };
        text: {
            type: StringConstructor;
            default: any;
        };
        icon: {
            type: StringConstructor;
            default: () => any;
        };
        size: {
            type: PropType<import("../../types/avatar.js").AvatarSize>;
            default: () => string;
            validator(value: string): boolean;
        };
        chipColor: {
            type: PropType<import("../../types/avatar.js").AvatarChipColor>;
            default: () => any;
            validator(value: string): boolean;
        };
        chipPosition: {
            type: PropType<import("../../types/avatar.js").AvatarChipPosition>;
            default: () => string;
            validator(value: string): boolean;
        };
        chipText: {
            type: (StringConstructor | NumberConstructor)[];
            default: any;
        };
        imgClass: {
            type: StringConstructor;
            default: string;
        };
        class: {
            type: PropType<any>;
            default: () => "";
        };
        ui: {
            type: PropType<import("../../types/utils.js").DeepPartial<{
                wrapper: string;
                background: string;
                rounded: string;
                text: string;
                placeholder: string;
                size: {
                    '3xs': string;
                    '2xs': string;
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                    xl: string;
                    '2xl': string;
                    '3xl': string;
                };
                chip: {
                    base: string;
                    background: string;
                    position: {
                        'top-right': string;
                        'bottom-right': string;
                        'top-left': string;
                        'bottom-left': string;
                    };
                    size: {
                        '3xs': string;
                        '2xs': string;
                        xs: string;
                        sm: string;
                        md: string;
                        lg: string;
                        xl: string;
                        '2xl': string;
                        '3xl': string;
                    };
                };
                icon: {
                    base: string;
                    size: {
                        '3xs': string;
                        '2xs': string;
                        xs: string;
                        sm: string;
                        md: string;
                        lg: string;
                        xl: string;
                        '2xl': string;
                        '3xl': string;
                    };
                };
                default: {
                    size: string;
                    icon: any;
                    chipColor: any;
                    chipPosition: string;
                };
            }> & {
                strategy?: import("../../types/utils.js").Strategy;
            }>;
            default: () => {};
        };
    }>> & Readonly<{}>, {
        size: import("../../types/avatar.js").AvatarSize;
        class: any;
        ui: {
            wrapper?: string;
            background?: string;
            rounded?: string;
            text?: string;
            placeholder?: string;
            size?: import("../../types/utils.js").DeepPartial<{
                '3xs': string;
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
                '2xl': string;
                '3xl': string;
            }, any>;
            chip?: import("../../types/utils.js").DeepPartial<{
                base: string;
                background: string;
                position: {
                    'top-right': string;
                    'bottom-right': string;
                    'top-left': string;
                    'bottom-left': string;
                };
                size: {
                    '3xs': string;
                    '2xs': string;
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                    xl: string;
                    '2xl': string;
                    '3xl': string;
                };
            }, any>;
            icon?: import("../../types/utils.js").DeepPartial<{
                base: string;
                size: {
                    '3xs': string;
                    '2xs': string;
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                    xl: string;
                    '2xl': string;
                    '3xl': string;
                };
            }, any>;
            default?: import("../../types/utils.js").DeepPartial<{
                size: string;
                icon: any;
                chipColor: any;
                chipPosition: string;
            }, any>;
        } & {
            [key: string]: any;
        } & {
            strategy?: import("../../types/utils.js").Strategy;
        };
        as: string | Record<string, any>;
        text: string;
        alt: string;
        icon: string;
        src: string | boolean;
        chipColor: import("../../types/avatar.js").AvatarChipColor;
        chipPosition: "top-right" | "bottom-right" | "top-left" | "bottom-left";
        chipText: string | number;
        imgClass: string;
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
    UKbd: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
        value: {
            type: StringConstructor;
            default: any;
        };
        size: {
            type: PropType<import("../../types/kbd.js").KbdSize>;
            default: () => string;
            validator(value: string): boolean;
        };
        class: {
            type: PropType<any>;
            default: () => "";
        };
        ui: {
            type: PropType<import("../../types/utils.js").DeepPartial<{
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
            }> & {
                strategy?: import("../../types/utils.js").Strategy;
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
            type: PropType<import("../../types/kbd.js").KbdSize>;
            default: () => string;
            validator(value: string): boolean;
        };
        class: {
            type: PropType<any>;
            default: () => "";
        };
        ui: {
            type: PropType<import("../../types/utils.js").DeepPartial<{
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
            }> & {
                strategy?: import("../../types/utils.js").Strategy;
            }>;
            default: () => {};
        };
    }>> & Readonly<{}>, {
        size: import("../../types/kbd.js").KbdSize;
        class: any;
        ui: {
            base?: string;
            padding?: string;
            size?: import("../../types/utils.js").DeepPartial<{
                xs: string;
                sm: string;
                md: string;
            }, any>;
            rounded?: string;
            font?: string;
            background?: string;
            ring?: string;
            default?: import("../../types/utils.js").DeepPartial<{
                size: string;
            }, any>;
        } & {
            [key: string]: any;
        } & {
            strategy?: import("../../types/utils.js").Strategy;
        };
        value: string;
    }, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
export default _default;
