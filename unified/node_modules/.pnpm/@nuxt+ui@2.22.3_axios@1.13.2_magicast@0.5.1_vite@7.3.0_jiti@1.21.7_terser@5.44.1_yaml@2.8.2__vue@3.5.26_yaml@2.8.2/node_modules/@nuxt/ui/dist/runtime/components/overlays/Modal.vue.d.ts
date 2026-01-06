import type { PropType } from 'vue';
import type { DeepPartial, Strategy } from '../../types/index.js';
declare const config: {
    wrapper: string;
    inner: string;
    container: string;
    padding: string;
    margin: string;
    base: string;
    overlay: {
        base: string;
        background: string;
        transition: {
            enter: string;
            enterFrom: string;
            enterTo: string;
            leave: string;
            leaveFrom: string;
            leaveTo: string;
        };
    };
    background: string;
    ring: string;
    rounded: string;
    shadow: string;
    width: string;
    height: string;
    fullscreen: string;
    transition: {
        enter: string;
        enterFrom: string;
        enterTo: string;
        leave: string;
        leaveFrom: string;
        leaveTo: string;
    };
};
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
    modelValue: {
        type: BooleanConstructor;
        default: boolean;
    };
    appear: {
        type: BooleanConstructor;
        default: boolean;
    };
    overlay: {
        type: BooleanConstructor;
        default: boolean;
    };
    transition: {
        type: BooleanConstructor;
        default: boolean;
    };
    preventClose: {
        type: BooleanConstructor;
        default: boolean;
    };
    fullscreen: {
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
        inner: string;
        container: string;
        padding: string;
        margin: string;
        base: string;
        overlay: {
            base: string;
            background: string;
            transition: {
                enter: string;
                enterFrom: string;
                enterTo: string;
                leave: string;
                leaveFrom: string;
                leaveTo: string;
            };
        };
        background: string;
        ring: string;
        rounded: string;
        shadow: string;
        width: string;
        height: string;
        fullscreen: string;
        transition: {
            enter: string;
            enterFrom: string;
            enterTo: string;
            leave: string;
            leaveFrom: string;
            leaveTo: string;
        };
    }>;
    attrs: import("vue").ComputedRef<Pick<{
        [x: string]: unknown;
    }, string>>;
    isOpen: import("vue").WritableComputedRef<boolean, boolean>;
    transitionClass: import("vue").ComputedRef<{
        enter: string;
        enterFrom: string;
        enterTo: string;
        leave: string;
        leaveFrom: string;
        leaveTo: string;
    }>;
    onAfterLeave: () => void;
    close: (value: boolean) => void;
}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, ("close" | "update:modelValue" | "close-prevented" | "after-leave")[], "close" | "update:modelValue" | "close-prevented" | "after-leave", import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    modelValue: {
        type: BooleanConstructor;
        default: boolean;
    };
    appear: {
        type: BooleanConstructor;
        default: boolean;
    };
    overlay: {
        type: BooleanConstructor;
        default: boolean;
    };
    transition: {
        type: BooleanConstructor;
        default: boolean;
    };
    preventClose: {
        type: BooleanConstructor;
        default: boolean;
    };
    fullscreen: {
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
    onClose?: (...args: any[]) => any;
    "onUpdate:modelValue"?: (...args: any[]) => any;
    "onClose-prevented"?: (...args: any[]) => any;
    "onAfter-leave"?: (...args: any[]) => any;
}>, {
    class: any;
    transition: boolean;
    ui: {
        wrapper?: string;
        inner?: string;
        container?: string;
        padding?: string;
        margin?: string;
        base?: string;
        overlay?: DeepPartial<{
            base: string;
            background: string;
            transition: {
                enter: string;
                enterFrom: string;
                enterTo: string;
                leave: string;
                leaveFrom: string;
                leaveTo: string;
            };
        }, any>;
        background?: string;
        ring?: string;
        rounded?: string;
        shadow?: string;
        width?: string;
        height?: string;
        fullscreen?: string;
        transition?: DeepPartial<{
            enter: string;
            enterFrom: string;
            enterTo: string;
            leave: string;
            leaveFrom: string;
            leaveTo: string;
        }, any>;
    } & {
        [key: string]: any;
    } & {
        strategy?: Strategy;
    };
    appear: boolean;
    overlay: boolean;
    preventClose: boolean;
    fullscreen: boolean;
    modelValue: boolean;
}, {}, {
    HDialog: import("vue").DefineComponent<{
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
        open: {
            type: (BooleanConstructor | StringConstructor)[];
            default: string;
        };
        initialFocus: {
            type: PropType<HTMLElement | null>;
            default: null;
        };
        id: {
            type: StringConstructor;
            default: null;
        };
        role: {
            type: PropType<"dialog" | "alertdialog">;
            default: string;
        };
    }, () => import("vue").VNode<import("vue").RendererNode, import("vue").RendererElement, {
        [key: string]: any;
    }>, unknown, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {
        close: (_close: boolean) => true;
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
        open: {
            type: (BooleanConstructor | StringConstructor)[];
            default: string;
        };
        initialFocus: {
            type: PropType<HTMLElement | null>;
            default: null;
        };
        id: {
            type: StringConstructor;
            default: null;
        };
        role: {
            type: PropType<"dialog" | "alertdialog">;
            default: string;
        };
    }>> & {
        onClose?: ((_close: boolean) => any) | undefined;
    }, {
        id: string;
        role: "dialog" | "alertdialog";
        as: string | Record<string, any>;
        unmount: boolean;
        static: boolean;
        open: string | boolean;
        initialFocus: HTMLElement | null;
    }, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
    HDialogPanel: import("vue").DefineComponent<{
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
    TransitionRoot: import("vue").DefineComponent<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
        show: {
            type: BooleanConstructor[];
            default: null;
        };
        unmount: {
            type: BooleanConstructor[];
            default: boolean;
        };
        appear: {
            type: BooleanConstructor[];
            default: boolean;
        };
        enter: {
            type: StringConstructor[];
            default: string;
        };
        enterFrom: {
            type: StringConstructor[];
            default: string;
        };
        enterTo: {
            type: StringConstructor[];
            default: string;
        };
        entered: {
            type: StringConstructor[];
            default: string;
        };
        leave: {
            type: StringConstructor[];
            default: string;
        };
        leaveFrom: {
            type: StringConstructor[];
            default: string;
        };
        leaveTo: {
            type: StringConstructor[];
            default: string;
        };
    }, () => import("vue").VNode<import("vue").RendererNode, import("vue").RendererElement, {
        [key: string]: any;
    }> | import("vue").VNode<import("vue").RendererNode, import("vue").RendererElement, {
        [key: string]: any;
    }>[] | null, unknown, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {
        beforeEnter: () => true;
        afterEnter: () => true;
        beforeLeave: () => true;
        afterLeave: () => true;
    }, string, import("vue").VNodeProps & import("vue").AllowedComponentProps & import("vue").ComponentCustomProps, Readonly<import("vue").ExtractPropTypes<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
        show: {
            type: BooleanConstructor[];
            default: null;
        };
        unmount: {
            type: BooleanConstructor[];
            default: boolean;
        };
        appear: {
            type: BooleanConstructor[];
            default: boolean;
        };
        enter: {
            type: StringConstructor[];
            default: string;
        };
        enterFrom: {
            type: StringConstructor[];
            default: string;
        };
        enterTo: {
            type: StringConstructor[];
            default: string;
        };
        entered: {
            type: StringConstructor[];
            default: string;
        };
        leave: {
            type: StringConstructor[];
            default: string;
        };
        leaveFrom: {
            type: StringConstructor[];
            default: string;
        };
        leaveTo: {
            type: StringConstructor[];
            default: string;
        };
    }>> & {
        onBeforeEnter?: (() => any) | undefined;
        onAfterEnter?: (() => any) | undefined;
        onBeforeLeave?: (() => any) | undefined;
        onAfterLeave?: (() => any) | undefined;
    }, {
        as: string | Record<string, any>;
        unmount: boolean;
        show: boolean;
        appear: boolean;
        enter: string;
        enterFrom: string;
        enterTo: string;
        entered: string;
        leave: string;
        leaveFrom: string;
        leaveTo: string;
    }, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
    TransitionChild: import("vue").DefineComponent<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
        show: {
            type: BooleanConstructor[];
            default: null;
        };
        unmount: {
            type: BooleanConstructor[];
            default: boolean;
        };
        appear: {
            type: BooleanConstructor[];
            default: boolean;
        };
        enter: {
            type: StringConstructor[];
            default: string;
        };
        enterFrom: {
            type: StringConstructor[];
            default: string;
        };
        enterTo: {
            type: StringConstructor[];
            default: string;
        };
        entered: {
            type: StringConstructor[];
            default: string;
        };
        leave: {
            type: StringConstructor[];
            default: string;
        };
        leaveFrom: {
            type: StringConstructor[];
            default: string;
        };
        leaveTo: {
            type: StringConstructor[];
            default: string;
        };
    }, () => import("vue").VNode<import("vue").RendererNode, import("vue").RendererElement, {
        [key: string]: any;
    }> | import("vue").VNode<import("vue").RendererNode, import("vue").RendererElement, {
        [key: string]: any;
    }>[] | null, unknown, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {
        beforeEnter: () => true;
        afterEnter: () => true;
        beforeLeave: () => true;
        afterLeave: () => true;
    }, string, import("vue").VNodeProps & import("vue").AllowedComponentProps & import("vue").ComponentCustomProps, Readonly<import("vue").ExtractPropTypes<{
        as: {
            type: (ObjectConstructor | StringConstructor)[];
            default: string;
        };
        show: {
            type: BooleanConstructor[];
            default: null;
        };
        unmount: {
            type: BooleanConstructor[];
            default: boolean;
        };
        appear: {
            type: BooleanConstructor[];
            default: boolean;
        };
        enter: {
            type: StringConstructor[];
            default: string;
        };
        enterFrom: {
            type: StringConstructor[];
            default: string;
        };
        enterTo: {
            type: StringConstructor[];
            default: string;
        };
        entered: {
            type: StringConstructor[];
            default: string;
        };
        leave: {
            type: StringConstructor[];
            default: string;
        };
        leaveFrom: {
            type: StringConstructor[];
            default: string;
        };
        leaveTo: {
            type: StringConstructor[];
            default: string;
        };
    }>> & {
        onBeforeEnter?: (() => any) | undefined;
        onAfterEnter?: (() => any) | undefined;
        onBeforeLeave?: (() => any) | undefined;
        onAfterLeave?: (() => any) | undefined;
    }, {
        as: string | Record<string, any>;
        unmount: boolean;
        show: boolean;
        appear: boolean;
        enter: string;
        enterFrom: string;
        enterTo: string;
        entered: string;
        leave: string;
        leaveFrom: string;
        leaveTo: string;
    }, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
export default _default;
