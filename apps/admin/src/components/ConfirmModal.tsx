import { Component, Show } from 'solid-js';

interface ConfirmModalProps {
  open: boolean;
  title: string;
  description: string;
  confirmLabel?: string;
  onConfirm: () => void;
  onCancel: () => void;
  danger?: boolean;
}

const ConfirmModal: Component<ConfirmModalProps> = (props) => {
  return (
    <Show when={props.open}>
      <div class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="fixed inset-0 bg-black/40" onClick={() => props.onCancel()} />
        <div class="relative bg-white rounded-lg shadow-xl p-6 w-full max-w-sm mx-4">
          <div class="flex items-center gap-3 mb-4">
            <div class={`w-10 h-10 rounded-full flex items-center justify-center ${props.danger !== false ? 'bg-red-100' : 'bg-gray-100'}`}>
              <svg class={`w-5 h-5 ${props.danger !== false ? 'text-red-600' : 'text-gray-600'}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 16.5c-.77.833.192 2.5 1.732 2.5z" />
              </svg>
            </div>
            <h2 class="text-lg font-semibold">{props.title}</h2>
          </div>
          <p class="text-sm text-gray-600 mb-6">{props.description}</p>
          <div class="flex justify-end gap-3">
            <button
              type="button"
              onClick={() => props.onCancel()}
              class="px-4 py-2 border rounded-md text-sm hover:bg-gray-50"
            >
              Cancel
            </button>
            <button
              type="button"
              onClick={() => props.onConfirm()}
              class={`px-4 py-2 rounded-md text-sm text-white ${props.danger !== false ? 'bg-red-600 hover:bg-red-700' : 'bg-blue-600 hover:bg-blue-700'}`}
            >
              {props.confirmLabel || 'Delete'}
            </button>
          </div>
        </div>
      </div>
    </Show>
  );
};

export default ConfirmModal;
