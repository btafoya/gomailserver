import type { ShallowRef, Component, InjectionKey } from 'vue';
import type { ComponentProps } from '../types/component.js';
import type { Modal, ModalState } from '../types/modal.js';
export declare const modalInjectionKey: InjectionKey<ShallowRef<ModalState>>;
declare function _useModal(): {
    open: <T extends Component>(component: T, props?: Modal & ComponentProps<T>) => void;
    close: () => Promise<void>;
    reset: () => void;
    patch: <T extends Component = {}>(props: Partial<Modal & ComponentProps<T>>) => void;
    isOpen: import("vue").Ref<boolean, boolean>;
};
export declare const useModal: typeof _useModal;
export {};
