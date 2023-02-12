import { useState } from 'react';
import { createGlobalStore } from 'hox';

const initialState = {
    type: null,
};

const [useDialog] = createGlobalStore(() => useState({}));

export default useDialog;
