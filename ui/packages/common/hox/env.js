import { createGlobalStore } from 'hox';
import { useState } from 'react';

export const [useEnv, useEnv2] = createGlobalStore(() => useState({}));

export default useEnv;

export const getEnv = () => {
    return useEnv2()[0];
}
