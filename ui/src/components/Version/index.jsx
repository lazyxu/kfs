import {Box, Typography} from "@mui/material";

export default () => (
    <Box sx={{
        position: 'absolute',
        bottom: "0",
        fontFamily: "KaiTi, STKaiti;",
    }}>
        <Typography>
            {process.env.REACT_APP_PLATFORM}.{process.env.NODE_ENV}
        </Typography>
    </Box>
);
