import {useRef} from 'react';
import AbsolutePath from "components/AbsolutePath";
import useResourceManager from 'hox/resourceManager';
import Driver from "./Driver";
import {Grid} from "@mui/material";
import DriverContextMenu from "../../components/ContextMenu/DriverContextMenu";
import DriversContextMenu from "../../components/ContextMenu/DriversContextMenu";
import useContextMenu from "../../hox/contextMenu";
import Dialog from "../../components/Dialog";

export default function () {
    const [resourceManager, setResourceManager] = useResourceManager();
    const [contextMenu, setContextMenu] = useContextMenu();
    const driversElm = useRef(null);

    return (
        <>
            <AbsolutePath/>
            <Grid container margin={1} spacing={1}
                  style={{flex: "auto", overflowY: "scroll"}}
                  ref={driversElm} onContextMenu={(e) => {
                e.preventDefault();
                // console.log(e.target, e.currentTarget, e.target === e.currentTarget);
                // if (e.target === e.currentTarget) {
                const {clientX, clientY} = e;
                let {x, y, width, height} = e.currentTarget.getBoundingClientRect();
                setContextMenu({
                    type: 'drivers',
                    clientX, clientY,
                    x, y, width, height,
                })
                // }
            }}>
                {resourceManager.drivers.map((driver, i) => (
                    <Grid item key={driver.name}>
                        <Driver driversElm={driversElm} driver={driver}>{driver.name}</Driver>
                    </Grid>
                ))}
            </Grid>
            <DriverContextMenu/>
            <DriversContextMenu/>
            <Dialog/>
        </>
    );
}
