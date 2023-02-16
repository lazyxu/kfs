import styles from './index.module.scss';
import useResourceManager from 'hox/resourceManager';
import {Box} from "@mui/material";
import DirItems from "./DirItems";
import Branches from "./Branches";

function App() {
    const [resourceManager, setResourceManager] = useResourceManager();
    return (
        <Box className={styles.right}>
            {resourceManager.branchName ?
                <DirItems/> :
                <Branches/>
            }
        </Box>
    );
}

export default App;
