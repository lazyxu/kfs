import { createGlobalStore } from 'hox';
import { useState } from 'react';

const initialState = {
  drivers: [],
  // driverName: "",
  dirPath: [],
  filePath: [],
  dirItems: null,
  content: null,
};

export const [useResourceManager, useResourceManager2] = createGlobalStore(() => useState(initialState));

export default useResourceManager;

export const getResourceManager = () => {
  return useResourceManager2()[0];
}

export async function openDrivers(setResourceManager) {
  console.log('openDrivers');
  setResourceManager({});
}

export async function openDir(setResourceManager, dirPath) {
  console.log('openDir', dirPath);
  setResourceManager({ dirPath });
}
