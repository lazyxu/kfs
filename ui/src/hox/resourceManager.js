import { useState } from 'react';
import { createGlobalStore } from 'hox';

const initialState = {
  branchName: 'master',
  filePath: [],
  dirItems: [],
  content: null,
};

const [useResourceManager] = createGlobalStore(() => useState(initialState));

export default useResourceManager;
