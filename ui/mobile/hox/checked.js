import { createGlobalStore } from 'hox';
import { useState } from 'react';

const [useCheckedType] = createGlobalStore(()=>useState({}));

export { useCheckedType };

const [useCheckedSuffix] = createGlobalStore(()=>useState({}));

export { useCheckedSuffix };

