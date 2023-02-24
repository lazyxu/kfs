import styles from './index.module.scss';
import useResourceManager from 'hox/resourceManager';
import {Stack} from "@mui/material";
import DirItems from "./DirItems";
import Branches from "./Branches";
import React, {useEffect} from "react";
import {listBranch} from "../../api/branch";

function App() {
    const [resourceManager, setResourceManager] = useResourceManager();
    useEffect(() => {
        listBranch(setResourceManager);
    }, []);
    return (
        <Stack className={styles.right}>
            {resourceManager.branches && <Branches/>}
            {resourceManager.branchName && <DirItems/>}
        </Stack>
    );
}

export default App;
