import { createGlobalStore } from 'hox';
import { useState } from 'react';

const initialState = {
    type: null,
    top: null,
    left: null,
};

const [useContextMenu] = createGlobalStore(() => useState(null));

export default useContextMenu;
