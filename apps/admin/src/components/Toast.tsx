import { Component, For } from 'solid-js';
import { useToast, dismissToast } from '../stores/toast';

const typeStyles: Record<string, string> = {
  success: 'bg-green-600',
  error: 'bg-red-600',
  info: 'bg-blue-600',
};

const typeIcons: Record<string, string> = {
  success: '✓',
  error: '✕',
  info: 'ⓘ',
};

const Toast: Component = () => {
  const toasts = useToast();

  return (
    <div class="fixed bottom-4 right-4 z-[100] flex flex-col gap-2 pointer-events-none">
      <For each={toasts()}>
        {(toast) => (
          <div
            class={`pointer-events-auto flex items-center gap-2 px-4 py-3 rounded-lg shadow-xl text-white text-sm animate-slide-up ${typeStyles[toast.type]}`}
          >
            <span class="text-base leading-none">{typeIcons[toast.type]}</span>
            <span>{toast.message}</span>
            <button
              onClick={() => dismissToast(toast.id)}
              class="ml-2 text-white/80 hover:text-white text-base leading-none"
            >
              ×
            </button>
          </div>
        )}
      </For>
    </div>
  );
};

export default Toast;
