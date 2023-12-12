import { Box, LinearProgress, Typography } from "@mui/material";

export default function (props) {
    return (
        <Box sx={{ width: '100%' }}>
            <Box sx={{ display: 'flex', alignItems: 'center' }}>
                <Box sx={{ width: '100%' }}>
                    <LinearProgress variant="determinate" {...props} />
                </Box>
                <Box sx={{ paddingLeft: "0.5em" }}>
                    <Typography variant="body2" color="text.secondary">{`${Math.round(
                        props.value,
                    )}%`}</Typography>
                </Box>
            </Box>
        </Box>
    );
}
