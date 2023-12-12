import { Add, Refresh } from "@mui/icons-material";
import { Grid, ListItemText, MenuItem } from "@mui/material";
import { listDriver } from "api/driver";
import Menu from 'components/Menu';
import { useEffect, useState } from 'react';
import DriverAttribute from "./Attribute/DriverAttribute";
import Driver from "./Driver";
import NewDriver from "./NewDriver/NewDriver";

export default function () {
    let [drivers, setDrivers] = useState([]);
    let [openNewDrive, setOpenNewDrive] = useState(false);
    let [driverAttribute, setDriverAttribute] = useState();
    const [contextMenu, setContextMenu] = useState(null);
    useEffect(() => {
        listDriver().then(setDrivers);
    }, []);
    return (
        <>
            <Grid container padding={1} spacing={1}
                style={{ flex: "1", overflowY: 'auto', alignContent: "flex-start" }}
                onContextMenu={(e) => {
                    e.preventDefault(); e.stopPropagation();
                    setContextMenu({ mouseX: e.clientX, mouseY: e.clientY });
                }}
            >
                {drivers.map(driver => (
                    <Grid item key={driver.name}>
                        <Driver driver={driver} setDriverAttribute={setDriverAttribute} onDelete={() => listDriver().then(setDrivers)}>{driver.name}</Driver>
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
                <MenuItem onClick={() => { listDriver().then(drivers => { setContextMenu(null); setDrivers(drivers); }) }}>
                    <Refresh />
                    <ListItemText>刷新</ListItemText>
                </MenuItem>
            </Menu>
            {openNewDrive && <NewDriver onClose={() => setOpenNewDrive(false)} onSucc={() => { setOpenNewDrive(false); listDriver().then(setDrivers); }} />}
            {driverAttribute && <DriverAttribute onClose={() => setDriverAttribute(false)} driver={driverAttribute} />}
        </>
    );
}
