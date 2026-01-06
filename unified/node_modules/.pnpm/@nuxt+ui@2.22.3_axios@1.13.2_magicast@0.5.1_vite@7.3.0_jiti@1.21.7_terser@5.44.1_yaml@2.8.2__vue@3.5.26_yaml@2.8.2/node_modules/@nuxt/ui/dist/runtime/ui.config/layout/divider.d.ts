import type { AvatarSize } from '../../types/index.js';
declare const _default: {
    wrapper: {
        base: string;
        horizontal: string;
        vertical: string;
    };
    container: {
        base: string;
        horizontal: string;
        vertical: string;
    };
    border: {
        base: string;
        horizontal: string;
        vertical: string;
        size: {
            horizontal: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
            };
            vertical: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
            };
        };
        type: {
            solid: string;
            dotted: string;
            dashed: string;
        };
    };
    icon: {
        base: string;
    };
    avatar: {
        base: string;
        size: AvatarSize;
    };
    label: string;
    default: {
        size: string;
        type: string;
    };
};
export default _default;
