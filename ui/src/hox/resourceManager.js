import { useState } from 'react';
import { createGlobalStore } from 'hox';

const initialState = {
  branches: [],
  // branchName: "",
  filePath: [],
  dirItems: null,
  content: null,
};

const [useResourceManager] = createGlobalStore(() => useState(initialState));

export default useResourceManager;
