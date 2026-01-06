import type { PropType } from 'vue';
import type { SelectSize, SelectColor, SelectVariant, PopperOptions, Strategy, DeepPartial } from '../../types/index.js';
declare const config: {
    form: string;
    placeholder: string;
    default: {
        size: string;
        color: string;
        variant: string;
        loadingIcon: string;
        trailingIcon: string;
    };
    wrapper: string;
    base: string;
    rounded: string;
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
declare const configMenu: {
    select: string;
    input: string;
    required: string;
    label: string;
    option: {
        create: string;
        base: string;
        rounded: string;
        padding: string;
        size: string;
        color: string;
        container: string;
        active: string;
        inactive: string;
        selected: string;
        disabled: string;
        empty: string;
        icon: {
            base: string;
            active: string;
            inactive: string;
        };
        selectedIcon: {
            wrapper: string;
            padding: string;
            base: string;
        };
        avatar: {
            base: string;
            size: import("../../types/avatar.js").AvatarSize;
        };
        chip: {
            base: string;
        };
    };
    transition: {
        leaveActiveClass: string;
        leaveFromClass: string;
        leaveToClass: string;
    };
    popper: {
        placement: string;
    };
    default: {
        selectedIcon: string;
        clearSearchOnClose: boolean;
        showCreateOptionWhen: string;
        searchablePlaceholder: {
            label: string;
        };
        empty: {
            label: string;
        };
        optionEmpty: {
            label: string;
        };
    };
    arrow: {
        ring: string;
        background: string;
        base: string;
        rounded: string;
        shadow: string;
        placement: string;
    };
    container: string;
    trigger: string;
    width: string;
    height: string;
    base: string;
    background: string;
    shadow: string;
    rounded: string;
    padding: string;
    ring: string;
    empty: string;
};
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
    modelValue: {
        type: PropType<string | number | object | Array<any> | boolean | null>;
        default: string;
    };
    query: {
        type: StringConstructor;
        default: any;
    };
    by: {
        type: StringConstructor;
        default: any;
    };
    options: {
        type: PropType<{
            [key: string]: any;
            disabled?: boolean;
        }[] | string[]>;
        default: () => any[];
    };
    id: {
        type: StringConstructor;
        default: any;
    };
    name: {
        type: StringConstructor;
        default: any;
    };
    required: {
        type: BooleanConstructor;
        default: boolean;
    };
    icon: {
        type: StringConstructor;
        default: any;
    };
    loadingIcon: {
        type: StringConstructor;
        default: () => string;
    };
    leadingIcon: {
        type: StringConstructor;
        default: any;
    };
    trailingIcon: {
        type: StringConstructor;
        default: () => string;
    };
    trailing: {
        type: BooleanConstructor;
        default: boolean;
    };
    leading: {
        type: BooleanConstructor;
        default: boolean;
    };
    loading: {
        type: BooleanConstructor;
        default: boolean;
    };
    selectedIcon: {
        type: StringConstructor;
        default: () => string;
    };
    disabled: {
        type: BooleanConstructor;
        default: boolean;
    };
    multiple: {
        type: BooleanConstructor;
        default: boolean;
    };
    searchable: {
        type: PropType<boolean | ((query: string) => Promise<any[]> | any[])>;
        default: boolean;
    };
    searchablePlaceholder: {
        type: StringConstructor;
        default: () => string;
    };
    searchableLazy: {
        type: BooleanConstructor;
        default: boolean;
    };
    clearSearchOnClose: {
        type: BooleanConstructor;
        default: () => boolean;
    };
    debounce: {
        type: NumberConstructor;
        default: number;
    };
    creatable: {
        type: BooleanConstructor;
        default: boolean;
    };
    showCreateOptionWhen: {
        type: PropType<"always" | "empty" | ((query: string, results: any[]) => boolean)>;
        default: () => string;
    };
    placeholder: {
        type: StringConstructor;
        default: any;
    };
    padded: {
        type: BooleanConstructor;
        default: boolean;
    };
    size: {
        type: PropType<SelectSize>;
        default: any;
        validator(value: string): boolean;
    };
    color: {
        type: PropType<SelectColor>;
        default: () => string;
        validator(value: string): boolean;
    };
    variant: {
        type: PropType<SelectVariant>;
        default: () => string;
        validator(value: string): boolean;
    };
    optionAttribute: {
        type: StringConstructor;
        default: string;
    };
    valueAttribute: {
        type: StringConstructor;
        default: any;
    };
    searchAttributes: {
        type: ArrayConstructor;
        default: any;
    };
    inputTargetForm: {
        type: StringConstructor;
        default: any;
    };
    popper: {
        type: PropType<PopperOptions>;
        default: () => {};
    };
    selectClass: {
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
    uiMenu: {
        type: PropType<DeepPartial<typeof configMenu> & {
            strategy?: Strategy;
        }>;
        default: () => {};
    };
}>, {
    ui: import("vue").ComputedRef<{
        form: string;
        placeholder: string;
        default: {
            size: string;
            color: string;
            variant: string;
            loadingIcon: string;
            trailingIcon: string;
        };
        wrapper: string;
        base: string;
        rounded: string;
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
    uiMenu: import("vue").ComputedRef<{
        select: string;
        input: string;
        required: string;
        label: string;
        option: {
            create: string;
            base: string;
            rounded: string;
            padding: string;
            size: string;
            color: string;
            container: string;
            active: string;
            inactive: string;
            selected: string;
            disabled: string;
            empty: string;
            icon: {
                base: string;
                active: string;
                inactive: string;
            };
            selectedIcon: {
                wrapper: string;
                padding: string;
                base: string;
            };
            avatar: {
                base: string;
                size: import("../../types/avatar.js").AvatarSize;
            };
            chip: {
                base: string;
            };
        };
        transition: {
            leaveActiveClass: string;
            leaveFromClass: string;
            leaveToClass: string;
        };
        popper: {
            placement: string;
        };
        default: {
            selectedIcon: string;
            clearSearchOnClose: boolean;
            showCreateOptionWhen: string;
            searchablePlaceholder: {
                label: string;
            };
            empty: {
                label: string;
            };
            optionEmpty: {
                label: string;
            };
        };
        arrow: {
            ring: string;
            background: string;
            base: string;
            rounded: string;
            shadow: string;
            placement: string;
        };
        container: string;
        trigger: string;
        width: string;
        height: string;
        base: string;
        background: string;
        shadow: string;
        rounded: string;
        padding: string;
        ring: string;
        empty: string;
    }>;
    attrs: import("vue").ComputedRef<Pick<{
        [x: string]: unknown;
    }, string>>;
    name: import("vue").ComputedRef<string>;
    inputId: import("vue").ComputedRef<string>;
    popper: import("vue").ComputedRef<PopperOptions>;
    trigger: import("vue").Ref<import("@vueuse/core").MaybeElement, import("@vueuse/core").MaybeElement>;
    container: import("vue").Ref<import("@vueuse/core").MaybeElement, import("@vueuse/core").MaybeElement>;
    selected: import("vue").ComputedRef<any>;
    label: import("vue").ComputedRef<any>;
    accessor: <T extends Record<string, any>>(obj: T, key: string) => any;
    isLeading: import("vue").ComputedRef<string | true>;
    isTrailing: import("vue").ComputedRef<string | true>;
    selectClass: import("vue").ComputedRef<string>;
    leadingIconName: import("vue").ComputedRef<string>;
    leadingIconClass: import("vue").ComputedRef<string>;
    leadingWrapperIconClass: import("vue").ComputedRef<string>;
    trailingIconName: import("vue").ComputedRef<string>;
    trailingIconClass: import("vue").ComputedRef<string>;
    trailingWrapperIconClass: import("vue").ComputedRef<string>;
    filteredOptions: import("vue").ComputedRef<any[]>;
    createOption: import("vue").ComputedRef<string | {
        [x: string]: string;
    }>;
    query: import("vue").WritableComputedRef<string, string>;
    onUpdate: (value: any) => void;
    onQueryChange: (event: any) => void;
    by: import("vue").ComputedRef<string | ((a: any, z: any) => boolean)>;
}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, ("change" | "close" | "update:modelValue" | "open" | "update:query")[], "change" | "close" | "update:modelValue" | "open" | "update:query", import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    modelValue: {
        type: PropType<string | number | object | Array<any> | boolean | null>;
        default: string;
    };
    query: {
        type: StringConstructor;
        default: any;
    };
    by: {
        type: StringConstructor;
        default: any;
    };
    options: {
        type: PropType<{
            [key: string]: any;
            disabled?: boolean;
        }[] | string[]>;
        default: () => any[];
    };
    id: {
        type: StringConstructor;
        default: any;
    };
    name: {
        type: StringConstructor;
        default: any;
    };
    required: {
        type: BooleanConstructor;
        default: boolean;
    };
    icon: {
        type: StringConstructor;
        default: any;
    };
    loadingIcon: {
        type: StringConstructor;
        default: () => string;
    };
    leadingIcon: {
        type: StringConstructor;
        default: any;
    };
    trailingIcon: {
        type: StringConstructor;
        default: () => string;
    };
    trailing: {
        type: BooleanConstructor;
        default: boolean;
    };
    leading: {
        type: BooleanConstructor;
        default: boolean;
    };
    loading: {
        type: BooleanConstructor;
        default: boolean;
    };
    selectedIcon: {
        type: StringConstructor;
        default: () => string;
    };
    disabled: {
        type: BooleanConstructor;
        default: boolean;
    };
    multiple: {
        type: BooleanConstructor;
        default: boolean;
    };
    searchable: {
        type: PropType<boolean | ((query: string) => Promise<any[]> | any[])>;
        default: boolean;
    };
    searchablePlaceholder: {
        type: StringConstructor;
        default: () => string;
    };
    searchableLazy: {
        type: BooleanConstructor;
        default: boolean;
    };
    clearSearchOnClose: {
        type: BooleanConstructor;
        default: () => boolean;
    };
    debounce: {
        type: NumberConstructor;
        default: number;
    };
    creatable: {
        type: BooleanConstructor;
        default: boolean;
    };
    showCreateOptionWhen: {
        type: PropType<"always" | "empty" | ((query: string, results: any[]) => boolean)>;
        default: () => string;
    };
    placeholder: {
        type: StringConstructor;
        default: any;
    };
    padded: {
        type: BooleanConstructor;
        default: boolean;
    };
    size: {
        type: PropType<SelectSize>;
        default: any;
        validator(value: string): boolean;
    };
    color: {
        type: PropType<SelectColor>;
        default: () => string;
        validator(value: string): boolean;
    };
    variant: {
        type: PropType<SelectVariant>;
        default: () => string;
        validator(value: string): boolean;
    };
    optionAttribute: {
        type: StringConstructor;
        default: string;
    };
    valueAttribute: {
        type: StringConstructor;
        default: any;
    };
    searchAttributes: {
        type: ArrayConstructor;
        default: any;
    };
    inputTargetForm: {
        type: StringConstructor;
        default: any;
    };
    popper: {
        type: PropType<PopperOptions>;
        default: () => {};
    };
    selectClass: {
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
    uiMenu: {
        type: PropType<DeepPartial<typeof configMenu> & {
            strategy?: Strategy;
        }>;
        default: () => {};
    };
}>> & Readonly<{
    onChange?: (...args: any[]) => any;
    onClose?: (...args: any[]) => any;
    "onUpdate:modelValue"?: (...args: any[]) => any;
    onOpen?: (...args: any[]) => any;
    "onUpdate:query"?: (...args: any[]) => any;
}>, {
    name: string;
    required: boolean;
    size: SelectSize;
    class: any;
    placeholder: string;
    ui: {
        form?: string;
        placeholder?: string;
        default?: DeepPartial<{
            size: string;
            color: string;
            variant: string;
            loadingIcon: string;
            trailingIcon: string;
        }, any>;
        wrapper?: string;
        base?: string;
        rounded?: string;
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
    color: SelectColor;
    variant: SelectVariant;
    leading: boolean;
    popper: PopperOptions;
    disabled: boolean;
    query: string;
    id: string;
    icon: string;
    options: string[] | {
        [key: string]: any;
        disabled?: boolean;
    }[];
    modelValue: string | number | boolean | object | any[];
    loading: boolean;
    padded: boolean;
    loadingIcon: string;
    leadingIcon: string;
    trailingIcon: string;
    trailing: boolean;
    by: string;
    multiple: boolean;
    selectedIcon: string;
    optionAttribute: string;
    valueAttribute: string;
    searchAttributes: unknown[];
    debounce: number;
    uiMenu: {
        select?: string;
        input?: string;
        required?: string;
        label?: string;
        option?: DeepPartial<{
            create: string;
            base: string;
            rounded: string;
            padding: string;
            size: string;
            color: string;
            container: string;
            active: string;
            inactive: string;
            selected: string;
            disabled: string;
            empty: string;
            icon: {
                base: string;
                active: string;
                inactive: string;
            };
            selectedIcon: {
                wrapper: string;
                padding: string;
                base: string;
            };
            avatar: {
                base: string;
                size: import("../../types/avatar.js").AvatarSize;
            };
            chip: {
                base: string;
            };
        }, any>;
        transition?: DeepPartial<{
            leaveActiveClass: string;
            leaveFromClass: string;
            leaveToClass: string;
        }, any>;
        popper?: DeepPartial<{
            placement: string;
        }, any>;
        default?: DeepPartial<{
            selectedIcon: string;
            clearSearchOnClose: boolean;
            showCreateOptionWhen: string;
            searchablePlaceholder: {
                label: string;
            };
            empty: {
                label: string;
            };
            optionEmpty: {
                label: string;
            };
        }, any>;
        arrow?: DeepPartial<{
            ring: string;
            background: string;
            base: string;
            rounded: string;
            shadow: string;
            placement: string;
        }, any>;
        container?: string;
        trigger?: string;
        width?: string;
        height?: string;
        base?: string;
        background?: string;
        shadow?: string;
        rounded?: string;
        padding?: string;
        ring?: string;
        empty?: string;
    } & {
        [key: string]: any;
    } & {
        strategy?: Strategy;
    };
    selectClass: string;
    searchable: boolean | ((query: string) => Promise<any[]> | any[]);
    searchablePlaceholder: string;
    searchableLazy: boolean;
    clearSearchOnClose: boolean;
    creatable: boolean;
    showCreateOptionWhen: "empty" | "always" | ((query: string, results: any[]) => boolean);
    inputTargetForm: string;
}, {}, {
    HCombobox: import("vue").DefineComponent<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
        disabled: {
            type: BooleanConstructor[];
            default: boolean;
        };
        by: {
            type: (StringConstructor | FunctionConstructor)[];
            nullable: boolean;
            default: null;
        };
        modelValue: {
            type: PropType<string | number | boolean | object | null>;
            default: undefined;
        };
        defaultValue: {
            type: PropType<string | number | boolean | object | null>;
            default: undefined;
        };
        form: {
            type: StringConstructor;
            optional: boolean;
        };
        name: {
            type: StringConstructor;
            optional: boolean;
        };
        nullable: {
            type: BooleanConstructor;
            default: boolean;
        };
        multiple: {
            type: BooleanConstructor[];
            default: boolean;
        };
        immediate: {
            type: BooleanConstructor[];
            default: boolean;
        };
        virtual: {
            type: PropType<{
                options: unknown[];
                disabled?: ((value: unknown) => boolean) | undefined;
            } | null>;
            default: null;
        };
    }, () => import("vue").VNode<import("vue").RendererNode, import("vue").RendererElement, {
        [key: string]: any;
    }>, unknown, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {
        'update:modelValue': (_value: any) => true;
    }, string, import("vue").VNodeProps & import("vue").AllowedComponentProps & import("vue").ComponentCustomProps, Readonly<import("vue").ExtractPropTypes<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
        disabled: {
            type: BooleanConstructor[];
            default: boolean;
        };
        by: {
            type: (StringConstructor | FunctionConstructor)[];
            nullable: boolean;
            default: null;
        };
        modelValue: {
            type: PropType<string | number | boolean | object | null>;
            default: undefined;
        };
        defaultValue: {
            type: PropType<string | number | boolean | object | null>;
            default: undefined;
        };
        form: {
            type: StringConstructor;
            optional: boolean;
        };
        name: {
            type: StringConstructor;
            optional: boolean;
        };
        nullable: {
            type: BooleanConstructor;
            default: boolean;
        };
        multiple: {
            type: BooleanConstructor[];
            default: boolean;
        };
        immediate: {
            type: BooleanConstructor[];
            default: boolean;
        };
        virtual: {
            type: PropType<{
                options: unknown[];
                disabled?: ((value: unknown) => boolean) | undefined;
            } | null>;
            default: null;
        };
    }>> & {
        "onUpdate:modelValue"?: ((_value: any) => any) | undefined;
    }, {
        as: string | Record<string, any>;
        disabled: boolean;
        by: string | Function;
        modelValue: string | number | boolean | object | null;
        defaultValue: string | number | boolean | object | null;
        nullable: boolean;
        multiple: boolean;
        immediate: boolean;
        virtual: {
            options: unknown[];
            disabled?: ((value: unknown) => boolean) | undefined;
        } | null;
    }, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
    HComboboxButton: import("vue").DefineComponent<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
        id: {
            type: StringConstructor;
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
        id: {
            type: StringConstructor;
            default: null;
        };
    }>>, {
        id: string;
        as: string | Record<string, any>;
    }, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
    HComboboxOptions: import("vue").DefineComponent<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
        static: {
            type: BooleanConstructor;
            default: boolean;
        };
        unmount: {
            type: BooleanConstructor;
            default: boolean;
        };
        hold: {
            type: BooleanConstructor[];
            default: boolean;
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
        static: {
            type: BooleanConstructor;
            default: boolean;
        };
        unmount: {
            type: BooleanConstructor;
            default: boolean;
        };
        hold: {
            type: BooleanConstructor[];
            default: boolean;
        };
    }>>, {
        as: string | Record<string, any>;
        unmount: boolean;
        static: boolean;
        hold: boolean;
    }, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
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
    HComboboxInput: import("vue").DefineComponent<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
        static: {
            type: BooleanConstructor;
            default: boolean;
        };
        unmount: {
            type: BooleanConstructor;
            default: boolean;
        };
        displayValue: {
            type: PropType<(item: unknown) => string>;
        };
        defaultValue: {
            type: StringConstructor;
            default: undefined;
        };
        id: {
            type: StringConstructor;
            default: null;
        };
    }, () => import("vue").VNode<import("vue").RendererNode, import("vue").RendererElement, {
        [key: string]: any;
    }> | import("vue").VNode<import("vue").RendererNode, import("vue").RendererElement, {
        [key: string]: any;
    }>[] | null, unknown, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {
        change: (_value: Event & {
            target: HTMLInputElement;
        }) => true;
    }, string, import("vue").VNodeProps & import("vue").AllowedComponentProps & import("vue").ComponentCustomProps, Readonly<import("vue").ExtractPropTypes<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
        static: {
            type: BooleanConstructor;
            default: boolean;
        };
        unmount: {
            type: BooleanConstructor;
            default: boolean;
        };
        displayValue: {
            type: PropType<(item: unknown) => string>;
        };
        defaultValue: {
            type: StringConstructor;
            default: undefined;
        };
        id: {
            type: StringConstructor;
            default: null;
        };
    }>> & {
        onChange?: ((_value: Event & {
            target: HTMLInputElement;
        }) => any) | undefined;
    }, {
        id: string;
        as: string | Record<string, any>;
        unmount: boolean;
        static: boolean;
        defaultValue: string;
    }, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
    HListbox: import("vue").DefineComponent<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
        disabled: {
            type: BooleanConstructor[];
            default: boolean;
        };
        by: {
            type: (StringConstructor | FunctionConstructor)[];
            default: () => <T>(a: T, z: T) => boolean;
        };
        horizontal: {
            type: BooleanConstructor[];
            default: boolean;
        };
        modelValue: {
            type: PropType<string | number | boolean | object | null>;
            default: undefined;
        };
        defaultValue: {
            type: PropType<string | number | boolean | object | null>;
            default: undefined;
        };
        form: {
            type: StringConstructor;
            optional: boolean;
        };
        name: {
            type: StringConstructor;
            optional: boolean;
        };
        multiple: {
            type: BooleanConstructor[];
            default: boolean;
        };
    }, () => import("vue").VNode<import("vue").RendererNode, import("vue").RendererElement, {
        [key: string]: any;
    }>, unknown, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {
        'update:modelValue': (_value: any) => true;
    }, string, import("vue").VNodeProps & import("vue").AllowedComponentProps & import("vue").ComponentCustomProps, Readonly<import("vue").ExtractPropTypes<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
        disabled: {
            type: BooleanConstructor[];
            default: boolean;
        };
        by: {
            type: (StringConstructor | FunctionConstructor)[];
            default: () => <T>(a: T, z: T) => boolean;
        };
        horizontal: {
            type: BooleanConstructor[];
            default: boolean;
        };
        modelValue: {
            type: PropType<string | number | boolean | object | null>;
            default: undefined;
        };
        defaultValue: {
            type: PropType<string | number | boolean | object | null>;
            default: undefined;
        };
        form: {
            type: StringConstructor;
            optional: boolean;
        };
        name: {
            type: StringConstructor;
            optional: boolean;
        };
        multiple: {
            type: BooleanConstructor[];
            default: boolean;
        };
    }>> & {
        "onUpdate:modelValue"?: ((_value: any) => any) | undefined;
    }, {
        as: string | Record<string, any>;
        disabled: boolean;
        horizontal: boolean;
        by: string | Function;
        modelValue: string | number | boolean | object | null;
        defaultValue: string | number | boolean | object | null;
        multiple: boolean;
    }, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
    HListboxButton: import("vue").DefineComponent<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
        id: {
            type: StringConstructor;
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
        id: {
            type: StringConstructor;
            default: null;
        };
    }>>, {
        id: string;
        as: string | Record<string, any>;
    }, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
    HListboxOptions: import("vue").DefineComponent<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
        static: {
            type: BooleanConstructor;
            default: boolean;
        };
        unmount: {
            type: BooleanConstructor;
            default: boolean;
        };
        id: {
            type: StringConstructor;
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
        static: {
            type: BooleanConstructor;
            default: boolean;
        };
        unmount: {
            type: BooleanConstructor;
            default: boolean;
        };
        id: {
            type: StringConstructor;
            default: null;
        };
    }>>, {
        id: string;
        as: string | Record<string, any>;
        unmount: boolean;
        static: boolean;
    }, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
    HListboxOption: import("vue").DefineComponent<{
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
        id: {
            type: StringConstructor;
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
        id: {
            type: StringConstructor;
            default: null;
        };
    }>>, {
        id: string;
        as: string | Record<string, any>;
        disabled: boolean;
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
            type: PropType<DeepPartial<{
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
                strategy?: Strategy;
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
            type: PropType<DeepPartial<{
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
                strategy?: Strategy;
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
            size?: DeepPartial<{
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
            chip?: DeepPartial<{
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
            icon?: DeepPartial<{
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
            default?: DeepPartial<{
                size: string;
                icon: any;
                chipColor: any;
                chipPosition: string;
            }, any>;
        } & {
            [key: string]: any;
        } & {
            strategy?: Strategy;
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
}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
export default _default;
