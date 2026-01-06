import type { PropType } from 'vue';
import type { AvatarSize, DeepPartial, Strategy } from '../../types/index.js';
declare const avatarGroupConfig: {
    wrapper: string;
    ring: string;
    margin: string;
};
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
    size: {
        type: PropType<AvatarSize>;
        default: any;
        validator(value: string): boolean;
    };
    max: {
        type: NumberConstructor;
        default: any;
    };
    class: {
        type: PropType<any>;
        default: () => "";
    };
    ui: {
        type: PropType<DeepPartial<typeof avatarGroupConfig> & {
            strategy?: Strategy;
        }>;
        default: () => {};
    };
}>, () => import("vue").VNode<import("vue").RendererNode, import("vue").RendererElement, {
    [key: string]: any;
}>, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    size: {
        type: PropType<AvatarSize>;
        default: any;
        validator(value: string): boolean;
    };
    max: {
        type: NumberConstructor;
        default: any;
    };
    class: {
        type: PropType<any>;
        default: () => "";
    };
    ui: {
        type: PropType<DeepPartial<typeof avatarGroupConfig> & {
            strategy?: Strategy;
        }>;
        default: () => {};
    };
}>> & Readonly<{}>, {
    ui: {
        wrapper?: string;
        ring?: string;
        margin?: string;
    } & {
        [key: string]: any;
    } & {
        strategy?: Strategy;
    };
    size: AvatarSize;
    class: any;
    max: number;
}, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
export default _default;
