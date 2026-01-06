import type { ShallowRef, Component, InjectionKey } from 'vue';
import type { ComponentProps } from '../types/component.js';
import type { Slideover, SlideoverState } from '../types/slideover.js';
export declare const slidOverInjectionKey: InjectionKey<ShallowRef<SlideoverState>>;
declare function _useSlideover(): {
    open: <T extends Component>(component: T, props?: Slideover & ComponentProps<T>) => void;
    close: () => Promise<void>;
    reset: () => void;
    patch: <T extends Component = {}>(props: Partial<Slideover & ComponentProps<T>>) => void;
    isOpen: import("vue").Ref<boolean, boolean>;
};
export declare const useSlideover: typeof _useSlideover;
export {};
