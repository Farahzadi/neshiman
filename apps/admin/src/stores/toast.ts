import { createSignal, createMemo } from 'solid-js';

type ToastType = 'success' | 'error' | 'info';

interface Toast {
  id: number;
  message: string;
  type: ToastType;
}

let nextId = 0;
const [toasts, setToasts] = createSignal<Toast[]>([]);

export function showToast(message: string, type: ToastType = 'info', duration = 4000) {
  const id = nextId++;
  setToasts((prev) => [...prev, { id, message, type }]);
  setTimeout(() => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, duration);
}

export function useToast() {
  return createMemo(() => toasts());
}

export function dismissToast(id: number) {
  setToasts((prev) => prev.filter((t) => t.id !== id));
}
