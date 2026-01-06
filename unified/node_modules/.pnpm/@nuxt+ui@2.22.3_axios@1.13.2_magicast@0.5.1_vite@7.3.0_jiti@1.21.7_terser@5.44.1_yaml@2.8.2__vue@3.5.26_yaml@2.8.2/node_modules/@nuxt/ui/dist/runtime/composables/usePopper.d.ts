import type { Ref } from 'vue';
import type { VirtualElement } from '@popperjs/core/lib/popper-lite';
import type { Instance } from '@popperjs/core';
import type { MaybeElement } from '@vueuse/core';
import type { PopperOptions } from '../types/popper.js';
export declare const createPopper: <TModifier extends Partial<import("@popperjs/core").Modifier<any, any>>>(reference: Element | VirtualElement, popper: HTMLElement, options?: Partial<import("@popperjs/core").OptionsGeneric<TModifier>>) => Instance;
export declare function usePopper({ locked, overflowPadding, offsetDistance, offsetSkid, gpuAcceleration, adaptive, scroll, resize, arrow, placement, strategy }: PopperOptions, virtualReference?: Ref<Element | VirtualElement>): readonly [Ref<MaybeElement, MaybeElement>, Ref<MaybeElement, MaybeElement>, Ref<{
    state: {
        elements: {
            reference: Element | {
                getBoundingClientRect: () => ClientRect | DOMRect;
                contextElement?: Element;
            };
            popper: HTMLElement;
            arrow?: HTMLElement;
        };
        options: {
            placement: import("@popperjs/core").Placement;
            modifiers: any[];
            strategy: import("@popperjs/core").PositioningStrategy;
            onFirstUpdate?: (arg0: Partial<import("@popperjs/core").State>) => void;
        };
        placement: import("@popperjs/core").Placement;
        strategy: import("@popperjs/core").PositioningStrategy;
        orderedModifiers: {
            name: any;
            enabled: boolean;
            phase: import("@popperjs/core").ModifierPhases;
            requires?: Array<string>;
            requiresIfExists?: Array<string>;
            fn: (arg0: import("@popperjs/core").ModifierArguments<any>) => import("@popperjs/core").State | void;
            effect?: (arg0: import("@popperjs/core").ModifierArguments<any>) => (() => void) | void;
            options?: Partial<any>;
            data?: import("@popperjs/core").Obj;
        }[];
        rects: {
            reference: {
                width: number;
                height: number;
                x: number;
                y: number;
            };
            popper: {
                width: number;
                height: number;
                x: number;
                y: number;
            };
        };
        scrollParents: {
            reference: (Element | {
                innerHeight: number;
                offsetHeight: number;
                innerWidth: number;
                offsetWidth: number;
                pageXOffset: number;
                pageYOffset: number;
                getComputedStyle: typeof getComputedStyle;
                addEventListener: (type: any, listener: any, optionsOrUseCapture?: any) => void;
                removeEventListener: (type: any, listener: any, optionsOrUseCapture?: any) => void;
                Element: Element;
                HTMLElement: HTMLElement;
                Node: Node;
                toString: () => "[object Window]";
                devicePixelRatio: number;
                visualViewport?: {
                    addEventListener: {
                        (type: string, callback: EventListenerOrEventListenerObject | null, options?: AddEventListenerOptions | boolean): void;
                        (type: string, callback: EventListenerOrEventListenerObject | null, options?: AddEventListenerOptions | boolean): void;
                    };
                    dispatchEvent: {
                        (event: Event): boolean;
                        (event: Event): boolean;
                    };
                    removeEventListener: {
                        (type: string, callback: EventListenerOrEventListenerObject | null, options?: EventListenerOptions | boolean): void;
                        (type: string, callback: EventListenerOrEventListenerObject | null, options?: EventListenerOptions | boolean): void;
                    };
                    width: number;
                    height: number;
                    offsetLeft: number;
                    offsetTop: number;
                    scale: number;
                };
                ShadowRoot: ShadowRoot;
            } | {
                addEventListener: {
                    (type: string, callback: EventListenerOrEventListenerObject | null, options?: AddEventListenerOptions | boolean): void;
                    (type: string, callback: EventListenerOrEventListenerObject | null, options?: AddEventListenerOptions | boolean): void;
                };
                dispatchEvent: {
                    (event: Event): boolean;
                    (event: Event): boolean;
                };
                removeEventListener: {
                    (type: string, callback: EventListenerOrEventListenerObject | null, options?: EventListenerOptions | boolean): void;
                    (type: string, callback: EventListenerOrEventListenerObject | null, options?: EventListenerOptions | boolean): void;
                };
                width: number;
                height: number;
                offsetLeft: number;
                offsetTop: number;
                scale: number;
            })[];
            popper: (Element | {
                innerHeight: number;
                offsetHeight: number;
                innerWidth: number;
                offsetWidth: number;
                pageXOffset: number;
                pageYOffset: number;
                getComputedStyle: typeof getComputedStyle;
                addEventListener: (type: any, listener: any, optionsOrUseCapture?: any) => void;
                removeEventListener: (type: any, listener: any, optionsOrUseCapture?: any) => void;
                Element: Element;
                HTMLElement: HTMLElement;
                Node: Node;
                toString: () => "[object Window]";
                devicePixelRatio: number;
                visualViewport?: {
                    addEventListener: {
                        (type: string, callback: EventListenerOrEventListenerObject | null, options?: AddEventListenerOptions | boolean): void;
                        (type: string, callback: EventListenerOrEventListenerObject | null, options?: AddEventListenerOptions | boolean): void;
                    };
                    dispatchEvent: {
                        (event: Event): boolean;
                        (event: Event): boolean;
                    };
                    removeEventListener: {
                        (type: string, callback: EventListenerOrEventListenerObject | null, options?: EventListenerOptions | boolean): void;
                        (type: string, callback: EventListenerOrEventListenerObject | null, options?: EventListenerOptions | boolean): void;
                    };
                    width: number;
                    height: number;
                    offsetLeft: number;
                    offsetTop: number;
                    scale: number;
                };
                ShadowRoot: ShadowRoot;
            } | {
                addEventListener: {
                    (type: string, callback: EventListenerOrEventListenerObject | null, options?: AddEventListenerOptions | boolean): void;
                    (type: string, callback: EventListenerOrEventListenerObject | null, options?: AddEventListenerOptions | boolean): void;
                };
                dispatchEvent: {
                    (event: Event): boolean;
                    (event: Event): boolean;
                };
                removeEventListener: {
                    (type: string, callback: EventListenerOrEventListenerObject | null, options?: EventListenerOptions | boolean): void;
                    (type: string, callback: EventListenerOrEventListenerObject | null, options?: EventListenerOptions | boolean): void;
                };
                width: number;
                height: number;
                offsetLeft: number;
                offsetTop: number;
                scale: number;
            })[];
        };
        styles: {
            [key: string]: Partial<CSSStyleDeclaration>;
        };
        attributes: {
            [key: string]: {
                [key: string]: string | boolean;
            };
        };
        modifiersData: {
            [x: string]: any;
            arrow?: {
                x?: number;
                y?: number;
                centerOffset: number;
            };
            hide?: {
                isReferenceHidden: boolean;
                hasPopperEscaped: boolean;
                referenceClippingOffsets: {
                    top: number;
                    left: number;
                    right: number;
                    bottom: number;
                };
                popperEscapeOffsets: {
                    top: number;
                    left: number;
                    right: number;
                    bottom: number;
                };
            };
            offset?: {
                "bottom-end"?: {
                    y: number;
                    x: number;
                };
                "bottom-start"?: {
                    y: number;
                    x: number;
                };
                bottom?: {
                    y: number;
                    x: number;
                };
                left?: {
                    y: number;
                    x: number;
                };
                right?: {
                    y: number;
                    x: number;
                };
                top?: {
                    y: number;
                    x: number;
                };
                auto?: {
                    y: number;
                    x: number;
                };
                "auto-start"?: {
                    y: number;
                    x: number;
                };
                "auto-end"?: {
                    y: number;
                    x: number;
                };
                "top-start"?: {
                    y: number;
                    x: number;
                };
                "top-end"?: {
                    y: number;
                    x: number;
                };
                "right-start"?: {
                    y: number;
                    x: number;
                };
                "right-end"?: {
                    y: number;
                    x: number;
                };
                "left-start"?: {
                    y: number;
                    x: number;
                };
                "left-end"?: {
                    y: number;
                    x: number;
                };
            };
            preventOverflow?: {
                y: number;
                x: number;
            };
            popperOffsets?: {
                y: number;
                x: number;
            };
        };
        reset: boolean;
    };
    destroy: () => void;
    forceUpdate: () => void;
    update: () => Promise<Partial<import("@popperjs/core").State>>;
    setOptions: (setOptionsAction: Partial<import("@popperjs/core").OptionsGeneric<any>> | ((prev: Partial<import("@popperjs/core").OptionsGeneric<any>>) => Partial<import("@popperjs/core").OptionsGeneric<any>>)) => Promise<Partial<import("@popperjs/core").State>>;
}, Instance | {
    state: {
        elements: {
            reference: Element | {
                getBoundingClientRect: () => ClientRect | DOMRect;
                contextElement?: Element;
            };
            popper: HTMLElement;
            arrow?: HTMLElement;
        };
        options: {
            placement: import("@popperjs/core").Placement;
            modifiers: any[];
            strategy: import("@popperjs/core").PositioningStrategy;
            onFirstUpdate?: (arg0: Partial<import("@popperjs/core").State>) => void;
        };
        placement: import("@popperjs/core").Placement;
        strategy: import("@popperjs/core").PositioningStrategy;
        orderedModifiers: {
            name: any;
            enabled: boolean;
            phase: import("@popperjs/core").ModifierPhases;
            requires?: Array<string>;
            requiresIfExists?: Array<string>;
            fn: (arg0: import("@popperjs/core").ModifierArguments<any>) => import("@popperjs/core").State | void;
            effect?: (arg0: import("@popperjs/core").ModifierArguments<any>) => (() => void) | void;
            options?: Partial<any>;
            data?: import("@popperjs/core").Obj;
        }[];
        rects: {
            reference: {
                width: number;
                height: number;
                x: number;
                y: number;
            };
            popper: {
                width: number;
                height: number;
                x: number;
                y: number;
            };
        };
        scrollParents: {
            reference: (Element | {
                innerHeight: number;
                offsetHeight: number;
                innerWidth: number;
                offsetWidth: number;
                pageXOffset: number;
                pageYOffset: number;
                getComputedStyle: typeof getComputedStyle;
                addEventListener: (type: any, listener: any, optionsOrUseCapture?: any) => void;
                removeEventListener: (type: any, listener: any, optionsOrUseCapture?: any) => void;
                Element: Element;
                HTMLElement: HTMLElement;
                Node: Node;
                toString: () => "[object Window]";
                devicePixelRatio: number;
                visualViewport?: {
                    addEventListener: {
                        (type: string, callback: EventListenerOrEventListenerObject | null, options?: AddEventListenerOptions | boolean): void;
                        (type: string, callback: EventListenerOrEventListenerObject | null, options?: AddEventListenerOptions | boolean): void;
                    };
                    dispatchEvent: {
                        (event: Event): boolean;
                        (event: Event): boolean;
                    };
                    removeEventListener: {
                        (type: string, callback: EventListenerOrEventListenerObject | null, options?: EventListenerOptions | boolean): void;
                        (type: string, callback: EventListenerOrEventListenerObject | null, options?: EventListenerOptions | boolean): void;
                    };
                    width: number;
                    height: number;
                    offsetLeft: number;
                    offsetTop: number;
                    scale: number;
                };
                ShadowRoot: ShadowRoot;
            } | {
                addEventListener: {
                    (type: string, callback: EventListenerOrEventListenerObject | null, options?: AddEventListenerOptions | boolean): void;
                    (type: string, callback: EventListenerOrEventListenerObject | null, options?: AddEventListenerOptions | boolean): void;
                };
                dispatchEvent: {
                    (event: Event): boolean;
                    (event: Event): boolean;
                };
                removeEventListener: {
                    (type: string, callback: EventListenerOrEventListenerObject | null, options?: EventListenerOptions | boolean): void;
                    (type: string, callback: EventListenerOrEventListenerObject | null, options?: EventListenerOptions | boolean): void;
                };
                width: number;
                height: number;
                offsetLeft: number;
                offsetTop: number;
                scale: number;
            })[];
            popper: (Element | {
                innerHeight: number;
                offsetHeight: number;
                innerWidth: number;
                offsetWidth: number;
                pageXOffset: number;
                pageYOffset: number;
                getComputedStyle: typeof getComputedStyle;
                addEventListener: (type: any, listener: any, optionsOrUseCapture?: any) => void;
                removeEventListener: (type: any, listener: any, optionsOrUseCapture?: any) => void;
                Element: Element;
                HTMLElement: HTMLElement;
                Node: Node;
                toString: () => "[object Window]";
                devicePixelRatio: number;
                visualViewport?: {
                    addEventListener: {
                        (type: string, callback: EventListenerOrEventListenerObject | null, options?: AddEventListenerOptions | boolean): void;
                        (type: string, callback: EventListenerOrEventListenerObject | null, options?: AddEventListenerOptions | boolean): void;
                    };
                    dispatchEvent: {
                        (event: Event): boolean;
                        (event: Event): boolean;
                    };
                    removeEventListener: {
                        (type: string, callback: EventListenerOrEventListenerObject | null, options?: EventListenerOptions | boolean): void;
                        (type: string, callback: EventListenerOrEventListenerObject | null, options?: EventListenerOptions | boolean): void;
                    };
                    width: number;
                    height: number;
                    offsetLeft: number;
                    offsetTop: number;
                    scale: number;
                };
                ShadowRoot: ShadowRoot;
            } | {
                addEventListener: {
                    (type: string, callback: EventListenerOrEventListenerObject | null, options?: AddEventListenerOptions | boolean): void;
                    (type: string, callback: EventListenerOrEventListenerObject | null, options?: AddEventListenerOptions | boolean): void;
                };
                dispatchEvent: {
                    (event: Event): boolean;
                    (event: Event): boolean;
                };
                removeEventListener: {
                    (type: string, callback: EventListenerOrEventListenerObject | null, options?: EventListenerOptions | boolean): void;
                    (type: string, callback: EventListenerOrEventListenerObject | null, options?: EventListenerOptions | boolean): void;
                };
                width: number;
                height: number;
                offsetLeft: number;
                offsetTop: number;
                scale: number;
            })[];
        };
        styles: {
            [key: string]: Partial<CSSStyleDeclaration>;
        };
        attributes: {
            [key: string]: {
                [key: string]: string | boolean;
            };
        };
        modifiersData: {
            [x: string]: any;
            arrow?: {
                x?: number;
                y?: number;
                centerOffset: number;
            };
            hide?: {
                isReferenceHidden: boolean;
                hasPopperEscaped: boolean;
                referenceClippingOffsets: {
                    top: number;
                    left: number;
                    right: number;
                    bottom: number;
                };
                popperEscapeOffsets: {
                    top: number;
                    left: number;
                    right: number;
                    bottom: number;
                };
            };
            offset?: {
                "bottom-end"?: {
                    y: number;
                    x: number;
                };
                "bottom-start"?: {
                    y: number;
                    x: number;
                };
                bottom?: {
                    y: number;
                    x: number;
                };
                left?: {
                    y: number;
                    x: number;
                };
                right?: {
                    y: number;
                    x: number;
                };
                top?: {
                    y: number;
                    x: number;
                };
                auto?: {
                    y: number;
                    x: number;
                };
                "auto-start"?: {
                    y: number;
                    x: number;
                };
                "auto-end"?: {
                    y: number;
                    x: number;
                };
                "top-start"?: {
                    y: number;
                    x: number;
                };
                "top-end"?: {
                    y: number;
                    x: number;
                };
                "right-start"?: {
                    y: number;
                    x: number;
                };
                "right-end"?: {
                    y: number;
                    x: number;
                };
                "left-start"?: {
                    y: number;
                    x: number;
                };
                "left-end"?: {
                    y: number;
                    x: number;
                };
            };
            preventOverflow?: {
                y: number;
                x: number;
            };
            popperOffsets?: {
                y: number;
                x: number;
            };
        };
        reset: boolean;
    };
    destroy: () => void;
    forceUpdate: () => void;
    update: () => Promise<Partial<import("@popperjs/core").State>>;
    setOptions: (setOptionsAction: Partial<import("@popperjs/core").OptionsGeneric<any>> | ((prev: Partial<import("@popperjs/core").OptionsGeneric<any>>) => Partial<import("@popperjs/core").OptionsGeneric<any>>)) => Promise<Partial<import("@popperjs/core").State>>;
}>];
