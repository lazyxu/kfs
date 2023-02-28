import {useRef} from 'react';
import AbsolutePath from "components/AbsolutePath";
import useResourceManager from 'hox/resourceManager';
import Branch from "./Branch";
import {Grid} from "@mui/material";
import BranchContextMenu from "../../components/ContextMenu/BranchContextMenu";
import BranchesContextMenu from "../../components/ContextMenu/BranchesContextMenu";
import useContextMenu from "../../hox/contextMenu";

export default function () {
    const [resourceManager, setResourceManager] = useResourceManager();
    const [contextMenu, setContextMenu] = useContextMenu();
    const branchesElm = useRef(null);

    return (
        <>
            <AbsolutePath/>
            <Grid container margin={1} spacing={1}
                  style={{flex: "auto", overflowY: "scroll"}}
                  ref={branchesElm} onContextMenu={(e) => {
                e.preventDefault();
                // console.log(e.target, e.currentTarget, e.target === e.currentTarget);
                // if (e.target === e.currentTarget) {
                const {clientX, clientY} = e;
                let {x, y, width, height} = e.currentTarget.getBoundingClientRect();
                setContextMenu({
                    type: 'branches',
                    clientX, clientY,
                    x, y, width, height,
                })
                // }
            }}>
                {resourceManager.branches.map((branch, i) => (
                    <Grid item key={branch.name}>
                        <Branch branchesElm={branchesElm} branch={branch}>{branch.name}</Branch>
                    </Grid>
                ))}
            </Grid>
            <BranchContextMenu/>
            <BranchesContextMenu/>
        </>
    );
}
