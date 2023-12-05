import { createGlobalStore } from 'hox';
import { useState } from 'react';

export const [useWindows, getWindows] = createGlobalStore(() => useState({}));

export default useWindows;

let windowId = 0;

export function newId() {
  return ++windowId;
}

export function newWindow(setWindows, app, props) {
  setWindows(prev => {
    let id = newId();
    return { ...prev, [id]: { id, app, props } }
  })
}

export function closeWindow(setWindows, id) {
  setWindows(prev => {
    let windows = { ...prev };
    console.log(windows, id);
    delete windows[id];
    return windows;
  })
}

export const APP_IMAGE_VIEWER = "APP_IMAGE_VIEWER";
export const APP_VIDEO_VIEWER = "APP_VIDEO_VIEWER";
export const APP_TEXT_VIEWER = "APP_TEXT_VIEWER";
