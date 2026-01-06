import type { ComputedRef, PropType, ComponentPublicInstance } from 'vue';
import type { UseFuseOptions } from '@vueuse/integrations/useFuse';
import type { Group, Command, Button, Strategy, DeepPartial } from '../../types/index.js';
import { commandPalette } from '#ui/ui.config';
declare const config: {
    wrapper: string;
    container: string;
    input: {
        wrapper: string;
        base: string;
        padding: string;
        height: string;
        size: string;
        icon: {
            base: string;
            loading: string;
            size: string;
            padding: string;
        };
        closeButton: {
            base: string;
            padding: string;
        };
    };
    emptyState: {
        wrapper: string;
        label: string;
        queryLabel: string;
        icon: string;
    };
    group: {
        wrapper: string;
        label: string;
        container: string;
        command: {
            base: string;
            active: string;
            inactive: string;
            label: string;
            prefix: string;
            suffix: string;
            container: string;
            icon: {
                base: string;
                active: string;
                inactive: string;
            };
            selectedIcon: {
                base: string;
            };
            avatar: {
                base: string;
                size: import("../../types/avatar.js").AvatarSize;
            };
            chip: {
                base: string;
            };
            disabled: string;
            shortcuts: string;
        };
        active: string;
        inactive: string;
    };
    default: {
        icon: string;
        loadingIcon: string;
        emptyState: {
            icon: string;
            label: string;
            queryLabel: string;
        };
        closeButton: any;
        selectedIcon: string;
    };
};
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
    modelValue: {
        type: (StringConstructor | NumberConstructor | ObjectConstructor | ArrayConstructor)[];
        default: any;
    };
    by: {
        type: StringConstructor;
        default: string;
    };
    multiple: {
        type: BooleanConstructor;
        default: boolean;
    };
    nullable: {
        type: BooleanConstructor;
        default: boolean;
    };
    searchable: {
        type: BooleanConstructor;
        default: boolean;
    };
    loading: {
        type: BooleanConstructor;
        default: boolean;
    };
    groups: {
        type: PropType<Group[]>;
        default: () => any[];
    };
    icon: {
        type: StringConstructor;
        default: () => string;
    };
    loadingIcon: {
        type: StringConstructor;
        default: () => string;
    };
    selectedIcon: {
        type: StringConstructor;
        default: () => string;
    };
    closeButton: {
        type: PropType<Button>;
        default: () => Button;
    };
    emptyState: {
        type: PropType<{
            icon: string;
            label: string;
            queryLabel: string;
        }>;
        default: () => {
            icon: string;
            label: string;
            queryLabel: string;
        };
    };
    placeholder: {
        type: StringConstructor;
        default: string;
    };
    groupAttribute: {
        type: StringConstructor;
        default: string;
    };
    commandAttribute: {
        type: StringConstructor;
        default: string;
    };
    autoselect: {
        type: BooleanConstructor;
        default: boolean;
    };
    autoclear: {
        type: BooleanConstructor;
        default: boolean;
    };
    debounce: {
        type: NumberConstructor;
        default: number;
    };
    fuse: {
        type: PropType<UseFuseOptions<Command>>;
        default: () => {};
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
    ui: ComputedRef<{
        wrapper: string;
        container: string;
        input: {
            wrapper: string;
            base: string;
            padding: string;
            height: string;
            size: string;
            icon: {
                base: string;
                loading: string;
                size: string;
                padding: string;
            };
            closeButton: {
                base: string;
                padding: string;
            };
        };
        emptyState: {
            wrapper: string;
            label: string;
            queryLabel: string;
            icon: string;
        };
        group: {
            wrapper: string;
            label: string;
            container: string;
            command: {
                base: string;
                active: string;
                inactive: string;
                label: string;
                prefix: string;
                suffix: string;
                container: string;
                icon: {
                    base: string;
                    active: string;
                    inactive: string;
                };
                selectedIcon: {
                    base: string;
                };
                avatar: {
                    base: string;
                    size: import("../../types/avatar.js").AvatarSize;
                };
                chip: {
                    base: string;
                };
                disabled: string;
                shortcuts: string;
            };
            active: string;
            inactive: string;
        };
        default: {
            icon: string;
            loadingIcon: string;
            emptyState: {
                icon: string;
                label: string;
                queryLabel: string;
            };
            closeButton: any;
            selectedIcon: string;
        };
    }>;
    attrs: ComputedRef<Pick<{
        [x: string]: unknown;
    }, string>>;
    groups: ComputedRef<(Group | {
        commands: Command[];
        key: string;
        active?: string;
        inactive?: string;
        search?: (...args: any[]) => any[] | Promise<any[]>;
        filter?: (...args: any[]) => Command[];
        static?: boolean;
    })[]>;
    comboboxInput: import("vue").Ref<ComponentPublicInstance<HTMLInputElement>, ComponentPublicInstance<HTMLInputElement>>;
    query: import("vue").Ref<string, string>;
    iconName: ComputedRef<string>;
    iconClass: ComputedRef<string>;
    emptyState: ComputedRef<{
        icon: string;
        label: string;
        queryLabel: string;
    }>;
    onSelect: (option: Command | Command[]) => void;
    onClear: () => void;
}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, ("close" | "update:modelValue")[], "close" | "update:modelValue", import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    modelValue: {
        type: (StringConstructor | NumberConstructor | ObjectConstructor | ArrayConstructor)[];
        default: any;
    };
    by: {
        type: StringConstructor;
        default: string;
    };
    multiple: {
        type: BooleanConstructor;
        default: boolean;
    };
    nullable: {
        type: BooleanConstructor;
        default: boolean;
    };
    searchable: {
        type: BooleanConstructor;
        default: boolean;
    };
    loading: {
        type: BooleanConstructor;
        default: boolean;
    };
    groups: {
        type: PropType<Group[]>;
        default: () => any[];
    };
    icon: {
        type: StringConstructor;
        default: () => string;
    };
    loadingIcon: {
        type: StringConstructor;
        default: () => string;
    };
    selectedIcon: {
        type: StringConstructor;
        default: () => string;
    };
    closeButton: {
        type: PropType<Button>;
        default: () => Button;
    };
    emptyState: {
        type: PropType<{
            icon: string;
            label: string;
            queryLabel: string;
        }>;
        default: () => {
            icon: string;
            label: string;
            queryLabel: string;
        };
    };
    placeholder: {
        type: StringConstructor;
        default: string;
    };
    groupAttribute: {
        type: StringConstructor;
        default: string;
    };
    commandAttribute: {
        type: StringConstructor;
        default: string;
    };
    autoselect: {
        type: BooleanConstructor;
        default: boolean;
    };
    autoclear: {
        type: BooleanConstructor;
        default: boolean;
    };
    debounce: {
        type: NumberConstructor;
        default: number;
    };
    fuse: {
        type: PropType<UseFuseOptions<Command>>;
        default: () => {};
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
}>> & Readonly<{
    onClose?: (...args: any[]) => any;
    "onUpdate:modelValue"?: (...args: any[]) => any;
}>, {
    class: any;
    placeholder: string;
    ui: {
        wrapper?: string;
        container?: string;
        input?: DeepPartial<{
            wrapper: string;
            base: string;
            padding: string;
            height: string;
            size: string;
            icon: {
                base: string;
                loading: string;
                size: string;
                padding: string;
            };
            closeButton: {
                base: string;
                padding: string;
            };
        }, any>;
        emptyState?: DeepPartial<{
            wrapper: string;
            label: string;
            queryLabel: string;
            icon: string;
        }, any>;
        group?: DeepPartial<{
            wrapper: string;
            label: string;
            container: string;
            command: {
                base: string;
                active: string;
                inactive: string;
                label: string;
                prefix: string;
                suffix: string;
                container: string;
                icon: {
                    base: string;
                    active: string;
                    inactive: string;
                };
                selectedIcon: {
                    base: string;
                };
                avatar: {
                    base: string;
                    size: import("../../types/avatar.js").AvatarSize;
                };
                chip: {
                    base: string;
                };
                disabled: string;
                shortcuts: string;
            };
            active: string;
            inactive: string;
        }, any>;
        default?: DeepPartial<{
            icon: string;
            loadingIcon: string;
            emptyState: {
                icon: string;
                label: string;
                queryLabel: string;
            };
            closeButton: any;
            selectedIcon: string;
        }, any>;
    } & {
        [key: string]: any;
    } & {
        strategy?: Strategy;
    };
    icon: string;
    closeButton: Button;
    modelValue: string | number | Record<string, any> | unknown[];
    loading: boolean;
    loadingIcon: string;
    emptyState: {
        icon: string;
        label: string;
        queryLabel: string;
    };
    by: string;
    multiple: boolean;
    nullable: boolean;
    selectedIcon: string;
    debounce: number;
    searchable: boolean;
    groupAttribute: string;
    commandAttribute: string;
    groups: Group[];
    autoselect: boolean;
    autoclear: boolean;
    fuse: UseFuseOptions<Command>;
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
    UButton: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
        type: {
            type: StringConstructor;
            default: string;
        };
        block: {
            type: BooleanConstructor;
            default: boolean;
        };
        label: {
            type: StringConstructor;
            default: any;
        };
        loading: {
            type: BooleanConstructor;
            default: boolean;
        };
        disabled: {
            type: BooleanConstructor;
            default: boolean;
        };
        padded: {
            type: BooleanConstructor;
            default: boolean;
        };
        size: {
            type: PropType<import("../../types/button.js").ButtonSize>;
            default: () => string;
            validator(value: string): boolean;
        };
        color: {
            type: PropType<import("../../types/button.js").ButtonColor>;
            default: () => string;
            validator(value: string): boolean;
        };
        variant: {
            type: PropType<import("../../types/button.js").ButtonVariant>;
            default: () => string;
            validator(value: string): boolean;
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
        square: {
            type: BooleanConstructor;
            default: boolean;
        };
        truncate: {
            type: BooleanConstructor;
            default: boolean;
        };
        class: {
            type: PropType<any>;
            default: () => "";
        };
        ui: {
            type: PropType<DeepPartial<{
                base: string;
                font: string;
                rounded: string;
                truncate: string;
                block: string;
                inline: string;
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
                square: {
                    '2xs': string;
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                    xl: string;
                };
                color: {
                    white: {
                        solid: string;
                        ghost: string;
                    };
                    gray: {
                        solid: string;
                        ghost: string;
                        link: string;
                    };
                    black: {
                        solid: string;
                        link: string;
                    };
                };
                variant: {
                    solid: string;
                    outline: string;
                    soft: string;
                    ghost: string;
                    link: string;
                };
                icon: {
                    base: string;
                    loading: string;
                    size: {
                        '2xs': string;
                        xs: string;
                        sm: string;
                        md: string;
                        lg: string;
                        xl: string;
                    };
                };
                default: {
                    size: string;
                    variant: string;
                    color: string;
                    loadingIcon: string;
                };
            }> & {
                strategy?: Strategy;
            }>;
            default: () => {};
        };
        to: {
            readonly type: PropType<import("vue-router").RouteLocationRaw>;
            readonly default: any;
            readonly required: false;
        };
        href: {
            readonly type: PropType<import("vue-router").RouteLocationRaw>;
            readonly default: any;
            readonly required: false;
        };
        target: {
            readonly type: PropType<import("#app").NuxtLinkProps["target"]>;
            readonly default: any;
            readonly required: false;
        };
        rel: {
            readonly type: PropType<any>;
            readonly default: any;
            readonly required: false;
        };
        noRel: {
            readonly type: PropType<import("#app").NuxtLinkProps["noRel"]>;
            readonly default: any;
            readonly required: false;
        };
        prefetch: {
            readonly type: PropType<import("#app").NuxtLinkProps["prefetch"]>;
            readonly default: any;
            readonly required: false;
        };
        noPrefetch: {
            readonly type: PropType<import("#app").NuxtLinkProps["noPrefetch"]>;
            readonly default: any;
            readonly required: false;
        };
        activeClass: {
            readonly type: PropType<import("#app").NuxtLinkProps["activeClass"]>;
            readonly default: any;
            readonly required: false;
        };
        exactActiveClass: {
            readonly type: PropType<import("#app").NuxtLinkProps["exactActiveClass"]>;
            readonly default: any;
            readonly required: false;
        };
        prefetchedClass: {
            readonly type: PropType<import("#app").NuxtLinkProps["prefetchedClass"]>;
            readonly default: any;
            readonly required: false;
        };
        replace: {
            readonly type: PropType<import("#app").NuxtLinkProps["replace"]>;
            readonly default: any;
            readonly required: false;
        };
        ariaCurrentValue: {
            readonly type: PropType<import("#app").NuxtLinkProps["ariaCurrentValue"]>;
            readonly default: any;
            readonly required: false;
        };
        external: {
            readonly type: PropType<import("#app").NuxtLinkProps["external"]>;
            readonly default: any;
            readonly required: false;
        };
    }>, {
        ui: ComputedRef<{
            base: string;
            font: string;
            rounded: string;
            truncate: string;
            block: string;
            inline: string;
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
            square: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
            };
            color: {
                white: {
                    solid: string;
                    ghost: string;
                };
                gray: {
                    solid: string;
                    ghost: string;
                    link: string;
                };
                black: {
                    solid: string;
                    link: string;
                };
            };
            variant: {
                solid: string;
                outline: string;
                soft: string;
                ghost: string;
                link: string;
            };
            icon: {
                base: string;
                loading: string;
                size: {
                    '2xs': string;
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                    xl: string;
                };
            };
            default: {
                size: string;
                variant: string;
                color: string;
                loadingIcon: string;
            };
        }>;
        attrs: ComputedRef<Pick<{
            [x: string]: unknown;
        }, string>>;
        isLeading: ComputedRef<string | true>;
        isTrailing: ComputedRef<string | true>;
        isSquare: ComputedRef<boolean>;
        buttonClass: ComputedRef<string>;
        leadingIconName: ComputedRef<string>;
        trailingIconName: ComputedRef<string>;
        leadingIconClass: ComputedRef<string>;
        trailingIconClass: ComputedRef<string>;
        linkProps: ComputedRef<{}>;
    }, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
        type: {
            type: StringConstructor;
            default: string;
        };
        block: {
            type: BooleanConstructor;
            default: boolean;
        };
        label: {
            type: StringConstructor;
            default: any;
        };
        loading: {
            type: BooleanConstructor;
            default: boolean;
        };
        disabled: {
            type: BooleanConstructor;
            default: boolean;
        };
        padded: {
            type: BooleanConstructor;
            default: boolean;
        };
        size: {
            type: PropType<import("../../types/button.js").ButtonSize>;
            default: () => string;
            validator(value: string): boolean;
        };
        color: {
            type: PropType<import("../../types/button.js").ButtonColor>;
            default: () => string;
            validator(value: string): boolean;
        };
        variant: {
            type: PropType<import("../../types/button.js").ButtonVariant>;
            default: () => string;
            validator(value: string): boolean;
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
        square: {
            type: BooleanConstructor;
            default: boolean;
        };
        truncate: {
            type: BooleanConstructor;
            default: boolean;
        };
        class: {
            type: PropType<any>;
            default: () => "";
        };
        ui: {
            type: PropType<DeepPartial<{
                base: string;
                font: string;
                rounded: string;
                truncate: string;
                block: string;
                inline: string;
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
                square: {
                    '2xs': string;
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                    xl: string;
                };
                color: {
                    white: {
                        solid: string;
                        ghost: string;
                    };
                    gray: {
                        solid: string;
                        ghost: string;
                        link: string;
                    };
                    black: {
                        solid: string;
                        link: string;
                    };
                };
                variant: {
                    solid: string;
                    outline: string;
                    soft: string;
                    ghost: string;
                    link: string;
                };
                icon: {
                    base: string;
                    loading: string;
                    size: {
                        '2xs': string;
                        xs: string;
                        sm: string;
                        md: string;
                        lg: string;
                        xl: string;
                    };
                };
                default: {
                    size: string;
                    variant: string;
                    color: string;
                    loadingIcon: string;
                };
            }> & {
                strategy?: Strategy;
            }>;
            default: () => {};
        };
        to: {
            readonly type: PropType<import("vue-router").RouteLocationRaw>;
            readonly default: any;
            readonly required: false;
        };
        href: {
            readonly type: PropType<import("vue-router").RouteLocationRaw>;
            readonly default: any;
            readonly required: false;
        };
        target: {
            readonly type: PropType<import("#app").NuxtLinkProps["target"]>;
            readonly default: any;
            readonly required: false;
        };
        rel: {
            readonly type: PropType<any>;
            readonly default: any;
            readonly required: false;
        };
        noRel: {
            readonly type: PropType<import("#app").NuxtLinkProps["noRel"]>;
            readonly default: any;
            readonly required: false;
        };
        prefetch: {
            readonly type: PropType<import("#app").NuxtLinkProps["prefetch"]>;
            readonly default: any;
            readonly required: false;
        };
        noPrefetch: {
            readonly type: PropType<import("#app").NuxtLinkProps["noPrefetch"]>;
            readonly default: any;
            readonly required: false;
        };
        activeClass: {
            readonly type: PropType<import("#app").NuxtLinkProps["activeClass"]>;
            readonly default: any;
            readonly required: false;
        };
        exactActiveClass: {
            readonly type: PropType<import("#app").NuxtLinkProps["exactActiveClass"]>;
            readonly default: any;
            readonly required: false;
        };
        prefetchedClass: {
            readonly type: PropType<import("#app").NuxtLinkProps["prefetchedClass"]>;
            readonly default: any;
            readonly required: false;
        };
        replace: {
            readonly type: PropType<import("#app").NuxtLinkProps["replace"]>;
            readonly default: any;
            readonly required: false;
        };
        ariaCurrentValue: {
            readonly type: PropType<import("#app").NuxtLinkProps["ariaCurrentValue"]>;
            readonly default: any;
            readonly required: false;
        };
        external: {
            readonly type: PropType<import("#app").NuxtLinkProps["external"]>;
            readonly default: any;
            readonly required: false;
        };
    }>> & Readonly<{}>, {
        type: string;
        size: import("../../types/button.js").ButtonSize;
        class: any;
        target: "_blank" | "_parent" | "_self" | "_top" | (string & {});
        ui: {
            base?: string;
            font?: string;
            rounded?: string;
            truncate?: string;
            block?: string;
            inline?: string;
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
            square?: DeepPartial<{
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
            }, any>;
            color?: DeepPartial<{
                white: {
                    solid: string;
                    ghost: string;
                };
                gray: {
                    solid: string;
                    ghost: string;
                    link: string;
                };
                black: {
                    solid: string;
                    link: string;
                };
            }, any>;
            variant?: DeepPartial<{
                solid: string;
                outline: string;
                soft: string;
                ghost: string;
                link: string;
            }, any>;
            icon?: DeepPartial<{
                base: string;
                loading: string;
                size: {
                    '2xs': string;
                    xs: string;
                    sm: string;
                    md: string;
                    lg: string;
                    xl: string;
                };
            }, any>;
            default?: DeepPartial<{
                size: string;
                variant: string;
                color: string;
                loadingIcon: string;
            }, any>;
        } & {
            [key: string]: any;
        } & {
            strategy?: Strategy;
        };
        color: import("../../types/button.js").ButtonColor;
        variant: import("../../types/button.js").ButtonVariant;
        truncate: boolean;
        to: string | import("vue-router").RouteLocationAsRelativeGeneric | import("vue-router").RouteLocationAsPathGeneric;
        activeClass: string;
        exactActiveClass: string;
        ariaCurrentValue: "page" | "step" | "location" | "date" | "time" | "true" | "false";
        replace: boolean;
        noRel: boolean;
        prefetch: boolean;
        noPrefetch: boolean;
        prefetchedClass: string;
        external: boolean;
        leading: boolean;
        disabled: boolean;
        href: string | import("vue-router").RouteLocationAsRelativeGeneric | import("vue-router").RouteLocationAsPathGeneric;
        rel: any;
        label: string;
        icon: string;
        block: boolean;
        square: boolean;
        loading: boolean;
        padded: boolean;
        loadingIcon: string;
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
        ULink: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
            as: {
                type: StringConstructor;
                default: string;
            };
            type: {
                type: StringConstructor;
                default: string;
            };
            disabled: {
                type: BooleanConstructor;
                default: any;
            };
            active: {
                type: BooleanConstructor;
                default: any;
            };
            exact: {
                type: BooleanConstructor;
                default: boolean;
            };
            exactQuery: {
                type: PropType<boolean | "partial">;
                default: boolean;
            };
            exactHash: {
                type: BooleanConstructor;
                default: boolean;
            };
            inactiveClass: {
                type: StringConstructor;
                default: any;
            };
            to: {
                readonly type: PropType<import("vue-router").RouteLocationRaw>;
                readonly default: any;
                readonly required: false;
            };
            href: {
                readonly type: PropType<import("vue-router").RouteLocationRaw>;
                readonly default: any;
                readonly required: false;
            };
            target: {
                readonly type: PropType<import("#app").NuxtLinkProps["target"]>;
                readonly default: any;
                readonly required: false;
            };
            rel: {
                readonly type: PropType<any>;
                readonly default: any;
                readonly required: false;
            };
            noRel: {
                readonly type: PropType<import("#app").NuxtLinkProps["noRel"]>;
                readonly default: any;
                readonly required: false;
            };
            prefetch: {
                readonly type: PropType<import("#app").NuxtLinkProps["prefetch"]>;
                readonly default: any;
                readonly required: false;
            };
            noPrefetch: {
                readonly type: PropType<import("#app").NuxtLinkProps["noPrefetch"]>;
                readonly default: any;
                readonly required: false;
            };
            activeClass: {
                readonly type: PropType<import("#app").NuxtLinkProps["activeClass"]>;
                readonly default: any;
                readonly required: false;
            };
            exactActiveClass: {
                readonly type: PropType<import("#app").NuxtLinkProps["exactActiveClass"]>;
                readonly default: any;
                readonly required: false;
            };
            prefetchedClass: {
                readonly type: PropType<import("#app").NuxtLinkProps["prefetchedClass"]>;
                readonly default: any;
                readonly required: false;
            };
            replace: {
                readonly type: PropType<import("#app").NuxtLinkProps["replace"]>;
                readonly default: any;
                readonly required: false;
            };
            ariaCurrentValue: {
                readonly type: PropType<import("#app").NuxtLinkProps["ariaCurrentValue"]>;
                readonly default: any;
                readonly required: false;
            };
            external: {
                readonly type: PropType<import("#app").NuxtLinkProps["external"]>;
                readonly default: any;
                readonly required: false;
            };
        }>, {
            resolveLinkClass: (route: any, $route: any, { isActive, isExactActive }: {
                isActive: boolean;
                isExactActive: boolean;
            }) => string;
        }, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
            as: {
                type: StringConstructor;
                default: string;
            };
            type: {
                type: StringConstructor;
                default: string;
            };
            disabled: {
                type: BooleanConstructor;
                default: any;
            };
            active: {
                type: BooleanConstructor;
                default: any;
            };
            exact: {
                type: BooleanConstructor;
                default: boolean;
            };
            exactQuery: {
                type: PropType<boolean | "partial">;
                default: boolean;
            };
            exactHash: {
                type: BooleanConstructor;
                default: boolean;
            };
            inactiveClass: {
                type: StringConstructor;
                default: any;
            };
            to: {
                readonly type: PropType<import("vue-router").RouteLocationRaw>;
                readonly default: any;
                readonly required: false;
            };
            href: {
                readonly type: PropType<import("vue-router").RouteLocationRaw>;
                readonly default: any;
                readonly required: false;
            };
            target: {
                readonly type: PropType<import("#app").NuxtLinkProps["target"]>;
                readonly default: any;
                readonly required: false;
            };
            rel: {
                readonly type: PropType<any>;
                readonly default: any;
                readonly required: false;
            };
            noRel: {
                readonly type: PropType<import("#app").NuxtLinkProps["noRel"]>;
                readonly default: any;
                readonly required: false;
            };
            prefetch: {
                readonly type: PropType<import("#app").NuxtLinkProps["prefetch"]>;
                readonly default: any;
                readonly required: false;
            };
            noPrefetch: {
                readonly type: PropType<import("#app").NuxtLinkProps["noPrefetch"]>;
                readonly default: any;
                readonly required: false;
            };
            activeClass: {
                readonly type: PropType<import("#app").NuxtLinkProps["activeClass"]>;
                readonly default: any;
                readonly required: false;
            };
            exactActiveClass: {
                readonly type: PropType<import("#app").NuxtLinkProps["exactActiveClass"]>;
                readonly default: any;
                readonly required: false;
            };
            prefetchedClass: {
                readonly type: PropType<import("#app").NuxtLinkProps["prefetchedClass"]>;
                readonly default: any;
                readonly required: false;
            };
            replace: {
                readonly type: PropType<import("#app").NuxtLinkProps["replace"]>;
                readonly default: any;
                readonly required: false;
            };
            ariaCurrentValue: {
                readonly type: PropType<import("#app").NuxtLinkProps["ariaCurrentValue"]>;
                readonly default: any;
                readonly required: false;
            };
            external: {
                readonly type: PropType<import("#app").NuxtLinkProps["external"]>;
                readonly default: any;
                readonly required: false;
            };
        }>> & Readonly<{}>, {
            type: string;
            target: "_blank" | "_parent" | "_self" | "_top" | (string & {});
            to: string | import("vue-router").RouteLocationAsRelativeGeneric | import("vue-router").RouteLocationAsPathGeneric;
            activeClass: string;
            exactActiveClass: string;
            ariaCurrentValue: "page" | "step" | "location" | "date" | "time" | "true" | "false";
            replace: boolean;
            noRel: boolean;
            prefetch: boolean;
            noPrefetch: boolean;
            prefetchedClass: string;
            external: boolean;
            as: string;
            disabled: boolean;
            active: boolean;
            exact: boolean;
            exactQuery: boolean | "partial";
            exactHash: boolean;
            inactiveClass: string;
            href: string | import("vue-router").RouteLocationAsRelativeGeneric | import("vue-router").RouteLocationAsPathGeneric;
            rel: any;
        }, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
    }, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
    CommandPaletteGroup: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
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
        label: ComputedRef<any>;
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
            ui: ComputedRef<{
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
            attrs: ComputedRef<Pick<{
                [x: string]: unknown;
            }, string>>;
            wrapperClass: ComputedRef<string>;
            imgClass: ComputedRef<string>;
            iconClass: ComputedRef<string>;
            chipClass: ComputedRef<string>;
            url: ComputedRef<string>;
            placeholder: ComputedRef<string>;
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
                type: PropType<DeepPartial<{
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
                    strategy?: Strategy;
                }>;
                default: () => {};
            };
        }>, {
            ui: ComputedRef<{
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
            attrs: ComputedRef<Pick<{
                [x: string]: unknown;
            }, string>>;
            kbdClass: ComputedRef<string>;
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
                type: PropType<DeepPartial<{
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
                    strategy?: Strategy;
                }>;
                default: () => {};
            };
        }>> & Readonly<{}>, {
            size: import("../../types/kbd.js").KbdSize;
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
    }, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
export default _default;
