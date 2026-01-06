import type { AvatarSize, ButtonColor, ButtonSize, ButtonVariant } from '../../types/index.js';
declare const _default: {
    wrapper: string;
    container: string;
    inner: string;
    title: string;
    description: string;
    descriptionOnly: string;
    actions: string;
    background: string;
    shadow: string;
    rounded: string;
    padding: string;
    gap: string;
    ring: string;
    icon: {
        base: string;
        color: string;
    };
    avatar: {
        base: string;
        size: AvatarSize;
    };
    progress: {
        base: string;
        background: string;
    };
    transition: {
        enterActiveClass: string;
        enterFromClass: string;
        enterToClass: string;
        leaveActiveClass: string;
        leaveFromClass: string;
        leaveToClass: string;
    };
    default: {
        color: string;
        icon: any;
        timeout: number;
        closeButton: {
            icon: string;
            color: ButtonColor;
            variant: ButtonVariant;
            padded: boolean;
        };
        actionButton: {
            size: ButtonSize;
            color: ButtonColor;
        };
    };
};
export default _default;
