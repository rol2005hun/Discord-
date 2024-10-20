import { writable } from 'svelte/store';

interface Toast {
    id: number;
    message: string;
    type: 'success' | 'error' | 'warning';
}

const { subscribe, update } = writable<Toast[]>([]);
let nextId = 1;

export function addToast(message: string, type: 'success' | 'error' | 'warning') {
    const id = nextId++;
    update((toasts) => [
        ...toasts,
        { id, message, type }
    ]);

    setTimeout(() => {
        removeToast(id);
    }, 5000);
}

function removeToast(id: number) {
    update((toasts) => toasts.filter((toast) => toast.id !== id));
}

export { subscribe, removeToast };