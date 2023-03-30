import {Box, Typography} from "@mui/material";
import useSysConfig from "../../hox/sysConfig";
import useWebSocket, {ReadyState} from "react-use-websocket";

export default function () {
    const {sysConfig, setSysConfig, resetSysConfig} = useSysConfig();
    const {
        sendJsonMessage,
        lastJsonMessage,
        readyState
    } = useWebSocket("ws://127.0.0.1:" + sysConfig.port + "/ws", {share: true});
    const connectionStatus = {
        [ReadyState.CONNECTING]: 'Connecting',
        [ReadyState.OPEN]: 'Open',
        [ReadyState.CLOSING]: 'Closing',
        [ReadyState.CLOSED]: 'Closed',
        [ReadyState.UNINSTANTIATED]: 'Uninstantiated',
    }[readyState];
    return (
        <Box sx={{
            position: 'absolute',
            bottom: "0",
            fontFamily: "KaiTi, STKaiti;",
        }}>
            <Typography>WebSocket: {connectionStatus}</Typography>
            <Typography>
                {process.env.REACT_APP_PLATFORM}.{process.env.NODE_ENV}
            </Typography>
        </Box>
    );
}
