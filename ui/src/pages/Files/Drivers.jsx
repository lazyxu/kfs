import { Add, Refresh } from "@mui/icons-material";
import { Grid, ListItemText, MenuItem } from "@mui/material";
import { listDriver } from "api/driver";
import Menu from 'components/Menu';
import useResourceManager from 'hox/resourceManager';
import { useState } from 'react';
import Dialog from "../../components/Dialog";
import Driver from "./Driver";
import DriverAttribute from "./DriverAttribute";
import NewDriver from "./NewDriver";

export default function () {
    const [resourceManager, setResourceManager] = useResourceManager();
    let [openNewDrive, setOpenNewDrive] = useState(false);
    let [driverAttribute, setDriverAttribute] = useState();
    const [contextMenu, setContextMenu] = useState(null);
    return (
        <>
            <Grid container padding={1} spacing={1}
                style={{ flex: "1", overflowY: 'auto', alignContent: "flex-start" }}
                onContextMenu={(e) => {
                    e.preventDefault(); e.stopPropagation();
                    setContextMenu({ mouseX: e.clientX, mouseY: e.clientY });
                }}
            >
                {resourceManager.drivers.map((driver, i) => (
                    <Grid item key={driver.name}>
                        <Driver driver={driver} setDriverAttribute={setDriverAttribute}>{driver.name}</Driver>
                    </Grid>
                ))}
            </Grid>
            <Menu
                contextMenu={contextMenu}
                open={contextMenu !== null}
                onClose={() => setContextMenu(null)}
            >
                <MenuItem onClick={() => { setContextMenu(null); setOpenNewDrive(true) }}>
                    <Add />
                    <ListItemText>新建云盘</ListItemText>
                </MenuItem>
                <MenuItem onClick={() => { setContextMenu(null); listDriver(setResourceManager) }}>
                    <Refresh />
                    <ListItemText>刷新</ListItemText>
                </MenuItem>
            </Menu>
            {openNewDrive && <NewDriver setOpen={setOpenNewDrive} />}
            {driverAttribute && <DriverAttribute setOpen={setDriverAttribute} driver={driverAttribute} />}
            <Dialog />
        </>
    );
}
