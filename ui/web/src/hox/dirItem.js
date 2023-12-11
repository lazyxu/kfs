import { useState } from 'react';
import { createGlobalStore } from 'hox';

const initialState = {
    hash: "",
    name: "",
    mode: 438,
    size: 913,
    createTime: 1661133306379099400,
    modifyTime: 1661133306379099400,
    changeTime: 1661133306379099400,
    accessTime: 1661133306379099400,
    content: "console.log(\"hello, world\\n\")",
    dirItems: [],
};

const [use] = createGlobalStore(() => useState(initialState));

export default use;
