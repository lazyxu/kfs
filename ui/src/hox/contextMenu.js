import { useState } from 'react';
import { createGlobalStore } from 'hox';

const initialState = {
    type: null,
    top: null,
    left: null,
};

const [useContextMenu] = createGlobalStore(() => useState(null));

export default useContextMenu;
