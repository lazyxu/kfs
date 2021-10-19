export const isElectron = typeof window.require === 'function';
export const isBrowser = !isElectron;
