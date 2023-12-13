import { noteWarning } from '@kfs/common/components/Notification/Notification';
import { createGlobalStore } from 'hox';
import { useState } from 'react';

export const [useWindows, useWindows2] = createGlobalStore(() => useState({}));

export default useWindows;

export const getWindows = () => {
  return useWindows2()[0];
}

let windowId = 0;

export function newId() {
  return ++windowId;
}

export function newWindow(setWindows, app, props) {
  if (!app) {
    noteWarning("不支持打开该文件");
    return;
  }
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

const appExtensions = {
  APP_TEXT_VIEWER: [".txt", ".md", ".log",
    ".go", ".js", ".jsx", ".java", ".c", ".h", ".c++", ".cpp"
  ],
  APP_IMAGE_VIEWER: [".jpg", ".jpeg", ".png", ".heic"],
  APP_VIDEO_VIEWER: [".mp4", ".mov"],
}

const defaultApp = {
}

let openAppByExt = {}

export function refreshOpenApp() {
  for (const app in appExtensions) {
    const extensions = appExtensions[app];
    for (const ext of extensions) {
      openAppByExt[ext] = app;
    }
  }
}

refreshOpenApp();

export function getOpenApp(name) {
  const ext = name.substring(name.lastIndexOf(".")).toLowerCase();
  const app = openAppByExt[ext];
  return app;
}
