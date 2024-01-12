import { createGlobalStore } from 'hox';
import { useState } from 'react';

function useFunc() {
  const [checked, setChecked] = useState({});
  return {
    checked,
    setChecked,
  };
}

const [useCheckedType] = createGlobalStore(useFunc);

export { useCheckedType };

const [useCheckedSuffix] = createGlobalStore(useFunc);

export { useCheckedSuffix };

