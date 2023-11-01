import { Add, Refresh } from "@mui/icons-material";
import { Grid, SpeedDial, SpeedDialAction, SpeedDialIcon } from "@mui/material";
import { listDriver } from "api/driver";
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

    const actions = [
        { icon: <Add />, name: '新建云盘', onClick: () => setOpenNewDrive(true) },
        { icon: <Refresh />, name: '刷新', onClick: () => listDriver(setResourceManager) },
    ];
    return (
        <>
            <Grid container padding={1} spacing={1}
                style={{ flex: "1", overflowY: 'auto', alignContent: "flex-start" }}>
                {resourceManager.drivers.map((driver, i) => (
                    <Grid item key={driver.name}>
                        <Driver driver={driver} setDriverAttribute={setDriverAttribute}>{driver.name}</Driver>
                    </Grid>
                ))}
            </Grid>
            <SpeedDial
                ariaLabel="SpeedDial basic example"
                sx={{ position: 'absolute', bottom: 16, right: 16 }}
                icon={<SpeedDialIcon />}
            >
                {actions.map((action) => (
                    <SpeedDialAction
                        key={action.name}
                        icon={action.icon}
                        tooltipTitle={action.name}
                        onClick={action.onClick}
                    />
                ))}
            </SpeedDial>
            {openNewDrive && <NewDriver setOpen={setOpenNewDrive} />}
            {driverAttribute && <DriverAttribute setOpen={setDriverAttribute} driver={driverAttribute} />}
            <Dialog />
        </>
    );
}
