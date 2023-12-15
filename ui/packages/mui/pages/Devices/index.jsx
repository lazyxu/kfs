import { listDevice } from "@kfs/mui/api/device";
import { Box, Grid } from "@mui/material";
import { useEffect, useState } from "react";
import Device from "./Device";

export default function () {
    const [devices, setDevices] = useState([]);
    useEffect(() => {
       listDevice(setDevices);
    }, []);
    return (
        <Box sx={{ display: 'flex', flex: "1", flexDirection: 'column', minHeight: '0' }}>
            <Grid container padding={1} spacing={1}
                style={{ flex: "1", overflowY: 'auto', alignContent: "flex-start" }}>
                {devices.map((device, i) => (
                    <Grid item key={i}>
                        <Device device={device} setDevices={setDevices}/>
                    </Grid>
                ))}
            </Grid>
        </Box>
    );
}
