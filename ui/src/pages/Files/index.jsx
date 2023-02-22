import styles from './index.module.scss';
import useResourceManager from 'hox/resourceManager';
import {Box} from "@mui/material";
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
        <Box className={styles.right}>
            {resourceManager.branches && <Branches/>}
            {resourceManager.branchName && <DirItems/>}
        </Box>
    );
}

export default App;
