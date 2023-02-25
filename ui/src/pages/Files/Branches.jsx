import {useRef} from 'react';
import AbsolutePath from "components/AbsolutePath";
import useResourceManager from 'hox/resourceManager';
import Branch from "./Branch";
import {Grid} from "@mui/material";
import BranchContextMenu from "../../components/ContextMenu/BranchContextMenu";

export default function () {
    const [resourceManager, setResourceManager] = useResourceManager();
    const branchesElm = useRef(null);

    return (
        <>
            <AbsolutePath/>
            <Grid container margin={1} spacing={1}
                  style={{flex: "auto", overflowY: "scroll"}}
                  ref={branchesElm}>
                {resourceManager.branches.map((branch, i) => (
                    <Grid item={true} key={branch.name}>
                        <Branch branchesElm={branchesElm} branch={branch}>{branch.name}</Branch>
                    </Grid>
                ))}
            </Grid>
            <BranchContextMenu/>
        </>
    );
}
