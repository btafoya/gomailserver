import type { AvatarSize, ButtonColor, ButtonSize, ButtonVariant } from '../../types/index.js';
declare const _default: {
    wrapper: string;
    inner: string;
    title: string;
    description: string;
    descriptionOnly: string;
    actions: string;
    shadow: string;
    rounded: string;
    padding: string;
    gap: string;
    icon: {
        base: string;
    };
    avatar: {
        base: string;
        size: AvatarSize;
    };
    color: {
        white: {
            solid: string;
        };
    };
    variant: {
        solid: string;
        outline: string;
        soft: string;
        subtle: string;
    };
    default: {
        color: string;
        variant: string;
        icon: any;
        closeButton: any;
        actionButton: {
            size: ButtonSize;
            color: ButtonColor;
            variant: ButtonVariant;
        };
    };
};
export default _default;
