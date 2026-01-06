import type { PropType } from 'vue';
import type { Strategy, ProgressSize, ProgressAnimation, ProgressColor, DeepPartial } from '../../types/index.js';
declare const config: {
    wrapper: string;
    indicator: {
        container: {
            base: string;
            width: string;
            transition: string;
        };
        align: string;
        width: string;
        color: string;
        size: {
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
            '2xl': string;
        };
    };
    progress: {
        base: string;
        width: string;
        size: {
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
            '2xl': string;
        };
        rounded: string;
        track: string;
        bar: string;
        color: string;
        background: string;
        indeterminate: {
            base: string;
            rounded: string;
        };
    };
    steps: {
        base: string;
        color: string;
        size: {
            '2xs': string;
            xs: string;
            sm: string;
            md: string;
            lg: string;
            xl: string;
            '2xl': string;
        };
    };
    step: {
        base: string;
        align: string;
        active: string;
        first: string;
    };
    animation: {
        carousel: string;
        'carousel-inverse': string;
        swing: string;
        elastic: string;
    };
    default: {
        color: string;
        size: string;
        animation: string;
    };
};
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
    value: {
        type: NumberConstructor;
        default: any;
    };
    max: {
        type: (NumberConstructor | {
            (arrayLength: number): any[];
            (...items: any[]): any[];
            new (arrayLength: number): any[];
            new (...items: any[]): any[];
            isArray(arg: any): arg is any[];
            readonly prototype: any[];
            from<T>(arrayLike: ArrayLike<T>): T[];
            from<T, U>(arrayLike: ArrayLike<T>, mapfn: (v: T, k: number) => U, thisArg?: any): U[];
            from<T>(iterable: Iterable<T> | ArrayLike<T>): T[];
            from<T, U>(iterable: Iterable<T> | ArrayLike<T>, mapfn: (v: T, k: number) => U, thisArg?: any): U[];
            of<T>(...items: T[]): T[];
            fromAsync<T>(iterableOrArrayLike: AsyncIterable<T> | Iterable<T | PromiseLike<T>> | ArrayLike<T | PromiseLike<T>>): Promise<T[]>;
            fromAsync<T, U>(iterableOrArrayLike: AsyncIterable<T> | Iterable<T> | ArrayLike<T>, mapFn: (value: Awaited<T>, index: number) => U, thisArg?: any): Promise<Awaited<U>[]>;
            readonly [Symbol.species]: ArrayConstructor;
        })[];
        default: number;
    };
    indicator: {
        type: BooleanConstructor;
        default: boolean;
    };
    animation: {
        type: PropType<ProgressAnimation>;
        default: () => string;
        validator(value: string): boolean;
    };
    size: {
        type: PropType<ProgressSize>;
        default: () => string;
        validator(value: string): boolean;
    };
    color: {
        type: PropType<ProgressColor>;
        default: () => string;
        validator(value: string): any;
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
    ui: import("vue").ComputedRef<{
        wrapper: string;
        indicator: {
            container: {
                base: string;
                width: string;
                transition: string;
            };
            align: string;
            width: string;
            color: string;
            size: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
                '2xl': string;
            };
        };
        progress: {
            base: string;
            width: string;
            size: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
                '2xl': string;
            };
            rounded: string;
            track: string;
            bar: string;
            color: string;
            background: string;
            indeterminate: {
                base: string;
                rounded: string;
            };
        };
        steps: {
            base: string;
            color: string;
            size: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
                '2xl': string;
            };
        };
        step: {
            base: string;
            align: string;
            active: string;
            first: string;
        };
        animation: {
            carousel: string;
            'carousel-inverse': string;
            swing: string;
            elastic: string;
        };
        default: {
            color: string;
            size: string;
            animation: string;
        };
    }>;
    attrs: import("vue").ComputedRef<Pick<{
        [x: string]: unknown;
    }, string>>;
    indicatorContainerClass: import("vue").ComputedRef<string>;
    indicatorClass: import("vue").ComputedRef<string>;
    progressClass: import("vue").ComputedRef<string>;
    stepsClass: import("vue").ComputedRef<string>;
    stepClasses: (index: string | number) => string;
    isIndeterminate: import("vue").ComputedRef<boolean>;
    isSteps: import("vue").ComputedRef<boolean>;
    realMax: import("vue").ComputedRef<number>;
    percent: import("vue").ComputedRef<number>;
}, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    value: {
        type: NumberConstructor;
        default: any;
    };
    max: {
        type: (NumberConstructor | {
            (arrayLength: number): any[];
            (...items: any[]): any[];
            new (arrayLength: number): any[];
            new (...items: any[]): any[];
            isArray(arg: any): arg is any[];
            readonly prototype: any[];
            from<T>(arrayLike: ArrayLike<T>): T[];
            from<T, U>(arrayLike: ArrayLike<T>, mapfn: (v: T, k: number) => U, thisArg?: any): U[];
            from<T>(iterable: Iterable<T> | ArrayLike<T>): T[];
            from<T, U>(iterable: Iterable<T> | ArrayLike<T>, mapfn: (v: T, k: number) => U, thisArg?: any): U[];
            of<T>(...items: T[]): T[];
            fromAsync<T>(iterableOrArrayLike: AsyncIterable<T> | Iterable<T | PromiseLike<T>> | ArrayLike<T | PromiseLike<T>>): Promise<T[]>;
            fromAsync<T, U>(iterableOrArrayLike: AsyncIterable<T> | Iterable<T> | ArrayLike<T>, mapFn: (value: Awaited<T>, index: number) => U, thisArg?: any): Promise<Awaited<U>[]>;
            readonly [Symbol.species]: ArrayConstructor;
        })[];
        default: number;
    };
    indicator: {
        type: BooleanConstructor;
        default: boolean;
    };
    animation: {
        type: PropType<ProgressAnimation>;
        default: () => string;
        validator(value: string): boolean;
    };
    size: {
        type: PropType<ProgressSize>;
        default: () => string;
        validator(value: string): boolean;
    };
    color: {
        type: PropType<ProgressColor>;
        default: () => string;
        validator(value: string): any;
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
    size: "sm" | "xs" | "2xs" | "md" | "lg" | "xl" | "2xl";
    class: any;
    ui: {
        wrapper?: string;
        indicator?: DeepPartial<{
            container: {
                base: string;
                width: string;
                transition: string;
            };
            align: string;
            width: string;
            color: string;
            size: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
                '2xl': string;
            };
        }, any>;
        progress?: DeepPartial<{
            base: string;
            width: string;
            size: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
                '2xl': string;
            };
            rounded: string;
            track: string;
            bar: string;
            color: string;
            background: string;
            indeterminate: {
                base: string;
                rounded: string;
            };
        }, any>;
        steps?: DeepPartial<{
            base: string;
            color: string;
            size: {
                '2xs': string;
                xs: string;
                sm: string;
                md: string;
                lg: string;
                xl: string;
                '2xl': string;
            };
        }, any>;
        step?: DeepPartial<{
            base: string;
            align: string;
            active: string;
            first: string;
        }, any>;
        animation?: DeepPartial<{
            carousel: string;
            'carousel-inverse': string;
            swing: string;
            elastic: string;
        }, any>;
        default?: DeepPartial<{
            color: string;
            size: string;
            animation: string;
        }, any>;
    } & {
        [key: string]: any;
    } & {
        strategy?: Strategy;
    };
    color: "primary" | "red" | "orange" | "amber" | "yellow" | "lime" | "green" | "emerald" | "teal" | "cyan" | "sky" | "blue" | "indigo" | "violet" | "purple" | "fuchsia" | "pink" | "rose";
    value: number;
    indicator: boolean;
    animation: "carousel" | "carousel-inverse" | "swing" | "elastic";
    max: number | any[];
}, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
export default _default;
