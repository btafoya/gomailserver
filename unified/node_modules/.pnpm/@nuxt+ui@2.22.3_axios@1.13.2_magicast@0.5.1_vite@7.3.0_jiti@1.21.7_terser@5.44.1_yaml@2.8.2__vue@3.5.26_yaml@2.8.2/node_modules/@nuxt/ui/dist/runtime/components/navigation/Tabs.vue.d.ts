import type { PropType } from 'vue';
import type { TabItem, Strategy, DeepPartial } from '../../types/index.js';
declare const config: {
    wrapper: string;
    container: string;
    base: string;
    list: {
        base: string;
        background: string;
        rounded: string;
        shadow: string;
        padding: string;
        height: string;
        width: string;
        marker: {
            wrapper: string;
            base: string;
            background: string;
            rounded: string;
            shadow: string;
        };
        tab: {
            base: string;
            background: string;
            active: string;
            inactive: string;
            height: string;
            padding: string;
            size: string;
            font: string;
            rounded: string;
            shadow: string;
            icon: string;
        };
    };
};
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
    modelValue: {
        type: NumberConstructor;
        default: any;
    };
    orientation: {
        type: PropType<"horizontal" | "vertical">;
        default: string;
        validator: (value: string) => boolean;
    };
    defaultIndex: {
        type: NumberConstructor;
        default: number;
    };
    items: {
        type: PropType<TabItem[]>;
        default: () => any[];
    };
    unmount: {
        type: BooleanConstructor;
        default: boolean;
    };
    content: {
        type: BooleanConstructor;
        default: boolean;
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
        wrapper: string;
        container: string;
        base: string;
        list: {
            base: string;
            background: string;
            rounded: string;
            shadow: string;
            padding: string;
            height: string;
            width: string;
            marker: {
                wrapper: string;
                base: string;
                background: string;
                rounded: string;
                shadow: string;
            };
            tab: {
                base: string;
                background: string;
                active: string;
                inactive: string;
                height: string;
                padding: string;
                size: string;
                font: string;
                rounded: string;
                shadow: string;
                icon: string;
            };
        };
    }>;
    attrs: import("vue").ComputedRef<Pick<{
        [x: string]: unknown;
    }, string>>;
    listRef: import("vue").Ref<HTMLElement, HTMLElement>;
    itemRefs: import("vue").Ref<HTMLElement[], HTMLElement[]>;
    markerRef: import("vue").Ref<HTMLElement, HTMLElement>;
    selectedIndex: import("vue").Ref<number, number>;
    onChange: (index: number) => void;
}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, ("change" | "update:modelValue")[], "change" | "update:modelValue", import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    modelValue: {
        type: NumberConstructor;
        default: any;
    };
    orientation: {
        type: PropType<"horizontal" | "vertical">;
        default: string;
        validator: (value: string) => boolean;
    };
    defaultIndex: {
        type: NumberConstructor;
        default: number;
    };
    items: {
        type: PropType<TabItem[]>;
        default: () => any[];
    };
    unmount: {
        type: BooleanConstructor;
        default: boolean;
    };
    content: {
        type: BooleanConstructor;
        default: boolean;
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
    onChange?: (...args: any[]) => any;
    "onUpdate:modelValue"?: (...args: any[]) => any;
}>, {
    class: any;
    ui: {
        wrapper?: string;
        container?: string;
        base?: string;
        list?: DeepPartial<{
            base: string;
            background: string;
            rounded: string;
            shadow: string;
            padding: string;
            height: string;
            width: string;
            marker: {
                wrapper: string;
                base: string;
                background: string;
                rounded: string;
                shadow: string;
            };
            tab: {
                base: string;
                background: string;
                active: string;
                inactive: string;
                height: string;
                padding: string;
                size: string;
                font: string;
                rounded: string;
                shadow: string;
                icon: string;
            };
        }, any>;
    } & {
        [key: string]: any;
    } & {
        strategy?: Strategy;
    };
    content: boolean;
    orientation: "horizontal" | "vertical";
    modelValue: number;
    unmount: boolean;
    items: TabItem[];
    defaultIndex: number;
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
    HTabGroup: import("vue").DefineComponent<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
        selectedIndex: {
            type: NumberConstructor[];
            default: null;
        };
        defaultIndex: {
            type: NumberConstructor[];
            default: number;
        };
        vertical: {
            type: BooleanConstructor[];
            default: boolean;
        };
        manual: {
            type: BooleanConstructor[];
            default: boolean;
        };
    }, () => import("vue").VNode<import("vue").RendererNode, import("vue").RendererElement, {
        [key: string]: any;
    }>, unknown, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {
        change: (_index: number) => true;
    }, string, import("vue").VNodeProps & import("vue").AllowedComponentProps & import("vue").ComponentCustomProps, Readonly<import("vue").ExtractPropTypes<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
        selectedIndex: {
            type: NumberConstructor[];
            default: null;
        };
        defaultIndex: {
            type: NumberConstructor[];
            default: number;
        };
        vertical: {
            type: BooleanConstructor[];
            default: boolean;
        };
        manual: {
            type: BooleanConstructor[];
            default: boolean;
        };
    }>> & {
        onChange?: ((_index: number) => any) | undefined;
    }, {
        as: string | Record<string, any>;
        vertical: boolean;
        manual: boolean;
        selectedIndex: number;
        defaultIndex: number;
    }, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
    HTabList: import("vue").DefineComponent<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
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
    }>>, {
        as: string | Record<string, any>;
    }, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
    HTab: import("vue").DefineComponent<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
        disabled: {
            type: BooleanConstructor[];
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
        disabled: {
            type: BooleanConstructor[];
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
    HTabPanels: import("vue").DefineComponent<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
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
    }>>, {
        as: string | Record<string, any>;
    }, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
    HTabPanel: import("vue").DefineComponent<{
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
        tabIndex: {
            type: NumberConstructor;
            default: number;
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
        tabIndex: {
            type: NumberConstructor;
            default: number;
        };
    }>>, {
        id: string;
        tabIndex: number;
        as: string | Record<string, any>;
        unmount: boolean;
        static: boolean;
    }, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
export default _default;
